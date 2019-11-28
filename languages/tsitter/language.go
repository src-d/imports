package tsitter

import (
	"runtime"
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
	runtime.KeepAlive(p)
	sort.Strings(out)
	return dedup(out), err
}

func dedup(sorted []string) []string {
	j, n := 0, len(sorted)
	if n == 0 {
		return []string{}
	}

	for i := 1; i < n; i++ {
		if sorted[j] == sorted[i] {
			continue
		}
		j++
		sorted[j] = sorted[i]
	}
	return sorted[:j+1]
}
