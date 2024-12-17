package se

import (
	"path"
	"runtime"
	se "symbolic-execution-2024"
	"testing"
)

func TestLoopWithSymbolicBoundAndComplexControlFlow(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "samples", "flow", "loops.go")

	results := se.AnalyseMethodDynamically(file, "Loops", "LoopWithSymbolicBoundAndComplexControlFlow")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}
