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
		Case{input: "SELECT", expected: Select},
		Case{input: "FROM", expected: From},
		Case{input: "WHERE", expected: Where},
		Case{input: "AS", expected: As},
		Case{input: "OR", expected: Or},
		Case{input: "AND", expected: And},
		Case{input: "NOT", expected: Not},
		Case{input: "IN", expected: In},
		Case{input: "IS", expected: Is},
		Case{input: "LIKE", expected: Like},
		Case{input: "RLIKE", expected: RLike},
		Case{input: "foo", expected: Identifier},
		Case{input: "(", expected: OpenParen},
		Case{input: ")", expected: CloseParen},
		Case{input: ",", expected: Comma},
		Case{input: "-", expected: Hyphen},
		Case{input: "=", expected: Equals},
		Case{input: "<>", expected: NotEquals},
		Case{input: "<", expected: LessThan},
		Case{input: "<=", expected: LessThanEquals},
		Case{input: ">", expected: GreaterThan},
		Case{input: ">=", expected: GreaterThanEquals},
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
		Case{input: "foo", expected: "foo"},
		Case{input: " foo ", expected: "foo"},
		Case{input: "\" foo \"", expected: " foo "},
		Case{input: "' foo '", expected: " foo "},
		Case{input: "` foo `", expected: " foo "},
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
		Token{Type: Select, Raw: "SELECT"},
		Token{Type: Identifier, Raw: "name"},
		Token{Type: Comma, Raw: ","},
		Token{Type: Identifier, Raw: "size"},
		Token{Type: From, Raw: "FROM"},
		Token{Type: Identifier, Raw: "~/Desktop"},
		Token{Type: Where, Raw: "WHERE"},
		Token{Type: Identifier, Raw: "name"},
		Token{Type: Like, Raw: "LIKE"},
		Token{Type: Identifier, Raw: "%go"},
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
		Token{Type: Select, Raw: "SELECT"},
		Token{Type: Identifier, Raw: "name"},
		Token{Type: Comma, Raw: ","},
		Token{Type: Identifier, Raw: "size"},
		Token{Type: From, Raw: "FROM"},
		Token{Type: Identifier, Raw: "~/Desktop"},
		Token{Type: Where, Raw: "WHERE"},
		Token{Type: Identifier, Raw: "name"},
		Token{Type: Like, Raw: "LIKE"},
		Token{Type: Identifier, Raw: "%go"},
		Token{Type: Or, Raw: "OR"},
		Token{Type: Identifier, Raw: "name"},
		Token{Type: In, Raw: "IN"},
		Token{Type: OpenParen, Raw: "("},
		Token{Type: Subquery, Raw: "SELECT name FROM $GOPATH/src/github.com WHERE name RLIKE .*_test\\.go"},
		Token{Type: CloseParen, Raw: ")"},
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
		Case{input: "foo", expected: "foo"},
		Case{input: "foo bar", expected: "foo"},
		Case{input: "", expected: ""},
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
