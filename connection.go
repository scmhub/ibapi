package ibapi

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxReconnectAttempts = 3
	reconnectDelay       = 500 * time.Millisecond
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

	// Reconnection control - prevents multiple concurrent reconnections
	reconnecting int32 // atomic: 0=not reconnecting, 1=reconnecting
}

func (c *Connection) Write(bs []byte) (int, error) {
	// Fast path: try write with current connection
	conn := c.getConn()
	if conn != nil {
		n, err := conn.Write(bs)
		if err == nil {
			// Lock-free atomic statistics update
			atomic.AddInt64(&c.numBytesSent, int64(n))
			atomic.AddInt64(&c.numMsgSent, 1)
			log.Trace().Int("nBytes", n).Msg("conn write")
			return n, nil
		}

		// Write failed, try to reconnect
		log.Warn().Err(err).Msg("Write error detected, attempting to reconnect...")
	}

	// Slow path: reconnect and retry
	if err := c.reconnect(); err != nil {
		return 0, fmt.Errorf("write failed and reconnection failed: %w", err)
	}

	// Retry write after reconnection
	conn = c.getConn()
	if conn == nil {
		return 0, fmt.Errorf("connection still not available after reconnect")
	}

	n, err := conn.Write(bs)
	if err != nil {
		return 0, fmt.Errorf("write retry after reconnect failed: %w", err)
	}

	// Lock-free atomic statistics update
	atomic.AddInt64(&c.numBytesSent, int64(n))
	atomic.AddInt64(&c.numMsgSent, 1)
	log.Trace().Int("nBytes", n).Msg("conn write (after reconnect)")
	return n, nil
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

// setConn sets the TCP connection in a lock-free way
func (c *Connection) setConn(conn *net.TCPConn) {
	c.tcpConn.Store(conn)
}

func (c *Connection) reset() {
	// Lock-free atomic reset of statistics
	atomic.StoreInt64(&c.numBytesSent, 0)
	atomic.StoreInt64(&c.numBytesRecv, 0)
	atomic.StoreInt64(&c.numMsgSent, 0)
	atomic.StoreInt64(&c.numMsgRecv, 0)
}

func (c *Connection) connect(host string, port int) error {
	// Protect host/port assignment with mutex to prevent races
	c.mu.Lock()
	c.host = host
	c.port = port
	c.mu.Unlock()

	c.reset()

	// Use the parameters directly instead of reading from struct to avoid races
	address := fmt.Sprintf("%v:%v", host, port)
	addr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Error().Err(err).Str("host", address).Msg("failed to resove tcp address")
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), FAIL_CREATE_SOCK.Code, FAIL_CREATE_SOCK.Msg, "")
		return err
	}

	newConn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		log.Error().Err(err).Any("address", addr).Msg("failed to dial tcp")
		c.wrapper.Error(NO_VALID_ID, currentTimeMillis(), FAIL_CREATE_SOCK.Code, FAIL_CREATE_SOCK.Msg, "")
		return err
	}

	// Atomically update connection state
	c.setConn(newConn)
	atomic.StoreInt32(&c.isConnected, 1)

	log.Debug().Any("address", newConn.RemoteAddr()).Msg("tcp socket connected")
	return nil
}

func (c *Connection) reconnect() error {
	// Use atomic CAS to prevent multiple concurrent reconnections
	if !atomic.CompareAndSwapInt32(&c.reconnecting, 0, 1) {
		// Another goroutine is already reconnecting, wait for it
		for atomic.LoadInt32(&c.reconnecting) == 1 {
			time.Sleep(10 * time.Millisecond)
		}
		// Check if the other goroutine succeeded
		if atomic.LoadInt32(&c.isConnected) == 1 {
			return nil
		}
		return fmt.Errorf("concurrent reconnection failed")
	}

	// Ensure we clear the reconnecting flag when done
	defer atomic.StoreInt32(&c.reconnecting, 0)

	var err error
	backoff := reconnectDelay // Start with base delay

	// try up to maxReconnectAttempts times
	for attempt := 1; attempt <= maxReconnectAttempts; attempt++ {
		log.Info().
			Int("attempt", attempt).
			Int("maxAttempts", maxReconnectAttempts).
			Msg("Attempting to reconnect")

		// Read host/port atomically to avoid race
		c.mu.RLock()
		host, port := c.host, c.port
		c.mu.RUnlock()

		err = c.connect(host, port)
		if err == nil {
			log.Info().Msg("Reconnection successful")
			atomic.StoreInt32(&c.isConnected, 1)
			return nil
		}

		// if this isn't our last try, wait and then loop again
		if attempt < maxReconnectAttempts {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	// if we get here, all attempts failed
	atomic.StoreInt32(&c.isConnected, 0)
	return fmt.Errorf("failed to reconnect after %d attempts: %w", maxReconnectAttempts, err)
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
	conn := c.getConn()
	if conn != nil {
		c.setConn(nil)
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
