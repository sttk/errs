//go:build github.sttk.errs.notify

package errs

import (
	"time"
)

var (
	syncErrHandlers    []func(Err, time.Time)
	asyncErrHandlers   []func(Err, time.Time)
	isErrHandlersFixed = false
)

// AddSyncErrHandler adds a new synchronous error handler to the global handler list.
// It will not add the handler if the handlers have been fixed using FixErrHandlers.
//
// NOTE: This function is enabled via the build tag: github.sttk.errs.notify
func AddSyncErrHandler(handler func(Err, time.Time)) {
	if isErrHandlersFixed {
		return
	}
	syncErrHandlers = append(syncErrHandlers, handler)
}

// AddAsyncErrHandler adds a new asynchronous error handler to the global handler list.
// It will not add the handler if the handlers have been fixed using FixErrHandlers.
//
// NOTE: This function is enabled via the build tag: github.sttk.errs.notify
func AddAsyncErrHandler(handler func(Err, time.Time)) {
	if isErrHandlersFixed {
		return
	}
	asyncErrHandlers = append(asyncErrHandlers, handler)
}

// FixErrHandlers prevents further modification of the error handler lists.
// Before this is called, no Err is notified to the handlers.
// After this is called, no new handlers can be added, and Err(s) is notified to the
// handlers.
//
// NOTE: This function is enabled via the build tag: github.sttk.errs.notify
func FixErrHandlers() {
	if isErrHandlersFixed {
		return
	}
	isErrHandlersFixed = true
	syncErrHandlers = clip(syncErrHandlers)
	asyncErrHandlers = clip(asyncErrHandlers)
}

func clip(s []func(Err, time.Time)) []func(Err, time.Time) {
	return s[:len(s):len(s)]
}

func notifyErr(e Err) {
	if !isErrHandlersFixed {
		return
	}

	if len(syncErrHandlers) == 0 && len(asyncErrHandlers) == 0 {
		return
	}

	tm := time.Now().UTC()

	for _, handler := range syncErrHandlers {
		handler(e, tm)
	}

	for _, handler := range asyncErrHandlers {
		go handler(e, tm)
	}
}
