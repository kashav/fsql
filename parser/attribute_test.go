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
		Case{
			input:    "name",
			expected: Expected{attributes: map[string]bool{"name": true}, err: nil},
		},
		Case{
			input: "name, size",
			expected: Expected{
				attributes: map[string]bool{"name": true, "size": true},
				err:        nil},
		},
		Case{
			input:    "*",
			expected: Expected{attributes: allAttributes, err: nil},
		},
		Case{
			input:    "all",
			expected: Expected{attributes: allAttributes, err: nil},
		},
		Case{
			input:    "format(time, iso)",
			expected: Expected{attributes: map[string]bool{"time": true}, err: nil},
		},

		Case{
			input:    "",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		Case{
			input:    "name,",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		Case{
			input:    "identifier",
			expected: Expected{err: &ErrUnknownToken{"identifier"}},
		},
	}

	for _, c := range cases {
		attributes := make(map[string]bool, 0)
		modifiers := make(map[string][]query.Modifier, 0)

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
		Case{
			input: "name",
			expected: Expected{
				modifiers: map[string][]query.Modifier{"name": []query.Modifier{}},
				err:       nil,
			},
		},
		Case{
			input: "upper(name)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": []query.Modifier{
						query.Modifier{
							Name:      "UPPER",
							Arguments: []string{},
						},
					},
				},
				err: nil,
			},
		},
		Case{
			input: "format(time, iso)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"time": []query.Modifier{
						query.Modifier{
							Name:      "FORMAT",
							Arguments: []string{"iso"},
						},
					},
				},
				err: nil,
			},
		},
		Case{
			input: "format(time, \"iso\")",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"time": []query.Modifier{
						query.Modifier{
							Name:      "FORMAT",
							Arguments: []string{"iso"},
						},
					},
				},
				err: nil,
			},
		},
		Case{
			input: "lower(name), format(size, mb)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": []query.Modifier{
						query.Modifier{
							Name:      "LOWER",
							Arguments: []string{},
						},
					},
					"size": []query.Modifier{
						query.Modifier{
							Name:      "FORMAT",
							Arguments: []string{"mb"},
						},
					},
				},
				err: nil,
			},
		},
		Case{
			input: "format(fullpath(name), lower)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": []query.Modifier{
						query.Modifier{
							Name:      "FULLPATH",
							Arguments: []string{},
						},
						query.Modifier{
							Name:      "FORMAT",
							Arguments: []string{"lower"},
						},
					},
				},
				err: nil,
			},
		},

		// No function/parameter validation yet!
		Case{
			input: "foo(name)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"name": []query.Modifier{query.Modifier{
						Name:      "FOO",
						Arguments: []string{},
					},
					},
				},
				err: nil,
			},
		},
		Case{
			input: "format(size, tb)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"size": []query.Modifier{
						query.Modifier{
							Name:      "FORMAT",
							Arguments: []string{"tb"},
						},
					},
				},
				err: nil,
			},
		},
		Case{
			input: "format(size, kb, mb)",
			expected: Expected{
				modifiers: map[string][]query.Modifier{
					"size": []query.Modifier{
						query.Modifier{
							Name:      "FORMAT",
							Arguments: []string{"kb", "mb"},
						},
					},
				},
				err: nil,
			},
		},

		Case{
			input:    "",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		Case{
			input:    "lower(name),",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},
		Case{
			input:    "identifier",
			expected: Expected{err: &ErrUnknownToken{"identifier"}},
		},
	}

	for _, c := range cases {
		attributes := make(map[string]bool, 0)
		modifiers := make(map[string][]query.Modifier, 0)

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
