package internal

import (
	"github.com/opsbl/gosmi/types"
	"unsafe"
)

type Macro struct {
	types.SmiMacro
	Module *Module
	Flags  Flags
	Next   *Macro
	Prev   *Macro
	Line   int
}

type MacroMap struct {
	First *Macro

	last *Macro
	m    map[types.SmiIdentifier]*Macro
}

func (x *MacroMap) Add(m *Macro) {
	m.Prev = x.last
	if x.First == nil {
		x.First = m
	} else {
		x.last.Next = m
	}
	x.last = m

	if x.m == nil {
		x.m = make(map[types.SmiIdentifier]*Macro)
	}
	x.m[m.Name] = m
}

func (x *MacroMap) Get(name types.SmiIdentifier) *Macro {
	if x.m == nil {
		return nil
	}
	return x.m[name]
}

func (x *MacroMap) GetName(name string) *Macro {
	return x.Get(types.SmiIdentifier(name))
}

// GetMacro -> C SmiMacro *smiGetMacro(SmiModule *smiModulePtr, char *macro)
func (h *Handle) GetMacro(smiModulePtr *types.SmiModule, macro string) *types.SmiMacro {
	if macro == "" {
		return nil
	}

	var modulePtr *Module
	if smiModulePtr != nil {
		modulePtr = (*Module)(unsafe.Pointer(smiModulePtr))
		macroPtr := modulePtr.Macros.GetName(macro)
		if macroPtr == nil {
			return nil
		}
		return &macroPtr.SmiMacro
	}
	for modulePtr = h.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
		macroPtr := modulePtr.Macros.GetName(macro)
		if macroPtr != nil {
			return &macroPtr.SmiMacro
		}
	}
	return nil
}
