package java

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestJavaImports(t *testing.T) {
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "simple",
			Src:  `import com.foo.bar;`,
			Exp: []string{
				"com/foo/bar",
			},
		},
		{
			Name: "static",
			Src:  `import static com.foo.Bar;`,
			Exp: []string{
				"com/foo/Bar",
			},
		},
		{
			Name: "wildcard",
			Src:  `import com.foo.*;`,
			Exp: []string{
				"com/foo",
			},
		},
		{
			Name: "multiple",
			Src: `
import com.foo.*;
import com.foo.B;
import com.foo.A;
`,
			Exp: []string{
				"com/foo",
				"com/foo/A",
				"com/foo/B",
			},
		},
	})
}
