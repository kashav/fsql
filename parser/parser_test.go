package parser

import (
	"io"
	"os/user"
	"reflect"
	"testing"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
)

func TestParser_ParseSelect(t *testing.T) {
	type Expected struct {
		attributes map[string]bool
		modifiers  map[string][]query.Modifier
		err        error
	}

	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{
			input: "all",
			expected: Expected{
				attributes: allAttributes,
				modifiers:  map[string][]query.Modifier{},
				err:        nil,
			},
		},

		{
			input: "SELECT",
			expected: Expected{
				attributes: allAttributes,
				modifiers:  map[string][]query.Modifier{},
				err:        nil,
			},
		},

		{
			input: "FROM",
			expected: Expected{
				attributes: allAttributes,
				modifiers:  map[string][]query.Modifier{},
				err:        nil,
			},
		},

		{
			input: "SELECT name",
			expected: Expected{
				attributes: map[string]bool{"name": true},
				modifiers:  map[string][]query.Modifier{"name": {}},
				err:        nil,
			},
		},

		{
			input: "SELECT format(size, kb)",
			expected: Expected{
				attributes: map[string]bool{"size": true},
				modifiers: map[string][]query.Modifier{
					"size": {
						{
							Name:      "FORMAT",
							Arguments: []string{"kb"},
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
	}

	for _, c := range cases {
		q := query.NewQuery()
		err := (&parser{tokenizer: tokenizer.NewTokenizer(c.input)}).parseSelectClause(q)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.attributes, q.Attributes) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.attributes, q.Attributes)
			}
			if !reflect.DeepEqual(c.expected.modifiers, q.Modifiers) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.modifiers, q.Modifiers)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestParser_ParseFrom(t *testing.T) {
	type Expected struct {
		sources map[string][]string
		aliases map[string]string
		err     error
	}

	type Case struct {
		input    string
		expected Expected
	}

	u, err := user.Current()
	if err != nil {
		// TODO: If we can't get the current user, should we fatal or just return?
		return
	}

	cases := []Case{
		{
			input: "WHERE",
			expected: Expected{
				sources: map[string][]string{
					"include": {"."},
					"exclude": {},
				},
				aliases: map[string]string{},
				err:     nil,
			},
		},

		{
			input: "FROM .",
			expected: Expected{
				sources: map[string][]string{
					"include": {"."},
					"exclude": {},
				},
				aliases: map[string]string{},
				err:     nil,
			},
		},

		{
			input: "FROM ~/foo, -./.git/",
			expected: Expected{
				sources: map[string][]string{
					"include": {u.HomeDir + "/foo"},
					"exclude": {".git"},
				},
				aliases: map[string]string{},
				err:     nil,
			},
		},

		{
			input: "FROM ./foo/ AS foo",
			expected: Expected{
				sources: map[string][]string{
					"include": {"foo"},
					"exclude": {},
				},
				aliases: map[string]string{"foo": "foo"},
				err:     nil,
			},
		},

		{
			input:    "FROM",
			expected: Expected{err: io.ErrUnexpectedEOF},
		},

		{
			input: "FROM WHERE",
			expected: Expected{
				err: &ErrUnexpectedToken{
					Actual:   tokenizer.Where,
					Expected: tokenizer.Identifier,
				},
			},
		},
	}

	for _, c := range cases {
		q := query.NewQuery()
		err := (&parser{tokenizer: tokenizer.NewTokenizer(c.input)}).parseFromClause(q)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.sources, q.Sources) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.sources, q.Sources)
			}
			if !reflect.DeepEqual(c.expected.aliases, q.SourceAliases) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.aliases, q.SourceAliases)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestParser_ParseWhere(t *testing.T) {
	type Expected struct {
		tree *query.ConditionNode
		err  error
	}

	type Case struct {
		input    string
		expected Expected
	}

	cases := []Case{
		{
			input: "WHERE name LIKE foo",
			expected: Expected{
				tree: &query.ConditionNode{
					Condition: &query.Condition{
						Attribute: "name",
						Operator:  tokenizer.Like,
						Value:     "foo",
					},
				},
				err: nil,
			},
		},

		// Our tree is fully-zeroed in this case, so it's easier just to give it
		// an empty Expected struct.
		{input: "", expected: Expected{}},

		{input: "WHERE", expected: Expected{err: io.ErrUnexpectedEOF}},

		{
			input: "name LIKE foo",
			expected: Expected{
				err: &ErrUnexpectedToken{
					Expected: tokenizer.Where,
					Actual:   tokenizer.Identifier,
				},
			},
		},
	}

	for _, c := range cases {
		q := query.NewQuery()
		err := (&parser{tokenizer: tokenizer.NewTokenizer(c.input)}).parseWhereClause(q)

		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.tree, q.ConditionTree) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.tree, q.ConditionTree)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}

