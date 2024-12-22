package se

import (
	"path"
	"runtime"
	"symbolic-execution-2024"
	"testing"
)

func TestIntegerOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}
	pc := &se.GT{Left: a, Right: b}

	resExpr := &se.BinaryOperation{
		Left:  a,
		Right: b,
		Type:  se.Add,
	}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestIntegerOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{a, b}}, &se.LT{a, b}, se.And}

	resExpr := &se.BinaryOperation{
		Left:  a,
		Right: b,
		Type:  se.Sub,
	}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestIntegerOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{a, b}}, &se.Not{&se.LT{a, b}}, se.And}

	resExpr := &se.BinaryOperation{
		Left:  a,
		Right: b,
		Type:  se.Mul,
	}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestIntegerOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "integerOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, false)
	}
}

func TestIntegerOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "integerOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestFloatOperationsFirstPath(t *testing.T) {
	x := &se.InputValue{Name: "x", Type: "float64"}
	y := &se.InputValue{Name: "y", Type: "float64"}
	pc := &se.GT{Left: x, Right: y}

	resExpr := &se.BinaryOperation{
		Left:  x,
		Right: y,
		Type:  se.Div,
	}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestFloatOperationsSecondPath(t *testing.T) {
	x := &se.InputValue{Name: "x", Type: "float64"}
	y := &se.InputValue{Name: "y", Type: "float64"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{x, y}}, &se.LT{x, y}, se.And}

	resExpr := &se.BinaryOperation{
		Left:  x,
		Right: y,
		Type:  se.Div,
	}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestFloatOperationsThirdPath(t *testing.T) {
	x := &se.InputValue{Name: "x", Type: "float64"}
	y := &se.InputValue{Name: "y", Type: "float64"}
	pc := &se.BinaryOperation{&se.Not{&se.GT{x, y}}, &se.Not{&se.LT{x, y}}, se.And}

	resExpr := &se.Literal[float64]{0.0}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestFloatOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "floatOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, true)
	}
}

func TestFloatOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "floatOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, true)
	}
}

func TestMixedOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "float64"}

	firstUpdate := &se.BinaryOperation{&se.Cast{a, "float64"}, b, se.Add}
	firstIf := &se.Equals{&se.BinaryOperation{a, &se.Literal[int]{2}, se.Mod}, &se.Literal[int]{0}}
	pc := &se.BinaryOperation{firstIf, &se.LT{firstUpdate, &se.Literal[float64]{10.0}}, se.And}

	resExpr := &se.BinaryOperation{firstUpdate, &se.Literal[float64]{2.0}, se.Mul}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "float64"}

	firstUpdate := &se.BinaryOperation{&se.Cast{a, "float64"}, b, se.Sub}
	firstIf := &se.Not{&se.Equals{&se.BinaryOperation{a, &se.Literal[int]{2}, se.Mod}, &se.Literal[int]{0}}}
	pc := &se.BinaryOperation{firstIf, &se.LT{firstUpdate, &se.Literal[float64]{10.0}}, se.And}

	resExpr := &se.BinaryOperation{firstUpdate, &se.Literal[float64]{2.0}, se.Mul}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "float64"}

	firstUpdate := &se.BinaryOperation{&se.Cast{a, "float64"}, b, se.Add}
	firstIf := &se.Equals{&se.BinaryOperation{a, &se.Literal[int]{2}, se.Mod}, &se.Literal[int]{0}}
	pc := &se.BinaryOperation{firstIf, &se.Not{&se.LT{firstUpdate, &se.Literal[float64]{10.0}}}, se.And}

	resExpr := &se.BinaryOperation{firstUpdate, &se.Literal[float64]{2.0}, se.Div}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsFourthPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "float64"}

	firstUpdate := &se.BinaryOperation{&se.Cast{a, "float64"}, b, se.Sub}
	firstIf := &se.Not{&se.Equals{&se.BinaryOperation{a, &se.Literal[int]{2}, se.Mod}, &se.Literal[int]{0}}}
	pc := &se.BinaryOperation{firstIf, &se.Not{&se.LT{firstUpdate, &se.Literal[float64]{10.0}}}, se.And}

	resExpr := &se.BinaryOperation{firstUpdate, &se.Literal[float64]{2.0}, se.Div}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestMixedOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "mixedOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, true)
	}
}

func TestMixedOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "mixedOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, true)
	}
}

func TestNestedConditionsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "float64"}

	pc := &se.BinaryOperation{&se.LT{a, &se.Literal[int]{0}}, &se.LT{b, &se.Literal[float64]{0.0}}, se.And}

	resExpr := &se.BinaryOperation{&se.Cast{&se.BinaryOperation{a, &se.Literal[int]{-1}, se.Mul}, "float64"}, b, se.Add}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestNestedConditionsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "float64"}

	pc := &se.BinaryOperation{&se.LT{a, &se.Literal[int]{0}}, &se.Not{&se.LT{b, &se.Literal[float64]{0.0}}}, se.And}

	resExpr := &se.BinaryOperation{&se.Cast{&se.BinaryOperation{a, &se.Literal[int]{-1}, se.Mul}, "float64"}, b, se.Sub}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestNestedConditionsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "float64"}

	pc := &se.Not{&se.LT{a, &se.Literal[int]{0}}}

	resExpr := &se.BinaryOperation{&se.Cast{a, "float64"}, b, se.Add}

	se.CheckResultWithPathCondition(t, pc, resExpr, true)
}

func TestNestedConditionsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "nestedConditions")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, true)
	}
}

func TestNestedConditionsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "nestedConditions")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, true)
	}
}

func TestBitwiseOperationsFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{&se.Equals{&se.BinaryOperation{a, &se.Literal[int]{1}, se.And},
		&se.Literal[int]{0}}, &se.Equals{&se.BinaryOperation{b, &se.Literal[int]{1}, se.And},
		&se.Literal[int]{0}}, se.And}

	resExpr := &se.BinaryOperation{a, b, se.Or}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestBitwiseOperationsSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{&se.Equals{&se.BinaryOperation{a, &se.Literal[int]{1}, se.And},
		&se.Literal[int]{1}}, &se.Equals{&se.BinaryOperation{b, &se.Literal[int]{1}, se.And},
		&se.Literal[int]{1}}, se.And}

	resExpr := &se.BinaryOperation{a, b, se.And}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestBitwiseOperationsThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{
		&se.Not{&se.BinaryOperation{&se.Equals{&se.BinaryOperation{a, &se.Literal[int]{1}, se.And},
			&se.Literal[int]{1}}, &se.Equals{&se.BinaryOperation{b, &se.Literal[int]{1}, se.And},
			&se.Literal[int]{1}}, se.And}},
		&se.Not{
			&se.BinaryOperation{&se.Equals{&se.BinaryOperation{a, &se.Literal[int]{1}, se.And},
				&se.Literal[int]{0}}, &se.Equals{&se.BinaryOperation{b, &se.Literal[int]{1}, se.And},
				&se.Literal[int]{0}}, se.And},
		},
		se.And,
	}

	resExpr := &se.BinaryOperation{a, b, se.Xor}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestBitwiseOperationsInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "bitwiseOperations")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, false)
	}
}

func TestBitwiseOperationsDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "bitwiseOperations")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestAdvancedBitwiseFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.GT{a, b}

	resExpr := &se.BinaryOperation{a, &se.Literal[int]{1}, se.LeftShift}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestAdvancedBitwiseSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{&se.Not{&se.GT{a, b}}, &se.LT{a, b}, se.And}

	resExpr := &se.BinaryOperation{b, &se.Literal[int]{1}, se.RightShift}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestAdvancedBitwiseThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{&se.Not{&se.GT{a, b}}, &se.Not{&se.LT{a, b}}, se.And}

	resExpr := &se.BinaryOperation{a, b, se.Xor}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestAdvancedBitwiseInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "advancedBitwise")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, false)
	}
}

func TestAdvancedBitwiseDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "advancedBitwise")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestCombinedBitwiseFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.Equals{&se.BinaryOperation{a, b, se.And}, &se.Literal[int]{0}}

	resExpr := &se.BinaryOperation{a, b, se.Or}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestCombinedBitwiseSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	firstIf := &se.Not{&se.Equals{&se.BinaryOperation{a, b, se.And}, &se.Literal[int]{0}}}

	resUpdate := &se.BinaryOperation{a, b, se.And}

	pc := &se.BinaryOperation{firstIf, &se.GT{resUpdate, &se.Literal[int]{10}}, se.And}

	resExpr := &se.BinaryOperation{resUpdate, b, se.Xor}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestCombinedBitwiseThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	firstIf := &se.Not{&se.Equals{&se.BinaryOperation{a, b, se.And}, &se.Literal[int]{0}}}

	resUpdate := &se.BinaryOperation{a, b, se.And}

	pc := &se.BinaryOperation{firstIf, &se.Not{&se.GT{resUpdate, &se.Literal[int]{10}}}, se.And}

	se.CheckResultWithPathCondition(t, pc, resUpdate, false)
}

func TestCombinedBitwiseInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "combinedBitwise")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, false)
	}
}

func TestCombinedBitwiseDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "combinedBitwise")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestNestedBitwiseFirstPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}

	pc := &se.LT{a, &se.Literal[int]{0}}

	resExpr := &se.Literal[int]{-1}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestNestedBitwiseSecondPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{&se.Not{&se.LT{a, &se.Literal[int]{0}}}, &se.LT{b, &se.Literal[int]{0}}, se.And}

	resExpr := &se.BinaryOperation{a, &se.Literal[int]{0}, se.Xor}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestNestedBitwiseThirdPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{&se.BinaryOperation{&se.Not{&se.LT{a, &se.Literal[int]{0}}},
		&se.Not{&se.LT{b, &se.Literal[int]{0}}}, se.And},
		&se.Equals{&se.BinaryOperation{a, b, se.And}, &se.Literal[int]{0}}, se.And}

	resExpr := &se.BinaryOperation{a, b, se.Or}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestNestedBitwiseFourthPath(t *testing.T) {
	a := &se.InputValue{Name: "a", Type: "int"}
	b := &se.InputValue{Name: "b", Type: "int"}

	pc := &se.BinaryOperation{&se.BinaryOperation{&se.Not{&se.LT{a, &se.Literal[int]{0}}},
		&se.Not{&se.LT{b, &se.Literal[int]{0}}}, se.And},
		&se.Not{&se.Equals{&se.BinaryOperation{a, b, se.And}, &se.Literal[int]{0}}}, se.And}

	resExpr := &se.BinaryOperation{a, b, se.And}

	se.CheckResultWithPathCondition(t, pc, resExpr, false)
}

func TestNestedBitwiseInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	conditional := se.AnalyseStatically(file, "nestedBitwise")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		se.CheckResultWithPathCondition(t, cond, value, false)
	}
}

func TestNestedBitwiseDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "numbers.go")

	results := se.AnalyseDynamically(file, "nestedBitwise")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}
