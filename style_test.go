package a11y

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestStyleNoColorRemovesColors(t *testing.T) {
	base := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Background(lipgloss.Color("#00ff00")).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#0000ff")).
		MarginBackground(lipgloss.Color("#111111")).
		Margin(0, 1)

	if rendered := base.Render("x"); !strings.Contains(rendered, "\x1b[") {
		t.Fatal("estilo de controle deveria conter ANSI")
	}

	styled := (Profile{NoColor: true, Unicode: true}).Style(base)
	if got := styled.GetForeground(); got != (lipgloss.NoColor{}) {
		t.Fatalf("Foreground = %#v, quer NoColor{}", got)
	}
	if got := styled.GetBackground(); got != (lipgloss.NoColor{}) {
		t.Fatalf("Background = %#v, quer NoColor{}", got)
	}
	if got := styled.GetBorderTopForeground(); got != (lipgloss.NoColor{}) {
		t.Fatalf("BorderTopForeground = %#v, quer NoColor{}", got)
	}
	if got := styled.GetBorderTopBackground(); got != (lipgloss.NoColor{}) {
		t.Fatalf("BorderTopBackground = %#v, quer NoColor{}", got)
	}

	rendered := styled.Render("x")
	if strings.Contains(rendered, "\x1b[") {
		t.Fatalf("estilo sem cor ainda contém ANSI: %q", rendered)
	}
	if !strings.Contains(rendered, "x") {
		t.Fatalf("estilo sem cor perdeu o texto: %q", rendered)
	}
	if !strings.Contains(rendered, "─") {
		t.Fatalf("Unicode deveria preservar a borda: %q", rendered)
	}

	// O estilo original deve permanecer intacto.
	if again := base.Render("x"); !strings.Contains(again, "\x1b[") {
		t.Fatal("Style modificou o estilo de entrada")
	}
}

func TestStyleNoColorRemovesAllColor(t *testing.T) {
	base := lipgloss.NewStyle().
		Underline(true).
		UnderlineColor(lipgloss.Color("#ff0000")).
		Border(lipgloss.NormalBorder()).
		BorderForegroundBlend(lipgloss.Color("#ff0000"), lipgloss.Color("#00ff00"))

	styled := (Profile{NoColor: true, Unicode: true}).Style(base)
	rendered := styled.Render("x")
	for _, prefix := range []string{"\x1b[38;", "\x1b[48;", "\x1b[58;"} {
		if strings.Contains(rendered, prefix) {
			t.Fatalf("NoColor deixou cor (%q): %q", prefix, rendered)
		}
	}
	if got := styled.GetUnderlineColor(); got != (lipgloss.NoColor{}) {
		t.Fatalf("UnderlineColor = %#v, quer NoColor{}", got)
	}
	if !styled.GetUnderline() {
		t.Fatal("NoColor não deveria remover o sublinhado")
	}
}

func TestStyleHighContrast(t *testing.T) {
	base := lipgloss.NewStyle().Faint(true)
	styled := (Profile{HighContrast: true, Unicode: true}).Style(base)

	if styled.GetFaint() {
		t.Fatal("HighContrast deveria remover faint")
	}
	if !styled.GetBold() {
		t.Fatal("HighContrast deveria aplicar negrito")
	}
	if got := styled.GetForeground(); got != (lipgloss.NoColor{}) {
		t.Fatalf("HighContrast introduziu cor de texto: %#v", got)
	}
	if got := styled.GetBackground(); got != (lipgloss.NoColor{}) {
		t.Fatalf("HighContrast introduziu cor de fundo: %#v", got)
	}

	rendered := styled.Render("x")
	if !strings.Contains(rendered, "\x1b[1m") {
		t.Fatalf("saída sem negrito: %q", rendered)
	}
	for _, prefix := range []string{"\x1b[38;", "\x1b[48;", "\x1b[58;"} {
		if strings.Contains(rendered, prefix) {
			t.Fatalf("HighContrast introduziu cor (%q): %q", prefix, rendered)
		}
	}
}

