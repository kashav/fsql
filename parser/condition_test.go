package parser

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
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
		{"name LIKE foo%", Expected{&query.Condition{
			Attribute: "name",
			Operator:  tokenizer.Like,
			Value:     "foo%"}, nil}},
		{"size = 10", Expected{&query.Condition{
			Attribute: "size",
			Operator:  tokenizer.Equals,
			Value:     "10"}, nil}},
		{"mode IS dir", Expected{&query.Condition{
			Attribute: "mode",
			Operator:  tokenizer.Is,
			Value:     "dir"}, nil}},
		{"format(time, iso) >= 2017-05-28T16:37:18Z", Expected{&query.Condition{
			Attribute:          "time",
			AttributeModifiers: []query.Modifier{query.Modifier{"FORMAT", []string{"iso"}}},
			Operator:           tokenizer.GreaterThanEquals,
			Value:              "2017-05-28T16:37:18Z"}, nil}},
		{"upper(name) != FOO", Expected{&query.Condition{
			Attribute:          "name",
			AttributeModifiers: []query.Modifier{query.Modifier{"UPPER", []string{}}},
			Operator:           tokenizer.NotEquals,
			Value:              "FOO"}, nil}},
		{"NOT name IN [foo,bar,baz]", Expected{&query.Condition{
			Attribute: "name",
			Operator:  tokenizer.In,
			Value:     []string{"foo", "bar", "baz"},
			Negate:    true}, nil}},

		// No attribute-operator validation yet (these should /eventually/ throw
		// some error)!
		{"time RLIKE '.*'", Expected{&query.Condition{
			Attribute: "time",
			Operator:  tokenizer.RLike,
			Value:     ".*"}, nil}},
		{"size LIKE foo", Expected{&query.Condition{
			Attribute: "size",
			Operator:  tokenizer.Like,
			Value:     "foo"}, nil}},
		{"time <> now", Expected{&query.Condition{
			Attribute: "time",
			Operator:  tokenizer.NotEquals,
			Value:     "now"}, nil}},

		{"name =", Expected{err: io.ErrUnexpectedEOF}},
		{"file IS dir", Expected{err: &ErrUnknownToken{"file"}}},
	}

	for _, c := range cases {
		actual, err := (&parser{
			tokenizer: tokenizer.NewTokenizer(c.input)}).parseNextCond()

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
		{"name LIKE foo%", Expected{&query.ConditionNode{
			Condition: &query.Condition{
				Attribute: "name",
				Operator:  tokenizer.Like,
				Value:     "foo%"}}, nil}},

		{"upper(name) = MAIN", Expected{&query.ConditionNode{
			Condition: &query.Condition{
				Attribute:          "name",
				AttributeModifiers: []query.Modifier{query.Modifier{"UPPER", []string{}}},
				Operator:           tokenizer.Equals,
				Value:              "MAIN"}}, nil}},

		{"name LIKE %foo AND name <> bar.foo", Expected{&query.ConditionNode{
			Type: &tmpAnd,
			Left: &query.ConditionNode{
				Condition: &query.Condition{
					Attribute: "name",
					Operator:  tokenizer.Like,
					Value:     "%foo"}},
			Right: &query.ConditionNode{
				Condition: &query.Condition{
					Attribute: "name",
					Operator:  tokenizer.NotEquals,
					Value:     "bar.foo"}}}, nil}},

		{"size <= 10 OR NOT mode IS dir", Expected{&query.ConditionNode{
			Type: &tmpOr,
			Left: &query.ConditionNode{
				Condition: &query.Condition{
					Attribute: "size",
					Operator:  tokenizer.LessThanEquals,
					Value:     "10"}},
			Right: &query.ConditionNode{
				Condition: &query.Condition{
					Attribute: "mode",
					Operator:  tokenizer.Is,
					Value:     "dir",
					Negate:    true}}}, nil}},

		{"size = 5 AND name = foo", Expected{&query.ConditionNode{
			Type: &tmpAnd,
			Left: &query.ConditionNode{
				Condition: &query.Condition{
					Attribute: "size",
					Operator:  tokenizer.Equals,
					Value:     "5"}},
			Right: &query.ConditionNode{
				Condition: &query.Condition{
					Attribute: "name",
					Operator:  tokenizer.Equals,
					Value:     "foo"}}}, nil}},

		{"format(size, mb) <= 2 AND (name = foo OR name = bar)", Expected{&query.ConditionNode{
			Type: &tmpAnd,
			Left: &query.ConditionNode{
				Condition: &query.Condition{
					Attribute:          "size",
					AttributeModifiers: []query.Modifier{query.Modifier{"FORMAT", []string{"mb"}}},
					Operator:           tokenizer.LessThanEquals,
					Value:              "2"}},
			Right: &query.ConditionNode{
				Type: &tmpOr,
				Left: &query.ConditionNode{
					Condition: &query.Condition{
						Attribute: "name",
						Operator:  tokenizer.Equals,
						Value:     "foo"}},
				Right: &query.ConditionNode{
					Condition: &query.Condition{
						Attribute: "name",
						Operator:  tokenizer.Equals,
						Value:     "bar"}}}}, nil}},

		{"name = foo AND NOT (name = bar OR name = baz)", Expected{
			&query.ConditionNode{}, &ErrUnexpectedToken{
				Expected: tokenizer.Identifier,
				Actual:   tokenizer.OpenParen}}},

		{"size = 5 AND ()", Expected{err: errors.New("failed to parse conditions")}},

		// FIXME: The following case /should/ throw EOF (it isn't right now).
		// {"name = foo AND", Expected{err: io.ErrUnexpectedEOF}},
	}

	for _, c := range cases {
		actual, err := (&parser{
			tokenizer: tokenizer.NewTokenizer(c.input)}).parseConditionTree()

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
