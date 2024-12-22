package se

import "strings"

type SymbolicMemory struct {
	heap       map[string]Array
	structs    map[string]map[int]string
	allocated  map[string]int
	SmtBuilder *SmtBuilder
}

func (memory *SymbolicMemory) Allocate(typ string) *Ref {
	id, ok := memory.allocated[typ]
	if !ok {
		memory.allocated[typ] = 1
		id = 1
	} else {
		id++
		memory.allocated[typ] = id
	}

	ref := &Ref{&Literal[int]{id}, typ}

	return ref
}

func (memory *SymbolicMemory) Deref(ref SymbolicExpression) SymbolicExpression {
	typ := GetType(ref)
	region, ok := memory.heap[typ]
	if !ok {
		region = memory.allocateRegion(typ)
	}

	value, ok := region.KnownValues[ref]
	if ok {
		return value
	}

	return &ArrayAccess{&region, ref}
}

func (memory *SymbolicMemory) AllocateStruct(typ string, fields map[int]string) *Ref {
	_, ok := memory.structs[typ]
	if ok {
		return memory.Allocate(typ)
	}

	ref := memory.Allocate(typ)
	for _, fieldType := range fields {
		memory.allocateRegion(fieldType)
		if strings.HasPrefix(fieldType, "[]") {
			memory.heap[fieldType].KnownValues[ref] = memory.AllocateArray(fieldType[2:])
		}
	}

	memory.structs[typ] = fields

	return ref
}

func (memory *SymbolicMemory) AssignField(fieldRef SymbolicExpression, field int, value SymbolicExpression) {
	structure := memory.structs[GetType(fieldRef)]
	fieldValue, ok := structure[field]
	if !ok {
		structure[field] = GetType(value)
	}
	region := memory.heap[fieldValue]
	region.KnownValues[fieldRef] = value
}

func (memory *SymbolicMemory) GetField(from SymbolicExpression, field int) SymbolicExpression {
	tpe := GetType(from)
	structure := memory.structs[tpe]
	fieldValue := structure[field]
	region := memory.heap[fieldValue]
	_, isRef := from.(*Ref)
	value, ok := region.KnownValues[from]
	if isRef && ok {
		return value
	}

	return &ArrayAccess{&region, from}
}

func (memory *SymbolicMemory) AllocateArray(elementType string) *Ref {
	arrayFields := make(map[int]string)
	return memory.AllocateStruct("[]"+elementType, arrayFields)
}

func (memory *SymbolicMemory) AssignToArray(array SymbolicExpression, index SymbolicExpression, value SymbolicExpression) {
	tpe := GetType(array)
	arrayRegion, ok := memory.heap[tpe]
	if !ok {
		arrayRegion = memory.allocateRegion(tpe)
	}
	arrayValue, ok := arrayRegion.KnownValues[array]
	if ok {
		arrayValue.(*Array).KnownValues[index] = value
	} else {
		arrayRegion.KnownValues[array] = &Array{KnownValues: map[SymbolicExpression]SymbolicExpression{index: value}}
	}
}

func (memory *SymbolicMemory) GetFromArray(array SymbolicExpression, index SymbolicExpression) SymbolicExpression {
	tpe := GetType(array)
	arrayRegion, ok := memory.heap[tpe]
	if !ok {
		arrayRegion = memory.allocateRegion(tpe)
	}
	arrayValue, ok := arrayRegion.KnownValues[array]
	if ok {
		return arrayValue.(*Array).KnownValues[index]
	} else {
		return &ArrayAccess{&ArrayAccess{&arrayRegion, array}, index}
	}
}

func (memory *SymbolicMemory) Copy() *SymbolicMemory {
	newHeap := make(map[string]Array)
	for k, v := range memory.heap {
		newHeap[k] = v
	}

	newStructs := make(map[string]map[int]string)
	for k, v := range memory.structs {
		newStructs[k] = copyMap(v)
	}

	newAllocated := make(map[string]int)
	for k, v := range memory.allocated {
		newAllocated[k] = v
	}

	return &SymbolicMemory{
		newHeap, newStructs, newAllocated, memory.SmtBuilder,
	}
}

func (memory *SymbolicMemory) allocateRegion(regionType string) Array {
	region, ok := memory.heap[regionType]
	if ok {
		return region
	}

	region = Array{regionType, make(map[SymbolicExpression]SymbolicExpression)}
	memory.heap[regionType] = region
	return region
}
