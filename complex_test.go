package main

import (
	"path"
	"runtime"
	"symbolic-execution-2024/z3"
	"testing"
)

func TestBasicComplexOperationsFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}
	realA := &FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}
	realB := &FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}
	pc := &GT{realA, realB}

	checkComplexResultAndPathCondition(t, pc, &BinaryOperation{a, b, Add}, false)
}

func TestBasicComplexOperationsSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}
	pc := &BinaryOperation{&Not{&GT{Left: &FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
		Right: &FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}}},
		&GT{Left: &FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			Right: &FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}}},
		And}

	checkComplexResultAndPathCondition(t, pc, &BinaryOperation{a, b, Sub}, false)
}

func TestBasicComplexOperationsThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}
	pc := &BinaryOperation{&Not{&GT{Left: &FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
		Right: &FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}}},
		&Not{&GT{Left: &FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			Right: &FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}}}},
		And}

	checkComplexResultAndPathCondition(t, pc, &BinaryOperation{a, b, Mul}, false)
}

func TestBasicComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	conditional := AnalyseStatically(file, "basicComplexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestBasicComplexOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	results := AnalyseDynamically(file, "basicComplexOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestComplexMagnitude(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	pc := &Literal[bool]{true}

	expression := &BinaryOperation{
		&BinaryOperation{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, Mul},
		&BinaryOperation{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}}, Mul},
		Add,
	}

	checkComplexResultAndPathCondition(t, pc, expression, true)
}

func TestComplexMagnitudeInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	conditional := AnalyseStatically(file, "complexMagnitude")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, true)
	}
}

func TestComplexMagnitudeDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	results := AnalyseDynamically(file, "complexMagnitude")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, true)
	}
}

func TestComplexComparisonFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	magA := &BinaryOperation{
		&BinaryOperation{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, Mul},
		&BinaryOperation{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}}, Mul},
		Add,
	}

	magB := &BinaryOperation{
		&BinaryOperation{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}, Mul},
		&BinaryOperation{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}},
			&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}}, Mul},
		Add,
	}

	pc := &GT{magA, magB}

	checkComplexResultAndPathCondition(t, pc, &Literal[float64]{1.0}, true)
}

func TestComplexComparisonSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	magA := &BinaryOperation{
		&BinaryOperation{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, Mul},
		&BinaryOperation{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}}, Mul},
		Add,
	}

	magB := &BinaryOperation{
		&BinaryOperation{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}, Mul},
		&BinaryOperation{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}},
			&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}}, Mul},
		Add,
	}

	pc := &BinaryOperation{&Not{&GT{magA, magB}}, &LT{magA, magB}, And}

	checkComplexResultAndPathCondition(t, pc, &Literal[float64]{2.0}, true)
}

func TestComplexComparisonThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	magA := &BinaryOperation{
		&BinaryOperation{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, Mul},
		&BinaryOperation{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}}, Mul},
		Add,
	}

	magB := &BinaryOperation{
		&BinaryOperation{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}, Mul},
		&BinaryOperation{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}},
			&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}}, Mul},
		Add,
	}

	pc := &BinaryOperation{&Not{&GT{magA, magB}}, &Not{&LT{magA, magB}}, And}

	checkComplexResultAndPathCondition(t, pc, &Literal[float64]{2.0}, true)
}

func TestComplexComparisonDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	results := AnalyseDynamically(file, "complexComparison")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CallStack[0].ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CallStack[0].ReturnValue, true)
	}
}

func TestComplexOperationsFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&Equals{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
		&Literal[float64]{0.0}},
		&Equals{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&Literal[float64]{0.0}}, And}

	expr := b

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&Not{&BinaryOperation{&Equals{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
		&Literal[float64]{0.0}},
		&Equals{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&Literal[float64]{0.0}}, And}},
		&BinaryOperation{&Equals{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}},
			&Literal[float64]{0.0}},
			&Equals{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}},
				&Literal[float64]{0.0}}, And},
		And}

	expr := a

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&BinaryOperation{&Not{&BinaryOperation{&Equals{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
		&Literal[float64]{0.0}},
		&Equals{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&Literal[float64]{0.0}}, And}},
		&Not{&BinaryOperation{&Equals{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}},
			&Literal[float64]{0.0}},
			&Equals{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}},
				&Literal[float64]{0.0}}, And}},
		And},
		&GT{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}},
		And}

	expr := &BinaryOperation{a, b, Div}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsFourthPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&BinaryOperation{&Not{&BinaryOperation{&Equals{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
		&Literal[float64]{0.0}},
		&Equals{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}},
			&Literal[float64]{0.0}}, And}},
		&Not{&BinaryOperation{&Equals{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}},
			&Literal[float64]{0.0}},
			&Equals{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}},
				&Literal[float64]{0.0}}, And}},
		And},
		&Not{&GT{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}},
			&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{b}}}},
		And}

	expr := &BinaryOperation{a, b, Add}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	conditional := AnalyseStatically(file, "complexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestComplexOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	results := AnalyseDynamically(file, "complexOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CallStack[0].ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CallStack[0].ReturnValue, false)
	}
}

func TestNestedComplexOperationsFirstPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&LT{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, &Literal[float64]{0.0}},
		&LT{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}}, &Literal[float64]{0.0}},
		And}

	expr := &BinaryOperation{a, b, Mul}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsSecondPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&LT{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, &Literal[float64]{0.0}},
		&Not{&LT{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{a}}, &Literal[float64]{0.0}}},
		And}

	expr := &BinaryOperation{a, b, Add}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsThirdPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&Not{&LT{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, &Literal[float64]{0.0}}},
		&LT{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}}, &Literal[float64]{0.0}},
		And}

	expr := &BinaryOperation{a, b, Sub}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsFourthPath(t *testing.T) {
	a := &InputValue{Name: "a", Type: "complex128"}
	b := &InputValue{Name: "b", Type: "complex128"}

	pc := &BinaryOperation{&Not{&LT{&FunctionCall{"builtin_real(ComplexType)", []SymbolicExpression{a}}, &Literal[float64]{0.0}}},
		&Not{&LT{&FunctionCall{"builtin_imag(ComplexType)", []SymbolicExpression{b}}, &Literal[float64]{0.0}}},
		And}

	expr := &BinaryOperation{a, b, Add}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	conditional := AnalyseStatically(file, "nestedComplexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestNestedComplexOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "complex.go")

	results := AnalyseDynamically(file, "nestedComplexOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CallStack[0].ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CallStack[0].ReturnValue, false)
	}
}

func checkComplexResultAndPathCondition(t *testing.T, pathCondition SymbolicExpression, resultExpression SymbolicExpression, isFloatResult bool) {
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
