package model

type ObjectType int

const (
	ObjectTypeStruct ObjectType = iota
	ObjectTypeUnion
	ObjectTypeEnum
	ObjectTypeField
	ObjectTypeEnumConstant
	ObjectTypeFunction
	ObjectTypeVariable
	ObjectTypeParam
	ObjectTypeTypedef
	ObjectTypeMacro
	ObjectTypeMacroFunction
	ObjectTypeBasic
)

type Object interface {
	GetName() string
	GetType() ObjectType
	GetComment() string
	GetUsr() string
}

type Base struct {
	Name    string     `json:"name"`
	Type    ObjectType `json:"type"`
	Comment string     `json:"comment"`
	Usr     string     `json:"usr"`

	File string `json:"file"`
	Line uint32 `json:"line"`
	Col  uint32 `json:"col"`
}

// TODO(konkers): Are any of these aside from GetType() necessary?
func (b *Base) GetName() string {
	return b.Name
}

func (b *Base) GetType() ObjectType {
	return b.Type
}

func (b *Base) GetComment() string {
	return b.Comment
}

func (b *Base) GetUsr() string {
	return b.Usr
}

type Basic struct {
	Base
}

type Enum struct {
	Base

	Constants []*EnumConstant `json:"constants"`
}

type EnumConstant struct {
	Base

	Value int64 `json:"value"`
}

type Field struct {
	Base

	CType string `json:"ctype"`
}

type Function struct {
	Base

	StorageClass string `json:"storage_class"`
	Inlined      bool   `json:"inlined"`

	Args       []*Param `json:"args"`
	ReturnType string   `json:"return_type"`
}

type Macro struct {
	Base

	Value string `json:"value"`
}

type MacroFunction struct {
	Base

	Args []string `json:"args"`
}

type Param struct {
	Base

	CType string `json:"ctype"`
}

type Struct struct {
	Base

	Fields []*Field `json:"fields"`
}

type Typedef struct {
	Base

	BaseType *Basic `json:"base_type"`
}

type Union struct {
	Base

	Fields []*Field `json:"fields"`
}

type Variable struct {
	Base

	CType        string `json:"ctype"`
	StorageClass string `json:"storage_class"`
}
