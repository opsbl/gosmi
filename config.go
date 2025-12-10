package gosmi

import (
	"os"

	"github.com/opsbl/gosmi/smi"
	"github.com/opsbl/gosmi/types"
)

type FS = smi.FS
type NamedFS = smi.NamedFS

func NewNamedFS(name string, fs FS) NamedFS { return smi.NewNamedFS(name, fs) }

// Instance is the gosmi entrypoint bound to a single SMI handle.
type Instance struct {
	smiInst *smi.Instance
}

func (i *Instance) createModule(m *types.SmiModule) SmiModule {
	return CreateModule(i, m)
}

func (i *Instance) createNode(n *types.SmiNode) SmiNode {
	return CreateNode(i, n)
}

func (i *Instance) createType(t *types.SmiType) SmiType {
	return CreateType(i, t)
}

// New creates a new gosmi instance.
func New(tag ...string) (*Instance, error) {
	inst, err := smi.NewInstance(tag...)
	if err != nil {
		return nil, err
	}
	return &Instance{smiInst: inst}, nil
}

// Must is a helper that panics on failure.
func Must(tag ...string) *Instance {
	inst, err := New(tag...)
	if err != nil {
		panic(err)
	}
	return inst
}

// Close releases resources for this instance.
func (i *Instance) Close() {
	if i == nil {
		return
	}
	i.smiInst.Close()
}

func (i *Instance) GetPath() string         { return i.smiInst.GetPath() }
func (i *Instance) SetPath(path string)     { i.smiInst.SetPath(path) }
func (i *Instance) AppendPath(path string)  { i.smiInst.SetPath(string(os.PathListSeparator) + path) }
func (i *Instance) PrependPath(path string) { i.smiInst.SetPath(path + string(os.PathListSeparator)) }

func (i *Instance) SetFS(fs ...smi.NamedFS)     { i.smiInst.SetFS(fs...) }
func (i *Instance) AppendFS(fs ...smi.NamedFS)  { i.smiInst.AppendFS(fs...) }
func (i *Instance) PrependFS(fs ...smi.NamedFS) { i.smiInst.PrependFS(fs...) }
