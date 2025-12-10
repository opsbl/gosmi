package gosmi

import (
	"fmt"

	"github.com/opsbl/gosmi/models"
	"github.com/opsbl/gosmi/types"
)

type SmiModule struct {
	models.Module
	smiModule *types.SmiModule
	instance  *Instance
}

func (m SmiModule) GetIdentityNode() (node SmiNode, ok bool) {
	smiIdentityNode := m.instance.smiInst.GetModuleIdentityNode(m.smiModule)
	if smiIdentityNode == nil {
		return
	}
	return m.instance.createNode(smiIdentityNode), true
}

func (m SmiModule) GetImports() (imports []models.Import) {
	for smiImport := m.instance.smiInst.GetFirstImport(m.smiModule); smiImport != nil; smiImport = m.instance.smiInst.GetNextImport(smiImport) {
		_import := models.Import{
			Module: string(smiImport.Module),
			Name:   string(smiImport.Name),
		}
		imports = append(imports, _import)
	}
	return
}

func (m SmiModule) GetNode(name string) (node SmiNode, err error) {
	return m.instance.GetNode(name, m)
}

func (m SmiModule) GetNodes(kind ...types.NodeKind) (nodes []SmiNode) {
	nodeKind := types.NodeAny
	if len(kind) > 0 && kind[0] != types.NodeUnknown {
		nodeKind = kind[0]
	}
	for smiNode := m.instance.smiInst.GetFirstNode(m.smiModule, nodeKind); smiNode != nil; smiNode = m.instance.smiInst.GetNextNode(smiNode, nodeKind) {
		nodes = append(nodes, m.instance.createNode(smiNode))
	}
	return
}

func (m SmiModule) GetRevisions() (revisions []models.Revision) {
	for smiRevision := m.instance.smiInst.GetFirstRevision(m.smiModule); smiRevision != nil; smiRevision = m.instance.smiInst.GetNextRevision(smiRevision) {
		revision := models.Revision{
			Date:        smiRevision.Date,
			Description: smiRevision.Description,
		}
		revisions = append(revisions, revision)
	}
	return
}

func (m SmiModule) GetType(name string) (outType SmiType, err error) {
	return m.instance.GetType(name, m)
}

func (m SmiModule) GetTypes() (types []SmiType) {
	for smiType := m.instance.smiInst.GetFirstType(m.smiModule); smiType != nil; smiType = m.instance.smiInst.GetNextType(smiType) {
		types = append(types, m.instance.createType(smiType))
	}
	return
}

func (m SmiModule) GetRaw() (module *types.SmiModule) {
	return m.smiModule
}

func (m *SmiModule) SetRaw(smiModule *types.SmiModule) {
	m.smiModule = smiModule
}

func CreateModule(inst *Instance, smiModule *types.SmiModule) (module SmiModule) {
	return SmiModule{
		Module: models.Module{
			ContactInfo:  smiModule.ContactInfo,
			Description:  smiModule.Description,
			Language:     smiModule.Language,
			Name:         string(smiModule.Name),
			Organization: smiModule.Organization,
			Path:         smiModule.Path,
			Reference:    smiModule.Reference,
		},
		smiModule: smiModule,
		instance:  inst,
	}
}

func (i *Instance) LoadModule(modulePath string) (string, error) {
	moduleName, err := i.smiInst.LoadModuleWithError(modulePath)
	if moduleName == "" || err != nil {
		if err == nil {
			err = fmt.Errorf("Could not load module at %s", modulePath)
		}
		return "", err
	}
	return moduleName, nil
}

func (i *Instance) GetLoadedModules() (modules []SmiModule) {
	for smiModule := i.smiInst.GetFirstModule(); smiModule != nil; smiModule = i.smiInst.GetNextModule(smiModule) {
		modules = append(modules, CreateModule(i, smiModule))
	}
	return
}

func (i *Instance) IsLoaded(moduleName string) bool {
	return i.smiInst.IsLoaded(moduleName)
}

func (i *Instance) GetModule(name string) (module SmiModule, err error) {
	smiModule := i.smiInst.GetModule(name)
	if smiModule == nil {
		err = fmt.Errorf("Could not find module named %s", name)
		return
	}
	return CreateModule(i, smiModule), nil
}
