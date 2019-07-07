package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/theliebeskind/go-genfig/generator"

	"github.com/theliebeskind/go-genfig/types"
	"github.com/theliebeskind/go-genfig/util"
)

// PROJECT is the name of this project
var PROJECT = "genfig"

var (
	dir = flag.String("dir", "config", "directory to write generated files into")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", PROJECT)
		flag.PrintDefaults()
	}
}

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
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"*"}
	}

	fmt.Printf("Called with \n\tdir:\t'%s'\n\targs:\t%s\n\n", *dir, strings.Join(args, ", "))

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
