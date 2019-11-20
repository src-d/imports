package golang

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestGoImports(t *testing.T) {
	const pre = "package foo\n\n"
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "simple",
			Src:  pre + `import "fmt"`,
			Exp: []string{
				"fmt",
			},
		},
		{
			Name: "multiple",
			Src: pre + `
import "fmt"
import "os"
`,
			Exp: []string{
				"fmt",
				"os",
			},
		},
		{
			Name: "group",
			Src: pre + `
import (
	"fmt"
	"os"
)
`,
			Exp: []string{
				"fmt",
				"os",
			},
		},
		{
			Name: "external",
			Src:  pre + `import "github.com/user/repo"`,
			Exp: []string{
				"github.com/user/repo",
			},
		},
		{
			Name: "relative",
			Src:  pre + `import "./subdir"`,
			Exp: []string{
				"./subdir",
			},
		},
		{
			Name: "named",
			Src:  pre + `import f "fmt"`,
			Exp: []string{
				"fmt",
			},
		},
		{
			Name: "dot",
			Src:  pre + `import . "fmt"`,
			Exp: []string{
				"fmt",
			},
		},
		{
			Name: "side effect",
			Src:  pre + `import _ "fmt"`,
			Exp: []string{
				"fmt",
			},
		},
	})
}
