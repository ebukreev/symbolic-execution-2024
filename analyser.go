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
	AnalyseDynamically("some-file", "some-name")
}

func AnalyseStatically(file string, functionName string) Conditional {
	cfg := BuildCfg(file)
	for _, member := range cfg.Members {
		function, ok := member.(*ssa.Function)
		if ok && function != nil && function.Name() == functionName {
			return InterpretStatically(function)
		}
	}
	return Conditional{}
}

type Analyser struct {
	StatesQueue PriorityQueue
	Results     []DynamicInterpreter
}

func AnalyseDynamically(file string, functionName string) []DynamicInterpreter {
	analyser := Analyser{make(PriorityQueue, 0), make([]DynamicInterpreter, 0)}
	cfg := BuildCfg(file)
	for _, member := range cfg.Members {
		function, ok := member.(*ssa.Function)
		if ok && function != nil && function.Name() == functionName {
			interpreter := DynamicInterpreter{
				Function:      function,
				PathCondition: &Literal[bool]{true},
				Memory:        map[string]SymbolicExpression{},
			}
			analyser.StatesQueue.Push(&Item{value: interpreter, priority: 1})
			for analyser.StatesQueue.Len() != 0 {
				interpretationResults := InterpretDynamically(analyser.StatesQueue.Pop().(*Item).value)
				for _, interpretationResult := range interpretationResults {
					if interpretationResult.ReturnValue != nil {
						analyser.Results = append(analyser.Results, interpretationResult)
					} else {
						// TODO calculate priority
						analyser.StatesQueue.Push(&Item{value: interpretationResult, priority: 1})
					}
				}
			}
			break
		}
	}
	return analyser.Results
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
