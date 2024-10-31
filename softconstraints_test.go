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

	solver.SmtSolver.Assert(smtBuilder.BuildSmt(pc)[0].(z3.Bool))

	assumption1 := &LT{&BinaryOperation{b, a, Sub}, &Literal[int]{5}}
	assumption2 := &GT{&BinaryOperation{b, a, Sub}, &Literal[int]{0}}
	assumption1Ast := smtBuilder.BuildSmt(assumption1)[0].(z3.Bool).AsAST()
	assumption2Ast := smtBuilder.BuildSmt(assumption2)[0].(z3.Bool).AsAST()

	sat := solver.SmtSolver.CheckAssumptions([]z3.AST{
		assumption1Ast,
		assumption2Ast,
	})

	if sat {
		t.Fatal("FALSE SAT")
	}

	unsatCore := solver.SmtSolver.UnsatCore()
	t.Log(unsatCore)

	sat = solver.SmtSolver.CheckAssumptions([]z3.AST{
		assumption1Ast,
	})

	if !sat {
		t.Fatal("FALSE UNSAT")
	}
	t.Log(solver.SmtSolver.Model())

	sat = solver.SmtSolver.CheckAssumptions([]z3.AST{
		assumption2Ast,
	})

	if !sat {
		t.Fatal("FALSE UNSAT")
	}
	t.Log(solver.SmtSolver.Model())
}
