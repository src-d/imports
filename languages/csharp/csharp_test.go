package csharp

import (
	"testing"

	"github.com/src-d/imports/languages/tsitter/tstest"
)

func TestCSharpImports(t *testing.T) {
	tstest.RunTest(t, language{}, []tstest.Case{
		{
			Name: "no imports",
			Src:  "class Test {}",
			Exp:  []string{},
		},
		{
			Name: "simple",
			Src:  `using Foo;`,
			Exp: []string{
				"Foo",
			},
		},
		{
			Name: "simple qualified",
			Src:  `using Foo.Bar;`,
			Exp: []string{
				"Foo.Bar",
			},
		},
		{
			Name: "static",
			Src:  `using static Foo.Bar;`,
			Exp: []string{
				"Foo.Bar",
			},
		},
		{
			Name: "multiple",
			Src: `
				using Foo.Bar.A;
				using Foo.Bar.B;
				`,
			Exp: []string{
				"Foo.Bar.A",
				"Foo.Bar.B",
			},
		},
		{
			Name: "named",
			Src:  `using f = Foo.Bar;`,
			Exp: []string{
				"Foo.Bar",
			},
		},
		{
			Name: "named template",
			Src:  `using f = Foo.Bar<int>;`,
			Exp: []string{
				"Foo.Bar",
			},
		},
		{
			Name: "named template",
			Src:  `using f = Foo.Bar /* comment */ <int> ;`,
			Exp: []string{
				"Foo.Bar",
			},
		},
		{
			Name: "named template",
			Src: `
			using MyStack = NameSpace.Stack<MyType>;
			using MyPopAlias = NameSpace.Stack<MyType>.Pop;
			`,
			Exp: []string{
				"NameSpace.Stack",
				"NameSpace.Stack.Pop",
			},
		},
		{
			Name: "use named",
			Src: `
				using f = Foo.Bar;
				using f.B;
				`,
			Exp: []string{
				"Foo.Bar",
				"f.B", // TODO: must be "Foo/Bar/B", but we need to resolve symbols for this
			},
		},
		{
			Name: "with class",
			Src: `
			using System;
			using System.Collections.Generic;

			using static System.Math;

			// Using alias directive for a class.
			using AliasToMyClass = NameSpace1.MyClass;

			// Using alias directive for a generic class.
			using UsingAlias = NameSpace2.MyClass<int>;

			using MyStack = NameSpace.Stack<MyType>;
			using MyPopAlias = NameSpace.Stack<MyType>.Pop;

			using f = Foo.Bar /* comment */ <int> ;

			class Test {
				static void main() {
					using (var font1 = new Font("Arial", 10.0f))
					{
						byte charset = font1.GdiCharSet;
					}
				}
			}
			`,
			Exp: []string{
				"Foo.Bar",
				"NameSpace.Stack",
				"NameSpace.Stack.Pop",
				"NameSpace1.MyClass",
				"NameSpace2.MyClass",
				"System",
				"System.Collections.Generic",
				"System.Math",
			},
		},
		{
			Name: "dedup",
			Src: `
			using System;
			using System.Collections.Generic;
			using System;
			using static System.Math;
			using System.Collections.Generic;
			using System;
			using System.Collections.Generic;
			`,
			Exp: []string{
				"System",
				"System.Collections.Generic",
				"System.Math",
			},
		},
	})
}
