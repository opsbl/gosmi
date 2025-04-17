package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// GetType C -> SmiType *smiGetType(SmiModule *smiModulePtr, char *type)
func GetType(smiModulePtr *types.SmiModule, typeName string) *types.SmiType {
	return DefaultSmiHandle.GetType(smiModulePtr, typeName)
}

// GetFirstType C -> SmiType *smiGetFirstType(SmiModule *smiModulePtr)
func GetFirstType(smiModulePtr *types.SmiModule) *types.SmiType {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	typePtr := modulePtr.Types.First
	if typePtr == nil {
		return nil
	}
	return &typePtr.SmiType
}

// GetNextType C -> SmiType *smiGetNextType(SmiType *smiTypePtr)
func GetNextType(smiTypePtr *types.SmiType) *types.SmiType {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Next == nil {
		return nil
	}
	return &typePtr.Next.SmiType
}

// GetParentType C -> SmiType *smiGetParentType(SmiType *smiTypePtr)
func GetParentType(smiTypePtr *types.SmiType) *types.SmiType {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Parent == nil {
		return nil
	}
	return &typePtr.Parent.SmiType
}

// GetTypeModule C -> SmiModule *smiGetTypeModule(SmiType *smiTypePtr)
func GetTypeModule(smiTypePtr *types.SmiType) *types.SmiModule {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Module == nil {
		return nil
	}
	return &typePtr.Module.SmiModule
}

// GetTypeLine C -> int smiGetTypeLine(SmiType *smiTypePtr)
func GetTypeLine(smiTypePtr *types.SmiType) int {
	if smiTypePtr == nil {
		return 0
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	return typePtr.Line
}
