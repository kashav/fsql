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
		{
			params: &FormatParams{
				Attribute: "size",
				Path:      "path",
				Info:      nil,
				Value:     int64(300),
				Name:      "format",
				Args:      []string{"kb"},
			},
			expected: Expected{
				val: fmt.Sprintf("%fkb", float64(300)/(1<<10)),
				err: nil,
			},
		},
		{
			params: &FormatParams{
				Attribute: "size",
				Path:      "path",
				Info:      nil,
				Value:     int64(300),
				Name:      "format",
				Args:      []string{"kilobytes"},
			},
			expected: Expected{
				val: nil,
				err: &ErrUnsupportedFormat{"kilobytes", "size"},
			},
		},
		{
			params: &FormatParams{
				Attribute: "name",
				Path:      "path",
				Info:      nil,
				Value:     "VALUE",
				Name:      "format",
				Args:      []string{"lower"},
			},
			expected: Expected{val: "value", err: nil},
		},
		{
			params: &FormatParams{
				Attribute: "name",
				Path:      "path",
				Info:      nil,
				Value:     "value",
				Name:      "upper",
				Args:      []string{},
			},
			expected: Expected{val: "VALUE", err: nil},
		},
		{
			params: &FormatParams{
				Attribute: "name",
				Path:      "path",
				Info:      nil,
				Value:     "value",
				Name:      "fullpath",
				Args:      []string{},
			},
			expected: Expected{val: "path", err: nil},
		},
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
