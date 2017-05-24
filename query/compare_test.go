package query

import (
	"os"
	"testing"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
)

type CompareInput struct {
	comp tokenizer.TokenType
	a, b interface{}
}

type CompareCase struct {
	input    CompareInput
	expected bool
}

func TestCmpAlpha(t *testing.T) {
	cases := []CompareCase{
		{CompareInput{tokenizer.Equals, "a", "a"}, true},
		{CompareInput{tokenizer.Equals, "a", "b"}, false},
		{CompareInput{tokenizer.Equals, "a", "A"}, false},

		{CompareInput{tokenizer.NotEquals, "a", "a"}, false},
		{CompareInput{tokenizer.NotEquals, "a", "b"}, true},
		{CompareInput{tokenizer.NotEquals, "a", "A"}, true},

		{CompareInput{tokenizer.Like, "abc", "%a%"}, true},
		{CompareInput{tokenizer.Like, "aaa", "%b%"}, false},
		{CompareInput{tokenizer.Like, "aaa", "%a"}, true},
		{CompareInput{tokenizer.Like, "abc", "%a"}, false},
		{CompareInput{tokenizer.Like, "abc", "a%"}, true},
		{CompareInput{tokenizer.Like, "cba", "a%"}, false},
		{CompareInput{tokenizer.Like, "a", "a"}, true},
		{CompareInput{tokenizer.Like, "a", "b"}, false},

		{CompareInput{tokenizer.RLike, "a", ".*a.*"}, true},
		{CompareInput{tokenizer.RLike, "a", "^$"}, false},
		{CompareInput{tokenizer.RLike, "", "^$"}, true},
		{CompareInput{tokenizer.RLike, "...", "[\\.]{3}"}, true},
		{CompareInput{tokenizer.RLike, "aaa", "\\s+"}, false},

		{CompareInput{tokenizer.In, "a", map[interface{}]bool{"a": true}}, true},
		{CompareInput{tokenizer.In, "a", map[interface{}]bool{"a": false}}, true},
		{CompareInput{tokenizer.In, "a", map[interface{}]bool{"b": true}}, false},
		{CompareInput{tokenizer.In, "a", map[interface{}]bool{}}, false},
	}

	for _, c := range cases {
		actual := cmpAlpha(c.input.comp, c.input.a, c.input.b)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.a, c.input.b, c.expected, actual)
		}
	}
}

func TestCmpNumeric(t *testing.T) {
	cases := []CompareCase{
		{CompareInput{tokenizer.Equals, int64(1), int64(1)}, true},
		{CompareInput{tokenizer.Equals, int64(1), int64(2)}, false},

		{CompareInput{tokenizer.NotEquals, int64(1), int64(1)}, false},
		{CompareInput{tokenizer.NotEquals, int64(1), int64(2)}, true},

		{CompareInput{tokenizer.GreaterThanEquals, int64(1), int64(1)}, true},
		{CompareInput{tokenizer.GreaterThanEquals, int64(2), int64(1)}, true},
		{CompareInput{tokenizer.GreaterThanEquals, int64(1), int64(2)}, false},

		{CompareInput{tokenizer.GreaterThan, int64(1), int64(1)}, false},
		{CompareInput{tokenizer.GreaterThan, int64(2), int64(1)}, true},
		{CompareInput{tokenizer.GreaterThan, int64(1), int64(2)}, false},

		{CompareInput{tokenizer.LessThanEquals, int64(1), int64(1)}, true},
		{CompareInput{tokenizer.LessThanEquals, int64(2), int64(1)}, false},
		{CompareInput{tokenizer.LessThanEquals, int64(1), int64(2)}, true},

		{CompareInput{tokenizer.LessThan, int64(1), int64(1)}, false},
		{CompareInput{tokenizer.LessThan, int64(2), int64(1)}, false},
		{CompareInput{tokenizer.LessThan, int64(1), int64(2)}, true},

		{CompareInput{tokenizer.In, int64(1), map[interface{}]bool{int64(1): true}}, true},
		{CompareInput{tokenizer.In, int64(1), map[interface{}]bool{int64(1): false}}, true},
		{CompareInput{tokenizer.In, int64(1), map[interface{}]bool{int64(2): true}}, false},
		{CompareInput{tokenizer.In, int64(1), map[interface{}]bool{}}, false},
	}

	for _, c := range cases {
		actual := cmpNumeric(c.input.comp, c.input.a, c.input.b)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.a, c.input.b, c.expected, actual)
		}
	}
}

func TestCmpTime(t *testing.T) {
	cases := []CompareCase{
		{CompareInput{tokenizer.Equals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, true},
		{CompareInput{tokenizer.Equals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, false},

		{CompareInput{tokenizer.NotEquals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, false},
		{CompareInput{tokenizer.NotEquals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, true},

		{CompareInput{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Round(time.Minute): true}}, true},
		{CompareInput{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Add(time.Hour): true}}, false},
	}

	for _, c := range cases {
		actual := cmpTime(c.input.comp, c.input.a, c.input.b)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.a, c.input.b, c.expected, actual)
		}
	}
}

func TestCmpFile(t *testing.T) {
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
		actual := cmpFile(c.input.comp, c.input.file, c.input.fileType)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.file, c.input.fileType, c.expected, actual)
		}
	}
}
