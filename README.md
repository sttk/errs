# [errs][repo-url] [![Release][release-img]][release-url] [![Go Reference][pkg-dev-img]][pkg-dev-url] [![CI Status][ci-img]][ci-url] [![MIT License][mit-img]][mit-url]

A library for handling errors with reasons for Go.

## Overview

`errs` is an error handling library for Go designed to focus on the "Reason" behind an error.

### Expressing "Why It Failed" via the Type System

Rather than treating errors as simple message strings or type-erased objects, it embraces a design that expresses "why it failed" through types, allowing for safe and clear propagation and determination.

For error reasons, you can use anything from lightweight types like `string` to type-safe definitions using `struct`, all handled flexibly with the same API.
By using an `struct` in particular, you can not only express failure factors within the type system but also hold contextual information in its fields, propagating the situation and relevant data at the time of the error as-is.
Furthermore, since reasons can be determined in a type-safe manner using type-switch statement, you can avoid fragile error handling that relies on string comparisons.

### Decentralized Error Definition and Traceability

`errs` encourages defining error reasons close to where they occur.
This eliminates the need to share a massive, monolithic error type across the entire application, enabling a highly maintenable design while keeping dependencies between modules clean.
Type information is utilized to identify the reason, and the type identifiers required for this determination are resolved statically at compile time.
This provides type-safe error handling with minimal runtime overhead.

The core `Err` type of the library implements `error`, allowing it to integrate naturally with standard Go error handling.
It can also retain lower-layer errors as causes, enabling you to manage the "Reason" of the upper layer and the "Cause" of the lower layer separately.
Additionally, it automatically records the file name and line number when an error is generated, making log output and failure analysis effortless.

### Powerful Error-Instantiation Notification & Monitoring Ecosystem

Furthermore, `errs` features a mechanism to notify error generation events.
By compiling with the build tag: `github.sttk.errs.notify`, an automatic notification can be sent to registered handlers the exact moment an `Err` is created.
It supports synchronous handlers and asynchronous handlers, and it accommodates registration within functions.
This makes it easy to implement logging, monitoring, metrics collection, and integration with telemetry systems.

While standard `fmt.Errorf` and traditional wrapper libraries focus primarily on annotating and propagating errors, `errs` emphasizes explicitly defining the reason for failure through types and reliably observing the exact moment it occurs.
This library is ideal for scenarios where you want to tightly manage the semantics of errors within your application while seamlessly integrating with production monitoring and operational infrastructure.

## Install

To get the latest version of this package, run the following command:

```bash
go get github.com/sttk/errs
```

To get a specific version of this package, run the following command:

```bash
go get github.com/sttk/errs@v0.1.1
```

## Usage

### Locally Defined Reasons and Instantiate an Err with Them

An `Err` struct can be instantiated with any arbitrary error reason.
Typically, a struct defined to indicate the cause or context of the error is used as the reason.
This reason does not need to belong to a single, centrally managed something; rather, it is preferable to define it close to where the error using it as a reason actually occurs.

```go
import "github.com/sttk/errs"

type /* error reasons */ (
  IllegalState struct { state string }
)

err := errs.New(IllegalState { state: "bad state" })
```

An Err can also be instantiated with the underlying cause error along with the reason.

```go
import (
  "fmt"
  "github.com/sttk/errs"
)

cause := fmt.Errorf("causal error")

err := errs.New(IllegalState { state: "bad state" }, cause)
```

### Type-Safe Reason Identification

By using the type-switch statement, you can extract the error reason as the specified type.

```go
type /* error reasons */ (
  FailToDoSomething struct {}

  FailToDoWithParams struct {
    Param1 string,
    Param2 int,
  }
)

switch reason := err.Reason().(type) {
case FailToDoSomething:
  fmt.Println("FailToDoSomething")
case FailToDoWithParams:
  fmt.Printf("FailToDoWithParam: Param1 = %s, Params = %d\n", reason.Param1, reason.Param2)
default:
  fmt.Println("Unknown reason")
}
```

### Error Handler Registration

This library optionally provides a feature to notify pre-registered error handlers when an `Err` is instantiated.
Multiple error handlers can be registered, and you can choose to receive notifications either synchronously or asynchronously.

To register handlers, you can use the following functions:

* `errs.AddSyncErrHandler`: For synchronous handlers.
* `errs.AddAsyncErrHandler`: For asynchronous handlers.

Error notifications will not occur until the `FixErrHandlers` function is called.
This function locks the current set of error handlers, preventing further additions and enabling notification processing.
```go
errs.AddSyncErrHandler(func(e errs.Err, tm time.Time) {
    fmt.Println("SYNC:", tm, e)
})

errs.AddAsyncErrHandler(func(e errs.Err, tm time.Time) {
    logToRemoteServer(e, tm)
})

// Fix the handlers to start receiving notifications.
errs.FixErrHandlers()
```

## Supporting Go versions

This framework supports Go 1.18 or later.

### Actual test results for each Go version:

```sh
% go-fav -tags=github.sttk.errs.notify 1.26.2 1.25.9 \
         -ldflags="-linkmode=external" 1.24.13 1.23.12 1.22.12 1.21.13 1.20.14 1.19.13 1.18.10
go version go1.26.2 darwin/amd64
ok  	github.com/sttk/errs	1.230s	coverage: 98.8% of statements

go version go1.25.9 darwin/amd64
ok  	github.com/sttk/errs	1.191s	coverage: 98.8% of statements

go version go1.24.13 darwin/amd64
ok  	github.com/sttk/errs	1.180s	coverage: 98.8% of statements

go version go1.23.12 darwin/amd64
ok  	github.com/sttk/errs	1.173s	coverage: 98.8% of statements

go version go1.22.12 darwin/amd64
ok  	github.com/sttk/errs	1.168s	coverage: 98.8% of statements

go version go1.21.13 darwin/amd64
ok  	github.com/sttk/errs	1.208s	coverage: 98.8% of statements

go version go1.20.14 darwin/amd64
ok  	github.com/sttk/errs	1.174s	coverage: 98.8% of statements

go version go1.19.13 darwin/amd64
ok  	github.com/sttk/errs	1.170s	coverage: 98.8% of statements

go version go1.18.10 darwin/amd64
ok  	github.com/sttk/errs	1.159s	coverage: 98.8% of statements
```

## License

Copyright (C) 2025-2026 Takayuki Sato

This program is free software under MIT License.<br>
See the file LICENSE in this distribution for more details.


[repo-url]: https://github.com/sttk/errs
[release-img]: https://img.shields.io/badge/release-0.1.1-0f9999.svg
[release-url]: https://github.com/sttk/errs/releases
[pkg-dev-img]: https://pkg.go.dev/badge/github.com/sttk/errs.svg
[pkg-dev-url]: https://pkg.go.dev/github.com/sttk/errs
[ci-img]: https://github.com/sttk/errs/actions/workflows/go.yml/badge.svg
[ci-url]: https://github.com/sttk/errs/actions?query=branch%3Amain
[mit-img]: https://img.shields.io/badge/license-MIT-green.svg
[mit-url]: https://opensource.org/licenses/MIT
