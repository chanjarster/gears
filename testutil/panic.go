package testutil

import (
	"reflect"
	"testing"
)

func ShouldPanic(t *testing.T, title string, f func()) {
	defer func() { recover() }()
	f()
	t.Errorf("%s did not panicked", title)
}

func ExpectPanic(t *testing.T, title string, expected interface{}, f func()) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("%s panic = [%v], want [%v]", title, r, expected)
		}
	}()
	f()
	t.Errorf("%s did not panicked", title)
}
func ShouldNotPanic(t *testing.T, title string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s panicked, message: %v", title, r)
		}
	}()
	f()
}
