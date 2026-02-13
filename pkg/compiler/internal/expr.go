package internal

import (
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Runtime functions
const (
	runtimeTypename = "TYPENAME"
	runtimeLength   = "LENGTH"
	runtimeWait     = "WAIT"
)

// ExprCompiler handles the compilation of expressions in FQL queries.
// It transforms expression operations from the AST into VM instructions.
type ExprCompiler struct {
	ctx *CompilerContext
}

// NewExprCompiler creates a new instance of ExprCompiler with the given compiler context.
func NewExprCompiler(ctx *CompilerContext) *ExprCompiler {
	return &ExprCompiler{ctx: ctx}
}

// Compile processes an expression from the FQL AST and delegates to the appropriate
// compilation method based on the expression type (unary, logical, ternary, or predicate).
// Parameters:
//   - ctx: The expression context from the AST
//
// Returns:
//   - An operand representing the compiled expression
//
// Panics if the expression type is not recognized.
func (c *ExprCompiler) Compile(ctx fql.IExpressionContext) vm.Operand {
	if uo := ctx.UnaryOperator(); uo != nil {
		return c.compileUnary(uo, ctx)
	}

	if and := ctx.LogicalAndOperator(); and != nil {
		return c.compileLogicalAnd(ctx)
	}

	if or := ctx.LogicalOrOperator(); or != nil {
		return c.compileLogicalOr(ctx)
	}

	if t := ctx.GetTernaryOperator(); t != nil {
		return c.compileTernary(ctx)
	}

	if p := ctx.Predicate(); p != nil {
		return c.compilePredicate(p)
	}

	return vm.NoopOperand
}

func (c *ExprCompiler) CompileIncDec(token antlr.Token, target vm.Operand) vm.Operand {
	if target.IsConstant() {
		panic("cannot increment/decrement a constant")
	}

	operator := token.GetText()

	if operator == "++" {
		c.ctx.Emitter.EmitA(vm.OpIncr, target)
	} else if operator == "--" {
		c.ctx.Emitter.EmitA(vm.OpDecr, target)
	} else {
		c.ctx.Errors.InvalidToken(token)

		return vm.NoopOperand
	}

	return target
}

// compileUnary processes a unary operation (NOT, minus, plus) from the FQL AST.
// It compiles the operand expression and applies the appropriate unary operation to it.
// Parameters:
//   - ctx: The unary operator context from the AST
//   - parent: The parent expression context containing the operand
//
// Returns:
//   - An operand representing the result of the unary operation
//
// Panics if the unary operator type is not recognized.
func (c *ExprCompiler) compileUnary(ctx fql.IUnaryOperatorContext, parent fql.IExpressionContext) vm.Operand {
	src := c.Compile(parent.GetRight())
	dst := c.ctx.Registers.Allocate()

	var op vm.Opcode

	if ctx.Not() != nil {
		op = vm.OpNot
	} else if ctx.Minus() != nil {
		op = vm.OpFlipNegative
	} else if ctx.Plus() != nil {
		op = vm.OpFlipPositive
	} else {
		return vm.NoopOperand
	}

	// We do not overwrite the source register
	c.ctx.Emitter.EmitAB(op, dst, src)

	return dst
}

// compileLogicalAnd processes a logical AND operation from the FQL AST.
// It implements short-circuit evaluation: if the left operand is falsy, it returns that value
// without evaluating the right operand. Otherwise, it evaluates and returns the right operand.
// Parameters:
//   - ctx: The expression context from the AST containing both operands
//
// Returns:
//   - An operand representing the result of the logical AND operation
func (c *ExprCompiler) compileLogicalAnd(ctx fql.IExpressionContext) vm.Operand {
	left := c.Compile(ctx.GetLeft())

	skip := c.ctx.Emitter.NewLabel("and.false")
	done := c.ctx.Emitter.NewLabel("and.done")
	dst := c.ctx.Registers.Allocate()

	// If left is falsy, jump to skip and use left
	c.ctx.Emitter.EmitJumpIfFalse(left, skip)

	// Otherwise evaluate right and use it
	right := c.Compile(ctx.GetRight())
	c.ctx.Emitter.EmitMove(dst, right)
	c.ctx.Emitter.EmitJump(done)

	// Short-circuit: use left as result
	c.ctx.Emitter.MarkLabel(skip)
	c.ctx.Emitter.EmitMove(dst, left)

	c.ctx.Emitter.MarkLabel(done)

	return dst
}

// compileLogicalOr processes a logical OR operation from the FQL AST.
// It implements short-circuit evaluation: if the left operand is truthy, it returns that value
// without evaluating the right operand. Otherwise, it evaluates and returns the right operand.
// Parameters:
//   - ctx: The expression context from the AST containing both operands
//
// Returns:
//   - An operand representing the result of the logical OR operation
func (c *ExprCompiler) compileLogicalOr(ctx fql.IExpressionContext) vm.Operand {
	left := c.Compile(ctx.GetLeft())

	next := c.ctx.Emitter.NewLabel("or.false")
	done := c.ctx.Emitter.NewLabel("or.done")
	dst := c.ctx.Registers.Allocate()

	// If left is truthy, short-circuit and skip right
	c.ctx.Emitter.EmitJumpIfTrue(left, next)

	// Otherwise evaluate right
	right := c.Compile(ctx.GetRight())
	c.ctx.Emitter.EmitMove(dst, right)
	c.ctx.Emitter.EmitJump(done)

	// Short-circuit: use left value
	c.ctx.Emitter.MarkLabel(next)
	c.ctx.Emitter.EmitMove(dst, left)

	// Common exit
	c.ctx.Emitter.MarkLabel(done)

	return dst
}

// compileTernary processes a ternary conditional operation (condition ? trueExpr : falseExpr) from the FQL AST.
// It evaluates the condition and then either the true expression or the false expression based on the result.
// Parameters:
//   - ctx: The expression context from the AST containing the condition, true expression, and false expression
//
// Returns:
//   - An operand representing the result of either the true or false expression
func (c *ExprCompiler) compileTernary(ctx fql.IExpressionContext) vm.Operand {
	dst := c.ctx.Registers.Allocate()

	// Compile condition and put result in dst
	condReg := c.Compile(ctx.GetCondition())
	c.ctx.Emitter.EmitMove(dst, condReg)

	// Define jump labels
	elseLabel := c.ctx.Emitter.NewLabel()
	endLabel := c.ctx.Emitter.NewLabel()

	// endLabel to 'false' branch if condition is false
	c.ctx.Emitter.EmitJumpIfFalse(dst, elseLabel)

	// True branch
	if onTrue := ctx.GetOnTrue(); onTrue != nil {
		trueReg := c.Compile(onTrue)
		// Move result of true branch to dst
		c.ctx.Emitter.EmitMove(dst, trueReg)
	}

	// endLabel over false branch
	c.ctx.Emitter.EmitJump(endLabel)
	// Mark label for 'else' branch
	c.ctx.Emitter.MarkLabel(elseLabel)

	// False branch
	if onFalse := ctx.GetOnFalse(); onFalse != nil {
		falseReg := c.Compile(onFalse)
		// Move result of false branch to dst
		c.ctx.Emitter.EmitMove(dst, falseReg)
	}

	// endLabel
	c.ctx.Emitter.MarkLabel(endLabel)

	return dst
}

// compilePredicate processes a predicate expression from the FQL AST.
// It handles both atomic expressions and comparison operations (equality, array operations, IN, LIKE).
// For atomic expressions, it also handles error catching with the TRY operator.
// Parameters:
//   - ctx: The predicate context from the AST
//
// Returns:
//   - An operand representing the result of the predicate expression
//
// Panics if the operator type is not recognized or not implemented.
func (c *ExprCompiler) compilePredicate(ctx fql.IPredicateContext) vm.Operand {
	if atom := ctx.ExpressionAtom(); atom != nil {
		startCatch := c.ctx.Emitter.Size()
		reg := c.compileAtom(atom)

		if atom.ErrorOperator() != nil {
			jump := -1
			endCatch := c.ctx.Emitter.Size()

			if fe := atom.ForExpression(); fe != nil {
				// Since FOR-IN loops depend on custom iterators,
				// We need to handle cleanup before exiting the loop.
				// TODO: Find a better way to handle this. The code assumes the knowledge of the internals of the FOR-IN loop.
				if fe.In() != nil {
					jump = endCatch - 1
				}
			}

			c.ctx.CatchTable.Push(startCatch, endCatch, jump)
		}

		return reg
	}

	var opcode vm.Opcode
	var isNegated bool
	dest := c.ctx.Registers.Allocate()
	left := c.compilePredicate(ctx.Predicate(0))
	right := c.compilePredicate(ctx.Predicate(1))

	if op := ctx.EqualityOperator(); op != nil {
		switch op.GetText() {
		case "==":
			opcode = vm.OpEq
		case "!=":
			opcode = vm.OpNe
		case ">":
			opcode = vm.OpGt
		case ">=":
			opcode = vm.OpGte
		case "<":
			opcode = vm.OpLt
		case "<=":
			opcode = vm.OpLte
		default:
			return vm.NoopOperand
		}
	} else if op := ctx.InOperator(); op != nil {
		opcode = vm.OpIn
		isNegated = op.Not() != nil
	} else if op := ctx.LikeOperator(); op != nil {
		opcode = vm.OpLike
		isNegated = op.Not() != nil
	} else if op := ctx.ArrayOperator(); op != nil {
		var pos int

		if op.All() != nil {
			pos = int(vm.OpAllEq)
		} else if op.Any() != nil {
			pos = int(vm.OpAnyEq)
		} else if op.None() != nil {
			pos = int(vm.OpNoneEq)
		}

		if eo := op.EqualityOperator(); eo != nil {
			switch eo.GetText() {
			case "!=":
				pos += int(vm.OpAllNe) - int(vm.OpAllEq)
			case ">":
				pos += int(vm.OpAllGt) - int(vm.OpAllEq)
			case ">=":
				pos += int(vm.OpAllGte) - int(vm.OpAllEq)
			case "<":
				pos += int(vm.OpAllLt) - int(vm.OpAllEq)
			case "<=":
				pos += int(vm.OpAllLte) - int(vm.OpAllEq)
			default:
				break
			}
		} else if inOp := op.InOperator(); inOp != nil {
			pos += int(vm.OpAllIn) - int(vm.OpAllEq)
		} else {
			return vm.NoopOperand
		}

		opcode = vm.Opcode(pos)
	}

	span := file.Span{Start: -1, End: -1}
	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = diagnostics.SpanFromRuleContext(prc)
	}
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABC(opcode, dest, left, right)

		if isNegated {
			// If the operator is negated, we need to invert the result
			c.ctx.Emitter.EmitAB(vm.OpNot, dest, dest)
		}
	})

	return dest
}

