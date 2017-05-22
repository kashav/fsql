package query

import (
	"os"
	"testing"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
)

type Input struct {
	comp tokenizer.TokenType
	a, b interface{}
}

type Case struct {
	input    Input
	expected bool
}

func TestCmpAlpha(t *testing.T) {
	cases := []Case{
		{Input{tokenizer.Equals, "a", "a"}, true},
		{Input{tokenizer.Equals, "a", "b"}, false},
		{Input{tokenizer.Equals, "a", "A"}, false},

		{Input{tokenizer.NotEquals, "a", "a"}, false},
		{Input{tokenizer.NotEquals, "a", "b"}, true},
		{Input{tokenizer.NotEquals, "a", "A"}, true},

		{Input{tokenizer.Like, "abc", "%a%"}, true},
		{Input{tokenizer.Like, "aaa", "%b%"}, false},
		{Input{tokenizer.Like, "aaa", "%a"}, true},
		{Input{tokenizer.Like, "abc", "%a"}, false},
		{Input{tokenizer.Like, "abc", "a%"}, true},
		{Input{tokenizer.Like, "cba", "a%"}, false},
		{Input{tokenizer.Like, "a", "a"}, true},
		{Input{tokenizer.Like, "a", "b"}, false},

		{Input{tokenizer.RLike, "a", ".*a.*"}, true},
		{Input{tokenizer.RLike, "a", "^$"}, false},
		{Input{tokenizer.RLike, "", "^$"}, true},
		{Input{tokenizer.RLike, "...", "[\\.]{3}"}, true},
		{Input{tokenizer.RLike, "aaa", "\\s+"}, false},

		{Input{tokenizer.In, "a", map[interface{}]bool{"a": true}}, true},
		{Input{tokenizer.In, "a", map[interface{}]bool{"a": false}}, true},
		{Input{tokenizer.In, "a", map[interface{}]bool{"b": true}}, false},
		{Input{tokenizer.In, "a", map[interface{}]bool{}}, false},
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
	cases := []Case{
		{Input{tokenizer.Equals, int64(1), int64(1)}, true},
		{Input{tokenizer.Equals, int64(1), int64(2)}, false},

		{Input{tokenizer.NotEquals, int64(1), int64(1)}, false},
		{Input{tokenizer.NotEquals, int64(1), int64(2)}, true},

		{Input{tokenizer.GreaterThanEquals, int64(1), int64(1)}, true},
		{Input{tokenizer.GreaterThanEquals, int64(2), int64(1)}, true},
		{Input{tokenizer.GreaterThanEquals, int64(1), int64(2)}, false},

		{Input{tokenizer.GreaterThan, int64(1), int64(1)}, false},
		{Input{tokenizer.GreaterThan, int64(2), int64(1)}, true},
		{Input{tokenizer.GreaterThan, int64(1), int64(2)}, false},

		{Input{tokenizer.LessThanEquals, int64(1), int64(1)}, true},
		{Input{tokenizer.LessThanEquals, int64(2), int64(1)}, false},
		{Input{tokenizer.LessThanEquals, int64(1), int64(2)}, true},

		{Input{tokenizer.LessThan, int64(1), int64(1)}, false},
		{Input{tokenizer.LessThan, int64(2), int64(1)}, false},
		{Input{tokenizer.LessThan, int64(1), int64(2)}, true},

		{Input{tokenizer.In, int64(1), map[interface{}]bool{int64(1): true}}, true},
		{Input{tokenizer.In, int64(1), map[interface{}]bool{int64(1): false}}, true},
		{Input{tokenizer.In, int64(1), map[interface{}]bool{int64(2): true}}, false},
		{Input{tokenizer.In, int64(1), map[interface{}]bool{}}, false},
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
	cases := []Case{
		{Input{tokenizer.Equals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, true},
		{Input{tokenizer.Equals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, false},

		{Input{tokenizer.NotEquals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, false},
		{Input{tokenizer.NotEquals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, true},

		{Input{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Round(time.Minute): true}}, true},
		{Input{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Add(time.Hour): true}}, false},
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
