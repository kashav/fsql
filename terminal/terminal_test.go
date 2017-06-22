package terminal

import (
	"bytes"
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
		query    *bytes.Buffer
		expected Expected
	}

	// We already test core functionality in the fsql package, so we can stick
	// with _simple_ queries for the following cases.
	cases := []Case{
		{
			query: bytes.NewBufferString("select name, hash from ../testdata where name = baz"),
			expected: Expected{
				out: "da39a3e\tbaz\n",
				err: nil,
			},
		},
		{
			query: bytes.NewBufferString("select all from"),
			expected: Expected{
				out: "",
				err: errors.New("unexpected EOF"),
			},
		},
	}

	for _, c := range cases {
		actual, err := run(*c.query)
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