// compileAtom processes an atomic expression from the FQL AST.
// It handles various types of expressions including arithmetic operations (*, /, %, +, -),
// regular expression operations (=~, !~), function calls, range operators, literals,
// variables, member expressions, parameters, and nested expressions.
// Parameters:
//   - ctx: The expression atom context from the AST
//
// Returns:
//   - An operand representing the result of the atomic expression
//
// Panics if the expression type is not recognized.
func (c *ExprCompiler) compileAtom(ctx fql.IExpressionAtomContext) vm.Operand {
	var opcode vm.Opcode
	var isSet bool
	var isNegated bool
	var isRegexp bool

	if op := ctx.MultiplicativeOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "*":
			opcode = vm.OpMulti
		case "/":
			opcode = vm.OpDiv
		case "%":
			opcode = vm.OpMod
		default:
			return vm.NoopOperand
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "+":
			opcode = vm.OpAdd
		case "-":
			opcode = vm.OpSub
		default:
			return vm.NoopOperand
		}

	} else if op := ctx.RegexpOperator(); op != nil {
		isSet = true
		opcode = vm.OpRegexp
		isNegated = op.GetText() == "!~"
		isRegexp = true
	}

	if isSet {
		regLeft := c.compileAtom(ctx.ExpressionAtom(0))
		regRight := c.compileAtom(ctx.ExpressionAtom(1))
		dst := c.ctx.Registers.Allocate()

		span := file.Span{Start: -1, End: -1}
		if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			span = diagnostics.SpanFromRuleContext(prc)
		}
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitABC(opcode, dst, regLeft, regRight)

			if isNegated {
				// If the operator is negated, we need to invert the result
				c.ctx.Emitter.EmitAB(vm.OpNot, dst, dst)
			}
		})

		if isRegexp {
			if rightCtx := ctx.ExpressionAtom(1); rightCtx != nil {
				if lit := rightCtx.Literal(); lit != nil {
					if sl := lit.StringLiteral(); sl != nil {
						exp := parseStringLiteral(sl).String()

						// Verify that the expression is a valid regular expression
						if _, err := regexp.Compile(exp); err != nil {
							c.ctx.Errors.InvalidRegexExpression(ctx, exp)
						}
					} else {
						c.ctx.Errors.InvalidRegexExpression(ctx, lit.GetText())
					}
				}
			}
		}

		return dst
	}

	if fex := ctx.FunctionCallExpression(); fex != nil {
		return c.CompileFunctionCallExpression(fex)
	} else if r := ctx.RangeOperator(); r != nil {
		return c.CompileRangeOperator(r)
	} else if l := ctx.Literal(); l != nil {
		return c.ctx.LiteralCompiler.Compile(l)
	} else if v := ctx.Variable(); v != nil {
		return c.CompileVariable(v)
	} else if me := ctx.MemberExpression(); me != nil {
		return c.CompileMemberExpression(me)
	} else if p := ctx.Param(); p != nil {
		return c.CompileParam(p)
	} else if fe := ctx.ForExpression(); fe != nil {
		return c.ctx.LoopCompiler.Compile(fe)
	} else if wfe := ctx.WaitForExpression(); wfe != nil {
		return c.ctx.WaitCompiler.Compile(wfe)
	} else if e := ctx.Expression(); e != nil {
		return c.Compile(e)
	}

	//c.ctx.Errors.UnexpectedToken(ctx)

	return vm.NoopOperand
}

