package csharp

import (
	"log"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/src-d/imports/languages/tsitter"
)

func init() {
	tsitter.RegisterLanguage(language{})
}

type language struct{}

func (language) Aliases() []string {
	return []string{"C#"}
}

func (l language) GetLanguage() *sitter.Language {
	return GetLanguage()
}

func (l language) Imports(content []byte, n *sitter.Node) ([]string, error) {
	log.Println(n)
	var out []string
	return out, nil
}
