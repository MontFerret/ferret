package internal

import (
	"errors"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// LiteralCompiler handles the compilation of literal values in FQL queries.
// It transforms literal expressions from the AST into VM instructions and constants.
type LiteralCompiler struct {
	ctx   *CompilationSession
	exprs *ExprCompiler
	facts *TypeFacts
}

// NewLiteralCompiler creates a new instance of LiteralCompiler with the given compiler context.
func NewLiteralCompiler(ctx *CompilationSession) *LiteralCompiler {
	return &LiteralCompiler{
		ctx: ctx,
	}
}

func (c *LiteralCompiler) bind(exprs *ExprCompiler, facts *TypeFacts) {
	if c == nil {
		return
	}

	c.exprs = exprs
	c.facts = facts
}

// Compile processes a literal expression from the FQL AST and delegates to the appropriate
// compilation method based on the literal type (string, integer, float, boolean, array, object, or none).
// Parameters:
//   - ctx: The literal context from the AST
//
// Returns:
//   - An operand representing the compiled literal value
//
// Panics if the literal type is not recognized.
func (c *LiteralCompiler) Compile(ctx fql.ILiteralContext) bytecode.Operand {
	if sl := ctx.StringLiteral(); sl != nil {
		return c.CompileStringLiteral(sl)
	} else if il := ctx.IntegerLiteral(); il != nil {
		return c.CompileIntegerLiteral(il)
	} else if fl := ctx.FloatLiteral(); fl != nil {
		return c.CompileFloatLiteral(fl)
	} else if dl := ctx.DurationLiteral(); dl != nil {
		return c.CompileDurationLiteral(dl)
	} else if bl := ctx.BooleanLiteral(); bl != nil {
		return c.CompileBooleanLiteral(bl)
	} else if al := ctx.ArrayLiteral(); al != nil {
		return c.CompileArrayLiteral(al)
	} else if ol := ctx.ObjectLiteral(); ol != nil {
		return c.CompileObjectLiteral(ol)
	} else if nl := ctx.NoneLiteral(); nl != nil {
		return c.CompileNoneLiteral(nl)
	}

	return bytecode.NoopOperand
}

// CompileDurationLiteral parses a duration without using a floating-point intermediate.
func (c *LiteralCompiler) CompileDurationLiteral(ctx fql.IDurationLiteralContext) bytecode.Operand {
	value, err := runtime.ParseDuration(ctx.GetText())
	if err != nil {
		c.reportInvalidDurationLiteral(ctx, err)

		return bytecode.NoopOperand
	}

	return c.facts.LoadConstant(value)
}

func (c *LiteralCompiler) reportInvalidDurationLiteral(ctx antlr.ParserRuleContext, err error) {
	if c == nil || c.ctx == nil || c.ctx.Program.Errors == nil || ctx == nil {
		core.PanicInvariant("cannot report invalid duration literal")
	}

	message := "Invalid duration literal"
	hint := "Use a valid duration, e.g. 100ms, 2s, or 1.5m."

	if errors.Is(err, runtime.ErrRange) {
		message = "Duration literal is out of range"
		hint = "Use a duration value that fits within the signed nanosecond range."
	}

	diag := c.ctx.Program.Errors.Create(parserd.SyntaxError, ctx, message)
	diag.Hint = hint
	c.ctx.Program.Errors.Add(diag)
}

// CompileStringLiteral processes a string literal from the FQL AST and converts it into a runtime string.
// It handles escape sequences like \n and \t, and properly extracts the string content without quotes.
// Parameters:
//   - ctx: The string literal context from the AST
//
// Returns:
//   - An operand representing the compiled string constant
func (c *LiteralCompiler) CompileStringLiteral(ctx fql.IStringLiteralContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if tmpl := ctx.TemplateLiteral(); tmpl != nil {
		return c.CompileTemplateLiteral(tmpl)
	}

	// Create a runtime string and load it as a constant
	return c.facts.LoadConstant(parseStringLiteral(ctx))
}

// CompileTemplateLiteral processes a template literal from the FQL AST.
// It concatenates literal chunks and interpolated expressions into a string.
func (c *LiteralCompiler) CompileTemplateLiteral(ctx fql.ITemplateLiteralContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	elements := ctx.AllTemplateElement()
	if len(elements) == 0 {
		return c.facts.LoadConstant(runtime.EmptyString)
	}

	parts := make([]concatOperandSegment, 0, len(elements))
	var literal strings.Builder

	flushLiteral := func() {
		if literal.Len() == 0 {
			return
		}

		parts = append(parts, concatOperandSegment{
			literal: runtime.NewString(literal.String()),
		})
		literal.Reset()
	}

	for _, el := range elements {
		if el == nil {
			continue
		}

		if chunk := el.TemplateChars(); chunk != nil {
			literal.WriteString(parseTemplateChunk(chunk.GetText()))

			continue
		}

		if expr := el.Expression(); expr != nil {
			flushLiteral()
			parts = append(parts, buildConcatOperandSegmentsFromExpression(c.exprs, expr)...)
		}
	}

	flushLiteral()

	return emitConcatOperandSegments(c.ctx, c.facts, parts)
}

// CompileIntegerLiteral processes an integer literal from the FQL AST and converts it into a runtime integer.
func (c *LiteralCompiler) CompileIntegerLiteral(ctx fql.IIntegerLiteralContext) bytecode.Operand {
	val, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		c.reportInvalidNumericLiteral(ctx, "integer", err)
		return bytecode.NoopOperand
	}

	return c.facts.LoadConstant(runtime.NewInt(val))
}

