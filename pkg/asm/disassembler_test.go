package asm

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestDisassemble(t *testing.T) {
	type args struct {
		p *vm.Program
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Disassemble(tt.args.p); got != tt.want {
				t.Errorf("Disassemble() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collectLabels(t *testing.T) {
	type args struct {
		bytecode []vm.Instruction
	}
	tests := []struct {
		name string
		args args
		want map[int]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collectLabels(tt.args.bytecode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_disasmLine(t *testing.T) {
	type args struct {
		ip     int
		instr  vm.Instruction
		p      *vm.Program
		labels map[int]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := disasmLine(tt.args.ip, tt.args.instr, tt.args.p, tt.args.labels); got != tt.want {
				t.Errorf("disasmLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
