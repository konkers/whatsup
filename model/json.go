package model

import (
	"encoding/json"
	"fmt"
)

var (
	objectTypeNameMap = make(map[ObjectType]string)
	objectTypeMap     = make(map[string]ObjectType)
)

func addObjectTypeNameMapping(objectType ObjectType, name string) {
	objectTypeNameMap[objectType] = name
	objectTypeMap[name] = objectType
}

func init() {
	addObjectTypeNameMapping(ObjectTypeStruct, "struct")
	addObjectTypeNameMapping(ObjectTypeUnion, "union")
	addObjectTypeNameMapping(ObjectTypeEnum, "enum")
	addObjectTypeNameMapping(ObjectTypeField, "field")
	addObjectTypeNameMapping(ObjectTypeEnumConstant, "enum_constant")
	addObjectTypeNameMapping(ObjectTypeFunction, "function")
	addObjectTypeNameMapping(ObjectTypeVariable, "variable")
	addObjectTypeNameMapping(ObjectTypeParam, "param")
	addObjectTypeNameMapping(ObjectTypeTypedef, "typedef")
	addObjectTypeNameMapping(ObjectTypeMacro, "macro")
	addObjectTypeNameMapping(ObjectTypeMacroFunction, "macro_function")
	addObjectTypeNameMapping(ObjectTypeBasic, "basic")
}

func (o *ObjectType) MarshalJSON() ([]byte, error) {
	name, ok := objectTypeNameMap[*o]
	if !ok {
		return nil, fmt.Errorf("Can't find name mapping for '%q'\n", o)
	}
	return json.Marshal(name)
}

func (o *ObjectType) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	objectType, ok := objectTypeMap[name]
	if !ok {
		return fmt.Errorf("Can't find object type mapping for '%s'\n", name)
	}
	*o = objectType
	return nil
}
