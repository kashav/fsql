package transform

import (
	"reflect"
	"testing"
)

type ParseOutput struct {
	val interface{}
	err error
}

type ParseCase struct {
	params   *ParseParams
	expected ParseOutput
}

func TestTransform_Parse(t *testing.T) {
	// TODO: Complete this.
	cases := []ParseCase{}

	for _, c := range cases {
		val, err := Parse(c.params)
		if !(reflect.DeepEqual(val, c.expected.val) &&
			reflect.DeepEqual(err, c.expected.err)) {
			t.Fatalf("\nExpected: %v, %v\n     Got: %v, %v",
				c.expected.val, c.expected.err,
				val, err)
		}
	}
}
