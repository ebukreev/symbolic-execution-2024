package se

import (
	"path"
	"runtime"
	se "symbolic-execution-2024"
	"testing"
)

func TestDefaultIntValues(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "samples", "arrays", "primitiveArrays.go")

	results := se.AnalyseMethodDynamically(file, "PrimitiveArrays", "DefaultIntValues")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestShortSizeAndIndex(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "samples", "arrays", "primitiveArrays.go")

	results := se.AnalyseMethodDynamically(file, "PrimitiveArrays", "ShortSizeAndIndex")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}
