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
	variable struct {
		name  string
		depth int
	}

	loopScope struct {
		result      int
		passThrough bool
		position    int
	}

	visitor struct {
		*fql.BaseFqlParserVisitor
		err error
		src string
		//funcs          core.Functions
		constantsIndex        map[uint64]int
		locations             []core.Location
		bytecode              []runtime.Opcode
		arguments             []int
		constants             []core.Value
		scope                 int
		loops                 []*loopScope
		globals               map[string]int
		locals                []variable
		catchTable            [][2]int
		operandsStackTracker  int
		variablesStackTracker int
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
	//v.funcs = funcs
	v.constantsIndex = make(map[uint64]int)
	v.locations = make([]core.Location, 0)
	v.bytecode = make([]runtime.Opcode, 0)
	v.arguments = make([]int, 0)
	v.constants = make([]core.Value, 0)
	v.scope = 0
	v.loops = make([]*loopScope, 0)
	v.globals = make(map[string]int)
	v.locals = make([]variable, 0)
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
		c.Accept(v)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		c.Accept(v)
		// remove un-used return value
		v.emitPop()
	} else if c := ctx.WaitForExpression(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitBodyExpression(ctx *fql.BodyExpressionContext) interface{} {
	if c := ctx.ForExpression(); c != nil {
		c.Accept(v)
	} else if c := ctx.ReturnExpression(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}

func (v *visitor) VisitForExpression(ctx *fql.ForExpressionContext) interface{} {
	v.beginScope()

	var passThrough bool
	var distinct bool
	var returnRuleCtx antlr.RuleContext
	var loopJump, exitJump int
	// identify whether it's WHILE or FOR loop
	isForInLoop := ctx.While() == nil
	returnCtx := ctx.ForExpressionReturn()

	if c := returnCtx.ReturnExpression(); c != nil {
		returnRuleCtx = c
		distinct = c.Distinct() != nil
	} else if c := returnCtx.ForExpression(); c != nil {
		returnRuleCtx = c
		passThrough = true
	}

	v.beginLoopScope(passThrough, distinct)

	if isForInLoop {
		// Loop data source to iterate over
		if c := ctx.ForExpressionSource(); c != nil {
			c.Accept(v)
		}

		v.emit(runtime.OpForLoopInitInput)
		loopJump = len(v.bytecode)
		v.emit(runtime.OpForLoopHasNext)
		exitJump = v.emitJump(runtime.OpJumpIfFalse)
		// pop the boolean value from the stack
		v.emitPop()

		valVar := ctx.GetValueVariable().GetText()
		counterVarCtx := ctx.GetCounterVariable()

		hasValVar := valVar != ignorePseudoVariable
		var hasCounterVar bool
		var counterVar string

		if counterVarCtx != nil {
			counterVar = counterVarCtx.GetText()
			hasCounterVar = true
		}

		var valVarIndex int

		// declare value variable
		if hasValVar {
			valVarIndex = v.declareVariable(valVar)
		}

		var counterVarIndex int

		if hasCounterVar {
			// declare counter variable
			counterVarIndex = v.declareVariable(counterVar)
		}

		if hasValVar && hasCounterVar {
			// we will calculate the index of the counter variable
			v.emit(runtime.OpForLoopNext)
		} else if hasValVar {
			v.emit(runtime.OpForLoopNextValue)
		} else if hasCounterVar {
			v.emit(runtime.OpForLoopNextCounter)
		} else {
			panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
		}

		if hasValVar {
			v.defineVariable(valVarIndex)
		}

		if hasCounterVar {
			v.defineVariable(counterVarIndex)
		}
	} else {
		// Create initial value for the loop counter
		v.emit(runtime.OpWhileLoopInitCounter)

		loopJump = len(v.bytecode)

		// Condition expression
		ctx.Expression().Accept(v)

		// Condition check
		exitJump = v.emitJump(runtime.OpJumpIfFalse)
		// pop the boolean value from the stack
		v.emitPop()

		counterVar := ctx.GetCounterVariable().GetText()

		// declare counter variable
		// and increment it by 1
		index := v.declareVariable(counterVar)
		v.emit(runtime.OpWhileLoopNext)
		v.defineVariable(index)
	}

	v.patchLoopScope(loopJump)

	// body
	if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
		for _, b := range body {
			b.Accept(v)
		}
	}

	// return
	returnRuleCtx.Accept(v)

	v.emitLoop(loopJump)
	v.patchJump(exitJump)
	v.endScope()
	// pop the boolean value from the stack
	v.emitPop()

	if isForInLoop {
		// pop the iterator
		v.emitPopAndClose()
	} else {
		// pop the counter
		v.emitPop()
	}

	v.endLoopScope()

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
	ctx.Expression().Accept(v)
	fwd := v.emitJump(runtime.OpJumpIfTrue)
	// Pop on false
	v.emitPop()
	// And jump back to the beginning of the loop
	bwd := v.emitJump(runtime.OpJumpBackward)
	v.patchJumpWith(bwd, len(v.bytecode)-v.resolveLoopPosition())

	// Otherwise, pop on true
	v.patchJump(fwd)
	v.emitPop()

	return nil
}

func (v *visitor) VisitForExpressionStatement(ctx *fql.ForExpressionStatementContext) interface{} {
	if c := ctx.VariableDeclaration(); c != nil {
		c.Accept(v)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		c.Accept(v)
		// remove un-used return value
		v.emitPop()
	}

	return nil
}

func (v *visitor) VisitFunctionCallExpression(ctx *fql.FunctionCallExpressionContext) interface{} {
	call := ctx.FunctionCall().(*fql.FunctionCallContext)

	var name string

	funcNS := call.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += call.FunctionName().GetText()

	isNonOptional := ctx.ErrorOperator() == nil

	v.emitConstant(values.String(name))

	var size int

	if args := call.ArgumentList(); args != nil {
		out := v.VisitArgumentList(args.(*fql.ArgumentListContext))
		size = out.(int)
	}

	switch size {
	case 0:
		if isNonOptional {
			v.emit(runtime.OpCall, 0)
		} else {
			v.emit(runtime.OpCallSafe, 0)
		}
	case 1:
		if isNonOptional {
			v.emit(runtime.OpCall1, 1)
		} else {
			v.emit(runtime.OpCall1Safe, 1)
		}
	case 2:
		if isNonOptional {
			v.emit(runtime.OpCall2, 2)
		} else {
			v.emit(runtime.OpCall2Safe, 2)
		}
	case 3:
		if isNonOptional {
			v.emit(runtime.OpCall3, 3)
		} else {
			v.emit(runtime.OpCall3Safe, 3)
		}
	case 4:
		if isNonOptional {
			v.emit(runtime.OpCall4, 4)
		} else {
			v.emit(runtime.OpCall4Safe, 4)
		}
	default:
		if isNonOptional {
			v.emit(runtime.OpCallN, size)
		} else {
			v.emit(runtime.OpCallNSafe, size)
		}
	}

	return nil
}

func (v *visitor) VisitMemberExpression(ctx *fql.MemberExpressionContext) interface{} {
	src := ctx.MemberExpressionSource().(*fql.MemberExpressionSourceContext)

	if c := src.Variable(); c != nil {
		c.Accept(v)
	} else if c := src.Param(); c != nil {
		c.Accept(v)
	} else if c := src.ObjectLiteral(); c != nil {
		c.Accept(v)
	} else if c := src.ArrayLiteral(); c != nil {
		c.Accept(v)
	} else if c := src.FunctionCall(); c != nil {
		c.Accept(v)
	}

	segments := ctx.AllMemberExpressionPath()

	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)

		if c := p.PropertyName(); c != nil {
			c.Accept(v)
		} else if c := p.ComputedPropertyName(); c != nil {
			c.Accept(v)
		}

		if p.ErrorOperator() != nil {
			v.emit(runtime.OpLoadPropertyOptional)
		} else {
			v.emit(runtime.OpLoadProperty)
		}
	}

	return nil
}

