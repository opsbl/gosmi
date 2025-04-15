package internal

import (
	"github.com/opsbl/gosmi/parser"
	"github.com/opsbl/gosmi/types"
	"os"
	"strings"
)

const WellKnownModuleName types.SmiIdentifier = "<well-known>"

var (
	WellKnownIdCcitt         types.SmiSubId = 0
	WellKnownIdIso           types.SmiSubId = 1
	WellKnownIdJointIsoCcitt types.SmiSubId = 2
)

type Handle struct {
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
}

func (h *Handle) SetErrorHandler(smiErrorHandler types.SmiErrorHandler) *Handle {
	h.ErrorHandler = smiErrorHandler
	return h
}

func (h *Handle) SetSeverity(pattern string, severity int) *Handle {
	return h
}

func (h *Handle) SetErrorLevel(level int) *Handle {
	h.ErrorLevel = level
	return h
}

func (h *Handle) SetFlags(flags Flags) *Handle {
	h.Flags = flags
	return h
}

func (h *Handle) GetFlags() Flags {
	return h.Flags
}

func (h *Handle) GetFirstModule() *Module {
	return h.Modules.First
}

func (h *Handle) Root() *Node {
	return h.RootNode
}

func (h *Handle) initData() bool {
	h.RootNode = &Node{
		Flags: FlagRoot,
		Oid:   types.Oid{},
	}
	wellKnownModule := &Module{
		SmiModule: types.SmiModule{
			Name: WellKnownModuleName,
		},
		Handle:     h,
		PrefixNode: h.RootNode,
		Objects:    ObjectMap{Handle: h},
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

func (h *Handle) GetPath() string {
	names := make([]string, len(h.Paths))
	for i, fs := range h.Paths {
		names[i] = fs.Name
	}
	return strings.Join(names, string(os.PathListSeparator))
}

func (h *Handle) SetPath(path ...string) {
	pathLen := len(path)
	if pathLen == 0 {
		return
	}
	if path[0] == "" {
		appendPath(h, path[1:]...)
	} else if path[pathLen-1] == "" {
		prependPath(h, path[:pathLen-1]...)
	} else {
		h.Paths = make([]NamedFS, 0, pathLen)
		for _, p := range path {
			if p, err := expandPath(p); err == nil {
				h.Paths = append(h.Paths, newPathFS(p))
			}
		}
	}
}

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

func NewHandle() *Handle {
	handle := new(Handle)
	handle.initData()
	return handle
}
