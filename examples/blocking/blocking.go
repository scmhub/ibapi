package main

import (
	"math/rand"
	"time"

	"github.com/scmhub/ibapi"
)

const (
	host = "localhost"
	port = 7497
)

var log = ibapi.Logger()

var tChan chan int64

func init() {
	tChan = make(chan int64)
}

// Wrapper
type Wrapper struct {
	ibapi.Wrapper
}

func (w Wrapper) CurrentTime(t int64) {
	log.Debug().Time("Server Time", time.Unix(t, 0)).Msg("<CurrentTime>")
	tChan <- t
}

// IB
type IB struct {
	ibapi.EClient
}

func NewIB() *IB {
	return &IB{
		EClient: *ibapi.NewEClient(&Wrapper{}),
	}
}

func (ib *IB) ReqCurrentTime() int64 {
	ib.EClient.ReqCurrentTime()
	return <-tChan
}

func main() {
	// Set the console writter
	ibapi.SetConsoleWriter()
	// Set log level
	//ibapi.SetLogLevel(int(zerolog.DebugLevel))

	// Creates IB CLient
	ib := NewIB()

	// Client connection
	err := ib.Connect(host, port, int64(rand.Intn(1e9)))
	if err != nil {
		log.Error().Err(err).Msg("Connect")
		return
	}
	defer ib.Disconnect()
	// Add a short delay to allow the connection to stabilize
	time.Sleep(100 * time.Millisecond)
	log.Info().Msg("Waited for connection to stabilize")

	// Request servert current time
	t := ib.ReqCurrentTime()
	log.Info().Time("current time", time.Unix(t, 0)).Msg("Requested Server Current time")
}
