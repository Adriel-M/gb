package assert

import (
	"github.com/Adriel-M/gb/gb/post"
	"testing"
)

type Assert struct {
	T *testing.T
}

func (a Assert) StringEqual(actual string, expected string) {
	if actual != expected {
		a.T.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

func (a Assert) BoolEqual(actual bool, expected bool) {
	if actual != expected {
		a.T.Fatalf("Expected %t, but got %t", expected, actual)
	}
}

func (a Assert) IntEqual(actual int, expected int) {
	if actual != expected {
		a.T.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func (a Assert) PostAddressEqual(actual *post.Post, expected *post.Post) {
	if actual != expected {
		a.T.Fatalf("Expected %#v, but got %#v", expected, actual)
	}
}
