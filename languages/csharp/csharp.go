package csharp

import (
	"sync"

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
	once sync.Once
	q    *sitter.Query
)

func init() {
	tsitter.RegisterLanguage(language{})
}

type language struct{}

func (language) Aliases() []string {
	return []string{"C#"}
}

func (l language) GetLanguage() *sitter.Language {
	return csharp.GetLanguage()
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
