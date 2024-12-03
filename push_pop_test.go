package main

import (
	"path"
	"runtime"
	"symbolic-execution-2024/z3"
	"testing"
)

func TestPushPopIncrementality(t *testing.T) {
	j := &InputValue{Name: "j", Type: "int"}
	result := &InputValue{Name: "res", Type: "int"}

	solver := CreateSolver(false)
	smtBuilder := SmtBuilder{Context: solver.Context}

	valueAssertion := &Equals{result, &BinaryOperation{&Literal[int]{10 * 11 / 2}, j, Add}}
	solver.SmtSolver.Assert(smtBuilder.BuildSmt(valueAssertion)[0].(z3.Bool))

	solver.SmtSolver.Push()

	mod := &Equals{&BinaryOperation{result, &Literal[int]{2}, Mod}, &Literal[int]{0}}
	solver.SmtSolver.Assert(smtBuilder.BuildSmt(mod)[0].(z3.Bool))

	sat, err := solver.SmtSolver.Check()
	if !sat {
		t.Fatal("UNSAT")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(solver.SmtSolver.Model().String())

	solver.SmtSolver.Pop()

	solver.SmtSolver.Assert(smtBuilder.BuildSmt(&Not{mod})[0].(z3.Bool))
	sat, err = solver.SmtSolver.Check()
	if !sat {
		t.Fatal("UNSAT")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(solver.SmtSolver.Model().String())
}

func TestPushPopDynamicInterpretation(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filename), "constraints", "push_pop.go")

	results := AnalyseDynamically(file, "pushPopIncrementality")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		checkResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}