func TestStyleNoColorAndHighContrast(t *testing.T) {
	base := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Faint(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00ff00"))

	styled := (Profile{NoColor: true, HighContrast: true, Unicode: true}).Style(base)
	if !styled.GetBold() {
		t.Fatal("negrito deveria ser mantido")
	}
	if styled.GetFaint() {
		t.Fatal("faint deveria ser removido")
	}

	rendered := styled.Render("x")
	for _, prefix := range []string{"\x1b[38;", "\x1b[48;", "\x1b[58;"} {
		if strings.Contains(rendered, prefix) {
			t.Fatalf("combinação introduziu cor (%q): %q", prefix, rendered)
		}
	}
}

func TestStyleASCIIBorder(t *testing.T) {
	base := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	styled := (Profile{Unicode: false}).Style(base)

	want := lipgloss.Border{
		Top:          "-",
		Bottom:       "-",
		Left:         "|",
		Right:        "|",
		TopLeft:      "+",
		TopRight:     "+",
		BottomLeft:   "+",
		BottomRight:  "+",
		MiddleLeft:   "+",
		MiddleRight:  "+",
		Middle:       "+",
		MiddleTop:    "+",
		MiddleBottom: "+",
	}
	if got := styled.GetBorderStyle(); got != want {
		t.Fatalf("borda ASCII = %#v, quer %#v", got, want)
	}

	rendered := styled.Render("x")
	if strings.ContainsAny(rendered, "╭╮╰╯─│") {
		t.Fatalf("saída ainda contém borda Unicode: %q", rendered)
	}
	for _, part := range []string{"+", "-", "|"} {
		if !strings.Contains(rendered, part) {
			t.Fatalf("borda ASCII sem %q: %q", part, rendered)
		}
	}
}

func TestStylePreservesBorderSides(t *testing.T) {
	base := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false, false)
	styled := (Profile{Unicode: false}).Style(base)

	if !styled.GetBorderTop() {
		t.Fatal("borda superior deveria continuar habilitada")
	}
	if styled.GetBorderRight() || styled.GetBorderBottom() || styled.GetBorderLeft() {
		t.Fatal("Style habilitou lados que estavam desligados")
	}
}

func TestStyleUnicodeBorderUnchanged(t *testing.T) {
	base := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	styled := (Profile{Unicode: true}).Style(base)
	if got, want := styled.GetBorderStyle(), lipgloss.RoundedBorder(); got != want {
		t.Fatalf("borda Unicode = %#v, quer %#v", got, want)
	}
}

func TestBorder(t *testing.T) {
	ascii := lipgloss.Border{
		Top:          "-",
		Bottom:       "-",
		Left:         "|",
		Right:        "|",
		TopLeft:      "+",
		TopRight:     "+",
		BottomLeft:   "+",
		BottomRight:  "+",
		MiddleLeft:   "+",
		MiddleRight:  "+",
		Middle:       "+",
		MiddleTop:    "+",
		MiddleBottom: "+",
	}

	tests := []struct {
		name    string
		profile Profile
		in      lipgloss.Border
		want    lipgloss.Border
	}{
		{name: "unicode normal", profile: Profile{Unicode: true}, in: lipgloss.NormalBorder(), want: lipgloss.NormalBorder()},
		{name: "unicode arredondada", profile: Profile{Unicode: true}, in: lipgloss.RoundedBorder(), want: lipgloss.RoundedBorder()},
		{name: "ascii normal", profile: Profile{}, in: lipgloss.NormalBorder(), want: ascii},
		{name: "ascii arredondada", profile: Profile{}, in: lipgloss.RoundedBorder(), want: ascii},
		{name: "ascii grossa", profile: Profile{}, in: lipgloss.ThickBorder(), want: ascii},
		{
			name:    "ascii dupla",
			profile: Profile{},
			in:      lipgloss.DoubleBorder(),
			want: lipgloss.Border{
				Top: "=", Bottom: "=", Left: "|", Right: "|",
				TopLeft: "+", TopRight: "+", BottomLeft: "+", BottomRight: "+",
				MiddleLeft: "+", MiddleRight: "+", Middle: "+",
				MiddleTop: "+", MiddleBottom: "+",
			},
		},
		{
			name:    "ascii blocos",
			profile: Profile{},
			in:      lipgloss.BlockBorder(),
			want: lipgloss.Border{
				Top: "#", Bottom: "#", Left: "#", Right: "#",
				TopLeft: "#", TopRight: "#", BottomLeft: "#", BottomRight: "#",
				MiddleLeft: "#", MiddleRight: "#", Middle: "#",
				MiddleTop: "#", MiddleBottom: "#",
			},
		},
		{
			name:    "ascii meio-bloco externo",
			profile: Profile{},
			in:      lipgloss.OuterHalfBlockBorder(),
			want: lipgloss.Border{
				Top: "^", Bottom: "_", Left: "[", Right: "]",
				TopLeft: "+", TopRight: "+", BottomLeft: "+", BottomRight: "+",
			},
		},
		{name: "borda zero", profile: Profile{}, in: lipgloss.Border{}, want: lipgloss.Border{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.profile.Border(tt.in); got != tt.want {
				t.Fatalf("Border() = %#v, quer %#v", got, tt.want)
			}
		})
	}
}

func TestStyleZeroValue(t *testing.T) {
	base := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Border(lipgloss.RoundedBorder())
	var p Profile
	got := p.Style(base)
	if got.GetForeground() == (lipgloss.NoColor{}) {
		t.Fatal("perfil zero não declara NoColor; a cor deve ser preservada")
	}
	if got.GetBorderStyle() == lipgloss.RoundedBorder() {
		t.Fatal("perfil zero deveria converter a borda para ASCII")
	}
}
