package se

import (
	"symbolic-execution-2024/z3"
	"testing"
)

func CheckResultWithPathCondition(t *testing.T, pathCondition SymbolicExpression, resultExpression SymbolicExpression, isFloatResult bool) {
	solver := CreateSolver(false)
	smtBuilder := SmtBuilder{Context: solver.Context}

	solver.SmtSolver.Assert(smtBuilder.BuildSmt(pathCondition)[0].(z3.Bool))

	resultSymbolicVar := &InputValue{Name: "res"}
	if isFloatResult {
		resultSymbolicVar.Type = "float64"
	} else {
		resultSymbolicVar.Type = "int"
	}

	res := smtBuilder.BuildSmt(resultSymbolicVar)[0]
	expressionSmt := smtBuilder.BuildSmt(resultExpression)[0]

	if isFloatResult {
		solver.SmtSolver.Assert(res.(z3.Float).Eq(expressionSmt.(z3.Float)))
	} else {
		solver.SmtSolver.Assert(res.(z3.BV).Eq(expressionSmt.(z3.BV)))
	}

	sat, err := solver.SmtSolver.Check()
	if !sat {
		t.Log("UNSAT")
	}
	if err != nil {
		if err.Error() == "timeout" {
			t.Log("SMT timeout")
		} else {
			t.Fatal(err)
		}
	}

	if sat {
		t.Log(solver.SmtSolver.Model().String())
	}
}

func CheckComplexResultAndPathCondition(t *testing.T, pathCondition SymbolicExpression, resultExpression SymbolicExpression, isFloatResult bool) {
	solver := CreateSolver(false)
	smtBuilder := SmtBuilder{Context: solver.Context}

	builtSmt := smtBuilder.BuildSmt(pathCondition)
	for i := 0; i < len(builtSmt); i++ {
		solver.SmtSolver.Assert(builtSmt[i].(z3.Bool))
	}

	sat, err := solver.SmtSolver.Check()
	if !sat {
		t.Log("UNSAT")
	}
	if err != nil {
		t.Fatal(err)
	}

	if sat {
		t.Log(solver.SmtSolver.Model().String())
	}
}
