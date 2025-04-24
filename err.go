// Copyright (C) 2025 Takayuki Sato. All Rights Reserved.
// This program is free software under MIT License.
// See the file LICENSE in this distribution for more details.

// Package errs is for error handling in Go programs, providing an Err struct which represents
// an error with a reason.
//
// The type of this reason is any, but typically an struct is used.
// The name of this struct indicates the reason for the error, and its fields store contextual
// information about the situation in which the error occurred.
// Since the path of an struct, including its package, is unique within a program, the struct
// representing the reason is useful for identifying the specific error, locating where it
// occurred, or generating appropriate error messages, etc.
//
// Optionally, by using the build tag: github.sttk.errs.notify and registering error handlers
// in advance, it is possible to receive notifications either synchronously or asynchronously
// at the time the error struct is created.
//
// # Install
//
// To use this package in your code, the following import declaration is necessary.
//
//	import "github.com/sttk/errs"
//
// # Usage
//
// # Err instantiation and identification of a reason
//
// The Err struct can be instantiated with New function.
// Then, the reason can be identified with Reason method and a type switch statement,
//
//	type /* error reasons */ (
//	    IllegalState struct { State string }
//	    // ...
//	}
//
//	err := errs.New(IllegalState{State: "bad state"})
//
//	switch r : err.Reason().(type) {
//	case nil:
//	    ....
//	case IllegalState:
//	    fmt.Printf("state = %s\n", r.State)
//	default:
//	    ...
//	}
//
// # Notification of Err instantiations
//
// This package optionally provides a feature to notify pre-registered error handlers when an Err
// is instantiated.
// Multiple error handlers can be registered, and you can choose to receive notifications either
// synchronously or asynchronously.
// To register error handlers that receive notifications synchronously, use the AddSyncErrHandler
// function.
// For asynchronous notifications, use the AddAsyncErrHandler function.
//
// Error notifications will not occur until the FixErrHandlers function is called.
// This function locks the current set of error handlers, preventing further additions and enabling
// notification processing.
//
//	errs.AddAsyncErrHandler(func(err errs.Err, tm time.Time) {
//	    fmt.Printf("%s (%s:%d) %v\n",
//	        tm.Format("2006-01-02T15:04:05Z"),
//	        err.File(), err.Line(), err)
//	});
//
//	errs.AddSyncErrHandler(func(err errs.Err, tm time.Time) {
//	    // ...
//	});
//
//	errs.FixErrHandlers()
//
// NOTE: To use this feature, it is necessary to specify the following build tag to go build
// command:
//
//	go build -tags github.sttk.errs.notify
package errs

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
)

// Err is the struct that represents an error with a reason.
//
// This struct encapsulates the reason for the error, which can be any data type.
// Typically, the reason is a struct, which makes it easy to uniquely identify the error kind and
// location in the source code.
// In addition, since a struct can store additional information as their fields, it is possible to
// provide more detailed information about the error.
//
// The reason for the error can be distinguished with a type switch statement, and type casting,
// so it is easy to handle the error in a type-safe manner.
//
// This struct also contains an optional cause error, which is the error caused the current error.
// This is useful for chaining errors.
//
// This struct is implements the Error method, so it can be used as an error object in Go programs.
// And since this struct implements the Unwrap method, it can be used as a wrapper error object in
// Go programs.
type Err struct {
	reason any
	cause  error
	file   string
	line   int
}

// Ok returns an instance of Err with no reason, indicating no error.
// This is the default "no error" state.
func Ok() Err {
	return Err{}
}

// New creates a new Err instance with the provided reason.
// Optionally, a cause can also be supplied, which represents a lower-level error.
func New(reason any, cause ...error) Err {
	var e Err
	e.reason = reason

	if len(cause) > 0 {
		e.cause = cause[0]
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		e.file = filepath.Base(file)
		e.line = line
	}

	notifyErr(e)

	return e
}

// Reason returns the reason for the error, which can be any type.
// This helps in analyzing why the error occurred.
func (e Err) Reason() any {
	return e.reason
}

// File returns the base name of the file where the error occurred.
func (e Err) File() string {
	return e.file
}

// Line returns the line number in the file where the error occurred.
func (e Err) Line() int {
	return e.line
}

// Error returns a string representation of the Err instance.
// It formats the error, including the package path, reason, and cause.
func (e Err) Error() string {
	var buf bytes.Buffer

	t := reflect.TypeOf(e)
	s := t.PkgPath()
	if len(s) > 0 {
		buf.WriteString(s)
		buf.WriteByte('.')
	}
	buf.WriteString(t.Name())

	buf.WriteString(" { reason = ")

	if e.reason == nil {
		buf.WriteString("nil")
	} else {
		v := reflect.ValueOf(e.reason)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			if v.CanInterface() {
				buf.WriteString(fmt.Sprintf("%v", v.Interface()))
			}
		} else {
			t := v.Type()

			s := t.PkgPath()
			if len(s) > 0 {
				buf.WriteString(s)
				buf.WriteByte('.')
			}
			buf.WriteString(t.Name())

			n := v.NumField()

			if n > 0 {
				buf.WriteString(" { ")

				for i := 0; i < n; i++ {
					if i > 0 {
						buf.WriteString(", ")
					}

					k := t.Field(i).Name

					f := v.Field(i)
					if f.CanInterface() { // false, if the field is not public
						buf.WriteString(k)
						buf.WriteString(": ")
						buf.WriteString(fmt.Sprintf("%v", f.Interface()))
					}
				}

				buf.WriteString(" }")
			}
		}
	}

	buf.WriteString(", file = ")
	buf.WriteString(e.file)
	buf.WriteString(", line = ")
	buf.WriteString(strconv.Itoa(e.line))

	if e.cause != nil {
		buf.WriteString(", cause = ")
		buf.WriteString(e.cause.Error())
	}

	buf.WriteString(" }")

	return buf.String()
}

// Unwrap returns the underlying cause of the error, allowing it to be chained.
// This helps in accessing the root cause when errors are wrapped.
func (e Err) Unwrap() error {
	return e.cause
}

// Cause returns the cause of the error.
// This is similar to Unwrap but provides a direct access method.
func (e Err) Cause() error {
	return e.cause
}

// IsOk returns true if the Err instance has no reason, indicating no error.
// This is used to check if the operation was successful.
func (e Err) IsOk() bool {
	return (e.reason == nil)
}

// IsNotOk returns true if the Err instance has a reason, indicating an error occurred.
// This is the inverse of IsOk, used to determine if an error is present.
func (e Err) IsNotOk() bool {
	return (e.reason != nil)
}

// IfOkThen executes the provided function if no error is present (IsOk).
// This is useful for chaining operations that only proceed if no error has occurred.
func (e Err) IfOkThen(fn func() Err) Err {
	if e.IsOk() {
		return fn()
	}
	return e
}
