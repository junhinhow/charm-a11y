package a11y

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

func ExampleStripANSI() {
	fmt.Println(StripANSI("\x1b[31mvermelho\x1b[0m"))
	// Output: vermelho
}

func ExampleProfile_Sanitize() {
	p := Profile{Unicode: false}
	fmt.Println(p.Sanitize("┌─┐"))
	fmt.Println(p.Sanitize("✓ ok…"))
	// Output:
	// +-+
	// v ok...
}

func ExampleProfile_Symbol() {
	ascii := Profile{Unicode: false}
	unicode := Profile{Unicode: true}
	fmt.Println(ascii.Symbol("✓", "ok"))
	fmt.Println(unicode.Symbol("✓", "ok"))
	// Output:
	// ok
	// ✓
}

func ExampleProfile_Announce() {
	p := Profile{ScreenReader: true}
	fmt.Println(p.Announce("\x1b[1m%s\x1b[0m\npronto", "build"))
	// Output: build pronto
}

func ExampleProfile_Style() {
	p := Profile{NoColor: true, Unicode: false}
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Border(lipgloss.RoundedBorder())
	fmt.Print(p.Style(style).Render("ok"))
	// Output:
	// +--+
	// |ok|
	// +--+
}
