package query

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents a Token's type.
type TokenType int8

const (
	// Unknown represents an unknown Token type.
	Unknown TokenType = iota
	// Subquery represents a subquery in this query string
	Subquery
	// Select represents the SELECT clause.
	Select
	// From represents the FROM clause.
	From
	// Where represents the WHERE clause.
	Where
	// Or represents the OR keyword for conditional disjunction.
	Or
	// And represents the AND keyword for conditonal conjunction.
	And
	// Not represents the NOT keyword for conditional negation.
	Not
	// In represents the IN keyword for list-based comparisons.
	In
	// Is represents the IS keyword for file type comparisons.
	Is
	// Like represents the LIKE keyword for string comparisons.
	Like
	// RLike represents the RLIKE keyword for string regexp comparisons.
	RLike
	// Identifier represents the value for each Query.
	Identifier
	// OpenParen represents an open parenthesis.
	OpenParen
	// CloseParen represents a closed parenthesis.
	CloseParen
	// Comma represents a comma.
	Comma
	// Minus represents the `-` operator for directory exclusion.
	Minus
	// Equals represents the `=` comparator for string/numeric comparisons.
	Equals
	// NotEquals represents the `<>` comparator for string/numeric comparisons.
	NotEquals
	// GreaterThanEquals represents the `>=` comparator for numeric comparisons.
	GreaterThanEquals
	// GreaterThan represents the `>` comparator for numeric comparisons.
	GreaterThan
	// LessThanEquals represents the `<=` comparator for numeric comparisons.
	LessThanEquals
	// LessThan represents the `<` comparator for numeric comparisons.
	LessThan
)

func (t TokenType) String() string {
	switch t {
	case Subquery:
		return "subquery"
	case Select:
		return "select"
	case From:
		return "from"
	case Where:
		return "where"
	case Or:
		return "or"
	case And:
		return "and"
	case Not:
		return "not"
	case In:
		return "in"
	case Is:
		return "is"
	case Like:
		return "like"
	case RLike:
		return "RLike"
	case Identifier:
		return "identifier"
	case OpenParen:
		return "open-parentheses"
	case CloseParen:
		return "close-parentheses"
	case Comma:
		return "comma"
	case Minus:
		return "minus"
	case Equals:
		return "equal"
	case NotEquals:
		return "not-equal"
	case GreaterThanEquals:
		return "greater-than-or-equal"
	case GreaterThan:
		return "greater-than"
	case LessThanEquals:
		return "less-than-or-equal"
	case LessThan:
		return "less-than"
	default:
		return "unknown"
	}
}

// Token represents a single token.
type Token struct {
	Type     TokenType
	Raw      string
	Previous *Token
}

func (t Token) String() string {
	return fmt.Sprintf("{type: %s, raw: \"%s\"}", t.Type.String(), t.Raw)
}

// Tokenizer represents a token worker.
type Tokenizer struct {
	input    []rune
	previous *Token
}

