package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/theliebeskind/genfig/generator"

	"github.com/theliebeskind/genfig/types"
	"github.com/theliebeskind/genfig/util"
)

var (
	dir = flag.String("dir", "config", "directory to write generated files into")
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\nERROR: %v\n", r)
			fmt.Println()
			flag.Usage()
			os.Exit(1)
		}
	}()
	exec()
}

func exec() {
	flag.Parse()
	fmt.Printf("Called with \n\tdir:\t'%s'\n\targs:\t%s\n", *dir, strings.Join(flag.Args(), ", "))
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"*"}
	}
	files := util.ResolveGlobs(args...)
	if len(files) == 0 {
		panic("No input files found")
	}

	params := types.Params{
		Dir: *dir,
	}
	fmt.Printf("Generating from files: %s\n", strings.Join(files, ", "))

	gofiles, err := generator.Generate(files, params)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	fmt.Printf("Successfully generated: %s", strings.Join(gofiles, ", "))
}
