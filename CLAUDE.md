# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

`ibapi` (module `github.com/scmhub/ibapi`) is an unofficial Go port of Interactive Brokers' TWS API. It mirrors the official Python/C++ `tws-api` client: it opens a TCP socket to TWS/IB Gateway, encodes outgoing requests, decodes incoming messages, and dispatches them to user-supplied callbacks. This is a library (no `main` package at the root); `examples/` contains runnable demonstrations.

## Common commands

```bash
go build ./...          # compile everything
go vet ./...            # static analysis
gofmt -l .              # check formatting
golangci-lint run ./... # linter (config: .golangci.yml)
go test ./... -run TestName   # run a single test by name
```

There is no Makefile in this repo — `go build`/`go vet`/`gofmt`/`golangci-lint run` are the checks to run before considering a change done. `.github/workflows/ci.yml` runs these same checks (plus the self-contained subset of tests that don't need a live TWS/Gateway, see below) on every push/PR to `main`.

### Tests require a live TWS/IB Gateway connection

Most tests in `client_test.go` are **integration tests**, not unit tests: they call `setupIBClient` which dials a real TWS/IB Gateway instance (default `localhost:7497`, overridable via `IB_HOST`/`IB_PORT`/`IB_ACCOUNT` env vars) and exercise live request/response round trips. Running `go test ./...` without a running TWS/Gateway will fail/hang on those tests. `connection_race_test.go` and `condition_proto_test.go` are self-contained (they spin up a dummy TCP server or test pure encoding logic) and can run without TWS.

## Architecture

### Request/response flow

1. **`EClient`** (`client.go`, ~6.4k lines) is the main entry point users instantiate via `NewEClient(wrapper EWrapper)`. It exposes one `Req*`/`Cancel*`/`Place*` method per IB API call (e.g. `ReqCurrentTime`, `ReqMktData`, `PlaceOrder`). Each method:
   - checks `c.useProtoBuf(msgID)` to decide between the legacy field-delimited wire format and protobuf, dispatching to a `xxxProtoBuf` sibling method when applicable,
   - validates connection state and `c.serverVersion` against `MIN_SERVER_VER_*` constants (`server_versions.go`), erroring via `c.wrapper.Error(...)` if the server is too old,
   - builds the message with `NewMsgEncoder` (in `client.go`) and pushes the encoded bytes onto `c.reqChan`, which a writer goroutine drains to the socket.
2. **`Connection`** (`connection.go`) wraps the raw `net.TCPConn` with lock-free atomic stats and automatic reconnect-with-backoff on write/read failure.
3. **`EReader`** (`reader.go`) runs two goroutines wired by a buffered channel: a scanner goroutine (`bufio.Scanner` with the custom `scanFields` split function from `utils.go`, which frames messages by a 4-byte big-endian length prefix) and a single decoder goroutine that must stay single-threaded (see the `// single worker and no go here!!` comment) because `EDecoder` is not safe for concurrent use.
4. **`EDecoder`** (`decoder.go`, ~4.4k lines) — `parseAndProcessMsg` reads the message ID, decides protobuf vs. legacy framing per message (same `PROTOBUF_MSG_ID` offset trick as the encoder), and switches on the IN-message ID (`message.go`) to a `processXxxMsg`/`processXxxMsgProtoBuf` handler in `decoder.go`/`decoder_utils.go`/`order_decoder.go`, which ultimately calls the matching method on the user's `EWrapper`.
5. **`EWrapper`** (`wrapper.go`) is the callback interface users implement to receive data (ticks, order status, account data, errors, etc.). `Wrapper` is a no-op default implementation (each method just logs) meant to be embedded so callers only override the callbacks they care about (see `examples/blocking/blocking.go`).

### Message IDs and versioning

- `message.go` defines the `IN` (incoming) and (further down) `OUT` (outgoing) message ID constants used to route decoding/encoding — this is the map to consult when adding support for a new message type.
- `server_versions.go` defines `MIN_SERVER_VER_*` constants gating features by TWS API server version; `EClient` checks these before sending newer request fields/messages.
- Both `EClient` (encode side) and `EDecoder` (decode side) support two wire formats side by side: the original delimited-field format and newer protobuf messages (`protobuf/*.pb.go`, generated from `proto/*.proto`). New code touching request/response encoding should preserve both paths (`useProtoBuf` branch + legacy branch) rather than assuming one.

### Supporting data types

- `contract.go`, `order.go`, `order_state.go`, `execution.go`, `bar.go`, `decimal.go` (built on `github.com/robaho/fixed`), and various `*_condition.go` files define the domain model shared between encoder and decoder.
- `*_samples.go` and `*_samples_proto.go` files (`contract_samples.go`, `order_samples.go`, `scanner_subscription_samples.go`, `fa_allocation_sample.go`, etc.) are example/helper constructors, not part of the wire protocol — useful references when building a `Contract`/`Order` correctly.
- `errors.go` defines the sentinel `CodeMsg` error values (e.g. `NOT_CONNECTED`, `FAIL_CREATE_SOCK`) reported through `wrapper.Error(...)`, mirroring the official API's error codes.
- `logger.go` wraps `zerolog`; use `ibapi.Logger()`, `ibapi.SetLogLevel()`, `ibapi.SetConsoleWriter()` rather than introducing a separate logging mechanism.

### Regenerating protobuf code

`proto/*.proto` are the IB protobuf schema sources; `protobuf/*.pb.go` are their generated Go bindings (imported as `github.com/scmhub/ibapi/protobuf`). Treat `protobuf/*.pb.go` as generated output — regenerate from `proto/` rather than hand-editing when the schema changes.
