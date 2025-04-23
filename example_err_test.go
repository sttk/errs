package errs_test

import (
	"errors"
	"fmt"

	"github.com/sttk/errs"
)

func ExampleNew() {
	type /* error reasons */ (
		FailToDoSomething  struct{}
		FailToDoWithParams struct {
			Param1 string
			Param2 int
		}
	)

	// (1) Creates an Err with no parameter.
	err := errs.New(FailToDoSomething{})
	fmt.Printf("(1) %v\n", err)

	// (2) Creates an Err with parameters.
	err = errs.New(FailToDoWithParams{
		Param1: "ABC",
		Param2: 123,
	})
	fmt.Printf("(2) %v\n", err)

	cause := errors.New("Causal error")

	// (3) Creates an Err with a causal error.
	err = errs.New(FailToDoSomething{}, cause)
	fmt.Printf("(3) %v\n", err)

	// (4) Creates an Err with parameters and a causal error.
	err = errs.New(FailToDoWithParams{
		Param1: "ABC",
		Param2: 123,
	}, cause)
	fmt.Printf("(4) %v\n", err)
	// Output:
	// (1) github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.FailToDoSomething, file = example_err_test.go, line = 20 }
	// (2) github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.FailToDoWithParams { Param1: ABC, Param2: 123 }, file = example_err_test.go, line = 24 }
	// (3) github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.FailToDoSomething, file = example_err_test.go, line = 33, cause = Causal error }
	// (4) github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.FailToDoWithParams { Param1: ABC, Param2: 123 }, file = example_err_test.go, line = 37, cause = Causal error }
}

func ExampleOk() {
	err := errs.Ok()
	fmt.Printf("err = %v\n", err)
	fmt.Printf("err.IsOk() = %v\n", err.IsOk())
	// Output:
	// err = github.com/sttk/errs.Err { reason = nil, file = , line = 0 }
	// err.IsOk() = true
}

func ExampleErr_Cause() {
	type FailToDoSomething struct{}

	cause := errors.New("Causal error")

	err := errs.New(FailToDoSomething{}, cause)
	fmt.Printf("%v\n", err.Cause())
	// Output:
	// Causal error
}

func ExampleErr_Error() {
	type FailToDoSomething struct {
		Param1 string
		Param2 int
	}

	cause := errors.New("Causal error")

	err := errs.New(FailToDoSomething{
		Param1: "ABC",
		Param2: 123,
	}, cause)
	fmt.Printf("%v\n", err.Error())
	// Output:
	// github.com/sttk/errs.Err { reason = github.com/sttk/errs_test.FailToDoSomething { Param1: ABC, Param2: 123 }, file = example_err_test.go, line = 77, cause = Causal error }
}

func ExampleErr_IsOk() {
	err := errs.Ok()
	fmt.Printf("%v\n", err.IsOk())

	type FailToDoSomething struct{}
	err = errs.New(FailToDoSomething{})
	fmt.Printf("%v\n", err.IsOk())
	// Output:
	// true
	// false
}

func ExampleErr_IsNotOk() {
	err := errs.Ok()
	fmt.Printf("%v\n", err.IsNotOk())

	type FailToDoSomething struct{}
	err = errs.New(FailToDoSomething{})
	fmt.Printf("%v\n", err.IsNotOk())
	// Output:
	// false
	// true
}

func ExampleErr_Reason() {
	type FailToDoSomething struct {
		Param1 string
	}

	err := errs.New(FailToDoSomething{Param1: "value1"})
	switch err.Reason().(type) {
	case FailToDoSomething:
		fmt.Println("The reason of the error is: FailToDoSomething")
		reason := err.Reason().(FailToDoSomething)
		fmt.Printf("The value of reason.Param1 is: %v\n", reason.Param1)
	}

	err = errs.New(&FailToDoSomething{Param1: "value2"})
	switch err.Reason().(type) {
	case *FailToDoSomething:
		fmt.Println("The reason of the error is: *FailToDoSomething")
		reason := err.Reason().(*FailToDoSomething)
		fmt.Printf("The value of reason.Param1 is: %v\n", reason.Param1)
	}
	// Output:
	// The reason of the error is: FailToDoSomething
	// The value of reason.Param1 is: value1
	// The reason of the error is: *FailToDoSomething
	// The value of reason.Param1 is: value2
}

func ExampleErr_Unwrap() {
	type FailToDoSomething struct{}

	cause1 := errors.New("Causal error 1")
	cause2 := errors.New("Causal error 2")

	err := errs.New(FailToDoSomething{}, cause1)

	fmt.Printf("err.Unwrap() = %v\n", err.Unwrap())
	fmt.Printf("errors.Unwrap(err) = %v\n", errors.Unwrap(err))
	fmt.Printf("errors.Is(err, cause1) = %v\n", errors.Is(err, cause1))
	fmt.Printf("errors.Is(err, cause2) = %v\n", errors.Is(err, cause2))
	// Output:
	// err.Unwrap() = Causal error 1
	// errors.Unwrap(err) = Causal error 1
	// errors.Is(err, cause1) = true
	// errors.Is(err, cause2) = false
}

func ExampleErr_IfOkThen() {
	type FailToDoSomething struct{}

	err := errs.Ok()
	err.IfOkThen(func() errs.Err {
		fmt.Println("execute if non error.")
		return errs.Ok()
	})

	err = errs.New(FailToDoSomething{})
	err.IfOkThen(func() errs.Err {
		fmt.Println("not execute if some error.")
		return errs.Ok()
	})
	// Output:
	// execute if non error.
}

func ExampleErr_Line() {
	type FailToDoSomething struct{}

	err := errs.New(FailToDoSomething{})
	fmt.Printf("line = %d\n", err.Line())
	// Output:
	// line = 177
}

func ExampleErr_File() {
	type FailToDoSomething struct{}

	err := errs.New(FailToDoSomething{})
	fmt.Printf("file = %s\n", err.File())
	// Output:
	// file = example_err_test.go
}
