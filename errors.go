package errors

import (
	"fmt"
	"io"
	"unsafe"
)

type werror struct {
	cause error
	msg   string
	*Stack
}

func (w *werror) Error() string {
	if w.cause == nil {
		return w.msg
	}
	if w.msg != "" {
		return w.msg + ": " + w.cause.Error()
	} else {
		return w.cause.Error()
	}
}

func (w *werror) Format(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, w.Error())
}

func (w *werror) Unwrap() error { return w.cause }

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) error {
	return &werror{
		msg:   message,
		Stack: Callers(0),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &werror{
		msg:   fmt.Sprintf(format, args...),
		Stack: Callers(0),
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &werror{
		err,
		"",
		Callers(0),
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &werror{
		err,
		message,
		Callers(0),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &werror{
		err,
		fmt.Sprintf(format, args...),
		Callers(0),
	}
}

type joinError struct {
	errs []error
	*Stack
}

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs:  make([]error, 0, n),
		Stack: Callers(0),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

func (e *joinError) Error() string {
	// Since Join returns nil if every value in errs is nil,
	// e.errs cannot be empty.
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	b := []byte(e.errs[0].Error())
	for _, err := range e.errs[1:] {
		b = append(b, '\n')
		b = append(b, err.Error()...)
	}
	// At this point, b has at least one byte '\n'.
	return unsafe.String(&b[0], len(b))
}

func (e *joinError) Unwrap() []error {
	return e.errs
}
