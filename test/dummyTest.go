package test

import "testing"

func TestForward(t *testing.T) {
	// System under testing

	got := Add(4, 6)
	want := 10
	// test
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestStateChange(t *testing.T) {
	// System under testing

	got := Add(4, 6)
	want := 10
	// test
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestForward(t *testing.T) {
	// System under testing

	got := Add(4, 6)
	want := 10
	// test
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
