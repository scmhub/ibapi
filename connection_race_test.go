package ibapi

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// DummyServer creates a simple TCP echo server for testing
type DummyServer struct {
	listener net.Listener
	addr     string
	port     int
}

func NewDummyServer() (*DummyServer, error) {
	// Listen on a random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	addr := listener.Addr().(*net.TCPAddr)
	server := &DummyServer{
		listener: listener,
		addr:     addr.IP.String(),
		port:     addr.Port,
	}

	return server, nil
}

func (s *DummyServer) Start() {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				return // Server stopped
			}
			go s.handleConnection(conn)
		}
	}()
}

func (s *DummyServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Simple echo server - read and write back data
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Server read error: %v\n", err)
			}
			return
		}

		// Echo the data back
		_, writeErr := conn.Write(buffer[:n])
		if writeErr != nil {
			fmt.Printf("Server write error: %v\n", writeErr)
			return
		}
	}
}

func (s *DummyServer) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *DummyServer) Address() (string, int) {
	return s.addr, s.port
}

// TestConnectionRaceConditions demonstrates race conditions in Connection
func TestConnectionRaceConditions(t *testing.T) {
	// Start dummy server
	server, err := NewDummyServer()
	if err != nil {
		t.Fatalf("Failed to create dummy server: %v", err)
	}
	defer server.Stop()

	server.Start()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	host, port := server.Address()

	// Create connection with a simple wrapper
	wrapper := &Wrapper{}
	conn := &Connection{
		wrapper: wrapper,
	}

	// Connect to dummy server
	err = conn.connect(host, port)
	if err != nil {
		t.Fatalf("Failed to connect to dummy server: %v", err)
	}

	// Test data
	testData := []byte("Hello, Race Condition Test!")
	numOperations := 100
	numGoroutines := 10

	var wg sync.WaitGroup

	// Start multiple writer goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				// This will race on numBytesSent and numMsgSent
				_, err := conn.Write(testData)
				if err != nil {
					// Expected during disconnect
					return
				}

				// Small delay to increase chance of race
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Start multiple reader goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			buffer := make([]byte, 1024)
			for j := 0; j < numOperations; j++ {
				// This will race on numBytesRecv and numMsgRecv
				_, err := conn.Read(buffer)
				if err != nil {
					// Expected during disconnect
					return
				}

				// Small delay to increase chance of race
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Start disconnect goroutines to trigger races
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Wait a bit then disconnect
			time.Sleep(time.Duration(goroutineID*10) * time.Millisecond)

			// This will race with the statistics updates and isConnected flag
			err := conn.disconnect()
			if err != nil {
				// Multiple disconnects expected to fail
			}
		}(i)
	}

	// Start statistics reader goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations*2; j++ {
				// Using atomic-safe methods to read statistics
				bytesSent, msgSent, bytesRecv, msgRecv := conn.GetStatistics()
				_ = bytesSent + msgSent + bytesRecv + msgRecv
				_ = conn.IsConnected()

				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Start reset goroutines
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			time.Sleep(time.Duration(goroutineID*20) * time.Millisecond)

			// This will race with ongoing statistics updates
			conn.reset()
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Final disconnect to clean up
	conn.disconnect()

	t.Logf("Test completed - check with 'go test -race' to detect race conditions")
	bytesSent, msgSent, bytesRecv, msgRecv := conn.GetStatistics()
	t.Logf("Final stats - Sent: %d msgs, %d bytes | Recv: %d msgs, %d bytes",
		msgSent, bytesSent, msgRecv, bytesRecv)
}

// TestConnectionConcurrentReconnect tests the reconnection logic under concurrent access
func TestConnectionConcurrentReconnect(t *testing.T) {
	// Start dummy server
	server, err := NewDummyServer()
	if err != nil {
		t.Fatalf("Failed to create dummy server: %v", err)
	}
	defer server.Stop()

	server.Start()
	time.Sleep(100 * time.Millisecond)

	host, port := server.Address()

	wrapper := &Wrapper{}
	conn := &Connection{
		wrapper: wrapper,
	}

	var wg sync.WaitGroup
	numGoroutines := 5

	// Multiple goroutines trying to write (which triggers reconnect on failure)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// First establish connection
			if err := conn.connect(host, port); err != nil {
				t.Errorf("Goroutine %d failed to connect: %v", id, err)
				return
			}

			// Try multiple writes - some may trigger reconnection
			for j := 0; j < 50; j++ {
				data := []byte(fmt.Sprintf("Message from goroutine %d, iteration %d", id, j))

				// This can race with other goroutines doing connect/disconnect/reconnect
				_, err := conn.Write(data)
				if err != nil {
					// Expected during concurrent access
				}

				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	// Goroutine that disconnects periodically
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			conn.disconnect() // Race with Write operations
		}
	}()

	wg.Wait()

	t.Logf("Concurrent reconnect test completed")
}

// TestConnectionStatisticsRace focuses specifically on the statistics counter races
func TestConnectionStatisticsRace(t *testing.T) {
	wrapper := &Wrapper{}
	conn := &Connection{
		wrapper: wrapper,
	}

	// Don't actually connect - just test the statistics
	// Simulate concurrent access to the counters

	var wg sync.WaitGroup
	iterations := 1000

	// Goroutines incrementing send stats using atomic operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				atomic.AddInt64(&conn.numBytesSent, 100) // Now atomic!
				atomic.AddInt64(&conn.numMsgSent, 1)     // Now atomic!
			}
		}()
	}

	// Goroutines incrementing recv stats using atomic operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				atomic.AddInt64(&conn.numBytesRecv, 50) // Now atomic!
				atomic.AddInt64(&conn.numMsgRecv, 1)    // Now atomic!
			}
		}()
	}

	// Goroutines reading stats using atomic operations
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations*2; j++ {
				// Reading using atomic-safe methods - no more race condition!
				bytesSent, msgSent, bytesRecv, msgRecv := conn.GetStatistics()
				_ = bytesSent + bytesRecv + msgSent + msgRecv
			}
		}()
	}

	// Goroutines resetting stats
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				time.Sleep(time.Millisecond)
				conn.reset() // Race condition with increments!
			}
		}()
	}

	wg.Wait()

	bytesSent, msgSent, bytesRecv, msgRecv := conn.GetStatistics()
	t.Logf("Final statistics after race: Sent=%d/%d, Recv=%d/%d",
		msgSent, bytesSent, msgRecv, bytesRecv)

	// Note: The final values will be unpredictable due to race conditions
	// Expected: 10 goroutines * 1000 iterations = 10,000 messages
	// Actual: Will be less due to lost updates from race conditions
}
