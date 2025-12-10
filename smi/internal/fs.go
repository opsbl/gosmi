package internal

import (
	"os"
	"path/filepath"
)

type NamedFS struct {
	Name string
	FS   FS
}

type pathFS string

func newPathFS(path string) NamedFS {
	return NamedFS{path, pathFS(path)}
}

func (p pathFS) Open(name string) (File, error) {
	filename := filepath.Join(string(p), name)
	return os.Open(filename)
}

func (h *Handle) SetFS(fs ...NamedFS)     { h.Paths = fs }
func (h *Handle) AppendFS(fs ...NamedFS)  { h.Paths = append(h.Paths, fs...) }
func (h *Handle) PrependFS(fs ...NamedFS) { h.Paths = append(fs, h.Paths...) }
