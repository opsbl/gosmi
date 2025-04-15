package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// GetNode C -> SmiNode *smiGetNode(SmiModule *smiModulePtr, const char *name)
func GetNode(smiModulePtr *types.SmiModule, name string) *types.SmiNode {
	return smiHandle.GetNode(smiModulePtr, name)
}

// GetNodeByOID C -> SmiNode *smiGetNodeByOID(unsigned int oidlen, SmiSubid oid[])
func GetNodeByOID(oid types.Oid) *types.SmiNode {
	return smiHandle.GetNodeByOID(oid)
}

// GetFirstNode C -> SmiNode *smiGetFirstNode(SmiModule *smiModulePtr, SmiNodekind nodekind)
func GetFirstNode(smiModulePtr *types.SmiModule, nodekind types.NodeKind) *types.SmiNode {
	return smiHandle.GetFirstNode(smiModulePtr, nodekind)
}

// GetNextNode C -> SmiNode *smiGetNextNode(SmiNode *smiNodePtr, SmiNodekind nodekind)
func GetNextNode(smiNodePtr *types.SmiNode, nodekind types.NodeKind) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Module == nil || objPtr.Node == nil {
		return nil
	}
	nodePtr := objPtr.Node
	modulePtr := objPtr.Module
	for nodePtr != nil {
		if nodePtr.Children.First != nil {
			nodePtr = nodePtr.Children.First
		} else if nodePtr.Next != nil {
			nodePtr = nodePtr.Next
		} else {
			for nodePtr.Parent != nil && nodePtr.Next == nil {
				nodePtr = nodePtr.Parent
			}
			nodePtr = nodePtr.Next
			if nodePtr == nil || !nodePtr.Oid.ChildOf(modulePtr.PrefixNode.Oid) {
				return nil
			}
		}
		objPtr = internal.GetNextChildObject(nodePtr, modulePtr, nodekind)
		if objPtr != nil {
			return objPtr.GetSmiNode()
		}
	}
	return nil
}

// GetParentNode C -> SmiNode *smiGetParentNode(SmiNode *smiNodePtr)
func GetParentNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	return smiHandle.GetParentNode(smiNodePtr)
}

// GetRelatedNode C -> SmiNode *smiGetRelatedNode(SmiNode *smiNodePtr)
func GetRelatedNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Related == nil {
		return nil
	}
	return objPtr.Related.GetSmiNode()
}

// GetFirstChildNode C -> SmiNode *smiGetFirstChildNode(SmiNode *smiNodePtr)
func GetFirstChildNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Node == nil || objPtr.Node.Children.First == nil {
		return nil
	}
	nodePtr := objPtr.Node.Children.First
	objPtr = internal.FindObjectByModuleAndNode(objPtr.Module, nodePtr)
	if objPtr == nil {
		objPtr = internal.FindObjectByNode(nodePtr)
	}
	if objPtr == nil {
		return nil
	}
	return objPtr.GetSmiNode()
}

// GetNextChildNode C -> SmiNode *smiGetNextChildNode(SmiNode *smiNodePtr)
func GetNextChildNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Node == nil || objPtr.Node.Next == nil {
		return nil
	}
	nodePtr := objPtr.Node.Next
	objPtr = internal.FindObjectByModuleAndNode(objPtr.Module, nodePtr)
	if objPtr == nil {
		objPtr = internal.FindObjectByNode(nodePtr)
	}
	if objPtr == nil {
		return nil
	}
	return objPtr.GetSmiNode()
}

// GetNodeModule C -> SmiModule *smiGetNodeModule(SmiNode *smiNodePtr)
func GetNodeModule(smiNodePtr *types.SmiNode) *types.SmiModule {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Module == nil {
		return nil
	}
	return &objPtr.Module.SmiModule
}

// GetNodeType C -> SmiType *smiGetNodeType(SmiNode *smiNodePtr)
func GetNodeType(smiNodePtr *types.SmiNode) *types.SmiType {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Type == nil {
		return nil
	}
	return &objPtr.Type.SmiType
}

// GetNodeLine C -> int smiGetNodeLine(SmiNode *smiNodePtr)
func GetNodeLine(smiNodePtr *types.SmiNode) int {
	if smiNodePtr == nil {
		return 0
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	return objPtr.Line
}
