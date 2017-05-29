package parser

import (
	"fmt"
	"path/filepath"

	"github.com/kshvmdn/fsql/tokenizer"
)

// parseSourceList parses the list of directories passed to the FROM clause. If
// a source is followed by the AS keyword, the following word is registered as
// an alias.
func (p *parser) parseSourceList(sources *map[string][]string,
	aliases *map[string]string) error {
	// If the next token is a hypen, exclude this directory.
	sourceType := "include"
	if p.expect(tokenizer.Hyphen) != nil {
		sourceType = "exclude"
	}

	source := p.expect(tokenizer.Identifier)
	if source == nil {
		return p.currentError()
	}

	// Normalize directory path.
	source.Raw = filepath.Clean(source.Raw)
	(*sources)[sourceType] = append((*sources)[sourceType], source.Raw)

	if p.expect(tokenizer.As) != nil {
		alias := p.expect(tokenizer.Identifier)
		if alias == nil {
			return p.currentError()
		}
		if sourceType == "exclude" {
			return fmt.Errorf("cannot alias excluded directory %s", source.Raw)
		}
		(*aliases)[alias.Raw] = source.Raw
	}

	// If next token is a comma, recurse!
	if p.expect(tokenizer.Comma) != nil {
		return p.parseSourceList(sources, aliases)
	}

	return nil
}