func TestParser_Expect(t *testing.T) {
	type Case struct {
		param    tokenizer.TokenType
		expected *tokenizer.Token
	}

	input := "SELECT all FROM . WHERE name = foo OR size <> 100"
	p := &parser{tokenizer: tokenizer.NewTokenizer(input)}

	cases := []Case{
		{
			param:    tokenizer.Select,
			expected: &tokenizer.Token{Type: tokenizer.Select, Raw: "SELECT"},
		},
		{
			param:    tokenizer.From,
			expected: nil,
		},
		{
			param:    tokenizer.Identifier,
			expected: &tokenizer.Token{Type: tokenizer.Identifier, Raw: "all"},
		},
		{
			param:    tokenizer.Identifier,
			expected: nil,
		},
		{
			param:    tokenizer.From,
			expected: &tokenizer.Token{Type: tokenizer.From, Raw: "FROM"},
		},
		{
			param:    tokenizer.Identifier,
			expected: &tokenizer.Token{Type: tokenizer.Identifier, Raw: "."},
		},
		{
			param:    tokenizer.Identifier,
			expected: nil,
		},
		{
			param:    tokenizer.Where,
			expected: &tokenizer.Token{Type: tokenizer.Where, Raw: "WHERE"},
		},
		{
			param:    tokenizer.Identifier,
			expected: &tokenizer.Token{Type: tokenizer.Identifier, Raw: "name"},
		},
		{
			param:    tokenizer.Equals,
			expected: &tokenizer.Token{Type: tokenizer.Equals, Raw: "="},
		},
		{
			param:    tokenizer.Identifier,
			expected: &tokenizer.Token{Type: tokenizer.Identifier, Raw: "foo"},
		},
		{
			param:    tokenizer.Or,
			expected: &tokenizer.Token{Type: tokenizer.Or, Raw: "OR"},
		},
		{
			param:    tokenizer.Identifier,
			expected: &tokenizer.Token{Type: tokenizer.Identifier, Raw: "size"},
		},
		{
			param:    tokenizer.Identifier,
			expected: nil,
		},
		{
			param:    tokenizer.NotEquals,
			expected: &tokenizer.Token{Type: tokenizer.NotEquals, Raw: "<>"},
		},
		{
			param:    tokenizer.Identifier,
			expected: &tokenizer.Token{Type: tokenizer.Identifier, Raw: "100"},
		},
	}

	for _, c := range cases {
		actual := p.expect(c.param)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}

func TestParser_SelectAllVariations(t *testing.T) {
	expected := &query.Query{
		Attributes: allAttributes,
		Sources: map[string][]string{
			"include": {"."},
			"exclude": {},
		},
		ConditionTree: &query.ConditionNode{
			Condition: &query.Condition{
				Attribute: "name",
				Operator:  tokenizer.Like,
				Value:     "foo",
			},
		},
		SourceAliases: map[string]string{},
		Modifiers:     map[string][]query.Modifier{},
	}

	cases := []string{
		"FROM . WHERE name LIKE foo",
		"all FROM . WHERE name LIKE foo",
		"SELECT FROM . WHERE name LIKE foo",
		"SELECT all FROM . WHERE name LIKE foo",
	}

	for _, c := range cases {
		actual, err := Run(c)
		if err != nil {
			t.Fatalf("\nExpected no error\n     Got %v", err)
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("\nExpected %v\n     Got %v", expected, actual)
		}
	}
}

func TestParser_Run(t *testing.T) {
	type Expected struct {
		q   *query.Query
		err error
	}

	type Case struct {
		input    string
		expected Expected
	}

	// TODO: Add more cases.
	cases := []Case{
		{
			input: "SELECT all FROM . WHERE name LIKE foo",
			expected: Expected{
				q: &query.Query{
					Attributes: allAttributes,
					Sources: map[string][]string{
						"include": {"."},
						"exclude": {},
					},
					ConditionTree: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "name",
							Operator:  tokenizer.Like,
							Value:     "foo",
						},
					},
					SourceAliases: map[string]string{},
					Modifiers:     map[string][]query.Modifier{},
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		actual, err := Run(c.input)
		if c.expected.err == nil {
			if err != nil {
				t.Fatalf("\nExpected no error\n     Got %v", err)
			}
			if !reflect.DeepEqual(c.expected.q, actual) {
				t.Fatalf("\nExpected %v\n     Got %v", c.expected.q, actual)
			}
		} else if !reflect.DeepEqual(c.expected.err, err) {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected.err, err)
		}
	}
}
