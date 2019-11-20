package tsitter

import (
	"sort"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/src-d/imports"
)

func RegisterLanguage(l Language) {
	imports.RegisterLanguage(NewLanguage(l))
}

type Language interface {
	Aliases() []string
	GetLanguage() *sitter.Language
	Imports(content []byte, root *sitter.Node) ([]string, error)
}

func NewLanguage(l Language) imports.Language {
	return language{l: l}
}

type language struct {
	l Language
}

func (l language) Aliases() []string {
	return l.l.Aliases()
}

func (l language) Imports(content []byte) ([]string, error) {
	p := sitter.NewParser()
	p.SetLanguage(l.l.GetLanguage())
	tree := p.Parse(content)
	out, err := l.l.Imports(content, tree.RootNode())
	sort.Strings(out)
	return out, err
}
