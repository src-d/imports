package langtest

import (
	"reflect"
	"testing"

	"github.com/src-d/imports"
)

type Case struct {
	Name string
	Skip bool
	Src  string
	Exp  []string
}

func RunTest(t *testing.T, l imports.Language, cases []Case) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Skip {
				t.SkipNow()
			}
			out, err := l.Imports([]byte(c.Src))
			if err != nil {
				t.Fatal(err)
			}
			if len(out) == 0 && len(c.Exp) == 0 {
				return
			}
			if !reflect.DeepEqual(c.Exp, out) {
				t.Fatal(out)
			}
		})
	}
}
