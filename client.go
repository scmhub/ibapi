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
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type ConnState int

const (
	DISCONNECTED ConnState = iota
	CONNECTING
	CONNECTED
	REDIRECT
)

func (cs ConnState) String() string {
	switch cs {
	case DISCONNECTED:
		return "disconnected"
	case CONNECTING:
		return "connecting"
	case CONNECTED:
		return "connected"
	case REDIRECT:
		return "redirect"
	default:
		return "unknown connection state"
	}
}

// EClient is the main struct to use from API user's point of view.
type EClient struct {
	host           string
	port           int
	clientID       int64
	connectOptions string
	conn           *Connection
	serverVersion  Version
	connTime       string
	connState      ConnState
	writer         *bufio.Writer
	scanner        *bufio.Scanner
	wrapper        EWrapper
	decoder        *EDecoder
	reqChan        chan []byte
	Ctx            context.Context
	Cancel         context.CancelFunc
	extraAuth      bool
	wg             sync.WaitGroup
	err            error
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
	c.extraAuth = false
	c.conn = &Connection{}
	c.serverVersion = -1
	c.connTime = ""

	// writer
	c.writer = bufio.NewWriter(c.conn)
	// init scanner
	c.scanner = bufio.NewScanner(c.conn)
	c.scanner.Split(scanFields)
	c.scanner.Buffer(make([]byte, 4096), MAX_MSG_LEN)

	c.reqChan = make(chan []byte, 10)

	c.Ctx, c.Cancel = context.WithCancel(context.Background())

	c.wg = sync.WaitGroup{}
	c.err = nil

	c.setConnState(DISCONNECTED)
	c.connectOptions = ""
}

func (c *EClient) setConnState(state ConnState) {
	cs := c.connState
	c.connState = state
	log.Debug().Stringer("from", cs).Stringer("to", c.connState).Msg("connection state changed")
}

// request is a goroutine that will get the req from reqChan and send it to TWS.
func (c *EClient) request() {
	log.Debug().Msg("requester started")
	defer log.Debug().Msg("requester ended")

	c.wg.Add(1)
	defer c.wg.Done()

	for {
		select {
		case <-c.Ctx.Done():
			return
		case req := <-c.reqChan:
			if !c.IsConnected() {
				c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
				break
			}
			nn, err := c.writer.Write(req)
			if err != nil {
				log.Error().Err(err).Int("nbytes", nn).Bytes("reqMsg", req).Msg("requester write error")
				break
			}
			err = c.writer.Flush()
			if err != nil {
				log.Error().Err(err).Bytes("reqMsg", req).Msg("requester flush error")
				c.writer.Reset(c.conn)
			}
		}
	}
}

// startAPI initiates the message exchange between the client application and the TWS/IB Gateway.
func (c *EClient) startAPI() error {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return NOT_CONNECTED
	}

	var msg []byte

	const VERSION = 2

	if c.serverVersion >= MIN_SERVER_VER_OPTIONAL_CAPABILITIES {
		msg = makeFields(START_API, VERSION, c.clientID, "")
	} else {
		msg = makeFields(START_API, VERSION, c.clientID)
	}

	if _, err := c.writer.Write(msg); err != nil {
		return err
	}
	if err := c.writer.Flush(); err != nil {
		return err
	}

	return nil
}

// Connect must be called before any other.
// There is no feedback for a successful connection, but a subsequent attempt to connect will return the message "Already connected.".
func (c *EClient) Connect(host string, port int, clientID int64) error {
	c.host = host
	c.port = port
	c.clientID = clientID

	log.Info().Str("host", host).Int("port", port).Int64("clientID", clientID).Msg("Connecting to IB server")
	if err := c.conn.connect(c.host, c.port); err != nil {
		log.Error().Err(CONNECT_FAIL).Msg("Connection fail")
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), CONNECT_FAIL.Code, CONNECT_FAIL.Msg, "")
		c.reset()
		return CONNECT_FAIL
	}

	// HandShake with the TWS or GateWay to ensure the version,
	log.Debug().Msg("HandShake with TWS or GateWay")

	head := []byte("API\x00")

	connectOptions := ""
	if c.connectOptions != "" {
		connectOptions = " " + c.connectOptions
	}
	sizeofCV := make([]byte, 4)
	clientVersion := []byte(fmt.Sprintf("v%d..%d%s", MIN_CLIENT_VER, MAX_CLIENT_VER, connectOptions))

	binary.BigEndian.PutUint32(sizeofCV, uint32(len(clientVersion)))

	var msg bytes.Buffer
	msg.Write(head)
	msg.Write(sizeofCV)
	msg.Write(clientVersion)

	log.Debug().Bytes("header", msg.Bytes()).Msg("send handShake header")

	if _, err := c.writer.Write(msg.Bytes()); err != nil {
		return err
	}

	if err := c.writer.Flush(); err != nil {
		return err
	}

	log.Debug().Msg("recv handShake Info")

	// scan once to get server info
	if !c.scanner.Scan() {
		return c.scanner.Err()
	}

	// Init server info
	msgBytes := c.scanner.Bytes()
	serverInfo := splitMsgBytes(msgBytes)
	v, _ := strconv.Atoi(string(serverInfo[0]))
	c.serverVersion = Version(v)
	c.connTime = string(serverInfo[1])
	log.Info().Int("serverVersion", v).Str("connectionTime", c.connTime).Msg("Handshake completed")

	// init decoder
	c.decoder = &EDecoder{wrapper: c.wrapper, serverVersion: c.serverVersion}

	//start Ereader
	go EReader(c.Ctx, c.scanner, c.decoder, &c.wg)

	// start requester
	go c.request()

	c.setConnState(CONNECTED)
	c.wrapper.ConnectAck()

	// startAPI
	if err := c.startAPI(); err != nil {
		return err
	}
	log.Debug().Msg("API started")

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info().Msg("detected KeyboardInterrupt, SystemExit")
		c.Disconnect()
		os.Exit(1)
	}()

	log.Debug().Msg("IB Client Connected!")

	return nil
}

// Disconnect terminates the connections with TWS.
// Calling this function does not cancel orders that have already been sent.
func (c *EClient) Disconnect() error {
	if !c.IsConnected() {
		return nil
	}

	c.Cancel()

	if err := c.conn.disconnect(); err != nil {
		return err
	}

	c.wg.Wait()

	defer c.reset()
	defer c.wrapper.ConnectionClosed()

	defer log.Debug().Msg("IB Client Disconnected!")

	return c.err
}

