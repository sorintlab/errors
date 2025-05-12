package errors

import (
	"fmt"
	"strings"
)

func unwrap(err error) []error {
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		err = x.Unwrap()
		if err == nil {
			return nil
		}
		return []error{err}
	case interface{ Unwrap() []error }:
		errs := []error{}
		for _, err := range x.Unwrap() {
			if err == nil {
				continue
			}
			errs = append(errs, err)
		}
		return errs
	default:
		return nil
	}
}

type stackTracer interface {
	StackTrace() StackTrace
}

type errWithStack struct {
	err   error
	msg   string
	stack StackTrace
	level int // level is increased when the wrapped errors are a more than 1
	index int // index is the wrapped errors index (in case of more than one wrapped error)
}

func getStackErr(errCause error, level int, index int) []errWithStack {
	var stackErrs []errWithStack

	stackErr := errWithStack{
		err:   errCause,
		msg:   errCause.Error(),
		level: level,
		index: index,
	}
	if s, ok := errCause.(stackTracer); ok {
		stackErr.stack = s.StackTrace()
	}
	stackErrs = append(stackErrs, stackErr)
	errCauses := unwrap(errCause)

	if len(errCauses) > 1 {
		level++
	}

	for i, errCause := range errCauses {
		stackErrs = append(stackErrs, getStackErr(errCause, level, i)...)
	}

	return stackErrs
}

func PrintErrorDetails(err error) []string {
	stackErrs := getStackErr(err, 0, 0)

	var lines []string
	curLevel := 0
	curIndex := 0
	for _, stackErr := range stackErrs {
		printIndex := false
		if curLevel != stackErr.level {
			printIndex = true
		}
		if curIndex != stackErr.index {
			printIndex = true
		}
		curLevel = stackErr.level
		curIndex = stackErr.index

		indexstr := ""
		indexstrph := ""
		if stackErr.level > 0 {
			indexstr = fmt.Sprintf("[%d] ", stackErr.index)
			indexstrph = strings.Repeat(" ", len(indexstr))
		}
		if !printIndex {
			indexstr = indexstrph
		}

		newlines := []string{}
		if len(stackErr.stack) > 0 {
			frame := stackErr.stack[0]
			newlines = append(newlines, fmt.Sprintf("%s(%T) %s", indexstr, stackErr.err, frame.name()))
			newlines = append(newlines, fmt.Sprintf("\t%s:%d: %s", frame.file(), frame.line(), stackErr.msg))
		} else {
			newlines = append(newlines, fmt.Sprintf("%s(%T) %s", indexstr, stackErr.err, stackErr.msg))
		}
		for _, newline := range newlines {
			newline = strings.Repeat("\t", stackErr.level) + newline
			lines = append(lines, newline)
		}
	}

	return lines
}
