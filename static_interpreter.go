package main

import (
	"go/types"
	"golang.org/x/tools/go/ssa"
	"strconv"
	"strings"
)

type StaticInterpreter struct {
	PathCondition SymbolicExpression
	Memory        map[string]*Conditional
	ReturnValue   Conditional
	BlocksStack   []*ssa.BasicBlock
}

func InterpretStatically(function *ssa.Function) Conditional {
	interpreter := StaticInterpreter{PathCondition: &Literal[bool]{true},
		Memory:      map[string]*Conditional{},
		ReturnValue: Conditional{make(map[SymbolicExpression]SymbolicExpression)},
		BlocksStack: []*ssa.BasicBlock{},
	}

	interpreter.interpretStatically(function.Blocks[0])

	return interpreter.ReturnValue
}

func (interpreter *StaticInterpreter) interpretStatically(element interface{}) SymbolicExpression {
	switch element.(type) {
	case *ssa.Alloc:
		return interpreter.interpretAllocStatically(element.(*ssa.Alloc))
	case *ssa.BinOp:
		return interpreter.interpretBinOpStatically(element.(*ssa.BinOp))
	case *ssa.Builtin:
		return interpreter.interpretBuiltinStatically(element.(*ssa.Builtin))
	case *ssa.Call:
		return interpreter.interpretCallStatically(element.(*ssa.Call))
	case *ssa.ChangeInterface:
		return interpreter.interpretChangeInterfaceStatically(element.(*ssa.ChangeInterface))
	case *ssa.ChangeType:
		return interpreter.interpretChangeTypeStatically(element.(*ssa.ChangeType))
	case *ssa.Const:
		return interpreter.interpretConstStatically(element.(*ssa.Const))
	case *ssa.Convert:
		return interpreter.interpretConvertStatically(element.(*ssa.Convert))
	case *ssa.DebugRef:
		return interpreter.interpretDebugRefStatically(element.(*ssa.DebugRef))
	case *ssa.Defer:
		return interpreter.interpretDeferStatically(element.(*ssa.Defer))
	case *ssa.Extract:
		return interpreter.interpretExtractStatically(element.(*ssa.Extract))
	case *ssa.Field:
		return interpreter.interpretFieldStatically(element.(*ssa.Field))
	case *ssa.FieldAddr:
		return interpreter.interpretFieldAddrStatically(element.(*ssa.FieldAddr))
	case *ssa.FreeVar:
		return interpreter.interpretFreeVarStatically(element.(*ssa.FreeVar))
	case *ssa.Global:
		return interpreter.interpretGlobalStatically(element.(*ssa.Global))
	case *ssa.Go:
		return interpreter.interpretGoStatically(element.(*ssa.Go))
	case *ssa.If:
		return interpreter.interpretIfStatically(element.(*ssa.If))
	case *ssa.Index:
		return interpreter.interpretIndexStatically(element.(*ssa.Index))
	case *ssa.IndexAddr:
		return interpreter.interpretIndexAddrStatically(element.(*ssa.IndexAddr))
	case *ssa.Jump:
		return interpreter.interpretJumpStatically(element.(*ssa.Jump))
	case *ssa.Lookup:
		return interpreter.interpretLookupStatically(element.(*ssa.Lookup))
	case *ssa.MakeChan:
		return interpreter.interpretMakeChanStatically(element.(*ssa.MakeChan))
	case *ssa.MakeClosure:
		return interpreter.interpretMakeClosureStatically(element.(*ssa.MakeClosure))
	case *ssa.MakeInterface:
		return interpreter.interpretMakeInterfaceStatically(element.(*ssa.MakeInterface))
	case *ssa.MakeMap:
		return interpreter.interpretMakeMapStatically(element.(*ssa.MakeMap))
	case *ssa.MakeSlice:
		return interpreter.interpretMakeSliceStatically(element.(*ssa.MakeSlice))
	case *ssa.MapUpdate:
		return interpreter.interpretMapUpdateStatically(element.(*ssa.MapUpdate))
	case *ssa.MultiConvert:
		return interpreter.interpretMultiConvertStatically(element.(*ssa.MultiConvert))
	case *ssa.NamedConst:
		return interpreter.interpretNamedConstStatically(element.(*ssa.NamedConst))
	case *ssa.Next:
		return interpreter.interpretNextStatically(element.(*ssa.Next))
	case *ssa.Panic:
		return interpreter.interpretPanicStatically(element.(*ssa.Panic))
	case *ssa.Parameter:
		return interpreter.interpretParameterStatically(element.(*ssa.Parameter))
	case *ssa.Phi:
		return interpreter.interpretPhiStatically(element.(*ssa.Phi))
	case *ssa.Range:
		return interpreter.interpretRangeStatically(element.(*ssa.Range))
	case *ssa.Return:
		return interpreter.interpretReturnStatically(element.(*ssa.Return))
	case *ssa.RunDefers:
		return interpreter.interpretRunDefersStatically(element.(*ssa.RunDefers))
	case *ssa.Select:
		return interpreter.interpretSelectStatically(element.(*ssa.Select))
	case *ssa.Send:
		return interpreter.interpretSendStatically(element.(*ssa.Send))
	case *ssa.Slice:
		return interpreter.interpretSliceStatically(element.(*ssa.Slice))
	case *ssa.SliceToArrayPointer:
		return interpreter.interpretSliceToArrayPointerStatically(element.(*ssa.SliceToArrayPointer))
	case *ssa.Store:
		return interpreter.interpretStoreStatically(element.(*ssa.Store))
	case *ssa.Type:
		return interpreter.interpretTypeStatically(element.(*ssa.Type))
	case *ssa.TypeAssert:
		return interpreter.interpretTypeAssertStatically(element.(*ssa.TypeAssert))
	case *ssa.UnOp:
		return interpreter.interpretUnOpStatically(element.(*ssa.UnOp))
	case *ssa.BasicBlock:
		return interpreter.interpretBasicBlockStatically(element.(*ssa.BasicBlock))
	default:
		panic("Unexpected element")
	}
}

