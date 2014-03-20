package foo

import (
	"testing"
)

func TestDoSomething(t *testing.T) {
	DoSomething()
	// Yup, something was done.
	if 0 == 1 {
		t.Fatal("Wow, that is unexpected.")
	}
}
