package main

import "math/rand"

type PathSelector interface {
	CalculatePriority(interpreter DynamicInterpreter) int
}

type DfsPathSelector struct{}

func (dfs DfsPathSelector) CalculatePriority(interpreter DynamicInterpreter) int {
	return 1
}

type RandomPathSelector struct{}

func (random RandomPathSelector) CalculatePriority(interpreter DynamicInterpreter) int {
	return rand.Int()
}
