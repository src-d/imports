package csharp

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestCSharpImports(t *testing.T) {
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "simple",
			Src:  `using Foo.Bar;`,
			Exp: []string{
				"Foo/Bar",
			},
		},
		{
			Name: "static",
			Src:  `using static Foo.Bar;`,
			Exp: []string{
				"Foo/Bar",
			},
		},
		{
			Name: "multiple",
			Src: `
using Foo.Bar.A;
using Foo.Bar.B;
`,
			Exp: []string{
				"Foo/Bar/A",
				"Foo/Bar/B",
			},
		},
		{
			Name: "named",
			Src:  `using f = Foo.Bar;`,
			Exp: []string{
				"Foo/Bar",
			},
		},
		{
			Name: "named template",
			Src:  `using f = Foo.Bar<int>;`,
			Exp: []string{
				"Foo/Bar",
			},
		},
		{
			Name: "use named",
			Src: `
using f = Foo.Bar;
using f.B;
`,
			Exp: []string{
				"Foo/Bar",
				"f/B", // TODO: must be "Foo/Bar/B", but we need to resolve symbols for this
			},
		},
	})
}
