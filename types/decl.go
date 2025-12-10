package types

//go:generate enumer -type=Decl -trimprefix=Decl -json -yaml -output=decl_string.go

type Decl int

const (
	DeclUnknown Decl = iota
	DeclImplicitType
	DeclTypeAssignment
	_
	DeclImplSequenceOf
	DeclValueAssignment
	DeclObjectType
	DeclObjectIdentity
	DeclModuleIdentity
	DeclNotificationType
	DeclTrapType
	DeclObjectGroup
	DeclNotificationGroup
	DeclModuleCompliance
	DeclAgentCapabilities
	DeclTextualConvention
	DeclMacro
	DeclComplGroup
	DeclComplObject
	DeclImplObject
	DeclModule Decl = iota + 13
	DeclExtension
	DeclTypedef
	DeclNode
	DeclScalar
	DeclTable
	DeclRow
	DeclColumn
	DeclNotification
	DeclGroup
	DeclCompliance
	DeclIdentity
	DeclClass
	DeclAttribute
	DeclEvent
)