// IsConnected checks connection to TWS or GateWay.
func (c *EClient) IsConnected() bool {
	return c.conn.IsConnected() && c.connState == CONNECTED
}

// SetConnectionOptions setup the Connection Options.
func (c *EClient) SetConnectionOptions(opts string) {
	c.connectOptions = opts
}

// ReqCurrentTime asks the current system time on the server side.
func (c *EClient) ReqCurrentTime() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION int64 = 1

	msg := makeFields(REQ_CURRENT_TIME, VERSION)

	c.reqChan <- msg
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

	msg := makeFields(SET_SERVER_LOGLEVEL, VERSION, logLevel)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 30)
	fields = append(fields,
		REQ_MKT_DATA,
		VERSION,
		reqID,
	)

	if c.serverVersion >= MIN_SERVER_VER_REQ_MKT_DATA_CONID {
		fields = append(fields, contract.ConID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier, // srv v15 and above
		contract.Exchange,
		contract.PrimaryExchange, // srv v14 and above
		contract.Currency,
		contract.LocalSymbol) // srv v2 and above

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	// Send combo legs for BAG requests (srv v8 and above)
	if contract.SecType == "BAG" {
		comboLegsCount := len(contract.ComboLegs)
		fields = append(fields, comboLegsCount)
		for _, comboLeg := range contract.ComboLegs {
			fields = append(fields,
				comboLeg.ConID,
				comboLeg.Ratio,
				comboLeg.Action,
				comboLeg.Exchange)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL {
		if contract.DeltaNeutralContract != nil {
			fields = append(fields,
				true,
				contract.DeltaNeutralContract.ConID,
				contract.DeltaNeutralContract.Delta,
				contract.DeltaNeutralContract.Price)
		} else {
			fields = append(fields, false)
		}
	}

	fields = append(fields,
		genericTickList, // srv v31 and above
		snapshot)        // srv v35 and above

	if c.serverVersion >= MIN_SERVER_VER_REQ_SMART_COMPONENTS {
		fields = append(fields, regulatorySnapshot)
	}

	// send mktDataOptions parameter
	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		//  current doc says this part if for "internal use only" -> won't support it
		if len(mktDataOptions) > 0 {
			log.Panic().Msg("not supported")
		}
		fields = append(fields, "")
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
}

// CancelMktData stops the market data flow for the specified TickerId.
func (c *EClient) CancelMktData(reqID TickerID) {

	if !c.IsConnected() {
		c.wrapper.Error(reqID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 2

	msg := makeFields(CANCEL_MKT_DATA, VERSION, reqID)

	c.reqChan <- msg
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

	msg := makeFields(REQ_MARKET_DATA_TYPE, VERSION, marketDataType)

	c.reqChan <- msg
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

	msg := makeFields(REQ_SMART_COMPONENTS, reqID, bboExchange)

	c.reqChan <- msg
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

	msg := makeFields(REQ_MARKET_RULE, marketRuleID)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 17)
	fields = append(fields, REQ_TICK_BY_TICK_DATA,
		reqID,
		contract.ConID,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
		contract.TradingClass,
		tickType)

	if c.serverVersion >= MIN_SERVER_VER_TICK_BY_TICK_IGNORE_SIZE {
		fields = append(fields, numberOfTicks, ignoreSize)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_TICK_BY_TICK_DATA, reqID)

	c.reqChan <- msg
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
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support calculateImpliedVolatility req.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS && contract.TradingClass != "" {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support tradingClass parameter in calculateImpliedVolatility.", "")
		return
	}

	const VERSION = 3

	fields := make([]interface{}, 0, 19)
	fields = append(fields,
		REQ_CALC_IMPLIED_VOLAT,
		VERSION,
		reqID,
		contract.ConID,
		contract.Symbol,
		contract.SecID,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	fields = append(fields, optionPrice, underPrice)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		var implVolOptBuffer bytes.Buffer
		tagValuesCount := len(impVolOptions)
		fields = append(fields, tagValuesCount)
		for _, tv := range impVolOptions {
			implVolOptBuffer.WriteString(tv.Tag)
			implVolOptBuffer.WriteString("=")
			implVolOptBuffer.WriteString(tv.Value)
			implVolOptBuffer.WriteString(";")
		}
		fields = append(fields, implVolOptBuffer.Bytes())
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_CALC_IMPLIED_VOLAT, VERSION, reqID)

	c.reqChan <- msg
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

	const VERSION = 3

	fields := make([]interface{}, 0, 19)
	fields = append(fields,
		REQ_CALC_OPTION_PRICE,
		VERSION,
		reqID,
		contract.ConID,
		contract.Symbol,
		contract.SecID,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	fields = append(fields, volatility, underPrice)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		var optPrcOptBuffer bytes.Buffer
		tagValuesCount := len(optPrcOptions)
		fields = append(fields, tagValuesCount)
		for _, tv := range optPrcOptions {
			optPrcOptBuffer.WriteString(tv.Tag)
			optPrcOptBuffer.WriteString("=")
			optPrcOptBuffer.WriteString(tv.Value)
			optPrcOptBuffer.WriteString(";")
		}

		fields = append(fields, optPrcOptBuffer.Bytes())
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_CALC_OPTION_PRICE, VERSION, reqID)

	c.reqChan <- msg
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
// override specifies whether your setting will override the system's natural action.
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

	fields := make([]interface{}, 0, 17)

	fields = append(fields, EXERCISE_OPTIONS, VERSION, reqID)

	// send contract fields
	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.ConID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.Currency,
		contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	fields = append(fields,
		exerciseAction,
		exerciseQuantity,
		account,
		override)

	if c.serverVersion >= MIN_SERVER_VER_MANUAL_ORDER_TIME_EXERCISE_OPTIONS {
		fields = append(fields, manualOrderTime)
	}

	if c.serverVersion >= MIN_SERVER_VER_CUSTOMER_ACCOUNT {
		fields = append(fields, customerAccount)
	}

	if c.serverVersion >= MIN_SERVER_VER_PROFESSIONAL_CUSTOMER {
		fields = append(fields, professionalCustomer)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	if c.serverVersion < MIN_SERVER_VER_ORDER_SOLICITED && order.Solictied {
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

	if c.serverVersion < MIN_SERVER_VER_PRICE_MGMT_ALGO && order.UsePriceMgmtAlgo {
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

	var VERSION int
	if c.serverVersion < MIN_SERVER_VER_NOT_HELD {
		VERSION = 27
	} else {
		VERSION = 45
	}

	// send place order msg
	fields := make([]interface{}, 0, 150)
	fields = append(fields, PLACE_ORDER)

	if c.serverVersion < MIN_SERVER_VER_ORDER_CONTAINER {
		fields = append(fields, VERSION)
	}

	fields = append(fields, orderID)

	// send contract fields
	if c.serverVersion >= MIN_SERVER_VER_PLACE_ORDER_CONID {
		fields = append(fields, contract.ConID)
	}
	fields = append(fields,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier, // srv v15 and above
		contract.Exchange,
		contract.PrimaryExchange, // srv v14 and above
		contract.Currency,
		contract.LocalSymbol) // srv v2 and above

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	if c.serverVersion >= MIN_SERVER_VER_SEC_ID_TYPE {
		fields = append(fields, contract.SecIDType, contract.SecID)
	}

	// send main order fields
	fields = append(fields, order.Action)

	if c.serverVersion >= MIN_SERVER_VER_FRACTIONAL_POSITIONS {
		fields = append(fields, order.TotalQuantity)
	} else {
		fields = append(fields, order.TotalQuantity.Int())
	}

	fields = append(fields, order.OrderType)

	if c.serverVersion < MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE {
		if order.LmtPrice != UNSET_FLOAT {
			fields = append(fields, order.LmtPrice)
		} else {
			fields = append(fields, float64(0))
		}
	} else {
		fields = append(fields, handleEmpty(order.LmtPrice))
	}

	if c.serverVersion < MIN_SERVER_VER_TRAILING_PERCENT {
		if order.AuxPrice != UNSET_FLOAT {
			fields = append(fields, order.AuxPrice)
		} else {
			fields = append(fields, float64(0))
		}
	} else {
		fields = append(fields, handleEmpty(order.AuxPrice))
	}

	// send extended order fields
	fields = append(fields,
		order.TIF,
		order.OCAGroup,
		order.Account,
		order.OpenClose,
		order.Origin,
		order.OrderRef,
		order.Transmit,
		order.ParentID,      // srv v4 and above
		order.BlockOrder,    // srv v5 and above
		order.SweepToFill,   // srv v5 and above
		order.DisplaySize,   // srv v5 and above
		order.TriggerMethod, // srv v5 and above
		order.OutsideRTH,    // srv v5 and above
		order.Hidden)        // srv v7 and above

	// Send combo legs for BAG requests (srv v8 and above)
	if contract.SecType == "BAG" {
		comboLegsCount := len(contract.ComboLegs)
		fields = append(fields, comboLegsCount)
		for _, comboLeg := range contract.ComboLegs {
			fields = append(fields,
				comboLeg.ConID,
				comboLeg.Ratio,
				comboLeg.Action,
				comboLeg.Exchange,
				comboLeg.OpenClose,
				comboLeg.ShortSaleSlot,      // srv v35 and above
				comboLeg.DesignatedLocation) // srv v35 and above
			if c.serverVersion >= MIN_SERVER_VER_SSHORTX_OLD {
				fields = append(fields, comboLeg.ExemptCode)
			}
		}
	}

	// Send order combo legs for BAG requests
	if c.serverVersion >= MIN_SERVER_VER_ORDER_COMBO_LEGS_PRICE && contract.SecType == "BAG" {
		orderComboLegsCount := len(order.OrderComboLegs)
		fields = append(fields, orderComboLegsCount)
		for _, orderComboLeg := range order.OrderComboLegs {
			fields = append(fields, handleEmpty(orderComboLeg.Price))
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_SMART_COMBO_ROUTING_PARAMS && contract.SecType == "BAG" {
		smartComboRoutingParamsCount := len(order.SmartComboRoutingParams)
		fields = append(fields, smartComboRoutingParamsCount)
		for _, tv := range order.SmartComboRoutingParams {
			fields = append(fields, tv.Tag, tv.Value)
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
	fields = append(fields,
		"",                     // srv v9 and above
		order.DiscretionaryAmt, //srv v10 and above
		order.GoodAfterTime,    //srv v11 and above
		order.GoodTillDate,     //srv v12 and above

		order.FAGroup,      //srv v13 and above
		order.FAMethod,     //srv v13 and above
		order.FAPercentage, //srv v13 and above
	)

	if c.serverVersion < MIN_SERVER_VER_FA_PROFILE_DESUPPORT {
		fields = append(fields, "") // send deprecated faProfile field
	}

	if c.serverVersion >= MIN_SERVER_VER_MODELS_SUPPORT {
		fields = append(fields, order.ModelCode)
	}

	// institutional short saleslot data (srv v18 and above)
	fields = append(fields,
		order.ShortSaleSlot,      // 0 for retail, 1 or 2 for institutions
		order.DesignatedLocation) // populate only when shortSaleSlot = 2.

	if c.serverVersion >= MIN_SERVER_VER_SSHORTX_OLD {
		fields = append(fields, order.ExemptCode)
	}

	// srv v19 and above fields
	fields = append(fields, order.OCAType)

	fields = append(fields,
		order.Rule80A,
		order.SettlingFirm,
		order.AllOrNone,
		handleEmpty(order.MinQty),
		handleEmpty(order.PercentOffset),
		false, // send deprecated order.ETradeOnly
		false, // send deprecated order.FirmQuoteOnly
		handleEmpty(UNSET_FLOAT),
		order.AuctionStrategy, // AUCTION_MATCH, AUCTION_IMPROVEMENT, AUCTION_TRANSPARENT
		handleEmpty(order.StartingPrice),
		handleEmpty(order.StockRefPrice),
		handleEmpty(order.Delta),
		handleEmpty(order.StockRangeLower),
		handleEmpty(order.StockRangeUpper),

		order.OverridePercentageConstraints, // srv v22 and above

		// Volatility orders (srv v26 and above)
		handleEmpty(order.Volatility),
		handleEmpty(order.VolatilityType),
		order.DeltaNeutralOrderType,             // srv v28 and above
		handleEmpty(order.DeltaNeutralAuxPrice)) // srv v28 and above

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL_CONID && order.DeltaNeutralOrderType != "" {
		fields = append(fields,
			order.DeltaNeutralConID,
			order.DeltaNeutralSettlingFirm,
			order.DeltaNeutralClearingAccount,
			order.DeltaNeutralClearingIntent)
	}

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL_OPEN_CLOSE && order.DeltaNeutralOrderType != "" {
		fields = append(fields,
			order.DeltaNeutralOpenClose,
			order.DeltaNeutralShortSale,
			order.DeltaNeutralShortSaleSlot,
			order.DeltaNeutralDesignatedLocation)
	}

	fields = append(fields,
		order.ContinuousUpdate,
		handleEmpty(order.ReferencePriceType),
		handleEmpty(order.TrailStopPrice)) // srv v30 and above

	if c.serverVersion >= MIN_SERVER_VER_TRAILING_PERCENT {
		fields = append(fields, handleEmpty(order.TrailingPercent))
	}

	// scale orders
	if c.serverVersion >= MIN_SERVER_VER_SCALE_ORDERS2 {
		fields = append(fields,
			handleEmpty(order.ScaleInitLevelSize),
			handleEmpty(order.ScaleSubsLevelSize))
	} else {
		// srv v35 and above
		fields = append(fields,
			"",                                    // for not supported scaleNumComponents
			handleEmpty(order.ScaleInitLevelSize)) // for scaleComponentSize
	}

	fields = append(fields, handleEmpty(order.ScalePriceIncrement))

	if c.serverVersion >= MIN_SERVER_VER_SCALE_ORDERS3 && order.ScalePriceIncrement != UNSET_FLOAT && order.ScalePriceIncrement > 0.0 {
		fields = append(fields,
			handleEmpty(order.ScalePriceAdjustValue),
			handleEmpty(order.ScalePriceAdjustInterval),
			handleEmpty(order.ScaleProfitOffset),
			order.ScaleAutoReset,
			handleEmpty(order.ScaleInitPosition),
			handleEmpty(order.ScaleInitFillQty),
			order.ScaleRandomPercent)
	}

	if c.serverVersion >= MIN_SERVER_VER_SCALE_TABLE {
		fields = append(fields,
			order.ScaleTable,
			order.ActiveStartTime,
			order.ActiveStopTime)
	}

	// hedge orders
	if c.serverVersion >= MIN_SERVER_VER_HEDGE_ORDERS {
		fields = append(fields, order.HedgeType)
		if order.HedgeType != "" {
			fields = append(fields, order.HedgeParam)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_OPT_OUT_SMART_ROUTING {
		fields = append(fields, order.OptOutSmartRouting)
	}

	if c.serverVersion >= MIN_SERVER_VER_PTA_ORDERS {
		fields = append(fields,
			order.ClearingAccount,
			order.ClearingIntent)
	}

	if c.serverVersion >= MIN_SERVER_VER_NOT_HELD {
		fields = append(fields, order.NotHeld)
	}

	if c.serverVersion >= MIN_SERVER_VER_DELTA_NEUTRAL {
		if contract.DeltaNeutralContract != nil {
			fields = append(fields,
				true,
				contract.DeltaNeutralContract.ConID,
				contract.DeltaNeutralContract.Delta,
				contract.DeltaNeutralContract.Price)
		} else {
			fields = append(fields, false)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_ALGO_ORDERS {
		fields = append(fields, order.AlgoStrategy)

		if order.AlgoStrategy != "" {
			algoParamsCount := len(order.AlgoParams)
			fields = append(fields, algoParamsCount)
			for _, tv := range order.AlgoParams {
				fields = append(fields, tv.Tag, tv.Value)
			}
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_ALGO_ID {
		fields = append(fields, order.AlgoID)
	}

	fields = append(fields, order.WhatIf) // srv v36 and above

	// send miscOptions parameter
	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		var miscOptionsBuffer bytes.Buffer
		for _, tv := range order.OrderMiscOptions {
			miscOptionsBuffer.WriteString(tv.Tag)
			miscOptionsBuffer.WriteString("=")
			miscOptionsBuffer.WriteString(tv.Value)
			miscOptionsBuffer.WriteString(";")
		}

		fields = append(fields, miscOptionsBuffer.Bytes())
	}

	if c.serverVersion >= MIN_SERVER_VER_ORDER_SOLICITED {
		fields = append(fields, order.Solictied)
	}

	if c.serverVersion >= MIN_SERVER_VER_RANDOMIZE_SIZE_AND_PRICE {
		fields = append(fields,
			order.RandomizeSize,
			order.RandomizePrice)
	}

	if c.serverVersion >= MIN_SERVER_VER_PEGGED_TO_BENCHMARK {
		if order.OrderType == "PEG BENCH" {
			fields = append(fields,
				order.ReferenceContractID,
				order.IsPeggedChangeAmountDecrease,
				order.PeggedChangeAmount,
				order.ReferenceChangeAmount,
				order.ReferenceExchangeID)
		}

		orderConditionsCount := len(order.Conditions)
		fields = append(fields, orderConditionsCount)
		for _, cond := range order.Conditions {
			fields = append(fields, cond.Type())
			fields = append(fields, cond.makeFields()...)
		}
		if orderConditionsCount > 0 {
			fields = append(fields,
				order.ConditionsIgnoreRth,
				order.ConditionsCancelOrder)
		}

		fields = append(fields,
			order.AdjustedOrderType,
			order.TriggerPrice,
			order.LmtPriceOffset,
			order.AdjustedStopPrice,
			order.AdjustedStopLimitPrice,
			order.AdjustedTrailingAmount,
			order.AdjustableTrailingUnit)
	}
	if c.serverVersion >= MIN_SERVER_VER_EXT_OPERATOR {
		fields = append(fields, order.ExtOperator)
	}

	if c.serverVersion >= MIN_SERVER_VER_SOFT_DOLLAR_TIER {
		fields = append(fields, order.SoftDollarTier.Name, order.SoftDollarTier.Value)
	}

	if c.serverVersion >= MIN_SERVER_VER_CASH_QTY {
		fields = append(fields, order.CashQty)
	}

	if c.serverVersion >= MIN_SERVER_VER_DECISION_MAKER {
		fields = append(fields, order.Mifid2DecisionMaker, order.Mifid2DecisionAlgo)
	}

	if c.serverVersion >= MIN_SERVER_VER_MIFID_EXECUTION {
		fields = append(fields, order.Mifid2ExecutionTrader, order.Mifid2ExecutionAlgo)
	}

	if c.serverVersion >= MIN_SERVER_VER_AUTO_PRICE_FOR_HEDGE {
		fields = append(fields, order.DontUseAutoPriceForHedge)
	}

	if c.serverVersion >= MIN_SERVER_VER_ORDER_CONTAINER {
		fields = append(fields, order.IsOmsContainer)
	}

	if c.serverVersion >= MIN_SERVER_VER_D_PEG_ORDERS {
		fields = append(fields, order.DiscretionaryUpToLimitPrice)
	}

	if c.serverVersion >= MIN_SERVER_VER_PRICE_MGMT_ALGO {
		fields = append(fields, order.UsePriceMgmtAlgo)
	}

	if c.serverVersion >= MIN_SERVER_VER_DURATION {
		fields = append(fields, handleEmpty(order.Duration))
	}

	if c.serverVersion >= MIN_SERVER_VER_POST_TO_ATS {
		fields = append(fields, handleEmpty(order.PostToAts))
	}

	if c.serverVersion >= MIN_SERVER_VER_AUTO_CANCEL_PARENT {
		fields = append(fields, order.AutoCancelParent)
	}

	if c.serverVersion >= MIN_SERVER_VER_ADVANCED_ORDER_REJECT {
		fields = append(fields, order.AdvancedErrorOverride)
	}

	if c.serverVersion >= MIN_SERVER_VER_MANUAL_ORDER_TIME {
		fields = append(fields, order.ManualOrderTime)
	}

	if c.serverVersion >= MIN_SERVER_VER_PEGBEST_PEGMID_OFFSETS {
		var sendMidOffsets bool
		if contract.Exchange == "IBKRATS" {
			fields = append(fields, handleEmpty(order.MinTradeQty))
		}
		if order.OrderType == "PEG BEST" {
			fields = append(fields,
				handleEmpty(order.MinCompeteSize),
				handleEmpty(order.CompeteAgainstBestOffset))
			if order.CompeteAgainstBestOffset == COMPETE_AGAINST_BEST_OFFSET_UP_TO_MID {
				sendMidOffsets = true
			}
		} else if order.OrderType == "PEG MID" {
			sendMidOffsets = true
		}
		if sendMidOffsets {
			fields = append(fields,
				handleEmpty(order.MidOffsetAtWhole),
				handleEmpty(order.MidOffsetAtHalf))
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_CUSTOMER_ACCOUNT {
		fields = append(fields, order.CustomerAccount)
	}

	if c.serverVersion >= MIN_SERVER_VER_PROFESSIONAL_CUSTOMER {
		fields = append(fields, order.ProfessionalCustomer)
	}

	if c.serverVersion >= MIN_SERVER_VER_RFQ_FIELDS && c.serverVersion < MIN_SERVER_VER_UNDO_RFQ_FIELDS {
		fields = append(fields, "")
		fields = append(fields, UNSET_INT)
	}

	if c.serverVersion >= MIN_SERVER_VER_INCLUDE_OVERNIGHT {
		fields = append(fields, order.IncludeOvernight)
	}

	if c.serverVersion >= MIN_SERVER_VER_CME_TAGGING_FIELDS {
		fields = append(fields, order.ManualOrderIndicator)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
}

// CancelOrder cancel an order by orderId.
// It can only be used to cancel an order that was placed originally by a client with the same client ID
func (c *EClient) CancelOrder(orderID OrderID, orderCancel OrderCancel) {

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

	fields := make([]interface{}, 0, 9)
	fields = append(fields, CANCEL_ORDER)

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS {
		fields = append(fields, VERSION)
	}

	fields = append(fields, orderID)

	if c.serverVersion >= MIN_SERVER_VER_MANUAL_ORDER_TIME {
		fields = append(fields, orderCancel.ManualOrderCancelTime)
	}

	if c.serverVersion >= MIN_SERVER_VER_RFQ_FIELDS && c.serverVersion < MIN_SERVER_VER_UNDO_RFQ_FIELDS {
		fields = append(fields, "")
		fields = append(fields, "")
		fields = append(fields, UNSET_INT)
	}

	if c.serverVersion >= MIN_SERVER_VER_CME_TAGGING_FIELDS {
		fields = append(fields, orderCancel.ExtOperator)
		fields = append(fields, orderCancel.ManualOrderIndicator)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
}

// ReqOpenOrders requests the open orders that were placed from this client.
// Each open order will be fed back through the openOrder() and orderStatus() functions on the EWrapper.
// The client with a clientId of 0 will also receive the TWS-owned open orders.
// These orders will be associated with the client and a new orderId will be generated.
// This association will persist over multiple API and TWS sessions.
func (c *EClient) ReqOpenOrders() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	msg := makeFields(REQ_OPEN_ORDERS, VERSION)

	c.reqChan <- msg
}

// ReqAutoOpenOrders requests that newly created TWS orders be implicitly associated with the client.
// When a new TWS order is created, the order will be associated with the client, and fed back through the openOrder() and orderStatus() functions on the EWrapper.
// This request can only be made from a client with clientId of 0.
// if autoBind is set to TRUE, newly created TWS orders will be implicitly associated with the client.
// If set to FALSE, no association will be made.
func (c *EClient) ReqAutoOpenOrders(autoBind bool) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	msg := makeFields(REQ_AUTO_OPEN_ORDERS, VERSION, autoBind)

	c.reqChan <- msg
}

// ReqAllOpenOrders request the open orders placed from all clients and also from TWS.
// Each open order will be fed back through the openOrder() and orderStatus() functions on the EWrapper.
// No association is made between the returned orders and the requesting client.
func (c *EClient) ReqAllOpenOrders() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	msg := makeFields(REQ_ALL_OPEN_ORDERS, VERSION)

	c.reqChan <- msg
}

// ReqGlobalCancel cancels all open orders globally. It cancels both API and TWS open orders.
func (c *EClient) ReqGlobalCancel(orderCancel OrderCancel) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS && (orderCancel.ExtOperator != "" || orderCancel.ManualOrderIndicator != UNSET_INT) {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+" It does not support ext operator and manual order indicator parameters.", "")
	}

	const VERSION = 1

	fields := make([]interface{}, 0, 4)
	fields = append(fields, REQ_GLOBAL_CANCEL)

	if c.serverVersion < MIN_SERVER_VER_CME_TAGGING_FIELDS {
		fields = append(fields, VERSION)
	}

	if c.serverVersion >= MIN_SERVER_VER_CME_TAGGING_FIELDS {
		fields = append(fields, orderCancel.ExtOperator)
		fields = append(fields, orderCancel.ManualOrderIndicator)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(REQ_IDS, VERSION, numIds)

	c.reqChan <- msg
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

	msg := makeFields(REQ_ACCT_DATA, VERSION, subscribe, accountName)

	c.reqChan <- msg
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

	msg := makeFields(REQ_ACCOUNT_SUMMARY, VERSION, reqID, groupName, tags)

	c.reqChan <- msg
}

// CancelAccountSummary cancels the request for Account Window Summary tab data.
// reqId is the ID of the data request being canceled.
func (c *EClient) CancelAccountSummary(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	msg := makeFields(CANCEL_ACCOUNT_SUMMARY, VERSION, reqID)

	c.reqChan <- msg
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

	msg := makeFields(REQ_POSITIONS, VERSION)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_POSITIONS, VERSION)

	c.reqChan <- msg
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

	msg := makeFields(REQ_POSITIONS_MULTI, VERSION, reqID, account, modelCode)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_POSITIONS_MULTI, VERSION, reqID)

	c.reqChan <- msg
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

	msg := makeFields(REQ_ACCOUNT_UPDATES_MULTI, VERSION, reqID, account, modelCode, ledgerAndNLV)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_ACCOUNT_UPDATES_MULTI, VERSION, reqID)

	c.reqChan <- msg
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

	msg := makeFields(REQ_PNL, reqID, account, modelCode)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_PNL, reqID)

	c.reqChan <- msg
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

	msg := makeFields(REQ_PNL_SINGLE, reqID, account, modelCode, contractID)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_PNL_SINGLE, reqID)

	c.reqChan <- msg
}

//	##########################################################################
//	#		Executions
// 	##########################################################################

// ReqExecutions downloads the execution reports that meet the filter criteria to the client via the execDetails() function.
// To view executions beyond the past 24 hours, open the Trade Log in TWS and, while the Trade Log is displayed, request the executions again from the API.
// reqId is the ID of the data request. Ensures that responses are matched to requests if several requests are in process.
// execFilter contains attributes that describe the filter criteria used to determine which execution reports are returned.
// NOTE: Time format must be 'yyyymmdd-hh:mm:ss' Eg: '20030702-14:55'
func (c *EClient) ReqExecutions(reqID int64, execFilter ExecutionFilter) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 3

	fields := make([]interface{}, 0, 10)
	fields = append(fields, REQ_EXECUTIONS, VERSION)

	if c.serverVersion >= MIN_SERVER_VER_EXECUTION_DATA_CHAIN {
		fields = append(fields, reqID)
	}

	fields = append(fields,
		execFilter.ClientID,
		execFilter.AcctCode,
		execFilter.Time,
		execFilter.Symbol,
		execFilter.SecType,
		execFilter.Exchange,
		execFilter.Side)
	msg := makeFields(fields...)

	c.reqChan <- msg
}

//	##########################################################################
//	#		Contract Details
// 	##########################################################################

// ReqContractDetails downloads all details for a particular underlying.
// The contract details will be received via the contractDetails() function on the EWrapper.
func (c *EClient) ReqContractDetails(reqID int64, contract *Contract) {

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

	fields := make([]interface{}, 0, 21)
	fields = append(fields, REQ_CONTRACT_DATA, VERSION)

	if c.serverVersion >= MIN_SERVER_VER_CONTRACT_DATA_CHAIN {
		fields = append(fields, reqID)
	}

	fields = append(fields,
		contract.ConID, // srv v37 and above
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier) // srv v15 and above

	if c.serverVersion >= MIN_SERVER_VER_PRIMARYEXCH {
		fields = append(fields, contract.Exchange, contract.PrimaryExchange)
	} else if c.serverVersion >= MIN_SERVER_VER_LINKING {
		if contract.PrimaryExchange != "" && (contract.Exchange == "BEST" || contract.Exchange == "SMART") {
			fields = append(fields, contract.Exchange+":"+contract.PrimaryExchange)
		} else {
			fields = append(fields, contract.Exchange)
		}
	}

	fields = append(fields, contract.Currency, contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass, contract.IncludeExpired) //  srv v31 and above
	}

	if c.serverVersion >= MIN_SERVER_VER_SEC_ID_TYPE {
		fields = append(fields, contract.SecIDType, contract.SecID)
	}

	if c.serverVersion >= MIN_SERVER_VER_BOND_ISSUERID {
		fields = append(fields, contract.IssuerID)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(REQ_MKT_DEPTH_EXCHANGES)

	c.reqChan <- msg
}

// ReqMktDepth requests the market depth for a specific contract.
// The market depth will be returned by the updateMktDepth() and updateMktDepthL2() events.
// Requests the contract's market depth (order book). Note this request must be direct-routed to an exchange and not smart-routed.
// The number of simultaneous market depth requests allowed in an account is calculated based on a formula
// that looks at an accounts equity, commissions, and quote booster packs.
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

	fields := make([]interface{}, 0, 17)
	fields = append(fields, REQ_MKT_DEPTH, VERSION, reqID)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.ConID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange)

	if c.serverVersion >= MIN_SERVER_VER_MKT_DEPTH_PRIM_EXCHANGE {
		fields = append(fields, contract.PrimaryExchange)
	}

	fields = append(fields,
		contract.Currency,
		contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	fields = append(fields, numRows)

	if c.serverVersion >= MIN_SERVER_VER_SMART_DEPTH {
		fields = append(fields, isSmartDepth)
	}

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		//current doc says this part if for "internal use only" -> won't support it
		if len(mktDepthOptions) > 0 {
			log.Panic().Msg("not supported")
		}

		fields = append(fields, "")
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 4)
	fields = append(fields, CANCEL_MKT_DEPTH, VERSION, reqID)

	if c.serverVersion >= MIN_SERVER_VER_SMART_DEPTH {
		fields = append(fields, isSmartDepth)
	}
	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(REQ_NEWS_BULLETINS, VERSION, allMsgs)

	c.reqChan <- msg
}

// CancelNewsBulletins cancels the news bulletins updates
func (c *EClient) CancelNewsBulletins() {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	msg := makeFields(CANCEL_NEWS_BULLETINS, VERSION)

	c.reqChan <- msg
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

	msg := makeFields(REQ_MANAGED_ACCTS, VERSION)

	c.reqChan <- msg
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

	msg := makeFields(REQ_FA, VERSION, int(faDataType))

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 5)
	fields = append(fields,
		REPLACE_FA,
		VERSION,
		int(faDataType),
		cxml,
	)

	if c.serverVersion >= MIN_SERVER_VER_REPLACE_FA_END {
		fields = append(fields, reqID)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 30)
	fields = append(fields, REQ_HISTORICAL_DATA)

	if c.serverVersion <= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		fields = append(fields, VERSION)
	}

	fields = append(fields, reqID)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.ConID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
	)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}
	fields = append(fields,
		contract.IncludeExpired,
		endDateTime,
		barSize,
		duration,
		useRTH,
		whatToShow,
		formatDate,
	)

	if contract.SecType == "BAG" {
		fields = append(fields, len(contract.ComboLegs))
		for _, comboLeg := range contract.ComboLegs {
			fields = append(fields,
				comboLeg.ConID,
				comboLeg.Ratio,
				comboLeg.Action,
				comboLeg.Exchange,
			)
		}
	}

	if c.serverVersion >= MIN_SERVER_VER_SYNT_REALTIME_BARS {
		fields = append(fields, keepUpToDate)
	}

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		chartOptionsStr := ""
		for _, tagValue := range chartOptions {
			chartOptionsStr += tagValue.Value
		}
		fields = append(fields, chartOptionsStr)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_HISTORICAL_DATA, VERSION, reqID)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 18)

	fields = append(fields,
		REQ_HEAD_TIMESTAMP,
		reqID,
		contract.ConID,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
		contract.TradingClass,
		contract.IncludeExpired,
		useRTH,
		whatToShow,
		formatDate)

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_HEAD_TIMESTAMP, reqID)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 18)
	fields = append(fields,
		REQ_HISTOGRAM_DATA,
		reqID,
		contract.ConID,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
		contract.TradingClass,
		contract.IncludeExpired,
		useRTH,
		timePeriod)

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_HISTOGRAM_DATA, reqID)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 22)
	fields = append(fields,
		REQ_HISTORICAL_TICKS,
		reqID,
		contract.ConID,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
		contract.TradingClass,
		contract.IncludeExpired,
		startDateTime,
		endDateTime,
		numberOfTicks,
		whatToShow,
		useRTH,
		ignoreSize)

	var miscOptionsBuffer bytes.Buffer
	for _, tv := range miscOptions {
		miscOptionsBuffer.WriteString(tv.String())
	}
	fields = append(fields, miscOptionsBuffer.Bytes())

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(REQ_SCANNER_PARAMETERS, VERSION)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 25)
	fields = append(fields, REQ_SCANNER_SUBSCRIPTION)

	if c.serverVersion < MIN_SERVER_VER_SCANNER_GENERIC_OPTS {
		fields = append(fields, VERSION)
	}

	fields = append(fields,
		reqID,
		handleEmpty(subscription.NumberOfRows),
		subscription.Instrument,
		subscription.LocationCode,
		subscription.ScanCode,
		handleEmpty(subscription.AbovePrice),
		handleEmpty(subscription.BelowPrice),
		handleEmpty(subscription.AboveVolume),
		handleEmpty(subscription.MarketCapAbove),
		handleEmpty(subscription.MarketCapBelow),
		subscription.MoodyRatingAbove,
		subscription.MoodyRatingBelow,
		subscription.SpRatingAbove,
		subscription.SpRatingBelow,
		subscription.MaturityDateAbove,
		subscription.MaturityDateBelow,
		handleEmpty(subscription.CouponRateAbove),
		handleEmpty(subscription.CouponRateBelow),
		subscription.ExcludeConvertible,
		handleEmpty(subscription.AverageOptionVolumeAbove),
		subscription.ScannerSettingPairs,
		subscription.StockTypeFilter)

	if c.serverVersion >= MIN_SERVER_VER_SCANNER_GENERIC_OPTS {
		var scannerSubscriptionFilterOptionsBuffer bytes.Buffer
		for _, tv := range scannerSubscriptionFilterOptions {
			scannerSubscriptionFilterOptionsBuffer.WriteString(tv.String())
		}
		fields = append(fields, scannerSubscriptionFilterOptionsBuffer.Bytes())
	}

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		var scannerSubscriptionOptionsBuffer bytes.Buffer
		for _, tv := range scannerSubscriptionOptions {
			scannerSubscriptionOptionsBuffer.WriteString(tv.String())
		}
		fields = append(fields, scannerSubscriptionOptionsBuffer.Bytes())

	}

	msg := makeFields(fields...)

	c.reqChan <- msg
}

