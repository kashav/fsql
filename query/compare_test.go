package query

import (
	"os"
	"testing"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
)

type CmpInput struct {
	comp tokenizer.TokenType
	a, b interface{}
}

type CmpCase struct {
	input    CmpInput
	expected bool
}

func TestCmpAlpha(t *testing.T) {
	cases := []CmpCase{
		{CmpInput{tokenizer.Equals, "a", "a"}, true},
		{CmpInput{tokenizer.Equals, "a", "b"}, false},
		{CmpInput{tokenizer.Equals, "a", "A"}, false},

		{CmpInput{tokenizer.NotEquals, "a", "a"}, false},
		{CmpInput{tokenizer.NotEquals, "a", "b"}, true},
		{CmpInput{tokenizer.NotEquals, "a", "A"}, true},

		{CmpInput{tokenizer.Like, "abc", "%a%"}, true},
		{CmpInput{tokenizer.Like, "aaa", "%b%"}, false},
		{CmpInput{tokenizer.Like, "aaa", "%a"}, true},
		{CmpInput{tokenizer.Like, "abc", "%a"}, false},
		{CmpInput{tokenizer.Like, "abc", "a%"}, true},
		{CmpInput{tokenizer.Like, "cba", "a%"}, false},
		{CmpInput{tokenizer.Like, "a", "a"}, true},
		{CmpInput{tokenizer.Like, "a", "b"}, false},

		{CmpInput{tokenizer.RLike, "a", ".*a.*"}, true},
		{CmpInput{tokenizer.RLike, "a", "^$"}, false},
		{CmpInput{tokenizer.RLike, "", "^$"}, true},
		{CmpInput{tokenizer.RLike, "...", "[\\.]{3}"}, true},
		{CmpInput{tokenizer.RLike, "aaa", "\\s+"}, false},

		{CmpInput{tokenizer.In, "a", map[interface{}]bool{"a": true}}, true},
		{CmpInput{tokenizer.In, "a", map[interface{}]bool{"a": false}}, true},
		{CmpInput{tokenizer.In, "a", map[interface{}]bool{"b": true}}, false},
		{CmpInput{tokenizer.In, "a", map[interface{}]bool{}}, false},
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
	cases := []CmpCase{
		{CmpInput{tokenizer.Equals, int64(1), int64(1)}, true},
		{CmpInput{tokenizer.Equals, int64(1), int64(2)}, false},

		{CmpInput{tokenizer.NotEquals, int64(1), int64(1)}, false},
		{CmpInput{tokenizer.NotEquals, int64(1), int64(2)}, true},

		{CmpInput{tokenizer.GreaterThanEquals, int64(1), int64(1)}, true},
		{CmpInput{tokenizer.GreaterThanEquals, int64(2), int64(1)}, true},
		{CmpInput{tokenizer.GreaterThanEquals, int64(1), int64(2)}, false},

		{CmpInput{tokenizer.GreaterThan, int64(1), int64(1)}, false},
		{CmpInput{tokenizer.GreaterThan, int64(2), int64(1)}, true},
		{CmpInput{tokenizer.GreaterThan, int64(1), int64(2)}, false},

		{CmpInput{tokenizer.LessThanEquals, int64(1), int64(1)}, true},
		{CmpInput{tokenizer.LessThanEquals, int64(2), int64(1)}, false},
		{CmpInput{tokenizer.LessThanEquals, int64(1), int64(2)}, true},

		{CmpInput{tokenizer.LessThan, int64(1), int64(1)}, false},
		{CmpInput{tokenizer.LessThan, int64(2), int64(1)}, false},
		{CmpInput{tokenizer.LessThan, int64(1), int64(2)}, true},

		{CmpInput{tokenizer.In, int64(1), map[interface{}]bool{int64(1): true}}, true},
		{CmpInput{tokenizer.In, int64(1), map[interface{}]bool{int64(1): false}}, true},
		{CmpInput{tokenizer.In, int64(1), map[interface{}]bool{int64(2): true}}, false},
		{CmpInput{tokenizer.In, int64(1), map[interface{}]bool{}}, false},
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
	cases := []CmpCase{
		{CmpInput{tokenizer.Equals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, true},
		{CmpInput{tokenizer.Equals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, false},

		{CmpInput{tokenizer.NotEquals, time.Now().Round(time.Minute), time.Now().Round(time.Minute)}, false},
		{CmpInput{tokenizer.NotEquals, time.Now().Add(time.Hour), time.Now().Add(-1 * time.Hour)}, true},

		{CmpInput{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Round(time.Minute): true}}, true},
		{CmpInput{tokenizer.In, time.Now().Round(time.Minute), map[interface{}]bool{time.Now().Add(time.Hour): true}}, false},
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
		actual := cmpMode(c.input.comp, c.input.file, c.input.fileType)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.file, c.input.fileType, c.expected, actual)
		}
	}
}