// CompileMemberExpression processes a member expression (e.g., obj.prop, arr[idx]) from the FQL AST.
// It handles property access for variables, parameters, object literals, array literals, and function calls.
// It supports both dot notation (obj.prop) and bracket notation (obj["prop"] or arr[idx]),
// as well as optional chaining with the ?. operator.
// Parameters:
//   - ctx: The member expression context from the AST
//
// Returns:
//   - An operand representing the value of the accessed property
func (c *ExprCompiler) CompileMemberExpression(ctx fql.IMemberExpressionContext) vm.Operand {
	mes := ctx.MemberExpressionSource()
	segments := ctx.AllMemberExpressionPath()

	if mes == nil || len(segments) == 0 {
		return vm.NoopOperand
	}

	src := c.compileMemberExpressionSource(mes, segments)
	if src == vm.NoopOperand {
		return src
	}

	return c.compileMemberExpressionSegments(src, segments)
}

func (c *ExprCompiler) compileMemberExpressionSource(mes fql.IMemberExpressionSourceContext, segments []fql.IMemberExpressionPathContext) vm.Operand {
	if v := mes.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if p := mes.Param(); p != nil {
		return c.CompileParam(p)
	}

	if ol := mes.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	if al := mes.ArrayLiteral(); al != nil {
		return c.ctx.LiteralCompiler.CompileArrayLiteral(al)
	}

	if fc := mes.FunctionCall(); fc != nil {
		// FOO()?.bar
		segment := segments[0]
		return c.CompileFunctionCall(fc, segment.ErrorOperator() != nil)
	}

	if fe := mes.ForExpression(); fe != nil {
		return c.ctx.LoopCompiler.Compile(fe)
	}

	if wfe := mes.WaitForExpression(); wfe != nil {
		return c.ctx.WaitCompiler.Compile(wfe)
	}

	if e := mes.Expression(); e != nil {
		return c.Compile(e)
	}

	return vm.NoopOperand
}