func (interpreter *StaticInterpreter) interpretAllocStatically(element *ssa.Alloc) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretBinOpStatically(element *ssa.BinOp) SymbolicExpression {
	left := interpreter.interpretStatically(element.X)
	right := interpreter.interpretStatically(element.Y)

	switch element.Op.String() {
	case "+":
		return &BinaryOperation{Left: left, Right: right, Type: Add}
	case "-":
		return &BinaryOperation{Left: left, Right: right, Type: Sub}
	case "*":
		return &BinaryOperation{Left: left, Right: right, Type: Mul}
	case "/":
		return &BinaryOperation{Left: left, Right: right, Type: Div}
	case "%":
		return &BinaryOperation{Left: left, Right: right, Type: Mod}
	case "&":
		return CreateAnd(left, right)
	case "|":
		return CreateOr(left, right)
	case "^":
		return &BinaryOperation{Left: left, Right: right, Type: Xor}
	case "<<":
		return &BinaryOperation{Left: left, Right: right, Type: LeftShift}
	case ">>":
		return &BinaryOperation{Left: left, Right: right, Type: RightShift}
	case "&^":
		return &BinaryOperation{Left: left, Right: right, Type: AndNot}
	case "==":
		return &Equals{Left: left, Right: right}
	case "!=":
		return &Not{&Equals{Left: left, Right: right}}
	case "<":
		return &LT{Left: left, Right: right}
	case "<=":
		return &BinaryOperation{Left: &LT{Left: left, Right: right}, Right: &Equals{Left: left, Right: right}, Type: Or}
	case ">":
		return &GT{Left: left, Right: right}
	case ">=":
		return &BinaryOperation{Left: &GT{Left: left, Right: right}, Right: &Equals{Left: left, Right: right}, Type: Or}
	default:
		panic("unexpected binOp")
	}
}

func (interpreter *StaticInterpreter) interpretBuiltinStatically(element *ssa.Builtin) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretCallStatically(element *ssa.Call) SymbolicExpression {
	signature := strings.ReplaceAll(element.Call.Value.String(), " ", "_") + "("
	var argValues []SymbolicExpression

	for _, arg := range element.Call.Args {
		argValues = append(argValues, interpreter.interpretStatically(arg))
		argType := arg.Type().String()
		if argType == "complex128" {
			signature += "ComplexType,"
		} else if strings.HasPrefix(signature, "builtin_len") {
			signature += "Type,"
		} else {
			signature += argType + ","
		}
	}
	if signature[len(signature)-1] == ',' {
		signature = signature[:len(signature)-1]
	}
	signature = signature + ")"

	return &FunctionCall{Signature: signature, Arguments: argValues}
}

func (interpreter *StaticInterpreter) interpretChangeInterfaceStatically(element *ssa.ChangeInterface) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretChangeTypeStatically(element *ssa.ChangeType) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretConstStatically(element *ssa.Const) SymbolicExpression {
	switch element.Type().String() {
	case "int":
		value, _ := strconv.Atoi(element.Value.ExactString())
		return &Literal[int]{value}
	case "uint":
		value, _ := strconv.ParseUint(element.Value.ExactString(), 10, 64)
		return &Literal[uint]{uint(value)}
	case "float64":
		value, _ := strconv.ParseFloat(element.Value.ExactString(), 64)
		return &Literal[float64]{value}
	}
	panic("unexpected const")
}

func (interpreter *StaticInterpreter) interpretConvertStatically(element *ssa.Convert) SymbolicExpression {
	return &Cast{interpreter.interpretStatically(element.X), element.Type().String()}
}

func (interpreter *StaticInterpreter) interpretDebugRefStatically(element *ssa.DebugRef) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretDeferStatically(element *ssa.Defer) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretExtractStatically(element *ssa.Extract) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretFieldStatically(element *ssa.Field) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretFieldAddrStatically(element *ssa.FieldAddr) SymbolicExpression {
	receiver := interpreter.interpretStatically(element.X)
	typeSignature := getTypeSignature(element.X.Type())
	return &FunctionCall{typeSignature + "_" + strconv.Itoa(element.Field), []SymbolicExpression{receiver}}
}