// NewTokenizer initializes a new Tokenizer.
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input: []rune(input), previous: nil}
}

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
		return t.setNextToken(&Token{Type: OpenParen, Raw: "("})

	case ')':
		t.input = t.input[1:]
		return t.setNextToken(&Token{Type: CloseParen, Raw: ")"})

	case ',':
		t.input = t.input[1:]
		return t.setNextToken(&Token{Type: Comma, Raw: ","})

	case '-':
		t.input = t.input[1:]
		return t.setNextToken(&Token{Type: Minus, Raw: "-"})

	case '=':
		t.input = t.input[1:]
		return t.setNextToken(&Token{Type: Equals, Raw: "="})

	case '>':
		if t.getRuneAtIndex(1) == '=' {
			t.input = t.input[2:]
			return t.setNextToken(&Token{Type: GreaterThanEquals, Raw: ">="})
		}

		t.input = t.input[1:]
		return t.setNextToken(&Token{Type: GreaterThan, Raw: ">"})

	case '<':
		if t.getRuneAtIndex(1) == '=' {
			t.input = t.input[2:]
			return t.setNextToken(&Token{Type: LessThanEquals, Raw: ">="})
		}

		if t.getRuneAtIndex(1) == '>' {
			t.input = t.input[2:]
			return t.setNextToken(&Token{Type: NotEquals, Raw: "<>"})
		}

		t.input = t.input[1:]
		return t.setNextToken(&Token{Type: LessThan, Raw: "<"})
	}

	if !t.currentIs(-1, ',', '\'', '"', '`', '(', ')') {
		word := t.readWord()
		tok := &Token{Raw: word}

		switch strings.ToUpper(word) {
		case "SELECT":
			tok.Type = Select
		case "FROM":
			tok.Type = From
		case "WHERE":
			tok.Type = Where
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
		case "RLIKE":
			tok.Type = RLike
		default:
			tok.Type = Identifier
		}

		// If the previous token was a `(`, and the one before was IN, then
		// this must be a subquery. Keep reading until we reach a `)`.
		if t.previous != nil && t.previous.Type == OpenParen &&
			t.previous.Previous != nil && t.previous.Previous.Type == In {
			tok.Type = Subquery
			tok.Raw = t.readQuery(word)
		}

		return t.setNextToken(tok)
	}

	tok := &Token{Type: Unknown, Raw: string(current)}

	// If the current rune is a single/double quote or backtick, we want to keep
	// reading until we reach the closing symbol.
	if t.currentIs('\'', '"', '`') {
		t.input = t.input[1:]
		tok.Raw = t.readUntil(t.readWord(), current)
		tok.Type = Identifier
	}

	t.input = t.input[1:]
	return t.setNextToken(tok)
}

// Return the rune at the ith index of the input.
func (t *Tokenizer) getRuneAtIndex(i int) rune {
	if len(t.input) == i {
		return -1
	}

	return t.input[i]
}

// Return the rune at the 0th index of the input.
func (t *Tokenizer) current() rune {
	return t.getRuneAtIndex(0)
}

// Returns true iff the input's current rune (at index 0) is in rs.
func (t *Tokenizer) currentIs(rs ...rune) bool {
	for _, r := range rs {
		if r == t.current() {
			return true
		}
	}
	return false
}

// Update the previous Token for both the Tokenizer and the supplied Token.
func (t *Tokenizer) setNextToken(token *Token) *Token {
	token.Previous, t.previous = t.previous, token
	return token
}

// Read a single word from the input. Returns when the next rune is any
// of: -1, " ", comma, single/double quote, backtick, or parenthesis.
func (t *Tokenizer) readWord() string {
	word := []rune{}

	for {
		if unicode.IsSpace(t.current()) ||
			t.currentIs(-1, ',', '\'', '"', '`', '(', ')') {
			return string(word)
		}

		word = append(word, t.current())
		t.input = t.input[1:]
	}
}

// Read a full string until we reaching a closing parentheses. Maintains a
// count of opening parens to ensure we don't return early.
func (t *Tokenizer) readQuery(start string) string {
	query := fmt.Sprintf("%s ", start)

	var count = 1
	for count > 0 {
		for unicode.IsSpace(t.current()) {
			t.input = t.input[1:]
		}

		word := fmt.Sprintf("%s ", t.readWord())

		if t.current() == '(' {
			count++
			word = "("
		} else if t.current() == ')' {
			count--
			word = ")"
		}

		if t.current() == -1 || count <= 0 {
			break
		}

		query += word
		t.input = t.input[1:]
	}

	return query
}

// Read the input starting at start, until reaching a rune in runes.
func (t *Tokenizer) readUntil(start string, runes ...rune) string {
	word := start
	for !t.currentIs(runes...) {
		for unicode.IsSpace(t.current()) {
			t.input = t.input[1:]
		}
		word = fmt.Sprintf("%s %s", word, t.readWord())
	}
	return word
}
