package imports

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/src-d/enry/v2"
)

type Config struct {
	Out     io.Writer
	Num     int
	MaxSize int64
}

func NewExtractor(c Config) *Extractor {
	if c.Num == 0 {
		c.Num = runtime.NumGoroutine()
	}
	if c.Out == nil {
		c.Out = os.Stdout
	}
	if c.MaxSize == 0 {
		c.MaxSize = 10 * 1024 * 1024
	}
	return &Extractor{
		enc:     json.NewEncoder(c.Out),
		num:     c.Num,
		maxSize: c.MaxSize,
	}
}

type File struct {
	Path    string   `json:"file"`
	Lang    string   `json:"lang"`
	Imports []string `json:"imports,omitempty"`
}

type Extractor struct {
	enc     *json.Encoder
	num     int // TODO(dennwc): use it!
	maxSize int64
}

func (e *Extractor) Extract(root, rel string) error {
	// TODO(dennwc): expand relative imports and use dependency manifests in the future
	return filepath.Walk(filepath.Join(root, rel), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil // continue
		}
		if info.Size() > e.maxSize {
			return nil // skip
		}
		fname, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		return e.processFile(fname, path)
	})
}

func (e *Extractor) processFile(fname, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	lang := enry.GetLanguage(path, data)
	if lang == enry.OtherLanguage {
		// unknown language - skip
		return nil
	}
	l := LanguageByName(lang)
	if l == nil {
		// emit the file-language mapping; no imports
		return e.enc.Encode(File{
			Path: fname,
			Lang: lang,
		})
	}
	// import extraction
	out, err := l.Imports(data)
	if err != nil {
		return err
	}
	sort.Strings(out)
	return e.enc.Encode(File{
		Path:    fname,
		Lang:    lang,
		Imports: out,
	})
}
