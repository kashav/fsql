package tokenizer

import "testing"

func TestToken_TokenTypeString(t *testing.T) {
	type Case struct {
		tt       TokenType
		expected string
	}

	cases := []Case{
		Case{tt: Identifier, expected: "identifier"},
		Case{tt: Subquery, expected: "subquery"},
		Case{tt: Select, expected: "select"},
		Case{tt: From, expected: "from"},
		Case{tt: As, expected: "as"},
		Case{tt: Where, expected: "where"},
		Case{tt: Or, expected: "or"},
		Case{tt: And, expected: "and"},
		Case{tt: Not, expected: "not"},
		Case{tt: In, expected: "in"},
		Case{tt: Is, expected: "is"},
		Case{tt: Like, expected: "like"},
		Case{tt: RLike, expected: "RLike"},
		Case{tt: Equals, expected: "equal"},
		Case{tt: NotEquals, expected: "not-equal"},
		Case{tt: GreaterThanEquals, expected: "greater-than-or-equal"},
		Case{tt: GreaterThan, expected: "greater-than"},
		Case{tt: LessThanEquals, expected: "less-than-or-equal"},
		Case{tt: LessThan, expected: "less-than"},
		Case{tt: Comma, expected: "comma"},
		Case{tt: Hyphen, expected: "hyphen"},
		Case{tt: ExclamationMark, expected: "exclamation-mark"},
		Case{tt: OpenParen, expected: "open-parentheses"},
		Case{tt: CloseParen, expected: "close-parentheses"},
		Case{tt: OpenBracket, expected: "open-bracket"},
		Case{tt: CloseBracket, expected: "close-bracket"},
		Case{tt: Unknown, expected: "unknown"},
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
		Case{
			token:    Token{Type: Identifier, Raw: "name"},
			expected: "{type: identifier, raw: \"name\"}",
		},
		Case{
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
