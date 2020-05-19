package errs

import (
	"testing"
)

func TestErrPair(t *testing.T) {
	tt := []struct {
		err error
		exp string
	}{
		{Pair(New("one"), New("two")), "one; two"},
		{Pair(New("one"), nil), "one"},
		{Pair(nil, New("two")), "two"},
	}
	for _, tc := range tt {
		if tc.err.Error() != tc.exp {
			t.Error("errpair gave wrong result")
		}
	}
	err := Pair(nil, nil)
	if err != nil {
		t.Error("a pair of nil errors should result in one nil error")
	}
}

func TestEat(t *testing.T) {
	f := func() (string, error) {
		return "testing", New("error")
	}
	if err := Eat(f()); err.Error() != "error" {
		t.Error("got wrong error message")
	}
}

func TestChain(t *testing.T) {
	var err error
	err = Chain()
	if err != nil {
		t.Error("empty chain should return nil")
	}
	if err = Chain(nil); err != nil {
		t.Error("should return nil")
	}
	err = New("test error")
	e := Chain(err)
	if e != err {
		t.Error("one error in the chain should return that error")
	}

	err = Chain(New("one"), New("two"), New("three"), New("four"))
	if err.Error() != "one; two; three; four" {
		t.Error("got wrong error message")
	}
	err = Chain(New("one"), nil, New("three"))
	if err.Error() != "one; three" {
		t.Error("got wrong error message")
	}
	err = Chain(nil, nil, nil, nil, nil)
	if err != nil {
		t.Error("all nil chain should return nil")
	}
}
