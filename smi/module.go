package smi

import (
	"fmt"
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// LoadModule C -> char *smiLoadModule(const char *module)
func LoadModule(module string) string {
	checkInit()
	modulePtr, err := DefaultSmiHandle.GetModule(module)
	if err != nil {
		fmt.Println(err)
	}
	if modulePtr == nil {
		return ""
	}
	return modulePtr.Name.String()
}

// IsLoaded C -> int smiIsLoaded(const char *module)
func IsLoaded(module string) bool {
	checkInit()
	return DefaultSmiHandle.FindModuleByName(module) != nil
}

// GetModule C -> SmiModule *smiGetModule(const char *module)
func GetModule(module string) *types.SmiModule {
	if module == "" {
		return nil
	}
	modulePtr, _ := DefaultSmiHandle.GetModule(module)
	if modulePtr == nil {
		return nil
	}
	return &modulePtr.SmiModule
}

// GetFirstModule C -> SmiModule *smiGetFirstModule(void)
func GetFirstModule() *types.SmiModule {
	modulePtr := DefaultSmiHandle.GetFirstModule()
	if modulePtr == nil {
		return nil
	}
	return &modulePtr.SmiModule
}

// GetNextModule C -> SmiModule *smiGetNextModule(SmiModule *smiModulePtr)
func GetNextModule(smiModulePtr *types.SmiModule) *types.SmiModule {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.Next == nil {
		return nil
	}
	return &modulePtr.Next.SmiModule

}

// GetModuleIdentityNode C -> SmiNode *smiGetModuleIdentityNode(SmiModule *smiModulePtr)
func GetModuleIdentityNode(smiModulePtr *types.SmiModule) *types.SmiNode {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.Identity == nil {
		return nil
	}
	return modulePtr.Identity.GetSmiNode()
}
