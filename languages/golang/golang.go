package golang

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/src-d/imports/languages/tsitter"
)

func init() {
	tsitter.RegisterLanguage(language{})
}

const query = `
(import_spec (interpreted_string_literal) @name)
`

var q *sitter.Query

func init() {
	var err error
	q, err = sitter.NewQuery([]byte(query), golang.GetLanguage())
	if err != nil {
		panic(err)
	}
}

type language struct{}

func (language) Aliases() []string {
	return []string{"Go"}
}

func (l language) GetLanguage() *sitter.Language {
	return golang.GetLanguage()
}

func (l language) Imports(content []byte, n *sitter.Node) ([]string, error) {
	var out []string
	c := sitter.NewQueryCursor()
	c.Exec(q, n)
	for {
		m, ok := c.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			str := string(content[c.Node.StartByte()+1 : c.Node.EndByte()-1])
			out = append(out, str)
		}
	}
	return out, nil
}