// CancelScannerSubscription cancel scanner.
// reqId is the unique ticker ID used for subscription.
func (c *EClient) CancelScannerSubscription(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	msg := makeFields(CANCEL_SCANNER_SUBSCRIPTION, VERSION, reqID)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 19)
	fields = append(fields, REQ_REAL_TIME_BARS, VERSION, reqID)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.ConID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecType,
		contract.LastTradeDateOrContractMonth,
		contract.Strike,
		contract.Right,
		contract.Multiplier,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.TradingClass)
	}

	fields = append(fields,
		barSize,
		whatToShow,
		useRTH)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		var realTimeBarsOptionsBuffer bytes.Buffer
		for _, tv := range realTimeBarsOptions {
			realTimeBarsOptionsBuffer.WriteString(tv.String())
		}
		fields = append(fields, realTimeBarsOptionsBuffer.Bytes())
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
}

// CancelRealTimeBars cancels realtime bars.
func (c *EClient) CancelRealTimeBars(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	const VERSION = 1

	msg := makeFields(CANCEL_REAL_TIME_BARS, VERSION, reqID)

	c.reqChan <- msg
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

	const VERSION = 2

	if c.serverVersion < MIN_SERVER_VER_FUNDAMENTAL_DATA {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support fundamental data request.", "")
		return
	}

	if c.serverVersion < MIN_SERVER_VER_TRADING_CLASS {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), UPDATE_TWS.Code, UPDATE_TWS.Msg+"  It does not support conId parameter in reqFundamentalData.", "")
		return
	}

	fields := make([]interface{}, 0, 12)
	fields = append(fields, REQ_FUNDAMENTAL_DATA, VERSION, reqID)

	if c.serverVersion >= MIN_SERVER_VER_TRADING_CLASS {
		fields = append(fields, contract.ConID)
	}

	fields = append(fields,
		contract.Symbol,
		contract.SecType,
		contract.Exchange,
		contract.PrimaryExchange,
		contract.Currency,
		contract.LocalSymbol,
		reportType)

	if c.serverVersion >= MIN_SERVER_VER_LINKING {
		var fundamentalDataOptionsBuffer bytes.Buffer
		for _, tv := range fundamentalDataOptions {
			fundamentalDataOptionsBuffer.WriteString(tv.String())
		}
		fields = append(fields, fundamentalDataOptionsBuffer.Bytes())

	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_FUNDAMENTAL_DATA, VERSION, reqID)

	c.reqChan <- msg

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

	msg := makeFields(REQ_NEWS_PROVIDERS)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 5)
	fields = append(fields,
		REQ_NEWS_ARTICLE,
		reqID,
		providerCode,
		articleID)

	if c.serverVersion >= MIN_SERVER_VER_NEWS_QUERY_ORIGINS {
		var newsArticleOptionsBuffer bytes.Buffer
		for _, tv := range newsArticleOptions {
			newsArticleOptionsBuffer.WriteString(tv.String())
		}
		fields = append(fields, newsArticleOptionsBuffer.Bytes())

	}
	msg := makeFields(fields...)

	c.reqChan <- msg
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

	fields := make([]interface{}, 0, 8)
	fields = append(fields,
		REQ_HISTORICAL_NEWS,
		reqID,
		contractID,
		providerCode,
		startDateTime,
		endDateTime,
		totalResults)

	if c.serverVersion >= MIN_SERVER_VER_NEWS_QUERY_ORIGINS {
		var historicalNewsOptionsBuffer bytes.Buffer
		for _, tv := range historicalNewsOptions {
			historicalNewsOptionsBuffer.WriteString(tv.String())
		}
		fields = append(fields, historicalNewsOptionsBuffer.Bytes())

	}
	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(QUERY_DISPLAY_GROUPS, VERSION, reqID)

	c.reqChan <- msg
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

	msg := makeFields(SUBSCRIBE_TO_GROUP_EVENTS, VERSION, reqID, groupID)

	c.reqChan <- msg
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

	msg := makeFields(UPDATE_DISPLAY_GROUP, VERSION, reqID, contractInfo)

	c.reqChan <- msg
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

	msg := makeFields(UPDATE_DISPLAY_GROUP, VERSION, reqID)

	c.reqChan <- msg
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

	msg := makeFields(VERIFY_REQUEST, VERSION, apiName, apiVersion)

	c.reqChan <- msg
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

	msg := makeFields(VERIFY_MESSAGE, VERSION, apiData)

	c.reqChan <- msg
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
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), BAD_MESSAGE.Code, BAD_MESSAGE.Msg+
			" Intent to authenticate needs to be expressed during initial connect request.", "")
		return
	}

	const VERSION = 1

	msg := makeFields(VERIFY_AND_AUTH_REQUEST, VERSION, apiName, apiVersion, opaqueIsvKey)

	c.reqChan <- msg
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

	msg := makeFields(VERIFY_MESSAGE, VERSION, apiData, xyzResponse)

	c.reqChan <- msg
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

	msg := makeFields(REQ_SEC_DEF_OPT_PARAMS, reqID, underlyingSymbol, futFopExchange, underlyingSecurityType, underlyingContractID)

	c.reqChan <- msg
}

