package main

import (
	"symbolic-execution-2024/z3"
	"testing"
)

func TestCompareAndIncrement(t *testing.T) {
	a := &InputValue{Name: "a", Type: "int"}
	b := &InputValue{Name: "b", Type: "int"}

	solver := CreateSolver(true)
	smtBuilder := SmtBuilder{Context: solver.Context}

	pc := &BinaryOperation{&GT{a, b},
		&Not{&GT{&BinaryOperation{a, &Literal[int]{15}, Add}, b}}, And}
	assumption := &BinaryOperation{&LT{&BinaryOperation{b, a, Sub}, &Literal[int]{5}},
		&GT{&BinaryOperation{b, a, Sub}, &Literal[int]{0}},
		And}
	assumptionVar := solver.Context.BoolConst("p1")

	solver.SmtSolver.Assert(smtBuilder.BuildSmt(pc)[0].(z3.Bool))
	solver.SmtSolver.Push()
	solver.SmtSolver.AssertAndTrack(smtBuilder.BuildSmt(assumption)[0].(z3.Bool), assumptionVar)

	sat, err := solver.SmtSolver.Check()
	if sat {
		t.Log("FALSE SAT")
	}
	if err != nil {
		t.Fatal(err)
	}

	unsatCore := solver.SmtSolver.UnsatCore()
	t.Log(unsatCore)

	for _, val := range unsatCore {
		if val.String() == assumptionVar.String() {
			solver.SmtSolver.Pop()
			sat, err := solver.SmtSolver.Check()
			if !sat {
				t.Log("UNSAT")
			}
			if err != nil {
				t.Fatal(err)
			}
			t.Log(solver.SmtSolver.Model().String())
			return
		}
	}
	t.Fatal("No assumption in unsat core")
}
