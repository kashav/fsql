package transform

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTransform_Format(t *testing.T) {
	type Expected struct {
		val interface{}
		err error
	}

	type Case struct {
		params   *FormatParams
		expected Expected
	}

	// TODO: Add tests for the time attribute!
	cases := []Case{
		{&FormatParams{"size", "path", nil, int64(300), "format", []string{"kb"}},
			Expected{fmt.Sprintf("%fkb", float64(300)/(1<<10)), nil}},
		{&FormatParams{"size", "path", nil, int64(300), "format", []string{"kilobytes"}},
			Expected{nil, &ErrUnsupportedFormat{"kilobytes", "size"}}},
		{&FormatParams{"name", "path", nil, "VALUE", "format", []string{"lower"}},
			Expected{"value", nil}},
		{&FormatParams{"name", "path", nil, "value", "upper", []string{}},
			Expected{"VALUE", nil}},
		{&FormatParams{"name", "path", nil, "value", "fullpath", []string{}},
			Expected{"path", nil}},
	}

	for _, c := range cases {
		val, err := Format(c.params)
		if !(reflect.DeepEqual(val, c.expected.val) &&
			reflect.DeepEqual(err, c.expected.err)) {
			t.Fatalf("\nExpected: %v, %v\n     Got: %v, %v",
				c.expected.val, c.expected.err,
				val, err)
		}
	}
}