// CompileFloatLiteral processes a float literal from the FQL AST and converts it into a runtime float.
func (c *LiteralCompiler) CompileFloatLiteral(ctx fql.IFloatLiteralContext) bytecode.Operand {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		c.reportInvalidNumericLiteral(ctx, "float", err)
		return bytecode.NoopOperand
	}

	return c.facts.LoadConstant(runtime.NewFloat(val))
}

func (c *LiteralCompiler) reportInvalidNumericLiteral(ctx antlr.ParserRuleContext, kind string, err error) {
	if c == nil || c.ctx == nil || c.ctx.Program.Errors == nil || ctx == nil {
		core.PanicInvariantf("cannot report invalid %s literal", kind)
	}

	message, hint := invalidNumericLiteralDetails(kind, err)
	diag := c.ctx.Program.Errors.Create(parserd.SyntaxError, ctx, message)
	diag.Hint = hint
	c.ctx.Program.Errors.Add(diag)
}

// CompileBooleanLiteral processes a boolean literal from the FQL AST and converts it into a runtime boolean.
// Parameters:
//   - ctx: The boolean literal context from the AST
//
// Returns:
//   - An operand representing the compiled boolean value
//
// Panics if the text is neither "true" nor "false".
func (c *LiteralCompiler) CompileBooleanLiteral(ctx fql.IBooleanLiteralContext) bytecode.Operand {
	// Allocate a temporary register for the boolean value
	reg := c.ctx.Function.Registers.Allocate()

	// Convert the text to lowercase and determine the boolean value
	switch strings.ToLower(ctx.GetText()) {
	case "true":
		c.ctx.Program.Emitter.EmitBoolean(reg, true)
	case "false":
		c.ctx.Program.Emitter.EmitBoolean(reg, false)
	default:
		reg = bytecode.NoopOperand
	}

	if reg.IsRegister() {
		c.ctx.Function.Types.Set(reg, core.TypeBool)
	}

	return reg
}

// CompileNoneLiteral processes a none literal (null/nil value) from the FQL AST.
// Parameters:
//   - _: The none literal context from the AST (unused)
//
// Returns:
//   - An operand representing the compiled none value
func (c *LiteralCompiler) CompileNoneLiteral(_ fql.INoneLiteralContext) bytecode.Operand {
	// Allocate a temporary register for the none value
	reg := c.ctx.Function.Registers.Allocate()
	// Emit instruction to load the none value into the register
	c.ctx.Program.Emitter.EmitA(bytecode.OpLoadNone, reg)
	c.ctx.Function.Types.Set(reg, core.TypeNone)

	return reg
}

// CompileArrayLiteral processes an array literal from the FQL AST and converts it into a runtime array.
// It compiles each element in the array and emits instructions to create the array.
// Parameters:
//   - ctx: The array literal context from the AST
//
// Returns:
//   - An operand representing the compiled array
func (c *LiteralCompiler) CompileArrayLiteral(ctx fql.IArrayLiteralContext) bytecode.Operand {
	destReg := c.ctx.Function.Registers.Allocate()
	entries := ctx.AllArrayEntry()
	c.ctx.Program.Emitter.EmitArray(destReg, len(entries))

	for _, entry := range entries {
		if spread := entry.SpreadEntry(); spread != nil {
			src := c.exprs.Compile(spread.Expression())

			c.ctx.Program.Emitter.WithSpan(parserd.SpanFromRuleContext(spread), func() {
				c.ctx.Program.Emitter.EmitArraySpread(destReg, src)
			})

			continue
		}

		item := c.exprs.Compile(entry.Expression())
		c.ctx.Program.Emitter.EmitArrayPush(destReg, item)
	}

	c.ctx.Function.Types.Set(destReg, core.TypeArray)

	return destReg
}

// CompileObjectLiteral processes an object literal from the FQL AST and converts it into a runtime object.
// It compiles each property-value pair in the object and emits instructions to create the object.
// Parameters:
//   - ctx: The object literal context from the AST
//
// Returns:
//   - An operand representing the compiled object
func (c *LiteralCompiler) CompileObjectLiteral(ctx fql.IObjectLiteralContext) bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()
	entries := ctx.AllObjectEntry()
	c.ctx.Program.Emitter.EmitObject(dst, len(entries))

	for _, entry := range entries {
		if spread := entry.SpreadEntry(); spread != nil {
			src := c.exprs.Compile(spread.Expression())

			c.ctx.Program.Emitter.WithSpan(parserd.SpanFromRuleContext(spread), func() {
				c.ctx.Program.Emitter.EmitObjectSpread(dst, src)
			})

			continue
		}

		c.compileObjectProperty(dst, entry.PropertyAssignment())
	}

	c.ctx.Function.Types.Set(dst, core.TypeObject)

	return dst
}

