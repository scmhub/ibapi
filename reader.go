package ibapi

import (
	"bufio"
	"context"
	"fmt"
	"sync"
)

// EReader starts the scan and decode goroutines
func EReader(ctx context.Context, scanner *bufio.Scanner, decoder *EDecoder, wg *sync.WaitGroup) {

	msgChan := make(chan []byte, 300)

	// Decode
	wg.Add(1)
	go func() {
		log.Debug().Msg("Decoder started")
		defer log.Debug().Msg("Decoder ended")
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgChan:
				if !ok {
					fmt.Println("msgChan closed, exiting decoder")
					return
				}
				decoder.interpret(msg) // single worker and no go here!!
			}
		}
	}()

	// Scan
	wg.Add(1)
	go func() {
		log.Debug().Msg("scanner started")
		defer log.Debug().Msg("scanner ended")
		defer wg.Done()
		for scanner.Scan() {
			msgBytes := make([]byte, len(scanner.Bytes()))
			copy(msgBytes, scanner.Bytes())
			msgChan <- msgBytes
			if err := scanner.Err(); err != nil {
				log.Error().Err(err).Msg("scanner error")
			}
		}
	}()
}
