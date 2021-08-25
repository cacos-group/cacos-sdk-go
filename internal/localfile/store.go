package localfile

import (
	"fmt"
	"os"
	"path/filepath"
)

type KV struct {
	Key string
	Val string
}

type Storer interface {
	Write(key string, val string) error
	ReadAll() ([]KV, error)
}

type store struct {
	path string
}

func New() Storer {
	return newStore()
}

func newStore() *store {
	s := &store{}
	s.path = "."
	return s
}

func (s *store) filename(key string) string {
	return fmt.Sprintf("%s/%s", s.path, key)
}

func (s *store) Write(key string, val string) error {
	if key == "" {
		return fmt.Errorf("key can't empty")
	}
	fname := filepath.Base(s.filename(key))
	if fname == "" {
		return fmt.Errorf("filename can't empty")
	}
	dir := filepath.Dir(s.path)
	fi, err := os.Stat(dir)
	if err == nil && !fi.IsDir() {
		return fmt.Errorf("%s already exists and not a directory", dir)
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create dir %s error: %s", dir, err.Error())
		}
	}

	current, err := newWrapFile(fname)
	if err != nil {
		return err
	}

	_, err = current.write([]byte(val))
	if err != nil {
		return err
	}

	current.fp.Close()

	return nil
}

func (s *store) ReadAll() ([]KV, error) {
	return nil, nil
}

type wrapFile struct {
	fsize int64
	fp    *os.File
}

func (w *wrapFile) size() int64 {
	return w.fsize
}

func (w *wrapFile) write(p []byte) (n int, err error) {
	n, err = w.fp.Write(p)
	w.fsize += int64(n)
	return
}

func newWrapFile(fpath string) (*wrapFile, error) {
	fp, err := os.OpenFile(fpath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	fi, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	return &wrapFile{fp: fp, fsize: fi.Size()}, nil
}
