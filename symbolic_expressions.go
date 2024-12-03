package main

import (
	"fmt"
	"strings"
)

type SymbolicExpression interface {
	String() string
}

type BinaryOperationType = int

const (
	Add BinaryOperationType = iota
	Sub
	Mul
	Div
	Mod
	And
	Or
	Xor
	AndNot
	LeftShift
	RightShift
)

type BinaryOperation struct {
	Left, Right SymbolicExpression
	Type        BinaryOperationType
}

func CreateAnd(left SymbolicExpression, right SymbolicExpression) SymbolicExpression {
	l, ok := left.(*Literal[bool])
	if ok && l.Value {
		return right
	}
	if ok && !l.Value {
		return l
	}

	r, ok := right.(*Literal[bool])
	if ok && r.Value {
		return left
	}
	if ok && !r.Value {
		return r
	}

	return &BinaryOperation{left, right, And}
}

func CreateOr(left SymbolicExpression, right SymbolicExpression) SymbolicExpression {
	l, ok := left.(*Literal[bool])
	if ok && l.Value {
		return l
	}
	if ok && !l.Value {
		return right
	}

	r, ok := right.(*Literal[bool])
	if ok && r.Value {
		return r
	}
	if ok && !r.Value {
		return left
	}

	return &BinaryOperation{left, right, Or}
}

func CreateAdd(left SymbolicExpression, right SymbolicExpression) SymbolicExpression {
	l, lok := left.(*Literal[int])
	r, rok := right.(*Literal[int])

	if lok && rok {
		return &Literal[int]{l.Value + r.Value}
	}

	return &BinaryOperation{left, right, Add}
}

func (binOp *BinaryOperation) String() string {
	var operator string
	switch binOp.Type {
	case Add:
		operator = " + "
	case Sub:
		operator = " - "
	case Mul:
		operator = " * "
	case Div:
		operator = " / "
	case Mod:
		operator = " % "
	case And:
		operator = " & "
	case Or:
		operator = " | "
	case Xor:
		operator = " ^ "
	case AndNot:
		operator = " &^ "
	case LeftShift:
		operator = " << "
	case RightShift:
		operator = " >> "
	}
	return fmt.Sprintf("%v%v%v", binOp.Left.String(), operator, binOp.Right.String())

}

type Equals struct {
	Left, Right SymbolicExpression
}

func (equals *Equals) String() string {
	return fmt.Sprintf("%v == %v", equals.Left.String(), equals.Right.String())
}

type Cast struct {
	Value SymbolicExpression
	To    string
}

func (cast *Cast) String() string {
	return fmt.Sprintf("(%v)%v", cast.To, cast.Value.String())
}

type Not struct {
	Operand SymbolicExpression
}

func (not *Not) String() string {
	return fmt.Sprintf("!(%v)", not.Operand.String())
}

type LT struct {
	Left, Right SymbolicExpression
}

func (lt *LT) String() string {
	return fmt.Sprintf("%v < %v", lt.Left.String(), lt.Right.String())
}

type GT struct {
	Left, Right SymbolicExpression
}

func (gt *GT) String() string {
	return fmt.Sprintf("%v > %v", gt.Left.String(), gt.Right.String())
}

type LiteralValue interface {
	bool |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		int8 |
		int16 |
		int32 |
		int64 |
		float32 |
		float64 |
		int |
		uint |
		uintptr
}

type Literal[T LiteralValue] struct {
	Value T
}

func (literal *Literal[any]) String() string {
	return fmt.Sprintf("%v", literal.Value)
}

type InputValue struct {
	Name string
	Type string
}

func (inputValue *InputValue) String() string {
	return inputValue.Name
}

type ComplexLiteral[T float32 | float64] struct {
	Real      T
	Imaginary T
}

func (complexLiteral *ComplexLiteral[any]) String() string {
	return fmt.Sprintf("%v + %vi", complexLiteral.Real, complexLiteral.Imaginary)
}

type FunctionCall struct {
	Signature string
	Arguments []SymbolicExpression
}

func (functionCall *FunctionCall) String() string {
	args := make([]string, len(functionCall.Arguments))
	for i := range functionCall.Arguments {
		args[i] = functionCall.Arguments[i].String()
	}
	return fmt.Sprintf("%v(%v)", functionCall.Signature, strings.Join(args, ", "))
}

type Array struct {
	ComponentType string
	Size          SymbolicExpression
	KnownValues   map[SymbolicExpression]SymbolicExpression
}

func (array *Array) String() string {
	return fmt.Sprintf("%v[%v]", array.ComponentType, array.Size)
}

type ArrayAccess struct {
	Array SymbolicExpression
	Index SymbolicExpression
}

func (arrayAccess *ArrayAccess) String() string {
	return fmt.Sprintf("%v[%v]", arrayAccess.Array, arrayAccess.Index)
}

type Conditional struct {
	Options map[SymbolicExpression]SymbolicExpression
}

func (conditional *Conditional) String() string {
	var options []string
	for key, value := range conditional.Options {
		options = append(options, fmt.Sprintf("%v => %v", key.String(), value.String()))
	}
	return strings.Join(options, "\n")
}
