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

// Set the supplied token's Previous to the t's previous and set t's previous
// to token.
func (t *Tokenizer) setAndReturnToken(token *Token) *Token {
	token.Previous, t.previous = t.previous, token
	return token
}

// Read input, starting at startWord, until reaching a rune in chars.
func (t *Tokenizer) readUntil(start *string, chars ...rune) {
	in := func(c rune) bool {
		for _, char := range chars {
			if c == char {
				return true
			}
		}
		return false
	}

	for !in(t.current()) {
		for unicode.IsSpace(t.current()) {
			t.input = t.input[1:]
		}

		*start = fmt.Sprintf("%s %s", *start, t.readWord())
	}
}

// Next gets the next Token in this Tokenizer.
func (t *Tokenizer) Next() *Token {
	for {
		if !unicode.IsSpace(t.current()) {
			break
		}

		t.input = t.input[1:]
	}

	current := t.current()
	if current == -1 {
		return nil
	}

	switch current {
	case '(':
		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: OpenParen, Raw: "("})

	case ')':
		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: CloseParen, Raw: ")"})

	case ',':
		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: Comma, Raw: ","})

	case '-':
		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: Minus, Raw: "-"})

	case '=':
		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: Equals, Raw: "="})

	case '>':
		if t.nextRune() == '=' {
			t.input = t.input[2:]
			return t.setAndReturnToken(&Token{Type: GreaterThanEquals, Raw: ">="})
		}

		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: GreaterThan, Raw: ">"})

	case '<':
		if t.nextRune() == '=' {
			t.input = t.input[2:]
			return t.setAndReturnToken(&Token{Type: LessThanEquals, Raw: ">="})
		}

		if t.nextRune() == '>' {
			t.input = t.input[2:]
			return t.setAndReturnToken(&Token{Type: NotEquals, Raw: "<>"})
		}

		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: LessThan, Raw: "<"})
	}

	if !(current == -1 || current == '`' || current == '\'' || current == '"' ||
		current == ',' || current == '(' || current == ')') {
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
		//
		// FIXME: We aren't counting parens, so as soon as we we reach a closing
		// paren, we stop. This breaks when the subquery also has parens in it.
		if t.previous != nil && t.previous.Type == OpenParen &&
			t.previous.Previous != nil && t.previous.Previous.Type == In {
			t.input = t.input[1:]
			t.readUntil(&word, ')')
			tok.Type = Subquery
			tok.Raw = word
		}

		return t.setAndReturnToken(tok)
	}

	t.input = t.input[1:]

	// If the current rune is a single/double quote or backtick, we want to keep
	// reading until we reach the closing symbol.
	//
	// FIXME: This doesn't actually check that the closing symbol matches the
	// opening one (i.e. <single-quote> ... <double-quote> works fine, which
	// is wrong!).
	if current == '\'' || current == '`' || current == '"' {
		word := t.readWord()
		t.readUntil(&word, '\'', '`', '"')
		t.input = t.input[1:]
		return t.setAndReturnToken(&Token{Type: Identifier, Raw: word})
	}

	return t.setAndReturnToken(&Token{
		Type: Unknown,
		Raw:  string([]rune{current})},
	)
}

func (t *Tokenizer) current() rune {
	if len(t.input) == 0 {
		return -1
	}

	return t.input[0]
}

func (t *Tokenizer) nextRune() rune {
	if len(t.input) == 1 {
		return -1
	}

	return t.input[1]
}

func (t *Tokenizer) readWord() string {
	word := []rune{}

	for {
		r := t.current()

		if r == -1 || unicode.IsSpace(r) || r == '`' || r == '\'' ||
			r == '"' || r == ',' || r == '(' || r == ')' {
			return string(word)
		}

		word = append(word, r)
		t.input = t.input[1:]
	}
}
