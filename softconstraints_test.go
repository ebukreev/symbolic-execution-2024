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

	assumptions := []z3.AST{
		assumption1Ast,
		assumption2Ast,
	}

	sat := solver.SmtSolver.CheckAssumptions(assumptions)

	if sat {
		t.Fatal("FALSE SAT")
	}

	unsatCore := solver.SmtSolver.UnsatCore()
	newAssumptions := []z3.AST{}

Outer:
	for _, assumption := range assumptions {
		stringAssumption := assumption.String()
		for _, unsat := range unsatCore {
			if stringAssumption == unsat.String() {
				continue Outer
			}
		}
		newAssumptions = append(newAssumptions, assumption)
	}

	t.Log(newAssumptions)

	if len(newAssumptions) > 0 {
		sat = solver.SmtSolver.CheckAssumptions(newAssumptions)
	} else {
		sat, _ = solver.SmtSolver.Check()
	}

	if !sat {
		t.Fatal("FALSE UNSAT")
	}

	t.Log(solver.SmtSolver.Model())
}
