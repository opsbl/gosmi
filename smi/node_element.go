package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// GetFirstElement C -> SmiElement *smiGetFirstElement(SmiNode *smiNodePtr)
func GetFirstElement(smiNodePtr *types.SmiNode) *types.SmiElement {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.List == nil {
		return nil
	}
	return &objPtr.List.SmiElement
}

// GetNextElement C -> SmiElement *smiGetNextElement(SmiElement *smiElementPtr)
func GetNextElement(smiElementPtr *types.SmiElement) *types.SmiElement {
	if smiElementPtr == nil {
		return nil
	}
	listPtr := (*internal.List)(unsafe.Pointer(smiElementPtr))
	if listPtr.Next == nil {
		return nil
	}
	return &listPtr.Next.SmiElement
}

// GetElementNode C -> SmiNode *smiGetElementNode(SmiElement *smiElementPtr)
func GetElementNode(smiElementPtr *types.SmiElement) *types.SmiNode {
	if smiElementPtr == nil {
		return nil
	}
	listPtr := (*internal.List)(unsafe.Pointer(smiElementPtr))
	if listPtr.Ptr == nil {
		return nil
	}
	return listPtr.Ptr.(*internal.Object).GetSmiNode()
}
