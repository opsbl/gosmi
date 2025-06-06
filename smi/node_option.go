package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// GetFirstOption C -> SmiOption *smiGetFirstOption(SmiNode *smiComplianceNodePtr)
func GetFirstOption(smiComplianceNodePtr *types.SmiNode) *types.SmiOption {
	if smiComplianceNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiComplianceNodePtr))
	if objPtr.NodeKind != types.NodeCompliance || objPtr.OptionList == nil || objPtr.OptionList.Ptr == nil {
		return nil
	}
	return &objPtr.OptionList.Ptr.(*internal.Option).SmiOption
}

// GetNextOption C -> SmiOption *smiGetNextOption(SmiOption *smiOptionPtr)
func GetNextOption(smiOptionPtr *types.SmiOption) *types.SmiOption {
	if smiOptionPtr == nil {
		return nil
	}
	optPtr := (*internal.Option)(unsafe.Pointer(smiOptionPtr))
	if optPtr.List == nil || optPtr.List.Next == nil || optPtr.List.Next.Ptr == nil {
		return nil
	}
	return &optPtr.List.Next.Ptr.(*internal.Option).SmiOption
}

// GetOptionNode C -> SmiNode *smiGetOptionNode(SmiOption *smiOptionPtr)
func GetOptionNode(smiOptionPtr *types.SmiOption) *types.SmiNode {
	if smiOptionPtr == nil {
		return nil
	}
	optionPtr := (*internal.Option)(unsafe.Pointer(smiOptionPtr))
	if optionPtr.Object == nil {
		return nil
	}
	return optionPtr.Object.GetSmiNode()
}

// GetOptionLine C ->  int smiGetOptionLine(SmiOption *smiOptionPtr)
func GetOptionLine(smiOptionPtr *types.SmiOption) int {
	if smiOptionPtr == nil {
		return 0
	}
	optionPtr := (*internal.Option)(unsafe.Pointer(smiOptionPtr))
	return optionPtr.Line
}
