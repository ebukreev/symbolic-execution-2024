package main

import (
	"go/types"
	"golang.org/x/tools/go/ssa"
	"strconv"
	"strings"
)

type Interpreter struct {
	PathCondition SymbolicExpression
	Memory        map[string]*Conditional
	ReturnValue   Conditional
	BlocksStack   []*ssa.BasicBlock
}

func Interpret(function *ssa.Function) Conditional {
	interpreter := Interpreter{PathCondition: &Literal[bool]{true},
		Memory:      map[string]*Conditional{},
		ReturnValue: Conditional{make(map[SymbolicExpression]SymbolicExpression)},
		BlocksStack: []*ssa.BasicBlock{},
	}

	interpreter.interpret(function.Blocks[0])

	return interpreter.ReturnValue
}

func (interpreter *Interpreter) interpret(element interface{}) SymbolicExpression {
	switch element.(type) {
	case *ssa.Alloc:
		return interpreter.interpretAlloc(element.(*ssa.Alloc))
	case *ssa.BinOp:
		return interpreter.interpretBinOp(element.(*ssa.BinOp))
	case *ssa.Builtin:
		return interpreter.interpretBuiltin(element.(*ssa.Builtin))
	case *ssa.Call:
		return interpreter.interpretCall(element.(*ssa.Call))
	case *ssa.ChangeInterface:
		return interpreter.interpretChangeInterface(element.(*ssa.ChangeInterface))
	case *ssa.ChangeType:
		return interpreter.interpretChangeType(element.(*ssa.ChangeType))
	case *ssa.Const:
		return interpreter.interpretConst(element.(*ssa.Const))
	case *ssa.Convert:
		return interpreter.interpretConvert(element.(*ssa.Convert))
	case *ssa.DebugRef:
		return interpreter.interpretDebugRef(element.(*ssa.DebugRef))
	case *ssa.Defer:
		return interpreter.interpretDefer(element.(*ssa.Defer))
	case *ssa.Extract:
		return interpreter.interpretExtract(element.(*ssa.Extract))
	case *ssa.Field:
		return interpreter.interpretField(element.(*ssa.Field))
	case *ssa.FieldAddr:
		return interpreter.interpretFieldAddr(element.(*ssa.FieldAddr))
	case *ssa.FreeVar:
		return interpreter.interpretFreeVar(element.(*ssa.FreeVar))
	case *ssa.Global:
		return interpreter.interpretGlobal(element.(*ssa.Global))
	case *ssa.Go:
		return interpreter.interpretGo(element.(*ssa.Go))
	case *ssa.If:
		return interpreter.interpretIf(element.(*ssa.If))
	case *ssa.Index:
		return interpreter.interpretIndex(element.(*ssa.Index))
	case *ssa.IndexAddr:
		return interpreter.interpretIndexAddr(element.(*ssa.IndexAddr))
	case *ssa.Jump:
		return interpreter.interpretJump(element.(*ssa.Jump))
	case *ssa.Lookup:
		return interpreter.interpretLookup(element.(*ssa.Lookup))
	case *ssa.MakeChan:
		return interpreter.interpretMakeChan(element.(*ssa.MakeChan))
	case *ssa.MakeClosure:
		return interpreter.interpretMakeClosure(element.(*ssa.MakeClosure))
	case *ssa.MakeInterface:
		return interpreter.interpretMakeInterface(element.(*ssa.MakeInterface))
	case *ssa.MakeMap:
		return interpreter.interpretMakeMap(element.(*ssa.MakeMap))
	case *ssa.MakeSlice:
		return interpreter.interpretMakeSlice(element.(*ssa.MakeSlice))
	case *ssa.MapUpdate:
		return interpreter.interpretMapUpdate(element.(*ssa.MapUpdate))
	case *ssa.MultiConvert:
		return interpreter.interpretMultiConvert(element.(*ssa.MultiConvert))
	case *ssa.NamedConst:
		return interpreter.interpretNamedConst(element.(*ssa.NamedConst))
	case *ssa.Next:
		return interpreter.interpretNext(element.(*ssa.Next))
	case *ssa.Panic:
		return interpreter.interpretPanic(element.(*ssa.Panic))
	case *ssa.Parameter:
		return interpreter.interpretParameter(element.(*ssa.Parameter))
	case *ssa.Phi:
		return interpreter.interpretPhi(element.(*ssa.Phi))
	case *ssa.Range:
		return interpreter.interpretRange(element.(*ssa.Range))
	case *ssa.Return:
		return interpreter.interpretReturn(element.(*ssa.Return))
	case *ssa.RunDefers:
		return interpreter.interpretRunDefers(element.(*ssa.RunDefers))
	case *ssa.Select:
		return interpreter.interpretSelect(element.(*ssa.Select))
	case *ssa.Send:
		return interpreter.interpretSend(element.(*ssa.Send))
	case *ssa.Slice:
		return interpreter.interpretSlice(element.(*ssa.Slice))
	case *ssa.SliceToArrayPointer:
		return interpreter.interpretSliceToArrayPointer(element.(*ssa.SliceToArrayPointer))
	case *ssa.Store:
		return interpreter.interpretStore(element.(*ssa.Store))
	case *ssa.Type:
		return interpreter.interpretType(element.(*ssa.Type))
	case *ssa.TypeAssert:
		return interpreter.interpretTypeAssert(element.(*ssa.TypeAssert))
	case *ssa.UnOp:
		return interpreter.interpretUnOp(element.(*ssa.UnOp))
	case *ssa.BasicBlock:
		return interpreter.interpretBasicBlock(element.(*ssa.BasicBlock))
	default:
		panic("Unexpected element")
	}
}

