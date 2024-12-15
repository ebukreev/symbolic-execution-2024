package se

import (
	"path"
	"runtime"
	"symbolic-execution-2024"
	"testing"
)

func TestCompareElementFirstPath(t *testing.T) {
	array := &se.InputValue{Name: "array", Type: "[]int"}
	index := &se.InputValue{Name: "index", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{array}}

	pc := &se.BinaryOperation{
		&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
		&se.Equals{index, arrLen},
		se.Or,
	}

	expression := &se.Literal[int]{-1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareElementSecondPath(t *testing.T) {
	array := &se.InputValue{Name: "array", Type: "[]int"}
	index := &se.InputValue{Name: "index", Type: "int"}
	value := &se.InputValue{Name: "value", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{array}}

	pc := &se.Not{&se.BinaryOperation{
		&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
		&se.Equals{index, arrLen},
		se.Or,
	}}

	element := &se.ArrayAccess{array, index}

	newPc := &se.BinaryOperation{pc, &se.GT{element, value}, se.And}

	checkResultWithPathCondition(t, newPc, &se.Literal[int]{1}, false)
}

func TestCompareElementThirdPath(t *testing.T) {
	array := &se.InputValue{Name: "array", Type: "[]int"}
	index := &se.InputValue{Name: "index", Type: "int"}
	value := &se.InputValue{Name: "value", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{array}}

	pc := &se.Not{&se.BinaryOperation{
		&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
		&se.Equals{index, arrLen},
		se.Or,
	}}

	element := &se.ArrayAccess{array, index}

	newPc := &se.BinaryOperation{&se.BinaryOperation{pc, &se.Not{&se.GT{element, value}}, se.And},
		&se.LT{element, value}, se.And}

	checkResultWithPathCondition(t, newPc, &se.Literal[int]{-1}, false)
}

func TestCompareElementFourthPath(t *testing.T) {
	array := &se.InputValue{Name: "array", Type: "[]int"}
	index := &se.InputValue{Name: "index", Type: "int"}
	value := &se.InputValue{Name: "value", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{array}}

	pc := &se.Not{&se.BinaryOperation{
		&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
		&se.Equals{index, arrLen},
		se.Or,
	}}

	element := &se.ArrayAccess{array, index}

	newPc := &se.BinaryOperation{&se.BinaryOperation{pc, &se.Not{&se.GT{element, value}}, se.And},
		&se.Not{&se.LT{element, value}}, se.And}

	checkResultWithPathCondition(t, newPc, &se.Literal[int]{0}, false)
}

func TestCompareElementInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "arrays.go")

	conditional := se.AnalyseStatically(file, "compareElement")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkResultWithPathCondition(t, cond, value, false)
	}
}

func TestCompareElementDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "arrays.go")

	results := se.AnalyseDynamically(file, "compareElement")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		checkResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestCompareAgeFirstPath(t *testing.T) {
	people := &se.InputValue{Name: "people", Type: "[]{string,int}"}
	index := &se.InputValue{Name: "index", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{people}}

	pc := &se.BinaryOperation{
		&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
		&se.Equals{index, arrLen},
		se.Or,
	}

	expression := &se.Literal[int]{-1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeSecondPath(t *testing.T) {
	people := &se.InputValue{Name: "people", Type: "[]{string,int}"}
	index := &se.InputValue{Name: "index", Type: "int"}
	value := &se.InputValue{Name: "index", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{people}}
	age := &se.FunctionCall{"{string,int}_1", []se.SymbolicExpression{&se.ArrayAccess{people, index}}}

	pc := &se.BinaryOperation{&se.Not{&se.BinaryOperation{
		&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
		&se.Equals{index, arrLen},
		se.Or,
	}},
		&se.GT{age, value},
		se.And,
	}

	expression := &se.Literal[int]{1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeThirdPath(t *testing.T) {
	people := &se.InputValue{Name: "people", Type: "[]{string,int}"}
	index := &se.InputValue{Name: "index", Type: "int"}
	value := &se.InputValue{Name: "index", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{people}}
	age := &se.FunctionCall{"{string,int}_1", []se.SymbolicExpression{&se.ArrayAccess{people, index}}}

	pc := &se.BinaryOperation{
		&se.BinaryOperation{&se.Not{&se.BinaryOperation{
			&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
			&se.Equals{index, arrLen},
			se.Or,
		}},
			&se.Not{&se.GT{age, value}},
			se.And,
		},
		&se.LT{age, value},
		se.And,
	}

	expression := &se.Literal[int]{-1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeFourthPath(t *testing.T) {
	people := &se.InputValue{Name: "people", Type: "[]{string,int}"}
	index := &se.InputValue{Name: "index", Type: "int"}
	value := &se.InputValue{Name: "index", Type: "int"}

	arrLen := &se.FunctionCall{"builtin_len(Type)", []se.SymbolicExpression{people}}
	age := &se.FunctionCall{"{string,int}_1", []se.SymbolicExpression{&se.ArrayAccess{people, index}}}

	pc := &se.BinaryOperation{
		&se.BinaryOperation{&se.Not{&se.BinaryOperation{
			&se.BinaryOperation{&se.LT{index, &se.Literal[int]{0}}, &se.GT{index, arrLen}, se.Or},
			&se.Equals{index, arrLen},
			se.Or,
		}},
			&se.Not{&se.GT{age, value}},
			se.And,
		},
		&se.Not{&se.LT{age, value}},
		se.And,
	}

	expression := &se.Literal[int]{0}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "arrays.go")

	conditional := se.AnalyseStatically(file, "compareAge")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkResultWithPathCondition(t, cond, value, false)
	}
}

func TestCompareAgeDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "constraints", "arrays.go")

	results := se.AnalyseDynamically(file, "compareAge")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		checkResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}
