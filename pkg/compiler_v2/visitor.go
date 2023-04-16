package compiler_v2

import (
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/runtime_v2"
)

type (
	variable struct {
		name  string
		depth int
	}

	visitor struct {
		*fql.BaseFqlParserVisitor
		err            error
		src            string
		funcs          *core.Functions
		constantsIndex map[uint64]int
		locations      []core.Location
		bytecode       []runtime_v2.Opcode
		arguments      []int
		constants      []core.Value
		scope          int
		globals        map[string]int
		locals         []variable
	}
)

const (
	ignoreVarPseudoVariable = "_"
	waitPseudoVariable      = "CURRENT"
	waitScope               = "waitfor"
	forScope                = "for"
)

func newVisitor(src string, funcs *core.Functions) *visitor {
	v := new(visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.src = src
	v.funcs = funcs
	v.constantsIndex = make(map[uint64]int)
	v.locations = make([]core.Location, 0)
	v.bytecode = make([]runtime_v2.Opcode, 0)
	v.arguments = make([]int, 0)
	v.constants = make([]core.Value, 0)
	v.scope = 0
	v.globals = make(map[string]int)
	v.locals = make([]variable, 0)

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

func (v *visitor) VisitHead(ctx *fql.HeadContext) interface{} {
	return nil
}

func (v *visitor) VisitVariableDeclaration(ctx *fql.VariableDeclarationContext) interface{} {
	name := ignoreVarPseudoVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	index := v.declareVariable(name)

	ctx.Expression().Accept(v)

	v.defineVariable(index)

	return nil
}

func (v *visitor) VisitVariable(ctx *fql.VariableContext) interface{} {
	name := ctx.GetText()

	if name == waitPseudoVariable {
		return nil
	}

	if v.scope == 0 {
		index, ok := v.globals[name]

		if !ok {
			panic(core.Error(ErrVariableNotFound, name))
		}

		v.emit(runtime_v2.OpGetGlobal, index)
	}

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

	v.emitConstant(runtime_v2.OpConstant, values.NewString(b.String()))

	return nil
}

func (v *visitor) VisitIntegerLiteral(ctx *fql.IntegerLiteralContext) interface{} {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	v.emitConstant(runtime_v2.OpConstant, values.NewInt(val))

	return nil
}

func (v *visitor) VisitFloatLiteral(ctx *fql.FloatLiteralContext) interface{} {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	v.emitConstant(runtime_v2.OpConstant, values.NewFloat(val))

	return nil
}

func (v *visitor) VisitBooleanLiteral(ctx *fql.BooleanLiteralContext) interface{} {
	switch strings.ToLower(ctx.GetText()) {
	case "true":
		v.emit(runtime_v2.OpTrue)
	case "false":
		v.emit(runtime_v2.OpFalse)
	default:
		panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return nil
}

func (v *visitor) VisitNoneLiteral(ctx *fql.NoneLiteralContext) interface{} {
	v.emit(runtime_v2.OpNone)

	return nil
}

func (v *visitor) VisitLiteral(ctx *fql.LiteralContext) interface{} {
	if c := ctx.StringLiteral(); c != nil {
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

	v.emit(runtime_v2.OpReturn)

	return nil
}

func (v *visitor) VisitExpression(ctx *fql.ExpressionContext) interface{} {
	if op := ctx.UnaryOperator(); op != nil {

	} else if op := ctx.LogicalAndOperator(); op != nil {

	} else if op := ctx.LogicalOrOperator(); op != nil {

	} else if op := ctx.GetTernaryOperator(); op != nil {

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
			v.emit(runtime_v2.OpEq)
		case "!=":
			v.emit(runtime_v2.OpNeq)
		case ">":
			v.emit(runtime_v2.OpGt)
		case ">=":
			v.emit(runtime_v2.OpGte)
		case "<":
			v.emit(runtime_v2.OpLt)
		case "<=":
			v.emit(runtime_v2.OpLte)
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.ArrayOperator(); op != nil {

	} else if op := ctx.InOperator(); op != nil {
		ctx.Predicate(0).Accept(v)
		ctx.Predicate(1).Accept(v)

		v.emit(runtime_v2.OpIn)
	} else if op := ctx.LikeOperator(); op != nil {
		ctx.Predicate(0).Accept(v)
		ctx.Predicate(1).Accept(v)

		v.emit(runtime_v2.OpLike)
	} else if c := ctx.ExpressionAtom(); c != nil {
		c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitExpressionAtom(ctx *fql.ExpressionAtomContext) interface{} {
	if op := ctx.MultiplicativeOperator(); op != nil {
		ctx.ExpressionAtom(0).Accept(v)
		ctx.ExpressionAtom(1).Accept(v)

		switch op.GetText() {
		case "*":
			v.emit(runtime_v2.OpMulti)
		case "/":
			v.emit(runtime_v2.OpDiv)
		case "%":
			v.emit(runtime_v2.OpMod)
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		ctx.ExpressionAtom(0).Accept(v)
		ctx.ExpressionAtom(1).Accept(v)

		switch op.GetText() {
		case "+":
			v.emit(runtime_v2.OpAdd)
		case "-":
			v.emit(runtime_v2.OpSub)
		}
	} else if op := ctx.RegexpOperator(); op != nil {
		ctx.ExpressionAtom(0).Accept(v)
		ctx.ExpressionAtom(1).Accept(v)

		switch op.GetText() {
		case "=~":
			v.emit(runtime_v2.OpRegexpPositive)
		case "!~":
			v.emit(runtime_v2.OpRegexpNegative)
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
}

func (v *visitor) declareVariable(name string) int {
	//if name == ignoreVarPseudoVariable {
	//	return
	//}
	//
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
		if v.locals[i].depth != v.scope {
			break
		}

		if v.locals[i].name == name {
			panic(core.Error(ErrVariableNotUnique, name))
		}
	}

	v.locals = append(v.locals, variable{name, v.scope})

	return len(v.locals) - 1
}

// defineVariable defines a variable in the current scope.
func (v *visitor) defineVariable(index int) {
	if v.scope == 0 {
		v.emit(runtime_v2.OpDefineGlobal, index)
	}
}

// emitConstant emits an opcode with a constant argument.
func (v *visitor) emitConstant(op runtime_v2.Opcode, constant core.Value) {
	v.emit(op, v.addConstant(constant))
}

func (v *visitor) emit(op runtime_v2.Opcode, args ...int) {
	v.bytecode = append(v.bytecode, op)

	var arg int

	if len(args) > 0 {
		arg = args[0]
	}

	v.arguments = append(v.arguments, arg)
}

// addConstant adds a constant to the constants pool and returns its index.
// If the constant is a scalar, it will be deduplicated.
// If the constant is not a scalar, it will be added to the pool without deduplication.
func (v *visitor) addConstant(constant core.Value) int {
	var hash uint64

	if types.IsScalar(constant.Type()) {
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
