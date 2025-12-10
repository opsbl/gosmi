package types

//go:generate enumer -type=IndexKind -trimprefix=IndexKind -json -yaml -output=indexkind_string.go

type IndexKind int

const (
	IndexUnknown IndexKind = iota
	IndexIndex
	IndexAugment
	IndexReorder
	IndexSparse
	IndexExpand
)
