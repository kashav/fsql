package query

import (
	"os"
	"testing"
	"time"

	"github.com/kshvmdn/fsql/tokenizer"
)

func TestCmpAlpha(t *testing.T) {
	type Input struct {
		comp tokenizer.TokenType
		a, b interface{}
	}

	type Case struct {
		input    Input
		expected bool
	}

	cases := []Case{
		{
			input:    Input{comp: tokenizer.Equals, a: "a", b: "a"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.Equals, a: "a", b: "b"},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.Equals, a: "a", b: "A"},
			expected: false,
		},

		{
			input:    Input{comp: tokenizer.NotEquals, a: "a", b: "a"},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.NotEquals, a: "a", b: "b"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.NotEquals, a: "a", b: "A"},
			expected: true,
		},

		{
			input:    Input{comp: tokenizer.Like, a: "abc", b: "%a%"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.Like, a: "aaa", b: "%b%"},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.Like, a: "aaa", b: "%a"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.Like, a: "abc", b: "%a"},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.Like, a: "abc", b: "a%"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.Like, a: "cba", b: "a%"},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.Like, a: "a", b: "a"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.Like, a: "a", b: "b"},
			expected: false,
		},

		{
			input:    Input{comp: tokenizer.RLike, a: "a", b: ".*a.*"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.RLike, a: "a", b: "^$"},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.RLike, a: "", b: "^$"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.RLike, a: "...", b: "[\\.]{3}"},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.RLike, a: "aaa", b: "\\s+"},
			expected: false,
		},

		{
			input:    Input{comp: tokenizer.In, a: "a", b: map[interface{}]bool{"a": true}},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.In, a: "a", b: map[interface{}]bool{"a": false}},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.In, a: "a", b: map[interface{}]bool{"b": true}},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.In, a: "a", b: map[interface{}]bool{}},
			expected: false,
		},
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
	type Input struct {
		comp tokenizer.TokenType
		a, b interface{}
	}

	type Case struct {
		input    Input
		expected bool
	}

	cases := []Case{
		{
			input:    Input{comp: tokenizer.Equals, a: int64(1), b: int64(1)},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.Equals, a: int64(1), b: int64(2)},
			expected: false,
		},

		{
			input:    Input{comp: tokenizer.NotEquals, a: int64(1), b: int64(1)},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.NotEquals, a: int64(1), b: int64(2)},
			expected: true,
		},

		{
			input:    Input{comp: tokenizer.GreaterThanEquals, a: int64(1), b: int64(1)},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.GreaterThanEquals, a: int64(2), b: int64(1)},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.GreaterThanEquals, a: int64(1), b: int64(2)},
			expected: false,
		},

		{
			input:    Input{comp: tokenizer.GreaterThan, a: int64(1), b: int64(1)},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.GreaterThan, a: int64(2), b: int64(1)},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.GreaterThan, a: int64(1), b: int64(2)},
			expected: false,
		},

		{
			input:    Input{comp: tokenizer.LessThanEquals, a: int64(1), b: int64(1)},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.LessThanEquals, a: int64(2), b: int64(1)},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.LessThanEquals, a: int64(1), b: int64(2)},
			expected: true,
		},

		{
			input:    Input{comp: tokenizer.LessThan, a: int64(1), b: int64(1)},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.LessThan, a: int64(2), b: int64(1)},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.LessThan, a: int64(1), b: int64(2)},
			expected: true,
		},

		{
			input:    Input{comp: tokenizer.In, a: int64(1), b: map[interface{}]bool{int64(1): true}},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.In, a: int64(1), b: map[interface{}]bool{int64(1): false}},
			expected: true,
		},
		{
			input:    Input{comp: tokenizer.In, a: int64(1), b: map[interface{}]bool{int64(2): true}},
			expected: false,
		},
		{
			input:    Input{comp: tokenizer.In, a: int64(1), b: map[interface{}]bool{}},
			expected: false,
		},
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
	type Input struct {
		comp tokenizer.TokenType
		a, b interface{}
	}

	type Case struct {
		input    Input
		expected bool
	}

	cases := []Case{
		{
			input: Input{
				comp: tokenizer.Equals,
				a:    time.Now().Round(time.Minute),
				b:    time.Now().Round(time.Minute),
			},
			expected: true,
		},
		{
			input: Input{
				comp: tokenizer.Equals,
				a:    time.Now().Add(time.Hour),
				b:    time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},

		{
			input: Input{
				comp: tokenizer.NotEquals,
				a:    time.Now().Round(time.Minute),
				b:    time.Now().Round(time.Minute),
			},
			expected: false,
		},
		{
			input: Input{
				comp: tokenizer.NotEquals,
				a:    time.Now().Add(time.Hour),
				b:    time.Now().Add(-1 * time.Hour),
			},
			expected: true,
		},

		{
			input: Input{
				comp: tokenizer.In,
				a:    time.Now().Round(time.Minute),
				b:    map[interface{}]bool{time.Now().Round(time.Minute): true},
			},
			expected: true,
		},
		{
			input: Input{
				comp: tokenizer.In,
				a:    time.Now().Round(time.Minute),
				b:    map[interface{}]bool{time.Now().Add(time.Hour): true},
			},
			expected: false,
		},
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
	type Input struct {
		comp     tokenizer.TokenType
		file     os.FileInfo
		fileType interface{}
	}

	type Case struct {
		input    Input
		expected bool
	}

	// TODO: Complete these cases.
	cases := []Case{}

	for _, c := range cases {
		actual := cmpMode(c.input.comp, c.input.file, c.input.fileType)
		if actual != c.expected {
			t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
				c.input.comp, c.input.file, c.input.fileType, c.expected, actual)
		}
	}
}
