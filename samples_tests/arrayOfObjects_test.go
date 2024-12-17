package se

import (
	"path"
	"runtime"
	se "symbolic-execution-2024"
	"testing"
)

func TestCreateArray(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "samples", "arrays", "arrayOfObjects.go")

	results := se.AnalyseMethodDynamically(file, "ArrayOfObjects", "CreateArray")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}

func TestObjectArray(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(path.Dir(filename)), "samples", "arrays", "arrayOfObjects.go")

	results := se.AnalyseMethodDynamically(file, "ArrayOfObjects", "ObjectArray")

	for _, result := range results {
		t.Log(result.PathCondition.String() + " => " + result.CurrentFrame().ReturnValue.String())
		se.CheckResultWithPathCondition(t, result.PathCondition, result.CurrentFrame().ReturnValue, false)
	}
}
