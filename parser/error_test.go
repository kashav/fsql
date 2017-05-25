package parser

import (
	"testing"

	"github.com/kshvmdn/fsql/tokenizer"
)

func TestErrUnexpectedToken(t *testing.T) {
	err := &ErrUnexpectedToken{
		Actual:   tokenizer.Select,
		Expected: tokenizer.Where,
	}
	expected := "expected where; got select"
	actual := err.Error()
	if expected != actual {
		t.Fatalf("\nExpected: %s\n     Got: %s", expected, actual)
	}
}

func TestErrUnknownTokent(t *testing.T) {
	err := &ErrUnknownToken{"r"}
	expected := "unknown token: r"
	actual := err.Error()
	if expected != actual {
		t.Fatalf("\nExpected: %s\n     Got: %s", expected, actual)
	}
}
