package tokenizer

import (
	"reflect"
	"testing"
)

func TestTokenizer_NextTokenType(t *testing.T) {
	type Case struct {
		input    string
		expected TokenType
	}

	cases := []Case{
		{input: "SELECT", expected: Select},
		{input: "FROM", expected: From},
		{input: "WHERE", expected: Where},
		{input: "AS", expected: As},
		{input: "OR", expected: Or},
		{input: "AND", expected: And},
		{input: "NOT", expected: Not},
		{input: "IN", expected: In},
		{input: "IS", expected: Is},
		{input: "LIKE", expected: Like},
		{input: "RLIKE", expected: RLike},
		{input: "foo", expected: Identifier},
		{input: "(", expected: OpenParen},
		{input: ")", expected: CloseParen},
		{input: ",", expected: Comma},
		{input: "-", expected: Hyphen},
		{input: "=", expected: Equals},
		{input: "<>", expected: NotEquals},
		{input: "<", expected: LessThan},
		{input: "<=", expected: LessThanEquals},
		{input: ">", expected: GreaterThan},
		{input: ">=", expected: GreaterThanEquals},
	}

	for _, c := range cases {
		actual := NewTokenizer(c.input).Next()
		expected := &Token{Type: c.expected, Raw: c.input}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
		}
	}
}

func TestTokenizer_NextRaw(t *testing.T) {
	type Case struct {
		input    string
		expected string
	}

	// TODO: Fix the last 2 cases, they're currently hanging.
	cases := []Case{
		{input: "foo", expected: "foo"},
		{input: " foo ", expected: "foo"},
		{input: "\" foo \"", expected: " foo "},
		{input: "' foo '", expected: " foo "},
		{input: "` foo `", expected: " foo "},
		// Case{input: "\"foo'bar\"", expected: "foo'bar"},
		// Case{input: "\"()\"", expected: "()"},
	}

	for _, c := range cases {
		actual := NewTokenizer(c.input).Next()
		expected := &Token{Type: Identifier, Raw: c.expected}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nExpected: %v\n     Got: %v", expected, actual)
		}
	}
}

func TestTokenizer_AllSimple(t *testing.T) {
	input := `
    SELECT
      name, size
    FROM
      ~/Desktop
    WHERE
      name LIKE %go
    `

	actual := NewTokenizer(input).All()
	expected := []Token{
		{Type: Select, Raw: "SELECT"},
		{Type: Identifier, Raw: "name"},
		{Type: Comma, Raw: ","},
		{Type: Identifier, Raw: "size"},
		{Type: From, Raw: "FROM"},
		{Type: Identifier, Raw: "~/Desktop"},
		{Type: Where, Raw: "WHERE"},
		{Type: Identifier, Raw: "name"},
		{Type: Like, Raw: "LIKE"},
		{Type: Identifier, Raw: "%go"},
	}

	for i := range expected {
		if !reflect.DeepEqual(actual[i], expected[i]) {
			t.Fatalf("\nExpected: %v\n     Got: %v", expected[i], actual[i])
		}
	}
}

func TestTokenizer_AllSubquery(t *testing.T) {
	input := `
  SELECT
    name, size
  FROM
    ~/Desktop
  WHERE
    name LIKE %go OR
    name IN (
      SELECT
        name
      FROM
        $GOPATH/src/github.com
      WHERE
        name RLIKE .*_test\.go)
  `

	actual := NewTokenizer(input).All()
	expected := []Token{
		{Type: Select, Raw: "SELECT"},
		{Type: Identifier, Raw: "name"},
		{Type: Comma, Raw: ","},
		{Type: Identifier, Raw: "size"},
		{Type: From, Raw: "FROM"},
		{Type: Identifier, Raw: "~/Desktop"},
		{Type: Where, Raw: "WHERE"},
		{Type: Identifier, Raw: "name"},
		{Type: Like, Raw: "LIKE"},
		{Type: Identifier, Raw: "%go"},
		{Type: Or, Raw: "OR"},
		{Type: Identifier, Raw: "name"},
		{Type: In, Raw: "IN"},
		{Type: OpenParen, Raw: "("},
		{Type: Subquery, Raw: "SELECT name FROM $GOPATH/src/github.com WHERE name RLIKE .*_test\\.go"},
		{Type: CloseParen, Raw: ")"},
	}

	for i := range expected {
		if !reflect.DeepEqual(actual[i], expected[i]) {
			t.Fatalf("\nExpected: %v\n     Got: %v", expected[i], actual[i])
		}
	}
}

func TestTokenizer_ReadWord(t *testing.T) {
	type Case struct {
		input    string
		expected string
	}

	cases := []Case{
		{input: "foo", expected: "foo"},
		{input: "foo bar", expected: "foo"},
		{input: "", expected: ""},
	}

	for _, c := range cases {
		actual := NewTokenizer(c.input).readWord()
		if !reflect.DeepEqual(actual, c.expected) {
			t.Fatalf("\nExpected: %v\n     Got: %v", c.expected, actual)
		}
	}
}

func TestTokenizer_ReadQuery(t *testing.T) {
	type Case struct {
		input    string
		expected string
	}

	// TODO: Complete these cases.
	cases := []Case{}

	for _, c := range cases {
		actual := NewTokenizer(c.input).readQuery()
		if !reflect.DeepEqual(actual, c.expected) {
			t.Fatalf("\nExpected: %v\n     Got: %v", c.expected, actual)
		}
	}
}

func TestTokenizer_ReadUntil(t *testing.T) {
	type Case struct {
		input    string
		until    []rune
		expected string
	}

	// TODO: Complete these cases.
	cases := []Case{}

	for _, c := range cases {
		actual := NewTokenizer(c.input).readUntil(c.until...)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Fatalf("\nExpected: %v\n     Got: %v", c.expected, actual)
		}
	}
}
