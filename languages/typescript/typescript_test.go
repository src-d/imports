package typescript

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestTypeScriptImports(t *testing.T) {
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "simple",
			Src:  `import { x } from "./file";`,
			Exp: []string{
				"./file",
			},
		},
		{
			Name: "simple 2",
			Src:  `import x from "file";`,
			Exp: []string{
				"file",
			},
		},
		{
			Name: "named",
			Src:  `import { x as y } from "./file";`,
			Exp: []string{
				"./file",
			},
		},
		{
			Name: "wildcard named",
			Src:  `import * as x from "./file";`,
			Exp: []string{
				"./file",
			},
		},
		{
			Name: "export",
			Src:  `export { x } from "./file";`,
			Exp: []string{
				"./file",
			},
		},
		{
			Name: "class",
			Src:  `export class foo {}`,
			Exp:  nil,
		},
		{
			Name: "export wildcard",
			Src:  `export * from "./file";`,
			Exp: []string{
				"./file",
			},
		},
		{
			Name: "export named",
			Src:  `export { x as y } from "./file";`,
			Exp: []string{
				"./file",
			},
		},
		{
			Name: "require",
			Src:  `import x = require("./file");`,
			Exp: []string{
				"./file",
			},
		},
		{
			Name: "side effect",
			Src:  `import "./file.js";`,
			Exp: []string{
				"./file.js",
			},
		},
	})
}
