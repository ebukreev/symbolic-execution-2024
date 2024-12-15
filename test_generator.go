package se

import (
	"fmt"
	"golang.org/x/tools/go/ssa"
	"strconv"
	"strings"
	"symbolic-execution-2024/z3"
)

func (interpreter *DynamicInterpreter) GenerateTest(testName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("func %v(t *testing.T) {\n", testName))

	smtSolver := interpreter.Analyser.Solver.SmtSolver
	smtBuilder := interpreter.Analyser.SmtBuilder
	smtSolver.Reset()
	smtSolver.Assert(smtBuilder.BuildSmt(interpreter.PathCondition)[0].(z3.Bool))

	returnSmt := smtBuilder.BuildSmt(interpreter.CurrentFrame().ReturnValue)[0]

	_, ok := returnSmt.(z3.Float)
	var typ string
	if ok {
		typ = "float64"
		smtSolver.Assert(smtBuilder.BuildSmt(&InputValue{Name: "res", Type: typ})[0].(z3.Float).Eq(returnSmt.(z3.Float)))
	} else {
		typ = "int"
		smtSolver.Assert(smtBuilder.BuildSmt(&InputValue{Name: "res", Type: typ})[0].(z3.BV).Eq(returnSmt.(z3.BV)))
	}

	_, err := smtSolver.Check()
	if err != nil {
		panic(err)
	}

	function := interpreter.CurrentFrame().Function
	sb.WriteString(fmt.Sprintf("    callResult := %v(%v)\n", function.Name(),
		strings.Join(parseArgs(function, smtSolver.Model(), smtBuilder.Context), ", ")))
	sb.WriteString(fmt.Sprintf("    returnValue := %v\n", getValue("res", typ, smtSolver.Model(),
		smtBuilder.Context)))
	sb.WriteString("    if callResult != returnValue {\n")
	sb.WriteString("        t.Fatalf(\"%v != %v\", callResult, returnValue)\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	return sb.String()
}

func getValue(name string, typ string, model *z3.Model, context *z3.Context) string {
	if typ == "int" {
		evaluatedExpr := model.Eval(context.Const(name, context.BVSort(64)), true).String()[2:]
		i, _ := strconv.ParseUint(evaluatedExpr, 16, 64)
		return strconv.Itoa(int(i))
	}
	if typ == "float64" {
		evaluatedExpr := model.Eval(context.Const(name, context.FloatSort(11, 53)), true)
		if strings.HasPrefix(evaluatedExpr.String(), "(_ NaN") {
			return "math.NaN()"
		}
		if strings.HasPrefix(evaluatedExpr.String(), "(_ +zero") {
			return "0"
		}
		evaluatedExpr = context.Simplify(evaluatedExpr.(z3.Float).ToIEEEBV(), z3.NewSimplifyConfig(context))
		i, _ := strconv.ParseFloat(evaluatedExpr.String()[2:], 64)
		return fmt.Sprintf("%f", i)
	}
	panic("unexpected type")
}

func parseArgs(fn *ssa.Function, model *z3.Model, context *z3.Context) []string {
	var args []string
	for _, param := range fn.Params {
		args = append(args, getValue(param.Name(), param.Type().String(), model, context))
	}
	return args
}