func (interpreter *Interpreter) interpretAlloc(element *ssa.Alloc) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretBinOp(element *ssa.BinOp) SymbolicExpression {
	left := interpreter.interpret(element.X)
	right := interpreter.interpret(element.Y)

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

func (interpreter *Interpreter) interpretBuiltin(element *ssa.Builtin) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretCall(element *ssa.Call) SymbolicExpression {
	signature := strings.ReplaceAll(element.Call.Value.String(), " ", "_") + "("
	var argValues []SymbolicExpression

	for _, arg := range element.Call.Args {
		argValues = append(argValues, interpreter.interpret(arg))
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

func (interpreter *Interpreter) interpretChangeInterface(element *ssa.ChangeInterface) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretChangeType(element *ssa.ChangeType) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretConst(element *ssa.Const) SymbolicExpression {
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

func (interpreter *Interpreter) interpretConvert(element *ssa.Convert) SymbolicExpression {
	return &Cast{interpreter.interpret(element.X), element.Type().String()}
}

func (interpreter *Interpreter) interpretDebugRef(element *ssa.DebugRef) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretDefer(element *ssa.Defer) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretExtract(element *ssa.Extract) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretField(element *ssa.Field) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretFieldAddr(element *ssa.FieldAddr) SymbolicExpression {
	receiver := interpreter.interpret(element.X)
	typeSignature := getTypeSignature(element.X.Type())
	return &FunctionCall{typeSignature + "_" + strconv.Itoa(element.Field), []SymbolicExpression{receiver}}
}

func (interpreter *Interpreter) interpretFreeVar(element *ssa.FreeVar) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretGlobal(element *ssa.Global) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretGo(element *ssa.Go) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretIf(element *ssa.If) SymbolicExpression {
	cond := interpreter.interpret(element.Cond)
	successors := element.Block().Succs
	enterBranch(*interpreter, cond, successors[0])
	if len(successors) == 1 {
		return nil
	}
	notCond := &Not{cond}
	enterBranch(*interpreter, notCond, successors[1])
	return nil
}

func (interpreter *Interpreter) interpretIndex(element *ssa.Index) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretIndexAddr(element *ssa.IndexAddr) SymbolicExpression {
	array := interpreter.interpret(element.X)
	index := interpreter.interpret(element.Index)
	return &ArrayAccess{array, index}
}

func (interpreter *Interpreter) interpretJump(element *ssa.Jump) SymbolicExpression {
	return interpreter.interpret(element.Block().Succs[0])
}

func (interpreter *Interpreter) interpretLookup(element *ssa.Lookup) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretMakeChan(element *ssa.MakeChan) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretMakeClosure(element *ssa.MakeClosure) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretMakeInterface(element *ssa.MakeInterface) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretMakeMap(element *ssa.MakeMap) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretMakeSlice(element *ssa.MakeSlice) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretMapUpdate(element *ssa.MapUpdate) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretMultiConvert(element *ssa.MultiConvert) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretNamedConst(element *ssa.NamedConst) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretNext(element *ssa.Next) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretPanic(element *ssa.Panic) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretParameter(element *ssa.Parameter) SymbolicExpression {
	return &InputValue{Name: element.Name(), Type: element.Type().String()}
}

func (interpreter *Interpreter) interpretPhi(element *ssa.Phi) SymbolicExpression {
	for i, pred := range element.Block().Preds {
		for j := len(interpreter.BlocksStack) - 2; j >= 0; j-- {
			if pred == interpreter.BlocksStack[j] {
				elementValue, ok := interpreter.Memory[element.Comment]
				edgeValue := interpreter.interpret(element.Edges[i])
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

func (interpreter *Interpreter) interpretRange(element *ssa.Range) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretReturn(element *ssa.Return) SymbolicExpression {
	returnExpr := interpreter.interpret(element.Results[0])
	interpreter.ReturnValue.Options[interpreter.PathCondition] = returnExpr
	return nil
}

func (interpreter *Interpreter) interpretRunDefers(element *ssa.RunDefers) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretSelect(element *ssa.Select) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretSend(element *ssa.Send) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretSlice(element *ssa.Slice) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretSliceToArrayPointer(element *ssa.SliceToArrayPointer) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretStore(element *ssa.Store) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretType(element *ssa.Type) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretTypeAssert(element *ssa.TypeAssert) SymbolicExpression {
	panic("TODO")
}

func (interpreter *Interpreter) interpretUnOp(element *ssa.UnOp) SymbolicExpression {
	operand := interpreter.interpret(element.X)
	switch element.Op.String() {
	case "*":
		return operand
	}
	panic("TODO")
}

func (interpreter *Interpreter) interpretBasicBlock(element *ssa.BasicBlock) SymbolicExpression {
	interpreter.BlocksStack = append(interpreter.BlocksStack, element)
	for _, instr := range element.Instrs {
		interpreter.interpret(instr)
	}
	interpreter.BlocksStack = interpreter.BlocksStack[:len(interpreter.BlocksStack)-1]
	return nil
}

func enterBranch(interpreter Interpreter, condition SymbolicExpression, body *ssa.BasicBlock) SymbolicExpression {
	interpreter.PathCondition = CreateAnd(interpreter.PathCondition, condition)
	interpreter.interpret(body)
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
