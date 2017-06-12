package parser

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
)

// Run parses the input string and returns the parsed AST (query).
func Run(input string) (*query.Query, error) {
	return (&parser{}).parse(input)
}

type parser struct {
	tokenizer *tokenizer.Tokenizer
	current   *tokenizer.Token
	expected  tokenizer.TokenType
}

// parse runs the respective parser function on each clause of the query.
func (p *parser) parse(input string) (*query.Query, error) {
	p.tokenizer = tokenizer.NewTokenizer(input)
	q := query.NewQuery()

	if err := p.parseSelectClause(q); err != nil {
		return nil, err
	}

	if err := p.parseFromClause(q); err != nil {
		return nil, err
	}

	if err := p.parseWhereClause(q); err != nil {
		return nil, err
	}

	return q, nil
}

// parseSelectClause parses the SELECT clause of the query.
func (p *parser) parseSelectClause(q *query.Query) error {
	// Determine if we should show all attributes. This is only true when
	// no attributes are provided (regardless of if the SELECT keyword is
	// provided or not).
	var showAll = true
	if p.expect(tokenizer.Select) == nil {
		if p.current == nil || p.current.Type == tokenizer.Identifier {
			showAll = false
		} else if p.current.Type == tokenizer.From || p.current.Type == tokenizer.Where {
			// No SELECT and next token is FROM/WHERE, show all!
			showAll = true
		} else {
			// No SELECT and next token is not Identifier nor FROM/WHERE -> malformed
			// input.
			return p.currentError()
		}
	} else if current := p.expect(tokenizer.Identifier); current != nil {
		p.current = current
		showAll = false
	}

	if showAll {
		q.Attributes = allAttributes
	} else if err := p.parseAttrs(&q.Attributes, &q.Modifiers); err != nil {
		return err
	}

	return nil
}

// parseFromClause parses the FROM clause of the query.
func (p *parser) parseFromClause(q *query.Query) error {
	if p.expect(tokenizer.From) == nil {
		err := p.currentError()
		if p.expect(tokenizer.Identifier) != nil {
			// No FROM, but an identifier -> malformed query.
			return err
		}

		// No specified directory, so we default to the CWD.
		q.Sources["include"] = append(q.Sources["include"], ".")
		return nil
	}

	if err := p.parseSourceList(&q.Sources, &q.SourceAliases); err != nil {
		return err
	}

	// Replace the tilde with the home directory in each source directory. This
	// is only required when the query is wrapped in quotes, since the shell
	// will automatically expand tildes otherwise.
	u, err := user.Current()
	if err != nil {
		return err
	}
	for _, sourceType := range []string{"include", "exclude"} {
		for i, src := range q.Sources[sourceType] {
			if strings.Contains(src, "~") {
				q.Sources[sourceType][i] = filepath.Join(u.HomeDir, src[1:])
			}
		}
	}

	return nil
}

// parseWhereClause parses the WHERE clause of the query.
func (p *parser) parseWhereClause(q *query.Query) error {
	if p.expect(tokenizer.Where) == nil {
		err := p.currentError()
		if p.expect(tokenizer.Identifier) == nil {
			return nil
		}
		return err
	}
	root, err := p.parseConditionTree()
	if err != nil {
		return err
	}
	q.ConditionTree = root

	return nil
}

// expect returns the next token if it matches the expectation t, and
// nil otherwise.
func (p *parser) expect(t tokenizer.TokenType) *tokenizer.Token {
	p.expected = t

	if p.current == nil {
		p.current = p.tokenizer.Next()
	}

	if p.current != nil && p.current.Type == t {
		tok := p.current
		p.current = nil
		return tok
	}

	return nil
}
