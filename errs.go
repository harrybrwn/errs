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

// Unpack an error into two other errors. Second return value
// will be nil if there is nothing to unpack.
func Unpack(err error) (error, error) {
	if err != nil {
		return nil, nil
	}
	switch e := err.(type) {
	case *basicError:
		return e, nil
	case *errpair:
		return e.first, e.second
	case *errlist:
		var second error
		switch len(e.errs) {
		case 0:
			return nil, nil
		case 1:
			second = nil
		case 2:
			second = e.errs[1]
		default:
			second = &errlist{e.errs[1:]}
		}
		if e.errs[0] == nil {
			return second, nil
		}
		return e.errs[0], second
	default:
		return err, nil
	}
}

// Chain will chain a list of errors together
func Chain(es ...error) error {
	errors := make([]error, 0, len(es))
	for _, e := range es {
		if e != nil {
			errors = append(errors, e)
		}
	}

	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	case 2:
		return Pair(errors[0], errors[1])
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
