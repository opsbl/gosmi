package internal

import (
	"github.com/opsbl/gosmi/types"
	"unsafe"
)

type Node struct {
	SubId       types.SmiSubId
	Flags       Flags
	OidLen      int
	Oid         types.Oid
	Parent      *Node
	Prev        *Node
	Next        *Node
	Children    NodeChildMap
	FirstObject *Object
	LastObject  *Object
}

func (x *Node) AddObject(obj *Object) {
	obj.Node = x
	obj.PrevSameNode = x.LastObject
	if x.LastObject == nil {
		x.FirstObject = obj
	} else {
		x.LastObject.NextSameNode = obj
	}
	x.LastObject = obj
}

func (x *Node) IsRoot() bool {
	return x != nil && x.Flags.Has(FlagRoot)
}

type NodeChildMap struct {
	First *Node

	last *Node
	m    map[types.SmiSubId]*Node
}

func (x *NodeChildMap) Add(n *Node) {
	existing := x.Get(n.SubId)
	if existing != nil {
		for obj := n.FirstObject; obj != nil; obj = obj.NextSameNode {
			existing.AddObject(obj)
		}
		return
	}
	if n.Parent != nil && n.Parent.Oid != nil {
		n.Oid = types.NewOid(n.Parent.Oid, n.SubId)
		n.OidLen = n.Parent.OidLen + 1
	}
	if x.last == nil {
		x.First = n
		x.last = n
	} else {
		c := x.First
		for c != nil && c.SubId < n.SubId {
			c = c.Next
		}
		if c == nil {
			n.Prev = x.last
			x.last.Next = n
			x.last = n
		} else {
			n.Prev = c.Prev
			n.Next = c
			if c.Prev == nil {
				x.First = n
			} else {
				c.Prev.Next = n
			}
			c.Prev = n
		}
	}

	if x.m == nil {
		x.m = make(map[types.SmiSubId]*Node)
	}
	x.m[n.SubId] = n
}

func (x *NodeChildMap) Get(id types.SmiSubId) *Node {
	if x.m == nil {
		return nil
	}
	return x.m[id]
}

func (h *Handle) FindNodeByOid(oidlen int, oid types.Oid) *Node {
	nodePtr := h.RootNode
	for i := 0; i < oidlen && nodePtr != nil; i++ {
		nodePtr = nodePtr.Children.Get(oid[i])
	}
	return nodePtr
}

// GetNodeByOID C -> SmiNode *smiGetNodeByOID(unsigned int oidlen, SmiSubid oid[])
func (h *Handle) GetNodeByOID(oid types.Oid) *types.SmiNode {
	if len(oid) == 0 || h.Root() == nil {
		return nil
	}
	var parentPtr, nodePtr *Node = nil, h.Root()
	for i := 0; i < len(oid) && nodePtr != nil; i++ {
		parentPtr, nodePtr = nodePtr, nodePtr.Children.Get(oid[i])
	}
	if nodePtr == nil {
		nodePtr = parentPtr
	}
	if nodePtr == nil || nodePtr.FirstObject == nil {
		return nil
	}
	return nodePtr.FirstObject.GetSmiNode()
}

// GetParentNode C -> SmiNode *smiGetParentNode(SmiNode *smiNodePtr)
func (h *Handle) GetParentNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Node == nil || objPtr.Node.Parent == nil || objPtr.Node.Flags.Has(FlagRoot) {
		return nil
	}
	var parentPtr *Object
	if objPtr.Module != nil {
		parentPtr = FindObjectByModuleAndNode(objPtr.Module, objPtr.Node.Parent)
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
		parentPtr = FindObjectByNode(objPtr.Node.Parent)
	}
	if parentPtr == nil {
		return nil
	}
	return parentPtr.GetSmiNode()
}

// GetFirstNode C -> SmiNode *smiGetFirstNode(SmiModule *smiModulePtr, SmiNodekind nodekind)
func (h *Handle) GetFirstNode(smiModulePtr *types.SmiModule, nodekind types.NodeKind) *types.SmiNode {
	if smiModulePtr == nil {
		return nil
	}
	var (
		modulePtr *Module
		nodePtr   *Node
		objPtr    *Object
	)
	modulePtr = (*Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.PrefixNode != nil {
		nodePtr = modulePtr.PrefixNode
	} else if h.Root() != nil {
		nodePtr = h.Root().Children.First
	}
	for nodePtr != nil {
		objPtr = GetNextChildObject(nodePtr, modulePtr, nodekind)
		if objPtr != nil {
			return objPtr.GetSmiNode()
		}
		if nodePtr.Children.First != nil {
			nodePtr = nodePtr.Children.First
		} else if nodePtr.Next != nil {
			nodePtr = nodePtr.Next
		} else {
			if nodePtr.Parent == nil {
				return nil
			}
			for nodePtr.Parent != nil && nodePtr.Next == nil {
				nodePtr = nodePtr.Parent
			}
			nodePtr = nodePtr.Next
		}
	}
	return nil
}

// GetNode C -> SmiNode *smiGetNode(SmiModule *smiModulePtr, const char *name)
func (h *Handle) GetNode(smiModulePtr *types.SmiModule, name string) *types.SmiNode {
	if name == "" {
		return nil
	}
	var modulePtr *Module
	if smiModulePtr != nil {
		modulePtr = (*Module)(unsafe.Pointer(smiModulePtr))
		objPtr := modulePtr.Objects.GetName(name)
		if objPtr == nil {
			return nil
		}
		return objPtr.GetSmiNode()
	}
	for modulePtr = h.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
		objPtr := modulePtr.Objects.GetName(name)
		if objPtr != nil {
			return objPtr.GetSmiNode()
		}
	}
	return nil
}

/*func createNodes(oid types.Oid) *Node {
	var nodePtr *Node
	parentNodePtr := smiHandle.RootNode
	for i := 0; i < len(oid); i++ {
		nodePtr = parentNodePtr.Children.Get(oid[i])
		if nodePtr != nil {
			continue
		}
		nodePtr = &Node{
			SubId:  oid[i],
			Oid:    oid[:i+1],
			Parent: parentNodePtr,
		}
		parentNodePtr.Children.Add(nodePtr)
		parentNodePtr = nodePtr
	}
	return parentNodePtr
}*/
