package main

import (
	"html/template"

	"github.com/konkers/whatsup/model"
)

type Project struct {
	Name string
}

type Module struct {
	Name          string
	Documentation string
}

type Constant struct {
}

type Variable struct {
}

type TypeBase struct {
	Name        string
	AnonTypedef bool
	TypeString  string // i.e. struct, enum, union...
	Comment     model.Comment
}

type Field struct {
	Name    string
	CType   template.HTML
	Comment model.Comment
}

type Struct struct {
	TypeBase

	Fields []*Field
}

type Typedef struct {
	TypeBase

	UnderlyingType string
}

type Arg struct {
	Name    string
	CType   template.HTML
	Comment string
}

type Function struct {
	Name       string
	Comment    model.Comment
	ReturnType template.HTML
	Prototype  template.HTML
	HasArgDocs bool
	Args       []*Arg
}

type TemplateData struct {
	Project   Project
	Module    Module
	Constants []*Constant
	Variables []*Variable
	Types     []Type
	Functions []*Function
}

type Type interface {
	GetName() string
	IsAnonTypedef() bool
	SetAnonTypedefName(string)
}

func (t *TypeBase) GetName() string {
	return t.Name
}

func (t *TypeBase) IsAnonTypedef() bool {
	return t.AnonTypedef
}

func (t *TypeBase) SetAnonTypedefName(name string) {
	t.Name = name
	t.AnonTypedef = true
}
