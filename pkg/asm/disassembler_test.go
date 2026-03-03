package asm

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDisassemble_UDFStartEndLabels(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{
			runtime.NewInt(1),
			runtime.NewInt(2),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Name: "A", Entry: 2, Registers: 2, Params: 0},
				{Name: "B", Entry: 4, Registers: 2, Params: 0},
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	expected := []string{
		"@udf.0.A.start:",
		"@udf.0.A.end:",
		"@udf.1.B.start:",
		"@udf.1.B.end:",
	}
	for _, label := range expected {
		if !strings.Contains(out, label) {
			t.Fatalf("missing label %q in output:\n%s", label, out)
		}
	}
}

func TestDisassemble_UDFEndLabelOnTailCall(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpTailCall, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{
			runtime.NewInt(0),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Name: "OUTER", Entry: 3, Registers: 2, Params: 0},
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	endPos := strings.Index(out, "@udf.0.OUTER.end:")
	tailPos := strings.Index(out, "4: TAILCALL")

	if endPos < 0 {
		t.Fatalf("missing UDF end label in output:\n%s", out)
	}

	if tailPos < 0 {
		t.Fatalf("missing TAILCALL instruction in output:\n%s", out)
	}

	if endPos > tailPos {
		t.Fatalf("UDF end label must appear before tail call line:\n%s", out)
	}
}

func TestDisassemble_UDFLabelsCoexistWithMetadataLabels(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{
			runtime.NewInt(1),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Name: "F", Entry: 2, Registers: 2, Params: 0},
			},
		},
		Metadata: bytecode.Metadata{
			Labels: map[int]string{
				2: "loop.1.1.start",
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	canonicalPos := strings.Index(out, "@loop.1.1.start:")
	udfPos := strings.Index(out, "@udf.0.F.start:")
	insPos := strings.Index(out, "2: LOADC")

	if canonicalPos < 0 || udfPos < 0 {
		t.Fatalf("expected both canonical and udf labels at ip 2:\n%s", out)
	}

	if insPos < 0 {
		t.Fatalf("missing instruction line at ip 2:\n%s", out)
	}

	if canonicalPos > udfPos {
		t.Fatalf("canonical label must be printed before udf labels:\n%s", out)
	}

	if udfPos > insPos {
		t.Fatalf("labels must be printed before instruction:\n%s", out)
	}
}

func TestDisassemble_UDFLabelsIgnoreInvalidEntries(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Constants: []runtime.Value{
			runtime.NewInt(1),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Name: "VALID", Entry: 1, Registers: 2, Params: 0},
				{Name: "NEG", Entry: -1, Registers: 2, Params: 0},
				{Name: "BIG", Entry: 10, Registers: 2, Params: 0},
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if !strings.Contains(out, "@udf.0.VALID.start:") {
		t.Fatalf("missing label for valid udf:\n%s", out)
	}

	if !strings.Contains(out, "@udf.0.VALID.end:") {
		t.Fatalf("missing end label for valid udf:\n%s", out)
	}

	if strings.Contains(out, "@udf.1.NEG.start:") || strings.Contains(out, "@udf.1.NEG.end:") {
		t.Fatalf("invalid negative entry must not produce boundary labels:\n%s", out)
	}

	if strings.Contains(out, "@udf.2.BIG.start:") || strings.Contains(out, "@udf.2.BIG.end:") {
		t.Fatalf("out-of-range entry must not produce boundary labels:\n%s", out)
	}
}
