/*
EClient is the main struct to use from API user's point of view.
It takes care of almost everything:
  - implementing the requests
  - creating the answer decoder
  - creating the connection to TWS/IBGW

The user just needs to override EWrapper methods to receive the answers.
*/
package ibapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/scmhub/ibapi/protobuf"
	"google.golang.org/protobuf/proto"
)

type ConnState int

const (
	DISCONNECTED ConnState = iota
	CONNECTING
	CONNECTED
)

func (cs ConnState) String() string {
	switch cs {
	case DISCONNECTED:
		return "disconnected"
	case CONNECTING:
		return "connecting"
	case CONNECTED:
		return "connected"
	default:
		return "unknown connection state"
	}
}

// MsgEncoder efficiently encodes messages for IB API
type MsgEncoder struct {
	buf     bytes.Buffer
	eClient *EClient
}

// NewMsgEncoder creates a new MsgEncoder with an initial number of fields and server version
func NewMsgEncoder(nFields int, eClient *EClient) *MsgEncoder {
	me := &MsgEncoder{
		eClient: eClient,
	}
	me.buf.Grow(8*nFields + 4)

	// Reserve 4 bytes for the message size header
	me.buf.Write([]byte{0, 0, 0, 0})
	return me
}

// encodeMsgID encodes a message ID with appropriate format based on server version
func (me *MsgEncoder) encodeMsgID(msgID int64) *MsgEncoder {
	if me.eClient.serverVersion >= MIN_SERVER_VER_PROTOBUF {
		// Encode as raw int (4 bytes, byte-swapped)
		me.encodeRawInt64(msgID)
	} else {
		// Encode as a regular field (string + delimiter)
		me.encodeInt64(msgID)
	}
	return me
}

// encodeField encodes a value of any supported type
func (me *MsgEncoder) encodeField(v any) *MsgEncoder {
	switch val := v.(type) {
	case int:
		return me.encodeInt(val)
	case int64:
		return me.encodeInt64(val)
	case float64:
		return me.encodeFloat64(val)
	case string:
		return me.encodeString(val)
	case bool:
		return me.encodeBool(val)
	case []byte:
		return me.encodeBytes(val)
	case Decimal:
		return me.encodeDecimal(val)
	default:
		// Convert to string as fallback
		return me.encodeString(fmt.Sprintf("%v", val))
	}
}

// encodeFileds encode many fields
func (me *MsgEncoder) encodeFields(v ...any) *MsgEncoder {
	for _, f := range v {
		me.encodeField(f)
	}
	return me
}

// encodeMax encodes a value that might be UNSET
func (me *MsgEncoder) encodeMax(v any) *MsgEncoder {
	switch val := v.(type) {
	case int64:
		return me.encodeIntMax(val)
	case float64:
		return me.encodeFloatMax(val)
	// case Decimal:
	// 	return me.encodeDecimalMax(val)
	default:
		return me.encodeField(v)
	}
}

// encodeInt adds an int value to the message
func (me *MsgEncoder) encodeInt(v int) *MsgEncoder {
	me.buf.WriteString(strconv.Itoa(v))
	me.buf.WriteByte(delim)
	return me
}

// encodeInt64 adds an int64 value to the message
func (me *MsgEncoder) encodeInt64(v int64) *MsgEncoder {
	me.buf.WriteString(strconv.FormatInt(v, 10))
	me.buf.WriteByte(delim)
	return me
}

// encodeInt64 adds a raw int64 value to the message
func (me *MsgEncoder) encodeRawInt64(v int64) *MsgEncoder {
	var arrayOfBytes [RAW_INT_LEN]byte
	binary.BigEndian.PutUint32(arrayOfBytes[:], uint32(v))
	me.buf.Write(arrayOfBytes[:])
	return me
}

// encodeFloat64 adds a float64 value to the message
func (me *MsgEncoder) encodeFloat64(v float64) *MsgEncoder {
	me.buf.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	me.buf.WriteByte(delim)
	return me
}

// encodeString adds a string value to the message
func (me *MsgEncoder) encodeString(v string) *MsgEncoder {
	me.buf.WriteString(v)
	me.buf.WriteByte(delim)
	return me
}

// encodeBool adds a boolean value to the message
func (me *MsgEncoder) encodeBool(v bool) *MsgEncoder {
	if v {
		me.buf.WriteByte('1')
	} else {
		me.buf.WriteByte('0')
	}
	me.buf.WriteByte(delim)
	return me
}

// encodeBytes adds raw bytes to the message
func (me *MsgEncoder) encodeBytes(v []byte) *MsgEncoder {
	me.buf.Write(v)
	me.buf.WriteByte(delim)
	return me
}

// encodeBytes adds raw bytes to the message
func (me *MsgEncoder) encodeProto(v []byte) *MsgEncoder {
	me.buf.Write(v)
	return me
}

// encodeDecimal adds a Decimal value to the message
func (me *MsgEncoder) encodeDecimal(v Decimal) *MsgEncoder {
	me.buf.WriteString(DecimalToString(v))
	me.buf.WriteByte(delim)
	return me
}

// encodeTagValues adds a slice of TagValue to the message
func (me *MsgEncoder) encodeTagValues(v []TagValue) *MsgEncoder {
	for _, tv := range v {
		me.buf.WriteString(tv.Tag)
		me.buf.WriteString("=")
		me.buf.WriteString(tv.Value)
		me.buf.WriteString(";")
	}
	me.buf.WriteByte(delim)
	return me
}

func (me *MsgEncoder) encodeContract(v *Contract) *MsgEncoder {
	me.encodeInt64(v.ConID)
	me.encodeString(v.Symbol)
	me.encodeString(v.SecType)
	me.encodeString(v.LastTradeDateOrContractMonth)
	me.encodeFloatMax(v.Strike)
	me.encodeString(v.Right)
	me.encodeString(v.Multiplier)
	me.encodeString(v.Exchange)
	me.encodeString(v.PrimaryExchange)
	me.encodeString(v.Currency)
	me.encodeString(v.LocalSymbol)
	me.encodeString(v.TradingClass)
	me.encodeBool(v.IncludeExpired)
	return me
}

// encodeIntMax adds an int64 value to the message, handling UNSET_INT
func (me *MsgEncoder) encodeIntMax(v int64) *MsgEncoder {
	if v == UNSET_INT {
		me.buf.WriteByte(delim)
		return me
	}
	return me.encodeInt64(v)
}

// encodeFloatMax adds a float64 value to the message, handling UNSET_FLOAT
func (me *MsgEncoder) encodeFloatMax(v float64) *MsgEncoder {
	if v == UNSET_FLOAT {
		me.buf.WriteByte(delim)
		return me
	}
	return me.encodeFloat64(v)
}

// // encodeDecimalMax adds a Decimal value to the message, handling UNSET_DECIMAL
// func (me *MsgEncoder) encodeDecimalMax(v Decimal) *MsgEncoder {
// 	if v == UNSET_DECIMAL {
// 		me.buf.WriteByte(delim)
// 		return me
// 	}
// 	return me.encodeDecimal(v)
// }

// Bytes finalizes the message by writing the size header and returning the complete message
func (me *MsgEncoder) Bytes() []byte {
	// Get the final buffer bytes
	result := me.buf.Bytes()

	// Calculate message size (excluding the 4-byte header)
	msgSize := len(result) - 4

	if msgSize > MAX_MSG_LEN {
		log.Error().Int("msgSize", msgSize).Msg("Message size exceeds maximum allowed size")
		me.eClient.wrapper.Error(NO_VALID_ID, currentTimeMillis(), BAD_LENGTH.Code, BAD_LENGTH.Msg, "")
		return nil
	}

	// Write the size back into the header
	binary.BigEndian.PutUint32(result[:4], uint32(msgSize))

	return result
}

// Reset resets the buffer for reuse while maintaining its capacity
func (me *MsgEncoder) Reset() {
	me.buf.Reset()
	// Reserve 4 bytes for the message size header
	me.buf.Write([]byte{0, 0, 0, 0})
}

// EClient is the main struct to use from API user's point of view.
type EClient struct {
	host                 string
	port                 int
	clientID             int64
	connectOptions       string
	optionalCapabilities string
	conn                 *Connection
	serverVersion        Version
	connTime             string
	connState            ConnState
	writer               *bufio.Writer
	scanner              *bufio.Scanner
	wrapper              EWrapper
	decoder              *EDecoder
	reqChan              chan []byte
	ctx                  context.Context
	cancel               context.CancelFunc
	extraAuth            bool
	wg                   sync.WaitGroup
	watchOnce            sync.Once
	err                  error
}

// NewEClient returns a new Eclient.
func NewEClient(wrapper EWrapper) *EClient {
	if wrapper == nil {
		wrapper = &Wrapper{}
	}
	c := &EClient{wrapper: wrapper}
	c.reset()

	return c
}

func (c *EClient) reset() {

	c.host = ""
	c.port = -1
	c.clientID = -1
	c.connectOptions = ""
	c.optionalCapabilities = ""
	c.extraAuth = false
	c.conn = &Connection{wrapper: c.wrapper}
	c.serverVersion = -1
	c.connTime = ""

	// writer
	c.writer = bufio.NewWriter(c.conn)
	// init scanner
	c.scanner = bufio.NewScanner(c.conn)
	c.scanner.Split(scanFields)
	c.scanner.Buffer(make([]byte, 4096), MAX_MSG_LEN)

	c.reqChan = make(chan []byte, 10)

	c.ctx, c.cancel = context.WithCancel(context.Background())

	c.wg = sync.WaitGroup{}
	c.err = nil

	c.watchOnce = sync.Once{}

	c.setConnState(DISCONNECTED)
	c.connectOptions = ""
}

func (c *EClient) setConnState(state ConnState) {
	cs := ConnState(atomic.LoadInt32((*int32)(unsafe.Pointer(&c.connState))))
	atomic.StoreInt32((*int32)(unsafe.Pointer(&c.connState)), int32(state))
	log.Debug().Stringer("from", cs).Stringer("to", state).Msg("connection state changed")
}

// request is a goroutine that will get the req from reqChan and send it to TWS.
func (c *EClient) request() {
	log.Debug().Msg("requester started")
	defer log.Debug().Msg("requester ended")

	c.wg.Add(1)
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		case req := <-c.reqChan:
			log.Trace().Bytes("req", req).Msg("sending request")
			if !c.IsConnected() {
				c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
				c.cancel()
				return
			}
			nn, err := c.writer.Write(req)
			if err != nil {
				log.Error().Err(err).Int("nbytes", nn).Bytes("reqMsg", req).Msg("requester write error")
				// Disconnect the client
				log.Info().Msg("Disconnecting client due to write error.")
				if disconnectErr := c.Disconnect(); disconnectErr != nil {
					log.Error().Err(disconnectErr).Msg("Error during disconnect.")
				}
				c.cancel()
				return
			}
			err = c.writer.Flush()
			if err != nil {
				log.Error().Err(err).Bytes("reqMsg", req).Msg("requester flush error")
				c.cancel()
				return
			}
		}
	}
}

func (c *EClient) validateInvalidSymbols(host string) error {
	if host != "" && !isASCIIPrintable(host) {
		return errors.New(host)
	}
	if c.connectOptions != "" && !isASCIIPrintable(c.connectOptions) {
		return errors.New(c.connectOptions)
	}
	if c.optionalCapabilities != "" && !isASCIIPrintable(c.optionalCapabilities) {
		return errors.New(c.optionalCapabilities)
	}
	return nil
}

func (c *EClient) useProtoBuf(msgID int64) bool {
	if version, exists := PROTOBUF_MSG_IDS[OUT(msgID)]; exists {
		return version <= c.serverVersion
	}
	return false
}

// startAPI initiates the message exchange between the client application and the TWS/IB Gateway.
func (c *EClient) startAPI() error {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return NOT_CONNECTED
	}

	const VERSION = 2

	msg := makeField(VERSION) + makeField(c.clientID)

	if c.serverVersion >= MIN_SERVER_VER_OPTIONAL_CAPABILITIES {
		msg += makeField(c.optionalCapabilities)
	}
	var payload []byte
	if c.serverVersion >= MIN_SERVER_VER_PROTOBUF {
		idBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(idBytes, uint32(START_API))
		payload = append(idBytes, []byte(msg)...)
	} else {
		payload = []byte(makeField(START_API) + msg)

	}

	msgLen := uint32(len(payload))
	bs := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(bs[:4], msgLen)
	copy(bs[4:], payload)

	log.Debug().Bytes("req", bs).Msg("sending startAPI")
	if _, err := c.writer.Write(bs); err != nil {
		return err
	}
	if err := c.writer.Flush(); err != nil {
		return err
	}

	return nil
}

