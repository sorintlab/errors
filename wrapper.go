package errors

import (
	"fmt"
)

// Wrapper error is an helper error type that (optionally) wraps an error and adds stack information starting at the frame where the error has been created
// It's meant to be embedded in custom errors to avoid the need to redefine the Error, Unwrap and StackTrace methods.
//
// Example usage:
//
//	type CustomError struct {
//		*WrapperError
//	}
//
//	func NewCustomError(err error) error {
//		return &CustomError{
//			util.NewWrapperError(err, util.WithWrapperErrorMsg("connection error")),
//		}
//	}
//
// Create the error
//
//	if err != nil {
//		return NewCustomError(err)
//	}
//
// Create the error without wrapping another error
//
//	return NewCustomError(nil)
//
// Detect error type
//
//	var werr *CustomError
//	if errors.As(err, &werr) {
//		fmt.Println("this is a CustomError")
//	}
type WrapperError struct {
	err error
	msg string

	stack *Stack
}

func NewWrapperError(err error, options ...WrapperErrorOption) *WrapperError {
	werr := &WrapperError{err: err}

	for _, opt := range options {
		opt(werr)
	}

	if werr.stack == nil {
		// skip one frame by default if the error is used as in the example
		werr.stack = Callers(1)
	}

	return werr
}

type WrapperErrorOption func(e *WrapperError)

func WithWrapperErrorMsgf(format string, args ...any) WrapperErrorOption {
	return func(e *WrapperError) {
		e.msg = fmt.Sprintf(format, args...)
	}
}

func WithWrapperErrorMsg(a ...any) WrapperErrorOption {
	return func(e *WrapperError) {
		e.msg = fmt.Sprint(a...)
	}
}

func WithWrapperErrorCallerDepth(depth int) WrapperErrorOption {
	return func(e *WrapperError) {
		e.stack = Callers(depth + 1)
	}
}

func (w *WrapperError) Error() string {
	if w.err == nil {
		return w.msg
	}
	if w.msg != "" {
		return w.msg + ": " + w.err.Error()
	} else {
		return w.err.Error()
	}
}

func (w *WrapperError) Unwrap() error { return w.err }

func (w *WrapperError) StackTrace() StackTrace {
	return w.stack.StackTrace()
}
