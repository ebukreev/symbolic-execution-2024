package se

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
	StatesQueue  PriorityQueue
	Results      []DynamicInterpreter
	Package      *ssa.Package
	Solver       *Solver
	SmtBuilder   SmtBuilder
	PathSelector PathSelector
}

func AnalyseMethodDynamically(file string, typeName string, functionName string) []DynamicInterpreter {
	pack := BuildCfg(file)
	for _, member := range pack.Members {
		typ, ok := member.(*ssa.Type)
		if ok && typ.Name() == typeName {
			named := typ.Type().(*types.Named)
			n := named.NumMethods()
			for i := 0; i < n; i++ {
				fun := pack.Prog.FuncValue(named.Method(i))
				if fun != nil && fun.Name() == functionName {
					return analyseDynamically(pack, fun)
				}
			}
		}
	}
	return nil
}

func AnalyseDynamically(file string, functionName string) []DynamicInterpreter {
	pack := BuildCfg(file)
	for _, member := range pack.Members {
		function, ok := member.(*ssa.Function)
		if ok && function != nil && function.Name() == functionName {
			return analyseDynamically(pack, function)
		}
	}
	return nil
}

func analyseDynamically(pack *ssa.Package, function *ssa.Function) []DynamicInterpreter {
	solver := CreateSolver(false)
	smtBuilder := SmtBuilder{Context: solver.Context}
	analyser := Analyser{make(PriorityQueue, 0), make([]DynamicInterpreter, 0), pack,
		solver, smtBuilder, &NursPathSelector{}}

	interpreter := DynamicInterpreter{
		CallStack:     []CallStackFrame{{Function: function, Memory: make(map[string]SymbolicExpression)}},
		Analyser:      &analyser,
		PathCondition: &Literal[bool]{true},
		Heap: &SymbolicMemory{make(map[string]Array),
			make(map[string]map[int]string), make(map[string]int), &smtBuilder},
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
				analyser.StatesQueue.Push(&Item{value: interpretationResult, priority: analyser.PathSelector.CalculatePriority(interpretationResult)})
			} else {
				analyser.StatesQueue.Push(&Item{value: interpretationResult, priority: analyser.PathSelector.CalculatePriority(interpretationResult)})
			}
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

	main, _, err := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, pkg, files, 0)
	if err != nil {
		panic(err)
	}

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
