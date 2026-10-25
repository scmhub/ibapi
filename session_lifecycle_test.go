package ibapi

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func lifecycleReadFrame(conn net.Conn) error {
	var header [4]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		return err
	}
	_, err := io.CopyN(io.Discard, conn, int64(binary.BigEndian.Uint32(header[:])))
	return err
}

func lifecycleWriteFrame(conn net.Conn, payload string) error {
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))
	if _, err := conn.Write(header[:]); err != nil {
		return err
	}
	_, err := io.WriteString(conn, payload)
	return err
}

// Only the initial version negotiation and StartAPI exchange are needed to
// reproduce shutdown. The fixture implements no order or account operations.
func lifecycleGateway(t *testing.T, handshakeReply string) (string, int, <-chan error, <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	result := make(chan error, 1)
	reached := make(chan struct{})
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			result <- err
			return
		}
		defer func() { _ = conn.Close() }()
		// Safety net only: normal completion is synchronized by the socket EOF.
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		var prefix [4]byte
		if _, err = io.ReadFull(conn, prefix[:]); err != nil {
			result <- err
			return
		}
		if string(prefix[:]) != "API\x00" {
			result <- io.ErrUnexpectedEOF
			return
		}
		if err = lifecycleReadFrame(conn); err != nil {
			result <- err
			return
		}
		close(reached)
		if handshakeReply == "blocked" {
			_, err = io.Copy(io.Discard, conn)
			result <- err
			return
		}
		if handshakeReply == "" {
			result <- nil
			return
		}
		if err = lifecycleWriteFrame(conn, handshakeReply); err != nil {
			result <- err
			return
		}
		if err = lifecycleReadFrame(conn); err != nil {
			result <- err
			return
		}
		_, err = io.Copy(io.Discard, conn)
		result <- err
	}()
	host, portText, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatal(err)
	}
	return host, port, result, reached
}

func TestLifecycleFreshSessionDisposal(t *testing.T) {
	for range 20 {
		host, port, finished, _ := lifecycleGateway(t, "180\x0020261026 09:30:00 EST\x00")
		client := NewEClient(&Wrapper{})
		if err := client.Connect(host, port, 15); err != nil {
			t.Fatal(err)
		}
		if err := client.Disconnect(); err != nil {
			t.Fatal(err)
		}
		if err := <-finished; err != nil {
			t.Fatal(err)
		}
	}
}

func TestLifecycleHandshakeEOFIsError(t *testing.T) {
	host, port, finished, _ := lifecycleGateway(t, "")
	client := NewEClient(&Wrapper{})
	if err := client.Connect(host, port, 15); err == nil {
		t.Error("Connect accepted EOF before server version negotiation")
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleHandshakeCancellationClosesSocket(t *testing.T) {
	host, port, finished, reached := lifecycleGateway(t, "blocked")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := NewEClient(&Wrapper{})
	connected := make(chan error, 1)
	go func() { connected <- client.ConnectContext(ctx, host, port, 15) }()
	select {
	case <-reached:
	case <-time.After(5 * time.Second):
		t.Fatal("client did not reach handshake")
	}
	cancel()
	select {
	case err := <-connected:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("handshake ignored context cancellation")
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleCanceledDial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := NewEClient(&Wrapper{})
	if err := client.ConnectContext(ctx, "127.0.0.1", 1, 15); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestLifecycleRejectsClientReuse(t *testing.T) {
	host, port, finished, _ := lifecycleGateway(t, "180\x0020261026 09:30:00 EST\x00")
	client := NewEClient(&Wrapper{})
	if err := client.Connect(host, port, 15); err != nil {
		t.Fatal(err)
	}
	if err := client.Disconnect(); err != nil {
		t.Fatal(err)
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
	if err := client.Connect(host, port, 15); err == nil {
		t.Fatal("reused a disposed session")
	}
}

type lifecycleWriteFailure struct{}

func (lifecycleWriteFailure) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestLifecycleRequesterWriteFailureDoesNotJoinItself(t *testing.T) {
	client := NewEClient(&Wrapper{})
	client.writer = bufio.NewWriterSize(lifecycleWriteFailure{}, 1)
	atomic.StoreInt32(&client.conn.isConnected, 1)
	client.setConnState(CONNECTED)
	exited := make(chan struct{})
	client.wg.Add(1)
	go func() { client.request(); close(exited) }()
	client.reqChan <- []byte("request larger than the writer buffer")
	select {
	case <-exited:
	case <-time.After(time.Second):
		t.Fatal("requester tried to join itself during shutdown")
	}
	if client.ctx.Err() == nil {
		t.Fatal("writer failure did not cancel session")
	}
	if err := client.Disconnect(); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleWriteDoesNotReconnectRawTCP(t *testing.T) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	host, portText, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatal(err)
	}
	conn := &Connection{wrapper: &Wrapper{}, host: host, port: port}
	defer func() { _ = conn.disconnect() }()
	if _, err := conn.Write([]byte("must not start a raw connection")); err == nil {
		t.Fatal("Write created a new TCP connection without an API handshake")
	}
	if conn.getConn() != nil {
		t.Fatal("Write changed the session socket")
	}
}

type lifecycleWrapper struct {
	Wrapper
	closed func()
}

func (w *lifecycleWrapper) ConnectionClosed() { w.closed() }

func TestLifecycleConcurrentAndReentrantDisconnect(t *testing.T) {
	host, port, finished, _ := lifecycleGateway(t, "180\x0020261026 09:30:00 EST\x00")
	var calls atomic.Int32
	wrapper := &lifecycleWrapper{}
	client := NewEClient(wrapper)
	wrapper.closed = func() {
		calls.Add(1)
		if err := client.Disconnect(); err != nil {
			t.Error(err)
		}
	}
	if err := client.Connect(host, port, 15); err != nil {
		t.Fatal(err)
	}
	var callers sync.WaitGroup
	for range 20 {
		callers.Go(func() {
			if err := client.Disconnect(); err != nil {
				t.Error(err)
			}
		})
	}
	callers.Wait()
	if calls.Load() != 1 {
		t.Fatalf("ConnectionClosed called %d times", calls.Load())
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleRejectsConnectAfterEarlyDisconnect(t *testing.T) {
	client := NewEClient(&Wrapper{})
	if err := client.Disconnect(); err != nil {
		t.Fatal(err)
	}
	if err := client.Connect("127.0.0.1", 1, 15); err == nil {
		t.Fatal("started a connection after client disposal")
	}
}

func TestLifecycleProtobufSessionDisposal(t *testing.T) {
	reply := fmt.Sprintf("%d\x0020261026 09:30:00 EST\x00", MAX_CLIENT_VER)
	host, port, finished, _ := lifecycleGateway(t, reply)
	client := NewEClient(&Wrapper{})
	if err := client.Connect(host, port, 15); err != nil {
		t.Fatal(err)
	}
	if err := client.Disconnect(); err != nil {
		t.Fatal(err)
	}
	if err := <-finished; err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleRejectsIncompleteHandshake(t *testing.T) {
	host, port, finished, _ := lifecycleGateway(t, "180\x00")
	client := NewEClient(&Wrapper{})
	if err := client.Connect(host, port, 15); err == nil {
		t.Fatal("accepted a handshake without a server connection time")
	}
	if err := <-finished; err != nil && !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}
}
