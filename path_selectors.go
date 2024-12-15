package se

import "math/rand"

type PathSelector interface {
	CalculatePriority(interpreter DynamicInterpreter) int
}

type DfsPathSelector struct {
	counter int // int.min_value
}

func (dfs *DfsPathSelector) CalculatePriority(interpreter DynamicInterpreter) int {
	dfs.counter++
	return dfs.counter
}

type BfsPathSelector struct {
	counter int // int.max_value
}

func (bfs *BfsPathSelector) CalculatePriority(interpreter DynamicInterpreter) int {
	bfs.counter--
	return bfs.counter
}

type RandomPathSelector struct{}

func (random *RandomPathSelector) CalculatePriority(interpreter DynamicInterpreter) int {
	return rand.Int()
}

type NursPathSelector struct{}

func (nurs *NursPathSelector) CalculatePriority(interpreter DynamicInterpreter) int {
	return len(interpreter.CurrentFrame().BlocksStack) + rand.Intn(5)
}
