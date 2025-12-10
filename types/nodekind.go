package types

//go:generate enumer -type=NodeKind -trimprefix=NodeKind -json -yaml -output=nodekind_string.go

type NodeKind int

const (
	NodeUnknown NodeKind = iota
	NodeNode    NodeKind = 1 << (iota - 1)
	NodeScalar
	NodeTable
	NodeRow
	NodeColumn
	NodeNotification
	NodeGroup
	NodeCompliance
	NodeCapabilities
	NodeAny NodeKind = 0xffff
)