// Connect must be called before any other.
// There is no feedback for a successful connection, but a subsequent attempt to connect will return the message "Already connected.".
// You should wait for the connection to be established and NextValidID to be returned before calling any other function. If you don't wait, you will get a broken pipe error.
func (c *EClient) Connect(host string, port int, clientID int64) error {

	if c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), ALREADY_CONNECTED.Code, ALREADY_CONNECTED.Msg, "")
		return NOT_CONNECTED
	}

	if err := c.validateInvalidSymbols(host); err != nil {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), INVALID_SYMBOL.Code, INVALID_SYMBOL.Msg+err.Error(), "")
		return err
	}

	c.host = host
	c.port = port
	c.clientID = clientID

	c.setConnState(CONNECTING)

	// Connecting to IB server
	log.Info().Str("host", host).Int("port", port).Int64("clientID", clientID).Msg("Connecting to IB server")
	if err := c.conn.connect(c.host, c.port); err != nil {
		log.Error().Err(CONNECT_FAIL).Msg("Connection fail")
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), CONNECT_FAIL.Code, CONNECT_FAIL.Msg, "")
		c.reset()
		return CONNECT_FAIL
	}

	// HandShake with the TWS or GateWay to ensure the version,
	log.Debug().Msg("Handshake with TWS or GateWay")

	head := []byte("API\x00")

	connectOptions := ""
	if c.connectOptions != "" {
		connectOptions = " " + c.connectOptions
	}
	sizeofCV := make([]byte, 4)
	clientVersion := fmt.Appendf(nil, "v%d..%d%s", MIN_CLIENT_VER, MAX_CLIENT_VER, connectOptions)

	binary.BigEndian.PutUint32(sizeofCV, uint32(len(clientVersion)))

	var msg bytes.Buffer
	msg.Write(head)
	msg.Write(sizeofCV)
	msg.Write(clientVersion)

	log.Debug().Bytes("header", msg.Bytes()).Msg("Sending handshake header")

	if _, err := c.writer.Write(msg.Bytes()); err != nil {
		return err
	}

	if err := c.writer.Flush(); err != nil {
		return err
	}

	log.Debug().Msg("Receiving handshake Info")

	// scan once to get server info
	if !c.scanner.Scan() {
		return c.scanner.Err()
	}

	// Init server info
	msgBytes := c.scanner.Bytes()
	serverInfo := splitMsgBytes(msgBytes)
	v, _ := strconv.Atoi(string(serverInfo[0]))
	c.serverVersion = Version(v)
	if c.serverVersion < MIN_SERVER_VER_SUPPORTED {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UNSUPPORTED_VERSION.Code, UNSUPPORTED_VERSION.Msg, "")
		return UNSUPPORTED_VERSION
	}

	c.connTime = string(serverInfo[1])
	log.Info().Int("serverVersion", v).Str("connectionTime", c.connTime).Msg("Handshake completed")

	// init decoder
	c.decoder = &EDecoder{wrapper: c.wrapper, serverVersion: c.serverVersion}

	//start Ereader
	go EReader(c.ctx, c.cancel, c.scanner, c.decoder, &c.wg)

	// start requester
	go c.request()

	c.setConnState(CONNECTED)
	c.wrapper.ConnectAck()

	// startAPI
	if err := c.startAPI(); err != nil {
		return err
	}

	// 4) Launch the shutdown watcher exactly once
	c.watchOnce.Do(func() {
		go func() {
			<-c.ctx.Done() // waits for c.cancel()
			if err := c.Disconnect(); err != nil {
				log.Error().Err(err).Msg("Disconnect error in watcher")
			}
		}()
	})

	log.Debug().Msg("IB Client Connected!")

	return nil
}

// ConnectWithGracefulShutdown connects and sets up signal handling for graceful shutdown.
// This is a convenience for simple apps. Advanced users should handle signals themselves.
func (c *EClient) ConnectWithGracefulShutdown(host string, port int, clientID int64) error {
	err := c.Connect(host, port, clientID)
	if err != nil {
		return err
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Warn().Msg("detected termination signal, shutting down gracefully")
		c.Disconnect()
		os.Exit(0)
	}()
	return nil
}

// Disconnect terminates the connections with TWS.
// Calling this function does not cancel orders that have already been sent.
func (c *EClient) Disconnect() error {
	if !c.IsConnected() {
		return nil
	}

	// Set Disconnected state realy so that new calls to Disconnect() will not block
	c.setConnState(DISCONNECTED)

	// 1) Cancel to unblock request Loop
	c.cancel()

	// 2) Close the socket to unblock reader loop
	if err := c.conn.disconnect(); err != nil {
		return err
	}

	// 3) Wait for loops to exit
	c.wg.Wait()

	// 4) Reset internal state
	c.reset()

	c.wrapper.ConnectionClosed()
	log.Debug().Msg("IB Client Disconnected!")

	return nil
}

func (c *EClient) Ctx() context.Context {
	return c.ctx
}

// IsConnected checks connection to TWS or GateWay.
func (c *EClient) IsConnected() bool {
	return c.conn.IsConnected() && ConnState(atomic.LoadInt32((*int32)(unsafe.Pointer(&c.connState)))) == CONNECTED
}

// OptionalCapabilities returns the Optional Capabilities.
func (c *EClient) OptionalCapabilities() string {
	return c.optionalCapabilities
}

// SetOptionalCapabilities setup the Optional Capabilities.
func (c *EClient) SetOptionalCapabilities(optCapts string) {
	c.optionalCapabilities = optCapts
}

// SetConnectionOptions setup the Connection Options.
func (c *EClient) SetConnectionOptions(connectOptions string) {

	if c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), ALREADY_CONNECTED.Code, ALREADY_CONNECTED.Msg, "")
		return
	}

	c.connectOptions = connectOptions
}

// ReqCurrentTime asks the current system time on the server side.
func (c *EClient) ReqCurrentTime() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_CURRENT_TIME).encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

// ServerVersion returns the version of the TWS instance to which the API application is connected.
func (c *EClient) ServerVersion() Version {
	return c.serverVersion
}

// SetServerLogLevel sets the log level of the server.
// logLevel can be:
// 1 = SYSTEM
// 2 = ERROR	(default)
// 3 = WARNING
// 4 = INFORMATION
// 5 = DETAIL
func (c *EClient) SetServerLogLevel(logLevel int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(SET_SERVER_LOGLEVEL).encodeInt(VERSION).encodeInt64(logLevel)

	c.reqChan <- me.Bytes()
}

// ConnectionTime is the time the API application made a connection to TWS.
func (c *EClient) TWSConnectionTime() string {
	return c.connTime
}

//	##########################################################################
//	#		Market Data
// 	##########################################################################

// ReqMktData Call this function to request market data.
// The market data will be returned by the tickPrice and tickSize events.
// reqID, the ticker id must be a unique value. When the market data returns it will be identified by this tag. This is also used when canceling the market data.
// contract contains a description of the Contract for which market data is being requested.
// genericTickList is a commma delimited list of generic tick types. Tick types can be found in the Generic Tick Types page.
// Prefixing w/ 'mdoff' indicates that top mkt data shouldn't tick. You can specify the news source by postfixing w/ ':<source>. Example: "mdoff,292:FLY+BRF"
// snapshot checks to return a single snapshot of Market data and have the market data subscription cancel.
// Do not enter any genericTicklist values if you use snapshots.
// regulatorySnapshot: With the US Value Snapshot Bundle for stocks, regulatory snapshots are available for 0.01 USD each.
// mktDataOptions is for internal use only.Use default value XYZ.
func (c *EClient) ReqMktData(reqID TickerID, contract *Contract, genericTickList string, snapshot bool, regulatorySnapshot bool, mktDataOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_DELTA_NEUTRAL && contract.DeltaNeutralContract != nil {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support delta-neutral orders.", "")
		return
	}
	if c.serverVersion < MIN_SERVER_VER_REQ_MKT_DATA_CONID && contract.ConID > 0 {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support conId parameter.", "")
		return
	}
	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tradingClass parameter in reqMktData.", "")
		return
	}

	const VERSION = 11

	me := NewMsgEncoder(30, c)

	me.encodeMsgID(REQ_MKT_DATA).encodeInt(VERSION).encodeInt64(reqID)

	if c.serverVersion >= MIN_SERVER_VER_REQ_MKT_DATA_CONID {
		me.encodeInt64(contract.ConID)
	}

	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier) // srv v15 and above
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	// Send combo legs for BAG requests (srv v8 and above)
	if contract.SecType == "BAG" {
		comboLegsCount := len(contract.ComboLegs)
		me.encodeInt(comboLegsCount)
		for _, comboLeg := range contract.ComboLegs {
			me.encodeInt64(comboLeg.ConID)
			me.encodeInt64(comboLeg.Ratio)
			me.encodeString(comboLeg.Action)
			me.encodeString(comboLeg.Exchange)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL {
		if contract.DeltaNeutralContract != nil {
			me.encodeBool(true)
			me.encodeInt64(contract.DeltaNeutralContract.ConID)
			me.encodeFloat64(contract.DeltaNeutralContract.Delta)
			me.encodeFloat64(contract.DeltaNeutralContract.Price)
		} else {
			me.encodeBool(false)
		}
	}

	me.encodeString(genericTickList)
	me.encodeBool(snapshot)

	if c.serverVersion >= MIN_SERVER_VER_REQ_SMART_COMPONENTS {
		me.encodeBool(regulatorySnapshot)
	}

	// send mktDataOptions parameter
	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(mktDataOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelMktData stops the market data flow for the specified TickerId.
func (c *EClient) CancelMktData(reqID TickerID) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 2

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_MKT_DATA).encodeInt(VERSION).encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqMarketDataType changes the market data type.
//
// The API can receive frozen market data from Trader Workstation. Frozen market data is the last data recorded in our system.
// During normal trading hours, the API receives real-time market data.
// If you use this function, you are telling TWS to automatically switch to frozen market data after the close. Then, before the opening of the next
// trading day, market data will automatically switch back to real-time market data.
// marketDataType:
//
//	1 -> realtime streaming market data
//	2 -> frozen market data
//	3 -> delayed market data
//	4 -> delayed frozen market data
func (c *EClient) ReqMarketDataType(marketDataType int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_MARKET_DATA_TYPE {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support market data type requests.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(REQ_MARKET_DATA_TYPE).encodeInt(VERSION).encodeInt64(marketDataType)

	c.reqChan <- me.Bytes()
}

// ReqSmartComponents request the smartComponents.
func (c *EClient) ReqSmartComponents(reqID int64, bboExchange string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_SMART_COMPONENTS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support smart components request.", "")
		return
	}

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(REQ_SMART_COMPONENTS).encodeInt64(reqID).encodeString(bboExchange)

	c.reqChan <- me.Bytes()
}

// ReqMarketRule requests the market rule.
func (c *EClient) ReqMarketRule(marketRuleID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MARKET_RULES {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support market rule requests.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_MARKET_RULE).encodeInt64(marketRuleID)

	c.reqChan <- me.Bytes()
}

// ReqTickByTickData request the tick-by-tick data.
// tickType is "Last", "AllLast", "BidAsk" or "MidPoint".
// numberOfTicks is the number of ticks or 0 for unlimited.
// ignoreSize will ignore bid/ask ticks that only update the size if true.
// Result will be delivered via wrapper.TickByTickAllLast() wrapper.TickByTickBidAsk() wrapper.TickByTickMidPoint().
func (c *EClient) ReqTickByTickData(reqID int64, contract *Contract, tickType string, numberOfTicks int64, ignoreSize bool) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TICK_BY_TICK {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tick-by-tick data requests.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TICK_BY_TICK_IGNORE_SIZE {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support ignoreSize and numberOfTicks parameters in tick-by-tick data requests.", "")
		return
	}

	me := NewMsgEncoder(17, c)

	me.encodeMsgID(REQ_TICK_BY_TICK_DATA)

	me.encodeInt64(reqID)
	me.encodeInt64(contract.ConID)
	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier)
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)
	me.encodeString(contract.TradingClass)
	me.encodeString(tickType)

	if c.serverVersion >= MIN_SERVER_VER_TICK_BY_TICK_IGNORE_SIZE {
		me.encodeInt64(numberOfTicks).encodeBool(ignoreSize)
	}

	c.reqChan <- me.Bytes()
}

