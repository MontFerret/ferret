package compiler

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Variable struct {
		Name     string
		FirstUse int // Instruction number of first use
		LastUse  int // Instruction number of last use
		Register runtime.Operand
		IsLive   bool
		Depth    int
	}

	loopScope struct {
		result      int
		passThrough bool
		position    int
	}

	visitor struct {
		*fql.BaseFqlParserVisitor
		err        error
		src        string
		emitter    *Emitter
		registers  *RegisterAllocator
		symbols    *SymbolTable
		catchTable [][2]int
	}
)

const (
	jumpPlaceholder      = -1
	undefinedVariable    = -1
	ignorePseudoVariable = "_"
	pseudoVariable       = "CURRENT"
	waitScope            = "waitfor"
	forScope             = "for"
)

func newVisitor(src string) *visitor {
	v := new(visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.src = src
	v.registers = NewRegisterAllocator()
	v.symbols = NewSymbolTable(v.registers)
	v.emitter = NewEmitter()
	v.catchTable = make([][2]int, 0)

	return v
}

func (v *visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	ctx.Body().Accept(v)

	return nil
}

func (v *visitor) VisitBody(ctx *fql.BodyContext) interface{} {
	for _, statement := range ctx.AllBodyStatement() {
		statement.Accept(v)
	}

	ctx.BodyExpression().Accept(v)

	return nil
}

func (v *visitor) VisitBodyStatement(ctx *fql.BodyStatementContext) interface{} {
	if c := ctx.VariableDeclaration(); c != nil {
		return c.Accept(v)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.WaitForExpression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitBodyExpression(ctx *fql.BodyExpressionContext) interface{} {
	if c := ctx.ForExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.ReturnExpression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}

func (v *visitor) VisitForExpression(ctx *fql.ForExpressionContext) interface{} {
	//v.beginScope()

	//var passThrough bool
	//var distinct bool
	//var returnRuleCtx antlr.RuleContext
	//var loopJump, exitJump int
	//// identify whether it's WHILE or FOR loop
	//isForInLoop := ctx.While() == nil
	//returnCtx := ctx.ForExpressionReturn()
	//
	//if c := returnCtx.ReturnExpression(); c != nil {
	//	returnRuleCtx = c
	//	distinct = c.Distinct() != nil
	//} else if c := returnCtx.ForExpression(); c != nil {
	//	returnRuleCtx = c
	//	passThrough = true
	//}
	//
	//v.beginLoopScope(passThrough, distinct)
	//
	//if isForInLoop {
	//	// Loop data source to iterate over
	//	if c := ctx.ForExpressionSource(); c != nil {
	//		c.Accept(v)
	//	}
	//
	//	v.emitter.EmitABC(runtime.OpForLoopInitInput)
	//	loopJump = len(v.instructions)
	//	v.emitter.EmitABC(runtime.OpForLoopHasNext)
	//	exitJump = v.emitJump(runtime.OpJumpIfFalse)
	//	// pop the boolean value from the stack
	//	v.emitPop()
	//
	//	valVar := ctx.GetValueVariable().GetText()
	//	counterVarCtx := ctx.GetCounterVariable()
	//
	//	hasValVar := valVar != ignorePseudoVariable
	//	var hasCounterVar bool
	//	var counterVar string
	//
	//	if counterVarCtx != nil {
	//		counterVar = counterVarCtx.GetText()
	//		hasCounterVar = true
	//	}
	//
	//	var valVarIndex int
	//
	//	// declare value variable
	//	if hasValVar {
	//		valVarIndex = v.declareVariable(valVar)
	//	}
	//
	//	var counterVarIndex int
	//
	//	if hasCounterVar {
	//		// declare counter variable
	//		counterVarIndex = v.declareVariable(counterVar)
	//	}
	//
	//	if hasValVar && hasCounterVar {
	//		// we will calculate the index of the counter variable
	//		v.emitter.EmitABC(runtime.OpForLoopNext)
	//	} else if hasValVar {
	//		v.emitter.EmitABC(runtime.OpForLoopNextValue)
	//	} else if hasCounterVar {
	//		v.emitter.EmitABC(runtime.OpForLoopNextCounter)
	//	} else {
	//		panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
	//	}
	//
	//	if hasValVar {
	//		v.defineVariable(valVarIndex)
	//	}
	//
	//	if hasCounterVar {
	//		v.defineVariable(counterVarIndex)
	//	}
	//} else {
	//	// Create initial value for the loop counter
	//	v.emitter.EmitABC(runtime.OpWhileLoopInitCounter)
	//
	//	loopJump = len(v.instructions)
	//
	//	// Condition expression
	//	ctx.Expression().Accept(v)
	//
	//	// Condition check
	//	exitJump = v.emitJump(runtime.OpJumpIfFalse)
	//	// pop the boolean value from the stack
	//	v.emitPop()
	//
	//	counterVar := ctx.GetCounterVariable().GetText()
	//
	//	// declare counter variable
	//	// and increment it by 1
	//	index := v.declareVariable(counterVar)
	//	v.emitter.EmitABC(runtime.OpWhileLoopNext)
	//	v.defineVariable(index)
	//}
	//
	//v.patchLoopScope(loopJump)
	//
	//// body
	//if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
	//	for _, b := range body {
	//		b.Accept(v)
	//	}
	//}
	//
	//// return
	//returnRuleCtx.Accept(v)
	//
	//v.emitLoop(loopJump)
	//v.patchJump(exitJump)
	//v.endScope()
	//// pop the boolean value from the stack
	//v.emitPop()
	//
	//if isForInLoop {
	//	// pop the iterator
	//	v.emitPopAndClose()
	//} else {
	//	// pop the counter
	//	v.emitPop()
	//}
	//
	//v.endLoopScope()

	return nil
}

func (v *visitor) VisitForExpressionSource(ctx *fql.ForExpressionSourceContext) interface{} {
	if c := ctx.FunctionCallExpression(); c != nil {
		c.Accept(v)
	} else if c := ctx.MemberExpression(); c != nil {
		c.Accept(v)
	} else if c := ctx.Variable(); c != nil {
		c.Accept(v)
	} else if c := ctx.Param(); c != nil {
		c.Accept(v)
	} else if c := ctx.RangeOperator(); c != nil {
		c.Accept(v)
	} else if c := ctx.ArrayLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.ObjectLiteral(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitForExpressionBody(ctx *fql.ForExpressionBodyContext) interface{} {
	if c := ctx.ForExpressionClause(); c != nil {
		c.Accept(v)
	}

	if c := ctx.ForExpressionStatement(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitForExpressionClause(ctx *fql.ForExpressionClauseContext) interface{} {
	if c := ctx.LimitClause(); c != nil {
		// TODO: Implement
		c.Accept(v)
	}

	if c := ctx.FilterClause(); c != nil {
		c.Accept(v)
	}

	if c := ctx.SortClause(); c != nil {
		// TODO: Implement
		c.Accept(v)
	}

	if c := ctx.CollectClause(); c != nil {
		// TODO: Implement
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitFilterClause(ctx *fql.FilterClauseContext) interface{} {
	//ctx.Expression().Accept(v)
	//fwd := v.emitJump(runtime.OpJumpIfTrue)
	//// Pop on false
	//v.emitPop()
	//// And jump back to the beginning of the loop
	//bwd := v.emitJump(runtime.OpJumpBackward)
	//v.patchJumpWith(bwd, len(v.instructions)-v.resolveLoopPosition())
	//
	//// Otherwise, pop on true
	//v.patchJump(fwd)
	//v.emitPop()

	return nil
}

func (v *visitor) VisitForExpressionStatement(ctx *fql.ForExpressionStatementContext) interface{} {
	//if c := ctx.VariableDeclaration(); c != nil {
	//	c.Accept(v)
	//} else if c := ctx.FunctionCallExpression(); c != nil {
	//	c.Accept(v)
	//	// remove un-used return value
	//	v.emitPop()
	//}

	return nil
}

func (v *visitor) VisitFunctionCallExpression(ctx *fql.FunctionCallExpressionContext) interface{} {
	//call := ctx.FunctionCall().(*fql.FunctionCallContext)
	//
	//var name string
	//
	//funcNS := call.Namespace()
	//
	//if funcNS != nil {
	//	name += funcNS.GetText()
	//}
	//
	//name += call.FunctionName().GetText()
	//
	//isNonOptional := ctx.ErrorOperator() == nil
	//
	//v.emitConstant(values.String(name))
	//
	//var size int
	//
	//if args := call.ArgumentList(); args != nil {
	//	out := v.VisitArgumentList(args.(*fql.ArgumentListContext))
	//	size = out.(int)
	//}
	//
	//switch size {
	//case 0:
	//	if isNonOptional {
	//		v.emitter.EmitABC(runtime.OpCall, 0)
	//	} else {
	//		v.emitter.EmitABC(runtime.OpCallSafe, 0)
	//	}
	//case 1:
	//	if isNonOptional {
	//		v.emitter.EmitABC(runtime.OpCall1, 1)
	//	} else {
	//		v.emitter.EmitABC(runtime.OpCall1Safe, 1)
	//	}
	//case 2:
	//	if isNonOptional {
	//		v.emitter.EmitABC(runtime.OpCall2, 2)
	//	} else {
	//		v.emitter.EmitABC(runtime.OpCall2Safe, 2)
	//	}
	//case 3:
	//	if isNonOptional {
	//		v.emitter.EmitABC(runtime.OpCall3, 3)
	//	} else {
	//		v.emitter.EmitABC(runtime.OpCall3Safe, 3)
	//	}
	//case 4:
	//	if isNonOptional {
	//		v.emitter.EmitABC(runtime.OpCall4, 4)
	//	} else {
	//		v.emitter.EmitABC(runtime.OpCall4Safe, 4)
	//	}
	//default:
	//	if isNonOptional {
	//		v.emitter.EmitABC(runtime.OpCallN, size)
	//	} else {
	//		v.emitter.EmitABC(runtime.OpCallNSafe, size)
	//	}
	//}

	return nil
}

func (v *visitor) VisitMemberExpression(ctx *fql.MemberExpressionContext) interface{} {
	mes := ctx.MemberExpressionSource().(*fql.MemberExpressionSourceContext)

	var mesOut interface{}

	if c := mes.Variable(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.Param(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.ObjectLiteral(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.ArrayLiteral(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.FunctionCall(); c != nil {
		mesOut = c.Accept(v)
	}

	var dst runtime.Operand
	src1 := v.toRegister(mesOut.(runtime.Operand))
	segments := ctx.AllMemberExpressionPath()

	for _, segment := range segments {
		var out2 interface{}
		p := segment.(*fql.MemberExpressionPathContext)

		if c := p.PropertyName(); c != nil {
			out2 = c.Accept(v)
		} else if c := p.ComputedPropertyName(); c != nil {
			out2 = c.Accept(v)
		}

		src2 := v.toRegister(out2.(runtime.Operand))
		dst = v.registers.Allocate(Temp)

		if p.ErrorOperator() != nil {
			v.emitter.EmitABC(runtime.OpLoadPropertyOptional, dst, src1, src2)
		} else {
			v.emitter.EmitABC(runtime.OpLoadProperty, dst, src1, src2)
		}

		src1 = dst
	}

	return dst
}

func (v *visitor) VisitRangeOperator(ctx *fql.RangeOperatorContext) interface{} {
	dst := v.registers.AllocateTempVar()
	start := v.toRegister(ctx.GetLeft().Accept(v).(runtime.Operand))
	end := v.toRegister(ctx.GetRight().Accept(v).(runtime.Operand))

	v.emitter.EmitABC(runtime.OpRange, dst, start, end)

	return dst
}

func (v *visitor) VisitRangeOperand(ctx *fql.RangeOperandContext) interface{} {
	if c := ctx.IntegerLiteral(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Param(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitVariableDeclaration(ctx *fql.VariableDeclarationContext) interface{} {
	name := ignorePseudoVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	src := ctx.Expression().Accept(v).(runtime.Operand)

	if name != ignorePseudoVariable {
		dest := v.symbols.DefineVariable(name)

		if src.IsConstant() {
			tmp := v.registers.Allocate(Temp)
			v.emitter.EmitAB(runtime.OpLoadConst, tmp, src)
			v.emitter.EmitAB(runtime.OpStoreGlobal, dest, tmp)
		} else if v.symbols.Scope() == 0 {
			v.emitter.EmitAB(runtime.OpStoreGlobal, dest, src)
		} else {
			v.emitter.EmitAB(runtime.OpMove, dest, src)
		}

		return dest
	}

	return nil
}

func (v *visitor) VisitVariable(ctx *fql.VariableContext) interface{} {
	// Just return the register / constant index
	op := v.symbols.LookupVariable(ctx.GetText())

	if op.IsRegister() {
		return op
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(runtime.OpLoadGlobal, reg, op)

	return reg
}

func (v *visitor) VisitArrayLiteral(ctx *fql.ArrayLiteralContext) interface{} {
	// Allocate destination register for the array
	destReg := v.registers.Allocate(Temp)

	if list := ctx.ArgumentList(); list != nil {
		// Get all array element expressions
		exps := list.(*fql.ArgumentListContext).AllExpression()
		size := len(exps)

		if size > 0 {
			// Allocate seq for array elements
			seq := v.registers.AllocateSequence(size, Temp)

			// Evaluate each element into seq registers
			for i, exp := range exps {
				// Compile expression and move to seq register
				srcReg := exp.Accept(v).(runtime.Operand)

				if srcReg.IsConstant() {
					v.emitter.EmitAB(runtime.OpLoadConst, seq.Registers[i], srcReg)
				} else {
					v.emitter.EmitAB(runtime.OpMove, seq.Registers[i], srcReg)
				}

				// Free source register if temporary
				if srcReg.IsRegister() {
					//v.registers.Free(srcReg)
				}
			}

			// Initialize an array
			v.emitter.EmitAs(runtime.OpArray, destReg, seq)

			// Free seq registers
			//v.registers.FreeSequence(seq)

			return destReg
		}
	}

	// Empty array
	v.emitter.EmitA(runtime.OpArray, destReg)

	return destReg
}

func (v *visitor) VisitObjectLiteral(ctx *fql.ObjectLiteralContext) interface{} {
	dst := v.registers.Allocate(Temp)
	assignments := ctx.AllPropertyAssignment()
	size := len(assignments)

	if size == 0 {
		v.emitter.EmitA(runtime.OpObject, dst)

		return dst
	}

	seq := v.registers.AllocateSequence(len(assignments)*2, Temp)

	for i := 0; i < size; i++ {
		var propOp runtime.Operand
		var valOp runtime.Operand
		pac := assignments[i].(*fql.PropertyAssignmentContext)

		if prop, ok := pac.PropertyName().(*fql.PropertyNameContext); ok {
			propOp = prop.Accept(v).(runtime.Operand)
			valOp = pac.Expression().Accept(v).(runtime.Operand)
		} else if comProp, ok := pac.ComputedPropertyName().(*fql.ComputedPropertyNameContext); ok {
			propOp = comProp.Accept(v).(runtime.Operand)
			valOp = pac.Expression().Accept(v).(runtime.Operand)
		} else if variable := pac.Variable(); variable != nil {
			propOp = v.symbols.AddConstant(values.NewString(variable.GetText()))
			valOp = variable.Accept(v).(runtime.Operand)
		}

		regIndex := i * 2

		if propOp.IsConstant() {
			v.emitter.EmitAB(runtime.OpLoadConst, seq.Registers[regIndex], propOp)
		} else {
			v.emitter.EmitAB(runtime.OpMove, seq.Registers[regIndex], propOp)
		}

		if valOp.IsConstant() {
			v.emitter.EmitAB(runtime.OpLoadConst, seq.Registers[regIndex+1], valOp)
		} else {
			v.emitter.EmitAB(runtime.OpMove, seq.Registers[regIndex+1], valOp)
		}

		// Free source register if temporary
		if propOp.IsRegister() {
			//v.registers.Free(propOp)
		}
	}

	v.emitter.EmitAs(runtime.OpObject, dst, seq)

	return dst
}

func (v *visitor) VisitPropertyName(ctx *fql.PropertyNameContext) interface{} {
	if id := ctx.Identifier(); id != nil {
		return v.symbols.AddConstant(values.NewString(ctx.GetText()))
	}

	if str := ctx.StringLiteral(); str != nil {
		return str.Accept(v)
	}

	if word := ctx.SafeReservedWord(); word != nil {
		return v.symbols.AddConstant(values.NewString(ctx.GetText()))
	}

	if word := ctx.UnsafeReservedWord(); word != nil {
		return v.symbols.AddConstant(values.NewString(ctx.GetText()))
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitComputedPropertyName(ctx *fql.ComputedPropertyNameContext) interface{} {
	return ctx.Expression().Accept(v)
}

func (v *visitor) VisitStringLiteral(ctx *fql.StringLiteralContext) interface{} {
	var b strings.Builder

	for _, child := range ctx.GetChildren() {
		tree := child.(antlr.TerminalNode)
		sym := tree.GetSymbol()
		input := sym.GetInputStream()

		if input == nil {
			continue
		}

		size := input.Size()
		// skip quotes
		start := sym.GetStart() + 1
		stop := sym.GetStop() - 1

		if stop >= size {
			stop = size - 1
		}

		if start < size && stop < size {
			for i := start; i <= stop; i++ {
				c := input.GetText(i, i)

				switch c {
				case "\\":
					c2 := input.GetText(i, i+1)

					switch c2 {
					case "\\n":
						b.WriteString("\n")
					case "\\t":
						b.WriteString("\t")
					default:
						b.WriteString(c2)
					}

					i++
				default:
					b.WriteString(c)
				}
			}
		}
	}

	return v.symbols.AddConstant(values.NewString(b.String()))
}

func (v *visitor) VisitIntegerLiteral(ctx *fql.IntegerLiteralContext) interface{} {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	return v.symbols.AddConstant(values.NewInt(val))
}

func (v *visitor) VisitFloatLiteral(ctx *fql.FloatLiteralContext) interface{} {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	return v.symbols.AddConstant(values.NewFloat(val))
}

func (v *visitor) VisitBooleanLiteral(ctx *fql.BooleanLiteralContext) interface{} {
	switch strings.ToLower(ctx.GetText()) {
	case "true":
		return v.symbols.AddConstant(values.True)
	case "false":
		return v.symbols.AddConstant(values.False)
	default:
		panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
	}
}

func (v *visitor) VisitNoneLiteral(_ *fql.NoneLiteralContext) interface{} {
	return v.symbols.AddConstant(values.None)
}

func (v *visitor) VisitLiteral(ctx *fql.LiteralContext) interface{} {
	if c := ctx.ArrayLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.ObjectLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.StringLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.IntegerLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.FloatLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.BooleanLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.NoneLiteral(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitReturnExpression(ctx *fql.ReturnExpressionContext) interface{} {
	valReg := ctx.Expression().Accept(v).(runtime.Operand)

	if valReg.IsConstant() {
		v.emitter.EmitAB(runtime.OpLoadGlobal, runtime.ResultOperand, valReg)
	} else {
		v.emitter.EmitAB(runtime.OpMove, runtime.ResultOperand, valReg)
	}

	v.emitter.Emit(runtime.OpReturn)

	//if len(v.loops) == 0 {
	//	v.emitter.EmitABC(runtime.OpReturn)
	//} else {
	//	v.emitter.EmitABC(runtime.OpLoopReturn, v.resolveLoopResultPosition())
	//}

	return runtime.ResultOperand
}

func (v *visitor) VisitExpression(ctx *fql.ExpressionContext) interface{} {
	if uo := ctx.UnaryOperator(); uo != nil {
		src := v.toRegister(ctx.GetRight().Accept(v).(runtime.Operand))
		dst := v.registers.Allocate(Temp)

		uoc := uo.(*fql.UnaryOperatorContext)
		var op runtime.Opcode

		if uoc.Not() != nil {
			op = runtime.OpNot
		} else if uoc.Minus() != nil {
			op = runtime.OpFlipNegative
		} else if uoc.Plus() != nil {
			op = runtime.OpFlipPositive
		} else {
			panic(core.Error(ErrUnexpectedToken, uoc.GetText()))
		}

		// We do not overwrite the source register
		v.emitter.EmitAB(op, dst, src)

		return dst
	}

	if op := ctx.LogicalAndOperator(); op != nil {
		dst := v.registers.Allocate(Temp)
		// Execute left expression
		left := v.toRegister(ctx.GetLeft().Accept(v).(runtime.Operand))
		v.emitter.EmitAB(runtime.OpMove, dst, left)
		// Test if left is false and jump to the end
		end := v.emitter.EmitJump(runtime.OpJumpIfFalse, dst)
		// If left is true, execute right expression
		right := v.toRegister(ctx.GetRight().Accept(v).(runtime.Operand))
		// And move the result to the destination register
		v.emitter.EmitAB(runtime.OpMove, dst, right)
		v.emitter.PatchJump(end)

		return dst
	}

	if op := ctx.LogicalOrOperator(); op != nil {
		dst := v.registers.Allocate(Temp)
		// Execute left expression
		left := v.toRegister(ctx.GetLeft().Accept(v).(runtime.Operand))
		// Move the result to the destination register
		v.emitter.EmitAB(runtime.OpMove, dst, left)
		// Test if left is true and jump to the end
		end := v.emitter.EmitJump(runtime.OpJumpIfTrue, dst)
		// If left is false, execute right expression
		right := v.toRegister(ctx.GetRight().Accept(v).(runtime.Operand))
		// And move the result to the destination register
		v.emitter.EmitAB(runtime.OpMove, dst, right)
		v.emitter.PatchJump(end)

		return dst
	}

	if op := ctx.GetTernaryOperator(); op != nil {
		dst := v.registers.Allocate(Temp)

		// Compile condition and put result in dst
		condReg := v.toRegister(ctx.GetCondition().Accept(v).(runtime.Operand))
		v.emitter.EmitAB(runtime.OpMove, dst, condReg)

		// If condition was temporary, free it
		if condReg.IsRegister() {
			//v.registers.Free(condReg)
		}

		// Jump to 'false' branch if condition is false
		otherwise := v.emitter.EmitJump(runtime.OpJumpIfFalse, dst)

		// True branch
		if onTrue := ctx.GetOnTrue(); onTrue != nil {
			trueReg := v.toRegister(onTrue.Accept(v).(runtime.Operand))
			v.emitter.EmitAB(runtime.OpMove, dst, trueReg)

			// Free temporary register if needed
			if trueReg.IsRegister() {
				//v.registers.Free(trueReg)
			}
		}

		// Jump over false branch
		end := v.emitter.EmitJump(runtime.OpJump, dst)
		v.emitter.PatchJump(otherwise)

		// False branch
		if onFalse := ctx.GetOnFalse(); onFalse != nil {
			falseReg := v.toRegister(onFalse.Accept(v).(runtime.Operand))
			v.emitter.EmitAB(runtime.OpMove, dst, falseReg)

			// Free temporary register if needed
			if falseReg.IsRegister() {
				//v.registers.Free(falseReg)
			}
		}

		v.emitter.PatchJump(end)

		return dst
	}

	if c := ctx.Predicate(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitPredicate(ctx *fql.PredicateContext) interface{} {
	if c := ctx.ExpressionAtom(); c != nil {
		startCatch := v.emitter.Size()
		reg := c.Accept(v)

		if c.ErrorOperator() != nil {
			endCatch := v.emitter.Size()
			v.catchTable = append(v.catchTable, [2]int{startCatch, endCatch})
		}

		return reg
	}

	var opcode runtime.Opcode
	dest := v.registers.Allocate(Temp)
	left := v.toRegister(ctx.Predicate(0).Accept(v).(runtime.Operand))
	right := v.toRegister(ctx.Predicate(1).Accept(v).(runtime.Operand))

	if op := ctx.EqualityOperator(); op != nil {
		switch op.GetText() {
		case "==":
			opcode = runtime.OpEq
		case "!=":
			opcode = runtime.OpNeq
		case ">":
			opcode = runtime.OpGt
		case ">=":
			opcode = runtime.OpGte
		case "<":
			opcode = runtime.OpLt
		case "<=":
			opcode = runtime.OpLte
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.ArrayOperator(); op != nil {
		// TODO: Implement me
		panic(core.Error(core.ErrNotImplemented, "array operator"))
	} else if op := ctx.InOperator(); op != nil {
		opcode = runtime.OpIn
	} else if op := ctx.LikeOperator(); op != nil {
		if op.(*fql.LikeOperatorContext).Not() != nil {
			opcode = runtime.OpNotLike
		} else {
			opcode = runtime.OpLike
		}
	}

	v.emitter.EmitABC(opcode, dest, left, right)

	return dest
}

func (v *visitor) VisitExpressionAtom(ctx *fql.ExpressionAtomContext) interface{} {
	var opcode runtime.Opcode
	var isSet bool

	if op := ctx.MultiplicativeOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "*":
			opcode = runtime.OpMulti
		case "/":
			opcode = runtime.OpDiv
		case "%":
			opcode = runtime.OpMod
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "+":
			opcode = runtime.OpAdd
		case "-":
			opcode = runtime.OpSub
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}

	} else if op := ctx.RegexpOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "=~":
			opcode = runtime.OpRegexpPositive
		case "!~":
			opcode = runtime.OpRegexpNegative
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	}

	if isSet {
		regLeft := v.toRegister(ctx.ExpressionAtom(0).Accept(v).(runtime.Operand))
		regRight := v.toRegister(ctx.ExpressionAtom(1).Accept(v).(runtime.Operand))
		dst := v.registers.Allocate(Temp)

		v.emitter.EmitABC(opcode, dst, regLeft, regRight)

		return dst
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.RangeOperator(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Literal(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	} else if c := ctx.MemberExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Param(); c != nil {
		return c.Accept(v)
	} else if c := ctx.ForExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.WaitForExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Expression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) toRegister(op runtime.Operand) runtime.Operand {
	if op.IsRegister() {
		return op
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(runtime.OpLoadConst, reg, op)

	return reg
}
