package errs

import (
	"fmt"
	"strings"
)

// New will return a new error.
func New(msg interface{}) error {
	return &basicError{msg}
}

// Eat the return value of any function call and return the error
func Eat(v interface{}, e error) error {
	return e
}

// Pair two errors together. If at least one error is nil, it will
// return the one that is not nil. If they are both non-nil it will
// return a new error pair which is itself an error.
func Pair(first, second error) error {
	if first == nil || second == nil {
		if second != nil {
			return second
		}
		return first
	}
	return &errpair{first, second}
}

// Chain will chain a list of errors together
func Chain(es ...error) error {
	var errors []error
	switch l := len(es); l {
	case 0:
		return nil
	case 1:
		if es[0] != nil {
			return es[0]
		}
		return nil
	default:
		errors = make([]error, 0, l)
	}

	for _, e := range es {
		if e != nil {
			errors = append(errors, e)
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return &errlist{errors}
}

type errpair struct {
	first, second error
}

func (e *errpair) Error() string {
	return fmt.Sprintf("%s; %s", e.first.Error(), e.second.Error())
}

type errlist struct {
	errs []error
}

func (el *errlist) Error() string {
	lst := make([]string, len(el.errs))
	for i, e := range el.errs {
		lst[i] = e.Error()
	}
	return strings.Join(lst, "; ")
}

type basicError struct {
	msg interface{}
}

func (e *basicError) Error() string {
	return fmt.Sprintf("%v", e.msg)
}
