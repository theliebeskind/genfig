package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go/parser"
	"go/token"

	"github.com/theliebeskind/genfig/generator"

	"github.com/theliebeskind/genfig/models"
	"github.com/theliebeskind/genfig/util"
)

// PROJECT is the name of this project
var PROJECT = "genfig"

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
	run()
}

func run() {
	var (
		help = flag.Bool("help", false, "print this usage help")
		dir  = flag.String("dir", "./config", "directory to write generated files into")
	)

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"*"}
	}

	fmt.Printf("Called with \n\tdir:\t'%s'\n\targs:\t%s\n\n", *dir, strings.Join(args, ", "))

	files := util.ResolveGlobs(args...)
	if len(files) == 0 {
		panic("No input files found")
	}

	params := models.Params{
		Dir: *dir,
	}
	fmt.Printf("Generating from files: %s\n", strings.Join(files, ", "))

	gofiles, err := generator.Generate(files, params)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	fmt.Println("\nChecking generaded code ...")
	path, _ := filepath.Abs(*dir)
	fset := token.NewFileSet()
	if _, err := parser.ParseDir(fset, path, nil, 0); err != nil {
		panic(fmt.Sprintf("At least one error in generated code: %v", err))
	}

	fmt.Println("\nFormatting generade code with gofmt ...")
	if err := exec.Command("gofmt", "-w", ".").Run(); err != nil {
		panic(fmt.Sprintf("Could not format code: %v", err))
	}

	fmt.Printf("\nSuccessfully generated %d files: %s\n", len(gofiles), strings.Join(gofiles, ", "))
}
