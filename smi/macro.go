package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// SmiMacro *smiGetMacro(SmiModule *smiModulePtr, char *macro)
func (i *Instance) GetMacro(smiModulePtr *types.SmiModule, macro string) *types.SmiMacro {
	if macro == "" {
		return nil
	}

	var result *types.SmiMacro
	i.withHandle(func(h *internal.Handle) {
		var modulePtr *internal.Module
		if smiModulePtr != nil {
			modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
			macroPtr := modulePtr.Macros.GetName(macro)
			if macroPtr != nil {
				result = &macroPtr.SmiMacro
			}
			return
		}
		for modulePtr = h.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
			macroPtr := modulePtr.Macros.GetName(macro)
			if macroPtr != nil {
				result = &macroPtr.SmiMacro
				return
			}
		}
	})
	return result
}

// SmiMacro *smiGetFirstMacro(SmiModule *smiModulePtr)
func (i *Instance) GetFirstMacro(smiModulePtr *types.SmiModule) *types.SmiMacro {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	macroPtr := modulePtr.Macros.First
	if macroPtr == nil {
		return nil
	}
	return &macroPtr.SmiMacro
}

// SmiMacro *smiGetNextMacro(SmiMacro *smiMacroPtr)
func (i *Instance) GetNextMacro(smiMacroPtr *types.SmiMacro) *types.SmiMacro {
	if smiMacroPtr == nil {
		return nil
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	if macroPtr.Next == nil {
		return nil
	}
	return &macroPtr.Next.SmiMacro
}

// SmiModule *smiGetMacroModule(SmiMacro *smiMacroPtr)
func (i *Instance) GetMacroModule(smiMacroPtr *types.SmiMacro) *types.SmiModule {
	if smiMacroPtr == nil {
		return nil
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	if macroPtr.Module == nil {
		return nil
	}
	return &macroPtr.Module.SmiModule
}

// int smiGetMacroLine(SmiMacro *smiMacroPtr)
func (i *Instance) GetMacroLine(smiMacroPtr *types.SmiMacro) int {
	if smiMacroPtr == nil {
		return 0
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	return macroPtr.Line
}
