package model

type Module struct {
	Name          string `json:"name"`
	Documentation string `json:"Documentation"`
	File          string `json:"File"`

	Enums          map[string]*Enum          `json:"enums"`
	Functions      map[string]*Function      `json:"functions"`
	Macros         map[string]*Macro         `json:"macros"`
	MacroFunctions map[string]*MacroFunction `json:"macro_functions"`
	Structs        map[string]*Struct        `json:"structs"`
	Typedefs       map[string]*Typedef       `json:"typedefs"`
	Unions         map[string]*Union         `json:"unions"`
	Variables      map[string]*Variable      `json:"variables"`
}

func NewModule() *Module {
	return &Module{
		Enums:          make(map[string]*Enum),
		Functions:      make(map[string]*Function),
		Macros:         make(map[string]*Macro),
		MacroFunctions: make(map[string]*MacroFunction),
		Structs:        make(map[string]*Struct),
		Typedefs:       make(map[string]*Typedef),
		Unions:         make(map[string]*Union),
		Variables:      make(map[string]*Variable),
	}
}
