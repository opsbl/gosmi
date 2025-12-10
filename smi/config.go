package smi

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

const (
	DefaultErrorLevel   = 3
	DefaultGlobalConfig = "/etc/smi.conf"
	DefaultUserConfig   = ".smirc"
)

var DefaultSmiPaths []string = []string{
	"/usr/local/share/mibs/ietf",
	"/usr/local/share/mibs/iana",
	"/usr/local/share/mibs/irtf",
	"/usr/local/share/mibs/site",
	"/usr/local/share/mibs/jacobs",
	"/usr/local/share/mibs/tubs",
}

type FS = internal.FS
type NamedFS = internal.NamedFS

func NewNamedFS(name string, fs FS) NamedFS { return NamedFS{Name: "[" + name + "]", FS: fs} }

type Instance struct {
	handle *internal.Handle
	mu     sync.Mutex
}

func (i *Instance) withHandle(fn func(h *internal.Handle)) {
	if i == nil || i.handle == nil {
		return
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	fn(i.handle)
}

func (i *Instance) bootstrap(configTag string) {
	if runtime.GOOS != "windows" {
		i.withHandle(func(h *internal.Handle) {
			h.SetPath(DefaultSmiPaths...)
		})
	}
	_ = i.ReadConfig(DefaultGlobalConfig, configTag)
	if homedir, err := os.UserHomeDir(); err == nil {
		_ = i.ReadConfig(filepath.Join(homedir, DefaultUserConfig), configTag)
	}
	if envPath := os.Getenv("SMIPATH"); envPath != "" {
		i.SetPath(envPath)
	}
}

func newInstance(tag ...string) (*Instance, error) {
	var configTag, handleName string
	if len(tag) > 0 {
		configTag = tag[0]
		handleName = strings.Join(tag, ":")
	}
	handle, err := internal.NewHandle(handleName)
	if err != nil {
		return nil, err
	}
	inst := &Instance{handle: handle}
	inst.bootstrap(configTag)
	return inst, nil
}

// Init creates a new SMI context (preferred over legacy global API).
func Init(tag ...string) (*Instance, error) { return newInstance(tag...) }

// NewInstance creates an isolated SMI context.
func NewInstance(tag ...string) (*Instance, error) { return newInstance(tag...) }

func (i *Instance) Close() {
	if i == nil {
		return
	}
	i.withHandle(func(h *internal.Handle) {
		h.Exit()
	})
	i.handle = nil
}

// void smiSetErrorLevel(int level)
func (i *Instance) SetErrorLevel(level int) {
	i.withHandle(func(h *internal.Handle) {
		h.SetErrorLevel(level)
	})
}

// int smiGetFlags(void)
func (i *Instance) GetFlags() (flags int) {
	i.withHandle(func(h *internal.Handle) {
		flags = h.GetFlags()
	})
	return
}

// void smiSetFlags(int userflags)
func (i *Instance) SetFlags(userflags int) {
	i.withHandle(func(h *internal.Handle) {
		h.SetFlags(userflags)
	})
}

// char *smiGetPath(void)
func (i *Instance) GetPath() (path string) {
	i.withHandle(func(h *internal.Handle) {
		path = h.GetPath()
	})
	return
}

// int smiSetPath(const char *path)
func (i *Instance) SetPath(path string) {
	paths := filepath.SplitList(path)
	if len(paths) == 0 {
		return
	}
	i.withHandle(func(h *internal.Handle) {
		h.SetPath(paths...)
	})
}

// void smiSetSeverity(char *pattern, int severity)
func (i *Instance) SetSeverity(pattern string, severity int) {
	i.withHandle(func(h *internal.Handle) {
		h.SetSeverity(pattern, severity)
	})
}

// int smiReadConfig(const char *filename, const char *tag)
func (i *Instance) ReadConfig(filename string, tag ...string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Open file: %w", err)
	}
	defer f.Close()
	// TODO: Parse file
	return nil
}

// void smiSetErrorHandler(SmiErrorHandler smiErrorHandler)
func (i *Instance) SetErrorHandler(smiErrorHandler types.SmiErrorHandler) {
	i.withHandle(func(h *internal.Handle) {
		h.SetErrorHandler(smiErrorHandler)
	})
}

func (i *Instance) SetFS(fs ...NamedFS) {
	i.withHandle(func(h *internal.Handle) {
		h.SetFS(fs...)
	})
}

func (i *Instance) AppendFS(fs ...NamedFS) {
	i.withHandle(func(h *internal.Handle) {
		h.AppendFS(fs...)
	})
}

func (i *Instance) PrependFS(fs ...NamedFS) {
	i.withHandle(func(h *internal.Handle) {
		h.PrependFS(fs...)
	})
}
