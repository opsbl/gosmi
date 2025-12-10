package gosmi

import (
	"fmt"

	"github.com/opsbl/gosmi/models"
	"github.com/opsbl/gosmi/types"
)

type SmiNode struct {
	models.Node
	smiNode  *types.SmiNode
	SmiType  *SmiType
	instance *Instance
}

func (n SmiNode) GetModule() (module SmiModule) {
	smiModule := n.instance.smiInst.GetNodeModule(n.smiNode)
	return CreateModule(n.instance, smiModule)
}

func (n SmiNode) GetSubtree() (nodes []SmiNode) {
	first := true
	smiNode := n.smiNode
	for oidlen := n.OidLen; smiNode != nil && (first || int(smiNode.OidLen) > oidlen); smiNode = n.instance.smiInst.GetNextNode(smiNode, types.NodeAny) {
		node := n.instance.createNode(smiNode)
		nodes = append(nodes, node)
		first = false
	}
	return
}

func (n SmiNode) Render(flags types.Render) string {
	return n.instance.smiInst.RenderNode(n.smiNode, flags)
}

func (n SmiNode) RenderNumeric() string {
	return n.instance.smiInst.RenderOID(n.smiNode.Oid, types.RenderNumeric)
}

func (n SmiNode) RenderQualified() string {
	return n.Render(types.RenderQualified)
}

func (n SmiNode) GetRaw() (node *types.SmiNode) {
	return n.smiNode
}

func (n *SmiNode) SetRaw(smiNode *types.SmiNode) {
	n.smiNode = smiNode
}

func CreateNode(inst *Instance, smiNode *types.SmiNode) SmiNode {
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
		smiNode:  smiNode,
		SmiType:  CreateTypeFromNode(inst, smiNode),
		instance: inst,
	}
	if node.SmiType != nil {
		node.Type = &node.SmiType.Type
	}
	return node
}

func (i *Instance) GetNode(name string, module ...SmiModule) (node SmiNode, err error) {
	var smiModule *types.SmiModule
	if len(module) > 0 {
		smiModule = module[0].GetRaw()
	}

	smiNode := i.smiInst.GetNode(smiModule, name)
	if smiNode == nil {
		if len(module) > 0 {
			err = fmt.Errorf("Could not find node named %s in module %s", name, module[0].Name)
		} else {
			err = fmt.Errorf("Could not find node named %s", name)
		}
		return
	}
	return i.createNode(smiNode), nil
}

func (i *Instance) GetNodeByOID(oid types.Oid) (node SmiNode, err error) {
	smiNode := i.smiInst.GetNodeByOID(oid)
	if smiNode == nil {
		err = fmt.Errorf("Could not find node for OID %s", oid)
		return
	}
	return i.createNode(smiNode), nil
}
