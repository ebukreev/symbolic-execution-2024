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
	"strings"
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
	Package     *ssa.Package
	Solver      *Solver
	SmtBuilder  SmtBuilder
}

func AnalyseDynamically(file string, functionName string) []DynamicInterpreter {
	solver := CreateSolver(false)
	smtBuilder := SmtBuilder{Context: solver.Context}
	analyser := Analyser{make(PriorityQueue, 0), make([]DynamicInterpreter, 0), BuildCfg(file),
		solver, smtBuilder}
	for _, member := range analyser.Package.Members {
		function, ok := member.(*ssa.Function)
		if ok && function != nil && function.Name() == functionName {
			interpreter := DynamicInterpreter{
				CallStack:     []CallStackFrame{{Function: function, Memory: make(map[string]SymbolicExpression)}},
				Analyser:      &analyser,
				PathCondition: &Literal[bool]{true},
			}
			analyser.StatesQueue.Push(&Item{value: interpreter, priority: 1})
			for analyser.StatesQueue.Len() != 0 {
				interpretationResults := InterpretDynamically(analyser.StatesQueue.Pop().(*Item).value)
				for _, interpretationResult := range interpretationResults {
					if len(interpretationResult.CallStack) == 1 && interpretationResult.CurrentFrame().ReturnValue != nil {
						analyser.Results = append(analyser.Results, interpretationResult)
					} else if interpretationResult.CurrentFrame().ReturnValue != nil {
						interpretationResult.CallStack[len(interpretationResult.CallStack)-2].ReturnValue =
							interpretationResult.CurrentFrame().ReturnValue
						interpretationResult.CallStack = interpretationResult.CallStack[:len(interpretationResult.CallStack)-1]
						interpretationResult.CurrentFrame().InstructionsPtr--
						// TODO calculate priority
						analyser.StatesQueue.Push(&Item{value: interpretationResult, priority: 1})
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

func (analyser *Analyser) ResolveFunctionDeclaration(signature string) *ssa.Function {
	for _, member := range analyser.Package.Members {
		function, ok := member.(*ssa.Function)
		if ok && function != nil && strings.HasPrefix(signature, function.String()) {
			return function
		}
	}
	panic("unexpected signature " + signature)
}
