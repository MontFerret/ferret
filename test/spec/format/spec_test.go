package format

import "testing"

func TestSExactOutput(t *testing.T) {
	const input = "  RETURN 1\r\n"

	for _, tc := range []struct{ name, expected string }{
		{"empty", ""},
		{"ordinary", "return 1"},
		{"spaces", "  return 1 "},
		{"one trailing LF", "return 1\n"},
		{"multiple trailing LFs", "return 1\n\n\n"},
		{"CRLF", "let x = 1\r\nreturn x\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := S(input, tc.expected, "first", "second")
			if !s.Output.Result.Defined() {
				t.Fatal("missing output assertion")
			}

			if got := s.Output.Result.Value; got != tc.expected {
				t.Errorf("expected output stored as %q, want %q", got, tc.expected)
			}

			if got := s.Base.Input.Expression; got != input {
				t.Errorf("input = %q, want %q", got, input)
			}

			if got := s.Base.Description; got != "first second" {
				t.Errorf("description = %q, want %q", got, "first second")
			}
		})
	}
}
