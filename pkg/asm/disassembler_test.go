package asm

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDisassemble_HeaderSectionsStackedLayout(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  4,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Params: []string{"value"},
		Constants: []runtime.Value{
			runtime.NewString("value"),
			runtime.NewInt(0),
			runtime.NewInt(1),
			runtime.NewString("X::DO"),
		},
		Functions: bytecode.Functions{
			Host: map[string]int{
				"X::DO": 1,
			},
			UserDefined: []bytecode.UDF{
				{Name: "FACT", Entry: 0, Registers: 10, Params: 1},
			},
		},
		Metadata: bytecode.Metadata{
			CompilerVersion:   "2.0.0",
			OptimizationLevel: 1,
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	metaPos := strings.Index(out, "\n.meta\n")
	paramsPos := strings.Index(out, "\n.params\n")
	constsPos := strings.Index(out, "\n.consts\n")
	udfPos := strings.Index(out, "\n.udf\n")
	funcPos := strings.Index(out, "\n.func\n")
	entryPos := strings.Index(out, "\n.entry\n")
	bodyPos := strings.Index(out, "0: RET R0")

	if metaPos < 0 || paramsPos < 0 || constsPos < 0 || udfPos < 0 || funcPos < 0 || entryPos < 0 || bodyPos < 0 {
		t.Fatalf("expected stacked header sections and body in output:\n%s", out)
	}

	if !(metaPos < paramsPos && paramsPos < constsPos && constsPos < udfPos && udfPos < funcPos && funcPos < entryPos && entryPos < bodyPos) {
		t.Fatalf("unexpected section order in output:\n%s", out)
	}

	expectedRows := []string{
		"\n  compiler 2.0.0\n",
		"\n  opt O1\n",
		"\n  value\n",
		"\n  \"value\"\n",
		"\n  0\n",
		"\n  1\n",
		"\n  \"X::DO\"\n",
		"\n  0 FACT 0 10 1  ; id name entry registers params\n",
		"\n  X::DO 1 ; name params\n",
	}

	for _, row := range expectedRows {
		if !strings.Contains(out, row) {
			t.Fatalf("missing expected row %q in output:\n%s", row, out)
		}
	}
}

func TestDisassemble_HeaderSections_OmitEmptySections(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			CompilerVersion:   "2.0.0",
			OptimizationLevel: 0,
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if !strings.Contains(out, "\n.meta\n") {
		t.Fatalf("expected .meta section in output:\n%s", out)
	}

	if strings.Contains(out, "\n.params\n") || strings.Contains(out, "\n.consts\n") || strings.Contains(out, "\n.udf\n") || strings.Contains(out, "\n.func\n") {
		t.Fatalf("expected empty sections to be omitted:\n%s", out)
	}
}

func TestDisassemble_EntryBeforeFirstLabelAndInstruction(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			Labels: map[int]string{
				0: "match.start",
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	entryPos := strings.Index(out, "\n.entry\n")
	labelPos := strings.Index(out, "\nmatch.start\n")
	insPos := strings.Index(out, "0: RET R0")

	if entryPos < 0 || labelPos < 0 || insPos < 0 {
		t.Fatalf("expected .entry, label, and first instruction in output:\n%s", out)
	}

	if !(entryPos < labelPos && labelPos < insPos) {
		t.Fatalf("expected .entry before label and label before instruction:\n%s", out)
	}
}

func TestDisassemble_EntryOmittedWhenNoBytecode(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  1,
		Metadata: bytecode.Metadata{
			CompilerVersion:   "2.0.0",
			OptimizationLevel: 0,
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if strings.Contains(out, "\n.entry\n") {
		t.Fatalf("expected .entry to be omitted for empty bytecode:\n%s", out)
	}
}

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
		"udf.0.A.start",
		"udf.0.A.end",
		"udf.1.B.start",
		"udf.1.B.end",
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

	endPos := strings.Index(out, "udf.0.OUTER.end")
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

	canonicalPos := strings.Index(out, "loop.1.1.start")
	udfPos := strings.Index(out, "udf.0.F.start")
	insPos := strings.Index(out, "2: LOADC")
	combined := "loop.1.1.start, udf.0.F.start"

	if canonicalPos < 0 || udfPos < 0 {
		t.Fatalf("expected both canonical and udf labels at ip 2:\n%s", out)
	}

	if !strings.Contains(out, combined) {
		t.Fatalf("expected stacked labels on one line %q:\n%s", combined, out)
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

	if !strings.Contains(out, "udf.0.VALID.start") {
		t.Fatalf("missing label for valid udf:\n%s", out)
	}

	if !strings.Contains(out, "udf.0.VALID.end") {
		t.Fatalf("missing end label for valid udf:\n%s", out)
	}

	if strings.Contains(out, "udf.1.NEG.start") || strings.Contains(out, "udf.1.NEG.end") {
		t.Fatalf("invalid negative entry must not produce boundary labels:\n%s", out)
	}

	if strings.Contains(out, "udf.2.BIG.start") || strings.Contains(out, "udf.2.BIG.end") {
		t.Fatalf("out-of-range entry must not produce boundary labels:\n%s", out)
	}
}

func TestDisassemble_LabelFormatting_NoColonAndBlankLineBetweenBlocks(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			Labels: map[int]string{
				0: "loop.1.1.start",
				2: "loop.1.1.end",
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if !strings.Contains(out, "\n\nloop.1.1.end\n") {
		t.Fatalf("expected blank line before second label block:\n%s", out)
	}

	definitionCount := 0
	for _, line := range strings.Split(out, "\n") {
		if line == "loop.1.1.start" || line == "loop.1.1.end" {
			definitionCount++
		}

		if strings.HasPrefix(line, "@") {
			t.Fatalf("label definition must not start with '@', got line %q\n%s", line, out)
		}

		if strings.HasPrefix(line, "@") && strings.HasSuffix(line, ":") {
			t.Fatalf("label definition must not end with colon, got line %q\n%s", line, out)
		}
	}

	if definitionCount == 0 {
		t.Fatalf("expected at least one label definition line in output:\n%s", out)
	}
}

func TestDisassemble_LabelDefinitionWithoutAt_LabelReferenceWithAt(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Metadata: bytecode.Metadata{
			Labels: map[int]string{
				2: "loop.1.1.end",
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if !strings.Contains(out, "\nloop.1.1.end\n") {
		t.Fatalf("expected label definition without '@':\n%s", out)
	}

	if !strings.Contains(out, "JMP @loop.1.1.end") {
		t.Fatalf("expected jump operand to keep '@' reference style:\n%s", out)
	}
}

func TestDisassemble_StackedUDFStartEndSingleInstruction(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
		},
		Functions: bytecode.Functions{
			UserDefined: []bytecode.UDF{
				{Name: "ONE", Entry: 0, Registers: 1, Params: 0},
			},
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	expected := "udf.0.ONE.start, udf.0.ONE.end"
	if !strings.Contains(out, expected) {
		t.Fatalf("expected stacked UDF start/end labels on one line %q:\n%s", expected, out)
	}
}

func TestDisassemble_HCallCommentWithHostFunctionName(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewRegister(2)),
		},
		Constants: []runtime.Value{
			runtime.NewString("X::DO"),
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if !strings.Contains(out, "1: HCALL R1 R2 R2 ; host X::DO") {
		t.Fatalf("expected HCALL host comment in output:\n%s", out)
	}
}

func TestDisassemble_ProtectedHCallCommentWithHostFunctionName(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpProtectedHCall, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewRegister(2)),
		},
		Constants: []runtime.Value{
			runtime.NewString("X::SAFE"),
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if !strings.Contains(out, "1: PHCALL R1 R2 R2 ; host X::SAFE") {
		t.Fatalf("expected protected HCALL host comment in output:\n%s", out)
	}
}

func TestDisassemble_HCallCommentMissingWhenNoMatchingLoadConst(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewRegister(2)),
		},
		Constants: []runtime.Value{
			runtime.NewString("X::DO"),
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if strings.Contains(out, "; host X::DO") {
		t.Fatalf("did not expect host comment when LOADC register does not match HCALL register:\n%s", out)
	}
}

func TestDisassemble_HCallCommentMissingWhenConstantIsNotString(t *testing.T) {
	prog := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewRegister(2)),
		},
		Constants: []runtime.Value{
			runtime.NewInt(42),
		},
	}

	out, err := Disassemble(prog)
	if err != nil {
		t.Fatalf("Disassemble() error: %v", err)
	}

	if strings.Contains(out, "; host ") {
		t.Fatalf("did not expect host comment when function-name constant is not a string:\n%s", out)
	}
}
