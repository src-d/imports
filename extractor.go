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
	// Out is a destination to write JSON output to during Extract.
	Out io.Writer
	// Num is the maximal number of goroutines when extracting imports.
	Num int
	// MaxSize is the maximal size of files in bytes that will be parsed.
	// For files larger than this, only a sample of this size will be used for language detection.
	// Library may use the sample to try extracting imports, or may return an empty list of imports.
	MaxSize int64
}

// NewExtractor creates an extractor with a given configuration. See Config for more details.
func NewExtractor(c Config) *Extractor {
	if c.Num == 0 {
		c.Num = runtime.NumGoroutine()
	}
	if c.Out == nil {
		c.Out = os.Stdout
	}
	if c.MaxSize == 0 {
		c.MaxSize = 1 * 1024 * 1024
	}
	return &Extractor{
		enc:     json.NewEncoder(c.Out),
		num:     c.Num,
		maxSize: c.MaxSize,
	}
}

type File struct {
	Path    string   `json:"file"`
	Lang    string   `json:"lang,omitempty"`
	Imports []string `json:"imports,omitempty"`
}

type Extractor struct {
	enc     *json.Encoder
	num     int // TODO(dennwc): use it!
	maxSize int64
}

// Extract imports recursively from a given directory. The root is a root of the project's repository and rel is the
// relative path inside it that will be processed. Two paths exists to allow the library to potentially parse dependency
// manifest files that are usually located in the root of the project.
func (e *Extractor) Extract(root, rel string) error {
	// TODO(dennwc): expand relative imports and use dependency manifests in the future
	return filepath.Walk(filepath.Join(root, rel), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil // continue
		}
		sample := info.Size() > e.maxSize
		fname, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		return e.processFile(fname, path, sample)
	})
}

// ExtractFrom extracts imports from a given file content, assuming it had a given path.
//
// The path is used for language detection only. It won't access any files locally and won't fetch dependency manifests
// as Extract may do.
func (e *Extractor) ExtractFrom(path string, content []byte) (*File, error) {
	lang := enry.GetLanguage(path, content)
	f := &File{Path: path, Lang: lang}
	if lang == enry.OtherLanguage {
		// unknown language - skip
		return f, nil
	}
	l := LanguageByName(lang)
	if l == nil {
		// emit the file-language mapping; no imports
		return f, nil
	}
	// import extraction
	list, err := l.Imports(content)
	if err != nil {
		return f, err
	}
	sort.Strings(list)
	f.Imports = list
	return f, nil
}

func (e *Extractor) processFile(fname, path string, sample bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var data []byte
	if sample {
		data = make([]byte, e.maxSize)
		_, err = io.ReadFull(f, data)
	} else {
		data, err = ioutil.ReadAll(f)
	}
	if err != nil {
		return err
	}
	_ = f.Close()
	return e.processAndEmit(fname, path, data)
}

func (e *Extractor) processAndEmit(fname, path string, data []byte) error {
	f, err := e.ExtractFrom(path, data)
	if err != nil {
		return err
	}
	f.Path = fname
	return e.enc.Encode(f)
}