func (c *LiteralCompiler) compileObjectProperty(dst bytecode.Operand, ctx fql.IPropertyAssignmentContext) {
	if prop := ctx.PropertyName(); prop != nil {
		// Evaluate the value first to shorten the live range of the key register.
		value := c.exprs.Compile(ctx.Expression())
		if key, ok := c.CompilePropertyNameConst(prop); ok {
			c.ctx.Program.Emitter.EmitObjectSetConst(dst, key, value)

			return
		}

		key := c.CompilePropertyName(prop)
		c.ctx.Program.Emitter.EmitObjectSet(dst, key, value)

		return
	}

	if computed := ctx.ComputedPropertyName(); computed != nil {
		if value, ok := c.facts.LiteralValueFromExpression(computed.Expression()); ok {
			switch value.(type) {
			case *runtime.Array, *runtime.Object:
				// Fall back to the generic path to preserve side effects.
			default:
				valueOp := c.exprs.Compile(ctx.Expression())
				key := c.ctx.Function.Symbols.AddConstant(runtime.ToString(value))
				c.ctx.Program.Emitter.EmitObjectSetConst(dst, key, valueOp)

				return
			}
		}

		key := c.CompileComputedPropertyName(computed)
		value := c.exprs.Compile(ctx.Expression())
		c.ctx.Program.Emitter.EmitObjectSet(dst, key, value)

		return
	}

	if variable := ctx.Variable(); variable != nil {
		value := c.exprs.CompileVariable(variable)
		key := c.ctx.Function.Symbols.AddConstant(runtime.NewString(variable.GetText()))
		c.ctx.Program.Emitter.EmitObjectSetConst(dst, key, value)
	}
}

// CompilePropertyName processes a property name from an object literal in the FQL AST.
// It handles different types of property names including string literals, identifiers,
// and reserved words (both safe and unsafe).
// Parameters:
//   - ctx: The property name context from the AST
//
// Returns:
//   - An operand representing the compiled property name as a string constant
//
// Panics if the property name type is not recognized.
func (c *LiteralCompiler) CompilePropertyName(ctx fql.IPropertyNameContext) bytecode.Operand {
	// Handle string literal property names (e.g., { "property": value })
	if str := ctx.StringLiteral(); str != nil {
		if val, ok := parseStringLiteralConst(str); ok {
			return c.facts.LoadConstant(val)
		}
		return c.CompileStringLiteral(str)
	}

	var name string

	// Handle different types of identifier property names
	if id := ctx.Identifier(); id != nil {
		// Regular identifier (e.g., { property: value })
		name = id.GetText()
	} else if word := ctx.SafeReservedWord(); word != nil {
		// Safe reserved word (e.g., { return: value })
		name = word.GetText()
	} else if word := ctx.UnsafeReservedWord(); word != nil {
		// Unsafe reserved word (e.g., { for: value })
		name = word.GetText()
	} else {
		return bytecode.NoopOperand
	}

	// Create a runtime string from the property name and load it as a constant
	return c.facts.LoadConstant(runtime.NewString(name))
}

// CompilePropertyNameConst compiles a property name into a constant operand without emitting instructions.
// It returns (operand, true) when a constant can be produced, otherwise (NoopOperand, false).
func (c *LiteralCompiler) CompilePropertyNameConst(ctx fql.IPropertyNameContext) (bytecode.Operand, bool) {
	if ctx == nil {
		return bytecode.NoopOperand, false
	}

	// Handle string literal property names (e.g., { "property": value })
	if str := ctx.StringLiteral(); str != nil {
		if value, ok := parseStringLiteralConst(str); ok {
			return c.ctx.Function.Symbols.AddConstant(value), true
		}
		return bytecode.NoopOperand, false
	}

	var name string

	// Handle different types of identifier property names
	if id := ctx.Identifier(); id != nil {
		// Regular identifier (e.g., { property: value })
		name = id.GetText()
	} else if word := ctx.SafeReservedWord(); word != nil {
		// Safe reserved word (e.g., { return: value })
		name = word.GetText()
	} else if word := ctx.UnsafeReservedWord(); word != nil {
		// Unsafe reserved word (e.g., { for: value })
		name = word.GetText()
	} else {
		return bytecode.NoopOperand, false
	}

	return c.ctx.Function.Symbols.AddConstant(runtime.NewString(name)), true
}

// CompileComputedPropertyName processes a computed property name from an object literal in the FQL AST.
// Computed property names are expressions enclosed in square brackets (e.g., { [expr]: value }).
// Parameters:
//   - ctx: The computed property name context from the AST
//
// Returns:
//   - An operand representing the compiled expression that will evaluate to the property name
func (c *LiteralCompiler) CompileComputedPropertyName(ctx fql.IComputedPropertyNameContext) bytecode.Operand {
	// Delegate to the expression compiler to compile the expression inside the brackets
	return c.exprs.Compile(ctx.Expression())
}
