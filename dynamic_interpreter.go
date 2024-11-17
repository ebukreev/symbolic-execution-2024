package main

import (
	"golang.org/x/tools/go/ssa"
	"strconv"
)

type DynamicInterpreter struct {
	Function        *ssa.Function
	PathCondition   SymbolicExpression
	Memory          map[string]SymbolicExpression
	ReturnValue     SymbolicExpression
	Instructions    []ssa.Instruction
	InstructionsPtr int
	BlocksStack     []*ssa.BasicBlock
}

func copyMap(m map[string]SymbolicExpression) map[string]SymbolicExpression {
	cp := make(map[string]SymbolicExpression)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func copySlice(slice []*ssa.BasicBlock) []*ssa.BasicBlock {
	tmp := make([]*ssa.BasicBlock, len(slice))
	copy(tmp, slice)
	return tmp
}

func (interpreter *DynamicInterpreter) copy() DynamicInterpreter {
	return DynamicInterpreter{
		interpreter.Function,
		interpreter.PathCondition,
		copyMap(interpreter.Memory),
		interpreter.ReturnValue,
		interpreter.Instructions,
		interpreter.InstructionsPtr,
		copySlice(interpreter.BlocksStack),
	}
}

func InterpretDynamically(interpreter DynamicInterpreter) []DynamicInterpreter {
	if interpreter.Instructions == nil {
		interpreter.BlocksStack = []*ssa.BasicBlock{interpreter.Function.Blocks[0]}
		interpreter.Instructions = interpreter.Function.Blocks[0].Instrs
		interpreter.InstructionsPtr = 0
	}

	front := interpreter.Instructions[interpreter.InstructionsPtr]
	interpreter.InstructionsPtr++
	return interpreter.interpretDynamically(front.(ssa.Instruction))
}

func (interpreter *DynamicInterpreter) interpretDynamically(element ssa.Instruction) []DynamicInterpreter {
	switch element.(type) {
	case *ssa.Alloc:
		return interpreter.interpretAllocDynamically(element.(*ssa.Alloc))
	case *ssa.BinOp:
		return interpreter.interpretBinOpDynamically(element.(*ssa.BinOp))
	case *ssa.Call:
		return interpreter.interpretCallDynamically(element.(*ssa.Call))
	case *ssa.ChangeInterface:
		return interpreter.interpretChangeInterfaceDynamically(element.(*ssa.ChangeInterface))
	case *ssa.ChangeType:
		return interpreter.interpretChangeTypeDynamically(element.(*ssa.ChangeType))
	case *ssa.Convert:
		return interpreter.interpretConvertDynamically(element.(*ssa.Convert))
	case *ssa.DebugRef:
		return interpreter.interpretDebugRefDynamically(element.(*ssa.DebugRef))
	case *ssa.Defer:
		return interpreter.interpretDeferDynamically(element.(*ssa.Defer))
	case *ssa.Extract:
		return interpreter.interpretExtractDynamically(element.(*ssa.Extract))
	case *ssa.Field:
		return interpreter.interpretFieldDynamically(element.(*ssa.Field))
	case *ssa.FieldAddr:
		return interpreter.interpretFieldAddrDynamically(element.(*ssa.FieldAddr))
	case *ssa.Go:
		return interpreter.interpretGoDynamically(element.(*ssa.Go))
	case *ssa.If:
		return interpreter.interpretIfDynamically(element.(*ssa.If))
	case *ssa.Index:
		return interpreter.interpretIndexDynamically(element.(*ssa.Index))
	case *ssa.IndexAddr:
		return interpreter.interpretIndexAddrDynamically(element.(*ssa.IndexAddr))
	case *ssa.Jump:
		return interpreter.interpretJumpDynamically(element.(*ssa.Jump))
	case *ssa.Lookup:
		return interpreter.interpretLookupDynamically(element.(*ssa.Lookup))
	case *ssa.MakeChan:
		return interpreter.interpretMakeChanDynamically(element.(*ssa.MakeChan))
	case *ssa.MakeClosure:
		return interpreter.interpretMakeClosureDynamically(element.(*ssa.MakeClosure))
	case *ssa.MakeInterface:
		return interpreter.interpretMakeInterfaceDynamically(element.(*ssa.MakeInterface))
	case *ssa.MakeMap:
		return interpreter.interpretMakeMapDynamically(element.(*ssa.MakeMap))
	case *ssa.MakeSlice:
		return interpreter.interpretMakeSliceDynamically(element.(*ssa.MakeSlice))
	case *ssa.MapUpdate:
		return interpreter.interpretMapUpdateDynamically(element.(*ssa.MapUpdate))
	case *ssa.MultiConvert:
		return interpreter.interpretMultiConvertDynamically(element.(*ssa.MultiConvert))
	case *ssa.Next:
		return interpreter.interpretNextDynamically(element.(*ssa.Next))
	case *ssa.Panic:
		return interpreter.interpretPanicDynamically(element.(*ssa.Panic))
	case *ssa.Phi:
		return interpreter.interpretPhiDynamically(element.(*ssa.Phi))
	case *ssa.Range:
		return interpreter.interpretRangeDynamically(element.(*ssa.Range))
	case *ssa.Return:
		return interpreter.interpretReturnDynamically(element.(*ssa.Return))
	case *ssa.RunDefers:
		return interpreter.interpretRunDefersDynamically(element.(*ssa.RunDefers))
	case *ssa.Select:
		return interpreter.interpretSelectDynamically(element.(*ssa.Select))
	case *ssa.Send:
		return interpreter.interpretSendDynamically(element.(*ssa.Send))
	case *ssa.Slice:
		return interpreter.interpretSliceDynamically(element.(*ssa.Slice))
	case *ssa.SliceToArrayPointer:
		return interpreter.interpretSliceToArrayPointerDynamically(element.(*ssa.SliceToArrayPointer))
	case *ssa.Store:
		return interpreter.interpretStoreDynamically(element.(*ssa.Store))
	case *ssa.TypeAssert:
		return interpreter.interpretTypeAssertDynamically(element.(*ssa.TypeAssert))
	case *ssa.UnOp:
		return interpreter.interpretUnOpDynamically(element.(*ssa.UnOp))
	default:
		panic("Unexpected element")
	}
}

func (interpreter *DynamicInterpreter) resolveExpression(value ssa.Value) SymbolicExpression {
	switch value.(type) {
	case *ssa.Const:
		return interpreter.resolveConst(value.(*ssa.Const))
	case *ssa.BinOp:
		return interpreter.resolveBinOp(value.(*ssa.BinOp))
	case *ssa.Parameter:
		return interpreter.resolveParameter(value.(*ssa.Parameter))
	case *ssa.Convert:
		return interpreter.resolveConvert(value.(*ssa.Convert))
	case *ssa.Phi:
		return interpreter.resolvePhi(value.(*ssa.Phi))
	default:
		panic("Unexpected element")
	}
}

func (interpreter *DynamicInterpreter) resolveConst(element *ssa.Const) SymbolicExpression {
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

func (interpreter *DynamicInterpreter) resolveBinOp(element *ssa.BinOp) SymbolicExpression {
	left := interpreter.resolveExpression(element.X)
	right := interpreter.resolveExpression(element.Y)

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

func (interpreter *DynamicInterpreter) resolveParameter(element *ssa.Parameter) SymbolicExpression {
	return &InputValue{Name: element.Name(), Type: element.Type().String()}
}

func (interpreter *DynamicInterpreter) resolveConvert(element *ssa.Convert) SymbolicExpression {
	return &Cast{interpreter.resolveExpression(element.X), element.Type().String()}
}

func (interpreter *DynamicInterpreter) resolvePhi(element *ssa.Phi) SymbolicExpression {
	for i, pred := range element.Block().Preds {
		for j := len(interpreter.BlocksStack) - 2; j >= 0; j-- {
			if pred == interpreter.BlocksStack[j] {
				edgeValue := interpreter.resolveExpression(element.Edges[i])
				interpreter.Memory[element.Comment] = edgeValue

				return edgeValue
			}
		}
	}

	panic("unexpected state")
}

func (interpreter *DynamicInterpreter) interpretAllocDynamically(element *ssa.Alloc) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretBinOpDynamically(element *ssa.BinOp) []DynamicInterpreter {
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretCallDynamically(element *ssa.Call) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretChangeInterfaceDynamically(element *ssa.ChangeInterface) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretChangeTypeDynamically(element *ssa.ChangeType) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretConvertDynamically(element *ssa.Convert) []DynamicInterpreter {
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretDebugRefDynamically(element *ssa.DebugRef) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretDeferDynamically(element *ssa.Defer) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretExtractDynamically(element *ssa.Extract) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretFieldDynamically(element *ssa.Field) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretFieldAddrDynamically(element *ssa.FieldAddr) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretGoDynamically(element *ssa.Go) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretIfDynamically(element *ssa.If) []DynamicInterpreter {
	cond := interpreter.resolveExpression(element.Cond)
	successors := element.Block().Succs

	thenBranch := interpreter.copy()
	elseBranch := interpreter.copy()

	thenBranch.PathCondition = CreateAnd(thenBranch.PathCondition, cond)
	thenBranch.BlocksStack = append(thenBranch.BlocksStack, successors[0])
	thenBranch.Instructions = successors[0].Instrs
	thenBranch.InstructionsPtr = 0

	elseBranch.PathCondition = CreateAnd(elseBranch.PathCondition, &Not{cond})
	elseBranch.BlocksStack = append(elseBranch.BlocksStack, successors[1])
	elseBranch.Instructions = successors[1].Instrs
	elseBranch.InstructionsPtr = 0

	return []DynamicInterpreter{
		thenBranch,
		elseBranch,
	}
}

func (interpreter *DynamicInterpreter) interpretIndexDynamically(element *ssa.Index) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretIndexAddrDynamically(element *ssa.IndexAddr) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretJumpDynamically(element *ssa.Jump) []DynamicInterpreter {
	interpreter.BlocksStack = append(interpreter.BlocksStack, element.Block().Succs[0])
	interpreter.Instructions = element.Block().Succs[0].Instrs
	interpreter.InstructionsPtr = 0
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretLookupDynamically(element *ssa.Lookup) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretMakeChanDynamically(element *ssa.MakeChan) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretMakeClosureDynamically(element *ssa.MakeClosure) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretMakeInterfaceDynamically(element *ssa.MakeInterface) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretMakeMapDynamically(element *ssa.MakeMap) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretMakeSliceDynamically(element *ssa.MakeSlice) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretMapUpdateDynamically(element *ssa.MapUpdate) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretMultiConvertDynamically(element *ssa.MultiConvert) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretNextDynamically(element *ssa.Next) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretPanicDynamically(element *ssa.Panic) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretPhiDynamically(element *ssa.Phi) []DynamicInterpreter {
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretRangeDynamically(element *ssa.Range) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretReturnDynamically(element *ssa.Return) []DynamicInterpreter {
	interpreter.ReturnValue = interpreter.resolveExpression(element.Results[0])
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretRunDefersDynamically(element *ssa.RunDefers) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretSelectDynamically(element *ssa.Select) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretSendDynamically(element *ssa.Send) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretSliceDynamically(element *ssa.Slice) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretSliceToArrayPointerDynamically(element *ssa.SliceToArrayPointer) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretStoreDynamically(element *ssa.Store) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretTypeAssertDynamically(element *ssa.TypeAssert) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretUnOpDynamically(element *ssa.UnOp) []DynamicInterpreter {
	panic("TODO")
}
