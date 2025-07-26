package internal

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// LiteralCompiler handles the compilation of literal values in FQL queries.
// It transforms literal expressions from the AST into VM instructions and constants.
type LiteralCompiler struct {
	ctx *CompilerContext
}

// NewLiteralCompiler creates a new instance of LiteralCompiler with the given compiler context.
func NewLiteralCompiler(ctx *CompilerContext) *LiteralCompiler {
	return &LiteralCompiler{
		ctx: ctx,
	}
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
func (lc *LiteralCompiler) Compile(ctx fql.ILiteralContext) vm.Operand {
	if c := ctx.StringLiteral(); c != nil {
		return lc.CompileStringLiteral(c)
	} else if c := ctx.IntegerLiteral(); c != nil {
		return lc.CompileIntegerLiteral(c)
	} else if c := ctx.FloatLiteral(); c != nil {
		return lc.CompileFloatLiteral(c)
	} else if c := ctx.BooleanLiteral(); c != nil {
		return lc.CompileBooleanLiteral(c)
	} else if c := ctx.ArrayLiteral(); c != nil {
		return lc.CompileArrayLiteral(c)
	} else if c := ctx.ObjectLiteral(); c != nil {
		return lc.CompileObjectLiteral(c)
	} else if c := ctx.NoneLiteral(); c != nil {
		return lc.CompileNoneLiteral(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

// CompileStringLiteral processes a string literal from the FQL AST and converts it into a runtime string.
// It handles escape sequences like \n and \t, and properly extracts the string content without quotes.
// Parameters:
//   - ctx: The string literal context from the AST
//
// Returns:
//   - An operand representing the compiled string constant
func (lc *LiteralCompiler) CompileStringLiteral(ctx fql.IStringLiteralContext) vm.Operand {
	var b strings.Builder

	// Process each child node in the string literal
	for _, child := range ctx.GetChildren() {
		tree := child.(antlr.TerminalNode)
		sym := tree.GetSymbol()
		input := sym.GetInputStream()

		if input == nil {
			continue
		}

		size := input.Size()
		// Skip the opening and closing quotes
		start := sym.GetStart() + 1
		stop := sym.GetStop() - 1

		// Ensure we don't go beyond the input size
		if stop >= size {
			stop = size - 1
		}

		if start < size && stop < size {
			// Process each character in the string
			for i := start; i <= stop; i++ {
				c := input.GetText(i, i)

				switch c {
				case "\\":
					// Handle escape sequences
					c2 := input.GetText(i, i+1)

					switch c2 {
					case "\\n":
						b.WriteString("\n")
					case "\\t":
						b.WriteString("\t")
					default:
						b.WriteString(c2)
					}

					// Skip the next character as it's part of the escape sequence
					i++
				default:
					// Add regular characters as-is
					b.WriteString(c)
				}
			}
		}
	}

	// Create a runtime string and load it as a constant
	return loadConstant(lc.ctx, runtime.NewString(b.String()))
}

// CompileIntegerLiteral processes an integer literal from the FQL AST and converts it into a runtime integer.
// Parameters:
//   - ctx: The integer literal context from the AST
//
// Returns:
//   - An operand representing the compiled integer constant
//
// Panics if the integer value cannot be parsed.
func (lc *LiteralCompiler) CompileIntegerLiteral(ctx fql.IIntegerLiteralContext) vm.Operand {
	// Parse the integer value from the text representation
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	// Create a runtime integer and load it as a constant
	return loadConstant(lc.ctx, runtime.NewInt(val))
}

// CompileFloatLiteral processes a float literal from the FQL AST and converts it into a runtime float.
// Parameters:
//   - ctx: The float literal context from the AST
//
// Returns:
//   - An operand representing the compiled float constant
//
// Panics if the float value cannot be parsed.
func (lc *LiteralCompiler) CompileFloatLiteral(ctx fql.IFloatLiteralContext) vm.Operand {
	// Parse the float value from the text representation with 64-bit precision
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	// Create a runtime float and load it as a constant
	return loadConstant(lc.ctx, runtime.NewFloat(val))
}

// CompileBooleanLiteral processes a boolean literal from the FQL AST and converts it into a runtime boolean.
// Parameters:
//   - ctx: The boolean literal context from the AST
//
// Returns:
//   - An operand representing the compiled boolean value
//
// Panics if the text is neither "true" nor "false".
func (lc *LiteralCompiler) CompileBooleanLiteral(ctx fql.IBooleanLiteralContext) vm.Operand {
	// Allocate a temporary register for the boolean value
	reg := lc.ctx.Registers.Allocate(core.Temp)

	// Convert the text to lowercase and determine the boolean value
	switch strings.ToLower(ctx.GetText()) {
	case "true":
		lc.ctx.Emitter.EmitBoolean(reg, true)
	case "false":
		lc.ctx.Emitter.EmitBoolean(reg, false)
	default:
		panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
	}

	return reg
}

// CompileNoneLiteral processes a none literal (null/nil value) from the FQL AST.
// Parameters:
//   - _: The none literal context from the AST (unused)
//
// Returns:
//   - An operand representing the compiled none value
func (lc *LiteralCompiler) CompileNoneLiteral(_ fql.INoneLiteralContext) vm.Operand {
	// Allocate a temporary register for the none value
	reg := lc.ctx.Registers.Allocate(core.Temp)
	// Emit instruction to load the none value into the register
	lc.ctx.Emitter.EmitA(vm.OpLoadNone, reg)

	return reg
}

// CompileArrayLiteral processes an array literal from the FQL AST and converts it into a runtime array.
// It compiles each element in the array and emits instructions to create the array.
// Parameters:
//   - ctx: The array literal context from the AST
//
// Returns:
//   - An operand representing the compiled array
func (lc *LiteralCompiler) CompileArrayLiteral(ctx fql.IArrayLiteralContext) vm.Operand {
	// Allocate destination register for the array
	destReg := lc.ctx.Registers.Allocate(core.Temp)
	// Compile the argument list (array elements) into a sequence of registers
	seq := lc.ctx.ExprCompiler.CompileArgumentList(ctx.ArgumentList())
	// Emit instruction to create an array from the sequence of registers
	lc.ctx.Emitter.EmitArray(destReg, seq)

	return destReg
}

// CompileObjectLiteral processes an object literal from the FQL AST and converts it into a runtime object.
// It compiles each property-value pair in the object and emits instructions to create the object.
// Parameters:
//   - ctx: The object literal context from the AST
//
// Returns:
//   - An operand representing the compiled object
func (lc *LiteralCompiler) CompileObjectLiteral(ctx fql.IObjectLiteralContext) vm.Operand {
	// Allocate destination register for the object
	dst := lc.ctx.Registers.Allocate(core.Temp)
	var seq core.RegisterSequence
	// Get all property assignments from the object literal
	assignments := ctx.AllPropertyAssignment()
	size := len(assignments)

	if size > 0 {
		// Allocate a sequence of registers for property-value pairs
		// We need two registers for each assignment (one for property name, one for value)
		seq = lc.ctx.Registers.AllocateSequence(len(assignments) * 2)

		// Process each property assignment
		for i := 0; i < size; i++ {
			var propOp vm.Operand
			var valOp vm.Operand
			pac := assignments[i]

			// Handle different types of property names
			if prop := pac.PropertyName(); prop != nil {
				// Regular property name (e.g., { name: value })
				propOp = lc.CompilePropertyName(prop)
				valOp = lc.ctx.ExprCompiler.Compile(pac.Expression())
			} else if comProp := pac.ComputedPropertyName(); comProp != nil {
				// Computed property name (e.g., { [expr]: value })
				propOp = lc.CompileComputedPropertyName(comProp)
				valOp = lc.ctx.ExprCompiler.Compile(pac.Expression())
			} else if variable := pac.Variable(); variable != nil {
				// Shorthand property (e.g., { variable })
				propOp = loadConstant(lc.ctx, runtime.NewString(variable.GetText()))
				valOp = lc.ctx.ExprCompiler.CompileVariable(variable)
			}

			// Calculate the index in the sequence for this property-value pair
			regIndex := i * 2

			// Move the property name and value to their respective registers in the sequence
			lc.ctx.Emitter.EmitMove(seq[regIndex], propOp)
			lc.ctx.Emitter.EmitMove(seq[regIndex+1], valOp)

			// Free source register if temporary
			// Note: This is commented out in the original code
			if propOp.IsRegister() {
				//lc.ctx.Registers.Free(propOp)
			}
		}
	}

	// Emit instruction to create an object from the sequence of property-value pairs
	lc.ctx.Emitter.EmitObject(dst, seq)

	return dst
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
func (lc *LiteralCompiler) CompilePropertyName(ctx fql.IPropertyNameContext) vm.Operand {
	// Handle string literal property names (e.g., { "property": value })
	if str := ctx.StringLiteral(); str != nil {
		return lc.CompileStringLiteral(str)
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
		panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
	}

	// Create a runtime string from the property name and load it as a constant
	return loadConstant(lc.ctx, runtime.NewString(name))
}

// CompileComputedPropertyName processes a computed property name from an object literal in the FQL AST.
// Computed property names are expressions enclosed in square brackets (e.g., { [expr]: value }).
// Parameters:
//   - ctx: The computed property name context from the AST
//
// Returns:
//   - An operand representing the compiled expression that will evaluate to the property name
func (lc *LiteralCompiler) CompileComputedPropertyName(ctx fql.IComputedPropertyNameContext) vm.Operand {
	// Delegate to the expression compiler to compile the expression inside the brackets
	return lc.ctx.ExprCompiler.Compile(ctx.Expression())
}
