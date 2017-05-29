package parser

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kshvmdn/fsql/tokenizer"
)

func TestSourceParser_ExpectCorrectSources(t *testing.T) {
	type Expected struct {
		sources map[string][]string
		err     error
	}

	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{".", Expected{map[string][]string{"include": []string{"."}}, nil}},
		{"., ~/foo",
			Expected{map[string][]string{"include": []string{".", "~/foo"}}, nil}},
		{"., -.bar",
			Expected{map[string][]string{
				"include": []string{"."}, "exclude": []string{".bar"}}, nil}},
		{"-.bar, ., ~/foo AS foo",
			Expected{map[string][]string{
				"include": []string{".", "~/foo"}, "exclude": []string{".bar"}}, nil}},

		{"", Expected{err: io.ErrUnexpectedEOF}},
		{"foo,", Expected{err: io.ErrUnexpectedEOF}},
	}

	for _, c := range cases {
		sources := make(map[string][]string, 0)
		aliases := make(map[string]string, 0)

		err := (&parser{
			tokenizer: tokenizer.NewTokenizer(c.input)}).parseSourceList(&sources, &aliases)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.sources, sources) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.sources, sources)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestSourceParser_ExpectCorrectAliases(t *testing.T) {
	type Expected struct {
		aliases map[string]string
		err     error
	}

	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{".", Expected{map[string]string{}, nil}},
		{". AS cwd", Expected{map[string]string{"cwd": "."}, nil}},
		{"., -.bar, ~/foo AS foo", Expected{map[string]string{"foo": "~/foo"}, nil}},

		{"-.bar AS bar", Expected{
			err: errors.New("cannot alias excluded directory .bar")}},
		{"", Expected{err: io.ErrUnexpectedEOF}},
		{"foo AS", Expected{err: io.ErrUnexpectedEOF}},
	}

	for _, c := range cases {
		sources := make(map[string][]string, 0)
		aliases := make(map[string]string, 0)

		err := (&parser{
			tokenizer: tokenizer.NewTokenizer(c.input)}).parseSourceList(&sources, &aliases)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.aliases, aliases) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.aliases, aliases)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}