func (v *visitor) VisitRangeOperator(ctx *fql.RangeOperatorContext) interface{} {
	ctx.GetLeft().Accept(v)
	ctx.GetRight().Accept(v)

	v.emit(runtime.OpRange)

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

	ctx.Expression().Accept(v)

	if name != ignorePseudoVariable {
		// we do not have custom functions, thus this feature is not needed at this moment
		index := v.declareVariable(name)
		v.defineVariable(index)
	} else {
		// if we ignore the result of the execution, we pop it from the stack
		v.emitPop()
	}

	return nil
}

func (v *visitor) VisitVariable(ctx *fql.VariableContext) interface{} {
	v.readVariable(ctx.GetText())

	return nil
}

func (v *visitor) VisitArrayLiteral(ctx *fql.ArrayLiteralContext) interface{} {
	var size int

	if args := ctx.ArgumentList(); args != nil {
		out := v.VisitArgumentList(args.(*fql.ArgumentListContext))
		size = out.(int)
	}

	v.emit(runtime.OpArray, size)

	return nil
}

func (v *visitor) VisitArgumentList(ctx *fql.ArgumentListContext) interface{} {
	exps := ctx.AllExpression()
	size := len(exps)

	for _, arg := range exps {
		arg.Accept(v)
	}

	return size
}

