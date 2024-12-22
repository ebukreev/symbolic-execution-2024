package se

import (
	"path"
	"runtime"
	"symbolic-execution-2024"
	"testing"
)

func TestBasicComplexOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}
	realA := &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}
	realB := &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}
	pc := &se.GT{realA, realB}

	se.CheckComplexResultAndPathCondition(t, pc, &se.BinaryOperation{a, b, se.Add}, false)
}

func TestBasicComplexOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{Left: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		Right: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}}},
		&se.GT{Left: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			Right: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}},
		se.And}

	se.CheckComplexResultAndPathCondition(t, pc, &se.BinaryOperation{a, b, se.Sub}, false)
}

func TestBasicComplexOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{Left: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		Right: &se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{b}}}},
		&se.Not{&se.GT{Left: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			Right: &se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}}},
		se.And}

	se.CheckComplexResultAndPathCondition(t, pc, &se.BinaryOperation{a, b, se.Mul}, false)
}

func TestBasicComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "basicComplexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestBasicComplexOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	results := se.AnalyseDynamically(file, "basicComplexOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckComplexResultAndPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
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

	se.CheckComplexResultAndPathCondition(t, pc, expression, true)
}

func TestComplexMagnitudeInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "complexMagnitude")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckComplexResultAndPathCondition(t, cond, value, true)
	}
}

func TestComplexMagnitudeDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	results := se.AnalyseDynamically(file, "complexMagnitude")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckComplexResultAndPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, true)
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

	se.CheckComplexResultAndPathCondition(t, pc, &se.Literal[float64]{1.0}, true)
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

	se.CheckComplexResultAndPathCondition(t, pc, &se.Literal[float64]{2.0}, true)
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

	se.CheckComplexResultAndPathCondition(t, pc, &se.Literal[float64]{2.0}, true)
}

func TestComplexOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.Equals{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}},
		&se.Literal[float64]{0.0}},
		&se.Equals{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}},
			&se.Literal[float64]{0.0}}, se.And}

	expr := b

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
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

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
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

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
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

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
}

func TestComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "complexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckComplexResultAndPathCondition(t, cond, value, false)
	}
}

func TestNestedComplexOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}},
		&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Mul}

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}},
		&se.Not{&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Add}

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.Not{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}}},
		&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}, &se.Literal[float64]{0.0}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Sub}

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsFourthPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "complex128"}
	b := &se.InputValue{Name: "b", Type: "complex128"}

	pc := &se.BinaryOperation{&se.Not{&se.LT{&se.FunctionCall{"builtin_real(ComplexType)", []se.SymbolicExpression{a}}, &se.Literal[float64]{0.0}}},
		&se.Not{&se.LT{&se.FunctionCall{"builtin_imag(ComplexType)", []se.SymbolicExpression{b}}, &se.Literal[float64]{0.0}}},
		se.And}

	expr := &se.BinaryOperation{a, b, se.Add}

	se.CheckComplexResultAndPathCondition(t, pc, expr, false)
}

func TestNestedComplexOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "complex.go")

	conditional := se.AnalyseStatically(file, "nestedComplexOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckComplexResultAndPathCondition(t, cond, value, false)
	}
}