func (c *ExprCompiler) compileMemberExpressionSegments(src vm.Operand, segments []fql.IMemberExpressionPathContext) vm.Operand {
	if len(segments) == 0 {
		return src
	}

	for idx, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)
		if contraction := p.ArrayContraction(); contraction != nil {
			inlineTail, restTail := splitArrayOperatorTail(segments[idx+1:])
			result := c.compileArrayContraction(src, contraction, inlineTail)
			if len(restTail) == 0 {
				return result
			}
			return c.compileMemberExpressionSegments(result, restTail)
		}
		if expansion := p.ArrayExpansion(); expansion != nil {
			return c.compileArrayExpansionChain(src, expansion, segments[idx+1:])
		}
		if question := p.ArrayQuestionMark(); question != nil {
			inlineTail, restTail := splitArrayOperatorTail(segments[idx+1:])
			result := c.compileArrayQuestionMark(src, question, inlineTail)
			if len(restTail) == 0 {
				return result
			}
			return c.compileMemberExpressionSegments(result, restTail)
		}
		if apply := p.ArrayApply(); apply != nil {
			return c.compileArrayApply(src, apply, segments[idx+1:])
		}

		var src2 vm.Operand
		var constOperand bool
		srcType := operandType(c.ctx, src)

		if pn := p.PropertyName(); pn != nil {
			if constOp, ok := c.ctx.LiteralCompiler.CompilePropertyNameConst(pn); ok {
				src2 = constOp
				constOperand = true
			} else {
				src2 = c.ctx.LiteralCompiler.CompilePropertyName(pn)
			}
		} else if cpn := p.ComputedPropertyName(); cpn != nil {
			if val, ok := literalValueFromExpression(cpn.Expression()); ok {
				switch val.(type) {
				case *runtime.Array, *runtime.Object:
					// Keep array/object literals dynamic to preserve their stringified key value.
					src2 = c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn)
				default:
					src2 = c.ctx.Symbols.AddConstant(val)
					constOperand = true
				}
			} else {
				src2 = c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn)
			}
		}

		dst := c.ctx.Registers.Allocate()

		span := diagnostics.SpanFromRuleContext(p)
		c.ctx.Emitter.WithSpan(span, func() {
			optional := p.ErrorOperator() != nil
			var op vm.Opcode

			switch srcType {
			case core.TypeArray:
				if constOperand {
					if optional {
						op = vm.OpLoadIndexOptionalConst
					} else {
						op = vm.OpLoadIndexConst
					}
				} else if optional {
					op = vm.OpLoadIndexOptional
				} else {
					op = vm.OpLoadIndex
				}
			case core.TypeObject:
				if constOperand {
					if optional {
						op = vm.OpLoadKeyOptionalConst
					} else {
						op = vm.OpLoadKeyConst
					}
				} else if optional {
					op = vm.OpLoadKeyOptional
				} else {
					op = vm.OpLoadKey
				}
			default:
				if constOperand {
					if optional {
						op = vm.OpLoadPropertyOptionalConst
					} else {
						op = vm.OpLoadPropertyConst
					}
				} else if optional {
					op = vm.OpLoadPropertyOptional
				} else {
					op = vm.OpLoadProperty
				}
			}

			c.ctx.Emitter.EmitABC(op, dst, src, src2)
		})

		src = dst
	}

	return src
}

func splitArrayOperatorTail(segments []fql.IMemberExpressionPathContext) ([]fql.IMemberExpressionPathContext, []fql.IMemberExpressionPathContext) {
	if len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)
		if p.ArrayContraction() != nil || p.ArrayExpansion() != nil || p.ArrayQuestionMark() != nil {
			return nil, segments
		}
	}

	return segments, nil
}

func (c *ExprCompiler) compileArrayExpansionChain(src vm.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext) vm.Operand {
	inline := expansion.InlineExpression()

	if inline == nil {
		if next, rest := nextArrayExpansion(tail); next != nil {
			return c.compileArrayExpansionChain(src, next, rest)
		}

		return c.compileArrayExpansionChainWithFilters(src, expansion, tail, nil)
	}

	if !isFilterOnlyInline(inline) {
		tail = dropIdentityExpansions(tail)
		return c.compileArrayExpansionChainWithFilters(src, expansion, tail, nil)
	}

	extraFilters, rest := collectFilterOnlyTail(tail)
	return c.compileArrayExpansionChainWithFilters(src, expansion, rest, extraFilters)
}

func (c *ExprCompiler) compileArrayExpansionWithFilters(src vm.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext, extraFilters []fql.IExpressionContext) vm.Operand {
	span := diagnostics.SpanFromRuleContext(expansion)
	inline := expansion.InlineExpression()
	return c.compileArrayIteration(src, span, tail, inline, extraFilters)
}

