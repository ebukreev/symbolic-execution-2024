package main

import (
	"github.com/aclements/go-z3/z3"
	"math/bits"
)

type SmtBuilder struct {
	Context *z3.Context
}

func (sb *SmtBuilder) BuildSmt(expression SymbolicExpression) z3.Value {
	switch expression.(type) {
	case *BinaryOperation:
		left := sb.BuildSmt(expression.(*BinaryOperation).Left)
		right := sb.BuildSmt(expression.(*BinaryOperation).Right)
		switch expression.(*BinaryOperation).Type {
		case Add:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).Add(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
				return left.(z3.Float).Add(right.(z3.Float))
			} else {
				panic("unexpected sort")
			}
		case Sub:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).Sub(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
				return left.(z3.Float).Sub(right.(z3.Float))
			} else {
				panic("unexpected sort")
			}
		case Mul:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).Mul(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
				return left.(z3.Float).Mul(right.(z3.Float))
			} else {
				panic("unexpected sort")
			}
		case Div:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).SDiv(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
				return left.(z3.Float).Div(right.(z3.Float))
			} else {
				panic("unexpected sort")
			}
		case Mod:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).SRem(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
				return left.(z3.Float).Rem(right.(z3.Float))
			} else {
				panic("unexpected sort")
			}
		case And:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).And(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindBool && right.Sort().Kind() == z3.KindBool {
				return left.(z3.Bool).And(right.(z3.Bool))
			} else {
				panic("unexpected sort")
			}
		case Or:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).Or(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindBool && right.Sort().Kind() == z3.KindBool {
				return left.(z3.Bool).Or(right.(z3.Bool))
			} else {
				panic("unexpected sort")
			}
		case Xor:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).Xor(right.(z3.BV))
			} else if left.Sort().Kind() == z3.KindBool && right.Sort().Kind() == z3.KindBool {
				return left.(z3.Bool).Xor(right.(z3.Bool))
			} else {
				panic("unexpected sort")
			}
		case AndNot:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).And(right.(z3.BV).Not())
			} else {
				panic("unexpected sort")
			}
		case LeftShift:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).Lsh(right.(z3.BV))
			} else {
				panic("unexpected sort")
			}
		case RightShift:
			if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
				return left.(z3.BV).URsh(right.(z3.BV))
			} else {
				panic("unexpected sort")
			}
		}

	case *Equals:
		left := sb.BuildSmt(expression.(*Equals).Left)
		right := sb.BuildSmt(expression.(*Equals).Right)

		if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
			return left.(z3.BV).Eq(right.(z3.BV))
		} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
			return left.(z3.Float).Eq(right.(z3.Float))
		} else if left.Sort().Kind() == z3.KindBool && right.Sort().Kind() == z3.KindBool {
			return left.(z3.Bool).Eq(right.(z3.Bool))
		} else {
			panic("unexpected sort")
		}

	case *Cast:
		value := sb.BuildSmt(expression.(*Cast).Value)
		if value.Sort().Kind() == z3.KindBV && expression.(*Cast).To == "float64" {
			return value.(z3.BV).IEEEToFloat(sb.Context.FloatSort(11, 53))
		} else {
			panic("unsupported cast")
		}

	case *Not:
		operand := sb.BuildSmt(expression.(*Not).Operand)

		if operand.Sort().Kind() == z3.KindBV {
			return operand.(z3.BV).Not()
		} else if operand.Sort().Kind() == z3.KindBool {
			return operand.(z3.Bool).Not()
		} else {
			panic("unexpected sort")
		}

	case *LT:
		left := sb.BuildSmt(expression.(*LT).Left)
		right := sb.BuildSmt(expression.(*LT).Right)

		if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
			return left.(z3.BV).SLT(right.(z3.BV))
		} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
			return left.(z3.Float).LT(right.(z3.Float))
		} else {
			panic("unexpected sort")
		}

	case *GT:
		left := sb.BuildSmt(expression.(*GT).Left)
		right := sb.BuildSmt(expression.(*GT).Right)

		if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
			return left.(z3.BV).SGT(right.(z3.BV))
		} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
			return left.(z3.Float).GT(right.(z3.Float))
		} else {
			panic("unexpected sort")
		}

	case *Literal[bool]:
		return sb.Context.FromBool(expression.(*Literal[bool]).Value)
	case *Literal[uint8]:
		return sb.Context.FromInt(int64(expression.(*Literal[uint8]).Value), sb.Context.BVSort(8))
	case *Literal[uint16]:
		return sb.Context.FromInt(int64(expression.(*Literal[uint16]).Value), sb.Context.BVSort(16))
	case *Literal[uint32]:
		return sb.Context.FromInt(int64(expression.(*Literal[uint32]).Value), sb.Context.BVSort(32))
	case *Literal[uint64]:
		return sb.Context.FromInt(int64(expression.(*Literal[uint64]).Value), sb.Context.BVSort(64))
	case *Literal[int8]:
		return sb.Context.FromInt(int64(expression.(*Literal[int8]).Value), sb.Context.BVSort(8))
	case *Literal[int16]:
		return sb.Context.FromInt(int64(expression.(*Literal[int16]).Value), sb.Context.BVSort(16))
	case *Literal[int32]:
		return sb.Context.FromInt(int64(expression.(*Literal[int32]).Value), sb.Context.BVSort(32))
	case *Literal[int64]:
		return sb.Context.FromInt(expression.(*Literal[int64]).Value, sb.Context.BVSort(64))
	case *Literal[float32]:
		return sb.Context.FromFloat32(expression.(*Literal[float32]).Value, sb.Context.FloatSort(8, 24))
	case *Literal[float64]:
		return sb.Context.FromFloat64(expression.(*Literal[float64]).Value, sb.Context.FloatSort(11, 53))
	case *Literal[int]:
		return sb.Context.FromInt(int64(expression.(*Literal[int]).Value), sb.Context.BVSort(bits.UintSize))
	case *Literal[uint]:
		return sb.Context.FromInt(int64(expression.(*Literal[uint]).Value), sb.Context.BVSort(bits.UintSize))
	case *Literal[uintptr]:
		return sb.Context.FromInt(int64(expression.(*Literal[uintptr]).Value), sb.Context.BVSort(bits.UintSize))

	case *InputValue:
		return sb.Context.Const(expression.(*InputValue).Name, sb.typeSignatureToSort(expression.(*InputValue).Type))
	}
	panic("unexpected expression")
}

func (sb *SmtBuilder) typeSignatureToSort(signature string) z3.Sort {
	switch signature {
	case "bool":
		return sb.Context.BoolSort()
	case "uint8", "int8":
		return sb.Context.BVSort(8)
	case "uint16", "int16":
		return sb.Context.BVSort(16)
	case "uint32", "int32":
		return sb.Context.BVSort(32)
	case "uint64", "int64":
		return sb.Context.BVSort(64)
	case "float32":
		return sb.Context.FloatSort(8, 24)
	case "float64":
		return sb.Context.FloatSort(11, 53)
	case "int", "uint", "uintptr":
		return sb.Context.BVSort(bits.UintSize)
	}
	panic("unexpected type signature")
}
