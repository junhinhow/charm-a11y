package a11y

import (
	"strings"
	"testing"
)

func TestAnnounce(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		format  string
		args    []any
		want    string
	}{
		{
			name:   "sem ANSI",
			format: "tudo pronto",
			want:   "tudo pronto",
		},
		{
			name:   "remove ANSI e lineariza",
			format: "\x1b[32m%s\x1b[0m\n%s\t!",
			args:   []any{"ok", "tudo certo"},
			want:   "ok tudo certo !",
		},
		{
			name:   "colapsa espaços",
			format: "a  b\n\nc",
			want:   "a b c",
		},
		{
			name:   "remove OSC",
			format: "\x1b]0;título\x07corpo",
			want:   "corpo",
		},
		{
			name:   "remove controles",
			format: "a\x00b\x07c",
			want:   "a b c",
		},
		{
			name:   "percentual literal sem argumentos",
			format: "100% [pronto]\x1b[0m",
			want:   "100% [pronto]",
		},
		{
			name:   "unicode preservado",
			format: "ação concluída: ✓",
			want:   "ação concluída: ✓",
		},
		{
			name:    "screen reader",
			profile: Profile{ScreenReader: true},
			format:  "etapa %d de %d",
			args:    []any{2, 3},
			want:    "etapa 2 de 3",
		},
		{
			name:   "vazio",
			format: "",
			want:   "",
		},
		{
			name:   "só espaços",
			format: " \n\t ",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.profile.Announce(tt.format, tt.args...)
			if got != tt.want {
				t.Fatalf("Announce() = %q, quer %q", got, tt.want)
			}
			if strings.Contains(got, "\x1b") {
				t.Fatalf("Announce() deixou ANSI: %q", got)
			}
		})
	}
}

func TestAnnounceDeterministic(t *testing.T) {
	p := Profile{ScreenReader: true}
	in := "\x1b[1m%s\x1b[0m\npronto"
	want := p.Announce(in, "build")
	for i := 0; i < 100; i++ {
		if got := p.Announce(in, "build"); got != want {
			t.Fatalf("iteração %d: Announce() = %q, quer %q", i, got, want)
		}
	}
}
