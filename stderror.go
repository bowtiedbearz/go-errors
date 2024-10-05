package errors

import (
	"fmt"
	"io"
)

// The standard error which supports code, cause, and stacktrace
type StdError struct {
	*stack
	msg   string
	cause error
	code  string
}

// Error returns the error message
func (e *StdError) Error() string {
	return e.msg
}

// Code returns the error code
func (e *StdError) Code() string {
	return e.code
}

// Cause returns the underlying cause of the error, if one
// exists. An error value has a cause if it implements the following
// interface:
func (e *StdError) Cause() error {
	return e.cause
}

// Unwrap returns the underlying cause of the error, if one
// exists.
func (e *StdError) Unwrap() error {
	return e.cause
}

// Is returns true if the error is the same as the target error
func (e *StdError) Is(target error) bool {
	if n, ok := target.(*StdError); ok {
		return e.code == n.code
	}

	return false
}

// Sets the underlying cause of the error
func (e *StdError) WithCause(err error) *StdError {
	e.cause = err
	return e
}

// Sets the error message
func (e *StdError) WithMessage(msg string) *StdError {
	e.msg = msg
	return e
}

// Sets the error message with a formatted string
func (e *StdError) WithMessageF(msg string, args ...interface{}) *StdError {
	e.msg = fmt.Sprintf(msg, args...)
	return e
}

// Formats the error. The verb is used to determine how to format the error
// '%+v' will print the error code, message, cause, and stacktrace
// '%v' will print the error code, message, and stacktrace
// '%s' will print the error message
// '%q' will print the error message using fmt.Sprintf("%q", e.Error())
func (e *StdError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.code+": "+e.msg)
			if e.cause != nil {
				s.Write([]byte{'\n'})
				fmt.Fprintf(s, "%+v", e.Cause())
			}

			e.stack.Format(s, verb)
			return
		} else {
			io.WriteString(s, e.code)
			io.WriteString(s, ": ")
			io.WriteString(s, e.Error())
			e.stack.Format(s, verb)
		}

	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}
