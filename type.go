package gosmi

import (
	"fmt"

	"github.com/opsbl/gosmi/models"
	"github.com/opsbl/gosmi/types"
)

type SmiType struct {
	models.Type
	smiType  *types.SmiType
	instance *Instance
}

func (t *SmiType) getEnum() {
	if t.BaseType == types.BaseTypeUnknown || !(t.BaseType == types.BaseTypeEnum || t.BaseType == types.BaseTypeBits) {
		return
	}

	smiNamedNumber := t.instance.smiInst.GetFirstNamedNumber(t.smiType)
	if smiNamedNumber == nil {
		return
	}

	enum := models.Enum{
		BaseType: types.BaseType(smiNamedNumber.Value.BaseType),
	}
	for ; smiNamedNumber != nil; smiNamedNumber = t.instance.smiInst.GetNextNamedNumber(smiNamedNumber) {
		namedNumber := models.NamedNumber{
			Name:  string(smiNamedNumber.Name),
			Value: convertValue(smiNamedNumber.Value),
		}
		enum.Values = append(enum.Values, namedNumber)
	}
	t.Enum = &enum
	return
}

func (t SmiType) GetModule() (module SmiModule) {
	smiModule := t.instance.smiInst.GetTypeModule(t.smiType)
	return CreateModule(t.instance, smiModule)
}

func (t *SmiType) getRanges() {
	if t.BaseType == types.BaseTypeUnknown {
		return
	}

	ranges := make([]models.Range, 0)
	// Workaround for libsmi bug that causes ranges to loop infinitely sometimes
	var currSmiRange *types.SmiRange
	for smiRange := t.instance.smiInst.GetFirstRange(t.smiType); smiRange != nil && smiRange != currSmiRange; smiRange = t.instance.smiInst.GetNextRange(smiRange) {
		r := models.Range{
			BaseType: smiRange.MinValue.BaseType,
			MinValue: convertValue(smiRange.MinValue),
			MaxValue: convertValue(smiRange.MaxValue),
		}
		ranges = append(ranges, r)
		currSmiRange = smiRange
	}
	t.Ranges = ranges
}

func (t SmiType) String() string {
	return t.Type.String()
}

func (t SmiType) GetRaw() (outType *types.SmiType) {
	return t.smiType
}

func (t *SmiType) SetRaw(smiType *types.SmiType) {
	t.smiType = smiType
}

func CreateType(inst *Instance, smiType *types.SmiType) (outType SmiType) {
	if smiType == nil {
		return
	}

	outType.SetRaw(smiType)
	outType.BaseType = smiType.BaseType
	outType.instance = inst

	if smiType.Name == "" {
		smiType = inst.smiInst.GetParentType(smiType)
	}

	outType.Decl = smiType.Decl
	outType.Description = smiType.Description
	outType.Format = smiType.Format
	outType.Name = string(smiType.Name)
	outType.Reference = smiType.Reference
	outType.Status = smiType.Status
	outType.Units = smiType.Units

	outType.getEnum()
	outType.getRanges()

	return
}

func CreateTypeFromNode(inst *Instance, smiNode *types.SmiNode) (outType *SmiType) {
	if inst == nil {
		return nil
	}
	smiType := inst.smiInst.GetNodeType(smiNode)

	if smiType == nil {
		return
	}

	tempType := CreateType(inst, smiType)
	outType = &tempType

	if smiNode.Format != "" {
		outType.Format = smiNode.Format
	}
	if smiNode.Units != "" {
		outType.Units = smiNode.Units
	}

	return
}

func (i *Instance) GetType(name string, module ...SmiModule) (outType SmiType, err error) {
	var smiModule *types.SmiModule
	if len(module) > 0 {
		smiModule = module[0].GetRaw()
	}

	smiType := i.smiInst.GetType(smiModule, name)
	if smiType == nil {
		if len(module) > 0 {
			err = fmt.Errorf("Could not find type named %s in module %s", name, module[0].Name)
		} else {
			err = fmt.Errorf("Could not find type named %s", name)
		}
		return
	}
	return CreateType(i, smiType), nil
}

func convertValue(value types.SmiValue) (outValue int64) {
	switch v := value.Value.(type) {
	case int32:
		outValue = int64(v)
	case int64:
		outValue = int64(v)
	case uint32:
		outValue = int64(v)
	case uint64:
		outValue = int64(v)
	}
	return
}
