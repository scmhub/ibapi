package ibapi

import (
	"bufio"
	"context"
	"slices"
	"sync"
)

// EReader starts the scan and decode goroutines
func EReader(ctx context.Context, cancel context.CancelFunc, scanner *bufio.Scanner, decoder *EDecoder, wg *sync.WaitGroup) {

	msgChan := make(chan []byte, 300)

	// Decode
	wg.Add(1)
	go func() {
		log.Debug().Msg("decoder started")
		defer log.Debug().Msg("decoder ended")
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgChan:
				if !ok {
					return
				}
				decoder.parseAndProcessMsg(msg) // single worker and no go here!!
			}
		}
	}()

	// Scan
	wg.Add(1)
	go func() {
		log.Debug().Msg("scanner started")
		defer log.Debug().Msg("scanner ended")
		defer wg.Done()
		defer close(msgChan) // close the channel so decoder exits

		// for scanner.Scan() {
		// 	msgBytes := make([]byte, len(scanner.Bytes()))
		// 	copy(msgBytes, scanner.Bytes())
		// 	msgChan <- msgBytes
		// 	if err := scanner.Err(); err != nil {
		// 		log.Error().Err(err).Msg("scanner error")
		// 		break
		// 	}
		// }

		for {
			select {
			case <-ctx.Done():
				// shutdown in flight
				return
			default:
				// block here until there's a token or an error/EOF
				if !scanner.Scan() {
					// only take action if we weren't already cancelled
					if ctx.Err() == nil {
						if err := scanner.Err(); err != nil {
							log.Error().Err(err).Msg("scanner error, triggering shutdown")
						} else {
							log.Debug().Msg("scanner reached EOF, triggering shutdown")
						}
						cancel()
					}
					return
				}

				// successful scan → queue for decode
				msgBytes := slices.Clone(scanner.Bytes())
				msgChan <- msgBytes
			}
		}
	}()
}
