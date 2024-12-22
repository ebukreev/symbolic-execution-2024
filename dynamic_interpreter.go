package se

import (
	"go/types"
	"golang.org/x/tools/go/ssa"
	"slices"
	"strconv"
	"strings"
	"symbolic-execution-2024/z3"
)

type CallStackFrame struct {
	Function        *ssa.Function
	Memory          map[string]SymbolicExpression
	ReturnValue     SymbolicExpression
	Instructions    []ssa.Instruction
	InstructionsPtr int
	BlocksStack     []*ssa.BasicBlock
}

type DynamicInterpreter struct {
	CallStack     []CallStackFrame
	Analyser      *Analyser
	PathCondition SymbolicExpression
	Heap          *SymbolicMemory
}

func (interpreter *DynamicInterpreter) CurrentFrame() *CallStackFrame {
	return &interpreter.CallStack[len(interpreter.CallStack)-1]
}

func copyMap[K string | int, V any](m map[K]V) map[K]V {
	cp := make(map[K]V)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func copySlice(slice []CallStackFrame) []CallStackFrame {
	tmp := make([]CallStackFrame, len(slice))
	for i, el := range slice {
		blocks := make([]*ssa.BasicBlock, len(el.BlocksStack))
		copy(blocks, el.BlocksStack)
		tmp[i] = CallStackFrame{
			el.Function,
			copyMap(el.Memory),
			el.ReturnValue,
			el.Instructions,
			el.InstructionsPtr,
			blocks,
		}
	}
	return tmp
}

func (interpreter *DynamicInterpreter) copy() DynamicInterpreter {
	return DynamicInterpreter{
		copySlice(interpreter.CallStack),
		interpreter.Analyser,
		interpreter.PathCondition,
		interpreter.Heap.Copy(),
	}
}

func InterpretDynamically(interpreter DynamicInterpreter) []DynamicInterpreter {
	if len(interpreter.CallStack) == 1 && interpreter.CurrentFrame().Instructions == nil {
		interpreter.CurrentFrame().BlocksStack = []*ssa.BasicBlock{interpreter.CurrentFrame().Function.Blocks[0]}
		interpreter.CurrentFrame().Instructions = interpreter.CurrentFrame().Function.Blocks[0].Instrs
		interpreter.CurrentFrame().InstructionsPtr = 0
	}

	front := interpreter.CurrentFrame().Instructions[interpreter.CurrentFrame().InstructionsPtr]
	interpreter.CurrentFrame().InstructionsPtr++
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
	res, ok := interpreter.CurrentFrame().Memory[value.Name()]
	if ok {
		return res
	}

	switch value.(type) {
	case *ssa.Const:
		res = interpreter.resolveConst(value.(*ssa.Const))
	case *ssa.BinOp:
		res = interpreter.resolveBinOp(value.(*ssa.BinOp))
		interpreter.CurrentFrame().Memory[value.Name()] = res
	case *ssa.Parameter:
		currentValue := interpreter.CurrentFrame().Memory[value.Name()]
		if currentValue != nil {
			return currentValue
		}
		res = interpreter.resolveParameter(value.(*ssa.Parameter))
		interpreter.CurrentFrame().Memory[value.Name()] = res
	case *ssa.Convert:
		res = interpreter.resolveConvert(value.(*ssa.Convert))
	case *ssa.Phi:
		res = interpreter.resolvePhi(value.(*ssa.Phi))
	case *ssa.Call:
		currentValue := interpreter.CurrentFrame().Memory[value.Name()]
		if currentValue != nil {
			return currentValue
		}
		res = interpreter.resolveCall(value.(*ssa.Call))
		_, ok := res.(*ssa.Function)
		if !ok {
			_, ok = res.(*InputValue)
			if ok {
				res.(*InputValue).Name = value.Name()
			}
			interpreter.CurrentFrame().Memory[value.Name()] = res
		}
	case *ssa.IndexAddr:
		res = interpreter.resolveIndexAddr(value.(*ssa.IndexAddr))
	case *ssa.UnOp:
		res = interpreter.resolveUnOp(value.(*ssa.UnOp))
	case *ssa.FieldAddr:
		res = interpreter.resolveFieldAddr(value.(*ssa.FieldAddr))
	case *ssa.MakeInterface:
		res = interpreter.resolveMakeInterface(value.(*ssa.MakeInterface))
	case *ssa.Alloc:
		res = interpreter.resolveAlloc(value.(*ssa.Alloc))
	case *ssa.Slice:
		res = interpreter.resolveSlice(value.(*ssa.Slice))
	default:
		panic("Unexpected element")
	}

	interpreter.CurrentFrame().Memory[value.Name()] = res

	return res
}

func (interpreter *DynamicInterpreter) resolveConst(element *ssa.Const) SymbolicExpression {
	switch element.Type().String() {
	case "int":
		value, _ := strconv.Atoi(element.Value.String())
		return &Literal[int]{value}
	case "uint":
		value, _ := strconv.ParseUint(element.Value.String(), 10, 64)
		return &Literal[uint]{uint(value)}
	case "float64":
		value, _ := strconv.ParseFloat(element.Value.String(), 64)
		return &Literal[float64]{value}
	case "bool":
		value, _ := strconv.ParseBool(element.Value.String())
		return &Literal[bool]{value}
	case "string":
		// TODO strings are not supported now
		return &Literal[float64]{1.0}
	}
	panic("unexpected const")
}

func (interpreter *DynamicInterpreter) resolveBinOp(element *ssa.BinOp) SymbolicExpression {
	left := interpreter.resolveExpression(element.X)
	right := interpreter.resolveExpression(element.Y)

	switch element.Op.String() {
	case "+":
		return CreateAdd(left, right)
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
	value, ok := interpreter.CurrentFrame().Memory[element.Name()]
	if ok {
		return value
	}
	ptr, isPtr := element.Type().(*types.Pointer)
	_, isBasic := element.Type().(*types.Basic)
	var ref SymbolicExpression
	if isPtr {
		named, ok := ptr.Elem().(*types.Named)
		if ok {
			structure := named.Underlying().(*types.Struct)
			fieldsNum := structure.NumFields()
			fields := make(map[int]string, fieldsNum)
			for i := 0; i < fieldsNum; i++ {
				fields[i] = structure.Field(i).Type().String()
			}
			ref = interpreter.Heap.Deref(interpreter.Heap.AllocateStruct(named.String(), fields))
		} else {
			ref = interpreter.Heap.Deref(interpreter.Heap.Allocate(element.Type().String()))
		}
	} else if element.Type().String() == "complex128" {
		ref = interpreter.Heap.AllocateArray("float64")
	} else if isBasic {
		ref = &InputValue{Name: element.Name(), Type: element.Type().String()}
	} else {
		ref = interpreter.Heap.Allocate(element.Type().String())
	}
	interpreter.CurrentFrame().Memory[element.Name()] = ref
	return ref
}

func (interpreter *DynamicInterpreter) resolveConvert(element *ssa.Convert) SymbolicExpression {
	return &Cast{interpreter.resolveExpression(element.X), element.Type().String()}
}

func (interpreter *DynamicInterpreter) resolvePhi(element *ssa.Phi) SymbolicExpression {
	for j := len(interpreter.CurrentFrame().BlocksStack) - 2; j >= 0; j-- {
		for i, pred := range element.Block().Preds {
			if pred == interpreter.CurrentFrame().BlocksStack[j] {
				edgeValue := interpreter.CurrentFrame().Memory[element.Edges[i].Name()]
				if edgeValue != nil {
					return edgeValue
				}

				return interpreter.resolveExpression(element.Edges[i])
			}
		}
	}

	panic("unexpected state")
}

func (interpreter *DynamicInterpreter) resolveIndexAddr(element *ssa.IndexAddr) SymbolicExpression {
	array := interpreter.resolveExpression(element.X)
	index := interpreter.resolveExpression(element.Index)
	return interpreter.Heap.GetFromArray(array, index)
}

func (interpreter *DynamicInterpreter) resolveUnOp(element *ssa.UnOp) SymbolicExpression {
	operand := interpreter.resolveExpression(element.X)
	switch element.Op.String() {
	case "*":
		return operand
	case "-":
		return &BinaryOperation{Left: &Literal[float64]{0.0}, Right: operand, Type: Sub}
	}
	panic("TODO")
}

func (interpreter *DynamicInterpreter) resolveFieldAddr(element *ssa.FieldAddr) SymbolicExpression {
	receiver := interpreter.resolveExpression(element.X)
	return interpreter.Heap.GetField(receiver, element.Field)
}

func (interpreter *DynamicInterpreter) resolveMakeInterface(element *ssa.MakeInterface) SymbolicExpression {
	return interpreter.resolveExpression(element.X)
}

func (interpreter *DynamicInterpreter) resolveAlloc(element *ssa.Alloc) SymbolicExpression {
	elementType := element.Type().(*types.Pointer).Elem()
	var result SymbolicExpression
	switch elementType.(type) {
	case *types.Array:
		componentType := elementType.(*types.Array).Elem()
		result = interpreter.Heap.AllocateArray(componentType.String())
	case *types.Named:
		structure := elementType.Underlying().(*types.Struct)
		fieldsNum := structure.NumFields()
		fields := make(map[int]string, fieldsNum)
		for i := 0; i < fieldsNum; i++ {
			fields[i] = structure.Field(i).Type().String()
		}
		result = interpreter.Heap.AllocateStruct(elementType.String(), fields)
		interpreter.CurrentFrame().Memory[element.Comment] = result
	default:
		panic("unexpected type")
	}
	return result
}

func (interpreter *DynamicInterpreter) resolveSlice(element *ssa.Slice) SymbolicExpression {
	return interpreter.resolveExpression(element.X)
}

func (interpreter *DynamicInterpreter) resolveCall(element *ssa.Call) SymbolicExpression {
	signature := strings.ReplaceAll(element.Call.Value.String(), " ", "_") + "("
	var argValues []SymbolicExpression

	for _, arg := range element.Call.Args {
		argValues = append(argValues, interpreter.resolveExpression(arg))
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

	supportedSignatures := []string{
		"builtin_len(Type)",
		"math.Sqrt(float64)",
		"math.IsNaN(float64)",
	}

	approximations := map[string]string{"math.Sqrt(float64)": "main.sqrt(float64)"}
	if approximation, found := approximations[signature]; found {
		signature = approximation
	}

	if strings.HasPrefix(signature, "main.makeSymbolic") {
		return &InputValue{Type: signature[18:strings.Index(signature, "]")]}
	} else if signature == "main.assume(any)" {
		cond := argValues[0]
		interpreter.PathCondition = CreateAnd(interpreter.PathCondition, cond)
		return cond
	} else if signature == "builtin_real(ComplexType)" {
		return interpreter.Heap.GetField(argValues[0], 0)
	} else if signature == "builtin_imag(ComplexType)" {
		return interpreter.Heap.GetField(argValues[0], 1)
	} else if slices.Contains(supportedSignatures, signature) {
		return &FunctionCall{Signature: signature, Arguments: argValues}
	} else {
		functionDecl := interpreter.Analyser.ResolveFunctionDeclaration(signature)
		memory := make(map[string]SymbolicExpression)

		for i, arg := range argValues {
			memory[functionDecl.Params[i].Name()] = arg
		}

		interpreter.CallStack = append(interpreter.CallStack,
			CallStackFrame{
				Function:        functionDecl,
				Memory:          memory,
				Instructions:    functionDecl.Blocks[0].Instrs,
				InstructionsPtr: 0,
				BlocksStack:     []*ssa.BasicBlock{functionDecl.Blocks[0]},
			},
		)

		return functionDecl
	}
}

func (interpreter *DynamicInterpreter) interpretAllocDynamically(element *ssa.Alloc) []DynamicInterpreter {
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretBinOpDynamically(element *ssa.BinOp) []DynamicInterpreter {
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretCallDynamically(element *ssa.Call) []DynamicInterpreter {
	if interpreter.CurrentFrame().ReturnValue != nil {
		interpreter.CurrentFrame().Memory[element.Name()] = interpreter.CurrentFrame().ReturnValue
		interpreter.CurrentFrame().ReturnValue = nil
		return []DynamicInterpreter{*interpreter}
	}

	interpreter.resolveExpression(element)

	return []DynamicInterpreter{*interpreter}
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
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
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
	thenBranch.CurrentFrame().BlocksStack = append(thenBranch.CurrentFrame().BlocksStack, successors[0])
	thenBranch.CurrentFrame().Instructions = successors[0].Instrs
	thenBranch.CurrentFrame().InstructionsPtr = 0

	elseBranch.PathCondition = CreateAnd(elseBranch.PathCondition, &Not{cond})
	elseBranch.CurrentFrame().BlocksStack = append(elseBranch.CurrentFrame().BlocksStack, successors[1])
	elseBranch.CurrentFrame().Instructions = successors[1].Instrs
	elseBranch.CurrentFrame().InstructionsPtr = 0

	var res []DynamicInterpreter
	if interpreter.Analyser.checkCondition(thenBranch.PathCondition) {
		res = append(res, thenBranch)
	}
	if interpreter.Analyser.checkCondition(elseBranch.PathCondition) {
		res = append(res, elseBranch)
	}

	return res
}

func (interpreter *DynamicInterpreter) interpretIndexDynamically(element *ssa.Index) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretIndexAddrDynamically(element *ssa.IndexAddr) []DynamicInterpreter {
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretJumpDynamically(element *ssa.Jump) []DynamicInterpreter {
	interpreter.CurrentFrame().BlocksStack = append(interpreter.CurrentFrame().BlocksStack, element.Block().Succs[0])
	interpreter.CurrentFrame().Instructions = element.Block().Succs[0].Instrs
	interpreter.CurrentFrame().InstructionsPtr = 0
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
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
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
	interpreter.CurrentFrame().ReturnValue = interpreter.resolveExpression(element.Results[0])
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
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretSliceToArrayPointerDynamically(element *ssa.SliceToArrayPointer) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretStoreDynamically(element *ssa.Store) []DynamicInterpreter {
	value := interpreter.resolveExpression(element.Val)
	indexAddr, ok := element.Addr.(*ssa.IndexAddr)
	if ok {
		array := interpreter.resolveExpression(indexAddr.X)
		index := interpreter.resolveExpression(indexAddr.Index)

		interpreter.Heap.AssignToArray(array, index, value)
	}

	fieldAddr, ok := element.Addr.(*ssa.FieldAddr)
	if ok {
		field := interpreter.resolveExpression(fieldAddr.X)
		interpreter.Heap.AssignField(field, fieldAddr.Field, value)
	}

	_, param := element.Val.(*ssa.Parameter)
	_, structure := element.Val.Type().Underlying().(*types.Struct)
	if param && structure {
		interpreter.CurrentFrame().Memory[element.Addr.Name()] = value
	}

	return []DynamicInterpreter{*interpreter}
}

func (interpreter *DynamicInterpreter) interpretTypeAssertDynamically(element *ssa.TypeAssert) []DynamicInterpreter {
	panic("TODO")
}

func (interpreter *DynamicInterpreter) interpretUnOpDynamically(element *ssa.UnOp) []DynamicInterpreter {
	interpreter.resolveExpression(element)
	return []DynamicInterpreter{*interpreter}
}

func (analyser *Analyser) checkCondition(condition SymbolicExpression) bool {
	smt := analyser.SmtBuilder.BuildSmt(condition)[0].(z3.Bool)
	analyser.Solver.SmtSolver.Assert(smt)
	sat, err := analyser.Solver.SmtSolver.Check()
	analyser.Solver.SmtSolver.Reset()
	return sat && err == nil
}
