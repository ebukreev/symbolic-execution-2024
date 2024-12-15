package se

import (
	"path"
	"runtime"
	"symbolic-execution-2024"
	"symbolic-execution-2024/z3"
	"testing"
)

func TestBasicComplexOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}
	realA := &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}
	realB := &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}
	pc := &se.GT{realA, realB}

	checkComplexResultAndPathCondition(t, pc, &se.BinaryOperation{a, b, se.Add}, false)
}

func TestBasicComplexOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{Left: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		Right: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}}},
		&se.GT{Left: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			Right: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}},
		se.And}

	checkComplexResultAndPathCondition(t, pc, &se.BinaryOperation{a, b, se.Sub}, false)
}

func TestBasicComplexOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{Left: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		Right: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}}},
		&se.Not{&se.GT{Left: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			Right: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}}},
		se.And}

	checkComplexResultAndPathCondition(t, pc, &se.BinaryOperation{a, b, se.Mul}, false)
}

func TestBasicComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "basicComplexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestBasicComplexOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	results := se.AnalyseDynamically(file, "basicComplexOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestComplexMagnitude(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	pc := &se.Literal[bool]{true}

	expression := &se.BinaryOperation{
		&se.BinaryOperation{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		&se.BinaryOperation{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		se.Add,
	}

	checkComplexResultAndPathCondition(t, pc, expression, true)
}

func TestComplexMagnitudeInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "complexMagnitude")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, true)
	}
}

func TestComplexMagnitudeDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	results := se.AnalyseDynamically(file, "complexMagnitude")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, true)
	}
}

func TestComplexComparisonFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	magA := &se.BinaryOperation{
		&se.BinaryOperation{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		&se.BinaryOperation{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		se.Add,
	}

	magB := &se.BinaryOperation{
		&se.BinaryOperation{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}, se.Mul},
		&se.BinaryOperation{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}},
			&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}, se.Mul},
		se.Add,
	}

	pc := &se.GT{magA, magB}

	checkComplexResultAndPathCondition(t, pc, &se.Literal[float64]{1.0}, true)
}

func TestComplexComparisonSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	magA := &se.BinaryOperation{
		&se.BinaryOperation{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		&se.BinaryOperation{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		se.Add,
	}

	magB := &se.BinaryOperation{
		&se.BinaryOperation{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}, se.Mul},
		&se.BinaryOperation{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}},
			&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}, se.Mul},
		se.Add,
	}

	pc := &se.BinaryOperation{&se.Not{&se.GT{magA, magB}}, &se.LT{magA, magB}, se.And}

	checkComplexResultAndPathCondition(t, pc, &se.Literal[float64]{2.0}, true)
}

func TestComplexComparisonThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	magA := &se.BinaryOperation{
		&se.BinaryOperation{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		&se.BinaryOperation{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, se.Mul},
		se.Add,
	}

	magB := &se.BinaryOperation{
		&se.BinaryOperation{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}, se.Mul},
		&se.BinaryOperation{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}},
			&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}, se.Mul},
		se.Add,
	}

	pc := &se.BinaryOperation{&se.Not{&se.GT{magA, magB}}, &se.Not{&se.LT{magA, magB}}, se.And}

	checkComplexResultAndPathCondition(t, pc, &se.Literal[float64]{2.0}, true)
}

func TestComplexComparisonDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	results := se.AnalyseDynamically(file, "complexComparison")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CallStack[0].ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CallStack[0].ReturnValue, true)
	}
}

func TestComplexOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		&se.Literal[float64]{0.0}},
		&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.Literal[float64]{0.0}}, se.And}

	expr := b

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.Not{&se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		&se.Literal[float64]{0.0}},
		&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.Literal[float64]{0.0}}, se.And}},
		&se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}},
			&se.Literal[float64]{0.0}},
			&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}},
				&se.Literal[float64]{0.0}}, se.And},
		se.And}

	expr := a

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.BinaryOperation{&se.Not{&se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		&se.Literal[float64]{0.0}},
		&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.Literal[float64]{0.0}}, se.And}},
		&se.Not{&se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}},
			&se.Literal[float64]{0.0}},
			&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}},
				&se.Literal[float64]{0.0}}, se.And}},
		se.And},
		&se.GT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Div}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsFourthPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.BinaryOperation{&se.Not{&se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		&se.Literal[float64]{0.0}},
		&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.Literal[float64]{0.0}}, se.And}},
		&se.Not{&se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}},
			&se.Literal[float64]{0.0}},
			&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}},
				&se.Literal[float64]{0.0}}, se.And}},
		se.And},
		&se.Not{&se.GT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
			&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Add}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "complexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestComplexOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	results := se.AnalyseDynamically(file, "complexOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CallStack[0].ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CallStack[0].ReturnValue, false)
	}
}

func TestNestedComplexOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}},
		&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Mul}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}},
		&se.Not{&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Add}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.Not{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}}},
		&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}, &se.Literal[float64]{0.0}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Sub}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsFourthPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.Not{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}}},
		&se.Not{&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}, &se.Literal[float64]{0.0}}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Add}

	checkComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "nestedComplexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestNestedComplexOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	results := se.AnalyseDynamically(file, "nestedComplexOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CallStack[0].ReturnValue.String())
		checkComplexResultAndPathCondition(t, result.PathCondition, result.CallStack[0].ReturnValue, false)
	}
}

func checkComplexResultAndPathCondition(t *testing.T, pathCondition se.SymbolicExpression, resultExpression se.SymbolicExpression, isFloatResult bool) {
	solver := se.CreateSolver(false)
	smtBuilder := se.SmtBuilder{Context: solver.Context}

	builtSmt := smtBuilder.BuildSmt(pathCondition)
	for i := 0; i < len(builtSmt); i++ {
		solver.SmtSolver.Assert(builtSmt[i].(z3.Bool))
	}

	if isFloatResult {
		resultSymbolicVar := &se.InputValue{Name: "res", Type: "float64"}

		res := smtBuilder.BuildSmt(resultSymbolicVar)[0]

		expressionSmt := smtBuilder.BuildSmt(resultExpression)
		solver.SmtSolver.Assert(res.(z3.Float).Eq(expressionSmt[0].(z3.Float)))
	} else {
		realResultSymbolicVar := &se.InputValue{Name: "$R_res", Type: "float64"}
		imaginaryResultSymbolicVar := &se.InputValue{Name: "$I_res", Type: "float64"}

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