// ReqSoftDollarTiers request pre-defined Soft Dollar Tiers.
// This is only supported for registered professional advisors and hedge and mutual funds
// who have configured Soft Dollar Tiers in Account Management.
func (c *EClient) ReqSoftDollarTiers(reqID int64) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	msg := makeFields(REQ_SOFT_DOLLAR_TIERS, reqID)

	c.reqChan <- msg
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

	msg := makeFields(REQ_FAMILY_CODES)

	c.reqChan <- msg
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

	msg := makeFields(REQ_MATCHING_SYMBOLS, reqID, pattern)

	c.reqChan <- msg
}

// ReqCompletedOrders request the completed orders.
// If apiOnly parameter is true, then only completed orders placed from API are requested.
// Result will be delivered via wrapper.CompletedOrder().
func (c *EClient) ReqCompletedOrders(apiOnly bool) {

	if !c.IsConnected() {
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), NOT_CONNECTED.Code, NOT_CONNECTED.Msg, "")
		return
	}

	msg := makeFields(REQ_COMPLETED_ORDERS, apiOnly)

	c.reqChan <- msg
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

	msg := makeFields(REQ_WSH_META_DATA, reqID)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_WSH_META_DATA, reqID)

	c.reqChan <- msg
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
	fields := make([]interface{}, 0, 10)
	fields = append(fields,
		REQ_WSH_EVENT_DATA,
		reqID,
		wshEventData.ConID)

	if c.serverVersion >= MIN_SERVER_VER_WSH_EVENT_DATA_FILTERS {
		fields = append(fields,
			wshEventData.Filter,
			wshEventData.FillWatchList,
			wshEventData.FillPortfolio,
			wshEventData.FillCompetitors)
	}

	if c.serverVersion >= MIN_SERVER_VER_WSH_EVENT_DATA_FILTERS_DATE {
		fields = append(fields,
			wshEventData.StartDate,
			wshEventData.EndDate,
			wshEventData.TotalLimit)
	}

	msg := makeFields(fields...)

	c.reqChan <- msg
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

	msg := makeFields(CANCEL_WSH_EVENT_DATA, reqID)

	c.reqChan <- msg
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

	msg := makeFields(REQ_USER_INFO, reqID)

	c.reqChan <- msg
}
