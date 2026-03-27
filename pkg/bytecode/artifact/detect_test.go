package artifact

import "testing"

func TestHasMagic(t *testing.T) {
	program := newArtifactTestProgram()
	artifactData, err := Marshal(program, Options{})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "marshal_output",
			data: artifactData,
			want: true,
		},
		{
			name: "exact_magic_prefix",
			data: []byte("FBC2"),
			want: true,
		},
		{
			name: "short_prefix",
			data: []byte("FBC"),
			want: false,
		},
		{
			name: "source_text",
			data: []byte("RETURN 1"),
			want: false,
		},
		{
			name: "unrelated_binary",
			data: []byte{0x00, 0x01, 0x02, 0x03, 0x04},
			want: false,
		},
		{
			name: "nil",
			data: nil,
			want: false,
		},
		{
			name: "empty",
			data: []byte{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasMagic(tt.data); got != tt.want {
				t.Fatalf("HasMagic() = %v, want %v", got, tt.want)
			}
		})
	}
}
