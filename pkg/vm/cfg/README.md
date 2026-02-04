# Control Flow Graph (CFG) Package

The `cfg` package provides a control flow graph generator and analyzer for Ferret bytecode. It can be used for bytecode optimization, analysis, and visualization.

## Overview

A Control Flow Graph (CFG) is a representation of all paths that might be traversed through a program during its execution. The CFG consists of:
- **Basic Blocks**: Sequences of instructions with a single entry and exit point
- **Edges**: Control flow between basic blocks

## Features

### CFG Construction
- **Automatic Basic Block Identification**: Identifies leaders (start of basic blocks) based on:
  - Program entry point
  - Jump targets
  - Instructions following control flow instructions
- **Edge Creation**: Creates edges between blocks based on:
  - Unconditional jumps (OpJump)
  - Conditional jumps (OpJumpIfFalse, OpJumpIfTrue)
  - Fall-through execution
  - Return statements

### Analysis Capabilities
- **Reachability Analysis**: Find reachable and unreachable blocks (dead code detection)
- **Loop Detection**: Identify back edges that form loops
- **Dominator Analysis**: Calculate dominator trees for optimization opportunities
- **Visualization**: Export to Graphviz DOT format

## Usage

### Building a CFG

```go
import "github.com/MontFerret/ferret/pkg/vm/cfg"

// Create a builder for your bytecode program
builder := cfg.NewBuilder(program)

// Build the control flow graph
graph, err := builder.Build()
if err != nil {
    log.Fatal(err)
}
```

### Analyzing the CFG

```go
// Create an analyzer
analyzer := cfg.NewAnalyzer(graph)

// Find unreachable code (dead code)
unreachable := analyzer.FindUnreachableBlocks()

// Detect loops
backEdges := analyzer.FindBackEdges()

// Calculate dominators for optimization
dominators := analyzer.CalculateDominators()
```

### Visualization

```go
// Generate DOT format for Graphviz
dot := graph.ToDOT()

// Save to file
os.WriteFile("cfg.dot", []byte(dot), 0644)

// Convert to image using Graphviz:
// dot -Tpng cfg.dot -o cfg.png
```

## Data Structures

### BasicBlock

Represents a sequence of instructions with:
- Unique ID
- Start and end instruction indices
- List of instructions
- Successor and predecessor blocks

```go
type BasicBlock struct {
    ID           int
    Start        int
    End          int
    Instructions []vm.Instruction
    Successors   []*BasicBlock
    Predecessors []*BasicBlock
}
```

### ControlFlowGraph

Represents the complete control flow structure:

```go
type ControlFlowGraph struct {
    Entry  *BasicBlock   // Entry point
    Exit   *BasicBlock   // Virtual exit block
    Blocks []*BasicBlock // All basic blocks
}
```

## Use Cases

### Optimization
- **Dead Code Elimination**: Remove unreachable blocks
- **Loop Optimization**: Identify loops for loop-invariant code motion
- **Common Subexpression Elimination**: Use dominator information

### Analysis
- **Code Coverage**: Identify all possible execution paths
- **Complexity Metrics**: Calculate cyclomatic complexity
- **Security Analysis**: Detect potential control flow vulnerabilities

### Debugging
- **Visualization**: Generate visual representations of program flow
- **Path Analysis**: Trace execution paths through the program

## Example

```go
// Create a simple conditional program
program := &vm.Program{
    Bytecode: []vm.Instruction{
        vm.NewInstruction(vm.OpLoadBool, 0, 1),
        vm.NewInstruction(vm.OpJumpIfFalse, 4, 0),
        vm.NewInstruction(vm.OpLoadConst, 1, 0),
        vm.NewInstruction(vm.OpJump, 5),
        vm.NewInstruction(vm.OpLoadConst, 2, 0),
        vm.NewInstruction(vm.OpReturn, 0),
    },
}

// Build and analyze
builder := cfg.NewBuilder(program)
graph, _ := builder.Build()
analyzer := cfg.NewAnalyzer(graph)

// Find optimization opportunities
unreachable := analyzer.FindUnreachableBlocks()
loops := analyzer.FindBackEdges()

// Visualize
fmt.Println(graph.String())
```

## Implementation Details

### Leader Identification
The builder uses a standard algorithm to identify basic block leaders:
1. The first instruction is always a leader
2. Any instruction that is a jump target is a leader
3. Any instruction following a control flow instruction is a leader

### Edge Construction
Edges are created based on the control flow semantics:
- **OpJump**: Single edge to target
- **OpJumpIfFalse/OpJumpIfTrue**: Two edges (target and fall-through)
- **OpReturn**: Edge to virtual exit block
- **Other instructions**: Edge to next instruction

### Dominator Calculation
Uses an iterative dataflow algorithm to compute dominator sets and immediate dominators.

## Testing

The package includes comprehensive tests covering:
- Empty programs
- Single blocks
- Conditional branches
- Loops
- Unreachable code
- Analysis functions

Run tests with:
```bash
go test ./pkg/vm/cfg/...
```

## Future Enhancements

Potential additions:
- Post-dominator analysis
- Natural loop identification
- CFG-based optimization passes
- Integration with the compiler for automatic CFG generation
- Live variable analysis
- Reaching definitions analysis