func (c *ExprCompiler) compileArrayExpansionChainWithFilters(src vm.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext, extraFilters []fql.IExpressionContext) vm.Operand {
	inlineTail, restTail := splitArrayOperatorTail(tail)
	result := c.compileArrayExpansionWithFilters(src, expansion, inlineTail, extraFilters)
	if len(restTail) == 0 {
		return result
	}

	return c.compileMemberExpressionSegments(result, restTail)
}

func isFilterOnlyInline(inline fql.IInlineExpressionContext) bool {
	if inline == nil {
		return false
	}

	return inline.InlineFilter() != nil && inline.InlineLimit() == nil && inline.InlineReturn() == nil
}

func nextArrayExpansion(segments []fql.IMemberExpressionPathContext) (fql.IArrayExpansionContext, []fql.IMemberExpressionPathContext) {
	if len(segments) == 0 {
		return nil, segments
	}

	p := segments[0].(*fql.MemberExpressionPathContext)
	if expansion := p.ArrayExpansion(); expansion != nil {
		return expansion, segments[1:]
	}

	return nil, segments
}

func dropIdentityExpansions(segments []fql.IMemberExpressionPathContext) []fql.IMemberExpressionPathContext {
	for len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()
		if expansion == nil {
			break
		}

		if expansion.InlineExpression() != nil {
			break
		}

		segments = segments[1:]
	}

	return segments
}

func collectFilterOnlyTail(segments []fql.IMemberExpressionPathContext) ([]fql.IExpressionContext, []fql.IMemberExpressionPathContext) {
	extraFilters := make([]fql.IExpressionContext, 0)
	rest := segments

	for len(rest) > 0 {
		p := rest[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()
		if expansion == nil {
			break
		}

		inline := expansion.InlineExpression()
		if inline == nil {
			rest = rest[1:]
			continue
		}

		if !isFilterOnlyInline(inline) {
			break
		}

		filter := inline.InlineFilter()
		if filter != nil {
			extraFilters = append(extraFilters, filter.Expression())
		}

		rest = rest[1:]
	}

	return extraFilters, rest
}

func (c *ExprCompiler) compileArrayQuestionMark(src vm.Operand, question fql.IArrayQuestionMarkContext, tail []fql.IMemberExpressionPathContext) vm.Operand {
	span := diagnostics.SpanFromRuleContext(question)

	loop := &core.Loop{
		Kind:     core.ForInLoop,
		Type:     core.NormalLoop,
		Distinct: false,
		Allocate: false,
		Dst:      vm.NoopOperand,
		Src:      src,
	}

	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	loop.DeclareValueVar(core.PseudoVariable, c.ctx.Symbols, core.TypeAny)
	if loop.Value.IsRegister() {
		c.ctx.Types.Set(loop.Value, core.TypeAny)
	}

	count := c.ctx.Registers.Allocate()
	total := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitA(vm.OpLoadZero, count)
		c.ctx.Emitter.EmitA(vm.OpLoadZero, total)
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
	})

	// Track total elements
	c.ctx.Emitter.EmitA(vm.OpIncr, total)

	// Apply optional filter
	if filter := question.Expression(); filter != nil {
		cond := c.ctx.ExprCompiler.Compile(filter)
		label := c.ctx.Loops.Current().ContinueLabel()
		c.ctx.Emitter.EmitJumpIfFalse(cond, label)
	}

	// Count matches
	c.ctx.Emitter.EmitA(vm.OpIncr, count)

	loop.EmitFinalization(c.ctx.Emitter)

	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	result := c.compileArrayQuestionQuantifier(question, count, total)
	if len(tail) > 0 {
		result = c.compileMemberExpressionSegments(result, tail)
	}

	if result.IsRegister() {
		c.ctx.Types.Set(result, core.TypeBool)
	}

	return result
}

func (c *ExprCompiler) compileArrayQuestionQuantifier(question fql.IArrayQuestionMarkContext, count, total vm.Operand) vm.Operand {
	quant := question.ArrayQuestionQuantifier()
	zero := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitA(vm.OpLoadZero, zero)

	if quant == nil || quant.Any() != nil {
		return c.emitComparison(vm.OpGt, count, zero)
	}

	if quant.None() != nil {
		return c.emitComparison(vm.OpEq, count, zero)
	}

	if quant.All() != nil {
		return c.emitComparison(vm.OpEq, count, total)
	}

	values := quant.AllArrayQuestionQuantifierValue()
	if quant.At() != nil {
		if len(values) == 0 {
			return c.emitComparison(vm.OpGt, count, zero)
		}

		value := c.compileArrayQuestionQuantifierValue(values[0])
		return c.emitComparison(vm.OpGte, count, value)
	}

	if quant.Range() != nil && len(values) >= 2 {
		min := c.compileArrayQuestionQuantifierValue(values[0])
		max := c.compileArrayQuestionQuantifierValue(values[1])

		left := c.emitComparison(vm.OpGte, count, min)
		right := c.emitComparison(vm.OpLte, count, max)
		return c.emitBooleanAnd(left, right)
	}

	if len(values) > 0 {
		value := c.compileArrayQuestionQuantifierValue(values[0])
		return c.emitComparison(vm.OpEq, count, value)
	}

	return c.emitComparison(vm.OpGt, count, zero)
}

