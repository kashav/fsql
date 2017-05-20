package tokenizer

import "fmt"

// TokenType represents a Token's type.
type TokenType int8

// All TokenTypes.
const (
	Unknown TokenType = iota

	Identifier
	Subquery

	Select
	From
	Where

	As
	Or
	And
	Not
	In
	Is
	Like
	RLike

	Equals
	NotEquals
	GreaterThanEquals
	GreaterThan
	LessThanEquals
	LessThan

	OpenParen
	CloseParen
	Comma
	Hyphen
)

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "identifier"
	case Subquery:
		return "subquery"
	case Select:
		return "select"
	case From:
		return "from"
	case As:
		return "As"
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
	case OpenParen:
		return "open-parentheses"
	case CloseParen:
		return "close-parentheses"
	case Comma:
		return "comma"
	case Hyphen:
		return "hyphen"
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
	return fmt.Sprintf("{type: %s, raw: \"%s\", previous: %s}",
		t.Type.String(), t.Raw, t.Previous.String())
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
