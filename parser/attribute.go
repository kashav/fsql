package parser

import (
	"strings"

	"github.com/kshvmdn/fsql/query"
	"github.com/kshvmdn/fsql/tokenizer"
)

var allAttributes = map[string]bool{
	"mode": true,
	"name": true,
	"size": true,
	"time": true,
}

// parseAttrList parses the list of attributes passed to the SELECT clause.
func (p *parser) parseAttrList(attributes *map[string]bool,
	modifiers *map[string][]query.Modifier) error {
	attribute := p.expect(tokenizer.Identifier)
	if attribute == nil {
		return p.currentError()
	}

	if attribute.Raw == "*" || attribute.Raw == "all" {
		*attributes = allAttributes
	} else {
		p.current = attribute
		attributeModifiers := make([]query.Modifier, 0)
		attribute, err := p.parseAttrModifiers(&attributeModifiers)
		if err != nil {
			return err
		}
		if attribute == nil {
			return p.currentError()
		}

		if _, ok := allAttributes[attribute.Raw]; !ok {
			return &ErrUnknownToken{attribute.Raw}
		}

		(*attributes)[attribute.Raw] = true
		(*modifiers)[attribute.Raw] = attributeModifiers
	}

	// If next token is a comma, recurse!
	if p.expect(tokenizer.Comma) != nil {
		return p.parseAttrList(attributes, modifiers)
	}

	return nil
}

// parseAttrModifiers parses an attribute's associated modifiers and
// returns the attribute.
func (p *parser) parseAttrModifiers(modifiers *[]query.Modifier) (*tokenizer.Token, error) {
	ident := p.expect(tokenizer.Identifier)
	if ident == nil {
		return nil, p.currentError()
	}

	if token := p.expect(tokenizer.OpenParen); token == nil {
		// No modifier on this attribute
		if _, ok := allAttributes[ident.Raw]; !ok {
			return nil, &ErrUnknownToken{ident.Raw}
		}
		return ident, nil
	}

	current := query.Modifier{
		Name:      strings.ToUpper(ident.Raw),
		Arguments: make([]string, 0),
	}

	attribute, err := p.parseAttrModifiers(modifiers)
	if err != nil {
		return nil, err
	}
	if attribute == nil {
		return nil, p.currentError()
	}

	for {
		if token := p.expect(tokenizer.Identifier); token != nil {
			current.Arguments = append(current.Arguments, token.Raw)
			continue
		}

		if token := p.expect(tokenizer.Comma); token != nil {
			continue
		}

		if token := p.expect(tokenizer.CloseParen); token != nil {
			*modifiers = append(*modifiers, current)
			return attribute, nil
		}
	}
}
