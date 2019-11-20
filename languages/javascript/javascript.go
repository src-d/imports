package javascript

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/src-d/imports/languages/tsitter"
)

func init() {
	tsitter.RegisterLanguage(language{})
}

type language struct{}

func (language) Aliases() []string {
	return []string{"JavaScript"}
}

func (l language) GetLanguage() *sitter.Language {
	return javascript.GetLanguage()
}

func (l language) Imports(content []byte, n *sitter.Node) ([]string, error) {
	var out []string
	tsitter.EachNodeTypes(n, func(n *sitter.Node) bool {
		c := n.ChildByFieldName("source")
		if c == nil || c.Type() != "string" {
			return false
		}
		src := string(content[c.StartByte()+1 : c.EndByte()-1])
		if !strings.HasPrefix(src, "./") && !strings.HasPrefix(src, "/") {
			src = "./" + src
		}
		out = append(out, src)
		return false // don't descend into the import
	}, "import_statement")
	return out, nil
}
