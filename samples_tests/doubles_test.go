package se

import (
	"path"
	"runtime"
	se "symbolic-execution-2024"
	"testing"
)

func TestCompareSum(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "samples", "primitives", "doubles.go")

	results := se.AnalyseMethodDynamically(file, "DoubleExamples", "CompareSum")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, true)
	}
}