// CancelTickByTickData cancel the tick-by-tick data
func (c *EClient) CancelTickByTickData(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TICK_BY_TICK {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tick-by-tick data requests.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_TICK_BY_TICK_DATA).encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Options
// 	##########################################################################

// CalculateImpliedVolatility calculates the implied volatility of the option.
// Result will be delivered via wrapper.TickOptionComputation().
func (c *EClient) CalculateImpliedVolatility(reqID int64, contract *Contract, optionPrice float64, underPrice float64, impVolOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_CALC_IMPLIED_VOLAT {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support calculateImpliedVolatility req.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tradingClass parameter in calculateImpliedVolatility.", "")
		return
	}

	const VERSION = 3

	me := NewMsgEncoder(19, c)

	me.encodeMsgID(REQ_CALC_IMPLIED_VOLAT)

	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	me.encodeInt64(contract.ConID)
	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier)
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	me.encodeFloat64(optionPrice)
	me.encodeFloat64(underPrice)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(impVolOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelCalculateImpliedVolatility cancels a request to calculate volatility for a supplied option price and underlying price.
func (c *EClient) CancelCalculateImpliedVolatility(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_CALC_IMPLIED_VOLAT {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support calculateImpliedVolatility req.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_CALC_IMPLIED_VOLAT)

	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// CalculateOptionPrice calculate the price of the option
// Call this function to calculate price for a supplied option volatility and underlying price.
// Result will be delivered via wrapper.TickOptionComputation().
func (c *EClient) CalculateOptionPrice(reqID int64, contract *Contract, volatility float64, underPrice float64, optPrcOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_CALC_IMPLIED_VOLAT {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support calculateImpliedVolatility req.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tradingClass parameter in calculateImpliedVolatility.", "")
		return
	}

	const VERSION = 2

	me := NewMsgEncoder(19, c)

	me.encodeMsgID(REQ_CALC_OPTION_PRICE)

	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	me.encodeInt64(contract.ConID)
	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier)
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	me.encodeFloat64(volatility)
	me.encodeFloat64(underPrice)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(optPrcOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelCalculateOptionPrice cancels the calculation of option price.
func (c *EClient) CancelCalculateOptionPrice(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_CALC_IMPLIED_VOLAT {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support calculateImpliedVolatility req.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_CALC_OPTION_PRICE)

	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ExerciseOptions exercises the option defined by the contract.
// reqId is the ticker id and must be a unique value.
// contract contains a description of the contract to be exercised.
// exerciseAction specifies whether you want the option to lapse or be exercised.
//
//	Values: 1 = exercise, 2 = lapse.
//
// exerciseQuantity is the quantity you want to exercise.
// account is the destination account.
// override	specifies whether your setting will override the system's natural action.
// For example, if your action is "exercise" and the option is not in-the-money, by natural action the option would not exercise.
// If you have override set to "yes" the natural action would be overridden	and the out-of-the money option would be exercised.
// Values: 0 = no, 1 = yes.
// manualOrderTime isthe manual order time.
// customerAccount is the customer account.
// professionalCustomer:bool - professional customer.
func (c *EClient) ExerciseOptions(reqID TickerID, contract *Contract, exerciseAction int, exerciseQuantity int, account string, override int, manualOrderTime string, customerAccount string, professionalCustomer bool) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS && (contract.TradingClass != "" || contract.ConID > 0) {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support conId, multiplier, tradingClass parameter in exerciseOptions.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MANUAL_ORDER_TIME_EXERCISE_OPTIONS && manualOrderTime != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support manual order time parameter in exerciseOptions.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CUSTOMER_ACCOUNT && customerAccount != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support customer account parameter in exerciseOptions.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PROFESSIONAL_CUSTOMER && professionalCustomer {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support professional customer parameter in exerciseOptions.", "")
		return
	}

	const VERSION = 2

	me := NewMsgEncoder(17, c)

	me.encodeMsgID(EXERCISE_OPTIONS)

	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	// send contract fields
	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeInt64(contract.ConID)
	}

	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier)
	me.encodeString(contract.Exchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	me.encodeInt(exerciseAction)
	me.encodeInt(exerciseQuantity)
	me.encodeString(account)
	me.encodeInt(override)

	if c.serverVersion >= MIN_SERVER_VER_MANUAL_ORDER_TIME_EXERCISE_OPTIONS {
		me.encodeString(manualOrderTime)
	}

	if c.serverVersion >= MIN_SERVER_VER_CUSTOMER_ACCOUNT {
		me.encodeString(customerAccount)
	}

	if c.serverVersion >= MIN_SERVER_VER_PROFESSIONAL_CUSTOMER {
		me.encodeBool(professionalCustomer)
	}

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Orders
// 	##########################################################################

// PlaceOrder places an order.
// The order status will be returned by the orderStatus event.
// The order id must specify a unique value. When the order status returns, it will be identified by this tag.
// This tag is also used when canceling the order.
// contract contains a description of the contract which is being traded.
// order contains the details of the traded order.
func (c *EClient) PlaceOrder(orderID OrderID, contract *Contract, order *Order) {

	if c.useProtoBuf(PLACE_ORDER) {
		placeOrderRequestProto, err := createPlaceOrderRequestProto(orderID, contract, order)
		if err != nil {
			c.wrapper.Error(orderID, currentTimeMillis(), ERROR_ENCODING_PROTOBUF.Code, ERROR_ENCODING_PROTOBUF.Msg+err.Error(), "")
		}
		c.placeOrderProtoBuf(placeOrderRequestProto)
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(orderID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_DELTA_NEUTRAL && contract.DeltaNeutralContract != nil {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support delta-neutral orders.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SCALE_ORDERS2 && order.ScaleSubsLevelSize != UNSET_INT {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support Subsequent Level Size for Scale orders.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_ALGO_ORDERS && order.AlgoStrategy != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support algo orders.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_NOT_HELD && order.NotHeld {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support notHeld parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SEC_ID_TYPE && (contract.SecType != "" || contract.SecID != "") {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support secIdType and secId parameters.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PLACE_ORDER_CONID && contract.ConID != UNSET_INT && contract.ConID > 0 {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support conId parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SSHORTX && order.ExemptCode != -1 {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support exemptCode parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SSHORTX {
		for _, comboLeg := range contract.ComboLegs {
			if comboLeg.ExemptCode != -1 {
				c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support exemptCode parameter.", "")
				return
			}
		}
	}

	if c.serverVersion < MIN_SERVER_VER_HEDGE_ORDERS && order.HedgeType != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support hedge orders.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_OPT_OUT_SMART_ROUTING && order.OptOutSmartRouting {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support optOutSmartRouting parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_DELTA_NEUTRAL_CONID {
		if order.DeltaNeutralConID > 0 || order.DeltaNeutralSettlingFirm != "" || order.DeltaNeutralClearingAccount != "" || order.DeltaNeutralClearingIntent != "" {
			c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support deltaNeutral parameters: ConId, SettlingFirm, ClearingAccount, ClearingIntent.", "")
			return
		}
	}

	if c.serverVersion < MIN_SERVER_VER_DELTA_NEUTRAL_OPEN_CLOSE {
		if order.DeltaNeutralOpenClose != "" ||
			order.DeltaNeutralShortSale ||
			order.DeltaNeutralShortSaleSlot > 0 ||
			order.DeltaNeutralDesignatedLocation != "" {
			c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support deltaNeutral parameters: OpenClose, ShortSale, ShortSaleSlot, DesignatedLocation.", "")
			return
		}
	}

	if c.serverVersion < MIN_SERVER_VER_SCALE_ORDERS3 {
		if (order.ScalePriceIncrement > 0 && order.ScalePriceIncrement != UNSET_FLOAT) &&
			(order.ScalePriceAdjustValue != UNSET_FLOAT ||
				order.ScalePriceAdjustInterval != UNSET_INT ||
				order.ScaleProfitOffset != UNSET_FLOAT ||
				order.ScaleAutoReset ||
				order.ScaleInitPosition != UNSET_INT ||
				order.ScaleInitFillQty != UNSET_INT ||
				order.ScaleRandomPercent) {
			c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+
				" It does not support Scale order parameters: PriceAdjustValue, PriceAdjustInterval, "+
				"ProfitOffset, AutoReset, InitPosition, InitFillQty and RandomPercent.", "")
			return
		}
	}

	if c.serverVersion < MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE && contract.SecType == "BAG" {
		for _, orderComboLeg := range order.OrderComboLegs {
			if orderComboLeg.Price != UNSET_FLOAT {
				c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support per-leg prices for order combo legs.", "")
				return
			}

		}
	}
	if c.serverVersion < MIN_SERVER_VER_TRAILING_PERCENT && order.TrailingPercent != UNSET_FLOAT {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support trailing percent parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tradingClass parameter in placeOrder.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SCALE_TABLE &&
		(order.ScaleTable != "" || order.ActiveStartTime != "" || order.ActiveStopTime != "") {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support scaleTable, activeStartTime and activeStopTime parameters.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_ALGO_ID && order.AlgoID != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support algoId parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_ORDER_SOLICITED && order.Solicited {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support order solicited parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT && order.ModelCode != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support model code parameter.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_EXT_OPERATOR && order.ExtOperator != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support ext operator parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SOFT_DOLLAR_TIER && (order.SoftDollarTier.Name != "" || order.SoftDollarTier.Value != "") {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support soft dollar tier", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CASH_QTY && order.CashQty != UNSET_FLOAT {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support cash quantity parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_DECISION_MAKER && (order.Mifid2DecisionMaker != "" || order.Mifid2DecisionAlgo != "") {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support MIFID II decision maker parameters", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MIFID_EXECUTION && (order.Mifid2ExecutionTrader != "" || order.Mifid2ExecutionAlgo != "") {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support MIFID II execution parameters", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE && order.DontUseAutoPriceForHedge {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support dontUseAutoPriceForHedge parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_ORDER_CONTAINER && order.IsOmsContainer {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support oms container parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PRICE_MGMT_ALGO && order.UsePriceMgmtAlgo != USE_PRICE_MGMT_ALGO_DEFAULT {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support Use price management algo requests", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_DURATION && order.Duration != UNSET_INT {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support duration attribute", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_POST_TO_ATS && order.PostToAts != UNSET_INT {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support postToAts attribute", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_AUTO_CANCEL_PARENT && order.AutoCancelParent {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support autoCancelParent attribute", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_ADVANCED_ORDER_REJECT && order.AdvancedErrorOverride != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support advanced error override attribute", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PEGBEST_PEGMID_OFFSETS {
		if order.MinTradeQty != UNSET_INT ||
			order.MinCompeteSize != UNSET_INT ||
			order.CompeteAgainstBestOffset != UNSET_FLOAT ||
			order.MidOffsetAtWhole != UNSET_FLOAT ||
			order.MidOffsetAtHalf != UNSET_FLOAT {
			c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+
				" It does not support PEG BEST / PEG MID order parameters: minTradeQty, minCompeteSize, "+
				"competeAgainstBestOffset, midOffsetAtWhole and midOffsetAtHalf.", "")
			return
		}
	}

	if c.serverVersion < MIN_SERVER_VER_CUSTOMER_ACCOUNT && order.CustomerAccount != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support customer account parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PROFESSIONAL_CUSTOMER && order.ProfessionalCustomer {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support professional customer parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_INCLUDE_OVERNIGHT && order.IncludeOvernight {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support include overnight parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS && order.ManualOrderIndicator != UNSET_INT {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support manual indicator parameter", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_IMBALANCE_ONLY && order.ImbalanceOnly {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support imbalance only parameter", "")
		return
	}

	var VERSION int
	if c.serverVersion < MIN_SERVER_VER_NOT_HELD {
		VERSION = 27
	} else {
		VERSION = 45
	}

	// send place order msg
	me := NewMsgEncoder(150, c)

	me.encodeMsgID(PLACE_ORDER)

	if c.serverVersion < MIN_SERVER_VER_ORDER_CONTAINER {
		me.encodeInt(VERSION)
	}

	me.encodeInt64(orderID)

	// send contract fields
	if c.serverVersion >= MIN_SERVER_VER_PLACE_ORDER_CONID {
		me.encodeInt64(contract.ConID)
	}
	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier) // srv v15 and above
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange) // srv v14 and above
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol) // srv v2 and above

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	if c.serverVersion >= MIN_SERVER_VER_SEC_ID_TYPE {
		me.encodeString(contract.SecIDType)
		me.encodeString(contract.SecID)
	}

	// send main order fields
	me.encodeString(order.Action)

	if c.serverVersion >= MIN_SERVER_VER_FRACTIONAL_POSITIONS {
		me.encodeDecimal(order.TotalQuantity)
	} else {
		me.encodeDecimal(order.TotalQuantity)
	}

	me.encodeString(order.OrderType)

	if c.serverVersion < MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE {
		if order.LmtPrice != UNSET_FLOAT {
			me.encodeFloat64(order.LmtPrice)
		} else {
			me.encodeFloat64(0)
		}
	} else {
		me.encodeFloatMax(order.LmtPrice)
	}

	if c.serverVersion < MIN_SERVER_VER_TRAILING_PERCENT {
		if order.AuxPrice != UNSET_FLOAT {
			me.encodeFloat64(order.AuxPrice)
		} else {
			me.encodeFloat64(0)
		}
	} else {
		me.encodeFloatMax(order.AuxPrice)
	}

	// send extended order fields
	me.encodeString(order.TIF)
	me.encodeString(order.OCAGroup)
	me.encodeString(order.Account)
	me.encodeString(order.OpenClose)
	me.encodeInt64(order.Origin)
	me.encodeString(order.OrderRef)
	me.encodeBool(order.Transmit)
	me.encodeInt64(order.ParentID) // srv v4 and above

	me.encodeBool(order.BlockOrder)   // srv v5 and above
	me.encodeBool(order.SweepToFill)  // srv v5 and above
	me.encodeInt64(order.DisplaySize) // srv v5 and above
	me.encodeInt64(order.TriggerMethod)
	// srv v5 and above
	me.encodeBool(order.OutsideRTH)
	// srv v5 and above
	me.encodeBool(order.Hidden) // srv v7 and above

	// Send combo legs for BAG requests (srv v8 and above)
	if contract.SecType == "BAG" {
		comboLegsCount := len(contract.ComboLegs)
		me.encodeInt(comboLegsCount)
		for _, comboLeg := range contract.ComboLegs {
			me.encodeInt64(comboLeg.ConID)
			me.encodeInt64(comboLeg.Ratio)
			me.encodeString(comboLeg.Action)
			me.encodeString(comboLeg.Exchange)
			me.encodeInt64(comboLeg.OpenClose)

			me.encodeInt64(comboLeg.ShortSaleSlot)       // srv v35 and above
			me.encodeString(comboLeg.DesignatedLocation) // srv v35 and above
			if c.serverVersion >= MIN_SERVER_VER_SSHORTX_OLD {
				me.encodeInt64(comboLeg.ExemptCode)
			}
		}
	}

	// Send order combo legs for BAG requests
	if c.serverVersion >= MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE && contract.SecType == "BAG" {
		orderComboLegsCount := len(order.OrderComboLegs)
		me.encodeInt(orderComboLegsCount)
		for _, orderComboLeg := range order.OrderComboLegs {
			me.encodeFloatMax(orderComboLeg.Price)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_SMART_COMBO_ROUTING_PARAMS && contract.SecType == "BAG" {
		smartComboRoutingParamsCount := len(order.SmartComboRoutingParams)
		me.encodeInt(smartComboRoutingParamsCount)
		for _, tv := range order.SmartComboRoutingParams {
			me.encodeString(tv.Tag)
			me.encodeString(tv.Value)
		}
	}

	//  Send the shares allocation.
	//
	//  This specifies the number of order shares allocated to each Financial
	//  Advisor managed account. The format of the allocation string is as
	//  follows:
	//                       <account_code1>/<number_shares1>,<account_code2>/<number_shares2>,...N
	//  E.g.
	//               To allocate 20 shares of a 100 share order to account 'U101' and the
	//       residual 80 to account 'U203' enter the following share allocation string:
	//       #          U101/20,U203/80

	// send deprecated sharesAllocation field
	me.encodeString("")

	me.encodeFloat64(order.DiscretionaryAmt) //srv v10 and above
	me.encodeString(order.GoodAfterTime)     //srv v11 and above
	me.encodeString(order.GoodTillDate)      //srv v12 and above

	me.encodeString(order.FAGroup)      //srv v13 and above
	me.encodeString(order.FAMethod)     //srv v13 and above
	me.encodeString(order.FAPercentage) //srv v13 and above

	if c.serverVersion < MIN_SERVER_VER_FA_PROFILE_DESUPPORT {
		me.encodeString("") // send deprecated faProfile field
	}

	if c.serverVersion >= MIN_SERVER_VER_MODELS_SUPPORT {
		me.encodeString(order.ModelCode)
	}

	// institutional short saleslot data (srv v18 and above)
	me.encodeInt64(order.ShortSaleSlot)       // 0 for retail, 1 or 2 for institutions
	me.encodeString(order.DesignatedLocation) // populate only when shortSaleSlot = 2.

	if c.serverVersion >= MIN_SERVER_VER_SSHORTX_OLD {
		me.encodeInt64(order.ExemptCode)
	}

	// srv v19 and above fields
	me.encodeInt64(order.OCAType)

	me.encodeString(order.Rule80A)
	me.encodeString(order.SettlingFirm)
	me.encodeBool(order.AllOrNone)
	me.encodeIntMax(order.MinQty)
	me.encodeFloatMax(order.PercentOffset)
	me.encodeBool(false) // send deprecated order.ETradeOnly
	me.encodeBool(false) // send deprecated order.FirmQuoteOnly
	me.encodeFloatMax(UNSET_FLOAT)
	me.encodeInt64(order.AuctionStrategy) // AUCTION_MATCH, AUCTION_IMPROVEMENT, AUCTION_TRANSPARENT
	me.encodeFloatMax(order.StartingPrice)
	me.encodeFloatMax(order.StockRefPrice)
	me.encodeFloatMax(order.Delta)
	me.encodeFloatMax(order.StockRangeLower)
	me.encodeFloatMax(order.StockRangeUpper)

	me.encodeBool(order.OverridePercentageConstraints) // srv v22 and above

	// Volatility orders (srv v26 and above)
	me.encodeFloatMax(order.Volatility)
	me.encodeIntMax(order.VolatilityType)
	me.encodeString(order.DeltaNeutralOrderType)  // srv v28 and above
	me.encodeFloatMax(order.DeltaNeutralAuxPrice) // srv v28 and above

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL_CONID && order.DeltaNeutralOrderType != "" {
		me.encodeInt64(order.DeltaNeutralConID)
		me.encodeString(order.DeltaNeutralSettlingFirm)
		me.encodeString(order.DeltaNeutralClearingAccount)
		me.encodeString(order.DeltaNeutralClearingIntent)
	}

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL_OPEN_CLOSE && order.DeltaNeutralOrderType != "" {
		me.encodeString(order.DeltaNeutralOpenClose)
		me.encodeBool(order.DeltaNeutralShortSale)
		me.encodeInt64(order.DeltaNeutralShortSaleSlot)
		me.encodeString(order.DeltaNeutralDesignatedLocation)
	}

	me.encodeBool(order.ContinuousUpdate)

	me.encodeIntMax(order.ReferencePriceType)

	me.encodeFloatMax(order.TrailStopPrice) // srv v30 and above

	if c.serverVersion >= MIN_SERVER_VER_TRAILING_PERCENT {
		me.encodeFloatMax(order.TrailingPercent)
	}

	// scale orders
	if c.serverVersion >= MIN_SERVER_VER_SCALE_ORDERS2 {
		me.encodeIntMax(order.ScaleInitLevelSize)
		me.encodeIntMax(order.ScaleSubsLevelSize)
	} else {
		// srv v35 and above
		me.encodeString("")                       // for not supported scaleNumComponents
		me.encodeIntMax(order.ScaleInitLevelSize) // for scaleComponentSize
	}

	me.encodeFloatMax(order.ScalePriceIncrement)

	if c.serverVersion >= MIN_SERVER_VER_SCALE_ORDERS3 && order.ScalePriceIncrement != UNSET_FLOAT && order.ScalePriceIncrement > 0.0 {
		me.encodeFloatMax(order.ScalePriceAdjustValue)
		me.encodeIntMax(order.ScalePriceAdjustInterval)
		me.encodeFloatMax(order.ScaleProfitOffset)
		me.encodeBool(order.ScaleAutoReset)
		me.encodeIntMax(order.ScaleInitPosition)
		me.encodeIntMax(order.ScaleInitFillQty)
		me.encodeBool(order.ScaleRandomPercent)
	}

	if c.serverVersion >= MIN_SERVER_VER_SCALE_TABLE {
		me.encodeString(order.ScaleTable)
		me.encodeString(order.ActiveStartTime)
		me.encodeString(order.ActiveStopTime)
	}

	// hedge orders
	if c.serverVersion >= MIN_SERVER_VER_HEDGE_ORDERS {
		me.encodeString(order.HedgeType)
		if order.HedgeType != "" {
			me.encodeString(order.HedgeParam)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_OPT_OUT_SMART_ROUTING {
		me.encodeBool(order.OptOutSmartRouting)
	}

	if c.serverVersion >= MIN_SERVER_VER_PTA_ORDERS {
		me.encodeString(order.ClearingAccount)
		me.encodeString(order.ClearingIntent)
	}

	if c.serverVersion >= MIN_SERVER_VER_NOT_HELD {
		me.encodeBool(order.NotHeld)
	}

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL {
		if contract.DeltaNeutralContract != nil {
			me.encodeBool(true)
			me.encodeInt64(contract.DeltaNeutralContract.ConID)
			me.encodeFloat64(contract.DeltaNeutralContract.Delta)
			me.encodeFloat64(contract.DeltaNeutralContract.Price)
		} else {
			me.encodeBool(false)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_ALGO_ORDERS {
		me.encodeString(order.AlgoStrategy)

		if order.AlgoStrategy != "" {
			algoParamsCount := len(order.AlgoParams)
			me.encodeInt(algoParamsCount)
			for _, tv := range order.AlgoParams {
				me.encodeString(tv.Tag)
				me.encodeString(tv.Value)
			}
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_ALGO_ID {
		me.encodeString(order.AlgoID)
	}

	me.encodeBool(order.WhatIf) // srv v36 and above

	// send miscOptions parameter
	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(order.OrderMiscOptions)
	}

	if c.serverVersion >= MIN_SERVER_VER_ORDER_SOLICITED {
		me.encodeBool(order.Solicited)
	}

	if c.serverVersion >= MIN_SERVER_VER_RANDOMIZE_SIZE_AND_PRICE {
		me.encodeBool(order.RandomizeSize)
		me.encodeBool(order.RandomizePrice)
	}

	if c.serverVersion >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		if order.OrderType == "PEG BENCH" {
			me.encodeInt64(order.ReferenceContractID)
			me.encodeBool(order.IsPeggedChangeAmountDecrease)
			me.encodeFloat64(order.PeggedChangeAmount)
			me.encodeFloat64(order.ReferenceChangeAmount)
			me.encodeString(order.ReferenceExchangeID)
		}

		orderConditionsCount := len(order.Conditions)
		me.encodeInt(orderConditionsCount)
		for _, cond := range order.Conditions {
			me.encodeInt64(cond.Type())
			me.encodeFields(cond.makeFields()...)
		}

		if orderConditionsCount > 0 {
			me.encodeBool(order.ConditionsIgnoreRth)
			me.encodeBool(order.ConditionsCancelOrder)
		}

		me.encodeString(order.AdjustedOrderType)
		me.encodeFloat64(order.TriggerPrice)
		me.encodeFloat64(order.LmtPriceOffset)
		me.encodeFloat64(order.AdjustedStopPrice)
		me.encodeFloat64(order.AdjustedStopLimitPrice)
		me.encodeFloat64(order.AdjustedTrailingAmount)
		me.encodeInt64(order.AdjustableTrailingUnit)
	}
	if c.serverVersion >= MIN_SERVER_VER_EXT_OPERATOR {
		me.encodeString(order.ExtOperator)
	}

	if c.serverVersion >= MIN_SERVER_VER_SOFT_DOLLAR_TIER {
		me.encodeString(order.SoftDollarTier.Name)
		me.encodeString(order.SoftDollarTier.Value)
	}

	if c.serverVersion >= MIN_SERVER_VER_CASH_QTY {
		me.encodeFloatMax(order.CashQty)
	}

	if c.serverVersion >= MIN_SERVER_VER_DECISION_MAKER {
		me.encodeString(order.Mifid2DecisionMaker)
		me.encodeString(order.Mifid2DecisionAlgo)
	}

	if c.serverVersion >= MIN_SERVER_VER_MIFID_EXECUTION {
		me.encodeString(order.Mifid2ExecutionTrader)
		me.encodeString(order.Mifid2ExecutionAlgo)
	}

	if c.serverVersion >= MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE {
		me.encodeBool(order.DontUseAutoPriceForHedge)
	}

	if c.serverVersion >= MIN_SERVER_VER_ORDER_CONTAINER {
		me.encodeBool(order.IsOmsContainer)
	}

	if c.serverVersion >= MIN_SERVER_VER_D_PEG_ORDERS {
		me.encodeBool(order.DiscretionaryUpToLimitPrice)
	}

	if c.serverVersion >= MIN_SERVER_VER_PRICE_MGMT_ALGO {
		me.encodeIntMax(order.UsePriceMgmtAlgo)
	}

	if c.serverVersion >= MIN_SERVER_VER_DURATION {
		me.encodeInt64(order.Duration)
	}

	if c.serverVersion >= MIN_SERVER_VER_POST_TO_ATS {
		me.encodeInt64(order.PostToAts)
	}

	if c.serverVersion >= MIN_SERVER_VER_AUTO_CANCEL_PARENT {
		me.encodeBool(order.AutoCancelParent)
	}

	if c.serverVersion >= MIN_SERVER_VER_ADVANCED_ORDER_REJECT {
		me.encodeString(order.AdvancedErrorOverride)
	}

	if c.serverVersion >= MIN_SERVER_VER_MANUAL_ORDER_TIME {
		me.encodeString(order.ManualOrderTime)
	}

	if c.serverVersion >= MIN_SERVER_VER_PEGBEST_PEGMID_OFFSETS {
		var sendMidOffsets bool
		if contract.Exchange == "IBKRATS" {
			me.encodeIntMax(order.MinTradeQty)
		}
		if order.OrderType == "PEG BEST" {
			me.encodeIntMax(order.MinCompeteSize)
			me.encodeFloatMax(order.CompeteAgainstBestOffset)
			if order.CompeteAgainstBestOffset == COMPETE_AGAINST_BEST_OFFSET_UP_TO_MID {
				sendMidOffsets = true
			}
		} else if order.OrderType == "PEG MID" {
			sendMidOffsets = true
		}
		if sendMidOffsets {
			me.encodeFloatMax(order.MidOffsetAtWhole)
			me.encodeFloatMax(order.MidOffsetAtHalf)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_CUSTOMER_ACCOUNT {
		me.encodeString(order.CustomerAccount)
	}

	if c.serverVersion >= MIN_SERVER_VER_PROFESSIONAL_CUSTOMER {
		me.encodeBool(order.ProfessionalCustomer)
	}

	if c.serverVersion >= MIN_SERVER_VER_RFQ_FIELDS && c.serverVersion < MIN_SERVER_VER_UNDO_RFQ_FIELDS {
		me.encodeString("")
		me.encodeInt64(UNSET_INT)
	}

	if c.serverVersion >= MIN_SERVER_VER_INCLUDE_OVERNIGHT {
		me.encodeBool(order.IncludeOvernight)
	}

	if c.serverVersion >= MIN_SERVER_VER_CME_TAGGING_FIELDS {
		me.encodeInt64(order.ManualOrderIndicator)
	}

	if c.serverVersion >= MIN_SERVER_VER_IMBALANCE_ONLY {
		me.encodeBool(order.ImbalanceOnly)
	}

	c.reqChan <- me.Bytes()
}

func (c *EClient) placeOrderProtoBuf(placeOrderRequestProto *protobuf.PlaceOrderRequest) {

	orderID := NO_VALID_ID
	if placeOrderRequestProto.OrderId != nil {
		orderID = int64(*placeOrderRequestProto.OrderId)
	}

	if !c.IsConnected() {
		c.wrapper.Error(orderID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(150, c)
	me.encodeMsgID(PLACE_ORDER + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(placeOrderRequestProto)
	if err != nil {
		c.wrapper.Error(orderID, currentTimeMillis(), 0, "Failed to marshal PlaceOrderRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

// CancelOrder cancel an order by orderId.
// It can only be used to cancel an order that was placed originally by a client with the same client ID
func (c *EClient) CancelOrder(orderID OrderID, orderCancel OrderCancel) {

	if c.useProtoBuf(CANCEL_ORDER) {
		c.cancelOrderProtoBuf(createCancelOrderRequestProto(orderID, &orderCancel))
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(orderID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MANUAL_ORDER_TIME && orderCancel.ManualOrderCancelTime != "" {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support manual order cancel time attribute.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS && (orderCancel.ExtOperator != "" || orderCancel.ManualOrderIndicator != UNSET_INT) {
		c.wrapper.Error(orderID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support ext operator and manual order indicator parameters.", "")
	}

	const VERSION = 1

	me := NewMsgEncoder(9, c)

	me.encodeMsgID(CANCEL_ORDER)

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS {
		me.encodeInt(VERSION)
	}

	me.encodeInt64(orderID)

	if c.serverVersion >= MIN_SERVER_VER_MANUAL_ORDER_TIME {
		me.encodeString(orderCancel.ManualOrderCancelTime)
	}

	if c.serverVersion >= MIN_SERVER_VER_RFQ_FIELDS && c.serverVersion < MIN_SERVER_VER_UNDO_RFQ_FIELDS {
		me.encodeString("")
		me.encodeString("")
		me.encodeInt64(UNSET_INT)
	}

	if c.serverVersion >= MIN_SERVER_VER_CME_TAGGING_FIELDS {
		me.encodeString(orderCancel.ExtOperator)
		me.encodeInt64(orderCancel.ManualOrderIndicator)
	}

	c.reqChan <- me.Bytes()
}

func (c *EClient) cancelOrderProtoBuf(cancelOrderRequestProto *protobuf.CancelOrderRequest) {

	orderID := NO_VALID_ID
	if cancelOrderRequestProto.OrderId != nil {
		orderID = int64(*cancelOrderRequestProto.OrderId)
	}

	if !c.IsConnected() {
		c.wrapper.Error(orderID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(9, c)
	me.encodeMsgID(CANCEL_ORDER + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(cancelOrderRequestProto)
	if err != nil {
		c.wrapper.Error(orderID, currentTimeMillis(), 0, "Failed to marshal CancelOrderRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

// CancelOrderAsync cancel an order by orderId.
// ReqOpenOrders requests the open orders that were placed from this client.
// Each open order will be fed back through the openOrder() and orderStatus() functions on the EWrapper.
// The client with a clientId of 0 will also receive the TWS-owned open orders.
// These orders will be associated with the client and a new orderId will be generated.
// This association will persist over multiple API and TWS sessions.
func (c *EClient) ReqOpenOrders() {

	if c.useProtoBuf(REQ_OPEN_ORDERS) {
		c.reqOpenOrdersProtoBuf(createOpenOrdersRequestProto())
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_OPEN_ORDERS)
	me.encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

func (c *EClient) reqOpenOrdersProtoBuf(openOrdersRequestProto *protobuf.OpenOrdersRequest) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(2, c)
	me.encodeMsgID(REQ_OPEN_ORDERS + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(openOrdersRequestProto)
	if err != nil {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), 0, "Failed to marshal OpenOrdersRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

// ReqAutoOpenOrders requests that newly created TWS orders be implicitly associated with the client.
// When a new TWS order is created, the order will be associated with the client, and fed back through the openOrder() and orderStatus() functions on the EWrapper.
// This request can only be made from a client with clientId of 0.
// if autoBind is set to TRUE, newly created TWS orders will be implicitly associated with the client.
// If set to FALSE, no association will be made.
func (c *EClient) ReqAutoOpenOrders(autoBind bool) {

	if c.useProtoBuf(REQ_AUTO_OPEN_ORDERS) {
		c.reqAutoOpenOrdersProtoBuf(createAutoOpenOrdersRequestProto(autoBind))
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(REQ_AUTO_OPEN_ORDERS)
	me.encodeInt(VERSION)
	me.encodeBool(autoBind)

	c.reqChan <- me.Bytes()
}

func (c *EClient) reqAutoOpenOrdersProtoBuf(autoOpenOrdersRequestProto *protobuf.AutoOpenOrdersRequest) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(3, c)
	me.encodeMsgID(REQ_AUTO_OPEN_ORDERS + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(autoOpenOrdersRequestProto)
	if err != nil {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), 0, "Failed to marshal AutoOpenOrdersRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

// ReqAllOpenOrders request the open orders placed from all clients and also from TWS.
// Each open order will be fed back through the openOrder() and orderStatus() functions on the EWrapper.
// No association is made between the returned orders and the requesting client.
func (c *EClient) ReqAllOpenOrders() {

	if c.useProtoBuf(REQ_ALL_OPEN_ORDERS) {
		c.reqAllOpenOrdersProtoBuf(createAllOpenOrdersRequestProto())
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_ALL_OPEN_ORDERS)
	me.encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

func (c *EClient) reqAllOpenOrdersProtoBuf(allOpenOrdersRequestProto *protobuf.AllOpenOrdersRequest) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(2, c)
	me.encodeMsgID(REQ_ALL_OPEN_ORDERS + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(allOpenOrdersRequestProto)
	if err != nil {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), 0, "Failed to marshal AllOpenOrdersRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

// ReqGlobalCancel cancels all open orders globally. It cancels both API and TWS open orders.
func (c *EClient) ReqGlobalCancel(orderCancel OrderCancel) {

	if c.useProtoBuf(REQ_GLOBAL_CANCEL) {
		c.reqGlobalCancelProtoBuf(createGlobalCancelRequestProto(&orderCancel))
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS && (orderCancel.ExtOperator != "" || orderCancel.ManualOrderIndicator != UNSET_INT) {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support ext operator and manual order indicator parameters.", "")
	}

	const VERSION = 1

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(REQ_GLOBAL_CANCEL)

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS {
		me.encodeInt(VERSION)
	}

	if c.serverVersion >= MIN_SERVER_VER_CME_TAGGING_FIELDS {
		me.encodeString(orderCancel.ExtOperator)
		me.encodeInt64(orderCancel.ManualOrderIndicator)
	}

	c.reqChan <- me.Bytes()
}

func (c *EClient) reqGlobalCancelProtoBuf(globalCancelRequestProto *protobuf.GlobalCancelRequest) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(4, c)
	me.encodeMsgID(REQ_GLOBAL_CANCEL + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(globalCancelRequestProto)
	if err != nil {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), 0, "Failed to marshal GlobalCancelRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

// ReqIDs request from TWS the next valid ID that can be used when placing an order.
// After calling this function, the nextValidId() event will be triggered, and the id returned is that next valid ID.
// That ID will reflect any autobinding that has occurred (which generates new IDs and increments the next valid ID therein).
// numIds is depreceted
func (c *EClient) ReqIDs(numIds int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(REQ_IDS)
	me.encodeInt(VERSION)
	me.encodeInt64(numIds)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Account and Portfolio
// 	##########################################################################

// ReqAccountUpdates will start getting account values, portfolio, and last update time information.
// it is returned via EWrapper.updateAccountValue(), EWrapperi.updatePortfolio() and Wrapper.updateAccountTime().
func (c *EClient) ReqAccountUpdates(subscribe bool, accountName string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 2

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(REQ_ACCT_DATA)
	me.encodeInt(VERSION)
	me.encodeBool(subscribe) // TRUE = subscribe, FALSE = unsubscribe.

	// Send the account code. This will only be used for FA clients
	me.encodeString(accountName) // srv v9 and above

	c.reqChan <- me.Bytes()
}

// ReqAccountSummary request and keep up to date the data that appears.
// on the TWS Account Window Summary tab. The data is returned by accountSummary().
// This request is designed for an FA managed account but can be used for any multi-account structure.
// reqId is the ID of the data request. it Ensures that responses are matched to requests If several requests are in process.
// groupName sets to All to return account summary data for all accounts, or set to a specific Advisor Account Group name that has
// already been created in TWS Global Configuration.
// tags:str - A comma-separated list of account tags.  Available tags are:
//
//	accountountType
//	NetLiquidation,
//	TotalCashValue - Total cash including futures pnl
//	SettledCash - For cash accounts, this is the same as
//	TotalCashValue
//	AccruedCash - Net accrued interest
//	BuyingPower - The maximum dollar value of securities that you can buy without depositing additional equity
//	EquityWithLoanValue - Cash + stocks + bonds + mutual funds
//	PreviousDayEquityWithLoanValue,
//	GrossPositionValue - The sum of the absolute value of all stock and equity option positions
//	RegTEquity,
//	RegTMargin,
//	SMA - Special Memorandum Account
//	InitMarginReq,
//	MaintMarginReq,
//	AvailableFunds,
//	ExcessLiquidity,
//	Cushion - Excess liquidity as a percentage of net liquidation value
//	FullInitMarginReq,
//	FullMaintMarginReq,
//	FullAvailableFunds,
//	FullExcessLiquidity,
//	LookAheadNextChange - Time when look-ahead values take effect
//	LookAheadInitMarginReq,
//	LookAheadMaintMarginReq,
//	LookAheadAvailableFunds,
//	LookAheadExcessLiquidity,
//	HighestSeverity - A measure of how close the account is to liquidation.
//	DayTradesRemaining - The Number of Open/Close trades a user could put on before Pattern Day Trading is detected.
//		A value of "-1"	means that the user can put on unlimited day trades.
//	Leverage - GrossPositionValue / NetLiquidation
//	$LEDGER - Single flag to relay all cash balance tags*, only in base	currency.
//	$LEDGER:CURRENCY - Single flag to relay all cash balance tags*, only in	the specified currency.
//	$LEDGER:ALL - Single flag to relay all cash balance tags* in all currencies.
func (c *EClient) ReqAccountSummary(reqID int64, groupName string, tags string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(5, c)

	me.encodeMsgID(REQ_ACCOUNT_SUMMARY)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)
	me.encodeString(groupName)
	me.encodeString(tags)

	c.reqChan <- me.Bytes()
}

// CancelAccountSummary cancels the request for Account Window Summary tab data.
// reqId is the ID of the data request being canceled.
func (c *EClient) CancelAccountSummary(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_ACCOUNT_SUMMARY)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqPositions requests real-time position data for all accounts.
func (c *EClient) ReqPositions() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_POSITIONS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support positions request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_POSITIONS)
	me.encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

// CancelPositions cancels real-time position updates.
func (c *EClient) CancelPositions() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_POSITIONS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support positions request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_POSITIONS)
	me.encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

// ReqPositionsMulti requests the positions for account and/or model.
// Results are delivered via EWrapper.positionMulti() and EWrapper.positionMultiEnd().
func (c *EClient) ReqPositionsMulti(reqID int64, account string, modelCode string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support positions multi request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(5, c)

	me.encodeMsgID(REQ_POSITIONS_MULTI)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)
	me.encodeString(account)
	me.encodeString(modelCode)

	c.reqChan <- me.Bytes()
}

// CancelPositionsMulti cancels the positions update of assigned account.
func (c *EClient) CancelPositionsMulti(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support cancel positions multi request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_POSITIONS_MULTI)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqAccountUpdatesMulti requests account updates for account and/or model.
func (c *EClient) ReqAccountUpdatesMulti(reqID int64, account string, modelCode string, ledgerAndNLV bool) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support account updates multi request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(6, c)

	me.encodeMsgID(REQ_ACCOUNT_UPDATES_MULTI)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)
	me.encodeString(account)
	me.encodeString(modelCode)
	me.encodeBool(ledgerAndNLV)

	c.reqChan <- me.Bytes()
}

// CancelAccountUpdatesMulti cancels account update for reqID.
func (c *EClient) CancelAccountUpdatesMulti(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MODELS_SUPPORT {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support cancel account updates multi request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_ACCOUNT_UPDATES_MULTI)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Daily PnL
// 	##########################################################################

// ReqPnL requests and subscribe the PnL of assigned account.
func (c *EClient) ReqPnL(reqID int64, account string, modelCode string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PNL {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support PnL request.", "")
		return
	}

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(REQ_PNL)
	me.encodeInt64(reqID)
	me.encodeString(account)
	me.encodeString(modelCode)

	c.reqChan <- me.Bytes()
}

// CancelPnL cancels the PnL update of assigned account.
func (c *EClient) CancelPnL(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PNL {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support PnL request.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_PNL)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqPnLSingle request and subscribe the single contract PnL of assigned account.
func (c *EClient) ReqPnLSingle(reqID int64, account string, modelCode string, contractID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PNL {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support PnL request.", "")
		return
	}

	me := NewMsgEncoder(5, c)

	me.encodeMsgID(REQ_PNL_SINGLE)
	me.encodeInt64(reqID)
	me.encodeString(account)
	me.encodeString(modelCode)
	me.encodeInt64(contractID)

	c.reqChan <- me.Bytes()
}

// CancelPnLSingle cancel the single contract PnL update of assigned account.
func (c *EClient) CancelPnLSingle(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PNL {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support PnL request.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_PNL_SINGLE)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Executions
// 	##########################################################################

// ReqExecutions downloads the execution reports that meet the filter criteria to the client via the execDetails() function.
// To view executions beyond the past 24 hours, open the Trade Log in TWS and, while the Trade Log is displayed, request the executions again from the API.
// reqId is the ID of the data request. Ensures that responses are matched to requests if several requests are in process.
// execFilter contains attributes that describe the filter criteria used to determine which execution reports are returned.
// NOTE: Time format must be 'yyyymmdd-hh:mm:ss' Eg: '20030702-14:55'
func (c *EClient) ReqExecutions(reqID int64, execFilter *ExecutionFilter) {

	if c.useProtoBuf(REQ_EXECUTIONS) {
		c.reqExecutionProtobuf(createExecutionRequestProto(reqID, execFilter))
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_PARAMETRIZED_DAYS_OF_EXECUTIONS && (execFilter.LastNDays != UNSET_INT || execFilter.SpecificDates != nil) {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support last N days and specific dates parameters.", "")
		return
	}

	const VERSION = 3

	me := NewMsgEncoder(14, c)

	me.encodeMsgID(REQ_EXECUTIONS)
	me.encodeInt(VERSION)

	if c.serverVersion >= MIN_SERVER_VER_EXECUTION_DATA_CHAIN {
		me.encodeInt64(reqID)
	}

	me.encodeInt64(execFilter.ClientID)
	me.encodeString(execFilter.AcctCode)
	me.encodeString(execFilter.Time)
	me.encodeString(execFilter.Symbol)
	me.encodeString(execFilter.SecType)
	me.encodeString(execFilter.Exchange)
	me.encodeString(execFilter.Side)

	if c.serverVersion >= MIN_SERVER_VER_PARAMETRIZED_DAYS_OF_EXECUTIONS {
		me.encodeInt64(execFilter.LastNDays)
		specificDatesCount := len(execFilter.SpecificDates)
		me.encodeInt(specificDatesCount)
		for _, date := range execFilter.SpecificDates {
			me.encodeInt64(date)
		}
	}

	c.reqChan <- me.Bytes()
}

func (c *EClient) reqExecutionProtobuf(executionRequestProto *protobuf.ExecutionRequest) {

	reqID := NO_VALID_ID
	if executionRequestProto.ReqId != nil {
		reqID = int64(*executionRequestProto.ReqId)
	}

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(0, c)
	me.encodeMsgID(REQ_EXECUTIONS + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(executionRequestProto)
	if err != nil {
		log.Panic().Err(err).Msg("reqExecutionProtobuf marshal error")
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Contract Details
// 	##########################################################################

// ReqContractDetails downloads all details for a particular underlying.
// The contract details will be received via the contractDetails() function on the EWrapper.
func (c *EClient) ReqContractDetails(reqID int64, contract *Contract) {

	if c.useProtoBuf(REQ_CONTRACT_DATA) {
		c.reqContractDataProtoBuf(createContractDataRequestProto(reqID, contract))
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SEC_ID_TYPE && (contract.SecIDType != "" || contract.SecID != "") {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support secIdType and secId parameters.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tradingClass parameter in reqContractDetails.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING && contract.PrimaryExchange != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support primaryExchange parameter in reqContractDetails.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_BOND_ISSUERID && contract.IssuerID != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support issuerId parameter in reqContractDetails.", "")
		return
	}

	const VERSION = 8

	me := NewMsgEncoder(21, c)

	me.encodeMsgID(REQ_CONTRACT_DATA)
	me.encodeInt(VERSION)

	if c.serverVersion >= MIN_SERVER_VER_CONTRACT_DATA_CHAIN {
		me.encodeInt64(reqID)
	}

	me.encodeInt64(contract.ConID) // srv v37 and above
	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier) // srv v15 and above

	if c.serverVersion >= MIN_SERVER_VER_PRIMARYEXCH {
		me.encodeString(contract.Exchange)
		me.encodeString(contract.PrimaryExchange)
	} else if c.serverVersion >= MIN_SERVER_VER_LINKING {
		if contract.PrimaryExchange != "" && (contract.Exchange == "BEST" || contract.Exchange == "SMART") {
			me.encodeString(contract.Exchange + ":" + contract.PrimaryExchange)
		} else {
			me.encodeString(contract.Exchange)
		}
	}

	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)

	}
	me.encodeBool(contract.IncludeExpired) //  srv v31 and above

	if c.serverVersion >= MIN_SERVER_VER_SEC_ID_TYPE {
		me.encodeString(contract.SecIDType)
		me.encodeString(contract.SecID)
	}

	if c.serverVersion >= MIN_SERVER_VER_BOND_ISSUERID {
		me.encodeString(contract.IssuerID)
	}

	c.reqChan <- me.Bytes()
}

func (c *EClient) reqContractDataProtoBuf(contractDataRequestProto *protobuf.ContractDataRequest) {

	reqID := NO_VALID_ID
	if contractDataRequestProto.ReqId != nil {
		reqID = int64(*contractDataRequestProto.ReqId)
	}

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(150, c)
	me.encodeMsgID(REQ_CONTRACT_DATA + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(contractDataRequestProto)
	if err != nil {
		c.wrapper.Error(reqID, currentTimeMillis(), 0, "Failed to marshal PlaceOrderRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Market Depth
// 	##########################################################################

// ReqMktDepthExchanges requests market depth exchanges.
func (c *EClient) ReqMktDepthExchanges() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_MKT_DEPTH_EXCHANGES {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support market depth exchanges request.", "")
		return
	}

	me := NewMsgEncoder(1, c)

	me.encodeMsgID(REQ_MKT_DEPTH_EXCHANGES)

	c.reqChan <- me.Bytes()
}

// ReqMktDepth requests the market depth for a specific contract.
// The market depth will be returned by the updateMktDepth() and updateMktDepthL2() events.
// Requests the contract's market depth (order book). Note this request must be direct-routed to an exchange and not smart-routed.
// The number of simultaneous market depth requests allowed in an account is calculated based on a formula
// that looks at an accounts equity, commissions and fees, and quote booster packs.
// reqId is the ticker id. It must be a unique value. When the market depth data returns, it will be identified by this tag.
// This is also used when canceling the market depth
// contract contains a description of the contract for which market depth data is being requested.
// numRows specifies the numRowsumber of market depth rows to display.
// isSmartDepth	specifies SMART depth request.
// mktDepthOptions is for internal use only. Use default value XYZ.
func (c *EClient) ReqMktDepth(reqID int64, contract *Contract, numRows int, isSmartDepth bool, mktDepthOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		if contract.TradingClass != "" || contract.ConID > 0 {
			c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support conId and tradingClass parameters in reqMktDepth.", "")
			return
		}
	}

	if c.serverVersion < MIN_SERVER_VER_SMART_DEPTH && isSmartDepth {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support SMART depth request.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_MKT_DEPTH_PRIM_EXCHANGE && contract.PrimaryExchange != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support primaryExchange parameter in reqMktDepth.", "")
		return
	}

	const VERSION = 5

	me := NewMsgEncoder(17, c)

	// send req mkt data msg
	me.encodeMsgID(REQ_MKT_DEPTH)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	// send contract fields
	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeInt64(contract.ConID)
	}

	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier) // srv v15 and above
	me.encodeString(contract.Exchange)

	if c.serverVersion >= MIN_SERVER_VER_MKT_DEPTH_PRIM_EXCHANGE {
		me.encodeString(contract.PrimaryExchange)
	}

	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	me.encodeInt(numRows) // srv v19 and above

	if c.serverVersion >= MIN_SERVER_VER_SMART_DEPTH {
		me.encodeBool(isSmartDepth)
	}

	// send mktDepthOptions parameter
	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(mktDepthOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelMktDepth cancels market depth updates.
func (c *EClient) CancelMktDepth(reqID int64, isSmartDepth bool) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SMART_DEPTH && isSmartDepth {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support SMART depth cancel.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(CANCEL_MKT_DEPTH)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	if c.serverVersion >= MIN_SERVER_VER_SMART_DEPTH {
		me.encodeBool(isSmartDepth)
	}

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		News Bulletins
// 	##########################################################################

// ReqNewsBulletins requests and subcribe the news bulletins.
// Each bulletin will be returned by the updateNewsBulletin() event.
//
// If allMsgs sets to TRUE, returns all the existing bulletins for the currencyent day and any new ones.
// If allMsgs sets to FALSE, will only return new bulletins.
func (c *EClient) ReqNewsBulletins(allMsgs bool) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(REQ_NEWS_BULLETINS)
	me.encodeInt(VERSION)
	me.encodeBool(allMsgs)

	c.reqChan <- me.Bytes()
}

// CancelNewsBulletins cancels the news bulletins updates
func (c *EClient) CancelNewsBulletins() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_NEWS_BULLETINS)
	me.encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Financials Advisor
// 	##########################################################################

// ReqManagedAccts requests the list of managed accounts.
// The result will be delivered via wrapper.ManagedAccounts().
// This request can only be made when connected to a FA managed account.
func (c *EClient) ReqManagedAccts() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_MANAGED_ACCTS)
	me.encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

// RequestFA requests fa.
// The data returns in an XML string via wrapper.ReceiveFA().
// faData is 1->"GROUPS", 3->"ALIASES".
func (c *EClient) RequestFA(faDataType FaDataType) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_FA_PROFILE_DESUPPORT && faDataType == 2 {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), FA_PROFILE_NOT_SUPPORTED.Code, FA_PROFILE_NOT_SUPPORTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(REQ_FA)
	me.encodeInt(VERSION)
	me.encodeInt(int(faDataType))

	c.reqChan <- me.Bytes()
}

// ReplaceFA replaces the FA configuration information from the API.
// Note that this can also be done manually in TWS itself.
// faData specifies the type of Financial Advisor configuration data being requested.
// 1 = GROUPS
// 3 = ACCOUNT ALIASES
// cxml is the XML string containing the new FA configuration information.
func (c *EClient) ReplaceFA(reqID int64, faDataType FaDataType, cxml string) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion >= MIN_SERVER_VER_REPLACE_FA_END && faDataType == 2 {
		c.wrapper.Error(reqID, currentTimeMillis(), FA_PROFILE_NOT_SUPPORTED.Code, FA_PROFILE_NOT_SUPPORTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(5, c)

	me.encodeMsgID(REPLACE_FA)
	me.encodeInt(VERSION)
	me.encodeInt(int(faDataType))
	me.encodeString(cxml)

	if c.serverVersion >= MIN_SERVER_VER_REPLACE_FA_END {
		me.encodeInt64(reqID)
	}

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Historical Data
// 	##########################################################################

// ReqHistoricalData requests historical data and subcribe the new data if keepUpToDate is assigned.
// Requests contracts' historical data. When requesting historical data, a finishing time and date is required along with a duration string.
// Result will be delivered via wrapper.HistoricalData()
// reqId the id of the request. Must be a unique value. When the market data returns, it whatToShowill be identified by this tag.
// This is also used when canceling the market data.
// contract contains a description of the contract for which market data is being requested.
// endDateTime defines a query end date and time at any point during the past six months in the format:
// yyyymmdd HH:mm:ss ttt where "ttt" is the optional time zone.
// durationStr set the query duration up to one week, using a time unit of seconds, days or weeks.
// Valid values include any integer followed by a space and then S (seconds), D (days) or W (week).
// If no unit is specified, seconds is used.
// barSizeSetting specifies the size of the bars that will be returned (within IB/TWS listimits).
// Valid values include:
// 	1 sec
// 	5 secs
// 	15 secs
// 	30 secs
// 	1 min
// 	2 mins
// 	3 mins
// 	5 mins
// 	15 mins
// 	30 mins
// 	1 hour
// 	1 day
// whatToShow determines the nature of data beinging extracted.
// Valid values include:
// 	TRADES
// 	MIDPOINT
// 	BID
// 	ASK
// 	BID_ASK
// 	HISTORICAL_VOLATILITY
// 	OPTION_IMPLIED_VOLATILITY
// useRTH determines whether to return all data available during the requested time span,
// or only data that falls within regular trading hours.
// Valid values include:
// 	0 - all data is returned even where the market in question was outside of its
// 	regular trading hours.
// 	1 - only data within the regular trading hours is returned, even if the
// 	requested time span falls partially or completely outside of the RTH.
// formatDate determines the date format applied to returned bars.
// Valid values include:
// 	1 - dates applying to bars returned in the format: yyyymmdd{space}{space}hh:mm:dd
// 	2 - dates are returned as a long integer specifying the number of seconds since
// 		1/1/1970 GMT.
// chartOptions is for internal use only. Use default value XYZ.

func (c *EClient) ReqHistoricalData(reqID int64, contract *Contract, endDateTime string, duration string, barSize string, whatToShow string, useRTH bool, formatDate int, keepUpToDate bool, chartOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		if contract.TradingClass != "" || contract.ConID > 0 {
			c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg, "")
		}
	}

	const VERSION = 6

	me := NewMsgEncoder(20, c)

	me.encodeMsgID(REQ_HISTORICAL_DATA)

	if c.serverVersion <= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		me.encodeInt(VERSION)
	}

	me.encodeInt64(reqID)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeInt64(contract.ConID)
	}

	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier)
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	me.encodeBool(contract.IncludeExpired)
	me.encodeString(endDateTime)
	me.encodeString(barSize)
	me.encodeString(duration)
	me.encodeBool(useRTH)
	me.encodeString(whatToShow)
	me.encodeInt(formatDate)

	if contract.SecType == "BAG" {
		me.encodeInt(len(contract.ComboLegs))
		for _, comboLeg := range contract.ComboLegs {
			me.encodeInt64(comboLeg.ConID)
			me.encodeInt64(comboLeg.Ratio)
			me.encodeString(comboLeg.Action)
			me.encodeString(comboLeg.Exchange)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		me.encodeBool(keepUpToDate)
	}

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(chartOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelHistoricalData cancels the update of historical data.
// Used if an internet disconnect has occurred or the results of a query are otherwise delayed and the application is no longer interested in receiving the data.
// reqId, the ticker ID, must be a unique value.
func (c *EClient) CancelHistoricalData(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_HISTORICAL_DATA)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqHeadTimeStamp request the head timestamp of assigned contract.
// call this func to get the headmost data you can get
func (c *EClient) ReqHeadTimeStamp(reqID int64, contract *Contract, whatToShow string, useRTH bool, formatDate int) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_HEAD_TIMESTAMP {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support head time stamp requests.", "")
		return
	}

	me := NewMsgEncoder(19, c)

	me.encodeMsgID(REQ_HEAD_TIMESTAMP)
	me.encodeInt64(reqID)
	me.encodeContract(contract)
	me.encodeBool(useRTH)
	me.encodeString(whatToShow)
	me.encodeInt(formatDate)

	c.reqChan <- me.Bytes()
}

// CancelHeadTimeStamp cancels the head timestamp data.
func (c *EClient) CancelHeadTimeStamp(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CANCEL_HEADTIMESTAMP {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support head time stamp requests.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_HEAD_TIMESTAMP)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqHistogramData requests histogram data.
func (c *EClient) ReqHistogramData(reqID int64, contract *Contract, useRTH bool, timePeriod string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_HISTOGRAM {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support histogram requests..", "")
		return
	}

	me := NewMsgEncoder(5, c)

	me.encodeMsgID(REQ_HISTOGRAM_DATA)
	me.encodeInt64(reqID)
	me.encodeContract(contract)
	me.encodeBool(useRTH)
	me.encodeString(timePeriod)

	c.reqChan <- me.Bytes()
}

// CancelHistogramData cancels histogram data.
func (c *EClient) CancelHistogramData(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_HISTOGRAM {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support histogram requests..", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_HISTOGRAM_DATA)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqHistoricalTicks requests historical ticks.
func (c *EClient) ReqHistoricalTicks(reqID int64, contract *Contract, startDateTime string, endDateTime string, numberOfTicks int, whatToShow string, useRTH bool, ignoreSize bool, miscOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_HISTORICAL_TICKS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support historical ticks requests..", "")
		return
	}

	me := NewMsgEncoder(22, c)

	me.encodeMsgID(REQ_HISTORICAL_TICKS)
	me.encodeInt64(reqID)
	me.encodeContract(contract)
	me.encodeString(startDateTime)
	me.encodeString(endDateTime)
	me.encodeInt(numberOfTicks)
	me.encodeString(whatToShow)
	me.encodeBool(useRTH)
	me.encodeBool(ignoreSize)
	me.encodeTagValues(miscOptions)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Market Scanners
// 	##########################################################################

// ReqScannerParameters requests an XML string that describes all possible scanner queries.
func (c *EClient) ReqScannerParameters() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_SCANNER_PARAMETERS)
	me.encodeInt(VERSION)

	c.reqChan <- me.Bytes()
}

// ReqScannerSubscription subcribes a scanner that matched the subcription.
// reqId, the ticker ID, must be a unique value.
// scannerSubscription contains possible parameters used to filter results.
// scannerSubscriptionOptions is for internal use only.Use default value XYZ.
func (c *EClient) ReqScannerSubscription(reqID int64, subscription *ScannerSubscription, scannerSubscriptionOptions []TagValue, scannerSubscriptionFilterOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SCANNER_GENERIC_OPTS && len(scannerSubscriptionFilterOptions) > 0 {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support API scanner subscription generic filter options", "")
		return
	}

	const VERSION = 4

	me := NewMsgEncoder(25, c)

	me.encodeMsgID(REQ_SCANNER_SUBSCRIPTION)

	if c.serverVersion < MIN_SERVER_VER_SCANNER_GENERIC_OPTS {
		me.encodeInt(VERSION)
	}

	me.encodeInt64(reqID)
	me.encodeIntMax(subscription.NumberOfRows)
	me.encodeString(subscription.Instrument)
	me.encodeString(subscription.LocationCode)
	me.encodeString(subscription.ScanCode)
	me.encodeFloatMax(subscription.AbovePrice)
	me.encodeFloatMax(subscription.BelowPrice)
	me.encodeIntMax(subscription.AboveVolume)
	me.encodeFloatMax(subscription.MarketCapAbove)
	me.encodeFloatMax(subscription.MarketCapBelow)
	me.encodeString(subscription.MoodyRatingAbove)
	me.encodeString(subscription.MoodyRatingBelow)
	me.encodeString(subscription.SpRatingAbove)
	me.encodeString(subscription.SpRatingBelow)
	me.encodeString(subscription.MaturityDateAbove)
	me.encodeString(subscription.MaturityDateBelow)
	me.encodeFloatMax(subscription.CouponRateAbove)
	me.encodeFloatMax(subscription.CouponRateBelow)
	me.encodeBool(subscription.ExcludeConvertible)
	me.encodeIntMax(subscription.AverageOptionVolumeAbove)
	me.encodeString(subscription.ScannerSettingPairs)
	me.encodeString(subscription.StockTypeFilter)

	if c.serverVersion >= MIN_SERVER_VER_SCANNER_GENERIC_OPTS {
		me.encodeTagValues(scannerSubscriptionFilterOptions)
	}

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(scannerSubscriptionOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelScannerSubscription cancel scanner.
// reqId is the unique ticker ID used for subscription.
func (c *EClient) CancelScannerSubscription(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_SCANNER_SUBSCRIPTION)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Real Time Bars
// 	##########################################################################

// ReqRealTimeBars requests realtime bars.
// Result will be delivered via wrapper.RealtimeBar().
// reqId, the ticker ID, must be a unique value. When the data is received, it will be identified by this Id.
// This is also used when canceling the request.
// contract contains a description of the contract for which real time bars are being requested.
// barSize, Currently only supports 5 second bars, if any other	value is used, an exception will be thrown.
// whatToShow determines the nature of the data extracted.
// Valid includes:
//
//	TRADES
//	BID
//	ASK
//	MIDPOINT
//
// useRTH sets regular Trading Hours only.
// Valid values include:
//
//	0 = all data available during the time span requested is returned,
//		including time intervals when the market in question was
//		outside of regular trading hours.
//	1 = only data within the regular trading hours for the product
//		requested is returned, even if the time time span falls
//		partially or completely outside.
//
// realTimeBarOptions is for internal use only. Use default value XYZ.
func (c *EClient) ReqRealTimeBars(reqID int64, contract *Contract, barSize int, whatToShow string, useRTH bool, realTimeBarsOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		c.wrapper.Error(reqID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support conId and tradingClass parameter in reqRealTimeBars.", "")
		return
	}

	const VERSION = 3

	me := NewMsgEncoder(19, c)

	me.encodeMsgID(REQ_REAL_TIME_BARS)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeInt64(contract.ConID)
	}

	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.LastTradeDateOrContractMonth)
	me.encodeFloatMax(contract.Strike)
	me.encodeString(contract.Right)
	me.encodeString(contract.Multiplier)
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeString(contract.TradingClass)
	}

	me.encodeInt(barSize)
	me.encodeString(whatToShow)
	me.encodeBool(useRTH)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(realTimeBarsOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelRealTimeBars cancels realtime bars.
func (c *EClient) CancelRealTimeBars(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_REAL_TIME_BARS)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Fundamental Data
// 	##########################################################################

// ReqFundamentalData requests fundamental data for stocks.
// The appropriate market data subscription must be set up in Account Management before you can receive this data.
// Result will be delivered via wrapper.FundamentalData().
// this func can handle conid specified in the Contract object, but not tradingClass or multiplier.
// This is because this func is used only for stocks and stocks do not have a multiplier and trading class.
// reqId is	the ID of the data request. Ensures that responses are matched to requests if several requests are in process.
// contract contains a description of the contract for which fundamental data is being requested.
// reportType is one of the following XML reports:
//
//	ReportSnapshot (company overview)
//	ReportsFinSummary (financial summary)
//	ReportRatios (financial ratios)
//	ReportsFinStatements (financial statements)
//	RESC (analyst estimates)
//	CalendarReport (company calendar)
func (c *EClient) ReqFundamentalData(reqID int64, contract *Contract, reportType string, fundamentalDataOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_FUNDAMENTAL_DATA {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support fundamental data request.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support conId parameter in reqFundamentalData.", "")
		return
	}

	const VERSION = 2

	me := NewMsgEncoder(12, c)

	me.encodeMsgID(REQ_FUNDAMENTAL_DATA)

	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		me.encodeInt64(contract.ConID)
	}

	me.encodeString(contract.Symbol)
	me.encodeString(contract.SecType)
	me.encodeString(contract.Exchange)
	me.encodeString(contract.PrimaryExchange)
	me.encodeString(contract.Currency)
	me.encodeString(contract.LocalSymbol)

	me.encodeString(reportType)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		me.encodeTagValues(fundamentalDataOptions)
	}

	c.reqChan <- me.Bytes()
}

// CancelFundamentalData cancels fundamental data.
func (c *EClient) CancelFundamentalData(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_FUNDAMENTAL_DATA {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support fundamental data request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(CANCEL_FUNDAMENTAL_DATA)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()

}

//	##########################################################################
//	#		News
// 	##########################################################################

// ReqNewsProviders request news providers.
func (c *EClient) ReqNewsProviders() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_NEWS_PROVIDERS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support news providers request.", "")
		return
	}

	me := NewMsgEncoder(1, c)

	me.encodeMsgID(REQ_NEWS_PROVIDERS)

	c.reqChan <- me.Bytes()
}

// ReqNewsArticle request news article.
func (c *EClient) ReqNewsArticle(reqID int64, providerCode string, articleID string, newsArticleOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_NEWS_ARTICLE {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support news article request.", "")
		return
	}

	me := NewMsgEncoder(5, c)

	me.encodeMsgID(REQ_NEWS_ARTICLE)
	me.encodeInt64(reqID)
	me.encodeString(providerCode)
	me.encodeString(articleID)

	if c.serverVersion >= MIN_SERVER_VER_NEWS_QUERY_ORIGINS {
		me.encodeTagValues(newsArticleOptions)
	}

	c.reqChan <- me.Bytes()
}

// ReqHistoricalNews request historical news.
func (c *EClient) ReqHistoricalNews(reqID int64, contractID int64, providerCode string, startDateTime string, endDateTime string, totalResults int64, historicalNewsOptions []TagValue) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_HISTORICAL_NEWS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support historical news request.", "")
		return
	}

	me := NewMsgEncoder(8, c)

	me.encodeMsgID(REQ_HISTORICAL_NEWS)
	me.encodeInt64(reqID)
	me.encodeInt64(contractID)
	me.encodeString(providerCode)
	me.encodeString(startDateTime)
	me.encodeString(endDateTime)
	me.encodeInt64(totalResults)

	if c.serverVersion >= MIN_SERVER_VER_NEWS_QUERY_ORIGINS {
		me.encodeTagValues(historicalNewsOptions)
	}

	c.reqChan <- me.Bytes()
}

//	##########################################################################
//	#		Display Groups
// 	##########################################################################

// QueryDisplayGroups request the display groups in TWS.
func (c *EClient) QueryDisplayGroups(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support queryDisplayGroups request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(QUERY_DISPLAY_GROUPS)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// SubscribeToGroupEvents subcribes the group events.
// reqId is the unique number associated with the notification.
// groupId is the ID of the group, currently it is a number from 1 to 7.
func (c *EClient) SubscribeToGroupEvents(reqID int64, groupID int) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support subscribeToGroupEvents request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(SUBSCRIBE_TO_GROUP_EVENTS)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)
	me.encodeInt(groupID)

	c.reqChan <- me.Bytes()
}

// UpdateDisplayGroup updates the display group in TWS.
// reqId is the requestId specified in subscribeToGroupEvents().
// contractInfo is the encoded value that uniquely represents the contract in IB.
// Possible values include:
//
//	none = empty selection
//	contractID@exchange - any non-combination contract.
//		Examples: 8314@SMART for IBM SMART; 8314@ARCA for IBM @ARCA.
//	combo = if any combo is selected.
func (c *EClient) UpdateDisplayGroup(reqID int64, contractInfo string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support updateDisplayGroup request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(UPDATE_DISPLAY_GROUP)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)
	me.encodeString(contractInfo)

	c.reqChan <- me.Bytes()
}

// UnsubscribeFromGroupEvents unsubcribes the display group events.
func (c *EClient) UnsubscribeFromGroupEvents(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support unsubscribeFromGroupEvents request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(UNSUBSCRIBE_FROM_GROUP_EVENTS)
	me.encodeInt(VERSION)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// VerifyRequest is just for IB's internal use.
// Allows to provide means of verification between the TWS and third party programs.
func (c *EClient) VerifyRequest(apiName string, apiVersion string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support verification request.", "")
		return
	}

	if c.extraAuth {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), BAD_MESSAGE.Code, BAD_MESSAGE.Msg+
			" Intent to authenticate needs to be expressed during initial connect request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(VERIFY_REQUEST)
	me.encodeInt(VERSION)
	me.encodeString(apiName)
	me.encodeString(apiVersion)

	c.reqChan <- me.Bytes()
}

// VerifyMessage is just for IB's internal use.
// Allows to provide means of verification between the TWS and third party programs.
func (c *EClient) VerifyMessage(apiData string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support verification request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(VERIFY_MESSAGE)
	me.encodeInt(VERSION)
	me.encodeString(apiData)

	c.reqChan <- me.Bytes()
}

// VerifyAndAuthRequest is just for IB's internal use.
// Allows to provide means of verification between the TWS and third party programs.
func (c *EClient) VerifyAndAuthRequest(apiName string, apiVersion string, opaqueIsvKey string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support verification request.", "")
		return
	}

	if c.extraAuth {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+
			" Intent to authenticate needs to be expressed during initial connect request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(VERIFY_AND_AUTH_REQUEST)
	me.encodeInt(VERSION)
	me.encodeString(apiName)
	me.encodeString(apiVersion)
	me.encodeString(opaqueIsvKey)

	c.reqChan <- me.Bytes()
}

// VerifyAndAuthMessage is just for IB's internal use.
// Allows to provide means of verification between the TWS and third party programs.
func (c *EClient) VerifyAndAuthMessage(apiData string, xyzResponse string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_LINKING {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support verification request.", "")
		return
	}

	const VERSION = 1

	me := NewMsgEncoder(4, c)

	me.encodeMsgID(VERIFY_MESSAGE)
	me.encodeInt(VERSION)
	me.encodeString(apiData)
	me.encodeString(xyzResponse)

	c.reqChan <- me.Bytes()
}

// ReqSecDefOptParams requests security definition option parameters.
// reqId the ID chosen for the request underlyingSymbol.
// futFopExchange is the exchange on which the returned options are trading. Can be set to the empty string "" for all exchanges.
// underlyingSecType is the type of the underlying security, i.e. STK.
// underlyingConId is the contract ID of the underlying security.
// Response comes via wrapper.SecurityDefinitionOptionParameter().
func (c *EClient) ReqSecDefOptParams(reqID int64, underlyingSymbol string, futFopExchange string, underlyingSecurityType string, underlyingContractID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_SEC_DEF_OPT_PARAMS_REQ {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support security definition option request.", "")
		return
	}

	me := NewMsgEncoder(6, c)

	me.encodeMsgID(REQ_SEC_DEF_OPT_PARAMS)
	me.encodeInt64(reqID)
	me.encodeString(underlyingSymbol)
	me.encodeString(futFopExchange)
	me.encodeString(underlyingSecurityType)
	me.encodeInt64(underlyingContractID)

	c.reqChan <- me.Bytes()
}

// ReqSoftDollarTiers request pre-defined Soft Dollar Tiers.
// This is only supported for registered professional advisors and hedge and mutual funds
// who have configured Soft Dollar Tiers in Account Management.
func (c *EClient) ReqSoftDollarTiers(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_SOFT_DOLLAR_TIERS)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqFamilyCodes requests family codes.
func (c *EClient) ReqFamilyCodes() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_FAMILY_CODES {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support family codes request.", "")
		return
	}

	me := NewMsgEncoder(1, c)

	me.encodeMsgID(REQ_FAMILY_CODES)

	c.reqChan <- me.Bytes()
}

// ReqMatchingSymbols requests matching symbols.
func (c *EClient) ReqMatchingSymbols(reqID int64, pattern string) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_REQ_MATCHING_SYMBOLS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support matching symbols request.", "")
		return
	}

	me := NewMsgEncoder(3, c)

	me.encodeMsgID(REQ_MATCHING_SYMBOLS)
	me.encodeInt64(reqID)
	me.encodeString(pattern)

	c.reqChan <- me.Bytes()
}

// ReqCompletedOrders request the completed orders.
// If apiOnly parameter is true, then only completed orders placed from API are requested.
// Result will be delivered via wrapper.CompletedOrder().
func (c *EClient) ReqCompletedOrders(apiOnly bool) {

	if c.useProtoBuf(REQ_COMPLETED_ORDERS) {
		c.reqCompletedOrdersProtoBuf(createCompletedOrdersRequestProto(apiOnly))
		return
	}

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_COMPLETED_ORDERS)
	me.encodeBool(apiOnly)

	c.reqChan <- me.Bytes()
}

func (c *EClient) reqCompletedOrdersProtoBuf(completedOrdersRequestProto *protobuf.CompletedOrdersRequest) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	me := NewMsgEncoder(3, c)
	me.encodeMsgID(REQ_COMPLETED_ORDERS + PROTOBUF_MSG_ID)

	msg, err := proto.Marshal(completedOrdersRequestProto)
	if err != nil {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), 0, "Failed to marshal CompletedOrdersRequest: "+err.Error(), "")
		return
	}

	me.encodeProto(msg)

	c.reqChan <- me.Bytes()
}

// ReqWshMetaData requests WSHE Meta data.
func (c *EClient) ReqWshMetaData(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_WSHE_CALENDAR {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support WSHE Calendar API.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_WSH_META_DATA)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// CancelWshMetaData cancels WSHE Meta data.
func (c *EClient) CancelWshMetaData(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_WSHE_CALENDAR {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support WSHE Calendar API.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_WSH_META_DATA)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqWshEventData requests WSHE Event data.
func (c *EClient) ReqWshEventData(reqID int64, wshEventData WshEventData) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_WSHE_CALENDAR {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support WSHE Calendar API.", "")
		return
	}
	if c.serverVersion < MIN_SERVER_VER_WSH_EVENT_DATA_FILTERS &&
		(wshEventData.Filter != "" || wshEventData.FillWatchList || wshEventData.FillPortfolio) {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support WSHE event data filters.", "")
		return
	}
	if c.serverVersion < MIN_SERVER_VER_WSH_EVENT_DATA_FILTERS_DATE &&
		(wshEventData.StartDate != "" || wshEventData.EndDate != "" || wshEventData.TotalLimit != UNSET_INT) {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support WSHE event data date filters.", "")
		return
	}
	me := NewMsgEncoder(10, c)

	me.encodeMsgID(REQ_WSH_EVENT_DATA)
	me.encodeInt64(reqID)
	me.encodeInt64(wshEventData.ConID)

	if c.serverVersion >= MIN_SERVER_VER_WSH_EVENT_DATA_FILTERS {
		me.encodeString(wshEventData.Filter)
		me.encodeBool(wshEventData.FillWatchList)
		me.encodeBool(wshEventData.FillPortfolio)
		me.encodeBool(wshEventData.FillCompetitors)
	}

	if c.serverVersion >= MIN_SERVER_VER_WSH_EVENT_DATA_FILTERS_DATE {
		me.encodeString(wshEventData.StartDate)
		me.encodeString(wshEventData.EndDate)
		me.encodeInt64(wshEventData.TotalLimit)
	}

	c.reqChan <- me.Bytes()
}

// CancelWshEventData cancels WSHE Event data.
func (c *EClient) CancelWshEventData(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_WSHE_CALENDAR {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support WSHE Calendar API.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(CANCEL_WSH_EVENT_DATA)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqUserInfo requests user info.
func (c *EClient) ReqUserInfo(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_USER_INFO {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support user info requests.", "")
		return
	}

	me := NewMsgEncoder(2, c)

	me.encodeMsgID(REQ_USER_INFO)
	me.encodeInt64(reqID)

	c.reqChan <- me.Bytes()
}

// ReqCurrentTimeInMillis requests the current system time in milliseconds on the server side.
func (c *EClient) ReqCurrentTimeInMillis() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CURRENT_TIME_IN_MILLIS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support current time in millis requests", "")
		return
	}

	me := NewMsgEncoder(1, c)

	me.encodeMsgID(REQ_CURRENT_TIME_IN_MILLIS)

	c.reqChan <- me.Bytes()
}
