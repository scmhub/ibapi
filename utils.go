package ibapi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"time"
)

const (
	delim byte = '\x00'
	// MAX_MSG_LEN is the max length that receiver could take.
	MAX_MSG_LEN int = 0xFFFFFF // 16Mb - 1byte
	RAW_INT_LEN int = 4
)

// scanFields defines how to unpack the buf
func scanFields(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, io.EOF
	}

	if len(data) < 4 {
		return 0, nil, nil // will try to read more data
	}

	totalSize := int(binary.BigEndian.Uint32(data[:4])) + 4

	if totalSize > len(data) {
		return 0, nil, nil
	}

	// msgBytes := make([]byte, totalSize-4, totalSize-4)
	// copy(msgBytes, data[4:totalSize])
	// not copy here, copied by callee more reasonable
	return totalSize, data[4:totalSize], nil
}

// MsgBuffer is the buffer that contains a whole msg.
type MsgBuffer struct {
	bytes.Buffer
	bs  []byte
	err error
}

// NewMsgBuffer create a new MsgBuffer.
func NewMsgBuffer(bs []byte) *MsgBuffer {
	return &MsgBuffer{*bytes.NewBuffer(bs), nil, nil}
}

func (m *MsgBuffer) decode() {
	_, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode read error")
	}
}

func (m *MsgBuffer) decodeInt64() int64 {
	var i int64
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode int64 read error")
	}

	m.bs = m.bs[:len(m.bs)-1]
	if bytes.Equal(m.bs, nil) {
		return 0
	}

	i, m.err = strconv.ParseInt(string(m.bs), 10, 64)
	if m.err != nil {
		fmt.Println(string(m.bs))
		log.Panic().Err(m.err).Msg("decode int64 parse error")
	}

	return i
}

func (m *MsgBuffer) decodeRawInt64() int64 {
	// Ensure there's enough data in the buffer
	if m.Len() < RAW_INT_LEN {
		log.Panic().Err(m.err).Msg("decode raw int64 read error")
	}

	// Read 4 bytes
	buf := make([]byte, RAW_INT_LEN)
	n, err := m.Read(buf)
	if err != nil || n != RAW_INT_LEN {
		log.Panic().Err(err).Msg("decode raw int64 read error")
	}

	// Update m.bs to contain the remaining buffer
	m.bs = m.Bytes()

	// Convert directly to int64 using LittleEndian
	return int64(binary.BigEndian.Uint32(buf))
}

func (m *MsgBuffer) decodeDecimal() Decimal {
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode decimal read error")
	}

	d := StringToDecimal(string(m.bs[:len(m.bs)-1]))
	return d
}

func (m *MsgBuffer) decodeInt64ShowUnset() int64 {
	var i int64
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode int64ShowUnset read error")
	}

	m.bs = m.bs[:len(m.bs)-1]
	if bytes.Equal(m.bs, nil) {
		return UNSET_INT
	}

	i, m.err = strconv.ParseInt(string(m.bs), 10, 64)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode int64ShowUnset parse error")
	}

	return i
}

func (m *MsgBuffer) decodeFloat64() float64 {
	var f float64
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode float64 read error")
	}

	m.bs = m.bs[:len(m.bs)-1]
	if bytes.Equal(m.bs, nil) {
		return 0.0
	}

	f, m.err = strconv.ParseFloat(string(m.bs), 64)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode float64 parse error")
	}

	return f
}

func (m *MsgBuffer) decodeFloat64ShowUnset() float64 {
	var f float64
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode float64ShowUnset read error")
	}

	m.bs = m.bs[:len(m.bs)-1]
	if bytes.Equal(m.bs, nil) || bytes.Equal(m.bs, []byte("None")) {
		return UNSET_FLOAT
	}

	f, m.err = strconv.ParseFloat(string(m.bs), 64)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode float64ShowUnset parse error")
	}

	return f
}

func (m *MsgBuffer) decodeBool() bool {
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode bool read error")
	}

	m.bs = m.bs[:len(m.bs)-1]

	if bytes.Equal(m.bs, []byte{'0'}) || bytes.Equal(m.bs, nil) {
		return false
	}
	return true
}

func (m *MsgBuffer) decodeString() string {
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode string read error")
	}

	return string(m.bs[:len(m.bs)-1])
}

func (m *MsgBuffer) decodeStringUnescaped() string {
	m.bs, m.err = m.ReadBytes(delim)
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode string read error")
	}
	var s string
	s, m.err = strconv.Unquote(fmt.Sprint("\"", m.bs[:len(m.bs)-1], "\""))
	if m.err != nil {
		log.Panic().Err(m.err).Msg("decode string unmarshal error")
	}
	return s
}

// Reset reset buffer, []byte, err.
func (m *MsgBuffer) Reset() {
	m.Buffer.Reset()
	m.bs = m.bs[:0]
	m.err = nil
}

func makeField(val any) string {
	return fmt.Sprintf("%v\x00", val)
}

func splitMsgBytes(data []byte) [][]byte {
	fields := bytes.Split(data, []byte{delim})
	return fields[:len(fields)-1]
}

func stringIsEmpty(s string) bool {
	return s == ""
}

func isValidFloat64Value(val float64) bool {
	return val != UNSET_FLOAT
}

func isValidInt64Value(val int64) bool {
	return val != UNSET_INT && val != UNSET_LONG
}

func isValidDecimalValue(val Decimal) bool {
	return val != UNSET_DECIMAL
}

func FloatMaxString(val float64) string {
	if val == UNSET_FLOAT {
		return ""
	}
	return strconv.FormatFloat(val, 'g', 10, 64)
}

func LongMaxString(val int64) string {
	if val == UNSET_LONG {
		return ""
	}
	return strconv.FormatInt(val, 10)
}

func IntMaxString(val int64) string {
	if val == UNSET_INT {
		return ""
	}
	return strconv.FormatInt(val, 10)
}

func DecimalMaxString(val Decimal) string {
	if val == UNSET_DECIMAL {
		return ""
	}
	return DecimalToString(val)
}

// CurrentTimeMillis returns the current time in milliseconds.
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// GetTimeStrFromMillis converts a timestamp in milliseconds to a formatted string.
// Returns an empty string if the input time is less than or equal to zero.
func GetTimeStrFromMillis(timestamp int64) string {
	if timestamp > 0 {
		return time.Unix(0, timestamp*int64(time.Millisecond)).Format("20060102-15:04:05")
	}
	return ""
}

// isASCIIPrintable checks if all characters in the given string are ASCII printable characters.
func isASCIIPrintable(s string) bool {
	for _, r := range s {
		if r == '\t' || r == '\n' || r == '\r' || r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func BoolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
