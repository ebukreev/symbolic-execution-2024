package main

import (
	"math"
	"testing"
)

func TestIntegerOperations0(t *testing.T) {
	callResult := integerOperations(13217643890510138, 13217643890510138)
	returnValue := 6778761081452962084
	if callResult != returnValue {
		t.Fatalf("%v != %v", callResult, returnValue)
	}
}

func TestIntegerOperations1(t *testing.T) {
	callResult := integerOperations(-561489442378677270, -416646377605233686)
	returnValue := -144843064773443584
	if callResult != returnValue {
		t.Fatalf("%v != %v", callResult, returnValue)
	}
}

func TestIntegerOperations2(t *testing.T) {
	callResult := integerOperations(-9223372036854775807, -9223372036854775808)
	returnValue := 1
	if callResult != returnValue {
		t.Fatalf("%v != %v", callResult, returnValue)
	}
}

func TestFloatOperations0(t *testing.T) {
	callResult := floatOperations(math.NaN(), math.NaN())
	returnValue := 0.000000
	if callResult != returnValue {
		t.Fatalf("%v != %v", callResult, returnValue)
	}
}

func TestFloatOperations1(t *testing.T) {
	callResult := floatOperations(0.000000, 0.000000)
	returnValue := 0.000000
	if callResult != returnValue {
		t.Fatalf("%v != %v", callResult, returnValue)
	}
}

func TestFloatOperations2(t *testing.T) {
	callResult := floatOperations(0.000000, 8000000000106039.000000)
	returnValue := 0.000000
	if callResult != returnValue {
		t.Fatalf("%v != %v", callResult, returnValue)
	}
}
