package sandbox

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const delim byte = '\x00'

// StartMockTWSServer starts a minimal mock TWS-like server that speaks enough of the
// handshake so tests can run without a real TWS/Gateway. It accepts multiple
// connections and, per-connection, handles:
// 1) Handshake header from client ("API\x00" + len + version string)
// 2) Responds with one server info line: "<serverVersion>;<connTime>\0"
// 3) Reads one framed message (START_API) and then drains/ignores further frames
func StartMockTWSServer(host string, port int) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handleMockConn(conn)
		}
	}()
	return nil
}

func handleMockConn(c net.Conn) {
	br := bufio.NewReader(c)
	bw := bufio.NewWriter(c)

	// Expect handshake magic "API\x00"
	magic := make([]byte, 4)
	if _, err := io.ReadFull(br, magic); err != nil {
		_ = c.Close()
		return
	}
	if !bytes.Equal(magic, []byte{'A', 'P', 'I', 0x00}) {
		_ = c.Close()
		return
	}

	// Read version string length (4 bytes BE) then the version string
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(br, lenBuf); err != nil {
		_ = c.Close()
		return
	}
	n := int(binary.BigEndian.Uint32(lenBuf))
	if n < 0 || n > 1<<20 {
		_ = c.Close()
		return
	}
	versionStr := make([]byte, n)
	if _, err := io.ReadFull(br, versionStr); err != nil {
		_ = c.Close()
		return
	}
	_ = versionStr // not used; just consumed

	// Write server info as a framed message expected by scanFields:
	// 4-byte big-endian length + payload "<serverVersion>;<connTime>\0"
	// Use serverVersion 200 (<201) to keep legacy framing but >=197 to support CurrentTimeInMillis
	// Handshake expects two fields separated by delim: "version\0connTime\0"
	payload := []byte("200" + string(delim) + "20250101 00:00:00" + string(delim))
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(payload)))
	if _, err := bw.Write(size); err != nil {
		_ = c.Close()
		return
	}
	if _, err := bw.Write(payload); err != nil {
		_ = c.Close()
		return
	}
	if err := bw.Flush(); err != nil {
		_ = c.Close()
		return
	}

	// Read one framed START_API message (len + payload)
	hdr := make([]byte, 4)
	if _, err := io.ReadFull(br, hdr); err != nil {
		_ = c.Close()
		return
	}
	l := int(binary.BigEndian.Uint32(hdr))
	if l < 0 || l > 1<<20 {
		_ = c.Close()
		return
	}
	startAPIPayload := make([]byte, l)
	if _, err := io.ReadFull(br, startAPIPayload); err != nil {
		_ = c.Close()
		return
	}

	// Immediately emit a couple of legacy responses to satisfy time-based tests
	// Helper to write a legacy framed message (fields joined by delim, ending with delim)
	writeLegacy := func(fields ...string) bool {
		data := strings.Join(fields, string(delim)) + string(delim)
		sz := make([]byte, 4)
		binary.BigEndian.PutUint32(sz, uint32(len(data)))
		if _, err := bw.Write(sz); err != nil {
			return false
		}
		if _, err := bw.WriteString(data); err != nil {
			return false
		}
		if err := bw.Flush(); err != nil {
			return false
		}
		return true
	}

	// CURRENT_TIME
	now := time.Now().Unix()
	_ = writeLegacy("49", "1", strconv.FormatInt(now, 10))

	// CURRENT_TIME_IN_MILLIS
	nowMs := time.Now().UnixNano() / int64(time.Millisecond)
	_ = writeLegacy("109", strconv.FormatInt(nowMs, 10))

	// Main loop: parse outbound requests (legacy) and emit minimal legacy responses
	for {
		if _, err := io.ReadFull(br, hdr); err != nil {
			_ = c.Close()
			return
		}
		l = int(binary.BigEndian.Uint32(hdr))
		if l < 0 || l > 1<<20 {
			_ = c.Close()
			return
		}
		req := make([]byte, l)
		if _, err := io.ReadFull(br, req); err != nil {
			_ = c.Close()
			return
		}

		// Legacy requests are delimited by `delim` (0x00). The first field is OUT msg id.
		parts := bytes.Split(req, []byte{delim})
		if len(parts) == 0 || len(parts[0]) == 0 {
			continue
		}
		msgID := string(parts[0])
		switch msgID {
		case "49": // REQ_CURRENT_TIME
			now := time.Now().Unix()
			_ = writeLegacy("49", "1", strconv.FormatInt(now, 10))
		case "105": // REQ_CURRENT_TIME_IN_MILLIS
			nowMs := time.Now().UnixNano() / int64(time.Millisecond)
			_ = writeLegacy("109", strconv.FormatInt(nowMs, 10))
		case "50": // REQ_REAL_TIME_BARS (ignored)
			// Could emit a synthetic REAL_TIME_BARS event if needed
		case "4": // REQ_MKT_DATA
			// Emit a couple of ticks: TICK_PRICE (msg 1) and TICK_SIZE (msg 2)
			// TickPrice legacy: msgID, version, reqID, tickType, price, size, attrMask
			_ = writeLegacy("1", "6", "1001", "1", "150.0", "0", "0")
			// TickSize legacy: msgID, version, reqID, sizeTickType, size
			_ = writeLegacy("2", "6", "1001", "0", "100")
		case "62": // REQ_MARKET_DATA_TYPE (ignored)
		case "21": // REQ_ACCOUNT_UPDATES (subscribe/unsubscribe)
			// ACCT_VALUE: version, tag, value, currency, account
			_ = writeLegacy("6", "1", "NetLiquidation", "100000", "USD", "DU000000")
			// PORTFOLIO_VALUE (v=7): version, conid, symbol, secType, lastTradeDateOrContractMonth, strike, right,
			// multiplier, primaryExch, currency, localSymbol, tradingClass, position, mktPrice, mktValue,
			// averageCost, unrealizedPNL, realizedPNL, accountName
			_ = writeLegacy(
				"7", "7",
				"265598", "AAPL", "STK", "", "0", "",
				"", "SMART", "USD", "AAPL", "NMS",
				"10", "190.0", "1900.0", "180.0", "100.0", "0.0", "DU000000",
			)
			// ACCT_UPDATE_TIME: version, time
			_ = writeLegacy("8", "1", "13:37")
			// ACCT_DOWNLOAD_END: version, account
			_ = writeLegacy("54", "1", "DU000000")
		case "56": // REQ_ACCOUNT_SUMMARY
			_ = writeLegacy("63", "1", "9001", "DU000000", "NetLiquidation", "100000", "USD") // ACCOUNT_SUMMARY
			_ = writeLegacy("64", "1", "9001")                                                // ACCOUNT_SUMMARY_END
		case "57": // CANCEL_ACCOUNT_SUMMARY
			// no-op
		case "104": // REQ_USER_INFO
			_ = writeLegacy("107", "1", "brand") // USER_INFO
		case "73": // REQ_ACCOUNT_UPDATES_MULTI
			// ACCOUNT_UPDATE_MULTI: version, reqId, account, modelCode, key, value, currency
			_ = writeLegacy("73", "1", "9005", "DU000000", "EUstocks", "NetLiquidation", "100000", "USD")
			// ACCOUNT_UPDATE_MULTI_END: version, reqId
			_ = writeLegacy("74", "1", "9005")
		case "74": // REQ_POSITIONS_MULTI
			// POSITION_MULTI: version, reqId, account, then contract fields, position, avgCost, modelCode
			_ = writeLegacy(
				"71", "1", "9006", "DU000000",
				"265598", "AAPL", "STK", "", "0", "", "SMART", "USD", "AAPL", "NMS", "NMS",
				"10", "180.0", "EUstocks",
			)
			// POSITION_MULTI_END: version, reqId
			_ = writeLegacy("72", "1", "9006")
		case "24": // REQ_SCANNER_PARAMETERS
			// Return a longer, safe XML to avoid wrapper slicing panic
			scannerXML := "<ScannerParameters>" +
				"<ScanType>TOP_PERC_GAIN</ScanType>" +
				"<Instruments>STK</Instruments>" +
				"<Location>STK.US.MAJOR</Location>" +
				"<Columns><Column>price</Column><Column>volume</Column></Columns>" +
				"</ScannerParameters>"
			_ = writeLegacy("19", "1", scannerXML)
		case "22": // REQ_SCANNER_SUBSCRIPTION
			// Emit end right away: SCANNER_DATA (IN=20) with zero elements
			_ = writeLegacy("20", "1", "7001", "0")
		case "79": // REQ_SOFT_DOLLAR_TIERS
			// SOFT_DOLLAR_TIERS: reqId, tiersCount, then name/value/displayName per tier
			// Return zero tiers to avoid extra parsing complexity in tests
			_ = writeLegacy("77", "4001", "0")
		case "25": // CANCEL_HISTORICAL_DATA (ignored)
		case "20": // REQ_POSITIONS
			// POSITION_DATA: version, account, full contract fields, position (decimal), avgCost (float)
			_ = writeLegacy(
				"61", "3",
				"DU000000",
				"265598", "AAPL", "STK", "", "0", "", "", "USD", "AAPL", "NMS", "NMS",
				"10", "180.0",
			)
			_ = writeLegacy("62", "1")
		case "15": // REQ_HISTORICAL_DATA
			// HISTORICAL_DATA legacy format used here:
			// serverVersion 200 >= SYNT_REALTIME_BARS (124), so no version field inside message
			// reqID, [no start/end fields], itemCount, then bar fields
			_ = writeLegacy("17", "4001", "1", "20250101 00:00:00", "100", "101", "99", "100.5", "1000", "100.2", "1")
			_ = writeLegacy("108", "4001", "20250101 00:00:00", "20250101 00:00:00")
		default:
			// Ignore unsupported requests
		}
	}
}
