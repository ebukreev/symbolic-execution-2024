package main

import (
	"golang.org/x/tools/go/ssa"
)

func Interpret(element interface{}) *SymbolicExpression {
	switch element.(type) {
	case *ssa.Package:
		return interpretPackage(element.(*ssa.Package))
	case *ssa.Alloc:
		return interpretAlloc(element.(*ssa.Alloc))
	case *ssa.BinOp:
		return interpretBinOp(element.(*ssa.BinOp))
	case *ssa.Builtin:
		return interpretBuiltin(element.(*ssa.Builtin))
	case *ssa.Call:
		return interpretCall(element.(*ssa.Call))
	case *ssa.ChangeInterface:
		return interpretChangeInterface(element.(*ssa.ChangeInterface))
	case *ssa.ChangeType:
		return interpretChangeType(element.(*ssa.ChangeType))
	case *ssa.Const:
		return interpretConst(element.(*ssa.Const))
	case *ssa.Convert:
		return interpretConvert(element.(*ssa.Convert))
	case *ssa.DebugRef:
		return interpretDebugRef(element.(*ssa.DebugRef))
	case *ssa.Defer:
		return interpretDefer(element.(*ssa.Defer))
	case *ssa.Extract:
		return interpretExtract(element.(*ssa.Extract))
	case *ssa.Field:
		return interpretField(element.(*ssa.Field))
	case *ssa.FieldAddr:
		return interpretFieldAddr(element.(*ssa.FieldAddr))
	case *ssa.FreeVar:
		return interpretFreeVar(element.(*ssa.FreeVar))
	case *ssa.Function:
		return interpretFunction(element.(*ssa.Function))
	case *ssa.Global:
		return interpretGlobal(element.(*ssa.Global))
	case *ssa.Go:
		return interpretGo(element.(*ssa.Go))
	case *ssa.If:
		return interpretIf(element.(*ssa.If))
	case *ssa.Index:
		return interpretIndex(element.(*ssa.Index))
	case *ssa.IndexAddr:
		return interpretIndexAddr(element.(*ssa.IndexAddr))
	case *ssa.Jump:
		return interpretJump(element.(*ssa.Jump))
	case *ssa.Lookup:
		return interpretLookup(element.(*ssa.Lookup))
	case *ssa.MakeChan:
		return interpretMakeChan(element.(*ssa.MakeChan))
	case *ssa.MakeClosure:
		return interpretMakeClosure(element.(*ssa.MakeClosure))
	case *ssa.MakeInterface:
		return interpretMakeInterface(element.(*ssa.MakeInterface))
	case *ssa.MakeMap:
		return interpretMakeMap(element.(*ssa.MakeMap))
	case *ssa.MakeSlice:
		return interpretMakeSlice(element.(*ssa.MakeSlice))
	case *ssa.MapUpdate:
		return interpretMapUpdate(element.(*ssa.MapUpdate))
	case *ssa.MultiConvert:
		return interpretMultiConvert(element.(*ssa.MultiConvert))
	case *ssa.NamedConst:
		return interpretNamedConst(element.(*ssa.NamedConst))
	case *ssa.Next:
		return interpretNext(element.(*ssa.Next))
	case *ssa.Panic:
		return interpretPanic(element.(*ssa.Panic))
	case *ssa.Parameter:
		return interpretParameter(element.(*ssa.Parameter))
	case *ssa.Phi:
		return interpretPhi(element.(*ssa.Phi))
	case *ssa.Range:
		return interpretRange(element.(*ssa.Range))
	case *ssa.Return:
		return interpretReturn(element.(*ssa.Return))
	case *ssa.RunDefers:
		return interpretRunDefers(element.(*ssa.RunDefers))
	case *ssa.Select:
		return interpretSelect(element.(*ssa.Select))
	case *ssa.Send:
		return interpretSend(element.(*ssa.Send))
	case *ssa.Slice:
		return interpretSlice(element.(*ssa.Slice))
	case *ssa.SliceToArrayPointer:
		return interpretSliceToArrayPointer(element.(*ssa.SliceToArrayPointer))
	case *ssa.Store:
		return interpretStore(element.(*ssa.Store))
	case *ssa.Type:
		return interpretType(element.(*ssa.Type))
	case *ssa.TypeAssert:
		return interpretTypeAssert(element.(*ssa.TypeAssert))
	case *ssa.UnOp:
		return interpretUnOp(element.(*ssa.UnOp))
	default:
		panic("Unexpected element")
	}
}

func interpretPackage(element *ssa.Package) *SymbolicExpression {
	for _, member := range element.Members {
		Interpret(member)
	}
	return nil
}

