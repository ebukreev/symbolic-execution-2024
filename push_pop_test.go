package main

import (
	"github.com/aclements/go-z3/z3"
	"testing"
)

func TestPushPopIncrementality(t *testing.T) {
	j := &InputValue{Name: "j", Type: "int"}
	result := &InputValue{Name: "res", Type: "int"}

	solver := CreateSolver(false)
	smtBuilder := SmtBuilder{Context: solver.Context}

	mod := &Equals{&BinaryOperation{result, &Literal[int]{2}, Mod}, &Literal[int]{0}}
	solver.SmtSolver.Assert(smtBuilder.BuildSmt(mod)[0].(z3.Bool))

	for i := 1; i <= 10; i++ {
		solver.SmtSolver.Push()
		valueAssertion := &Equals{result, &BinaryOperation{&Literal[int]{i * (i + 1) / 2}, j, Add}}
		solver.SmtSolver.Assert(smtBuilder.BuildSmt(valueAssertion)[0].(z3.Bool))

		sat, err := solver.SmtSolver.Check()
		if !sat {
			t.Log("UNSAT")
			solver.SmtSolver.Pop()
			continue
		}
		if err != nil {
			t.Fatal(err)
		}

		t.Log(solver.SmtSolver.Model().String())
		solver.SmtSolver.Pop()
	}
}
