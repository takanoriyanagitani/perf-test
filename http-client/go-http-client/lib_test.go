package http_client_perf

import (
	"testing"
)

func assertEqNew[T any](comp func(a, b T) (same bool)) func(a, b T) func(*testing.T) {
	return func(a, b T) func(*testing.T) {
		return func(t *testing.T) {
			var same bool = comp(a, b)
			if same {
				return
			}

			t.Errorf("Unexpected value got.\n")
			t.Errorf("Expected: %v\n", b)
			t.Fatalf("Got:      %v\n", a)
		}
	}
}

func assertEq[T comparable](a, b T) func(*testing.T) {
	return assertEqNew(
		func(a, b T) (same bool) { return a == b },
	)(a, b)
}

func assertTrue(b bool) func(*testing.T) {
	return func(t *testing.T) {
		if !b {
			t.Fatalf("Must be true: %v", b)
		}
	}
}
