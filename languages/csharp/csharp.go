package csharp

import (
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
	q *sitter.Query

	// include filters
	inFilters = map[string]struct{}{
		"identifier_name": struct{}{},
		".":               struct{}{},
	}

	// exclude filters
	exFilters = map[string]struct{}{
		"type_argument_list": struct{}{},
		"<":                  struct{}{},
		">":                  struct{}{},
	}
)

func init() {
	tsitter.RegisterLanguage(language{})

	var err error
	if q, err = sitter.NewQuery([]byte(query), csharp.GetLanguage()); err != nil {
		panic(err)
	}
}

type language struct{}

func (language) Aliases() []string {
	return []string{"C#"}
}

func (l language) GetLanguage() *sitter.Language {
	return csharp.GetLanguage()
}

func (l language) Imports(content []byte, root *sitter.Node) ([]string, error) {
	var out []string

	c := sitter.NewQueryCursor()
	c.Exec(q, root)
	for {
		m, ok := c.NextMatch()
		if !ok {
			break
		}

		for _, cap := range m.Captures {
			out = append(out, filter(content, cap.Node))
		}
	}

	return out, nil
}

func filter(content []byte, node *sitter.Node) string {
	var (
		out string
		fn  func(n *sitter.Node)
	)

	fn = func(n *sitter.Node) {
		// try to include the whole node (identifiers)
		if _, in := inFilters[n.Type()]; in {
			str := string(content[n.StartByte():n.EndByte()])
			out += str
		} else {
			// otherwise, try include children
			cnt := int(n.ChildCount())
			for i := 0; i < cnt; i++ {
				// exclude the whole node if not interested
				if _, ex := exFilters[n.Type()]; ex {
					continue
				}
				fn(n.Child(i))
			}
		}
	}

	fn(node)
	return out
}
