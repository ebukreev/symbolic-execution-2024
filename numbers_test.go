package main

import (
	"symbolic-execution-2024/z3"
	"testing"
)

func TestIntegerOperationsFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}
	pc := &GT{Left: a, Right: b}

	resExpr := &BinaryOperation{
		Left:  a,
		Right: b,
		Type:  Add,
	}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestIntegerOperationsSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}
	pc := &BinaryOperation{&Not{&GT{a, b}}, &LT{a, b}, And}

	resExpr := &BinaryOperation{
		Left:  a,
		Right: b,
		Type:  Sub,
	}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestIntegerOperationsThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}
	pc := &BinaryOperation{&Not{&GT{a, b}}, &Not{&LT{a, b}}, And}

	resExpr := &BinaryOperation{
		Left:  a,
		Right: b,
		Type:  Mul,
	}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestFloatOperationsFirstPath(t *testing.T) {
	x := &InputValue{Name: "x", Type: "float64"}
	y := &InputValue{Name: "y", Type: "float64"}
	pc := &GT{Left: x, Right: y}

	resExpr := &BinaryOperation{
		Left:  x,
		Right: y,
		Type:  Div,
	}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestFloatOperationsSecondPath(t *testing.T) {
	x := &InputValue{Name: "x", Type: "float64"}
	y := &InputValue{Name: "y", Type: "float64"}
	pc := &BinaryOperation{&Not{&GT{x, y}}, &LT{x, y}, And}

	resExpr := &BinaryOperation{
		Left:  x,
		Right: y,
		Type:  Div,
	}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestFloatOperationsThirdPath(t *testing.T) {
	x := &InputValue{Name: "x", Type: "float64"}
	y := &InputValue{Name: "y", Type: "float64"}
	pc := &BinaryOperation{&Not{&GT{x, y}}, &Not{&LT{x, y}}, And}

	resExpr := &Literal[float64]{0.0}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "float64"}

	firstUpdate := &BinaryOperation{&Cast{a, "float64"}, b, Add}
	firstIf := &Equals{&BinaryOperation{a, &Literal[int]{2}, Mod}, &Literal[int]{0}}
	pc := &BinaryOperation{firstIf, &LT{firstUpdate, &Literal[float64]{10.0}}, And}

	resExpr := &BinaryOperation{firstUpdate, &Literal[float64]{2.0}, Mul}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "float64"}

	firstUpdate := &BinaryOperation{&Cast{a, "float64"}, b, Sub}
	firstIf := &Not{&Equals{&BinaryOperation{a, &Literal[int]{2}, Mod}, &Literal[int]{0}}}
	pc := &BinaryOperation{firstIf, &LT{firstUpdate, &Literal[float64]{10.0}}, And}

	resExpr := &BinaryOperation{firstUpdate, &Literal[float64]{2.0}, Mul}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "float64"}

	firstUpdate := &BinaryOperation{&Cast{a, "float64"}, b, Add}
	firstIf := &Equals{&BinaryOperation{a, &Literal[int]{2}, Mod}, &Literal[int]{0}}
	pc := &BinaryOperation{firstIf, &Not{&LT{firstUpdate, &Literal[float64]{10.0}}}, And}

	resExpr := &BinaryOperation{firstUpdate, &Literal[float64]{2.0}, Div}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsFourthPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "float64"}

	firstUpdate := &BinaryOperation{&Cast{a, "float64"}, b, Sub}
	firstIf := &Not{&Equals{&BinaryOperation{a, &Literal[int]{2}, Mod}, &Literal[int]{0}}}
	pc := &BinaryOperation{firstIf, &Not{&LT{firstUpdate, &Literal[float64]{10.0}}}, And}

	resExpr := &BinaryOperation{firstUpdate, &Literal[float64]{2.0}, Div}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestNestedConditionsFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "float64"}

	pc := &BinaryOperation{&LT{a, &Literal[int]{0}}, &LT{b, &Literal[float64]{0.0}}, And}

	resExpr := &BinaryOperation{&Cast{&BinaryOperation{a, &Literal[int]{-1}, Mul}, "float64"}, b, Add}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestNestedConditionsSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "float64"}

	pc := &BinaryOperation{&LT{a, &Literal[int]{0}}, &Not{&LT{b, &Literal[float64]{0.0}}}, And}

	resExpr := &BinaryOperation{&Cast{&BinaryOperation{a, &Literal[int]{-1}, Mul}, "float64"}, b, Sub}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestNestedConditionsThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "float64"}

	pc := &Not{&LT{a, &Literal[int]{0}}}

	resExpr := &BinaryOperation{&Cast{a, "float64"}, b, Add}

	checkResultWithPathCondition(t, pc, resExpr, true)
}

func TestBitwiseOperationsFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{&Equals{&BinaryOperation{a, &Literal[int]{1}, And},
		&Literal[int]{0}}, &Equals{&BinaryOperation{b, &Literal[int]{1}, And},
		&Literal[int]{0}}, And}

	resExpr := &BinaryOperation{a, b, Or}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestBitwiseOperationsSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{&Equals{&BinaryOperation{a, &Literal[int]{1}, And},
		&Literal[int]{1}}, &Equals{&BinaryOperation{b, &Literal[int]{1}, And},
		&Literal[int]{1}}, And}

	resExpr := &BinaryOperation{a, b, And}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestBitwiseOperationsThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{
		&Not{&BinaryOperation{&Equals{&BinaryOperation{a, &Literal[int]{1}, And},
			&Literal[int]{1}}, &Equals{&BinaryOperation{b, &Literal[int]{1}, And},
			&Literal[int]{1}}, And}},
		&Not{
			&BinaryOperation{&Equals{&BinaryOperation{a, &Literal[int]{1}, And},
				&Literal[int]{0}}, &Equals{&BinaryOperation{b, &Literal[int]{1}, And},
				&Literal[int]{0}}, And},
		},
		And,
	}

	resExpr := &BinaryOperation{a, b, Xor}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestAdvancedBitwiseFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &GT{a, b}

	resExpr := &BinaryOperation{a, &Literal[int]{1}, LeftShift}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestAdvancedBitwiseSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{&Not{&GT{a, b}}, &LT{a, b}, And}

	resExpr := &BinaryOperation{b, &Literal[int]{1}, RightShift}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestAdvancedBitwiseThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{&Not{&GT{a, b}}, &Not{&LT{a, b}}, And}

	resExpr := &BinaryOperation{a, b, Xor}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestCombinedBitwiseFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &Equals{&BinaryOperation{a, b, And}, &Literal[int]{0}}

	resExpr := &BinaryOperation{a, b, Or}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestCombinedBitwiseSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	firstIf := &Not{&Equals{&BinaryOperation{a, b, And}, &Literal[int]{0}}}

	resUpdate := &BinaryOperation{a, b, And}

	pc := &BinaryOperation{firstIf, &GT{resUpdate, &Literal[int]{10}}, And}

	resExpr := &BinaryOperation{resUpdate, b, Xor}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestCombinedBitwiseThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	firstIf := &Not{&Equals{&BinaryOperation{a, b, And}, &Literal[int]{0}}}

	resUpdate := &BinaryOperation{a, b, And}

	pc := &BinaryOperation{firstIf, &Not{&GT{resUpdate, &Literal[int]{10}}}, And}

	checkResultWithPathCondition(t, pc, resUpdate, false)
}

func TestNestedBitwiseFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}

	pc := &LT{a, &Literal[int]{0}}

	resExpr := &Literal[int]{-1}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestNestedBitwiseSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{&Not{&LT{a, &Literal[int]{0}}}, &LT{b, &Literal[int]{0}}, And}

	resExpr := &BinaryOperation{a, &Literal[int]{0}, Xor}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestNestedBitwiseThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{&BinaryOperation{&Not{&LT{a, &Literal[int]{0}}},
		&Not{&LT{b, &Literal[int]{0}}}, And},
		&Equals{&BinaryOperation{a, b, And}, &Literal[int]{0}}, And}

	resExpr := &BinaryOperation{a, b, Or}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func TestNestedBitwiseFourthPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	pc := &BinaryOperation{&BinaryOperation{&Not{&LT{a, &Literal[int]{0}}},
		&Not{&LT{b, &Literal[int]{0}}}, And},
		&Not{&Equals{&BinaryOperation{a, b, And}, &Literal[int]{0}}}, And}

	resExpr := &BinaryOperation{a, b, And}

	checkResultWithPathCondition(t, pc, resExpr, false)
}

func checkResultWithPathCondition(t *testing.T, pathCondition SymbolicExpression, resultExpression SymbolicExpression, isFloatResult bool) {
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
		t.Fatal("UNSAT")
	}
	if err != nil {
		t.Fatal(err)
	}

	t.Log(solver.SmtSolver.Model().String())
}