func (c *ExprCompiler) compileArrayQuestionQuantifierValue(ctx fql.IArrayQuestionQuantifierValueContext) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	if pm := ctx.Param(); pm != nil {
		return c.CompileParam(pm)
	}

	return vm.NoopOperand
}

func (c *ExprCompiler) compileArrayApply(src vm.Operand, apply fql.IArrayApplyContext, tail []fql.IMemberExpressionPathContext) vm.Operand {
	if apply == nil {
		return src
	}

	query := c.compileQueryLiteral(apply.QueryLiteral())
	if query == vm.NoopOperand {
		return vm.NoopOperand
	}

	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(apply)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABC(vm.OpApplyQuery, dst, src, query)
	})

	if len(tail) > 0 {
		dst = c.compileMemberExpressionSegments(dst, tail)
	}

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeAny)
	}

	return dst
}

func (c *ExprCompiler) compileQueryLiteral(ctx fql.IQueryLiteralContext) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	kind := ""
	if ident := ctx.Identifier(); ident != nil {
		kind = strings.ToLower(ident.GetText())
	}

	payload := runtime.EmptyString
	if str := ctx.StringLiteral(); str != nil {
		payload = parseStringLiteral(str)
	}
	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(ctx)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArray(dst, 3)
	})

	kindReg := loadConstant(c.ctx, runtime.NewString(kind))
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, kindReg)
	})

	payloadReg := loadConstant(c.ctx, payload)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, payloadReg)
	})

	params := ctx.Expression()
	var paramsReg vm.Operand
	if params == nil {
		paramsReg = loadConstant(c.ctx, runtime.None)
	} else {
		paramsReg = c.Compile(params)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, paramsReg)
	})

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeAny)
	}

	return dst
}

func (c *ExprCompiler) emitComparison(op vm.Opcode, left, right vm.Operand) vm.Operand {
	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(op, dst, left, right)
	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeBool)
	}
	return dst
}

func (c *ExprCompiler) emitBooleanAnd(left, right vm.Operand) vm.Operand {
	dst := c.ctx.Registers.Allocate()
	skip := c.ctx.Emitter.NewLabel("and.false")
	done := c.ctx.Emitter.NewLabel("and.done")

	c.ctx.Emitter.EmitJumpIfFalse(left, skip)
	c.ctx.Emitter.EmitJumpIfFalse(right, skip)
	c.ctx.Emitter.EmitAb(vm.OpLoadBool, dst, true)
	c.ctx.Emitter.EmitJump(done)

	c.ctx.Emitter.MarkLabel(skip)
	c.ctx.Emitter.EmitAb(vm.OpLoadBool, dst, false)
	c.ctx.Emitter.MarkLabel(done)

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeBool)
	}

	return dst
}

func (c *ExprCompiler) compileArrayExpansion(src vm.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext) vm.Operand {
	span := diagnostics.SpanFromRuleContext(expansion)
	inline := expansion.InlineExpression()
	return c.compileArrayIteration(src, span, tail, inline, nil)
}

func (c *ExprCompiler) compileArrayContraction(src vm.Operand, contraction fql.IArrayContractionContext, tail []fql.IMemberExpressionPathContext) vm.Operand {
	depth := arrayContractionDepth(contraction)
	if depth < 1 {
		depth = 1
	}

	span := diagnostics.SpanFromRuleContext(contraction)

	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABx(vm.OpFlatten, dst, src, depth)
	})

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeList)
	}

	inline := contraction.InlineExpression()

	if len(tail) == 0 && inline == nil {
		return dst
	}

	return c.compileArrayIteration(dst, span, tail, inline, nil)
}

func arrayContractionDepth(ctx fql.IArrayContractionContext) int {
	if ctx == nil {
		return 1
	}

	count := len(ctx.GetStars())
	if count > 1 {
		return count - 1
	}

	return 1
}

func (c *ExprCompiler) compileArrayIteration(src vm.Operand, span file.Span, tail []fql.IMemberExpressionPathContext, inline fql.IInlineExpressionContext, extraFilters []fql.IExpressionContext) vm.Operand {
	loop := &core.Loop{
		Kind:     core.ForInLoop,
		Type:     core.NormalLoop,
		Distinct: false,
		Allocate: true,
		Dst:      c.ctx.Registers.Allocate(),
		Src:      src,
	}

	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	loop.DeclareValueVar(core.PseudoVariable, c.ctx.Symbols, core.TypeAny)
	if loop.Value.IsRegister() {
		c.ctx.Types.Set(loop.Value, core.TypeAny)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
	})

	if inline != nil {
		c.compileInlineFilter(inline)
	}

	for _, expr := range extraFilters {
		c.compileInlineFilterExpr(expr)
	}

	if inline != nil {
		c.compileInlineLimit(inline)
	}

	projection := loop.Value
	if inline != nil {
		if ret := inline.InlineReturn(); ret != nil {
			projection = c.ctx.ExprCompiler.Compile(ret.Expression())
		}
	}

	if len(tail) > 0 {
		projection = c.compileMemberExpressionSegments(projection, tail)
	}

	c.ctx.Emitter.EmitAB(vm.OpPush, loop.Dst, projection)
	loop.EmitFinalization(c.ctx.Emitter)

	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	if loop.Dst.IsRegister() {
		c.ctx.Types.Set(loop.Dst, core.TypeList)
	}

	return loop.Dst
}

