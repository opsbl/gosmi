package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handle) Exit() {
	if h == nil {
		return
	}
	h.freeData()
}

func (h *Handle) GetPath() string {
	names := make([]string, len(h.Paths))
	for i, fs := range h.Paths {
		names[i] = fs.Name
	}
	return strings.Join(names, string(os.PathListSeparator))
}

func expandPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("Path is empty")
	}
	if path[0] == '~' {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("Cannot expand homedir")
		}
		path = filepath.Join(homedir, path[1:])
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Get absolute path for '%s': %w", path, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("Cannot stat '%s': %w", path, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("'%s' is not a directory", path)
	}
	return path, nil
}

func (h *Handle) SetPath(path ...string) {
	pathLen := len(path)
	if pathLen == 0 {
		return
	}
	if path[0] == "" {
		h.appendPath(path[1:]...)
	} else if path[pathLen-1] == "" {
		h.prependPath(path[:pathLen-1]...)
	} else {
		h.Paths = make([]NamedFS, 0, pathLen)
		for _, p := range path {
			if p, err := expandPath(p); err == nil {
				h.Paths = append(h.Paths, newPathFS(p))
			}
		}
	}
}

func (h *Handle) appendPath(path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]NamedFS, len(h.Paths), len(h.Paths)+len(path))
	copy(paths, h.Paths)
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, newPathFS(p))
		}
	}
	h.Paths = paths
}

func (h *Handle) prependPath(path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]NamedFS, 0, len(h.Paths)+len(path))
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, newPathFS(p))
		}
	}
	paths = append(paths, h.Paths...)
	h.Paths = paths
}
