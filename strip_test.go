package a11y

import (
	"strings"
	"testing"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "vazia", in: "", want: ""},
		{name: "texto simples", in: "olá, mundo", want: "olá, mundo"},
		{name: "SGR", in: "\x1b[31mvermelho\x1b[0m", want: "vermelho"},
		{name: "SGR composto", in: "\x1b[1;38;5;196mX\x1b[0m", want: "X"},
		{name: "CSI de cursor", in: "\x1b[2J\x1b[1;1Htela", want: "tela"},
		{name: "CSI privado", in: "\x1b[?25lsegredo\x1b[?25h", want: "segredo"},
		{name: "OSC com BEL", in: "\x1b]0;título\x07corpo", want: "corpo"},
		{name: "OSC com ST", in: "\x1b]8;;https://example.com\x1b\\link\x1b]8;;\x1b\\", want: "link"},
		{name: "sequência de dois bytes", in: "\x1b(Btexto", want: "texto"},
		{name: "sequência de dois bytes sem texto", in: "\x1b#8", want: ""},
		{name: "colchetes literais", in: "a[0] = [31m]", want: "a[0] = [31m]"},
		{name: "percentual", in: "100% [ok]", want: "100% [ok]"},
		{name: "multibyte", in: "\x1b[1mOlá, 世界! 🎉\x1b[0m", want: "Olá, 世界! 🎉"},
		{name: "ANSI no meio", in: "a\x1b[31mb\x1b[0mc", want: "abc"},
		{name: "CSI sem terminador", in: "\x1b[31", want: ""},
		{name: "OSC sem terminador", in: "\x1b]0;título", want: ""},
		{name: "ESC solto", in: "\x1b", want: ""},
		{name: "texto antes de OSC sem terminador", in: "ok\x1b]0;título", want: "ok"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripANSI(tt.in); got != tt.want {
				t.Fatalf("StripANSI(%q) = %q, quer %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestStripANSIIdempotent(t *testing.T) {
	inputs := []string{
		"\x1b[31mvermelho\x1b[0m",
		"\x1b]8;;https://example.com\x1b\\link\x1b]8;;\x1b\\",
		"sem ansi",
		"",
	}
	for _, in := range inputs {
		once := StripANSI(in)
		twice := StripANSI(once)
		if once != twice {
			t.Fatalf("StripANSI não é idempotente para %q: %q != %q", in, once, twice)
		}
	}
}

func TestStripANSIMalformedDoesNotPanic(t *testing.T) {
	sequences := []string{
		"\x1b[31m",
		"\x1b]0;título\x07",
		"\x1b]8;;url\x1b\\",
		"\x1b(B",
		"\x1b[?25l",
		"\x1bPq\x1b\\",
		"\x1b[38;2;255",
	}
	for _, seq := range sequences {
		for i := 0; i <= len(seq); i++ {
			_ = StripANSI(seq[:i])
		}
	}

	// Bytes inválidos não devem causar panic nem perda de texto ASCII.
	if got := StripANSI("\xff\xfe\x1b[31mok"); !strings.Contains(got, "ok") {
		t.Fatalf("StripANSI com UTF-8 inválido perdeu texto: %q", got)
	}
}

func TestStripANSIEmpty(t *testing.T) {
	if got := StripANSI(""); got != "" {
		t.Fatalf("StripANSI(\"\") = %q, quer vazio", got)
	}
}
