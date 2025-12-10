package smi

import (
	"strconv"
	"strings"

	"github.com/opsbl/gosmi/smi/internal"
	"github.com/opsbl/gosmi/types"
)

func (i *Instance) RenderNode(smiNodePtr *types.SmiNode, flags types.Render) string {
	if smiNodePtr == nil {
		if flags&types.RenderUnknown > 0 {
			return internal.UnknownLabel
		}
		return ""
	}
	modulePtr := i.GetNodeModule(smiNodePtr)
	if flags&types.RenderQualified == 0 || modulePtr == nil || modulePtr.Name == "" {
		return smiNodePtr.Name.String()
	}
	return modulePtr.Name.String() + "::" + smiNodePtr.Name.String()
}

func (i *Instance) RenderOID(oid types.Oid, flags types.Render) string {
	if len(oid) == 0 {
		if flags&types.RenderUnknown > 0 {
			return internal.UnknownLabel
		}
		return ""
	}
	var idx int
	var b strings.Builder
	if flags&(types.RenderName|types.RenderQualified) > 0 {
		nodePtr := i.GetNodeByOID(oid)
		if nodePtr != nil && nodePtr.Name != "" {
			idx = nodePtr.OidLen
			b.WriteString(i.RenderNode(nodePtr, flags))
		}
	}
	for ; idx < len(oid); idx++ {
		if b.Len() > 0 {
			b.WriteRune('.')
		}
		b.WriteString(strconv.FormatUint(uint64(oid[idx]), 10))
	}
	return b.String()
}