func (interpreter *StaticInterpreter) interpretFreeVarStatically(element *ssa.FreeVar) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretGlobalStatically(element *ssa.Global) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretGoStatically(element *ssa.Go) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretIfStatically(element *ssa.If) SymbolicExpression {
	cond := interpreter.interpretStatically(element.Cond)
	successors := element.Block().Succs
	pc := interpreter.PathCondition
	interpreter.enterBranch(CreateAnd(pc, cond), successors[0])
	if len(successors) == 1 {
		return nil
	}
	notCond := &Not{cond}
	interpreter.enterBranch(CreateAnd(pc, notCond), successors[1])
	return nil
}

func (interpreter *StaticInterpreter) interpretIndexStatically(element *ssa.Index) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretIndexAddrStatically(element *ssa.IndexAddr) SymbolicExpression {
	array := interpreter.interpretStatically(element.X)
	index := interpreter.interpretStatically(element.Index)
	return &ArrayAccess{array, index}
}

func (interpreter *StaticInterpreter) interpretJumpStatically(element *ssa.Jump) SymbolicExpression {
	return interpreter.interpretStatically(element.Block().Succs[0])
}

func (interpreter *StaticInterpreter) interpretLookupStatically(element *ssa.Lookup) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretMakeChanStatically(element *ssa.MakeChan) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretMakeClosureStatically(element *ssa.MakeClosure) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretMakeInterfaceStatically(element *ssa.MakeInterface) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretMakeMapStatically(element *ssa.MakeMap) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretMakeSliceStatically(element *ssa.MakeSlice) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretMapUpdateStatically(element *ssa.MapUpdate) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretMultiConvertStatically(element *ssa.MultiConvert) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretNamedConstStatically(element *ssa.NamedConst) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretNextStatically(element *ssa.Next) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretPanicStatically(element *ssa.Panic) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretParameterStatically(element *ssa.Parameter) SymbolicExpression {
	return &InputValue{Name: element.Name(), Type: element.Type().String()}
}

func (interpreter *StaticInterpreter) interpretPhiStatically(element *ssa.Phi) SymbolicExpression {
	for i, pred := range element.Block().Preds {
		for j := len(interpreter.BlocksStack) - 2; j >= 0; j-- {
			if pred == interpreter.BlocksStack[j] {
				elementValue, ok := interpreter.Memory[element.Comment]
				edgeValue := interpreter.interpretStatically(element.Edges[i])
				if ok {
					elementValue.Options[interpreter.PathCondition] = edgeValue
				} else {
					options := map[SymbolicExpression]SymbolicExpression{}
					options[interpreter.PathCondition] = edgeValue
					interpreter.Memory[element.Comment] = &Conditional{options}
				}
				return edgeValue
			}
		}
	}

	panic("unexpected state")
}

func (interpreter *StaticInterpreter) interpretRangeStatically(element *ssa.Range) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretReturnStatically(element *ssa.Return) SymbolicExpression {
	returnExpr := interpreter.interpretStatically(element.Results[0])
	interpreter.ReturnValue.Options[interpreter.PathCondition] = returnExpr
	return nil
}

func (interpreter *StaticInterpreter) interpretRunDefersStatically(element *ssa.RunDefers) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretSelectStatically(element *ssa.Select) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretSendStatically(element *ssa.Send) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretSliceStatically(element *ssa.Slice) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretSliceToArrayPointerStatically(element *ssa.SliceToArrayPointer) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretStoreStatically(element *ssa.Store) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretTypeStatically(element *ssa.Type) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretTypeAssertStatically(element *ssa.TypeAssert) SymbolicExpression {
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretUnOpStatically(element *ssa.UnOp) SymbolicExpression {
	operand := interpreter.interpretStatically(element.X)
	switch element.Op.String() {
	case "*":
		return operand
	}
	panic("TODO")
}

func (interpreter *StaticInterpreter) interpretBasicBlockStatically(element *ssa.BasicBlock) SymbolicExpression {
	interpreter.BlocksStack = append(interpreter.BlocksStack, element)
	for _, instr := range element.Instrs {
		interpreter.interpretStatically(instr)
	}
	interpreter.BlocksStack = interpreter.BlocksStack[:len(interpreter.BlocksStack)-1]
	return nil
}

func (interpreter *StaticInterpreter) enterBranch(pathCondition SymbolicExpression, body *ssa.BasicBlock) SymbolicExpression {
	interpreter.PathCondition = pathCondition
	interpreter.interpretStatically(body)
	return &interpreter.ReturnValue
}

func getTypeSignature(tpe types.Type) string {
	pointer, ok := tpe.(*types.Pointer)
	if ok {
		return getTypeSignature(pointer.Elem())
	}
	named, ok := tpe.(*types.Named)
	if ok {
		return getTypeSignature(named.Underlying())
	}
	structure, ok := tpe.(*types.Struct)
	if ok {
		signature := "{"
		for i := 0; i < structure.NumFields(); i++ {
			signature += structure.Field(i).Type().String() + ","
		}
		if signature[len(signature)-1] == ',' {
			return signature[:len(signature)-1] + "}"
		}
		return signature + "}"
	}
	panic("TODO")
}
