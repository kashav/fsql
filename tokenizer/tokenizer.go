package tokenizer

import (
	"fmt"
	"strings"
	"unicode"
)

// All parses all tokens for this Tokenizer.
func (t *Tokenizer) All() []Token {
	tokens := []Token{}
	for tok := t.Next(); tok != nil; tok = t.Next() {
		tokens = append(tokens, *tok)
	}
	return tokens
}

// Next finds and returns the next Token in the input string.
func (t *Tokenizer) Next() *Token {
	for unicode.IsSpace(t.current()) {
		t.input = t.input[1:]
	}

	current := t.current()
	if current == -1 {
		return nil
	}

	switch current {
	case '(':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: OpenParen, Raw: "("})
	case ')':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: CloseParen, Raw: ")"})
	case '[':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: OpenBracket, Raw: "["})
	case ']':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: CloseBracket, Raw: "]"})
	case ',':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: Comma, Raw: ","})
	case '-':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: Hyphen, Raw: "-"})
	case '!':
		if t.getRuneAt(1) == '=' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: NotEquals, Raw: "!="})
		}
		return t.setToken(&Token{Type: ExclamationMark, Raw: "!"})
	case '=':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: Equals, Raw: "="})
	case '>':
		if t.getRuneAt(1) == '=' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: GreaterThanEquals, Raw: ">="})
		}
		t.input = t.input[1:]
		return t.setToken(&Token{Type: GreaterThan, Raw: ">"})
	case '<':
		if t.getRuneAt(1) == '=' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: LessThanEquals, Raw: "<="})
		}
		if t.getRuneAt(1) == '>' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: NotEquals, Raw: "<>"})
		}
		t.input = t.input[1:]
		return t.setToken(&Token{Type: LessThan, Raw: "<"})
	}

	if !t.currentIs(-1, ',', '\'', '"', '`', '(', ')', '[', ']') {
		word := t.readWord()
		tok := &Token{Raw: word}

		switch strings.ToUpper(word) {
		case "SELECT":
			tok.Type = Select
		case "FROM":
			tok.Type = From
		case "WHERE":
			tok.Type = Where
		case "AS":
			tok.Type = As
		case "OR":
			tok.Type = Or
		case "AND":
			tok.Type = And
		case "NOT":
			tok.Type = Not
		case "IN":
			tok.Type = In
		case "IS":
			tok.Type = Is
		case "LIKE":
			tok.Type = Like
		case "REGEXP", "RLIKE":
			tok.Type = RLike
		default:
			tok.Type = Identifier
		}

		if t.getPreviousToken() != nil && t.getPreviousToken().Type == OpenParen &&
			t.getTokenAt(1) != nil && t.getTokenAt(1).Type == In {
			// The two previous tokens were: `IN` and `(`, so we're at a subquery.
			tok.Type = Subquery
			tok.Raw = fmt.Sprintf("%s %s", word, t.readQuery())
		}

		return t.setToken(tok)
	}

	tok := &Token{Type: Unknown, Raw: string(current)}

	// If the current rune is a single/double quote or backtick, we want to keep
	// reading until we reach the matching closing symbol.
	if t.currentIs('\'', '"', '`') {
		t.input = t.input[1:]
		tok.Raw = t.readWord() + t.readUntil(current)
		tok.Type = Identifier
	}

	t.input = t.input[1:]
	return t.setToken(tok)
}

// setToken adds token to the list of this Tokenizer's tokens.
func (t *Tokenizer) setToken(token *Token) *Token {
	t.tokens = append(t.tokens, token)
	return token
}

// getPreviousToken returns the token that was most-recently read.
func (t *Tokenizer) getPreviousToken() *Token {
	return t.getTokenAt(0)
}

// getTokenAt returns the token at index i from the end of the tokens slice.
func (t *Tokenizer) getTokenAt(i int) *Token {
	j := len(t.tokens) - 1 - i
	if j < 0 {
		return nil
	}

	return t.tokens[j]
}

// current returns the run at the 0th index of the input.
func (t *Tokenizer) current() rune {
	return t.getRuneAt(0)
}

// getRuneAt returns the rune at the ith index of the input.
func (t *Tokenizer) getRuneAt(i int) rune {
	if len(t.input) == i {
		return -1
	}

	return t.input[i]
}

// currentIs returns true iff the input's current rune (at index 0) is in runes.
func (t *Tokenizer) currentIs(runes ...rune) bool {
	for _, r := range runes {
		if r == t.current() {
			return true
		}
	}
	return false
}

// readWord reads a single word from the input. Returns when the next rune is
// any of: nil (-1), empty space, comma, single/double quote, backtick,
// or opening/closing parenthesis/bracket.
func (t *Tokenizer) readWord() string {
	word := []rune{}

	for {
		if unicode.IsSpace(t.current()) ||
			t.currentIs(-1, ',', '\'', '"', '`', '(', ')', '[', ']') {
			return string(word)
		}

		word = append(word, t.current())
		t.input = t.input[1:]
	}
}

// readQuery reads a full string until reaching a closing parentheses. Counts
// opening parens to ensure that balance is maintained.
func (t *Tokenizer) readQuery() string {
	var query string

	var count = 1
	for count > 0 {
		for unicode.IsSpace(t.current()) {
			t.input = t.input[1:]
		}

		word := fmt.Sprintf("%s", t.readWord())

		if t.current() == -1 {
			break
		}

		if t.current() == '(' {
			count++
			word = "("
		} else if t.current() == ')' {
			count--
			if count <= 0 {
				query += word
				break
			}
			word = ")"
		} else if t.currentIs('\'', '`') {
			word += string(t.current())
		} else {
			word += " "
		}

		query += word
		t.input = t.input[1:]
	}

	return query
}

// readUnitl reads the input starting at start, until reaching a rune in runes.
func (t *Tokenizer) readUntil(runes ...rune) string {
	var word string
	for !t.currentIs(runes...) {
		for unicode.IsSpace(t.current()) {
			t.input = t.input[1:]
		}
		word = fmt.Sprintf("%s %s", word, t.readWord())
	}
	return word
}
