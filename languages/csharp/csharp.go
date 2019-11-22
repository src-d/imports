package csharp

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
	"github.com/src-d/imports/languages/tsitter"
)

const (
	query = `
		(using_directive (identifier_name) @name)
		(using_directive (qualified_name) @name)
	`
)

var (
	q *sitter.Query
)

func init() {
	tsitter.RegisterLanguage(language{})

	var err error
	if q, err = sitter.NewQuery([]byte(query), csharp.GetLanguage()); err != nil {
		panic(err)
	}
}

type language struct{}

func (language) Aliases() []string {
	return []string{"C#"}
}

func (l language) GetLanguage() *sitter.Language {
	return csharp.GetLanguage()
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
