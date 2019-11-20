package tsitter

import sitter "github.com/smacker/go-tree-sitter"

func EachNode(root *sitter.Node, fnc func(n *sitter.Node) bool) {
	if !fnc(root) {
		return
	}
	for i := 0; i < int(root.ChildCount()); i++ {
		EachNode(root.Child(i), fnc)
	}
}

func EachNodeTypes(root *sitter.Node, fnc func(n *sitter.Node) bool, types ...string) {
	EachNode(root, func(n *sitter.Node) bool {
		typ := n.Type()
		if typ == "" {
			return true // continue
		}
		for _, t := range types {
			if t == typ {
				return fnc(n)
			}
		}
		return true // continue
	})
}
