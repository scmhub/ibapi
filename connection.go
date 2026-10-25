package ibapi

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

// Connection is a TCPConn wrapper with lock-free statistics and minimal contention.
type Connection struct {
	// Connection state - protected by mutex for host/port coordination only
	mu          sync.RWMutex
	tcpConn     atomic.Pointer[net.TCPConn] // Lock-free pointer for maximum performance
	wrapper     EWrapper
	host        string
	port        int
	isConnected int32 // atomic: 0=disconnected, 1=connected

	// Statistics - lock-free atomic counters for maximum performance
	numBytesSent int64 // atomic
	numMsgSent   int64 // atomic
	numBytesRecv int64 // atomic
	numMsgRecv   int64 // atomic
}

func (c *Connection) Write(bs []byte) (int, error) {
	conn := c.getConn()
	if conn == nil {
		return 0, fmt.Errorf("connection not available")
	}
	n, err := conn.Write(bs)
	if err == nil {
		atomic.AddInt64(&c.numBytesSent, int64(n))
		atomic.AddInt64(&c.numMsgSent, 1)
	}
	// A TCP reconnect cannot restore an API session. Its owner must create
	// a fresh client and perform the complete API handshake before retrying.
	return n, err
}

func (c *Connection) Read(bs []byte) (int, error) {
	conn := c.getConn()
	if conn == nil {
		return 0, fmt.Errorf("connection not available")
	}

	n, err := conn.Read(bs)

	// Lock-free atomic statistics update
	atomic.AddInt64(&c.numBytesRecv, int64(n))
	atomic.AddInt64(&c.numMsgRecv, 1)

	log.Trace().Int("nBytes", n).Msg("conn read")

	return n, err
}

// getConn returns the current TCP connection in a lock-free way
func (c *Connection) getConn() *net.TCPConn {
	return c.tcpConn.Load()
}

// setConn swaps the TCP connection and returns the previous value.
func (c *Connection) setConn(conn *net.TCPConn) *net.TCPConn {
	return c.tcpConn.Swap(conn)
}

func (c *Connection) reset() {
	// Lock-free atomic reset of statistics
	atomic.StoreInt64(&c.numBytesSent, 0)
	atomic.StoreInt64(&c.numBytesRecv, 0)
	atomic.StoreInt64(&c.numMsgSent, 0)
	atomic.StoreInt64(&c.numMsgRecv, 0)
}

func (c *Connection) connectContext(ctx context.Context, host string, port int) error {
	// Protect host/port assignment with mutex to prevent races
	c.mu.Lock()
	c.host = host
	c.port = port
	c.mu.Unlock()

	c.reset()

	// Use the parameters directly instead of reading from struct to avoid races
	address := fmt.Sprintf("%v:%v", host, port)
	dialed, err := (&net.Dialer{}).DialContext(ctx, "tcp4", address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("failed to dial tcp")
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), FAIL_CREATE_SOCK.Code, FAIL_CREATE_SOCK.Msg, "")
		return err
	}

	newConn := dialed.(*net.TCPConn)

	// Atomically update connection state (close previous socket if necessary)
	if oldConn := c.setConn(newConn); oldConn != nil {
		_ = oldConn.Close()
	}
	atomic.StoreInt32(&c.isConnected, 1)

	log.Debug().Any("address", newConn.RemoteAddr()).Msg("tcp socket connected")
	return nil
}

func (c *Connection) disconnect() error {
	// Load statistics atomically for logging
	msgSent := atomic.LoadInt64(&c.numMsgSent)
	bytesSent := atomic.LoadInt64(&c.numBytesSent)
	msgRecv := atomic.LoadInt64(&c.numMsgRecv)
	bytesRecv := atomic.LoadInt64(&c.numBytesRecv)

	log.Trace().
		Int64("nMsgSent", msgSent).Int64("nBytesSent", bytesSent).
		Int64("nMsgRecv", msgRecv).Int64("nBytesRecv", bytesRecv).
		Msg("conn disconnect")

	// Atomically mark as disconnected
	atomic.StoreInt32(&c.isConnected, 0)

	// Close the connection
	if conn := c.setConn(nil); conn != nil {
		return conn.Close()
	}
	return nil
}

func (c *Connection) IsConnected() bool {
	return atomic.LoadInt32(&c.isConnected) == 1
}

// GetStatistics returns current connection statistics atomically
func (c *Connection) GetStatistics() (bytesSent, msgSent, bytesRecv, msgRecv int64) {
	return atomic.LoadInt64(&c.numBytesSent),
		atomic.LoadInt64(&c.numMsgSent),
		atomic.LoadInt64(&c.numBytesRecv),
		atomic.LoadInt64(&c.numMsgRecv)
}
