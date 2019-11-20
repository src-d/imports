package python

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestPythonImports(t *testing.T) {
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "simple",
			Src:  `import a.b`,
			Exp: []string{
				"a/b",
			},
		},
		{
			Name: "named",
			Src:  `import a.b as c`,
			Exp: []string{
				"a.b",
			},
		},
		{
			Name: "relative",
			Src:  `from a.b import x`,
			Exp: []string{
				"./a/b",
			},
		},
		{
			Name: "relative dot",
			Src:  `from .a.b import x`,
			Exp: []string{
				"./a/b",
			},
		},
		{
			Name: "relative dot 2",
			Src:  `from ..a.b import x`,
			Exp: []string{
				"../a/b",
			},
		},
		{
			Name: "relative dot 3",
			Src:  `from ...a.b import x`,
			Exp: []string{
				"../../a/b",
			},
		},
		{
			Name: "relative dots",
			Src:  `from ... import x`,
			Exp: []string{
				"../../",
			},
		},
		{
			Name: "relative named",
			Src:  `from a.b import x as b`,
			Exp: []string{
				"./a/b",
			},
		},
		{
			Name: "relative symbols",
			Src:  `from a.b import x, y`,
			Exp: []string{
				"./a/b",
			},
		},
		{
			Name: "wildcard",
			Src:  `from a.b import *`,
			Exp: []string{
				"./a/b",
			},
		},
		{
			Name: "dynamic",
			Src:  `__import__('/a/b/c')`,
			Exp:  []string{
				//"/a/b/c", // TODO: dynamic imports
			},
		},
	})
}
