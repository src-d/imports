package javascript

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestJavaScriptImports(t *testing.T) {
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "simple",
			Src:  `import 'file';`,
			Exp: []string{
				"file",
			},
		},
		{
			Name: "relative",
			Src:  `import './local/file';`,
			Exp: []string{
				"./local/file",
			},
		},
		{
			Name: "symbol",
			Src:  `import foo from 'file';`,
			Exp: []string{
				"file",
			},
		},
		{
			Name: "multiple symbols",
			Src:  `import { foo, bar } from 'file';`,
			Exp: []string{
				"file",
			},
		},
		{
			Name: "symbol aliases",
			Src:  `import { foo as f, bar as b } from 'file';`,
			Exp: []string{
				"file",
			},
		},
		{
			Name: "wildcard",
			Src:  `import * as ns from 'file';`,
			Exp: []string{
				"file",
			},
		},
		{
			Name: "type",
			Src:  `import type Foo from 'file';`,
			Exp: []string{
				"file",
			},
		},
	})
}
