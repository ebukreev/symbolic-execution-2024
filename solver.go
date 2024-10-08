package main

import "github.com/aclements/go-z3/z3"

type Solver struct {
	Context   *z3.Context
	SmtSolver *z3.Solver
}

func CreateSolver() *Solver {
	config := z3.NewContextConfig()
	context := z3.NewContext(config)

	return &Solver{Context: context, SmtSolver: z3.NewSolver(context)}
}
