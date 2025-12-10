package types

//go:generate enumer -type=BaseType -trimprefix=BaseType -json -yaml -output=basetype_string.go
type BaseType int

const (
	BaseTypeUnknown BaseType = iota
	BaseTypeInteger32
	BaseTypeOctetString
	BaseTypeObjectIdentifier
	BaseTypeUnsigned32
	BaseTypeInteger64
	BaseTypeUnsigned64
	BaseTypeFloat32
	BaseTypeFloat64
	BaseTypeFloat128
	BaseTypeEnum
	BaseTypeBits
	BaseTypePointer
)
