package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// GetMacro C -> SmiMacro *smiGetMacro(SmiModule *smiModulePtr, char *macro)
func GetMacro(smiModulePtr *types.SmiModule, macro string) *types.SmiMacro {
	return smiHandle.GetMacro(smiModulePtr, macro)
}

// GetFirstMacro C -> SmiMacro *smiGetFirstMacro(SmiModule *smiModulePtr)
func GetFirstMacro(smiModulePtr *types.SmiModule) *types.SmiMacro {
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

// GetNextMacro C -> SmiMacro *smiGetNextMacro(SmiMacro *smiMacroPtr)
func GetNextMacro(smiMacroPtr *types.SmiMacro) *types.SmiMacro {
	if smiMacroPtr == nil {
		return nil
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	if macroPtr.Next == nil {
		return nil
	}
	return &macroPtr.Next.SmiMacro
}

// GetMacroModule C -> SmiModule *smiGetMacroModule(SmiMacro *smiMacroPtr)
func GetMacroModule(smiMacroPtr *types.SmiMacro) *types.SmiModule {
	if smiMacroPtr == nil {
		return nil
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	if macroPtr.Module == nil {
		return nil
	}
	return &macroPtr.Module.SmiModule
}

// GetMacroLine C -> int smiGetMacroLine(SmiMacro *smiMacroPtr)
func GetMacroLine(smiMacroPtr *types.SmiMacro) int {
	if smiMacroPtr == nil {
		return 0
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	return macroPtr.Line
}
