package se

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
	"symbolic-execution-2024/z3"
)

type SmtBuilder struct {
	Context    *z3.Context
	callsCount map[string]int
}

func (sb *SmtBuilder) BuildSmt(expression SymbolicExpression) []z3.Value {
	switch expression.(type) {
	case *BinaryOperation:
		left := sb.BuildSmt(expression.(*BinaryOperation).Left)
		right := sb.BuildSmt(expression.(*BinaryOperation).Right)
		switch expression.(*BinaryOperation).Type {
		case Add:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.Add(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					return left.Add(right)
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		case Sub:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.Sub(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					return left.Sub(right)
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		case Mul:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.Mul(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					return left.Mul(right)
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		case Div:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.SDiv(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					return left.Div(right)
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		case Mod:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.SRem(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					return left.Rem(right)
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		case And:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.And(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					panic("unexpected sort")
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					return left.And(right)
				})
		case Or:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.Or(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					panic("unexpected sort")
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					return left.Or(right)
				})
		case Xor:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.Xor(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					panic("unexpected sort")
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					return left.Xor(right)
				})
		case AndNot:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.And(right.Not())
				}, func(left z3.Float, right z3.Float) z3.Value {
					panic("unexpected sort")
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		case LeftShift:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.Lsh(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					panic("unexpected sort")
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		case RightShift:
			return processExpressions(left, right,
				func(left z3.BV, right z3.BV) z3.Value {
					return left.URsh(right)
				}, func(left z3.Float, right z3.Float) z3.Value {
					panic("unexpected sort")
				}, func(left z3.Bool, right z3.Bool) z3.Value {
					panic("unexpected sort")
				})
		}

	case *Equals:
		left := sb.BuildSmt(expression.(*Equals).Left)
		right := sb.BuildSmt(expression.(*Equals).Right)

		return processExpressions(left, right,
			func(left z3.BV, right z3.BV) z3.Value {
				return left.Eq(right)
			}, func(left z3.Float, right z3.Float) z3.Value {
				return left.Eq(right)
			}, func(left z3.Bool, right z3.Bool) z3.Value {
				return left.Eq(right)
			})

	case *Cast:
		value := sb.BuildSmt(expression.(*Cast).Value)[0]
		if value.Sort().Kind() == z3.KindBV && expression.(*Cast).To == "float64" {
			return []z3.Value{value.(z3.BV).SToFloat(sb.Context.FloatSort(11, 53))}
		} else if value.Sort().Kind() == z3.KindBV {
			return []z3.Value{value.(z3.BV)}
		} else {
			panic("unsupported cast")
		}

	case *Not:
		operand := sb.BuildSmt(expression.(*Not).Operand)[0]

		if operand.Sort().Kind() == z3.KindBV {
			return []z3.Value{operand.(z3.BV).Not()}
		} else if operand.Sort().Kind() == z3.KindBool {
			return []z3.Value{operand.(z3.Bool).Not()}
		} else {
			panic("unexpected sort")
		}

	case *LT:
		left := sb.BuildSmt(expression.(*LT).Left)
		right := sb.BuildSmt(expression.(*LT).Right)

		return processExpressions(left, right,
			func(left z3.BV, right z3.BV) z3.Value {
				return left.SLT(right)
			}, func(left z3.Float, right z3.Float) z3.Value {
				return left.LT(right)
			}, func(left z3.Bool, right z3.Bool) z3.Value {
				panic("unexpected sort")
			})

	case *GT:
		left := sb.BuildSmt(expression.(*GT).Left)
		right := sb.BuildSmt(expression.(*GT).Right)

		return processExpressions(left, right,
			func(left z3.BV, right z3.BV) z3.Value {
				return left.SGT(right)
			}, func(left z3.Float, right z3.Float) z3.Value {
				return left.GT(right)
			}, func(left z3.Bool, right z3.Bool) z3.Value {
				panic("unexpected sort")
			})

	case *Literal[bool]:
		return []z3.Value{sb.Context.FromBool(expression.(*Literal[bool]).Value)}
	case *Literal[uint8]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[uint8]).Value), sb.Context.BVSort(8))}
	case *Literal[uint16]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[uint16]).Value), sb.Context.BVSort(16))}
	case *Literal[uint32]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[uint32]).Value), sb.Context.BVSort(32))}
	case *Literal[uint64]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[uint64]).Value), sb.Context.BVSort(64))}
	case *Literal[int8]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[int8]).Value), sb.Context.BVSort(8))}
	case *Literal[int16]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[int16]).Value), sb.Context.BVSort(16))}
	case *Literal[int32]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[int32]).Value), sb.Context.BVSort(32))}
	case *Literal[int64]:
		return []z3.Value{sb.Context.FromInt(expression.(*Literal[int64]).Value, sb.Context.BVSort(64))}
	case *Literal[float32]:
		return []z3.Value{sb.Context.FromFloat32(expression.(*Literal[float32]).Value, sb.Context.FloatSort(8, 24))}
	case *Literal[float64]:
		return []z3.Value{sb.Context.FromFloat64(expression.(*Literal[float64]).Value, sb.Context.FloatSort(11, 53))}
	case *Literal[int]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[int]).Value), sb.Context.BVSort(bits.UintSize))}
	case *Literal[uint]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[uint]).Value), sb.Context.BVSort(bits.UintSize))}
	case *Literal[uintptr]:
		return []z3.Value{sb.Context.FromInt(int64(expression.(*Literal[uintptr]).Value), sb.Context.BVSort(bits.UintSize))}

	case *InputValue:
		expressionName := expression.(*InputValue).Name
		expressionType := expression.(*InputValue).Type
		if expressionType == "complex64" {
			return []z3.Value{
				sb.Context.Const("$R_"+expressionName, sb.Context.FloatSort(8, 24)),
				sb.Context.Const("$I_"+expressionName, sb.Context.FloatSort(8, 24)),
			}
		}
		if expressionType == "complex128" {
			return []z3.Value{
				sb.Context.Const("$R_"+expressionName, sb.Context.FloatSort(11, 53)),
				sb.Context.Const("$I_"+expressionName, sb.Context.FloatSort(11, 53)),
			}
		}
		if strings.HasPrefix(expressionType, "[]") {
			return []z3.Value{
				sb.Context.Const(expressionName, sb.typeSignatureToSort(expressionType)),
			}
		}
		return []z3.Value{sb.Context.Const(expressionName, sb.typeSignatureToSort(expressionType))}

	case *ComplexLiteral[float32]:
		return []z3.Value{
			sb.Context.FromFloat32(expression.(*ComplexLiteral[float32]).Real, sb.Context.FloatSort(8, 24)),
			sb.Context.FromFloat32(expression.(*ComplexLiteral[float32]).Imaginary, sb.Context.FloatSort(8, 24)),
		}
	case *ComplexLiteral[float64]:
		return []z3.Value{
			sb.Context.FromFloat64(expression.(*ComplexLiteral[float64]).Real, sb.Context.FloatSort(11, 53)),
			sb.Context.FromFloat64(expression.(*ComplexLiteral[float64]).Imaginary, sb.Context.FloatSort(11, 53)),
		}

	case *FunctionCall:
		args := expression.(*FunctionCall).Arguments
		signature := expression.(*FunctionCall).Signature
		if strings.HasPrefix(signature, "{") {
			receiverType := strings.Split(signature, "_")[0]
			argTypes := strings.Split(receiverType[1:len(receiverType)-1], ",")
			index, _ := strconv.Atoi(strings.SplitAfter(signature, "_")[1])

			tpe := GetType(args[0])
			funcDecl := sb.uninterpretedFunction(strings.Split(strings.SplitAfter(signature, "_")[1], "(")[0], []string{tpe}, argTypes[index])
			return []z3.Value{funcDecl.Apply(sb.BuildSmt(args[0])[0])}
		}
		switch signature {
		case "builtin_real(ComplexType)":
			return []z3.Value{sb.BuildSmt(args[0])[0]}
		case "builtin_imag(ComplexType)":
			return []z3.Value{sb.BuildSmt(args[0])[1]}
		case "builtin_len(Type)":
			var tpe string
			_, isRef := args[0].(*Ref)
			if isRef {
				tpe = "int"
			} else {
				tpe = GetType(args[0])
			}
			funcDecl := sb.uninterpretedFunction(strings.Split(strings.SplitAfter(expression.(*FunctionCall).Signature, "_")[1], "(")[0], []string{tpe}, "int")
			return []z3.Value{funcDecl.Apply(sb.BuildSmt(args[0])[0])}
		case "math.IsNaN(float64)":
			return []z3.Value{sb.BuildSmt(&Equals{args[0], &Literal[float64]{math.NaN()}})[0]}
		case "math/rand.Int()":
			num, ok := sb.callsCount[signature]
			if !ok {
				num = 1
			}
			sb.callsCount[signature] = num + 1
			return []z3.Value{sb.uninterpretedFunction(signature+"_"+strconv.Itoa(num), []string{}, "int").Apply()}
		}

	case *Array:
		valueSort := sb.typeSignatureToSort(expression.(*Array).ComponentType)
		array := sb.Context.FreshConst("arr_", sb.Context.ArraySort(sb.typeSignatureToSort("int"), valueSort))

		for key, value := range expression.(*Array).KnownValues {
			keySmt := sb.BuildSmt(key)[0]
			valueSmt := sb.BuildSmt(value)[0]
			if keySmt.Sort().String() == sb.typeSignatureToSort("int").String() && valueSmt.Sort().String() == valueSort.String() {
				array = array.(z3.Array).Store(keySmt, valueSmt)
			}
		}
		return []z3.Value{array}

	case *ArrayAccess:
		array := sb.BuildSmt(expression.(*ArrayAccess).Array)[0]
		index := sb.BuildSmt(expression.(*ArrayAccess).Index)[0]
		return []z3.Value{array.(z3.Array).Select(index)}

	case *Conditional:
		var result z3.Value
		for key, value := range expression.(*Conditional).Options {
			cond := sb.BuildSmt(key)[0].(z3.Bool)
			expr := sb.BuildSmt(value)[0]
			if result == nil {
				result = cond.IfThenElse(expr, z3.Uninterpreted{})
			} else {
				result = cond.IfThenElse(expr, result)
			}
		}
		return []z3.Value{result}
	case *Ref:
		return sb.BuildSmt(expression.(*Ref).Representation)
	}
	panic("unexpected expression")
}

