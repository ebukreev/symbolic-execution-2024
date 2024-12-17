package se

import "symbolic-execution-2024/z3"

type Solver struct {
	Context   *z3.Context
	SmtSolver *z3.Solver
}

func CreateSolver(unsatCore bool) *Solver {
	config := z3.NewContextConfig()
	config.SetBool("unsat_core", unsatCore)
	config.SetUint("timeout", 1000)
	context := z3.NewContext(config)

	return &Solver{Context: context, SmtSolver: z3.NewSolver(context)}
}
