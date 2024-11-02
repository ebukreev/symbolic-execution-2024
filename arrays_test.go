package main

import (
	"path"
	"runtime"
	"testing"
)

func TestCompareElementFirstPath(t *testing.T) {
	array := &InputValue{Name: "array", Type: "[]int"}
	index := &InputValue{Name: "index", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{array}}

	pc := &BinaryOperation{
		&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
		&Equals{index, arrLen},
		Or,
	}

	expression := &Literal[int]{-1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareElementSecondPath(t *testing.T) {
	array := &InputValue{Name: "array", Type: "[]int"}
	index := &InputValue{Name: "index", Type: "int"}
	value := &InputValue{Name: "value", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{array}}

	pc := &Not{&BinaryOperation{
		&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
		&Equals{index, arrLen},
		Or,
	}}

	element := &ArrayAccess{array, index}

	newPc := &BinaryOperation{pc, &GT{element, value}, And}

	checkResultWithPathCondition(t, newPc, &Literal[int]{1}, false)
}

func TestCompareElementThirdPath(t *testing.T) {
	array := &InputValue{Name: "array", Type: "[]int"}
	index := &InputValue{Name: "index", Type: "int"}
	value := &InputValue{Name: "value", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{array}}

	pc := &Not{&BinaryOperation{
		&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
		&Equals{index, arrLen},
		Or,
	}}

	element := &ArrayAccess{array, index}

	newPc := &BinaryOperation{&BinaryOperation{pc, &Not{&GT{element, value}}, And},
		&LT{element, value}, And}

	checkResultWithPathCondition(t, newPc, &Literal[int]{-1}, false)
}

func TestCompareElementFourthPath(t *testing.T) {
	array := &InputValue{Name: "array", Type: "[]int"}
	index := &InputValue{Name: "index", Type: "int"}
	value := &InputValue{Name: "value", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{array}}

	pc := &Not{&BinaryOperation{
		&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
		&Equals{index, arrLen},
		Or,
	}}

	element := &ArrayAccess{array, index}

	newPc := &BinaryOperation{&BinaryOperation{pc, &Not{&GT{element, value}}, And},
		&Not{&LT{element, value}}, And}

	checkResultWithPathCondition(t, newPc, &Literal[int]{0}, false)
}

func TestCompareElementInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "arrays.go")

	conditional := Analyse(file, "compareElement")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkResultWithPathCondition(t, cond, value, false)
	}
}

func TestCompareAgeFirstPath(t *testing.T) {
	people := &InputValue{Name: "people", Type: "[]{string,int}"}
	index := &InputValue{Name: "index", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{people}}

	pc := &BinaryOperation{
		&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
		&Equals{index, arrLen},
		Or,
	}

	expression := &Literal[int]{-1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeSecondPath(t *testing.T) {
	people := &InputValue{Name: "people", Type: "[]{string,int}"}
	index := &InputValue{Name: "index", Type: "int"}
	value := &InputValue{Name: "index", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{people}}
	age := &FunctionCall{"{string,int}_1", []SymbolicExpression{&ArrayAccess{people, index}}}

	pc := &BinaryOperation{&Not{&BinaryOperation{
		&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
		&Equals{index, arrLen},
		Or,
	}},
		&GT{age, value},
		And,
	}

	expression := &Literal[int]{1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeThirdPath(t *testing.T) {
	people := &InputValue{Name: "people", Type: "[]{string,int}"}
	index := &InputValue{Name: "index", Type: "int"}
	value := &InputValue{Name: "index", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{people}}
	age := &FunctionCall{"{string,int}_1", []SymbolicExpression{&ArrayAccess{people, index}}}

	pc := &BinaryOperation{
		&BinaryOperation{&Not{&BinaryOperation{
			&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
			&Equals{index, arrLen},
			Or,
		}},
			&Not{&GT{age, value}},
			And,
		},
		&LT{age, value},
		And,
	}

	expression := &Literal[int]{-1}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeFourthPath(t *testing.T) {
	people := &InputValue{Name: "people", Type: "[]{string,int}"}
	index := &InputValue{Name: "index", Type: "int"}
	value := &InputValue{Name: "index", Type: "int"}

	arrLen := &FunctionCall{"builtin_len(Type)", []SymbolicExpression{people}}
	age := &FunctionCall{"{string,int}_1", []SymbolicExpression{&ArrayAccess{people, index}}}

	pc := &BinaryOperation{
		&BinaryOperation{&Not{&BinaryOperation{
			&BinaryOperation{&LT{index, &Literal[int]{0}}, &GT{index, arrLen}, Or},
			&Equals{index, arrLen},
			Or,
		}},
			&Not{&GT{age, value}},
			And,
		},
		&Not{&LT{age, value}},
		And,
	}

	expression := &Literal[int]{0}

	checkResultWithPathCondition(t, pc, expression, false)
}

func TestCompareAgeInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "arrays.go")

	conditional := Analyse(file, "compareAge")

	t.Log((&conditional).String())

	for cond, value := range conditional.Options {
		checkResultWithPathCondition(t, cond, value, false)
	}
}