func (sb *SmtBuilder) typeSignatureToSort(signature string) z3.Sort {
	if strings.HasPrefix(signature, "[]") {
		valueSort := sb.typeSignatureToSort("int")
		return sb.Context.ArraySort(sb.typeSignatureToSort("int"), valueSort)
	}
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
	return sb.Context.BVSort(bits.UintSize) // ref
}

func GetType(se SymbolicExpression) string {
	switch se.(type) {
	case *InputValue:
		return se.(*InputValue).Type
	case *ArrayAccess:
		return strings.SplitN(GetType(se.(*ArrayAccess).Array), "]", 2)[1]
	case *Literal[int]:
		return "int"
	case *Literal[bool]:
		return "bool"
	case *Literal[uint8]:
		return "uint8"
	case *Literal[uint16]:
		return "uint16"
	case *Literal[uint32]:
		return "uint32"
	case *Literal[uint64]:
		return "uint64"
	case *Literal[int8]:
		return "int8"
	case *Literal[int16]:
		return "int16"
	case *Literal[int32]:
		return "int32"
	case *Literal[int64]:
		return "int64"
	case *Literal[float32]:
		return "float32"
	case *Literal[float64]:
		return "float64"
	case *Literal[uint]:
		return "uint"
	case *Literal[uintptr]:
		return "uintptr"
	case *Array:
		return se.(*Array).String()
	case *Ref:
		return se.(*Ref).Type
	}
	panic("unexpected expression")
}

