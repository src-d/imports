package tstest

import (
	"testing"

	"github.com/src-d/imports/languages/langtest"
	"github.com/src-d/imports/languages/tsitter"
)

type Case = langtest.Case

func RunTest(t *testing.T, l tsitter.Language, cases []Case) {
	t.Helper()
	langtest.RunTest(t, tsitter.NewLanguage(l), cases)
}
