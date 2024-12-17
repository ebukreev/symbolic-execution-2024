package se

import (
	"path"
	"runtime"
	se "symbolic-execution-2024"
	"testing"
)

func TestFib(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "samples", "flow", "recursion.go")

	results := se.AnalyseMethodDynamically(file, "MathExamples", "Fib")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}
