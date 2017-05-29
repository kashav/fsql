package parser

import (
	"io"
	"reflect"
	"testing"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
)

func TestAttributeParser_ExpectCorrectAttributes(t *testing.T) {

	type Expected struct {
		attributes map[string]bool
		err        error
	}

	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{"name", Expected{map[string]bool{"name": true}, nil}},
		{"name, size", Expected{map[string]bool{"name": true, "size": true}, nil}},
		{"*", Expected{allAttributes, nil}},
		{"all", Expected{allAttributes, nil}},
		{"format(time, iso)", Expected{map[string]bool{"time": true}, nil}},

		{"", Expected{err: io.ErrUnexpectedEOF}},
		{"name,", Expected{err: io.ErrUnexpectedEOF}},
		{"identifier", Expected{err: &ErrUnknownToken{"identifier"}}},
	}

	for _, c := range cases {
		attributes := make(map[string]bool, 0)
		modifiers := make(map[string][]query.Modifier, 0)

		err := (&parser{
			tokenizer: tokenizer.NewTokenizer(c.input)}).parseAttrList(&attributes, &modifiers)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.attributes, attributes) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.attributes, attributes)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestAttributeParser_ExpectCorrectModifiers(t *testing.T) {
	type Expected struct {
		modifiers map[string][]query.Modifier
		err       error
	}
	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{"name", Expected{map[string][]query.Modifier{"name": []query.Modifier{}}, nil}},
		{"upper(name)", Expected{map[string][]query.Modifier{
			"name": []query.Modifier{query.Modifier{"UPPER", []string{}}}}, nil}},
		{"format(time, iso)", Expected{map[string][]query.Modifier{
			"time": []query.Modifier{query.Modifier{
				"FORMAT", []string{"iso"}}}}, nil}},
		{"format(time, \"iso\")", Expected{map[string][]query.Modifier{
			"time": []query.Modifier{query.Modifier{
				"FORMAT", []string{"iso"}}}}, nil}},
		{"lower(name), format(size, mb)", Expected{map[string][]query.Modifier{
			"name": []query.Modifier{query.Modifier{"LOWER", []string{}}},
			"size": []query.Modifier{query.Modifier{"FORMAT", []string{"mb"}}}}, nil}},
		{"format(fullpath(name), lower)", Expected{map[string][]query.Modifier{
			"name": []query.Modifier{
				query.Modifier{"FULLPATH", []string{}},
				query.Modifier{"FORMAT", []string{"lower"}}}}, nil}},

		// No function/parameter validation yet!
		{"foo(name)", Expected{map[string][]query.Modifier{
			"name": []query.Modifier{query.Modifier{"FOO", []string{}}}}, nil}},
		{"format(size, tb)", Expected{map[string][]query.Modifier{
			"size": []query.Modifier{query.Modifier{"FORMAT", []string{"tb"}}}}, nil}},
		{"format(size, kb, mb)", Expected{map[string][]query.Modifier{
			"size": []query.Modifier{query.Modifier{"FORMAT", []string{"kb", "mb"}}}}, nil}},

		{"", Expected{err: io.ErrUnexpectedEOF}},
		{"lower(name),", Expected{err: io.ErrUnexpectedEOF}},
		{"identifier", Expected{err: &ErrUnknownToken{"identifier"}}},
	}

	for _, c := range cases {
		attributes := make(map[string]bool, 0)
		modifiers := make(map[string][]query.Modifier, 0)

		err := (&parser{
			tokenizer: tokenizer.NewTokenizer(c.input)}).parseAttrList(&attributes, &modifiers)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.modifiers, modifiers) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.modifiers, modifiers)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}
