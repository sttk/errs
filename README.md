# [errs][repo-url] [![Go Reference][pkg-dev-img]][pkg-dev-url] [![CI Status][ci-img]][ci-url] [![MIT License][mit-img]][mit-url]

`errs` is a package for handling errors with reasons for Golang programs.

This package provides the structure type `Err` which takes a parameter of any type as a reason for an error.
This parameter is typically a structure type, and its name represents the reason, and its fields represent the situation where the error occurred.

## Features

- **Structured error representation** using any Go type as the error reason.
- **Type-safe error handling** through type switches or type assertions.
- **Optional error cause** to support error chaining.
- **Captured file and line number** where the error was created.
- **Optional error notification system** enabled via build tag: `github.sttk.errs.notify`.

## Installation

```sh
go get github.com/sttk/errs
```

## Usage

### Creates `Err`(s)

First, imports `errs` package as follows:

```go
import "github.com/sttk/errs"
```

Next, defines structure types which represent reasons of errors.

```go
type /* error reasons */ (
  FailToDoSomething struct {}

  FailToDoWithParams struct {
    Param1 string,
    Param2 int,
  }
)
```

Then, creates `Err`(s).

```go
func f0() errs.Err {
    ...
    return errs.Ok()
}
```

```go
func f1() errs.Err {
    ...
    return errs.New(FailToDoSomething{})
}

func f2() errs.Err {
    ...
    return errs.New(FailToDoWithParams{Param1: "abc", Param2: 123})
}
```

It is enabled to use an `Err` as an `error`.

```go
func f3() error {
    ...
    return errs.New(FailToDoSomething{})
}
```

It is enabled to take a cause error.

```go
var cause = errors.New("I/O timeout")

func f4() errs.Err {
    ...
    return errs.New(FailToDoWithParams{Param1: "abc", Param2: 123}, cause)
}
```

### Operating a `Err`

```go
err := f4()

fmt.Println(err.Reason())  // => path/to/pkg.FailToDoWithParams { Param1: abc, Param2: 123 }
fmt.Println(err.File())    // e.g. source_file.go
fmt.Println(err.Line())    // e.g. 123
fmt.Println(err.Cause())   // => I/O timeout
fmt.Println(err.Error())   // => github.com/sttk/errs.Err { reason = path/to/pkg.FailToDoWithParams { Param1: abc, Param2: 123 }, file = source_file.go, line = 123, cause = I/O timeout }

fmt.Println(err.IsOk())    // => true
fmt.Println(err.IsNotOk()) // => false

fmt.Println(err.Unwrap())  // => I/O timeout
fmt.Println(errors.Is(err, cause)) // => true
```
```go
err := f0().IfOkThen(func() errs.Err {
    // This function is executed.
})
```
```go
err := f1().IsOkThen(func() errs.Err {
    // This function is not executed.
})
```

### Type-safe error handling

A reason of an `Err` can be identified with a type-switch statement.

```go
switch reason := err.Reason().(type) {
case FailToDoSomething:
  fmt.Println("FailToDoSomething")
case FailToDoWithParams:
  fmt.Printf("FailToDoWithParam: Param1 = %s, Param = %d\n", reason.Param1, reason.Param2)
default:
  fmt.Println("Unknown reason")
}
```

### Error notification (Optional)

> To use the error notification feature, build with the tag: `github.sttk.errs.notify`.

Adds synchronous/asynchronous handlers.

```go
errs.AddSyncErrHandler(func(e errs.Err, tm time.Time) {
    fmt.Println("SYNC:", tm, e)
})

errs.AddAsyncErrHandler(func(e errs.Err, tm time.Time) {
    logToRemoteServer(e, tm)
})
```

Prevents the addition of extra error handlers and enables error notifications.

```go
errs.FixErrHandlers()
```

After this point, each time an `Err` is created, it will be notified to the registered error handlers.

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
[pkg-dev-img]: https://pkg.go.dev/badge/github.com/sttk/errs.svg
[pkg-dev-url]: https://pkg.go.dev/github.com/sttk/errs
[ci-img]: https://github.com/sttk/errs/actions/workflows/go.yml/badge.svg
[ci-url]: https://github.com/sttk/errs/actions?query=branch%3Amain
[mit-img]: https://img.shields.io/badge/license-MIT-green.svg
[mit-url]: https://opensource.org/licenses/MIT
