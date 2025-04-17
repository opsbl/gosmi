package gosmi

import (
	"fmt"
	"github.com/opsbl/gosmi/smi"
	"github.com/opsbl/gosmi/types"
)

type HandleHelper struct {
	*smi.Handle
}

func (x *HandleHelper) LoadModule(module string) (string, error) {
	modulePtr, err := x.Handle.GetModule(module)
	if err != nil {
		return "", err
	}
	return modulePtr.Name.String(), nil
}

func (x *HandleHelper) GetLoadedModules() (modules []SmiModule) {
	firstModule := x.Handle.GetFirstModule()
	for smiModule := &firstModule.SmiModule; smiModule != nil; smiModule = smi.GetNextModule(smiModule) {
		modules = append(modules, CreateModule(smiModule, x.Handle))
	}
	return
}

func (x *HandleHelper) IsLoaded(moduleName string) bool {
	return x.Handle.FindModuleByName(moduleName) != nil
}

func (x *HandleHelper) GetModule(name string) (SmiModule, error) {
	var module SmiModule
	smiModule, err := x.Handle.GetModule(name)
	if err != nil {
		return module, err
	}
	return CreateModule(&smiModule.SmiModule, x.Handle), nil
}

func (x *HandleHelper) GetNode(name string, module ...SmiModule) (node SmiNode, err error) {
	var smiModule *types.SmiModule
	if len(module) > 0 {
		smiModule = module[0].GetRaw()
	}
	smiNode := x.Handle.GetNode(smiModule, name)
	if smiNode == nil {
		if len(module) > 0 {
			err = fmt.Errorf("could not find node named %s in module %s", name, module[0].Name)
		} else {
			err = fmt.Errorf("could not find node named %s", name)
		}
		return
	}
	return CreateNode(smiNode, x.Handle), nil
}

func (x *HandleHelper) GetNodeByOID(oid types.Oid) (node SmiNode, err error) {
	smiNode := x.Handle.GetNodeByOID(oid)
	if smiNode == nil {
		err = fmt.Errorf("could not find node for OID %s", oid)
		return
	}
	return CreateNode(smiNode, x.Handle), nil
}

func (x *HandleHelper) GetType(name string, module ...SmiModule) (outType SmiType, err error) {
	var smiModule *types.SmiModule
	if len(module) > 0 {
		smiModule = module[0].GetRaw()
	}

	smiType := x.Handle.GetType(smiModule, name)
	if smiType == nil {
		if len(module) > 0 {
			err = fmt.Errorf("could not find type named %s in module %s", name, module[0].Name)
		} else {
			err = fmt.Errorf("could not find type named %s", name)
		}
		return
	}
	return CreateType(smiType, x.Handle), nil
}

func NewHandleHelper(handle *smi.Handle) *HandleHelper {
	return &HandleHelper{Handle: handle}
}
