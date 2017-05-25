package transform

import "testing"

func TestErrNotImplemented(t *testing.T) {
	err := &ErrNotImplemented{"n", "a"}
	expected := "function N is not implemented for attribute a"
	actual := err.Error()
	if expected != actual {
		t.Fatalf("\nExpected: %s\n     Got: %s", expected, actual)
	}
}

func TestErrUnsupportedFormat(t *testing.T) {
	err := &ErrUnsupportedFormat{"f", "a"}
	expected := "unsupported format type f for attribute a"
	actual := err.Error()
	if expected != actual {
		t.Fatalf("\nExpected: %s\n     Got: %s", expected, actual)
	}

}
