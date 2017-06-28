package terminal

import (
	"errors"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	type Expected struct {
		out string
		err error
	}

	type Case struct {
		query    string
		expected Expected
	}

	// We already test core functionality in the fsql package, so we can stick
	// with _simple_ queries for the following cases.
	cases := []Case{
		{
			query: "select name, hash from ../testdata where name = baz",
			expected: Expected{
				out: "baz\tda39a3e\n",
				err: nil,
			},
		},
		{
			query: "select all from",
			expected: Expected{
				out: "",
				err: errors.New("unexpected EOF"),
			},
		},
	}

	for _, c := range cases {
		actual, err := run(c.query)
		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.out, actual) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.out, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}
