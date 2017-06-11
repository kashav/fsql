package query

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kshvmdn/fsql/transform"
)

// Modifier represents an attribute modifier.
type Modifier struct {
	Name      string
	Arguments []string
}

func (m *Modifier) String() string {
	return fmt.Sprintf("%s(%s)", m.Name, strings.Join(m.Arguments, ", "))
}

// applyModifiers iterates through each SELECT attribute for this query
// and applies the associated modifier to the attribute's output value.
func (q *Query) applyModifiers(path string, info os.FileInfo) map[string]interface{} {
	results := make(map[string]interface{}, len(q.Attributes))

	for attribute := range q.Attributes {
		var err error
		value, err := transform.DefaultFormatValue(attribute, path, info)

		if _, ok := q.Modifiers[attribute]; !ok {
			results[attribute] = value
			continue
		}

		for _, m := range q.Modifiers[attribute] {
			value, err = transform.Format(&transform.FormatParams{
				Attribute: attribute,
				Path:      path,
				Info:      info,
				Value:     value,
				Name:      m.Name,
				Args:      m.Arguments,
			})
			if err != nil {
				log.Fatalln(err.Error())
			}
		}

		results[attribute] = value
	}

	return results
}
