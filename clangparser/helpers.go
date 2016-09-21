package clangparser

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

var (
	commentRegexp = regexp.MustCompile(`^/{2,} ?`)
	newlineRegexp = regexp.MustCompile(`\n`)
)

func getCanonicalCType(ctype clang.Type) string {
	switch ctype.Spelling() {
	default:
		return ctype.Spelling()
	case "_Bool":
		return "bool"
	}
}

func getCursorContents(cursor clang.Cursor) (string, error) {
	startFile, _, _, startOffset := cursor.Extent().Start().FileLocation()
	endFile, _, _, endOffset := cursor.Extent().End().FileLocation()
	if startFile != endFile {
		return "", fmt.Errorf("Start and end files differ. %s != %s.",
			startFile.Name(), endFile.Name())
	}

	f, err := os.Open(startFile.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, endOffset-startOffset)
	_, err = f.ReadAt(buf, int64(startOffset))
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func stripComment(comment string) string {
	comments := newlineRegexp.Split(comment, -1)
	for i := range comments {
		comments[i] = commentRegexp.ReplaceAllString(comments[i], "")
	}
	return strings.Join(comments, "\n")
}

func storageClassToString(class clang.StorageClass) string {
	switch class {
	default:
		return "invalid"
	case clang.SC_None:
		return "none"
	case clang.SC_Extern:
		return "extern"
	case clang.SC_Static:
		return "static"
	case clang.SC_PrivateExtern:
		return "private_extern"
	case clang.SC_OpenCLWorkGroupLocal:
		return "open_cl_work_group_local"
	case clang.SC_Auto:
		return "auto"
	case clang.SC_Register:
		return "register"
	}
}
