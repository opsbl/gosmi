package smi

import (
	"fmt"
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// char *smiLoadModule(const char *module)
func (i *Instance) LoadModuleWithError(module string) (string, error) {
	if module == "" {
		return "", nil
	}
	var moduleName string
	var loadErr error
	i.withHandle(func(h *internal.Handle) {
		modulePtr, err := h.GetModule(module)
		if err != nil {
			loadErr = err
			return
		}
		if modulePtr != nil {
			moduleName = modulePtr.Name.String()
		}
	})
	return moduleName, loadErr
}

func (i *Instance) LoadModule(module string) string {
	moduleName, err := i.LoadModuleWithError(module)
	if err != nil {
		fmt.Println(err)
	}
	return moduleName
}

// int smiIsLoaded(const char *module)
func (i *Instance) IsLoaded(module string) bool {
	if module == "" {
		return false
	}
	var loaded bool
	i.withHandle(func(h *internal.Handle) {
		loaded = h.FindModuleByName(module) != nil
	})
	return loaded
}

// SmiModule *smiGetModule(const char *module)
func (i *Instance) GetModule(module string) *types.SmiModule {
	if module == "" {
		return nil
	}
	var out *types.SmiModule
	i.withHandle(func(h *internal.Handle) {
		modulePtr, _ := h.GetModule(module)
		if modulePtr != nil {
			out = &modulePtr.SmiModule
		}
	})
	return out
}

// SmiModule *smiGetFirstModule(void)
func (i *Instance) GetFirstModule() *types.SmiModule {
	var modulePtr *internal.Module
	i.withHandle(func(h *internal.Handle) {
		modulePtr = h.GetFirstModule()
	})
	if modulePtr == nil {
		return nil
	}
	return &modulePtr.SmiModule
}

// SmiModule *smiGetNextModule(SmiModule *smiModulePtr)
func (i *Instance) GetNextModule(smiModulePtr *types.SmiModule) *types.SmiModule {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.Next == nil {
		return nil
	}
	return &modulePtr.Next.SmiModule

}

// SmiNode *smiGetModuleIdentityNode(SmiModule *smiModulePtr)
func (i *Instance) GetModuleIdentityNode(smiModulePtr *types.SmiModule) *types.SmiNode {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.Identity == nil {
		return nil
	}
	return modulePtr.Identity.GetSmiNode()
}