func (v *visitor) VisitObjectLiteral(ctx *fql.ObjectLiteralContext) interface{} {
	assignments := ctx.AllPropertyAssignment()

	for _, pa := range assignments {
		pac := pa.(*fql.PropertyAssignmentContext)

		if prop, ok := pac.PropertyName().(*fql.PropertyNameContext); ok {
			prop.Accept(v)
			pac.Expression().Accept(v)
		} else if comProp, ok := pac.ComputedPropertyName().(*fql.ComputedPropertyNameContext); ok {
			comProp.Accept(v)
			pac.Expression().Accept(v)
		} else if variable := pac.Variable(); variable != nil {
			v.emitConstant(values.NewString(variable.GetText()))
			variable.Accept(v)
		}
	}

	v.emit(runtime.OpObject, len(assignments))

	return nil
}

func (v *visitor) VisitPropertyName(ctx *fql.PropertyNameContext) interface{} {
	if id := ctx.Identifier(); id != nil {
		v.emitConstant(values.NewString(ctx.GetText()))
	} else if str := ctx.StringLiteral(); str != nil {
		str.Accept(v)
	} else if word := ctx.SafeReservedWord(); word != nil {
		v.emitConstant(values.NewString(ctx.GetText()))
	} else if word := ctx.UnsafeReservedWord(); word != nil {
		v.emitConstant(values.NewString(ctx.GetText()))
	}

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

	v.emitConstant(values.NewString(b.String()))

	return nil
}

func (v *visitor) VisitIntegerLiteral(ctx *fql.IntegerLiteralContext) interface{} {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	v.emitConstant(values.NewInt(val))

	return nil
}

func (v *visitor) VisitFloatLiteral(ctx *fql.FloatLiteralContext) interface{} {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	v.emitConstant(values.NewFloat(val))

	return nil
}