func (c *ExprCompiler) compileInlineFilter(inline fql.IInlineExpressionContext) {
	if inline == nil {
		return
	}

	filter := inline.InlineFilter()
	if filter == nil {
		return
	}

	src := c.ctx.ExprCompiler.Compile(filter.Expression())
	label := c.ctx.Loops.Current().ContinueLabel()
	c.ctx.Emitter.EmitJumpIfFalse(src, label)
}

func (c *ExprCompiler) compileInlineFilterExpr(expr fql.IExpressionContext) {
	if expr == nil {
		return
	}

	src := c.ctx.ExprCompiler.Compile(expr)
	label := c.ctx.Loops.Current().ContinueLabel()
	c.ctx.Emitter.EmitJumpIfFalse(src, label)
}

func (c *ExprCompiler) compileInlineLimit(inline fql.IInlineExpressionContext) {
	if inline == nil {
		return
	}

	limit := inline.InlineLimit()
	if limit == nil {
		return
	}

	clauses := limit.AllLimitClauseValue()
	if len(clauses) == 0 {
		return
	}

	if len(clauses) == 1 {
		c.ctx.LoopCompiler.compileLimit(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[0]))
		return
	}

	c.ctx.LoopCompiler.compileOffset(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[0]))
	c.ctx.LoopCompiler.compileLimit(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[1]))
}

// CompileVariable processes a variable reference from the FQL AST.
// It resolves the variable name in the symbol table and returns the register or constant
// containing its value. If the variable is not found, it panics with an error.
// Parameters:
//   - ctx: The variable context from the AST
//
// Returns:
//   - An operand representing the variable's value
//
// Panics if the variable is not found in the symbol table.
func (c *ExprCompiler) CompileVariable(ctx fql.IVariableContext) vm.Operand {
	// Check if the context is valid (in case of parser errors)
	var name string
	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else {
		return vm.NoopOperand
	}
	// Just return the register / constant index
	op, _, found := c.ctx.Symbols.Resolve(name)

	if !found {
		c.ctx.Errors.VariableNotFound(ctx.Identifier().GetSymbol(), name)

		return vm.NoopOperand
	}

	if op.IsRegister() {
		return op
	}

	reg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitMove(reg, op)

	return reg
}

// CompileParam processes a parameter reference (e.g., @paramName) from the FQL AST.
// It binds the parameter name in the symbol table and emits instructions to load
// the parameter value at runtime.
// Parameters:
//   - ctx: The parameter context from the AST
//
// Returns:
//   - An operand representing the parameter's value
func (c *ExprCompiler) CompileParam(ctx fql.IParamContext) vm.Operand {
	var name string
	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else {
		return vm.NoopOperand
	}
	reg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadParam(reg, c.ctx.Symbols.BindParam(name))
	c.ctx.Types.Set(reg, core.TypeAny)

	return reg
}

// CompileFunctionCallExpression processes a function call expression from the FQL AST.
// It handles both regular function calls and protected function calls (with TRY).
// Parameters:
//   - ctx: The function call expression context from the AST
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) vm.Operand {
	protected := ctx.ErrorOperator() != nil
	call := ctx.FunctionCall()

	return c.CompileFunctionCall(call, protected)
}

// CompileFunctionCall processes a function call from the FQL AST.
// It compiles the function arguments and delegates to CompileFunctionCallWith.
// Parameters:
//   - ctx: The function call context from the AST
//   - protected: Whether this is a protected call (with TRY)
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCall(ctx fql.IFunctionCallContext, protected bool) vm.Operand {
	return c.CompileFunctionCallWith(ctx, protected, c.CompileArgumentList(ctx.ArgumentList()))
}

// CompileFunctionCallWith processes a function call with pre-compiled arguments from the FQL AST.
// It extracts the function name and delegates to CompileFunctionCallByNameWith.
// Parameters:
//   - ctx: The function call context from the AST
//   - protected: Whether this is a protected call
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) CompileFunctionCallWith(ctx fql.IFunctionCallContext, protected bool, seq core.RegisterSequence) vm.Operand {
	name := getFunctionName(ctx, c.ctx.UseAliases)
	span := file.Span{Start: -1, End: -1}

	if ctx != nil {
		if fn := ctx.FunctionName(); fn != nil {
			span = diagnostics.SpanFromRuleContext(fn)

			if ns := ctx.Namespace(); ns != nil && ns.GetStart() != nil {
				span.Start = ns.GetStart().GetStart()
			}
		} else if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			span = diagnostics.SpanFromRuleContext(prc)
		}
	}

	var out vm.Operand

	c.ctx.Emitter.WithSpan(span, func() {
		out = c.CompileFunctionCallByNameWith(name, protected, seq)
	})

	return out
}

