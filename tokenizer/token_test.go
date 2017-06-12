package tokenizer

import "testing"

func TestToken_TokenTypeString(t *testing.T) {
	type Case struct {
		tt       TokenType
		expected string
	}

	cases := []Case{
		{tt: Identifier, expected: "identifier"},
		{tt: Subquery, expected: "subquery"},
		{tt: Select, expected: "select"},
		{tt: From, expected: "from"},
		{tt: As, expected: "as"},
		{tt: Where, expected: "where"},
		{tt: Or, expected: "or"},
		{tt: And, expected: "and"},
		{tt: Not, expected: "not"},
		{tt: In, expected: "in"},
		{tt: Is, expected: "is"},
		{tt: Like, expected: "like"},
		{tt: RLike, expected: "RLike"},
		{tt: Equals, expected: "equal"},
		{tt: NotEquals, expected: "not-equal"},
		{tt: GreaterThanEquals, expected: "greater-than-or-equal"},
		{tt: GreaterThan, expected: "greater-than"},
		{tt: LessThanEquals, expected: "less-than-or-equal"},
		{tt: LessThan, expected: "less-than"},
		{tt: Comma, expected: "comma"},
		{tt: Hyphen, expected: "hyphen"},
		{tt: ExclamationMark, expected: "exclamation-mark"},
		{tt: OpenParen, expected: "open-parentheses"},
		{tt: CloseParen, expected: "close-parentheses"},
		{tt: OpenBracket, expected: "open-bracket"},
		{tt: CloseBracket, expected: "close-bracket"},
		{tt: Unknown, expected: "unknown"},
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
		{
			token:    Token{Type: Identifier, Raw: "name"},
			expected: "{type: identifier, raw: \"name\"}",
		},
		{
			token:    Token{Type: Comma, Raw: ","},
			expected: "{type: comma, raw: \",\"}",
		},
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
