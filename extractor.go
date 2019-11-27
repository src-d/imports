package imports

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/src-d/enry/v2"
)

// Extract imports from a given file content, assuming it had a given path.
//
// The path is used for language detection only. It won't access any files locally and won't fetch dependency manifests
// as Extract may do.
//
// The function does no sampling for the file, it means if the file size is large, it may take a lot of time to parse
// it or even detect the language. Instead, passing a part of the file may be an option.
func Extract(path string, content []byte) (*File, error) {
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

type Config struct {
	// Out is a destination to write JSON output to during Extract.
	Out io.Writer
	// Num is the maximal number of goroutines when extracting imports.
	// Zero value means use NumCPU.
	Num int
	// MaxSize is the maximal size of files in bytes that will be parsed.
	// For files larger than this, only a sample of this size will be used for language detection.
	// Library may use the sample to try extracting imports, or may return an empty list of imports.
	MaxSize int64
	// SymLinks is option that allows traversal over sym-links, be aware of potentials loops
	// if false - all symlinks will be skipped
	SymLinks bool
}

// NewExtractor creates an extractor with a given configuration. See Config for more details.
func NewExtractor(c Config) *Extractor {
	if c.Num <= 0 {
		c.Num = runtime.NumCPU()
	}
	if c.Out == nil {
		c.Out = os.Stdout
	}
	if c.MaxSize == 0 {
		c.MaxSize = 1 * 1024 * 1024
	}
	return &Extractor{
		enc:      json.NewEncoder(c.Out),
		num:      c.Num,
		maxSize:  c.MaxSize,
		symLinks: c.SymLinks,
	}
}

type File struct {
	Path    string   `json:"file"`
	Lang    string   `json:"lang,omitempty"`
	Imports []string `json:"imports,omitempty"`
}

type Extractor struct {
	mu  sync.Mutex
	enc *json.Encoder

	num      int
	maxSize  int64
	symLinks bool
}

type extractJob struct {
	fname string
	path  string
	buf   *bytes.Buffer // sample buffer
}

// Extract imports recursively from a given directory. The root is a root of the project's repository and rel is the
// relative path inside it that will be processed. Two paths exists to allow the library to potentially parse dependency
// manifest files that are usually located in the root of the project.
func (e *Extractor) Extract(root, rel string) error {
	var (
		jobs      chan *extractJob
		errs      chan error
		wg        sync.WaitGroup
		sampleBuf *bytes.Buffer
	)
	if e.num != 1 {
		jobs = make(chan *extractJob, e.num)
		errs = make(chan error, e.num)
		for i := 0; i < e.num; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				e.worker(jobs, errs)
			}()
		}
		// no sample buffer, it's per-routine
	} else {
		sampleBuf = bytes.NewBuffer(nil)
		sampleBuf.Grow(int(e.maxSize))
	}
	// TODO(dennwc): expand relative imports and use dependency manifests in the future
	err := filepath.Walk(filepath.Join(root, rel), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() {
			return nil // continue
		}
		// get a relative path, since we must print it in the output
		fname, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink == 0 {
			// regular file
			args := extractJob{fname: fname, path: path, buf: sampleBuf}
			if e.num == 1 {
				// no concurrency - process inline
				return e.processFile(args)
			}
			// send a job to workers running in the background
			select {
			case jobs <- &args:
			case err = <-errs:
				return err
			}
			return nil
		}
		if !e.symLinks {
			// skip symlinks
			return nil
		}
		// read the symlink destination
		dst, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}

		// we calculate the symlink destination to make sure it's still inside the directory
		dst, err = filepath.Rel(root, dst)
		if err != nil {
			return err
		}
		if strings.HasPrefix(dst, "../") {
			return nil // skip
		}

		// TODO(lwsanty): infinite loop detection
		return e.Extract(root, dst)
	})
	if e.num == 1 {
		// no concurrency, don't wait for workers
		return err
	}
	// close the queue, wait for workers, check errors
	close(jobs)
	if err != nil {
		return err
	}
	wg.Wait()
	select {
	case err = <-errs:
		return err
	default:
	}
	return nil
}

// ExtractFrom extracts imports from a given file content, assuming it had a given path.
//
// The path is used for language detection only. It won't access any files locally and won't fetch dependency manifests
// as Extract may do.
func (e *Extractor) ExtractFrom(path string, content []byte) (*File, error) {
	return Extract(path, content)
}

func (e *Extractor) worker(jobs <-chan *extractJob, errc chan<- error) {
	// each worker has it's own buffer
	buf := bytes.NewBuffer(nil)
	buf.Grow(int(e.maxSize))

	for args := range jobs {
		args.buf = buf
		if err := e.processFile(*args); err != nil {
			errc <- err
			return
		}
	}
}

func (e *Extractor) processFile(args extractJob) error {
	if args.buf == nil {
		panic("buffer must be set")
	}
	f, err := os.Open(args.path)
	if err != nil {
		return err
	}
	defer f.Close()

	var r io.Reader
	if e.maxSize >= 0 {
		// don't read the whole content, only a part of it
		r = io.LimitReader(f, e.maxSize)
	}

	buf := args.buf
	buf.Reset()

	_, err = buf.ReadFrom(r)
	if err != nil {
		return err
	}
	_ = f.Close()

	return e.processAndEmit(args.fname, args.path, buf.Bytes())
}

func (e *Extractor) processAndEmit(fname, path string, data []byte) error {
	f, err := e.ExtractFrom(path, data)
	if err != nil {
		return err
	}
	f.Path = fname
	if e.num != 1 {
		e.mu.Lock()
		defer e.mu.Unlock()
	}
	return e.enc.Encode(f)
}
