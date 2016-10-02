package main

import (
	"fmt"
	"html/template"

	"github.com/konkers/blackfriday"
)

func handleMarkdown(args ...interface{}) template.HTML {
	commonHtmlFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	commonExtensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	renderer := blackfriday.HtmlRendererWithParameters(commonHtmlFlags, "", "",
		blackfriday.HtmlRendererParameters{
			InitialHeaderLevel: 5,
		})
	s := blackfriday.Markdown([]byte(fmt.Sprintf("%s", args...)), renderer, commonExtensions)
	return template.HTML(s)
}

func isAnonTypedef(arg interface{}) bool {
	t, ok := arg.(Type)
	if !ok {
		return false
	}
	return t.IsAnonTypedef()
}

func isStruct(arg interface{}) bool {
	_, ok := arg.(*Struct)
	return ok
}

func isTypedef(arg interface{}) bool {
	_, ok := arg.(*Typedef)
	return ok
}

var (
	templateFuncs = template.FuncMap{
		"isAnonTypedef": isAnonTypedef,
		"isStruct":      isStruct,
		"isTypedef":     isTypedef,
		"markdown":      handleMarkdown,
	}
)
