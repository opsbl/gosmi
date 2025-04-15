package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func expandPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("path is empty")
	}
	if path[0] == '~' {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("cannot expand homedir")
		}
		path = filepath.Join(homedir, path[1:])
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("get absolute path for '%s': %w", path, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("cannot stat '%s': %w", path, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("'%s' is not a directory", path)
	}
	return path, nil
}

func appendPath(smiHandle *Handle, path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]NamedFS, len(smiHandle.Paths), len(smiHandle.Paths)+len(path))
	copy(paths, smiHandle.Paths)
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, newPathFS(p))
		}
	}
	smiHandle.Paths = paths
}

func prependPath(smiHandle *Handle, path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]NamedFS, 0, len(smiHandle.Paths)+len(path))
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, newPathFS(p))
		}
	}
	paths = append(paths, smiHandle.Paths...)
	smiHandle.Paths = paths
}
