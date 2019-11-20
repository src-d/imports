package imports

const Separator = "/"

type Language interface {
	Aliases() []string
	Imports(content []byte) ([]string, error)
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