func interpretAlloc(element *ssa.Alloc) *SymbolicExpression {
	panic("TODO")
}

func interpretBinOp(element *ssa.BinOp) *SymbolicExpression {
	panic("TODO")
}

func interpretBuiltin(element *ssa.Builtin) *SymbolicExpression {
	panic("TODO")
}

func interpretCall(element *ssa.Call) *SymbolicExpression {
	panic("TODO")
}

func interpretChangeInterface(element *ssa.ChangeInterface) *SymbolicExpression {
	panic("TODO")
}

func interpretChangeType(element *ssa.ChangeType) *SymbolicExpression {
	panic("TODO")
}

func interpretConst(element *ssa.Const) *SymbolicExpression {
	panic("TODO")
}

func interpretConvert(element *ssa.Convert) *SymbolicExpression {
	panic("TODO")
}

func interpretDebugRef(element *ssa.DebugRef) *SymbolicExpression {
	panic("TODO")
}

func interpretDefer(element *ssa.Defer) *SymbolicExpression {
	panic("TODO")
}

func interpretExtract(element *ssa.Extract) *SymbolicExpression {
	panic("TODO")
}

func interpretField(element *ssa.Field) *SymbolicExpression {
	panic("TODO")
}

func interpretFieldAddr(element *ssa.FieldAddr) *SymbolicExpression {
	panic("TODO")
}

func interpretFreeVar(element *ssa.FreeVar) *SymbolicExpression {
	panic("TODO")
}

func interpretFunction(element *ssa.Function) *SymbolicExpression {
	panic("TODO")
}

func interpretGlobal(element *ssa.Global) *SymbolicExpression {
	panic("TODO")
}

func interpretGo(element *ssa.Go) *SymbolicExpression {
	panic("TODO")
}

func interpretIf(element *ssa.If) *SymbolicExpression {
	panic("TODO")
}

func interpretIndex(element *ssa.Index) *SymbolicExpression {
	panic("TODO")
}

func interpretIndexAddr(element *ssa.IndexAddr) *SymbolicExpression {
	panic("TODO")
}

func interpretJump(element *ssa.Jump) *SymbolicExpression {
	panic("TODO")
}

func interpretLookup(element *ssa.Lookup) *SymbolicExpression {
	panic("TODO")
}

func interpretMakeChan(element *ssa.MakeChan) *SymbolicExpression {
	panic("TODO")
}

func interpretMakeClosure(element *ssa.MakeClosure) *SymbolicExpression {
	panic("TODO")
}

func interpretMakeInterface(element *ssa.MakeInterface) *SymbolicExpression {
	panic("TODO")
}

func interpretMakeMap(element *ssa.MakeMap) *SymbolicExpression {
	panic("TODO")
}

func interpretMakeSlice(element *ssa.MakeSlice) *SymbolicExpression {
	panic("TODO")
}

func interpretMapUpdate(element *ssa.MapUpdate) *SymbolicExpression {
	panic("TODO")
}

func interpretMultiConvert(element *ssa.MultiConvert) *SymbolicExpression {
	panic("TODO")
}

func interpretNamedConst(element *ssa.NamedConst) *SymbolicExpression {
	panic("TODO")
}

func interpretNext(element *ssa.Next) *SymbolicExpression {
	panic("TODO")
}

func interpretPanic(element *ssa.Panic) *SymbolicExpression {
	panic("TODO")
}

func interpretParameter(element *ssa.Parameter) *SymbolicExpression {
	panic("TODO")
}

func interpretPhi(element *ssa.Phi) *SymbolicExpression {
	panic("TODO")
}

func interpretRange(element *ssa.Range) *SymbolicExpression {
	panic("TODO")
}

func interpretReturn(element *ssa.Return) *SymbolicExpression {
	panic("TODO")
}

func interpretRunDefers(element *ssa.RunDefers) *SymbolicExpression {
	panic("TODO")
}

func interpretSelect(element *ssa.Select) *SymbolicExpression {
	panic("TODO")
}

func interpretSend(element *ssa.Send) *SymbolicExpression {
	panic("TODO")
}

func interpretSlice(element *ssa.Slice) *SymbolicExpression {
	panic("TODO")
}

func interpretSliceToArrayPointer(element *ssa.SliceToArrayPointer) *SymbolicExpression {
	panic("TODO")
}

func interpretStore(element *ssa.Store) *SymbolicExpression {
	panic("TODO")
}

func interpretType(element *ssa.Type) *SymbolicExpression {
	panic("TODO")
}

func interpretTypeAssert(element *ssa.TypeAssert) *SymbolicExpression {
	panic("TODO")
}

func interpretUnOp(element *ssa.UnOp) *SymbolicExpression {
	panic("TODO")
}
