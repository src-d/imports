package python

import (
	"log"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/src-d/imports/languages/tsitter"
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
	log.Println(n)
	var out []string
	return out, nil
}
