package internal

import (
	"errors"

	"github.com/opsbl/gosmi/parser"
	"github.com/opsbl/gosmi/types"
)

const WellKnownModuleName types.SmiIdentifier = "<well-known>"

var (
	WellKnownIdCcitt         types.SmiSubId = 0
	WellKnownIdIso           types.SmiSubId = 1
	WellKnownIdJointIsoCcitt types.SmiSubId = 2
)

type Handle struct {
	Name                 string
	Modules              ModuleMap
	RootNode             *Node
	TypeBits             *Type
	TypeEnum             *Type
	TypeInteger32        *Type
	TypeInteger64        *Type
	TypeObjectIdentifier *Type
	TypeOctetString      *Type
	TypeUnsigned32       *Type
	TypeUnsigned64       *Type
	Flags                Flags
	Paths                []NamedFS
	Cache                string
	CacheProg            string
	ErrorLevel           int
	ErrorHandler         types.SmiErrorHandler

	loading map[types.SmiIdentifier]bool
	loadErr error
}

func NewHandle(name string) (*Handle, error) {
	handlePtr := &Handle{Name: name}
	ok := handlePtr.initData()
	if !ok {
		return nil, errors.New("failed to initialize handle")
	}
	return handlePtr, nil
}

func SetErrorHandler(smiErrorHandler types.SmiErrorHandler) {
	// Deprecated global; no-op
}

func SetSeverity(pattern string, severity int) {}

func SetErrorLevel(level int) {}

func GetFlags() int { return 0 }

func SetFlags(userflags int)                                    {}
func (h *Handle) SetErrorHandler(handler types.SmiErrorHandler) { h.ErrorHandler = handler }
func (h *Handle) SetSeverity(pattern string, severity int)      {}
func (h *Handle) SetErrorLevel(level int)                       {}
func (h *Handle) GetFlags() int                                 { return int(h.Flags) }
func (h *Handle) SetFlags(userflags int)                        { h.Flags = Flags(userflags) }

func oidFromSubId(subId types.SmiSubId) parser.Oid {
	return parser.Oid{
		SubIdentifiers: []parser.SubIdentifier{
			{Number: &subId},
		},
	}
}

func createBaseType(module *Module, baseType types.BaseType) *Type {
	return &Type{
		SmiType: types.SmiType{
			Name:     types.SmiIdentifier(baseType.String()),
			BaseType: baseType,
			Decl:     types.DeclImplicitType,
		},
		Module: module,
	}
}

func (h *Handle) GetFirstModule() *Module {
	if h == nil {
		return nil
	}
	return h.Modules.First
}

func (h *Handle) Root() *Node {
	if h == nil {
		return nil
	}
	return h.RootNode
}

func (h *Handle) initData() bool {
	h.RootNode = &Node{Flags: FlagRoot, Oid: types.Oid{}}

	wellKnownModule := &Module{
		SmiModule: types.SmiModule{
			Name: WellKnownModuleName,
		},
		PrefixNode: h.RootNode,
		Handle:     h,
	}

	// Create ccitt well-known node
	ccitt := &Object{
		SmiNode: types.SmiNode{
			Name:     "ccitt",
			Decl:     types.DeclImplObject,
			NodeKind: types.NodeNode,
		},
		Module: wellKnownModule,
	}
	wellKnownModule.Objects.AddWithOid(ccitt, oidFromSubId(WellKnownIdCcitt))

	// Create iso well-known node
	iso := &Object{
		SmiNode: types.SmiNode{
			Name:     "iso",
			Decl:     types.DeclImplObject,
			NodeKind: types.NodeNode,
		},
		Module: wellKnownModule,
	}
	wellKnownModule.Objects.AddWithOid(iso, oidFromSubId(WellKnownIdIso))

	// Create joint-iso-ccitt well-known node
	jointIsoCcitt := &Object{
		SmiNode: types.SmiNode{
			Name:     "joint-iso-ccitt",
			Decl:     types.DeclImplObject,
			NodeKind: types.NodeNode,
		},
		Module: wellKnownModule,
	}
	wellKnownModule.Objects.AddWithOid(jointIsoCcitt, oidFromSubId(WellKnownIdJointIsoCcitt))

	h.Modules.Add(wellKnownModule)

	h.TypeBits = createBaseType(wellKnownModule, types.BaseTypeBits)
	h.TypeEnum = createBaseType(wellKnownModule, types.BaseTypeEnum)
	h.TypeInteger32 = createBaseType(wellKnownModule, types.BaseTypeInteger32)
	h.TypeInteger64 = createBaseType(wellKnownModule, types.BaseTypeInteger64)
	h.TypeObjectIdentifier = createBaseType(wellKnownModule, types.BaseTypeObjectIdentifier)
	h.TypeOctetString = createBaseType(wellKnownModule, types.BaseTypeOctetString)
	h.TypeUnsigned32 = createBaseType(wellKnownModule, types.BaseTypeUnsigned32)
	h.TypeUnsigned64 = createBaseType(wellKnownModule, types.BaseTypeUnsigned64)

	return true
}

func (h *Handle) freeData() {
}