func processExpressions(leftArgs []z3.Value, rightArgs []z3.Value,
	bvOperation func(left z3.BV, right z3.BV) z3.Value,
	fpOperation func(left z3.Float, right z3.Float) z3.Value,
	boolOperation func(left z3.Bool, right z3.Bool) z3.Value) []z3.Value {

	if len(leftArgs) == 1 && len(rightArgs) == 1 {
		left := leftArgs[0]
		right := rightArgs[0]

		if left.Sort().Kind() == z3.KindBV && right.Sort().Kind() == z3.KindBV {
			return []z3.Value{
				bvOperation(left.(z3.BV), right.(z3.BV)),
			}
		} else if left.Sort().Kind() == z3.KindFloatingPoint && right.Sort().Kind() == z3.KindFloatingPoint {
			return []z3.Value{
				fpOperation(left.(z3.Float), right.(z3.Float)),
			}
		} else if left.Sort().Kind() == z3.KindBool && right.Sort().Kind() == z3.KindBool {
			return []z3.Value{
				boolOperation(left.(z3.Bool), right.(z3.Bool)),
			}
		} else {
			panic("unexpected sort")
		}
	} else if len(leftArgs) == 2 && len(rightArgs) == 2 {
		return []z3.Value{
			processExpressions([]z3.Value{leftArgs[0]}, []z3.Value{rightArgs[0]}, bvOperation, fpOperation, boolOperation)[0],
			processExpressions([]z3.Value{leftArgs[1]}, []z3.Value{rightArgs[1]}, bvOperation, fpOperation, boolOperation)[0],
		}
	}

	panic("unexpected state")
}

func (sb *SmtBuilder) uninterpretedFunction(signature string, args []string, returnType string) z3.FuncDecl {
	var argSorts []z3.Sort
	for _, arg := range args {
		argSorts = append(argSorts, sb.typeSignatureToSort(arg))
	}
	return sb.Context.FuncDecl(
		signature,
		argSorts,
		sb.typeSignatureToSort(returnType),
	)
}
