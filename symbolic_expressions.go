package main

type SymbolicExpression interface{}

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

type Equals struct {
	Left, Right SymbolicExpression
}

type Cast struct {
	Value SymbolicExpression
	To    string
}

type Not struct {
	Operand SymbolicExpression
}

type LT struct {
	Left, Right SymbolicExpression
}

type GT struct {
	Left, Right SymbolicExpression
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

type InputValue struct {
	Name string
	Type string
}

type ComplexLiteral[T float32 | float64] struct {
	Real      T
	Imaginary T
}

type FunctionCall struct {
	Signature string
	Arguments []SymbolicExpression
}
