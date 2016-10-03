package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/konkers/whatsup/model"
)

var (
	inputDir    = flag.String("input_dir", "", "Directory containing module json files.")
	outputDir   = flag.String("output_dir", "", "Directory to write html docs.")
	templateDir = flag.String("template_dir", "", "Directory containing templates.")
	projectName = flag.String("project_name", "", "Name of the project.")

	paramRefRegexp = regexp.MustCompile(`<(.*?)>`)
)

// These functions directly output HTML which violates the abstraction keeping
// all HTML in templates.  This could be resolved by creating custom markdown
// tags and letting the markdown processor handle this at template time.
//
// Alternately, and probably better is to create a template for each of these.
func decorate(elem string, class string) string {
	return fmt.Sprintf(`<span class="%s">%s</span>`, class, elem)
}

func decorateType(t string) string {
	return decorate(t, "hljs-keyword")
}

func decorateFunction(f string) string {
	return decorate(f, "hljs-title")
}

func decorateParam(p string) string {
	return decorate(p, "hljs-params")
}

func markupComment(comment string) string {
	comment = paramRefRegexp.ReplaceAllString(comment,
		`<code><span class="param_ref mdl-shadow--1dp">$1</span></code>`)

	return comment
}

func isPointer(t string) bool {
	return t[len(t)-1] == '*'
}

func transmogrifyFunction(in *model.Function) *Function {
	comment := model.ParseComment(in.Comment)

	returnType := decorateType(in.ReturnType)
	if !isPointer(in.ReturnType) {
		returnType += " "
	}

	prototype := decorateFunction(in.Name) + "("
	var argStrings []string
	var args []*Arg
	hasArgDocs := false

	for _, arg := range in.Args {
		ctype := arg.CType
		argType := decorateType(ctype)
		if !isPointer(ctype) {
			argType += " "
		}
		argStr := argType + decorateParam(arg.Name)
		argStrings = append(argStrings, argStr)

		commentStr, ok := comment.Args[arg.Name]
		hasArgDocs = hasArgDocs || ok

		a := &Arg{
			Name:    arg.Name,
			CType:   template.HTML(argType),
			Comment: markupComment(commentStr),
		}
		args = append(args, a)
	}

	if len(args) == 0 {
		prototype += decorateType("void")
	} else {
		prototype += strings.Join(argStrings, ", ")
	}
	prototype += ")"

	comment.Title = markupComment(comment.Title)
	comment.Body = markupComment(comment.Body)

	return &Function{
		Name:       in.Name,
		ReturnType: template.HTML(returnType),
		Prototype:  template.HTML(prototype),
		HasArgDocs: hasArgDocs,
		Args:       args,
		Comment:    *comment,
	}
}

func transmogrifyFields(in []*model.Field) []*Field {
	var fields []*Field
	for _, f := range in {
		fields = append(fields, &Field{
			Name:    f.Name,
			CType:   template.HTML(decorateType(f.CType)),
			Comment: *model.ParseComment(f.Comment),
		})
	}
	return fields
}

func transmogrifyStruct(in *model.Struct) *Struct {
	return &Struct{
		TypeBase: TypeBase{
			Name:       in.Name,
			TypeString: "struct",
			Comment:    *model.ParseComment(in.Comment),
		},
		Fields: transmogrifyFields(in.Fields),
	}
}

func transmogrifyUnion(in *model.Union) *Union {
	return &Union{
		TypeBase: TypeBase{
			Name:       in.Name,
			TypeString: "union",
			Comment:    *model.ParseComment(in.Comment),
		},
		Fields: transmogrifyFields(in.Fields),
	}
}

func transmogrifyEnum(in *model.Enum) *Enum {
	var constants []*EnumConstant
	for _, c := range in.Constants {

		constants = append(constants, &EnumConstant{
			Name:    c.Name,
			Value:   c.Value,
			Comment: *model.ParseComment(c.Comment),
		})
	}

	return &Enum{
		TypeBase: TypeBase{
			Name:       in.Name,
			TypeString: "enum",
			Comment:    *model.ParseComment(in.Comment),
		},
		Constants: constants,
	}
}

func transmogrifyTypedef(in *model.Typedef) *Typedef {
	return &Typedef{
		TypeBase: TypeBase{
			Name:       in.Name,
			TypeString: "typedef",
			Comment:    *model.ParseComment(in.Comment),
		},
		UnderlyingType: in.BaseType.Name,
	}
}

type byType []Type

func (t byType) Len() int      { return len(t) }
func (t byType) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t byType) Less(i, j int) bool {
	return strings.Compare(t[i].GetName(), t[j].GetName()) <= 0
}

func main() {
	flag.Parse()

	templates := template.Must(template.New("main").Funcs(templateFuncs).ParseGlob(filepath.Join(*templateDir, "*")))

	var jsonFiles []string
	filepath.Walk(*inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".json" {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})

	for _, path := range jsonFiles {
		module, err := loadModule(path)
		if err != nil {
			panic(err.Error())
		}
		data := TemplateData{
			Project: Project{
				Name: *projectName,
			},
			Module: Module{
				Name:          module.Name,
				Documentation: module.Documentation,
			},
		}

		types := make(map[string]Type)

		for _, function := range module.Functions {
			data.Functions = append(data.Functions, transmogrifyFunction(function))
		}

		for _, s := range module.Structs {
			structType := transmogrifyStruct(s)
			data.Types = append(data.Types, structType)
			types[s.Usr] = structType
		}

		for _, u := range module.Unions {
			unionType := transmogrifyUnion(u)
			data.Types = append(data.Types, unionType)
			types[u.Usr] = unionType
		}

		for _, e := range module.Enums {
			enumType := transmogrifyEnum(e)
			data.Types = append(data.Types, enumType)
			types[e.Usr] = enumType
		}

		for _, typedef := range module.Typedefs {
			if t, ok := types[typedef.BaseType.Usr]; ok {
				if t.GetName() == "" {
					t.SetAnonTypedefName(typedef.Name)
				} else {
					data.Types = append(data.Types, transmogrifyTypedef(typedef))
				}
			}
		}

		sort.Sort(byType(data.Types))

		err = templates.ExecuteTemplate(os.Stdout, "index.html", data)
		if err != nil {
			panic(err.Error())
		}
	}
}

func loadModule(path string) (*model.Module, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var module model.Module
	json.Unmarshal(bytes, &module)

	return &module, nil
}
