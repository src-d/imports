package java

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/src-d/imports/languages/tsitter"
)

func init() {
	tsitter.RegisterLanguage(language{})
}

const query = `
((import_declaration) @imp)
`

var q *sitter.Query

func init() {
	var err error
	q, err = sitter.NewQuery([]byte(query), java.GetLanguage())
	if err != nil {
		panic(err)
	}
}

type language struct{}

func (language) Aliases() []string {
	return []string{"Java"}
}

func (l language) GetLanguage() *sitter.Language {
	return java.GetLanguage()
}

func (l language) Imports(content []byte, n *sitter.Node) ([]string, error) {
	c := sitter.NewQueryCursor()
	c.Exec(q, n)
	var out []string
	for {
		m, ok := c.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			n := c.Node
			cnt := int(n.ChildCount())
			parts := make([]string, 0, cnt)
			for i := 0; i < cnt; i++ {
				c := n.Child(i)
				switch c.Type() {
				case "identifier":
					name := content[c.StartByte():c.EndByte()]
					parts = append(parts, string(name))
				case "asterisk":
					// we add an empty string, so instead of "a.b.*" we return an "a.b." import
					parts = append(parts, "")
				}
			}
			out = append(out, strings.Join(parts, "."))
		}
	}
	return out, nil
}