// CompileFunctionCallByNameWith processes a function call by name with pre-compiled arguments.
// It handles both built-in runtime functions (TYPENAME, LENGTH, WAIT) and user-defined functions.
// For built-in functions, it emits specialized instructions. For user-defined functions,
// it delegates to compileUserFunctionCallWith.
// Parameters:
//   - name: The function name
//   - protected: Whether this is a protected call (with TRY)
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
//
// Panics if a built-in function is called with an incorrect number of arguments.
func (c *ExprCompiler) CompileFunctionCallByNameWith(name runtime.String, protected bool, seq core.RegisterSequence) vm.Operand {
	switch name {
	case runtimeLength:
		dst := c.ctx.Registers.Allocate()

		if seq == nil || len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeLength+": expected 1 argument"))
		}

		c.ctx.Emitter.EmitAB(vm.OpLength, dst, seq[0])

		return dst
	case runtimeTypename:
		dst := c.ctx.Registers.Allocate()

		if seq == nil || len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeTypename+": expected 1 argument"))
		}

		c.ctx.Emitter.EmitAB(vm.OpType, dst, seq[0])

		return dst
	case runtimeWait:
		if len(seq) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeWait+": expected 1 argument"))
		}

		c.ctx.Emitter.EmitA(vm.OpSleep, seq[0])

		return seq[0]
	default:
		return c.compileUserFunctionCallWith(name, protected, seq)
	}
}

// compileUserFunctionCallWith processes a user-defined function call with pre-compiled arguments.
// It loads the function name as a constant, binds the function in the symbol table,
// and emits the appropriate call instruction based on the number of arguments and whether
// the call is protected.
// Parameters:
//   - name: The function name
//   - protected: Whether this is a protected call
//   - seq: The pre-compiled function arguments as a sequence of registers
//
// Returns:
//   - An operand representing the function call result
func (c *ExprCompiler) compileUserFunctionCallWith(name runtime.String, protected bool, seq core.RegisterSequence) vm.Operand {
	dest := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(name))
	c.ctx.Symbols.BindFunction(name.String(), len(seq))

	var opcode vm.Opcode
	var protectedOpcode vm.Opcode

	switch len(seq) {
	case 0:
		opcode = vm.OpCall0
		protectedOpcode = vm.OpProtectedCall0
	case 1:
		opcode = vm.OpCall1
		protectedOpcode = vm.OpProtectedCall1
	case 2:
		opcode = vm.OpCall2
		protectedOpcode = vm.OpProtectedCall2
	case 3:
		opcode = vm.OpCall3
		protectedOpcode = vm.OpProtectedCall3
	case 4:
		opcode = vm.OpCall4
		protectedOpcode = vm.OpProtectedCall4
	default:
		opcode = vm.OpCall
		protectedOpcode = vm.OpProtectedCall
	}

	if !protected {
		c.ctx.Emitter.EmitAs(opcode, dest, seq)
	} else {
		c.ctx.Emitter.EmitAs(protectedOpcode, dest, seq)
	}

	c.ctx.Types.Set(dest, core.TypeAny)
	return dest
}

// CompileArgumentList processes a list of function arguments from the FQL AST.
// It compiles each argument expression and moves the results into a contiguous sequence
// of registers, which is required for function calls, array literals, and object literals.
// Parameters:
//   - ctx: The argument list context from the AST
//
// Returns:
//   - A sequence of registers containing the compiled argument values,
//     or an empty sequence if there are no arguments
func (c *ExprCompiler) CompileArgumentList(ctx fql.IArgumentListContext) core.RegisterSequence {
	var seq core.RegisterSequence

	if ctx == nil {
		return seq
	}

	// Get all array element expressions
	exps := ctx.AllExpression()
	size := len(exps)

	if size > 0 {
		// Allocate seq for function arguments
		seq = c.ctx.Registers.AllocateSequence(size)

		// Evaluate each element into seq Registers
		for i, exp := range exps {
			// Compile expression and move to seq register
			srcReg := c.Compile(exp)

			// The reason we move is that the argument list must be a contiguous sequence of registers
			// Otherwise, we cannot compileInitialization neither a list nor an object literal with arguments
			c.ctx.Emitter.EmitMove(seq[i], srcReg)
			c.ctx.Types.Set(seq[i], operandType(c.ctx, srcReg))
		}
	}

	return seq
}

// CompileRangeOperator processes a range operator expression (e.g., 1..10) from the FQL AST.
// It compiles the start and end operands and emits instructions to create a range.
// Parameters:
//   - ctx: The range operator context from the AST
//
// Returns:
//   - An operand representing the compiled range
func (c *ExprCompiler) CompileRangeOperator(ctx fql.IRangeOperatorContext) vm.Operand {
	dst := c.ctx.Registers.Allocate()
	start := c.compileRangeOperand(ctx.GetLeft())
	end := c.compileRangeOperand(ctx.GetRight())

	span := file.Span{Start: -1, End: -1}
	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = diagnostics.SpanFromRuleContext(prc)
	}
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitRange(dst, start, end)
	})

	c.ctx.Types.Set(dst, core.TypeList)
	return dst
}

// compileRangeOperand processes a range operand (start or end value) from the FQL AST.
// It handles variables, parameters, and integer literals as valid range operands.
// Parameters:
//   - ctx: The range operand context from the AST
//
// Returns:
//   - An operand representing the compiled range operand
//
// Panics if the operand type is not recognized.
func (c *ExprCompiler) compileRangeOperand(ctx fql.IRangeOperandContext) vm.Operand {
	if v := ctx.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.CompileParam(p)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.CompileMemberExpression(me)
	}

	if fc := ctx.FunctionCallExpression(); fc != nil {
		return c.CompileFunctionCallExpression(fc)
	}

	return vm.NoopOperand
}
