package main

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
)

func main() {
	Analyse("some-file")
}

func Analyse(file string) {
	cfg := BuildCfg(file)
	Interpret(cfg)
}

func BuildCfg(file string) *ssa.Package {
	fset := token.NewFileSet()
	fileContent, _ := os.ReadFile(file)
	f, _ := parser.ParseFile(fset, file, string(fileContent), 0)
	files := []*ast.File{f}

	pkg := types.NewPackage("main", "")

	main, _, _ := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, pkg, files, 0)

	return main
}
