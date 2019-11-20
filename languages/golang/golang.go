package golang

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/src-d/imports/languages/tsitter"
)

func init() {
	tsitter.RegisterLanguage(language{})
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
	tsitter.EachNodeTypes(n, func(n *sitter.Node) bool {
		c := n.Child(1)
		if c != nil && c.Type() == "import_spec_list" {
			n = c
		}
		for i := 0; i < int(n.ChildCount()); i++ {
			n := n.Child(i)
			if n.Type() != "import_spec" {
				continue
			}
			for j := 0; j < int(n.ChildCount()); j++ {
				n := n.Child(j)
				if n.Type() != "interpreted_string_literal" {
					continue
				}
				out = append(out, string(content[n.StartByte()+1:n.EndByte()-1]))
			}
		}
		return false
	}, "import_declaration")
	return out, nil
}
