package smi

import (
	"unsafe"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

// SmiNode *smiGetNode(SmiModule *smiModulePtr, const char *name)
func (i *Instance) GetNode(smiModulePtr *types.SmiModule, name string) *types.SmiNode {
	if name == "" {
		return nil
	}
	var result *types.SmiNode
	i.withHandle(func(h *internal.Handle) {
		var modulePtr *internal.Module
		if smiModulePtr != nil {
			modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
			objPtr := modulePtr.Objects.GetName(name)
			if objPtr != nil {
				result = objPtr.GetSmiNode()
			}
			return
		}
		for modulePtr = h.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
			objPtr := modulePtr.Objects.GetName(name)
			if objPtr != nil {
				result = objPtr.GetSmiNode()
				return
			}
		}
	})
	return result
}

// SmiNode *smiGetNodeByOID(unsigned int oidlen, SmiSubid oid[])
func (i *Instance) GetNodeByOID(oid types.Oid) *types.SmiNode {
	if len(oid) == 0 {
		return nil
	}
	var result *types.SmiNode
	i.withHandle(func(h *internal.Handle) {
		if h.Root() == nil {
			return
		}
		var parentPtr, nodePtr *internal.Node = nil, h.Root()
		for idx := 0; idx < len(oid) && nodePtr != nil; idx++ {
			parentPtr, nodePtr = nodePtr, nodePtr.Children.Get(oid[idx])
		}
		if nodePtr == nil {
			nodePtr = parentPtr
		}
		if nodePtr != nil && nodePtr.FirstObject != nil {
			result = nodePtr.FirstObject.GetSmiNode()
		}
	})
	return result
}

// SmiNode *smiGetFirstNode(SmiModule *smiModulePtr, SmiNodekind nodekind)
func (i *Instance) GetFirstNode(smiModulePtr *types.SmiModule, nodekind types.NodeKind) *types.SmiNode {
	if smiModulePtr == nil {
		return nil
	}
	var result *types.SmiNode
	i.withHandle(func(h *internal.Handle) {
		var (
			modulePtr *internal.Module
			nodePtr   *internal.Node
			objPtr    *internal.Object
		)
		modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
		if modulePtr.PrefixNode != nil {
			nodePtr = modulePtr.PrefixNode
		} else if h.Root() != nil {
			nodePtr = h.Root().Children.First
		}
		for nodePtr != nil {
			objPtr = h.GetNextChildObject(nodePtr, modulePtr, nodekind)
			if objPtr != nil {
				result = objPtr.GetSmiNode()
				return
			}
			if nodePtr.Children.First != nil {
				nodePtr = nodePtr.Children.First
			} else if nodePtr.Next != nil {
				nodePtr = nodePtr.Next
			} else {
				if nodePtr.Parent == nil {
					return
				}
				for nodePtr.Parent != nil && nodePtr.Next == nil {
					nodePtr = nodePtr.Parent
				}
				nodePtr = nodePtr.Next
			}
		}
	})
	return result
}

// SmiNode *smiGetNextNode(SmiNode *smiNodePtr, SmiNodekind nodekind)
func (i *Instance) GetNextNode(smiNodePtr *types.SmiNode, nodekind types.NodeKind) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	var result *types.SmiNode
	i.withHandle(func(h *internal.Handle) {
		objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
		if objPtr.Module == nil || objPtr.Node == nil {
			return
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
					return
				}
			}
			objPtr = h.GetNextChildObject(nodePtr, modulePtr, nodekind)
			if objPtr != nil {
				result = objPtr.GetSmiNode()
				return
			}
		}
	})
	return result
}

// SmiNode *smiGetParentNode(SmiNode *smiNodePtr)
func (i *Instance) GetParentNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	var result *types.SmiNode
	i.withHandle(func(h *internal.Handle) {
		objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
		if objPtr.Node == nil || objPtr.Node.Parent == nil || objPtr.Node.Flags.Has(internal.FlagRoot) {
			return
		}
		var parentPtr *internal.Object
		if objPtr.Module != nil {
			parentPtr = h.FindObjectByModuleAndNode(objPtr.Module, objPtr.Node.Parent)
			if parentPtr != nil {
				importPtr := objPtr.Module.Imports.Get(parentPtr.Name)
				if importPtr != nil {
					parentPtr = h.FindObjectByModuleNameAndNode(string(importPtr.Module), objPtr.Node.Parent)
				} else {
					parentPtr = nil
				}
			}
		}
		if parentPtr == nil {
			parentPtr = h.FindObjectByNode(objPtr.Node.Parent)
		}
		if parentPtr != nil {
			result = parentPtr.GetSmiNode()
		}
	})
	return result
}

// SmiNode *smiGetRelatedNode(SmiNode *smiNodePtr)
func (i *Instance) GetRelatedNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Related == nil {
		return nil
	}
	return objPtr.Related.GetSmiNode()
}

// SmiNode *smiGetFirstChildNode(SmiNode *smiNodePtr)
func (i *Instance) GetFirstChildNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	var result *types.SmiNode
	i.withHandle(func(h *internal.Handle) {
		objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
		if objPtr.Node == nil || objPtr.Node.Children.First == nil {
			return
		}
		nodePtr := objPtr.Node.Children.First
		objPtr = h.FindObjectByModuleAndNode(objPtr.Module, nodePtr)
		if objPtr == nil {
			objPtr = h.FindObjectByNode(nodePtr)
		}
		if objPtr != nil {
			result = objPtr.GetSmiNode()
		}
	})
	return result
}

// SmiNode *smiGetNextChildNode(SmiNode *smiNodePtr)
func (i *Instance) GetNextChildNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	var result *types.SmiNode
	i.withHandle(func(h *internal.Handle) {
		objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
		if objPtr.Node == nil || objPtr.Node.Next == nil {
			return
		}
		nodePtr := objPtr.Node.Next
		objPtr = h.FindObjectByModuleAndNode(objPtr.Module, nodePtr)
		if objPtr == nil {
			objPtr = h.FindObjectByNode(nodePtr)
		}
		if objPtr != nil {
			result = objPtr.GetSmiNode()
		}
	})
	return result
}

// SmiModule *smiGetNodeModule(SmiNode *smiNodePtr)
func (i *Instance) GetNodeModule(smiNodePtr *types.SmiNode) *types.SmiModule {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Module == nil {
		return nil
	}
	return &objPtr.Module.SmiModule
}

// SmiType *smiGetNodeType(SmiNode *smiNodePtr)
func (i *Instance) GetNodeType(smiNodePtr *types.SmiNode) *types.SmiType {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Type == nil {
		return nil
	}
	return &objPtr.Type.SmiType
}

// int smiGetNodeLine(SmiNode *smiNodePtr)
func (i *Instance) GetNodeLine(smiNodePtr *types.SmiNode) int {
	if smiNodePtr == nil {
		return 0
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	return objPtr.Line
}
