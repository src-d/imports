package python

import (
	"sync"

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
	once sync.Once
	q    *sitter.Query
)

func init() {
	tsitter.RegisterLanguage(language{})
}

type language struct{}

func (language) Aliases() []string {
	return []string{"Python"}
}

func (l language) GetLanguage() *sitter.Language {
	return python.GetLanguage()
}

func (l language) Imports(content []byte, n *sitter.Node) ([]string, error) {
	var err error
	once.Do(func() { q, err = sitter.NewQuery([]byte(query), l.GetLanguage()) })
	if err != nil {
		return nil, err
	}

	var out []string
	c := sitter.NewQueryCursor()
	c.Exec(q, n)
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
