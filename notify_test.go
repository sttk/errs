package errs

import (
	"container/list"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ClearErrHandlers() {
	syncErrHandlers = nil
	asyncErrHandlers = nil
	isErrHandlersFixed = false
}

func TestAddErrSyncHandler(t *testing.T) {
	const fn_sig string = "func(errs.Err, time.Time)"

	t.Run("add zero handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		assert.Empty(t, syncErrHandlers)
		assert.Empty(t, syncErrHandlers)
	})

	t.Run("add one handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddSyncErrHandler(func(e Err, tm time.Time) {})

		assert.Len(t, syncErrHandlers, 1)
		assert.Equal(t, reflect.TypeOf(syncErrHandlers[0]).String(), fn_sig)

		assert.Empty(t, asyncErrHandlers)
	})

	t.Run("add two handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddSyncErrHandler(func(e Err, tm time.Time) {})
		AddSyncErrHandler(func(e Err, tm time.Time) {})

		assert.Len(t, syncErrHandlers, 2)
		assert.Equal(t, reflect.TypeOf(syncErrHandlers[0]).String(), fn_sig)
		assert.Equal(t, reflect.TypeOf(syncErrHandlers[1]).String(), fn_sig)

		assert.Empty(t, asyncErrHandlers)
	})
}

func TestAddErrAsyncHandler(t *testing.T) {
	const fn_sig string = "func(errs.Err, time.Time)"

	t.Run("add zero handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		assert.Empty(t, asyncErrHandlers)
	})

	t.Run("add one handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddAsyncErrHandler(func(e Err, tm time.Time) {})

		assert.Empty(t, syncErrHandlers)

		assert.Len(t, asyncErrHandlers, 1)
		assert.Equal(t, reflect.TypeOf(asyncErrHandlers[0]).String(), fn_sig)
	})

	t.Run("add two handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddAsyncErrHandler(func(e Err, tm time.Time) {})
		AddAsyncErrHandler(func(e Err, tm time.Time) {})

		assert.Empty(t, syncErrHandlers)

		assert.Len(t, asyncErrHandlers, 2)
		assert.Equal(t, reflect.TypeOf(asyncErrHandlers[0]).String(), fn_sig)
		assert.Equal(t, reflect.TypeOf(asyncErrHandlers[1]).String(), fn_sig)
	})
}

func TestFixErrHandlers(t *testing.T) {
	t.Run("cannot add any more handlers after fixed", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		AddSyncErrHandler(func(e Err, tm time.Time) {})
		AddAsyncErrHandler(func(e Err, tm time.Time) {})

		assert.Len(t, syncErrHandlers, 1)
		assert.Len(t, asyncErrHandlers, 1)

		assert.False(t, isErrHandlersFixed)

		FixErrHandlers()

		assert.True(t, isErrHandlersFixed)

		AddSyncErrHandler(func(e Err, tm time.Time) {})
		AddAsyncErrHandler(func(e Err, tm time.Time) {})

		assert.Len(t, syncErrHandlers, 1)
		assert.Len(t, asyncErrHandlers, 1)
	})
}

func TestNotifyErr(t *testing.T) {
	t.Run("when there is no error handler", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		type FailToDoSomething struct{}

		assert.False(t, isErrHandlersFixed)
		New(FailToDoSomething{})

		FixErrHandlers()
		assert.True(t, isErrHandlersFixed)
		New(FailToDoSomething{})
	})

	t.Run("notify when there are error handlers", func(t *testing.T) {
		ClearErrHandlers()
		defer ClearErrHandlers()

		syncLogs := list.New()
		asyncLogs := list.New()

		type FailToDoSomething struct{}

		AddSyncErrHandler(func(e Err, tm time.Time) {
			syncLogs.PushBack(fmt.Sprintf("%s-1:%s", e.Error(), tm.String()))
		})
		AddSyncErrHandler(func(e Err, tm time.Time) {
			syncLogs.PushBack(fmt.Sprintf("%s-2:%s", e.Error(), tm.String()))
		})
		AddAsyncErrHandler(func(e Err, tm time.Time) {
			time.Sleep(100 * time.Millisecond)
			asyncLogs.PushBack(fmt.Sprintf("%s-3:%s", e.Error(), tm.String()))
		})
		AddAsyncErrHandler(func(e Err, tm time.Time) {
			time.Sleep(10 * time.Millisecond)
			asyncLogs.PushBack(fmt.Sprintf("%s-4:%s", e.Error(), tm.String()))
		})

		assert.False(t, isErrHandlersFixed)

		New(FailToDoSomething{})

		assert.Equal(t, syncLogs.Len(), 0)
		assert.Equal(t, asyncLogs.Len(), 0)

		FixErrHandlers()

		assert.True(t, isErrHandlersFixed)

		New(FailToDoSomething{})

		assert.Equal(t, syncLogs.Len(), 2)
		log := syncLogs.Front()
		assert.Contains(t, log.Value, "github.com/sttk/errs.Err { reason = github.com/sttk/errs.FailToDoSomething, file = notify_test.go, line = 169 }-1:")
		log = log.Next()
		assert.Contains(t, log.Value, "github.com/sttk/errs.Err { reason = github.com/sttk/errs.FailToDoSomething, file = notify_test.go, line = 169 }-2:")
		log = log.Next()
		assert.Nil(t, log)

		time.Sleep(500 * time.Millisecond)

		assert.Equal(t, asyncLogs.Len(), 2)
		log = asyncLogs.Front()
		assert.Contains(t, log.Value, "github.com/sttk/errs.Err { reason = github.com/sttk/errs.FailToDoSomething, file = notify_test.go, line = 169 }-4:")
		log = log.Next()
		assert.Contains(t, log.Value, "github.com/sttk/errs.Err { reason = github.com/sttk/errs.FailToDoSomething, file = notify_test.go, line = 169 }-3:")
		log = log.Next()
		assert.Nil(t, log)
	})
}
