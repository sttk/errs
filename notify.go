//go:build github.sttk.errs.notify

package errs

import (
	"time"
)

type errHandlerListItem struct {
	handler func(Err, time.Time)
	next    *errHandlerListItem
}

type errHandlerList struct {
	head *errHandlerListItem
	last *errHandlerListItem
}

var (
	syncErrHandlers    = errHandlerList{nil, nil}
	asyncErrHandlers   = errHandlerList{nil, nil}
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

	last := syncErrHandlers.last
	syncErrHandlers.last = &errHandlerListItem{handler, nil}

	if last != nil {
		last.next = syncErrHandlers.last
	}

	if syncErrHandlers.head == nil {
		syncErrHandlers.head = syncErrHandlers.last
	}
}

// AddAsyncErrHandler adds a new asynchronous error handler to the global handler list.
// It will not add the handler if the handlers have been fixed using FixErrHandlers.
//
// NOTE: This function is enabled via the build tag: github.sttk.errs.notify
func AddAsyncErrHandler(handler func(Err, time.Time)) {
	if isErrHandlersFixed {
		return
	}

	last := asyncErrHandlers.last
	asyncErrHandlers.last = &errHandlerListItem{handler, nil}

	if last != nil {
		last.next = asyncErrHandlers.last
	}

	if asyncErrHandlers.head == nil {
		asyncErrHandlers.head = asyncErrHandlers.last
	}
}

// FixErrHandlers prevents further modification of the error handler lists.
// Before this is called, no Err is notified to the handlers.
// After this is called, no new handlers can be added, and Err(s) is notified to the
// handlers.
//
// NOTE: This function is enabled via the build tag: github.sttk.errs.notify
func FixErrHandlers() {
	isErrHandlersFixed = true
}

func notifyErr(e Err) {
	if !isErrHandlersFixed {
		return
	}

	if syncErrHandlers.head == nil && asyncErrHandlers.head == nil {
		return
	}

	tm := time.Now().UTC()

	for item := syncErrHandlers.head; item != nil; item = item.next {
		item.handler(e, tm)
	}

	for item := asyncErrHandlers.head; item != nil; item = item.next {
		go item.handler(e, tm)
	}
}
