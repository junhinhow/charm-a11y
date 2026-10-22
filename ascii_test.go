package a11y

import (
	"strings"
	"testing"
)

func TestTransliterateRequiredGlyphs(t *testing.T) {
	tests := []struct {
		glyph string
		want  string
	}{
		{"─", "-"},
		{"│", "|"},
		{"┌", "+"},
		{"┐", "+"},
		{"└", "+"},
		{"┘", "+"},
		{"├", "+"},
		{"┤", "+"},
		{"┬", "+"},
		{"┴", "+"},
		{"┼", "+"},
		{"╭", "+"},
		{"╮", "+"},
		{"╰", "+"},
		{"╯", "+"},
		{"▀", "^"},
		{"▄", "_"},
		{"█", "#"},
		{"░", "."},
		{"▒", ":"},
		{"▓", "#"},
		{"▸", ">"},
		{"▾", "v"},
		{"●", "o"},
		{"○", "o"},
		{"✓", "v"},
		{"✗", "x"},
		{"…", "..."},
		{"→", "->"},
		{"←", "<-"},
		{"↑", "^"},
		{"↓", "v"},
		// Glifos afins, além da lista mínima.
		{"━", "-"},
		{"┃", "|"},
		{"═", "="},
		{"║", "|"},
		{"▶", ">"},
		{"◀", "<"},
		{"•", "*"},
		{"≤", "<="},
		{"≥", ">="},
		{"⚠", "(!)"},
	}

	p := Profile{Unicode: false}
	for _, tt := range tests {
		t.Run(tt.glyph, func(t *testing.T) {
			if got := transliterate(tt.glyph); got != tt.want {
				t.Fatalf("transliterate(%q) = %q, quer %q", tt.glyph, got, tt.want)
			}
			if got := p.Sanitize("[" + tt.glyph + "]"); got != "["+tt.want+"]" {
				t.Fatalf("Sanitize([%s]) = %q, quer %q", tt.glyph, got, "["+tt.want+"]")
			}
		})
	}
}

func TestASCIITableIsCompleteAndASCII(t *testing.T) {
	for glyph, want := range asciiFallbacks {
		if got := transliterate(glyph); got != want {
			t.Errorf("transliterate(%q) = %q, quer %q", glyph, got, want)
		}
		for _, r := range want {
			if r > 0x7f {
				t.Errorf("equivalente de %q não é ASCII: %q", glyph, want)
			}
		}
	}
}

func TestSanitizeMixedText(t *testing.T) {
	p := Profile{Unicode: false}
	in := "\x1b[1m┌─┐\x1b[0m\n│✓│\n└─┘ → 100% …"
	want := "+-+\n|v|\n+-+ -> 100% ..."
	if got := p.Sanitize(in); got != want {
		t.Fatalf("Sanitize() = %q, quer %q", got, want)
	}
}

func TestSanitizePreservesUnicodeAndStripsANSI(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
	}{
		{name: "unicode", profile: Profile{Unicode: true}},
		{name: "sem cor mas unicode", profile: Profile{NoColor: true, Unicode: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.profile.Sanitize("\x1b[31m┌─┐ ✓ 世界\x1b[0m")
			want := "┌─┐ ✓ 世界"
			if got != want {
				t.Fatalf("Sanitize() = %q, quer %q", got, want)
			}
		})
	}
}

func TestSanitizeUnknownGlyphIsPreserved(t *testing.T) {
	p := Profile{Unicode: false}
	if got := p.Sanitize("世界"); got != "世界" {
		t.Fatalf("Sanitize() = %q, quer %q", got, "世界")
	}
}

func TestSanitizeEmpty(t *testing.T) {
	if got := (Profile{}).Sanitize(""); got != "" {
		t.Fatalf("Sanitize(\"\") = %q, quer vazio", got)
	}
}

func TestSanitizeIsDeterministic(t *testing.T) {
	p := Profile{Unicode: false}
	in := "┌─┐ ✓ → … █"
	want := p.Sanitize(in)
	for i := 0; i < 100; i++ {
		if got := p.Sanitize(in); got != want {
			t.Fatalf("iteração %d: Sanitize() = %q, quer %q", i, got, want)
		}
	}
}

func TestSanitizeDoesNotGrowControlText(t *testing.T) {
	p := Profile{Unicode: false}
	got := p.Sanitize("ASCII puro: a[0] -> b")
	if strings.Contains(got, "\x1b") {
		t.Fatalf("Sanitize deixou ESC: %q", got)
	}
	if got != "ASCII puro: a[0] -> b" {
		t.Fatalf("Sanitize alterou texto ASCII: %q", got)
	}
}
