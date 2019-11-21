package typescript

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
	"github.com/src-d/imports/languages/tsitter"
)

func init() {
	tsitter.RegisterLanguage(language{})
}

type language struct{}

func (language) Aliases() []string {
	return []string{"TypeScript"}
}

func (l language) GetLanguage() *sitter.Language {
	return typescript.GetLanguage()
}

func (l language) Imports(content []byte, n *sitter.Node) ([]string, error) {
	var out []string
	tsitter.EachNodeTypes(n, func(n *sitter.Node) bool {
		if n.Type() == "export_statement" {
			for i := 0; i < int(n.ChildCount()); i++ {
				n := n.Child(i)
				if n.Type() != "string" {
					continue
				}
				imp := string(content[n.StartByte()+1 : n.EndByte()-1])
				out = append(out, imp)
			}
			return false
		}
		for i := 0; i < int(n.ChildCount()); i++ {
			n := n.Child(i)
			switch n.Type() {
			case "string":
				imp := string(content[n.StartByte()+1 : n.EndByte()-1])
				out = append(out, imp)
			case "import_require_clause":
				for i := 0; i < int(n.ChildCount()); i++ {
					c := n.Child(i)
					if c.Type() != "string" {
						continue
					}
					imp := string(content[c.StartByte()+1 : c.EndByte()-1])
					out = append(out, imp)
				}
			}
		}
		return false // don't descend into the import
	}, "import_statement", "export_statement")
	return out, nil
}
