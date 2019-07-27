package evaluate

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kashav/fsql/tokenizer"
)

func TestCmpAlpha(t *testing.T) {
	type Input struct {
		o    Opts
		a, b interface{}
	}

	type Expected struct {
		result bool
		err    error
	}

	type Case struct {
		input    Input
		expected Expected
	}

	// TODO: Test for errors.
	cases := []Case{
		{
			input:    Input{o: Opts{Operator: tokenizer.Equals}, a: "a", b: "a"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Equals}, a: "a", b: "b"},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Equals}, a: "a", b: "A"},
			expected: Expected{result: false, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.NotEquals}, a: "a", b: "a"},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.NotEquals}, a: "a", b: "b"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.NotEquals}, a: "a", b: "A"},
			expected: Expected{result: true, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "abc", b: "%a%"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "aaa", b: "%b%"},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "aaa", b: "%a"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "abc", b: "%a"},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "abc", b: "a%"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "cba", b: "a%"},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "a", b: "a"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Like}, a: "a", b: "b"},
			expected: Expected{result: false, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.RLike}, a: "a", b: ".*a.*"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.RLike}, a: "a", b: "^$"},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.RLike}, a: "", b: "^$"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.RLike}, a: "...", b: "[\\.]{3}"},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.RLike}, a: "aaa", b: "\\s+"},
			expected: Expected{result: false, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: "a", b: map[interface{}]bool{"a": true}},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: "a", b: map[interface{}]bool{"a": false}},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: "a", b: map[interface{}]bool{"b": true}},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: "a", b: map[interface{}]bool{}},
			expected: Expected{result: false, err: nil},
		},
	}

	for _, c := range cases {
		actual, err := cmpAlpha(&c.input.o, c.input.a, c.input.b)
		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.result, actual) {
				t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
					c.input.o.Operator, c.input.a, c.input.b, c.expected.result, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestCmpNumeric(t *testing.T) {
	type Input struct {
		o    Opts
		a, b interface{}
	}

	type Expected struct {
		result bool
		err    error
	}

	type Case struct {
		input    Input
		expected Expected
	}

	// TODO: Test for errors.
	cases := []Case{
		{
			input:    Input{o: Opts{Operator: tokenizer.Equals}, a: int64(1), b: int64(1)},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.Equals}, a: int64(1), b: int64(2)},
			expected: Expected{result: false, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.NotEquals}, a: int64(1), b: int64(1)},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.NotEquals}, a: int64(1), b: int64(2)},
			expected: Expected{result: true, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.GreaterThanEquals}, a: int64(1), b: int64(1)},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.GreaterThanEquals}, a: int64(2), b: int64(1)},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.GreaterThanEquals}, a: int64(1), b: int64(2)},
			expected: Expected{result: false, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.GreaterThan}, a: int64(1), b: int64(1)},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.GreaterThan}, a: int64(2), b: int64(1)},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.GreaterThan}, a: int64(1), b: int64(2)},
			expected: Expected{result: false, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.LessThanEquals}, a: int64(1), b: int64(1)},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.LessThanEquals}, a: int64(2), b: int64(1)},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.LessThanEquals}, a: int64(1), b: int64(2)},
			expected: Expected{result: true, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.LessThan}, a: int64(1), b: int64(1)},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.LessThan}, a: int64(2), b: int64(1)},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.LessThan}, a: int64(1), b: int64(2)},
			expected: Expected{result: true, err: nil},
		},

		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: int64(1), b: map[interface{}]bool{int64(1): true}},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: int64(1), b: map[interface{}]bool{int64(1): false}},
			expected: Expected{result: true, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: int64(1), b: map[interface{}]bool{int64(2): true}},
			expected: Expected{result: false, err: nil},
		},
		{
			input:    Input{o: Opts{Operator: tokenizer.In}, a: int64(1), b: map[interface{}]bool{}},
			expected: Expected{result: false, err: nil},
		},
	}

	for _, c := range cases {
		actual, err := cmpNumeric(&c.input.o, c.input.a, c.input.b)
		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.result, actual) {
				t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
					c.input.o.Operator, c.input.a, c.input.b, c.expected.result, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestCmpTime(t *testing.T) {
	type Input struct {
		o    Opts
		a, b interface{}
	}

	type Expected struct {
		result bool
		err    error
	}

	type Case struct {
		input    Input
		expected Expected
	}

	// TODO: Test for errors.
	cases := []Case{
		{
			input: Input{
				o: Opts{Operator: tokenizer.Equals},
				a: time.Now().Round(time.Minute),
				b: time.Now().Round(time.Minute),
			},
			expected: Expected{result: true, err: nil},
		},
		{
			input: Input{
				o: Opts{Operator: tokenizer.Equals},
				a: time.Now().Add(time.Hour),
				b: time.Now().Add(-1 * time.Hour),
			},
			expected: Expected{result: false, err: nil},
		},

		{
			input: Input{
				o: Opts{Operator: tokenizer.NotEquals},
				a: time.Now().Round(time.Minute),
				b: time.Now().Round(time.Minute),
			},
			expected: Expected{result: false, err: nil},
		},
		{
			input: Input{
				o: Opts{Operator: tokenizer.NotEquals},
				a: time.Now().Add(time.Hour),
				b: time.Now().Add(-1 * time.Hour),
			},
			expected: Expected{result: true, err: nil},
		},

		{
			input: Input{
				o: Opts{Operator: tokenizer.In},
				a: time.Now().Round(time.Minute),
				b: map[interface{}]bool{time.Now().Round(time.Minute): true},
			},
			expected: Expected{result: true, err: nil},
		},
		{
			input: Input{
				o: Opts{Operator: tokenizer.In},
				a: time.Now().Round(time.Minute),
				b: map[interface{}]bool{time.Now().Add(time.Hour): true},
			},
			expected: Expected{result: false, err: nil},
		},
	}

	for _, c := range cases {
		actual, err := cmpTime(&c.input.o, c.input.a, c.input.b)
		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.result, actual) {
				t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
					c.input.o.Operator, c.input.a, c.input.b, c.expected.result, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}

}

func TestCmpMode(t *testing.T) {
	type Input struct {
		o    Opts
		file os.FileInfo
		typ  interface{}
	}

	type Expected struct {
		result bool
		err    error
	}

	type Case struct {
		input    Input
		expected Expected
	}

	// TODO: Complete these cases.
	cases := []Case{}

	for _, c := range cases {
		actual, err := cmpMode(&c.input.o)
		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.result, actual) {
				t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
					c.input.o.Operator, c.input.file, c.input.typ, c.expected.result, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestCmpHash(t *testing.T) {
	type Input struct {
		o    Opts
		file os.FileInfo
		typ  interface{}
	}

	type Expected struct {
		result bool
		err    error
	}

	type Case struct {
		input    Input
		expected Expected
	}

	// TODO: Complete these cases.
	cases := []Case{}

	for _, c := range cases {
		actual, err := cmpHash(&c.input.o)
		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.result, actual) {
				t.Fatalf("%v, %v, %v\nExpected: %v\n     Got: %v",
					c.input.o.Operator, c.input.file, c.input.typ, c.expected.result, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}
