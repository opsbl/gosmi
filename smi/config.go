package smi

import (
	"fmt"
	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

var DefaultSmiHandle *internal.Handle

type Handle = internal.Handle
type FS = internal.FS
type NamedFS = internal.NamedFS

func NewNamedFS(name string, fs FS) NamedFS { return NamedFS{Name: "[" + name + "]", FS: fs} }

func NewSmiHandle() *internal.Handle {
	return internal.NewHandle()
}

func checkInit() {
	if DefaultSmiHandle == nil {
		Init()
	}
}

// Init C -> int smiInit(const char *tag)
func Init(tag ...string) bool {
	if DefaultSmiHandle != nil {
		panic("repeat initialization")
	}
	DefaultSmiHandle = NewSmiHandle()
	var configTag, handleName string
	_ = handleName // Not used yet
	if len(tag) > 0 {
		configTag = tag[0]
		handleName = strings.Join(tag, ":")
	}
	// Set to built-in default path, if not Windows
	if runtime.GOOS != "windows" {
		DefaultSmiHandle.SetPath(DefaultSmiPaths...)
	}

	// Read global config file, if we can
	_ = ReadConfig(DefaultGlobalConfig, configTag)

	// Read user config file, if we can
	if homedir, err := os.UserHomeDir(); err == nil {
		_ = ReadConfig(filepath.Join(homedir, DefaultUserConfig), configTag)
	}
	// Use SMIPATH environment variable, if set
	SetPath(os.Getenv("SMIPATH"))
	return true
}

// Exit C -> void smiExit(void)
func Exit() {

}

// SetErrorLevel C -> void smiSetErrorLevel(int level)
func SetErrorLevel(level int) {
	checkInit()
	DefaultSmiHandle.SetErrorLevel(level)
}

// GetFlags C -> int smiGetFlags(void)
func GetFlags() int {
	checkInit()
	return int(DefaultSmiHandle.GetFlags())
}

// SetFlags C -> void smiSetFlags(int userflags)
func SetFlags(userflags int) {
	checkInit()
	DefaultSmiHandle.SetFlags(internal.Flags(userflags))
}

// GetPath C -> char *smiGetPath(void)
func GetPath() string {
	checkInit()
	return DefaultSmiHandle.GetPath()
}

// SetPath -> int smiSetPath(const char *path)
func SetPath(path string) {
	paths := filepath.SplitList(path)
	if len(paths) == 0 {
		return
	}
	DefaultSmiHandle.SetPath(paths...)
}

// SetSeverity C -> void smiSetSeverity(char *pattern, int severity)
func SetSeverity(pattern string, severity int) {
	checkInit()
	DefaultSmiHandle.SetSeverity(pattern, severity)
}

// ReadConfig -> int smiReadConfig(const char *filename, const char *tag)
func ReadConfig(filename string, tag ...string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Open file: %w", err)
	}
	defer f.Close()
	// TODO: Parse file
	return nil
}

// SetErrorHandler C -> void smiSetErrorHandler(SmiErrorHandler smiErrorHandler)
func SetErrorHandler(smiErrorHandler types.SmiErrorHandler) {
	checkInit()
	DefaultSmiHandle.SetErrorHandler(smiErrorHandler)
}

func SetFS(fs ...NamedFS)     { DefaultSmiHandle.SetFS(fs...) }
func AppendFS(fs ...NamedFS)  { DefaultSmiHandle.AppendFS(fs...) }
func PrependFS(fs ...NamedFS) { DefaultSmiHandle.PrependFS(fs...) }
