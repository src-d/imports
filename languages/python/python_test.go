package python

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestPythonImports(t *testing.T) {
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "no imports",
			Src:  "",
			Exp:  []string{},
		},
		{
			Name: "simple",
			Src:  `import a.b`,
			Exp: []string{
				"a.b",
			},
		},
		{
			Name: "named",
			Src: `
				import a.b as c
				import A.B as C
			`,
			Exp: []string{
				"A.B",
				"a.b",
			},
		},
		{
			Name: "relative",
			Src:  `from a.b import x`,
			Exp: []string{
				"a.b",
			},
		},
		{
			Name: "relative dot",
			Src:  `from .a.b import x`,
			Exp: []string{
				".a.b",
			},
		},
		{
			Name: "relative dot 2",
			Src:  `from ..a.b import x`,
			Exp: []string{
				"..a.b",
			},
		},
		{
			Name: "relative dot 3",
			Src:  `from ...a.b import x`,
			Exp: []string{
				"...a.b",
			},
		},
		{
			Name: "relative dots",
			Src:  `from ... import x`,
			Exp: []string{
				"...",
			},
		},
		{
			Name: "relative named",
			Src: `
				from a.b import x as b
			`,
			Exp: []string{
				"a.b",
			},
		},
		{
			Name: "relative named",
			Src: `
				from a.b import x as b
				from A.B import X as B
			`,
			Exp: []string{
				"A.B",
				"a.b",
			},
		},
		{
			Name: "relative symbols",
			Src:  `from a.b import x, y`,
			Exp: []string{
				"a.b",
			},
		},
		{
			Name: "dedup",
			Src: `
			from a.b import foo,bar
			from a.b import baz
			`,
			Exp: []string{
				"a.b",
			},
		},
		{
			Name: "dynamic",
			Src: `
				__import__('/a/b/c')
				another = import_file('relative_subdir/another.py')
				`,
			Exp: []string{
				"/a/b/c",
				"relative_subdir/another.py",
			},
		},
		{
			Name: "all",
			Src: `
				import a.b
				import c.d as x
				from ... import x
				from e.f import x as g
				from h.i import x, y
				from j.k import x,y
				from a.b import baz
				__import__('/a/b/c')
				another = import_file('d/e/f/g')
				`,
			Exp: []string{
				"...",
				"/a/b/c",
				"a.b",
				"c.d",
				"d/e/f/g",
				"e.f",
				"h.i",
				"j.k",
			},
		},
	})
}
