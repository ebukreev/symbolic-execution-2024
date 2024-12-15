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
		t.Fatal(err)
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

	if isFloatResult {
		resultSymbolicVar := &InputValue{Name: "res", Type: "float64"}

		res := smtBuilder.BuildSmt(resultSymbolicVar)[0]

		expressionSmt := smtBuilder.BuildSmt(resultExpression)
		solver.SmtSolver.Assert(res.(z3.Float).Eq(expressionSmt[0].(z3.Float)))
	} else {
		realResultSymbolicVar := &InputValue{Name: "$R_res", Type: "float64"}
		imaginaryResultSymbolicVar := &InputValue{Name: "$I_res", Type: "float64"}

		resReal := smtBuilder.BuildSmt(realResultSymbolicVar)[0]
		resImaginary := smtBuilder.BuildSmt(imaginaryResultSymbolicVar)[0]

		expressionSmt := smtBuilder.BuildSmt(resultExpression)
		solver.SmtSolver.Assert(resReal.(z3.Float).Eq(expressionSmt[0].(z3.Float)))
		solver.SmtSolver.Assert(resImaginary.(z3.Float).Eq(expressionSmt[1].(z3.Float)))
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