func (v *visitor) VisitBooleanLiteral(ctx *fql.BooleanLiteralContext) interface{} {
	switch strings.ToLower(ctx.GetText()) {
	case "true":
		v.emit(runtime.OpTrue)
	case "false":
		v.emit(runtime.OpFalse)
	default:
		panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return nil
}

func (v *visitor) VisitNoneLiteral(ctx *fql.NoneLiteralContext) interface{} {
	v.emit(runtime.OpNone)

	return nil
}

func (v *visitor) VisitLiteral(ctx *fql.LiteralContext) interface{} {
	if c := ctx.ArrayLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.ObjectLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.StringLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.IntegerLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.FloatLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.BooleanLiteral(); c != nil {
		c.Accept(v)
	} else if c := ctx.NoneLiteral(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitReturnExpression(ctx *fql.ReturnExpressionContext) interface{} {
	ctx.Expression().Accept(v)

	if len(v.loops) == 0 {
		v.emit(runtime.OpReturn)
	} else {
		v.emit(runtime.OpLoopReturn, v.resolveLoopResultPosition())
	}

	return nil
}

func (v *visitor) VisitExpression(ctx *fql.ExpressionContext) interface{} {
	if op := ctx.UnaryOperator(); op != nil {
		ctx.GetRight().Accept(v)

		op := op.(*fql.UnaryOperatorContext)

		if op.Not() != nil {
			v.emit(runtime.OpNot)
		} else if op.Minus() != nil {
			v.emit(runtime.OpFlipNegative)
		} else if op.Plus() != nil {
			v.emit(runtime.OpFlipPositive)
		} else {
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.LogicalAndOperator(); op != nil {
		ctx.GetLeft().Accept(v)
		end := v.emitJump(runtime.OpJumpIfFalse)
		v.emitPop()
		ctx.GetRight().Accept(v)
		v.patchJump(end)
	} else if op := ctx.LogicalOrOperator(); op != nil {
		ctx.GetLeft().Accept(v)
		end := v.emitJump(runtime.OpJumpIfTrue)
		v.emitPop()
		ctx.GetRight().Accept(v)
		v.patchJump(end)
	} else if op := ctx.GetTernaryOperator(); op != nil {
		ctx.GetCondition().Accept(v)

		otherwise := v.emitJump(runtime.OpJumpIfFalse)

		if onTrue := ctx.GetOnTrue(); onTrue != nil {
			// Remove the top value from the stack (the condition)
			v.emitPop()
			onTrue.Accept(v)
		}

		end := v.emitJump(runtime.OpJump)
		v.patchJump(otherwise)

		// Remove the top value from the stack (the condition)
		v.emitPop()
		ctx.GetOnFalse().Accept(v)
		v.patchJump(end)
	} else if c := ctx.Predicate(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitPredicate(ctx *fql.PredicateContext) interface{} {
	if op := ctx.EqualityOperator(); op != nil {
		ctx.Predicate(0).Accept(v)
		ctx.Predicate(1).Accept(v)

		switch op.GetText() {
		case "==":
			v.emit(runtime.OpEq)
		case "!=":
			v.emit(runtime.OpNeq)
		case ">":
			v.emit(runtime.OpGt)
		case ">=":
			v.emit(runtime.OpGte)
		case "<":
			v.emit(runtime.OpLt)
		case "<=":
			v.emit(runtime.OpLte)
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.ArrayOperator(); op != nil {
		// TODO: Implement me
	} else if op := ctx.InOperator(); op != nil {
		ctx.Predicate(0).Accept(v)
		ctx.Predicate(1).Accept(v)

		v.emit(runtime.OpIn)
	} else if op := ctx.LikeOperator(); op != nil {
		ctx.Predicate(0).Accept(v)
		ctx.Predicate(1).Accept(v)

		if op.(*fql.LikeOperatorContext).Not() != nil {
			v.emit(runtime.OpNotLike)
		} else {
			v.emit(runtime.OpLike)
		}
	} else if c := ctx.ExpressionAtom(); c != nil {
		startCatch := len(v.bytecode)
		c.Accept(v)

		if c.ErrorOperator() != nil {
			endCatch := len(v.bytecode)
			v.catchTable = append(v.catchTable, [2]int{startCatch, endCatch})
		}
	}

	return nil
}

func (v *visitor) VisitExpressionAtom(ctx *fql.ExpressionAtomContext) interface{} {
	if op := ctx.MultiplicativeOperator(); op != nil {
		ctx.ExpressionAtom(0).Accept(v)
		ctx.ExpressionAtom(1).Accept(v)

		switch op.GetText() {
		case "*":
			v.emit(runtime.OpMulti)
		case "/":
			v.emit(runtime.OpDiv)
		case "%":
			v.emit(runtime.OpMod)
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		ctx.ExpressionAtom(0).Accept(v)
		ctx.ExpressionAtom(1).Accept(v)

		switch op.GetText() {
		case "+":
			v.emit(runtime.OpAdd)
		case "-":
			v.emit(runtime.OpSub)
		}
	} else if op := ctx.RegexpOperator(); op != nil {
		ctx.ExpressionAtom(0).Accept(v)
		ctx.ExpressionAtom(1).Accept(v)

		switch op.GetText() {
		case "=~":
			v.emit(runtime.OpRegexpPositive)
		case "!~":
			v.emit(runtime.OpRegexpNegative)
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if c := ctx.FunctionCallExpression(); c != nil {
		c.Accept(v)
	} else if c := ctx.RangeOperator(); c != nil {
		c.Accept(v)
	} else if c := ctx.Literal(); c != nil {
		c.Accept(v)
	} else if c := ctx.Variable(); c != nil {
		c.Accept(v)
	} else if c := ctx.MemberExpression(); c != nil {
		c.Accept(v)
	} else if c := ctx.Param(); c != nil {
		c.Accept(v)
	} else if c := ctx.ForExpression(); c != nil {
		c.Accept(v)
	} else if c := ctx.WaitForExpression(); c != nil {
		c.Accept(v)
	} else if c := ctx.Expression(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) beginScope() {
	v.scope++
}

func (v *visitor) endScope() {
	v.scope--

	// Pop all local variables from the stack within the closed scope.
	for len(v.locals) > 0 && v.locals[len(v.locals)-1].depth > v.scope {
		v.emit(runtime.OpPopLocal)
		v.locals = v.locals[:len(v.locals)-1]
	}
}

func (v *visitor) beginLoopScope(passThrough, distinct bool) {
	var allocate bool
	var prevResult int

	// top loop
	if len(v.loops) == 0 {
		allocate = true
	} else if !passThrough {
		// nested with explicit RETURN expression
		prev := v.loops[len(v.loops)-1]
		// if the loop above does not do pass through
		// we allocate a new array for this loop
		allocate = !prev.passThrough
		prevResult = prev.result
	}

	var resultPos int

	if allocate {
		var arg int

		if distinct {
			arg = 1
		}

		resultPos = v.operandsStackTracker
		v.emit(runtime.OpLoopInitOutput, arg)
	} else {
		resultPos = prevResult
	}

	v.loops = append(v.loops, &loopScope{
		result:      resultPos,
		passThrough: passThrough,
	})
}

func (v *visitor) patchLoopScope(jump int) {
	v.loops[len(v.loops)-1].position = jump
}

func (v *visitor) resolveLoopResultPosition() int {
	return v.loops[len(v.loops)-1].result
}

func (v *visitor) resolveLoopPosition() int {
	return v.loops[len(v.loops)-1].position
}

func (v *visitor) endLoopScope() {
	v.loops = v.loops[:len(v.loops)-1]

	var unwrap bool

	if len(v.loops) == 0 {
		unwrap = true
	} else if !v.loops[len(v.loops)-1].passThrough {
		unwrap = true
	}

	if unwrap {
		v.emit(runtime.OpLoopUnwrapOutput)
	}
}

func (v *visitor) resolveLocalVariable(name string) int {
	for i := len(v.locals) - 1; i >= 0; i-- {
		if v.locals[i].name == name {
			return i
		}
	}

	return -1
}

func (v *visitor) readVariable(name string) {
	if name == pseudoVariable {
		return
	}

	// Resolve the variable name to an index.
	arg := v.resolveLocalVariable(name)

	if arg > -1 {
		v.emit(runtime.OpLoadLocal, arg)

		return
	}

	index, ok := v.globals[name]

	if !ok {
		panic(core.Error(ErrVariableNotFound, name))
	}

	v.emit(runtime.OpLoadGlobal, index)
}

func (v *visitor) declareVariable(name string) int {
	if name == ignorePseudoVariable {
		return -1
	}

	if v.scope == 0 {
		// Check for duplicate global variable names.
		_, ok := v.globals[name]

		if ok {
			panic(core.Error(ErrVariableNotUnique, name))
		}

		index := v.addConstant(values.String(name))
		v.globals[name] = index

		return index
	}

	// Check for duplicate variable names in the current scope.
	for i := len(v.locals) - 1; i >= 0; i-- {
		local := v.locals[i]

		if local.depth > -1 && local.depth < v.scope {
			break
		}

		if local.name == name {
			panic(core.Error(ErrVariableNotUnique, name))
		}
	}

	v.locals = append(v.locals, variable{name, undefinedVariable})

	return len(v.locals) - 1
}

// defineVariable defines a variable in the current scope.
func (v *visitor) defineVariable(index int) {
	if v.scope == 0 {
		v.emit(runtime.OpStoreGlobal, index)

		return
	}

	v.emit(runtime.OpStoreLocal, index)
	v.locals[index].depth = v.scope
}

// emitConstant emits an opcode with a constant argument.
func (v *visitor) emitConstant(constant core.Value) {
	v.emit(runtime.OpPush, v.addConstant(constant))
}

// emitLoop emits a loop instruction.
func (v *visitor) emitLoop(loopStart int) {
	pos := v.emitJump(runtime.OpJumpBackward)
	jump := pos - loopStart
	v.arguments[pos-1] = jump
}

// emitJump emits an opcode with a jump result argument.
func (v *visitor) emitJump(op runtime.Opcode) int {
	v.emit(op, jumpPlaceholder)

	return len(v.bytecode)
}

// patchJump patches a jump result argument.
func (v *visitor) patchJump(offset int) {
	jump := len(v.bytecode) - offset
	v.arguments[offset-1] = jump
}

func (v *visitor) patchJumpWith(offset, jump int) {
	v.arguments[offset-1] = jump
}

func (v *visitor) emitPop() {
	v.emit(runtime.OpPop)
}

func (v *visitor) emitPopAndClose() {
	v.emit(runtime.OpPopClose)
}

func (v *visitor) emit(op runtime.Opcode, args ...int) {
	v.bytecode = append(v.bytecode, op)

	var arg int

	if len(args) > 0 {
		arg = args[0]
	}

	v.arguments = append(v.arguments, arg)
	v.updateStackTracker(op, arg)
}

func (v *visitor) updateStackTracker(op runtime.Opcode, arg int) {
	switch op {
	case runtime.OpPush:
		v.operandsStackTracker++

	case runtime.OpPop:
		v.operandsStackTracker--

	case runtime.OpPopClose:
		v.operandsStackTracker--

	case runtime.OpStoreGlobal:
		v.operandsStackTracker--

	case runtime.OpLoadGlobal:
		v.operandsStackTracker++

	case runtime.OpStoreLocal:
		v.operandsStackTracker--
		v.variablesStackTracker++

	case runtime.OpPopLocal:
		v.variablesStackTracker--

	case runtime.OpLoadLocal:
		v.operandsStackTracker++

	case runtime.OpNone:
		v.operandsStackTracker++

	case runtime.OpCastBool:
		break

	case runtime.OpTrue:
		v.operandsStackTracker++

	case runtime.OpFalse:
		v.operandsStackTracker++

	case runtime.OpArray:
		v.operandsStackTracker++
		v.operandsStackTracker -= arg

	case runtime.OpObject:
		v.operandsStackTracker++
		v.operandsStackTracker -= arg * 2

	case runtime.OpLoadProperty, runtime.OpLoadPropertyOptional:
		v.operandsStackTracker--

	case runtime.OpNegate, runtime.OpFlipPositive, runtime.OpFlipNegative, runtime.OpNot:
		break

	case runtime.OpEq, runtime.OpNeq:
		v.operandsStackTracker--

	case runtime.OpGt, runtime.OpLt, runtime.OpGte, runtime.OpLte:
		v.operandsStackTracker--

	case runtime.OpIn, runtime.OpNotIn:
		v.operandsStackTracker--

	case runtime.OpLike, runtime.OpNotLike:
		v.operandsStackTracker--

	case runtime.OpAdd, runtime.OpSub, runtime.OpMulti, runtime.OpDiv, runtime.OpMod:
		v.operandsStackTracker--

	case runtime.OpIncr, runtime.OpDecr:
		break

	case runtime.OpRegexpPositive, runtime.OpRegexpNegative:
		break

	case runtime.OpCall, runtime.OpCallSafe:
		break

	case runtime.OpCall1, runtime.OpCall1Safe:
		v.operandsStackTracker--

	case runtime.OpCall2, runtime.OpCall2Safe:
		v.operandsStackTracker -= 2

	case runtime.OpCall3, runtime.OpCall3Safe:
		v.operandsStackTracker -= 3

	case runtime.OpCall4, runtime.OpCall4Safe:
		v.operandsStackTracker -= 4

	case runtime.OpCallN, runtime.OpCallNSafe:
		v.operandsStackTracker -= arg

	case runtime.OpRange:
		v.operandsStackTracker--

	case runtime.OpLoopInitOutput:
		v.operandsStackTracker++

	case runtime.OpLoopUnwrapOutput, runtime.OpForLoopInitInput:
		break

	case runtime.OpForLoopHasNext:
		v.operandsStackTracker++

	case runtime.OpForLoopNext:
		v.operandsStackTracker += 2

	case runtime.OpForLoopNextValue, runtime.OpForLoopNextCounter:
		v.operandsStackTracker++

	case runtime.OpWhileLoopInitCounter:
		v.operandsStackTracker++

	case runtime.OpWhileLoopNext:
		v.operandsStackTracker += 2

	case runtime.OpJump, runtime.OpJumpBackward, runtime.OpJumpIfFalse, runtime.OpJumpIfTrue:
		break

	case runtime.OpLoopReturn, runtime.OpReturn:
		break
	}
}

// addConstant adds a constant to the constants pool and returns its index.
// If the constant is a scalar, it will be deduplicated.
// If the constant is not a scalar, it will be added to the pool without deduplication.
func (v *visitor) addConstant(constant core.Value) int {
	var hash uint64

	if values.IsScalar(constant) {
		hash = constant.Hash()
	}

	if hash > 0 {
		if p, ok := v.constantsIndex[hash]; ok {
			return p
		}
	}

	v.constants = append(v.constants, constant)
	p := len(v.constants) - 1

	if hash > 0 {
		v.constantsIndex[hash] = p
	}

	return p
}
