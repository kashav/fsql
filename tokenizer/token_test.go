package tokenizer

import "testing"

func TestToken_TokenTypeString(t *testing.T) {
	type Case struct {
		tt       TokenType
		expected string
	}

	cases := []Case{
		{Identifier, "identifier"},
		{Subquery, "subquery"},
		{Select, "select"},
		{From, "from"},
		{As, "as"},
		{Where, "where"},
		{Or, "or"},
		{And, "and"},
		{Not, "not"},
		{In, "in"},
		{Is, "is"},
		{Like, "like"},
		{RLike, "RLike"},
		{Equals, "equal"},
		{NotEquals, "not-equal"},
		{GreaterThanEquals, "greater-than-or-equal"},
		{GreaterThan, "greater-than"},
		{LessThanEquals, "less-than-or-equal"},
		{LessThan, "less-than"},
		{Comma, "comma"},
		{Hyphen, "hyphen"},
		{ExclamationMark, "exclamation-mark"},
		{OpenParen, "open-parentheses"},
		{CloseParen, "close-parentheses"},
		{OpenBracket, "open-bracket"},
		{CloseBracket, "close-bracket"},
		{Unknown, "unknown"},
	}

	for _, c := range cases {
		actual := c.tt.String()
		if c.expected != actual {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}

func TestToken_String(t *testing.T) {
	type Case struct {
		token    Token
		expected string
	}

	cases := []Case{
		{Token{Type: Identifier, Raw: "name"}, "{type: identifier, raw: \"name\"}"},
		{Token{Type: Comma, Raw: ","}, "{type: comma, raw: \",\"}"},
	}

	for _, c := range cases {
		actual := c.token.String()
		if c.expected != actual {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}

func TestToken_NewTokenizer(t *testing.T) {
	input := "SELECT all FROM ."
	inputLength, tokenLength := len([]rune(input)), 0

	tokenizer := NewTokenizer(input)
	if len(tokenizer.input) != inputLength {
		t.Fatalf("len(tokenizer.input)\nExpected %v\n     Got %v", inputLength,
			len(tokenizer.input))
	}
	if len(tokenizer.tokens) != tokenLength {
		t.Fatalf("len(tokenizer.tokens)\nExpected %v\n     Got %v", tokenLength,
			len(tokenizer.tokens))
	}
}
