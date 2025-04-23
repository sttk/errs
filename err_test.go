package errs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sttk/errs"
)

type /* error reasons */ (
	InvalidValue struct {
		Name  string
		Value string
	}

	FailToGetValue struct {
		Name string
	}
)

type InvalidValueError struct {
	Name  string
	Value string
}

func (e InvalidValueError) Error() string {
	return "InvalidValue { Name: " + e.Name + ", Value: " + e.Value + " }"
}

///

func TestErr(t *testing.T) {

	t.Run("New", func(t *testing.T) {
		t.Run("reason is a value", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"})

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.InvalidValue { Name: foo, Value: abc }, file = err_test.go, line = 38 }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a pointer", func(t *testing.T) {
			err := errs.New(&InvalidValue{Name: "foo", Value: "abc"})

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.InvalidValue { Name: foo, Value: abc }, file = err_test.go, line = 45 }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is nil", func(t *testing.T) {
			err := errs.New(nil)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = nil, file = err_test.go, line = 52 }")
			assert.Nil(t, err.Cause())
		})

		t.Run("cause is an error", func(t *testing.T) {
			cause := errors.New("def")
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.InvalidValue { Name: foo, Value: abc }, file = err_test.go, line = 60, cause = def }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("cause is a custom error", func(t *testing.T) {
			cause := InvalidValueError{Name: "bar", Value: "def"}
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.InvalidValue { Name: foo, Value: abc }, file = err_test.go, line = 68, cause = InvalidValue { Name: bar, Value: def } }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("cause is also an errs.Err", func(t *testing.T) {
			cause := errs.New(FailToGetValue{Name: "foo"})
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.InvalidValue { Name: foo, Value: abc }, file = err_test.go, line = 76, cause = github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.FailToGetValue { Name: foo }, file = err_test.go, line = 75 } }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("reason is nil but cause is given", func(t *testing.T) {
			cause := errors.New("def")
			err := errs.New(nil, cause)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = nil, file = err_test.go, line = 84, cause = def }")
			assert.Equal(t, err.Cause(), cause)
		})

		t.Run("reason is pointer to nil", func(t *testing.T) {
			var reason error = nil
			err := errs.New(&reason)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = <nil>, file = err_test.go, line = 92 }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a boolean", func(t *testing.T) {
			err := errs.New(true)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = true, file = err_test.go, line = 99 }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a number", func(t *testing.T) {
			err := errs.New(123)

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = 123, file = err_test.go, line = 106 }")
			assert.Nil(t, err.Cause())
		})

		t.Run("reason is a string", func(t *testing.T) {
			err := errs.New("abc")

			assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = abc, file = err_test.go, line = 113 }")
			assert.Nil(t, err.Cause())
		})
	})

	t.Run("Ok", func(t *testing.T) {
		err := errs.Ok()

		assert.Equal(t, err.Error(), "github.com/sttk/errs.Err { reason = nil, file = , line = 0 }")
		assert.Nil(t, err.Cause())
	})

	t.Run("Switch expression for reason", func(t *testing.T) {
		t.Run("reason is a value", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"})

			switch r := err.Reason().(type) {
			case InvalidValue:
				assert.Equal(t, r.Name, "foo")
				assert.Equal(t, r.Value, "abc")
			default:
				assert.Fail(t, err.Error())
			}
		})

		t.Run("reason is a pointer", func(t *testing.T) {
			err := errs.New(&InvalidValue{Name: "foo", Value: "abc"})

			switch r := err.Reason().(type) {
			case *InvalidValue:
				assert.Equal(t, r.Name, "foo")
				assert.Equal(t, r.Value, "abc")
			default:
				assert.Fail(t, err.Error())
			}
		})
	})

	t.Run("IsOk, IsNotOk", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			err := errs.Ok()

			assert.True(t, err.IsOk())
			assert.False(t, err.IsNotOk())
		})

		t.Run("reason is a value", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"})

			assert.False(t, err.IsOk())
			assert.True(t, err.IsNotOk())
		})

		t.Run("reason is a pointer", func(t *testing.T) {
			err := errs.New(&InvalidValue{Name: "foo", Value: "abc"})

			assert.False(t, err.IsOk())
			assert.True(t, err.IsNotOk())
		})
	})

	t.Run("apply errors.Is", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			err := errs.Ok()
			assert.Nil(t, err.Unwrap())

			err0 := errs.Ok()
			err1 := errs.New("def")
			err2 := InvalidValueError{Value: "def"}
			err3 := errs.New(InvalidValue{Name: "foo", Value: "abc"})
			err4 := errs.New(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.True(t, errors.Is(err, err0))
			assert.False(t, errors.Is(err, err1))
			assert.False(t, errors.Is(err, err2))
			assert.False(t, errors.Is(err, err3))
			assert.False(t, errors.Is(err, err4))
		})

		t.Run("reason is a value and with no cause", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"})
			assert.Nil(t, err.Unwrap())

			err0 := errs.Ok()
			err1 := errs.New("def")
			err2 := InvalidValueError{Value: "def"}
			err3 := errs.New(InvalidValue{Name: "foo", Value: "abc"})
			err4 := errs.New(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, err0))
			assert.False(t, errors.Is(err, err1))
			assert.False(t, errors.Is(err, err2))
			assert.False(t, errors.Is(err, err3))
			assert.False(t, errors.Is(err, err4))
		})

		t.Run("reason is a pointer and with no cause", func(t *testing.T) {
			err := errs.New(&InvalidValue{Name: "foo", Value: "abc"})
			assert.Nil(t, err.Unwrap())

			err0 := errs.Ok()
			err1 := errs.New("def")
			err2 := InvalidValueError{Value: "def"}
			err3 := errs.New(InvalidValue{Name: "foo", Value: "abc"})
			err4 := errs.New(&InvalidValue{Name: "foo", Value: "abc"})

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, err0))
			assert.False(t, errors.Is(err, err1))
			assert.False(t, errors.Is(err, err2))
			assert.False(t, errors.Is(err, err3))
			assert.False(t, errors.Is(err, err4))
		})

		t.Run("reason is a value and with cause", func(t *testing.T) {
			cause := errors.New("def")
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			err0 := errs.Ok()
			err1 := errs.New("def")
			err2 := InvalidValueError{Value: "def"}
			err3 := errs.New(InvalidValue{Name: "foo", Value: "abc"})
			err4 := errs.New(&InvalidValue{Name: "foo", Value: "abc"})
			err5 := errs.New(InvalidValue{Name: "foo", Value: "abc"}, err1)
			err6 := errs.New(&InvalidValue{Name: "foo", Value: "abc"}, err1)
			err7 := errs.New(InvalidValue{Name: "foo", Value: "abc"}, cause)
			err8 := errs.New(&InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, err0))
			assert.False(t, errors.Is(err, err1))
			assert.False(t, errors.Is(err, err2))
			assert.False(t, errors.Is(err, err3))
			assert.False(t, errors.Is(err, err4))
			assert.False(t, errors.Is(err, err5))
			assert.False(t, errors.Is(err, err6))
			assert.False(t, errors.Is(err, err7))
			assert.False(t, errors.Is(err, err8))

			assert.False(t, errors.Is(err, err1))
			assert.False(t, errors.Is(err0, err1))
			assert.True(t, errors.Is(err1, err1))
			assert.False(t, errors.Is(err2, err1))
			assert.False(t, errors.Is(err3, err1))
			assert.False(t, errors.Is(err4, err1))
			assert.True(t, errors.Is(err5, err1))
			assert.True(t, errors.Is(err6, err1))
			assert.False(t, errors.Is(err7, err1))
			assert.False(t, errors.Is(err8, err1))

			assert.True(t, errors.Is(err, cause))
			assert.False(t, errors.Is(err0, cause))
			assert.False(t, errors.Is(err1, cause))
			assert.False(t, errors.Is(err2, cause))
			assert.False(t, errors.Is(err3, cause))
			assert.False(t, errors.Is(err4, cause))
			assert.False(t, errors.Is(err5, cause))
			assert.False(t, errors.Is(err6, cause))
			assert.True(t, errors.Is(err7, cause))
			assert.True(t, errors.Is(err8, cause))
		})

		t.Run("reason is a pointer and with cause", func(t *testing.T) {
			cause := errors.New("def")
			err := errs.New(&InvalidValue{Name: "foo", Value: "abc"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			err0 := errs.Ok()
			err1 := errs.New("def")
			err2 := InvalidValueError{Value: "def"}
			err3 := errs.New(InvalidValue{Name: "foo", Value: "abc"})
			err4 := errs.New(&InvalidValue{Name: "foo", Value: "abc"})
			err5 := errs.New(InvalidValue{Name: "foo", Value: "abc"}, err1)
			err6 := errs.New(&InvalidValue{Name: "foo", Value: "abc"}, err1)
			err7 := errs.New(InvalidValue{Name: "foo", Value: "abc"}, cause)
			err8 := errs.New(&InvalidValue{Name: "foo", Value: "abc"}, cause)

			assert.True(t, errors.Is(err, err))
			assert.False(t, errors.Is(err, err0))
			assert.False(t, errors.Is(err, err1))
			assert.False(t, errors.Is(err, err2))
			assert.False(t, errors.Is(err, err3))
			assert.False(t, errors.Is(err, err4))
			assert.False(t, errors.Is(err, err5))
			assert.False(t, errors.Is(err, err6))
			assert.False(t, errors.Is(err, err7))
			assert.False(t, errors.Is(err, err8))

			assert.False(t, errors.Is(err, err1))
			assert.False(t, errors.Is(err0, err1))
			assert.True(t, errors.Is(err1, err1))
			assert.False(t, errors.Is(err2, err1))
			assert.False(t, errors.Is(err3, err1))
			assert.False(t, errors.Is(err4, err1))
			assert.True(t, errors.Is(err5, err1))
			assert.True(t, errors.Is(err6, err1))
			assert.False(t, errors.Is(err7, err1))
			assert.False(t, errors.Is(err8, err1))

			assert.True(t, errors.Is(err, cause))
			assert.False(t, errors.Is(err0, cause))
			assert.False(t, errors.Is(err1, cause))
			assert.False(t, errors.Is(err2, cause))
			assert.False(t, errors.Is(err3, cause))
			assert.False(t, errors.Is(err4, cause))
			assert.False(t, errors.Is(err5, cause))
			assert.False(t, errors.Is(err6, cause))
			assert.True(t, errors.Is(err7, cause))
			assert.True(t, errors.Is(err8, cause))
		})
	})

	t.Run("apply errors.As", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			err := errs.Ok()
			assert.Nil(t, err.Unwrap())

			var err0 errs.Err
			//var err1 error
			var err2 InvalidValueError

			assert.True(t, errors.As(err, &err0))
			assert.Equal(t, err0.Error(), err.Error())

			//assert.False(t, errors.As(err, err1)) // --> compile error
			assert.False(t, errors.As(err, &err2))
		})

		t.Run("reason is a value and with no cause", func(t *testing.T) {
			err := errs.New(FailToGetValue{Name: "foo"})
			assert.Nil(t, err.Unwrap())

			var err0 errs.Err
			// var err1 error
			var err2 InvalidValueError

			assert.True(t, errors.As(err, &err0))
			assert.Equal(t, err0.Error(), err.Error())

			//assert.False(t, errors.As(err, err1)) // --> compile error
			assert.False(t, errors.As(err, &err2))
		})

		t.Run("reason is a pointer and with no cause", func(t *testing.T) {
			err := errs.New(&FailToGetValue{Name: "foo"})
			assert.Nil(t, err.Unwrap())

			var err0 errs.Err
			// var err1 error
			var err2 InvalidValueError

			assert.True(t, errors.As(err, &err0))
			assert.Equal(t, err0.Error(), err.Error())

			//assert.False(t, errors.As(err, err1)) // --> compile error
			assert.False(t, errors.As(err, &err2))
		})

		t.Run("reason is a value and with cause", func(t *testing.T) {
			cause := InvalidValueError{Name: "a", Value: "b"}
			err := errs.New(InvalidValue{Name: "foo", Value: "abc"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			var err0 errs.Err
			// var err1 error
			var err2 InvalidValueError

			assert.True(t, errors.As(err, &err0))
			assert.Equal(t, err0.Error(), err.Error())

			//assert.False(t, errors.As(err, err1)) // --> compile error

			assert.True(t, errors.As(err, &err2))
			assert.Equal(t, err2.Name, cause.Name)
			assert.Equal(t, err2.Value, cause.Value)
		})

		t.Run("reason is a pointer and with cause", func(t *testing.T) {
			cause := InvalidValueError{Name: "a", Value: "b"}
			err := errs.New(&FailToGetValue{Name: "foo"}, cause)
			assert.Equal(t, err.Unwrap(), cause)

			var err0 errs.Err
			// var err1 error
			var err2 InvalidValueError

			assert.True(t, errors.As(err, &err0))
			assert.Equal(t, err0.Error(), err.Error())

			//assert.False(t, errors.As(err, err1)) // --> compile error

			assert.True(t, errors.As(err, &err2))
			assert.Equal(t, err2.Name, cause.Name)
			assert.Equal(t, err2.Value, cause.Value)
		})
	})

	t.Run("IfOkThen", func(t *testing.T) {
		t.Run("ok1 & ok2 -> ok2", func(t *testing.T) {
			err := errs.Ok()

			var done bool
			err2 := err.IfOkThen(func() errs.Err {
				done = true
				return errs.Ok()
			})
			assert.True(t, err2.IsOk())
			assert.True(t, done)
		})

		t.Run("ok & error -> error", func(t *testing.T) {
			err := errs.Ok()

			var done bool
			err2 := err.IfOkThen(func() errs.Err {
				done = true
				return errs.New(InvalidValue{Name: "abc", Value: "def"})
			})
			assert.True(t, err2.IsNotOk())
			assert.True(t, done)

			switch r := err2.Reason().(type) {
			case InvalidValue:
				assert.Equal(t, r.Name, "abc")
				assert.Equal(t, r.Value, "def")
			default:
				assert.Fail(t, err2.Error())
			}
		})

		t.Run("error1 & error2 -> error1", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "abc", Value: "def"})

			var done bool
			err2 := err.IfOkThen(func() errs.Err {
				done = true
				return errs.New(FailToGetValue{Name: "abc"})
			})

			assert.True(t, err2.IsNotOk())
			assert.False(t, done)

			switch r := err2.Reason().(type) {
			case InvalidValue:
				assert.Equal(t, r.Name, "abc")
				assert.Equal(t, r.Value, "def")
			default:
				assert.Fail(t, err2.Error())
			}
		})

		t.Run("error1 & ok -> error1", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "abc", Value: "def"})

			var done bool
			err2 := err.IfOkThen(func() errs.Err {
				done = true
				return errs.Ok()
			})

			assert.True(t, err2.IsNotOk())
			assert.False(t, done)

			switch r := err2.Reason().(type) {
			case InvalidValue:
				assert.Equal(t, r.Name, "abc")
				assert.Equal(t, r.Value, "def")
			default:
				assert.Fail(t, err2.Error())
			}
		})
	})

	t.Run("Print", func(t *testing.T) {
		t.Run("%v", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "abc", Value: "def"})
			assert.Equal(t, fmt.Sprintf("%v", err), `github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.InvalidValue { Name: abc, Value: def }, file = err_test.go, line = 493 }`)
		})

		t.Run("%w", func(t *testing.T) {
			err := errs.New(InvalidValue{Name: "abc", Value: "def"})
			assert.Equal(t, fmt.Errorf("%w", err).Error(), `github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.InvalidValue { Name: abc, Value: def }, file = err_test.go, line = 498 }`)
		})
	})
}
