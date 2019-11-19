//go:generate sh -c "printf 'package main\n\nfunc init() {\n\tversion = \"%s\"\n}\n' $(cat VERSION) > version.go"

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

	"github.com/thlcodes/genfig/generator"

	"github.com/thlcodes/genfig/models"
	"github.com/thlcodes/genfig/util"
)

const project = "genfig"

var version = "v0.0.0-dev"

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s %s:\n", project, version)
		flag.PrintDefaults()
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\nERROR: %v\n\n", r)
			flag.Usage()
			os.Exit(1)
		}
	}()
	run()
}

func run() {
	var (
		helpFlag    = flag.Bool("help", false, "print this usage help")
		versionFlag = flag.Bool("version", false, "print version")
		dir         = flag.String("dir", "./config", "directory to write generated files into")
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s %s", project, version)
		return
	}
	if *helpFlag {
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
		fmt.Printf("\nCould not format code: %v, continuing anyway.", err)
	}

	fmt.Printf("\nSuccessfully generated %d files: %s\n", len(gofiles), strings.Join(gofiles, ", "))
}
