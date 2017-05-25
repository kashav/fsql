package query

import (
	"os"
	"testing"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
)

type EvaluateInput struct {
	comp tokenizer.TokenType
	a, b interface{}
}

type EvaluateCase struct {
	input    EvaluateInput
	expected bool
}

func TestEvalAlpha(t *testing.T) {
	cases := []EvaluateCase{
		{EvaluateInput{tokenizer.Equals, "a", "a"}, true},
		{EvaluateInput{tokenizer.Equals, "a", "b"}, false},
		{EvaluateInput{tokenizer.Equals, "a", "A"}, false},

		{EvaluateInput{tokenizer.NotEquals, "a", "a"}, false},
		{EvaluateInput{tokenizer.NotEquals, "a", "b"}, true},
		{EvaluateInput{tokenizer.NotEquals, "a", "A"}, true},

		{EvaluateInput{tokenizer.Like, "abc", "%a%"}, true},
		{EvaluateInput{tokenizer.Like, "aaa", "%b%"}, false},
		{EvaluateInput{tokenizer.Like, "aaa", "%a"}, true},
		{EvaluateInput{tokenizer.Like, "abc", "%a"}, false},
		{EvaluateInput{tokenizer.Like, "abc", "a%"}, true},
		{EvaluateInput{tokenizer.Like, "cba", "a%"}, false},
		{EvaluateInput{tokenizer.Like, "a", "a"}, true},
		{EvaluateInput{tokenizer.Like, "a", "b"}, false},

		{EvaluateInput{tokenizer.RLike, "a", ".*a.*"}, true},
		{EvaluateInput{tokenizer.RLike, "a", "^$"}, false},
		{EvaluateInput{tokenizer.RLike, "", "^$"}, true},
		{EvaluateInput{tokenizer.RLike, "...", "[\\.]{3}"}, true},
		{EvaluateInput{tokenizer.RLike, "aaa", "\\s+"}, false},

		{EvaluateInput{tokenizer.In, "a", map[interface{}]bool{"a": true}}, true},
		{EvaluateInput{tokenizer.In, "a", map[interface{}]bool{"a": false}}, true},
		{EvaluateInput{tokenizer.In, "a", map[interface{}]bool{"b": true}}, false},
		{EvaluateInput{tokenizer.In, "a", map[interface{}]bool{}}, false},
	}

	for _, c := range cases {
		actual := evalAlpha(c.input.comp, c.input.a, c.input.b)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.a, c.input.b, c.expected, actual)
		}
	}
}

func TestEvalNumeric(t *testing.T) {
	cases := []EvaluateCase{
		{EvaluateInput{tokenizer.Equals, int64(1), int64(1)}, true},
		{EvaluateInput{tokenizer.Equals, int64(1), int64(2)}, false},

		{EvaluateInput{tokenizer.NotEquals, int64(1), int64(1)}, false},
		{EvaluateInput{tokenizer.NotEquals, int64(1), int64(2)}, true},

		{EvaluateInput{tokenizer.GreaterThanEquals, int64(1), int64(1)}, true},
		{EvaluateInput{tokenizer.GreaterThanEquals, int64(2), int64(1)}, true},
		{EvaluateInput{tokenizer.GreaterThanEquals, int64(1), int64(2)}, false},

		{EvaluateInput{tokenizer.GreaterThan, int64(1), int64(1)}, false},
		{EvaluateInput{tokenizer.GreaterThan, int64(2), int64(1)}, true},
		{EvaluateInput{tokenizer.GreaterThan, int64(1), int64(2)}, false},

		{EvaluateInput{tokenizer.LessThanEquals, int64(1), int64(1)}, true},
		{EvaluateInput{tokenizer.LessThanEquals, int64(2), int64(1)}, false},
		{EvaluateInput{tokenizer.LessThanEquals, int64(1), int64(2)}, true},

		{EvaluateInput{tokenizer.LessThan, int64(1), int64(1)}, false},
		{EvaluateInput{tokenizer.LessThan, int64(2), int64(1)}, false},
		{EvaluateInput{tokenizer.LessThan, int64(1), int64(2)}, true},

		{EvaluateInput{tokenizer.In, int64(1), map[interface{}]bool{int64(1): true}}, true},
		{EvaluateInput{tokenizer.In, int64(1), map[interface{}]bool{int64(1): false}}, true},
		{EvaluateInput{tokenizer.In, int64(1), map[interface{}]bool{int64(2): true}}, false},
		{EvaluateInput{tokenizer.In, int64(1), map[interface{}]bool{}}, false},
	}

	for _, c := range cases {
		actual := evalNumeric(c.input.comp, c.input.a, c.input.b)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.a, c.input.b, c.expected, actual)
		}
	}
}

func TestEvalTime(t *testing.T) {
	cases := []EvaluateCase{
		{EvaluateInput{tokenizer.Equals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, true},
		{EvaluateInput{tokenizer.Equals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, false},

		{EvaluateInput{tokenizer.NotEquals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, false},
		{EvaluateInput{tokenizer.NotEquals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, true},

		{EvaluateInput{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Round(time.Minute): true}}, true},
		{EvaluateInput{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Add(time.Hour): true}}, false},
	}

	for _, c := range cases {
		actual := evalTime(c.input.comp, c.input.a, c.input.b)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.a, c.input.b, c.expected, actual)
		}
	}
}

func TestEvalFile(t *testing.T) {
	type FileInput struct {
		comp     tokenizer.TokenType
		file     os.FileInfo
		fileType interface{}
	}

	type FileCase struct {
		input    FileInput
		expected bool
	}

	// TODO: Complete these cases.
	cases := []FileCase{}

	for _, c := range cases {
		actual := evalFile(c.input.comp, c.input.file, c.input.fileType)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.file, c.input.fileType, c.expected, actual)
		}
	}
}
