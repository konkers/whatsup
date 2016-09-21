package clangparser

// TODO(konkers): Convert panics into returned errors.

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/go-clang/v3.9/clang"

	"github.com/konkers/whatsup/model"
)

var (
	argSplitRegexp        = regexp.MustCompile(`\s*,\s*`)
	macroFunctionRegexp   = regexp.MustCompile(`\(\s*(.*?)\s*\)`)
	whitespaceSplitRegexp = regexp.MustCompile(`\s+`)
)

type Parser struct {
	fileName string
	args     []string

	verbose bool

	comments *commentDb

	// This holds types until the we're done with parsing.  We keep them out of
	// the module because we a type is repeated twice when it's inline into a
	// typedef or a structure.  This lets us easily re-parent then under their
	// parent type.
	types     map[string]model.Object
	objectMap map[string]model.Object

	module *model.Module
}

func NewParser(fileName string, args []string, verbose bool) *Parser {
	return &Parser{
		fileName:  fileName,
		args:      args,
		verbose:   verbose,
		comments:  newCommentDb(),
		types:     make(map[string]model.Object),
		objectMap: make(map[string]model.Object),
		module:    model.NewModule(),
	}
}

func (p *Parser) newBase(cursor clang.Cursor, objectType model.ObjectType) model.Base {
	file, line, col, _ := cursor.Location().FileLocation()
	return model.Base{
		Name:    cursor.Spelling(),
		Type:    objectType,
		Comment: stripComment(p.comments.find(cursor)),
		Usr:     cursor.USR(),

		File: file.Name(),
		Line: line,
		Col:  col,
	}
}

func (p *Parser) handleEnumCursor(cursor, parent clang.Cursor) bool {
	enum := &model.Enum{
		Base: p.newBase(cursor, model.ObjectTypeEnum),
	}

	p.module.Enums[enum.GetUsr()] = enum
	p.objectMap[enum.GetUsr()] = enum

	return true
}

func (p *Parser) handleEnumConstantCursor(cursor, parent clang.Cursor) bool {
	parentObject, ok := p.objectMap[parent.USR()].(*model.Enum)
	if !ok {
		panic("Got an enum constant in not an enum!\n")
	}

	constant := &model.EnumConstant{
		Base:  p.newBase(cursor, model.ObjectTypeEnumConstant),
		Value: cursor.EnumConstantDeclValue(),
	}

	parentObject.Constants = append(parentObject.Constants, constant)

	return false
}

func (p *Parser) handleFieldCursor(cursor, parent clang.Cursor) bool {
	field := &model.Field{
		Base:  p.newBase(cursor, model.ObjectTypeField),
		CType: getCanonicalCType(cursor.Type()),
	}

	parentStruct, ok := p.objectMap[parent.USR()].(*model.Struct)
	if ok {
		parentStruct.Fields = append(parentStruct.Fields, field)
		return false
	}
	parentUnion, ok := p.objectMap[parent.USR()].(*model.Union)
	if ok {
		parentUnion.Fields = append(parentUnion.Fields, field)
		return false
	}

	panic("Got a field not in an Union or struct!\n")
}

func (p *Parser) handleFunctionCursor(cursor, parent clang.Cursor) bool {
	function := &model.Function{
		Base:         p.newBase(cursor, model.ObjectTypeFunction),
		StorageClass: storageClassToString(cursor.StorageClass()),
		Inlined:      cursor.IsFunctionInlined(),
		ReturnType:   getCanonicalCType(cursor.ResultType()),
	}

	p.module.Functions[function.GetUsr()] = function
	p.objectMap[function.GetUsr()] = function

	return true
}

func (p *Parser) handleMacroCursor(cursor, parent clang.Cursor) bool {
	content, err := getCursorContents(cursor)
	if err != nil {
		panic(err.Error())
	}

	if cursor.IsMacroFunctionLike() {
		arg_string := macroFunctionRegexp.FindStringSubmatch(content)
		args := argSplitRegexp.Split(arg_string[1], -1)
		macroFunction := &model.MacroFunction{
			Base: p.newBase(cursor, model.ObjectTypeMacroFunction),
			Args: args,
		}
		p.module.MacroFunctions[macroFunction.GetUsr()] = macroFunction
		p.objectMap[macroFunction.GetUsr()] = macroFunction

	} else {
		vals := whitespaceSplitRegexp.Split(content, 2)
		if len(vals) == 2 {
			macro := &model.Macro{
				Base:  p.newBase(cursor, model.ObjectTypeMacro),
				Value: vals[1],
			}
			p.module.Macros[macro.GetUsr()] = macro
			p.objectMap[macro.GetUsr()] = macro
		}
	}

	return false
}

func (p *Parser) handleParmCursor(cursor, parent clang.Cursor) bool {
	param := &model.Param{
		Base:  p.newBase(cursor, model.ObjectTypeParam),
		CType: getCanonicalCType(cursor.Type()),
	}

	function, ok := p.objectMap[parent.USR()].(*model.Function)
	if !ok {
		panic("Got a parm not in a function!")
	}

	function.Args = append(function.Args, param)
	return false
}

