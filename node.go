package gosmi

import (
	"fmt"

	"github.com/opsbl/gosmi/models"
	"github.com/opsbl/gosmi/smi"
	"github.com/opsbl/gosmi/types"
)

type SmiNode struct {
	models.Node
	smiNode *types.SmiNode
	SmiType *SmiType
	handle  *smi.Handle
}

func (n *SmiNode) GetModule() (module SmiModule) {
	smiModule := smi.GetNodeModule(n.smiNode)
	return CreateModule(smiModule, n.handle)
}

func (n *SmiNode) GetSubtree() (nodes []SmiNode) {
	first := true
	smiNode := n.smiNode
	for oidlen := n.OidLen; smiNode != nil && (first || int(smiNode.OidLen) > oidlen); smiNode = smi.GetNextNode(smiNode, types.NodeAny) {
		node := CreateNode(smiNode, n.handle)
		nodes = append(nodes, node)
		first = false
	}
	return
}

func (n *SmiNode) Render(flags types.Render) string {
	return smi.RenderNode(n.smiNode, flags)
}

func (n *SmiNode) RenderNumeric() string {
	return smi.RenderOID(n.smiNode.Oid, types.RenderNumeric)
}

func (n *SmiNode) RenderQualified() string {
	return n.Render(types.RenderQualified)
}

func (n *SmiNode) GetRaw() (node *types.SmiNode) {
	return n.smiNode
}

func (n *SmiNode) SetRaw(smiNode *types.SmiNode) {
	n.smiNode = smiNode
}

func CreateNode(smiNode *types.SmiNode, handle *smi.Handle) SmiNode {
	node := SmiNode{
		Node: models.Node{
			Access:      smiNode.Access,
			Decl:        smiNode.Decl,
			Description: smiNode.Description,
			Kind:        smiNode.NodeKind,
			Name:        string(smiNode.Name),
			OidLen:      smiNode.OidLen,
			Oid:         smiNode.Oid,
			Status:      smiNode.Status,
		},
		smiNode: smiNode,
		SmiType: CreateTypeFromNode(smiNode, handle),
		handle:  handle,
	}
	if node.SmiType != nil {
		node.Type = &node.SmiType.Type
	}
	return node
}

func GetNode(name string, module ...SmiModule) (node SmiNode, err error) {
	var smiModule *types.SmiModule
	if len(module) > 0 {
		smiModule = module[0].GetRaw()
	}

	smiNode := smi.GetNode(smiModule, name)
	if smiNode == nil {
		if len(module) > 0 {
			err = fmt.Errorf("could not find node named %s in module %s", name, module[0].Name)
		} else {
			err = fmt.Errorf("could not find node named %s", name)
		}
		return
	}
	return CreateNode(smiNode, smi.DefaultSmiHandle), nil
}

func GetNodeByOID(oid types.Oid) (node SmiNode, err error) {
	smiNode := smi.GetNodeByOID(oid)
	if smiNode == nil {
		err = fmt.Errorf("could not find node for OID %s", oid)
		return
	}
	return CreateNode(smiNode, smi.DefaultSmiHandle), nil
}
