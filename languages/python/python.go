package python

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/src-d/imports/languages/tsitter"
)

const (
	query = `
		((import_statement name: (dotted_name (identifier)) @name))
		((import_from_statement module_name: (dotted_name (identifier)) @name))
		((import_from_statement module_name: (relative_import (import_prefix)) @name))
		((import_from_statement module_name: (relative_import (import_prefix) (dotted_name (identifier))) @name))
		((aliased_import name: (dotted_name (identifier)) @name))
	`
)

var (
	q *sitter.Query
)

func init() {
	tsitter.RegisterLanguage(language{})

	var err error
	if q, err = sitter.NewQuery([]byte(query), python.GetLanguage()); err != nil {
		panic(err)
	}
}

type language struct{}

func (language) Aliases() []string {
	return []string{"Python"}
}

func (l language) GetLanguage() *sitter.Language {
	return python.GetLanguage()
}

func (l language) Imports(content []byte, root *sitter.Node) ([]string, error) {
	var out []string
	c := sitter.NewQueryCursor()
	c.Exec(q, root)
	for {
		m, ok := c.NextMatch()
		if !ok {
			break
		}

		for _, cap := range m.Captures {
			str := string(content[cap.Node.StartByte():cap.Node.EndByte()])
			out = append(out, str)
		}
	}

	return out, nil
}
