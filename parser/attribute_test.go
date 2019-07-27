package parser

import (
	"io"
	"reflect"
	"testing"

	"github.com/kashav/fsql/query"
	"github.com/kashav/fsql/tokenizer"
)

func TestAttributeParser_ExpectCorrectAttributes(t *testing.T) {
	type Expected struct {
		attributes []string
		err        error
	}

	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{
			input:    "name",
			expected: Expected{attributes: []string{"name"}, err: nil},
		},
		{
			input: "name, size",
			expected: Expected{
				attributes: []string{"name", "size"},
				err:        nil,
			},
		},
		{
			input:    "*",
			expected: Expected{attributes: allAttributes, err: nil},
		},
		{
			input:    "all",
			expected: Expected{attributes: allAttributes, err: nil},
		},
		{
			input:    "format(time, iso)",
			expected: Expected{attributes: []string{"time"}, err: nil},
		},

		{
			input:    "",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		{
			input:    "name,",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		{
			input:    "identifier",
			expected: Expected{err: &ErrUnknownToken{"identifier"}},
		},
	}

	for _, c := range cases {
		attributes := make([]string, 0)
		modifiers := make(map[string][]query.Modifier)

		p := &parser{tokenizer: tokenizer.NewTokenizer(c.input)}
		err := p.parseAttrs(&attributes, &modifiers)

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
		{
			input: "name",
			expected: Expected{
				modifiers: map[string][]query.Modifier{"name": {}},
				err:       nil,
			},
		},
		{
			input: "upper(name)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": {
						{
							Name:      "UPPER",
							Arguments: []string{},
						},
					},
				},
				err: nil,
			},
		},
		{
			input: "format(time, iso)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"time": {
						{
							Name:      "FORMAT",
							Arguments: []string{"iso"},
						},
					},
				},
				err: nil,
			},
		},
		{
			input: "format(time, \"iso\")",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"time": {
						{
							Name:      "FORMAT",
							Arguments: []string{"iso"},
						},
					},
				},
				err: nil,
			},
		},
		{
			input: "lower(name), format(size, mb)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": {
						{
							Name:      "LOWER",
							Arguments: []string{},
						},
					},
					"size": {
						{
							Name:      "FORMAT",
							Arguments: []string{"mb"},
						},
					},
				},
				err: nil,
			},
		},
		{
			input: "format(fullpath(name), lower)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": {
						{
							Name:      "FULLPATH",
							Arguments: []string{},
						},
						{
							Name:      "FORMAT",
							Arguments: []string{"lower"},
						},
					},
				},
				err: nil,
			},
		},

		// No function/parameter validation yet!
		{
			input: "foo(name)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": {{
						Name:      "FOO",
						Arguments: []string{},
					},
					},
				},
				err: nil,
			},
		},
		{
			input: "format(size, tb)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"size": {
						{
							Name:      "FORMAT",
							Arguments: []string{"tb"},
						},
					},
				},
				err: nil,
			},
		},
		{
			input: "format(size, kb, mb)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"size": {
						{
							Name:      "FORMAT",
							Arguments: []string{"kb", "mb"},
						},
					},
				},
				err: nil,
			},
		},

		{
			input:    "",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		{
			input:    "lower(name),",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		{
			input:    "identifier",
			expected: Expected{err: &ErrUnknownToken{"identifier"}},
		},
	}

	for _, c := range cases {
		attributes := make([]string, 0)
		modifiers := make(map[string][]query.Modifier)

		p := &parser{tokenizer: tokenizer.NewTokenizer(c.input)}
		err := p.parseAttrs(&attributes, &modifiers)

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
