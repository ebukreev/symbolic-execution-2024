package se

import (
	"symbolic-execution-2024"
	"symbolic-execution-2024/z3"
	"testing"
)

func TestPushPopIncrementality(t *testing.T) {
	j := &se.InputValue{Name: "j", Type: "int"}
	result := &se.InputValue{Name: "res", Type: "int"}

	solver := se.CreateSolver(false)
	smtBuilder := se.SmtBuilder{Context: solver.Context}

	valueAssertion := &se.Equals{result, &se.BinaryOperation{&se.Literal[int]{10 * 11 / 2}, j, se.Add}}
	solver.SmtSolver.Assert(smtBuilder.BuildSmt(valueAssertion)[0].(z3.Bool))

	solver.SmtSolver.Push()

	mod := &se.Equals{&se.BinaryOperation{result, &se.Literal[int]{2}, se.Mod}, &se.Literal[int]{0}}
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

	solver.SmtSolver.Assert(smtBuilder.BuildSmt(&se.Not{mod})[0].(z3.Bool))
	sat, err = solver.SmtSolver.Check()
	if !sat {
		t.Fatal("UNSAT")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(solver.SmtSolver.Model().String())
}