func (p *Parser) handleStructCursor(cursor, parent clang.Cursor) bool {
	structObj := &model.Struct{
		Base: p.newBase(cursor, model.ObjectTypeStruct),
	}

	p.module.Structs[structObj.GetUsr()] = structObj
	p.objectMap[structObj.GetUsr()] = structObj

	return true
}

func (p *Parser) handleTypedefCursor(cursor, parent clang.Cursor) bool {
	typedef := &model.Typedef{
		Base: p.newBase(cursor, model.ObjectTypeTypedef),
		BaseType: &model.Basic{model.Base{
			Name: cursor.TypedefDeclUnderlyingType().Spelling(),
			Usr:  cursor.TypedefDeclUnderlyingType().Declaration().USR(),
			Type: model.ObjectTypeBasic,
		}},
	}
	p.module.Typedefs[typedef.GetUsr()] = typedef
	p.objectMap[typedef.GetUsr()] = typedef
	return true
}

func (p *Parser) handleUnionCursor(cursor, parent clang.Cursor) bool {
	union := &model.Union{
		Base: p.newBase(cursor, model.ObjectTypeUnion),
	}
	p.module.Unions[union.GetUsr()] = union
	p.objectMap[union.GetUsr()] = union
	return true
}

func (p *Parser) handleVarCursor(cursor, parent clang.Cursor) bool {
	variable := &model.Variable{
		Base:         p.newBase(cursor, model.ObjectTypeVariable),
		CType:        getCanonicalCType(cursor.Type()),
		StorageClass: storageClassToString(cursor.StorageClass()),
	}
	p.module.Variables[variable.GetUsr()] = variable
	p.objectMap[variable.GetUsr()] = variable
	return true
}

func (p *Parser) handleCursor(cursor, parent clang.Cursor) bool {
	var recurse bool

	// Only process a cursor if we don't already have an object for it's USR.
	if _, ok := p.objectMap[cursor.USR()]; !ok {
		switch cursor.Kind() {
		case clang.Cursor_EnumDecl:
			recurse = p.handleEnumCursor(cursor, parent)
		case clang.Cursor_EnumConstantDecl:
			recurse = p.handleEnumConstantCursor(cursor, parent)
		case clang.Cursor_FieldDecl:
			recurse = p.handleFieldCursor(cursor, parent)
		case clang.Cursor_FunctionDecl:
			recurse = p.handleFunctionCursor(cursor, parent)
		case clang.Cursor_MacroDefinition:
			recurse = p.handleMacroCursor(cursor, parent)
		case clang.Cursor_ParmDecl:
			recurse = p.handleParmCursor(cursor, parent)
		case clang.Cursor_StructDecl:
			recurse = p.handleStructCursor(cursor, parent)
		case clang.Cursor_TypedefDecl:
			recurse = p.handleTypedefCursor(cursor, parent)
		case clang.Cursor_UnionDecl:
			recurse = p.handleUnionCursor(cursor, parent)
		case clang.Cursor_VarDecl:
			recurse = p.handleVarCursor(cursor, parent)
		}
	}

	// Here we fix up typedefs.
	//	if parentObject, ok := p.objectMap[parent.USR()]; ok &&
	//		parentObject.GetType() == model.ObjectTypeTypedef {
	//		obj := p.objectMap[cursor.USR()]
	//		parentObject.(*model.Typedef).BaseType = obj
	//
	//		delete(p.types, cursor.USR())
	//	}

	return recurse
}

func (p *Parser) visit(cursor clang.Cursor, indent string) {
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() || !cursor.Location().IsFromMainFile() {
			return clang.ChildVisit_Continue
		}

		if p.verbose {
			fmt.Printf("%s%s: %s (%s) [%s]\n", indent,
				cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(),
				cursor.LexicalParent().Kind().Spelling())
		}
		recurse := p.handleCursor(cursor.CanonicalCursor(), parent)
		if recurse {
			// Handle recursion ourselves so that we can nicely indent debug output.
			p.visit(cursor, indent+"  ")
		}

		return clang.ChildVisit_Continue

	})
}

func (p *Parser) Parse() {
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit(p.fileName, p.args, nil,
		clang.TranslationUnit_DetailedPreprocessingRecord)
	defer tu.Dispose()

	diagnostics := tu.Diagnostics()
	for _, d := range diagnostics {
		fmt.Println("PROBLEM:", d.Spelling())
	}

	p.comments.populate(tu, p.verbose)

	cursor := tu.TranslationUnitCursor()
	p.visit(cursor, "")

	comment := model.ParseComment(stripComment(p.comments.TopComment))
	p.module.Name = comment.Title
	p.module.File = cursor.Spelling()
	p.module.Documentation = comment.Body

	output, _ := json.MarshalIndent(p.module, "", "  ")
	fmt.Println(string(output))

	if len(diagnostics) > 0 {
		fmt.Println("NOTE: There were problems while analyzing the given file")
	}
}
