package gosmi

import (
	"github.com/opsbl/gosmi/smi"
	"github.com/opsbl/gosmi/types"
)

type Table struct {
	SmiNode
	Columns     map[string]SmiNode
	ColumnOrder []string
	Implied     bool
	Index       []SmiNode
}

func (n *SmiNode) AsTable() Table {
	columns, columnOrder := n.GetColumns()
	return Table{
		SmiNode:     *n,
		Columns:     columns,
		ColumnOrder: columnOrder,
		Implied:     n.GetImplied(),
		Index:       n.GetIndex(),
	}
}

func (n *SmiNode) getRow() (row *types.SmiNode) {
	switch n.Kind {
	case types.NodeRow:
		row = n.GetRaw()
	case types.NodeTable:
		row = smi.GetFirstChildNode(n.smiNode)
		if row == nil {
			return
		}
	default:
		return
	}

	if row.NodeKind != types.NodeRow {
		// TODO: error
		return nil
	}

	return
}

func (n *SmiNode) GetRow() (row SmiNode) {
	smiRow := n.getRow()
	if smiRow == nil {
		return
	}
	return CreateNode(smiRow, n.handle)
}

func (n *SmiNode) GetColumns() (columns map[string]SmiNode, columnOrder []string) {
	row := n.getRow()
	if row == nil {
		return
	}

	columns = make(map[string]SmiNode)
	columnOrder = make([]string, 0, 2)

	for smiColumn := smi.GetFirstChildNode(row); smiColumn != nil; smiColumn = smi.GetNextChildNode(smiColumn) {
		if smiColumn.NodeKind != types.NodeColumn {
			// TODO: error
			return
		}
		column := CreateNode(smiColumn, n.handle)
		columns[column.Name] = column
		columnOrder = append(columnOrder, column.Name)
	}
	return
}

func (n *SmiNode) GetImplied() bool {
	row := n.getRow()
	if row == nil {
		return false
	}

	return row.Implied
}

func (n *SmiNode) GetAugment() (row SmiNode) {
	smiRow := n.getRow()
	if smiRow == nil {
		return
	}

	if smiRow.IndexKind != types.IndexAugment {
		return
	}

	smiRow = smi.GetRelatedNode(smiRow)
	if smiRow == nil {
		return
	}

	if smiRow.NodeKind != types.NodeRow {
		// TODO: error
		return
	}

	return CreateNode(smiRow, n.handle)
}

func (n *SmiNode) GetIndex() (index []SmiNode) {
	row := n.getRow()
	if row == nil {
		return
	}

	if row.IndexKind == types.IndexAugment {
		row = smi.GetRelatedNode(row)
		if row == nil {
			return
		}

		if row.NodeKind != types.NodeRow {
			// TODO: error
			return
		}
	} else if row.IndexKind != types.IndexIndex {
		// TODO: unsupported
		return
	}

	for smiElement := smi.GetFirstElement(row); smiElement != nil; smiElement = smi.GetNextElement(smiElement) {
		smiColumn := smi.GetElementNode(smiElement)
		if smiColumn == nil {
			// TODO: error
			return
		}
		if smiColumn.NodeKind != types.NodeColumn {
			// TODO: error
			return
		}
		index = append(index, CreateNode(smiColumn, n.handle))
	}
	return
}
