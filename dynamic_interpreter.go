package main

import "golang.org/x/tools/go/ssa"

type DynamicInterpreter struct {
	PathCondition   SymbolicExpression
	Memory          map[string]*SymbolicExpression
	ReturnValue     SymbolicExpression
	BlocksStack     []*ssa.BasicBlock
	NextInstruction interface{}
}

func InterpretDynamically(interpreter DynamicInterpreter) []DynamicInterpreter {
	return interpreter.interpretDynamically(interpreter.NextInstruction)
}

func (interpreter *DynamicInterpreter) interpretDynamically(element interface{}) []DynamicInterpreter {
	switch element.(type) {
	case *ssa.Alloc:
		return interpreter.interpretAllocDynamically(element.(*ssa.Alloc))
	case *ssa.BinOp:
		return interpreter.interpretBinOpDynamically(element.(*ssa.BinOp))
	case *ssa.Builtin:
		return interpreter.interpretBuiltinDynamically(element.(*ssa.Builtin))
	case *ssa.Call:
		return interpreter.interpretCallDynamically(element.(*ssa.Call))
	case *ssa.ChangeInterface:
		return interpreter.interpretChangeInterfaceDynamically(element.(*ssa.ChangeInterface))
	case *ssa.ChangeType:
		return interpreter.interpretChangeTypeDynamically(element.(*ssa.ChangeType))
	case *ssa.Const:
		return interpreter.interpretConstDynamically(element.(*ssa.Const))
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
	case *ssa.FreeVar:
		return interpreter.interpretFreeVarDynamically(element.(*ssa.FreeVar))
	case *ssa.Global:
		return interpreter.interpretGlobalDynamically(element.(*ssa.Global))
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
	case *ssa.NamedConst:
		return interpreter.interpretNamedConstDynamically(element.(*ssa.NamedConst))
	case *ssa.Next:
		return interpreter.interpretNextDynamically(element.(*ssa.Next))
	case *ssa.Panic:
		return interpreter.interpretPanicDynamically(element.(*ssa.Panic))
	case *ssa.Parameter:
		return interpreter.interpretParameterDynamically(element.(*ssa.Parameter))
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
	case *ssa.Type:
		return interpreter.interpretTypeDynamically(element.(*ssa.Type))
	case *ssa.TypeAssert:
		return interpreter.interpretTypeAssertDynamically(element.(*ssa.TypeAssert))
	case *ssa.UnOp:
		return interpreter.interpretUnOpDynamically(element.(*ssa.UnOp))
	case *ssa.BasicBlock:
		return interpreter.interpretBasicBlockDynamically(element.(*ssa.BasicBlock))
	default:
		panic("Unexpected element")
	}
}

func (interpreter *DynamicInterpreter) interpretAllocDynamically(element *ssa.Alloc) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretBinOpDynamically(element *ssa.BinOp) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretBuiltinDynamically(element *ssa.Builtin) []DynamicInterpreter {
	panic("TODO")
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

func (interpreter *DynamicInterpreter) interpretConstDynamically(element *ssa.Const) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretConvertDynamically(element *ssa.Convert) []DynamicInterpreter {
	panic("TODO")
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

func (interpreter *DynamicInterpreter) interpretFreeVarDynamically(element *ssa.FreeVar) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretGlobalDynamically(element *ssa.Global) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretGoDynamically(element *ssa.Go) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretIfDynamically(element *ssa.If) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretIndexDynamically(element *ssa.Index) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretIndexAddrDynamically(element *ssa.IndexAddr) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretJumpDynamically(element *ssa.Jump) []DynamicInterpreter {
	panic("TODO")
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

func (interpreter *DynamicInterpreter) interpretNamedConstDynamically(element *ssa.NamedConst) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretNextDynamically(element *ssa.Next) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretPanicDynamically(element *ssa.Panic) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretParameterDynamically(element *ssa.Parameter) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretPhiDynamically(element *ssa.Phi) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretRangeDynamically(element *ssa.Range) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretReturnDynamically(element *ssa.Return) []DynamicInterpreter {
	panic("TODO")
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

func (interpreter *DynamicInterpreter) interpretTypeDynamically(element *ssa.Type) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretTypeAssertDynamically(element *ssa.TypeAssert) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretUnOpDynamically(element *ssa.UnOp) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretBasicBlockDynamically(element *ssa.BasicBlock) []DynamicInterpreter {
	panic("TODO")
}
