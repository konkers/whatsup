// go-clang-dump shows how to dump the AST of a C/C++ file via the Cursor
// visitor API.
//
// ex:
// $ go-clang-dump -fname=foo.cxx
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/konkers/whatsup/clangparser"
)

var fname = flag.String("fname", "", "the file to analyze")
var verbose = flag.Bool("verbose", false, "verbose debugging output")

func main() {
	flag.Parse()

	if *fname == "" {
		flag.Usage()
		fmt.Printf("please provide a file name to analyze\n")
		os.Exit(1)
	}

	args := []string{}
	if len(flag.Args()) > 0 && flag.Args()[0] == "-" {
		args = make([]string, len(flag.Args()[1:]))
		copy(args, flag.Args()[1:])
	}

	parser := clangparser.NewParser(*fname, args, *verbose)
	parser.Parse()
}
