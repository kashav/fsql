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
		Case{
			input: ".",
			expected: Expected{
				sources: map[string][]string{"include": []string{"."}},
				err:     nil,
			},
		},
		Case{
			input: "., ~/foo",
			expected: Expected{
				sources: map[string][]string{"include": []string{".", "~/foo"}},
				err:     nil,
			},
		},
		Case{
			input: "., -.bar",
			expected: Expected{
				sources: map[string][]string{
					"include": []string{"."},
					"exclude": []string{".bar"},
				},
				err: nil,
			},
		},
		Case{
			input: "-.bar, ., ~/foo AS foo",
			expected: Expected{
				sources: map[string][]string{
					"include": []string{".", "~/foo"},
					"exclude": []string{".bar"},
				},
				err: nil,
			},
		},

		Case{input: "", expected: Expected{err: io.ErrUnexpectedEOF}},
		Case{input: "foo,", expected: Expected{err: io.ErrUnexpectedEOF}},
	}

	for _, c := range cases {
		sources := make(map[string][]string, 0)
		aliases := make(map[string]string, 0)

		p := &parser{tokenizer: tokenizer.NewTokenizer(c.input)}
		err := p.parseSourceList(&sources, &aliases)

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
		Case{
			input: ".",
			expected: Expected{
				aliases: map[string]string{},
				err:     nil,
			},
		},
		Case{
			input: ". AS cwd",
			expected: Expected{
				aliases: map[string]string{"cwd": "."},
				err:     nil,
			},
		},
		Case{
			input: "., -.bar, ~/foo AS foo",
			expected: Expected{
				aliases: map[string]string{"foo": "~/foo"},
				err:     nil,
			},
		},

		Case{
			input:    "-.bar AS bar",
			expected: Expected{err: errors.New("cannot alias excluded directory .bar")},
		},
		Case{
			input:    "",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		Case{
			input:    "foo AS",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
	}

	for _, c := range cases {
		sources := make(map[string][]string, 0)
		aliases := make(map[string]string, 0)

		p := &parser{tokenizer: tokenizer.NewTokenizer(c.input)}
		err := p.parseSourceList(&sources, &aliases)

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
