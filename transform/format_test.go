package transform

import (
	"fmt"
	"reflect"
	"testing"
)

type FormatOutput struct {
	val interface{}
	err error
}

type FormatCase struct {
	params   *FormatParams
	expected FormatOutput
}

func TestTransform_Format(t *testing.T) {
	// TODO: Add case with time format to unix/iso (might need a fixture for
	// this).
	cases := []FormatCase{
		{
			&FormatParams{"size", "path", nil, int64(300), "format", []string{"kb"}},
			FormatOutput{fmt.Sprintf("%fkb", float64(300)/(1<<10)), nil},
		},
		{
			&FormatParams{"size", "path", nil, int64(300), "format", []string{"kilobytes"}},
			FormatOutput{nil, &ErrUnsupportedFormat{"kilobytes", "size"}},
		},
		{
			&FormatParams{"time", "path", nil, nil, "format", []string{""}},
			FormatOutput{nil, &ErrUnsupportedFormat{"", "time"}},
		},
		{
			&FormatParams{"name", "path", nil, "VALUE", "format", []string{"lower"}},
			FormatOutput{"value", nil},
		},
		{
			&FormatParams{"name", "path", nil, "value", "upper", []string{}},
			FormatOutput{"VALUE", nil},
		},
		{
			&FormatParams{"name", "path", nil, "value", "fullpath", []string{}},
			FormatOutput{"path", nil},
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
