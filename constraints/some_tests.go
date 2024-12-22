package main

type Test struct{}

func (a *Test) AllocArray1() int {
	arr := []int{1, 2, 3}
	if arr[1] == 2 {
		return 3
	}

	return 5
}

func (a *Test) AllocArray2() int {
	arr := []int{1, 2, 3}
	if arr[1] == 3 {
		return 4
	}
	return 5
}

func (a *Test) ParamArray(arr []int) int {
	if arr[1] == 3 {
		return 4
	}
	return 5
}

func (a *Test) AllocArrayOfArrays() int {
	arr1 := []int{1, 2, 3}
	arr2 := []int{4, 5, 6}
	arr3 := []int{7, 8, 9}
	arr := [][]int{arr1, arr2, arr3}
	if arr[0][1] == 2 {
		return 3
	}

	return 5
}

func (a *Test) ParamArrayOfArrays(arr [][]int) int {
	if arr[1][2] == 3 {
		return 4
	}
	return 5
}

type Foo struct {
	a int
	b bool
}

func (a *Test) SimpleStructures1() int {
	foo := Foo{a: 3, b: true}
	if foo.a == 3 {
		return 4
	}

	return 5
}

func (a *Test) SimpleStructures2() int {
	foo := Foo{a: 3, b: true}
	if foo.a == 4 {
		return 4
	}

	return 5
}

func (a *Test) ParamStructure(foo Foo) int {
	if foo.a == 3 {
		return 4
	}

	return 5
}

type Bar struct {
	fooRef *Foo
	foo    Foo
}

func (a *Test) StructOfStruct1() int {
	foo := Foo{a: 3, b: true}
	bar := Bar{fooRef: &foo, foo: foo}
	if bar.foo.a == 3 {
		return 4
	}

	return 5
}

func (a *Test) StructOfStruct2() int {
	foo := Foo{a: 3, b: true}
	bar := Bar{fooRef: &foo, foo: foo}
	if bar.foo.a == 3 {
		return 4
	}

	return 5
}

func (a *Test) ParamStructOfStruct(bar Bar) int {
	if bar.fooRef.a == 3 {
		return 4
	}

	return 5
}

type Baz struct {
	array []int
}

func (a *Test) StructOfArray() int {
	arr := []int{1, 2, 3}
	baz := Baz{array: arr}
	if baz.array[2] == 3 {
		return 4
	}

	return 5
}

func (a *Test) ParamStructOfArray(baz Baz) int {
	if baz.array[2] == 3 {
		return 4
	}

	return 5
}

func (a *Test) ArrayOfStructs() int {
	foo1 := Foo{a: 3, b: false}
	foo2 := Foo{a: 4, b: true}
	foo3 := Foo{a: 5, b: false}
	arr := []Foo{foo1, foo2, foo3}

	if arr[1].a == 4 {
		return 4
	}

	return 5
}

func (a *Test) ParamArrayOfStructs(arr []Foo) int {
	if arr[1].a == 4 {
		return 4
	}

	return 5
}

func (a *Test) ArrayAssign() int {
	arr := []int{1, 2, 3}
	arr[1] = 3
	if arr[1] == 3 {
		return 4
	}
	return 5
}

func (a *Test) ParamArrayAssign(arr []int) int {
	arr[1] = 3
	if arr[1] == 3 {
		return 4
	}
	return 5
}

func (a *Test) Aliasing(foo1 *Foo, foo2 *Foo) int {
	foo2.a = 5
	foo1.a = 2
	if foo2.a == 2 {
		return 4
	}
	return 5
}
