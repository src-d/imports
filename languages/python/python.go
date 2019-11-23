package python

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/src-d/imports/languages/tsitter"
)

var (
	query = []string{
		`((import_statement name: (dotted_name (identifier)) @name))`,
		`((import_statement name: (aliased_import name: (dotted_name (identifier)) @name)))`,
		`((import_from_statement module_name: (dotted_name (identifier)) @name))`,
		`((import_from_statement module_name: (relative_import (import_prefix)) @name))`,
	}
	q []*sitter.Query

	dQuery    = `((call function: (identifier) @fn arguments: (argument_list (string) @arg)))`
	dFnFilter = map[string]struct{}{
		"__import__":  struct{}{},
		"import_file": struct{}{},
	}
	dq *sitter.Query
)

func init() {
	tsitter.RegisterLanguage(language{})

	var err error
	q = make([]*sitter.Query, len(query))
	for i, qi := range query {
		if q[i], err = sitter.NewQuery([]byte(qi), python.GetLanguage()); err != nil {
			panic(err)
		}
	}

	if dq, err = sitter.NewQuery([]byte(dQuery), python.GetLanguage()); err != nil {
		panic(err)
	}
}

type language struct{}

func (language) Aliases() []string {
	return []string{"Python"}
}

func (l language) GetLanguage() *sitter.Language {
	return python.GetLanguage()
}

func (l language) Imports(content []byte, root *sitter.Node) ([]string, error) {
	var out = []string{}

	// regular imports (import_statement, import_from_statement)
	for _, qi := range q {
		out = append(out, imports(content, root, qi)...)
	}

	// dynamic imports (call function)
	out = append(out, dImports(content, root, dq, dFnFilter)...)

	return out, nil
}

func imports(content []byte, root *sitter.Node, query *sitter.Query) []string {
	var out = []string{}

	c := sitter.NewQueryCursor()
	c.Exec(query, root)
	for {
		m, ok := c.NextMatch()
		if !ok {
			break
		}

		if len(m.Captures) == 1 {
			n := m.Captures[0].Node
			str := string(content[n.StartByte():n.EndByte()])
			out = append(out, str)
		}
	}

	return out
}

func dImports(content []byte, root *sitter.Node, query *sitter.Query, filter map[string]struct{}) []string {
	var out = []string{}

	c := sitter.NewQueryCursor()
	c.Exec(query, root)
	for {
		m, ok := c.NextMatch()
		if !ok {
			break
		}
		if len(m.Captures) == 2 {
			capFn, capArg := m.Captures[0], m.Captures[1]

			fn := string(content[capFn.Node.StartByte():capFn.Node.EndByte()])
			if _, ok := filter[fn]; ok {
				arg := string(content[capArg.Node.StartByte():capArg.Node.EndByte()])
				out = append(out, strings.Trim(arg, `'"`))
			}
		}
	}

	return out
}
