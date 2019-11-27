package imports

import "context"

const Separator = "/"

type Language interface {
	Aliases() []string
	Imports(ctx context.Context, content []byte) ([]string, error)
}

var (
	languages = make(map[string]Language)
)

func RegisterLanguage(lang Language) {
	for _, name := range lang.Aliases() {
		languages[name] = lang
	}
}

func LanguageByName(name string) Language {
	return languages[name]
}
