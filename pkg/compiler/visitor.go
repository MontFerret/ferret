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

const (
	VarTemporary VarType = iota // Short-lived intermediate results
	VarLocal                    // Local variables
	VarIterator                 // FOR loop iterators
	VarResult                   // Final result variables
)

type (
	VarType int

	Variable struct {
		Name     string
		FirstUse int // Instruction number of first use
		LastUse  int // Instruction number of last use
		Register runtime.Operand
		IsLive   bool
		Type     VarType
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
		// remove un-used return value
		//v.emitPop()
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
	v.beginScope()

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
	v.endScope()
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
	//src := ctx.MemberExpressionSource().(*fql.MemberExpressionSourceContext)
	//
	//if c := src.Variable(); c != nil {
	//	c.Accept(v)
	//} else if c := src.Param(); c != nil {
	//	c.Accept(v)
	//} else if c := src.ObjectLiteral(); c != nil {
	//	c.Accept(v)
	//} else if c := src.ArrayLiteral(); c != nil {
	//	c.Accept(v)
	//} else if c := src.FunctionCall(); c != nil {
	//	c.Accept(v)
	//}
	//
	//segments := ctx.AllMemberExpressionPath()
	//
	//for _, segment := range segments {
	//	p := segment.(*fql.MemberExpressionPathContext)
	//
	//	if c := p.PropertyName(); c != nil {
	//		c.Accept(v)
	//	} else if c := p.ComputedPropertyName(); c != nil {
	//		c.Accept(v)
	//	}
	//
	//	if p.ErrorOperator() != nil {
	//		v.emitter.EmitABC(runtime.OpLoadPropertyOptional)
	//	} else {
	//		v.emitter.EmitABC(runtime.OpLoadProperty)
	//	}
	//}

	return nil
}

func (v *visitor) VisitRangeOperator(ctx *fql.RangeOperatorContext) interface{} {
	//ctx.GetLeft().Accept(v)
	//ctx.GetRight().Accept(v)
	//
	//v.emitter.EmitABC(runtime.OpRange)

	return nil
}

func (v *visitor) VisitRangeOperand(ctx *fql.RangeOperandContext) interface{} {
	if c := ctx.IntegerLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.Variable(); c != nil {
		c.Accept(v)
	} else if c := ctx.Param(); c != nil {
		c.Accept(v)
	} else {
		panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return nil
}

func (v *visitor) VisitVariableDeclaration(ctx *fql.VariableDeclarationContext) interface{} {
	name := ignorePseudoVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	valReg := ctx.Expression().Accept(v).(runtime.Operand)

	if name != ignorePseudoVariable {
		varReg := v.symbols.DefineVariable(name)

		if valReg.IsConstant() {
			v.emitter.EmitAB(runtime.OpLoadConst, varReg, valReg)
		} else if v.symbols.Scope() == 0 {
			v.emitter.EmitAB(runtime.OpStoreGlobal, varReg, valReg)
		} else {
			v.emitter.EmitAB(runtime.OpMove, varReg, valReg)
		}

		return varReg
	}

	return nil
}

func (v *visitor) VisitVariable(ctx *fql.VariableContext) interface{} {
	// Just return the register / constant index
	op := v.symbols.LookupVariable(ctx.GetText())

	if op.IsRegister() {
		return op
	}

	reg := v.registers.Allocate(VarTemporary)
	v.emitter.EmitAB(runtime.OpLoadGlobal, reg, op)

	return reg
}

func (v *visitor) VisitArrayLiteral(ctx *fql.ArrayLiteralContext) interface{} {
	// Allocate destination register for the array
	destReg := v.registers.Allocate(VarTemporary)

	if list := ctx.ArgumentList(); list != nil {
		// Get all array element expressions
		exps := list.(*fql.ArgumentListContext).AllExpression()
		size := len(exps)

		if size > 0 {
			// Allocate sequence for array elements
			sequence := v.registers.AllocateSequence(size, VarTemporary)

			// Evaluate each element into sequence registers
			for i, exp := range exps {
				// Compile expression and move to sequence register
				srcReg := exp.Accept(v).(runtime.Operand)

				if srcReg.IsConstant() {
					v.emitter.EmitAB(runtime.OpLoadConst, sequence.Registers[i], srcReg)
				} else {
					v.emitter.EmitAB(runtime.OpMove, sequence.Registers[i], srcReg)
				}

				// Free source register if temporary
				if srcReg.IsRegister() {
					v.registers.Free(srcReg)
				}
			}

			// Initialize an array
			v.emitter.EmitAs(runtime.OpArray, destReg, sequence)

			// Free sequence registers
			v.registers.FreeSequence(sequence)

			return destReg
		}
	}

	// Empty array
	v.emitter.EmitA(runtime.OpArray, destReg)

	return destReg
}

func (v *visitor) VisitArgumentList(ctx *fql.ArgumentListContext) interface{} {
	//exps := ctx.AllExpression()
	//size := len(exps)
	//
	//for _, arg := range exps {
	//	arg.Accept(v)
	//}

	return nil
}

func (v *visitor) VisitObjectLiteral(ctx *fql.ObjectLiteralContext) interface{} {
	//assignments := ctx.AllPropertyAssignment()
	//
	//for _, pa := range assignments {
	//	pac := pa.(*fql.PropertyAssignmentContext)
	//
	//	if prop, ok := pac.PropertyName().(*fql.PropertyNameContext); ok {
	//		prop.Accept(v)
	//		pac.Expression().Accept(v)
	//	} else if comProp, ok := pac.ComputedPropertyName().(*fql.ComputedPropertyNameContext); ok {
	//		comProp.Accept(v)
	//		pac.Expression().Accept(v)
	//	} else if variable := pac.Variable(); variable != nil {
	//		v.emitConstant(values.NewString(variable.GetText()))
	//		variable.Accept(v)
	//	}
	//}
	//
	//v.emitter.EmitABC(runtime.OpObject, len(assignments))

	return nil
}

func (v *visitor) VisitPropertyName(ctx *fql.PropertyNameContext) interface{} {
	//if id := ctx.Identifier(); id != nil {
	//	v.emitConstant(values.NewString(ctx.GetText()))
	//} else if str := ctx.StringLiteral(); str != nil {
	//	str.Accept(v)
	//} else if word := ctx.SafeReservedWord(); word != nil {
	//	v.emitConstant(values.NewString(ctx.GetText()))
	//} else if word := ctx.UnsafeReservedWord(); word != nil {
	//	v.emitConstant(values.NewString(ctx.GetText()))
	//}

	return nil
}

func (v *visitor) VisitComputedPropertyName(ctx *fql.ComputedPropertyNameContext) interface{} {
	ctx.Expression().Accept(v)

	return nil
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
	if op := ctx.UnaryOperator(); op != nil {
		ctx.GetRight().Accept(v)

		op := op.(*fql.UnaryOperatorContext)

		if op.Not() != nil {
			//v.emitter.EmitABC(runtime.OpNot)
		} else if op.Minus() != nil {
			// v.emitter.EmitABC(runtime.OpFlipNegative)
		} else if op.Plus() != nil {
			//v.emitter.EmitABC(runtime.OpFlipPositive)
		} else {
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.LogicalAndOperator(); op != nil {
		ctx.GetLeft().Accept(v)
		end := v.emitJump(runtime.OpJumpIfFalse)
		ctx.GetRight().Accept(v)
		v.patchJump(end)
	} else if op := ctx.LogicalOrOperator(); op != nil {
		ctx.GetLeft().Accept(v)
		end := v.emitJump(runtime.OpJumpIfTrue)
		ctx.GetRight().Accept(v)
		v.patchJump(end)
	} else if op := ctx.GetTernaryOperator(); op != nil {
		ctx.GetCondition().Accept(v)

		otherwise := v.emitJump(runtime.OpJumpIfFalse)

		if onTrue := ctx.GetOnTrue(); onTrue != nil {
			onTrue.Accept(v)
		}

		end := v.emitJump(runtime.OpJump)
		v.patchJump(otherwise)

		ctx.GetOnFalse().Accept(v)
		v.patchJump(end)
	} else if c := ctx.Predicate(); c != nil {
		return c.Accept(v)
	}

	return nil
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
	dest := v.registers.Allocate(VarTemporary)

	if op := ctx.EqualityOperator(); op != nil {
		src1 := ctx.Predicate(0).Accept(v).(runtime.Operand)
		src2 := ctx.Predicate(1).Accept(v).(runtime.Operand)

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

		v.emitter.EmitABC(opcode, dest, src1, src2)
	} else if op := ctx.ArrayOperator(); op != nil {
		// TODO: Implement me
		panic(core.Error(core.ErrNotImplemented, "array operator"))
	} else if op := ctx.InOperator(); op != nil {
		src1 := ctx.Predicate(0).Accept(v).(runtime.Operand)
		src2 := ctx.Predicate(1).Accept(v).(runtime.Operand)
		opcode = runtime.OpIn

		v.emitter.EmitABC(opcode, dest, src1, src2)
	} else if op := ctx.LikeOperator(); op != nil {
		src1 := ctx.Predicate(0).Accept(v).(runtime.Operand)
		src2 := ctx.Predicate(1).Accept(v).(runtime.Operand)

		if op.(*fql.LikeOperatorContext).Not() != nil {
			opcode = runtime.OpNotLike
		} else {
			opcode = runtime.OpLike
		}

		v.emitter.EmitABC(opcode, dest, src1, src2)
	}

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
		regLeft := ctx.ExpressionAtom(0).Accept(v).(runtime.Operand)
		regRight := ctx.ExpressionAtom(1).Accept(v).(runtime.Operand)
		dst := v.registers.Allocate(VarTemporary)

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

	return nil
}

func (v *visitor) beginScope() {
	v.symbols.EnterScope()
}

func (v *visitor) endScope() {
	v.symbols.ExitScope()
}

func (v *visitor) beginLoopScope(passThrough, distinct bool) {
	//var allocate bool
	//var prevResult int
	//
	//// top loop
	//if len(v.loops) == 0 {
	//	allocate = true
	//} else if !passThrough {
	//	// nested with explicit RETURN expression
	//	prev := v.loops[len(v.loops)-1]
	//	// if the loop above does not do pass through
	//	// we allocate a new array for this loop
	//	allocate = !prev.passThrough
	//	prevResult = prev.result
	//}
	//
	//var resultPos int
	//
	//if allocate {
	//	var arg int
	//
	//	if distinct {
	//		arg = 1
	//	}
	//
	//	resultPos = v.operandsStackTracker
	//	v.emitter.EmitABC(runtime.OpLoopInitOutput, arg)
	//} else {
	//	resultPos = prevResult
	//}
	//
	//v.loops = append(v.loops, &loopScope{
	//	result:      resultPos,
	//	passThrough: passThrough,
	//})
}

func (v *visitor) patchLoopScope(jump int) {
	//v.loops[len(v.loops)-1].position = jump
}

func (v *visitor) resolveLoopResultPosition() int {
	//return v.loops[len(v.loops)-1].result
	return 0
}

func (v *visitor) resolveLoopPosition() int {
	//return v.loops[len(v.loops)-1].position
	return 0
}

func (v *visitor) endLoopScope() {
	//v.loops = v.loops[:len(v.loops)-1]
	//
	//var unwrap bool
	//
	//if len(v.loops) == 0 {
	//	unwrap = true
	//} else if !v.loops[len(v.loops)-1].passThrough {
	//	unwrap = true
	//}
	//
	//if unwrap {
	//	v.emitter.EmitABC(runtime.OpLoopUnwrapOutput)
	//}
}

// emitLoop emits a loop instruction.
func (v *visitor) emitLoop(loopStart int) {
	//pos := v.emitJump(runtime.OpJumpBackward)
	//jump := pos - loopStart
	//v.arguments[pos-1] = jump
}

// emitJump emits an opcode with a jump result argument.
func (v *visitor) emitJump(op runtime.Opcode) int {
	//v.emitter.EmitABC(op, jumpPlaceholder)
	//
	//return len(v.instructions)

	return 0
}

// patchJump patches a jump result argument.
func (v *visitor) patchJump(offset int) {
	//jump := len(v.instructions) - offset
	//v.arguments[offset-1] = jump
}

func (v *visitor) patchJumpWith(offset, jump int) {
	//v.arguments[offset-1] = jump
}

func (v *visitor) emitPopAndClose() {
	//v.emitter.EmitABC(runtime.OpPopClose)
}
