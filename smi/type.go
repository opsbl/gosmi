package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// SmiType *smiGetType(SmiModule *smiModulePtr, char *type)
func (i *Instance) GetType(smiModulePtr *types.SmiModule, typeName string) *types.SmiType {
	if typeName == "" {
		return nil
	}

	var result *types.SmiType
	i.withHandle(func(h *internal.Handle) {
		var modulePtr *internal.Module
		if smiModulePtr != nil {
			modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
			typePtr := modulePtr.Types.GetName(typeName)
			if typePtr != nil {
				result = &typePtr.SmiType
			}
			return
		}
		for modulePtr = h.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
			typePtr := modulePtr.Types.GetName(typeName)
			if typePtr != nil {
				result = &typePtr.SmiType
				return
			}
		}
	})
	return result
}

// SmiType *smiGetFirstType(SmiModule *smiModulePtr)
func (i *Instance) GetFirstType(smiModulePtr *types.SmiModule) *types.SmiType {
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

// SmiType *smiGetNextType(SmiType *smiTypePtr)
func (i *Instance) GetNextType(smiTypePtr *types.SmiType) *types.SmiType {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Next == nil {
		return nil
	}
	return &typePtr.Next.SmiType
}

// SmiType *smiGetParentType(SmiType *smiTypePtr)
func (i *Instance) GetParentType(smiTypePtr *types.SmiType) *types.SmiType {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Parent == nil {
		return nil
	}
	return &typePtr.Parent.SmiType
}

// SmiModule *smiGetTypeModule(SmiType *smiTypePtr)
func (i *Instance) GetTypeModule(smiTypePtr *types.SmiType) *types.SmiModule {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Module == nil {
		return nil
	}
	return &typePtr.Module.SmiModule
}

// int smiGetTypeLine(SmiType *smiTypePtr)
func (i *Instance) GetTypeLine(smiTypePtr *types.SmiType) int {
	if smiTypePtr == nil {
		return 0
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	return typePtr.Line
}
