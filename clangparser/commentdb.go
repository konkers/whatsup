package clangparser

import (
	"fmt"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

type commentDb struct {
	TopComment string
	comments   map[clang.Cursor]string
}

func newCommentDb() *commentDb {
	return &commentDb{
		comments: make(map[clang.Cursor]string),
	}
}

func (db *commentDb) populate(tu clang.TranslationUnit, verbose bool) {
	tokens := tu.Tokenize(tu.TranslationUnitCursor().Extent())
	defer tu.DisposeTokens(tokens)

	var topComments []string

	for index, token := range tokens {
		if verbose {
			_, tokLine, _, _ := tu.TokenLocation(token).FileLocation()
			fmt.Printf("%s %d %d\n", token.Kind().Spelling(), index, tokLine)
		}
		if token.Kind() == clang.Token_Comment {
			commentLocation := tu.TokenLocation(token)
			_, commentLine, _, _ := commentLocation.FileLocation()

			if int(commentLine) == index+1 {
				topComments = append(topComments, tu.TokenSpelling(token))
			}

			prevCursor, prevLine := findInterestingCursor(&tu, tokens, index, -1)
			nextCursor, nextLine := findInterestingCursor(&tu, tokens, index, 1)
			comments := collapseComments(&tu, tokens, index)

			//fmt.Printf("comment(%d:%d) '%s'\n", commentLine, commentCol, comments)

			if prevCursor != nil && prevLine == commentLine {
				db.comments[prevCursor.CanonicalCursor()] = comments
			} else if nextCursor != nil && nextLine == commentLine+1 {
				db.comments[nextCursor.CanonicalCursor()] = comments
			}
		}
	}
	db.TopComment = strings.Join(topComments, "\n")
}

func (db *commentDb) find(cursor clang.Cursor) string {
	comment, ok := db.comments[cursor]
	if ok {
		return comment
	} else {
		return ""
	}
}

func isCursorInteresting(cursor clang.Cursor) bool {
	kind := cursor.Kind()
	// Only support C decls at the moment.
	if kind == clang.Cursor_StructDecl ||
		kind == clang.Cursor_UnionDecl ||
		kind == clang.Cursor_ClassDecl ||
		kind == clang.Cursor_EnumDecl ||
		kind == clang.Cursor_FieldDecl ||
		kind == clang.Cursor_EnumConstantDecl ||
		kind == clang.Cursor_FunctionDecl ||
		kind == clang.Cursor_VarDecl ||
		kind == clang.Cursor_ParmDecl ||
		kind == clang.Cursor_MacroDefinition {
		return true
	}
	return false
}

func findInterestingCursor(
	tu *clang.TranslationUnit,
	tokens []clang.Token,
	index int,
	dir int) (*clang.Cursor, uint32) {

	for {
		index += dir
		if index < 0 || index >= len(tokens) {
			return nil, 0
		}

		token := tokens[index]
		tokenLocation := tu.TokenLocation(token)
		cursor := tu.Cursor(tokenLocation)
		if !cursor.Location().Equal(tokenLocation) {
			continue
		}
		if isCursorInteresting(cursor) {
			_, line, _, _ := cursor.Location().FileLocation()
			return &cursor, line
		}
	}
}

func collapseComments(tu *clang.TranslationUnit,
	tokens []clang.Token, index int) string {

	token := tokens[index]
	_, lastLine, _, _ := tu.TokenLocation(token).FileLocation()
	comment := tu.TokenSpelling(token)

	for index -= 1; index >= 0; index -= 1 {
		token := tokens[index]
		_, curLine, _, _ := tu.TokenLocation(token).FileLocation()

		if token.Kind() != clang.Token_Comment || curLine+1 != lastLine {
			break
		}

		comment = tu.TokenSpelling(token) + "\n" + comment

		lastLine = curLine
	}
	return comment
}
