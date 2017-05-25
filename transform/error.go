package transform

import (
	"fmt"
	"strings"
)

// ErrNotImplemented used for non-implemented modifier functions.
type ErrNotImplemented struct {
	Name      string
	Attribute string
}

func (e *ErrNotImplemented) Error() string {
	return fmt.Sprintf("function %s is not implemented for attribute %s",
		strings.ToUpper(e.Name), e.Attribute)
}

// ErrUnsupportedFormat used for unsupport arguments for FORMAT functions.
type ErrUnsupportedFormat struct {
	Format    string
	Attribute string
}

func (e *ErrUnsupportedFormat) Error() string {
	return fmt.Sprintf("unsupported format type %s for attribute %s",
		e.Format, e.Attribute)
}
