package parser

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kashav/fsql/query"
	"github.com/kashav/fsql/tokenizer"
)

func TestConditionParser_ExpectCorrectCondition(t *testing.T) {
	type Expected struct {
		condition *query.Condition
		err       error
	}

	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{
			input: "name LIKE foo%",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "name",
					Operator:  tokenizer.Like,
					Value:     "foo%",
				},
				err: nil,
			},
		},

		{
			input: "size = 10",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "size",
					Operator:  tokenizer.Equals,
					Value:     "10",
				},
				err: nil,
			},
		},

		{
			input: "mode IS dir",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "mode",
					Operator:  tokenizer.Is,
					Value:     "dir",
				},
				err: nil,
			},
		},

		{
			input: "format(time, iso) >= 2017-05-28T16:37:18Z",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "time",
					AttributeModifiers: []query.Modifier{
						{
							Name:      "FORMAT",
							Arguments: []string{"iso"},
						},
					},
					Operator: tokenizer.GreaterThanEquals,
					Value:    "2017-05-28T16:37:18Z",
				},
				err: nil,
			},
		},

		{
			input: "upper(name) != FOO",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "name",
					AttributeModifiers: []query.Modifier{
						{
							Name:      "UPPER",
							Arguments: []string{},
						},
					},
					Operator: tokenizer.NotEquals,
					Value:    "FOO",
				},
				err: nil,
			},
		},

		{
			input: "NOT name IN [foo,bar,baz]",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "name",
					Operator:  tokenizer.In,
					Value:     []string{"foo", "bar", "baz"},
					Negate:    true,
				},
				err: nil,
			},
		},

		// No attribute-operator validation yet (these 3 should /eventually/ throw
		// some error)!
		{
			input: "time RLIKE '.*'",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "time",
					Operator:  tokenizer.RLike,
					Value:     ".*",
				},
				err: nil,
			},
		},

		{
			input: "size LIKE foo",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "size",
					Operator:  tokenizer.Like,
					Value:     "foo",
				},
				err: nil,
			},
		},

		{
			input: "time <> now",
			expected: Expected{
				condition: &query.Condition{
					Attribute: "time",
					Operator:  tokenizer.NotEquals,
					Value:     "now",
				},
				err: nil,
			},
		},

		{
			input:    "name =",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},

		{
			input:    "file IS dir",
			expected: Expected{err: &ErrUnknownToken{"file"}},
		},
	}

	for _, c := range cases {
		p := &parser{tokenizer: tokenizer.NewTokenizer(c.input)}
		actual, err := p.parseCondition()

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.condition, actual) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.condition, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestConditionParser_ExpectCorrectConditionTree(t *testing.T) {
	type Expected struct {
		node *query.ConditionNode
		err  error
	}

	type Case struct {
		input    string
		expected Expected
	}

	// Not sure why, but compiler throws when attempting to take the memory
	// address of any tokenizer.TokenType.
	var (
		tmpAnd = tokenizer.And
		tmpOr  = tokenizer.Or
	)

	cases := []Case{
		{
			input: "name LIKE foo%",
			expected: Expected{
				node: &query.ConditionNode{
					Condition: &query.Condition{
						Attribute: "name",
						Operator:  tokenizer.Like,
						Value:     "foo%",
					},
				},
				err: nil,
			},
		},

		{
			input: "upper(name) = MAIN",
			expected: Expected{
				node: &query.ConditionNode{
					Condition: &query.Condition{
						Attribute: "name",
						AttributeModifiers: []query.Modifier{
							{
								Name:      "UPPER",
								Arguments: []string{},
							},
						},
						Operator: tokenizer.Equals,
						Value:    "MAIN",
					},
				},
				err: nil,
			},
		},

		{
			input: "name LIKE %foo AND name <> bar.foo",
			expected: Expected{
				node: &query.ConditionNode{
					Type: &tmpAnd,
					Left: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "name",
							Operator:  tokenizer.Like,
							Value:     "%foo",
						},
					},
					Right: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "name",
							Operator:  tokenizer.NotEquals,
							Value:     "bar.foo",
						},
					},
				},
				err: nil,
			},
		},

		{
			input: "size <= 10 OR NOT mode IS dir",
			expected: Expected{
				node: &query.ConditionNode{
					Type: &tmpOr,
					Left: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "size",
							Operator:  tokenizer.LessThanEquals,
							Value:     "10",
						},
					},
					Right: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "mode",
							Operator:  tokenizer.Is,
							Value:     "dir",
							Negate:    true,
						},
					},
				},
				err: nil,
			},
		},

		{
			input: "size = 5 AND name = foo",
			expected: Expected{
				node: &query.ConditionNode{
					Type: &tmpAnd,
					Left: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "size",
							Operator:  tokenizer.Equals,
							Value:     "5",
						},
					},
					Right: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "name",
							Operator:  tokenizer.Equals,
							Value:     "foo",
						},
					},
				},
				err: nil,
			},
		},

		{
			input: "format(size, mb) <= 2 AND (name = foo OR name = bar)",
			expected: Expected{
				node: &query.ConditionNode{
					Type: &tmpAnd,
					Left: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "size",
							AttributeModifiers: []query.Modifier{
								{
									Name:      "FORMAT",
									Arguments: []string{"mb"},
								},
							},
							Operator: tokenizer.LessThanEquals,
							Value:    "2",
						},
					},
					Right: &query.ConditionNode{
						Type: &tmpOr,
						Left: &query.ConditionNode{
							Condition: &query.Condition{
								Attribute: "name",
								Operator:  tokenizer.Equals,
								Value:     "foo",
							},
						},
						Right: &query.ConditionNode{
							Condition: &query.Condition{
								Attribute: "name",
								Operator:  tokenizer.Equals,
								Value:     "bar",
							},
						},
					},
				},
				err: nil,
			},
		},

		{
			input: "name = foo AND NOT (name = bar OR name = baz)",
			expected: Expected{
				node: &query.ConditionNode{},
				err: &ErrUnexpectedToken{
					Expected: tokenizer.Identifier,
					Actual:   tokenizer.OpenParen,
				},
			},
		},

		{
			input:    "size = 5 AND ()",
			expected: Expected{err: errors.New("failed to parse conditions")},
		},

		// FIXME: The following case /should/ throw EOF (it doesn't right now).
		// Case{input: "name = foo AND", expected: Expected{err: io.ErrUnexpectedEOF}},
	}

	for _, c := range cases {
		p := &parser{tokenizer: tokenizer.NewTokenizer(c.input)}
		actual, err := p.parseConditionTree()

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.node, actual) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.node, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestConditionParser_ExpectCorrectSubquery(t *testing.T) {
	type Expected struct {
		condition *query.Condition
		err       error
	}

	type Case struct {
		input    *query.Condition
		expected Expected
	}

	// TODO: Complete these cases. This test relies on testdata fixtures.
	cases := []Case{}

	for _, c := range cases {
		p := &parser{tokenizer: tokenizer.NewTokenizer("")}
		err := p.parseSubquery(c.input)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.condition, c.input) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.condition, c.input)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}
