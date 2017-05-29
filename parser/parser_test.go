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
		{"all", Expected{allAttributes, map[string][]query.Modifier{}, nil}},
		{"SELECT", Expected{allAttributes, map[string][]query.Modifier{}, nil}},
		{"FROM", Expected{allAttributes, map[string][]query.Modifier{}, nil}},

		{"SELECT name", Expected{
			map[string]bool{"name": true},
			map[string][]query.Modifier{"name": []query.Modifier{}},
			nil}},

		{"SELECT format(size, kb)", Expected{
			map[string]bool{"size": true},
			map[string][]query.Modifier{
				"size": []query.Modifier{query.Modifier{"FORMAT", []string{"kb"}}}},
			nil}},

		{"", Expected{err: io.ErrUnexpectedEOF}},
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
		{"WHERE", Expected{
			map[string][]string{"include": []string{"."}, "exclude": []string{}},
			map[string]string{},
			nil}},

		{"FROM .", Expected{
			map[string][]string{"include": []string{"."}, "exclude": []string{}},
			map[string]string{},
			nil}},

		{"FROM ~/foo, -./.git/", Expected{
			map[string][]string{
				"include": []string{u.HomeDir + "/foo"},
				"exclude": []string{".git"}},
			map[string]string{},
			nil}},

		{"FROM ./foo/ AS foo", Expected{
			map[string][]string{"include": []string{"foo"}, "exclude": []string{}},
			map[string]string{"foo": "foo"},
			nil}},

		{"FROM", Expected{err: io.ErrUnexpectedEOF}},

		{"FROM WHERE", Expected{
			err: &ErrUnexpectedToken{
				Actual:   tokenizer.Where,
				Expected: tokenizer.Identifier}}},
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
		{"WHERE name LIKE foo", Expected{
			&query.ConditionNode{
				Condition: &query.Condition{
					Attribute: "name",
					Operator:  tokenizer.Like,
					Value:     "foo"}}, nil}},

		// Our tree is fully-zeroed in this case, so it's easier just to give it
		// an empty Expected struct.
		{"", Expected{}},

		{"WHERE", Expected{err: io.ErrUnexpectedEOF}},

		{"name LIKE foo", Expected{
			err: &ErrUnexpectedToken{
				Expected: tokenizer.Where, Actual: tokenizer.Identifier}}},
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
		{tokenizer.Select, &tokenizer.Token{tokenizer.Select, "SELECT"}},
		{tokenizer.From, nil},
		{tokenizer.Identifier, &tokenizer.Token{tokenizer.Identifier, "all"}},
		{tokenizer.Identifier, nil},
		{tokenizer.From, &tokenizer.Token{tokenizer.From, "FROM"}},
		{tokenizer.Identifier, &tokenizer.Token{tokenizer.Identifier, "."}},
		{tokenizer.Identifier, nil},
		{tokenizer.Where, &tokenizer.Token{tokenizer.Where, "WHERE"}},
		{tokenizer.Identifier, &tokenizer.Token{tokenizer.Identifier, "name"}},
		{tokenizer.Equals, &tokenizer.Token{tokenizer.Equals, "="}},
		{tokenizer.Identifier, &tokenizer.Token{tokenizer.Identifier, "foo"}},
		{tokenizer.Or, &tokenizer.Token{tokenizer.Or, "OR"}},
		{tokenizer.Identifier, &tokenizer.Token{tokenizer.Identifier, "size"}},
		{tokenizer.Identifier, nil},
		{tokenizer.NotEquals, &tokenizer.Token{tokenizer.NotEquals, "<>"}},
		{tokenizer.Identifier, &tokenizer.Token{tokenizer.Identifier, "100"}},
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
		Sources:    map[string][]string{"include": []string{"."}, "exclude": []string{}},
		ConditionTree: &query.ConditionNode{
			Condition: &query.Condition{
				Attribute: "name",
				Operator:  tokenizer.Like,
				Value:     "foo"}},
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
					Sources:    map[string][]string{"include": []string{"."}, "exclude": []string{}},
					ConditionTree: &query.ConditionNode{
						Condition: &query.Condition{
							Attribute: "name",
							Operator:  tokenizer.Like,
							Value:     "foo"}},
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
