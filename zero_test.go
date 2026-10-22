package a11y

import (
	"testing"

	"charm.land/lipgloss/v2"
)

func TestZeroValueDoesNotPanic(t *testing.T) {
	var p Profile

	_ = p.String()
	_ = p.Describe()
	_ = p.Symbol("✓", "ok")
	_ = p.Sanitize("\x1b[31m┌─┐\x1b[0m")
	_ = p.Style(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Border(lipgloss.RoundedBorder()))
	_ = p.Border(lipgloss.RoundedBorder())
	_ = p.Frames(nil)
	_ = p.Frames([]string{"a"})
	_ = p.Announce("\x1b[1m%s\x1b[0m", "ok")
}

func TestDescribeIsStable(t *testing.T) {
	p := Profile{Unicode: true, Width: 100, HighContrast: true}
	want := "a11y: color=on unicode=on width=100 high-contrast"
	for i := 0; i < 100; i++ {
		if got := p.Describe(); got != want {
			t.Fatalf("iteração %d: Describe() = %q, quer %q", i, got, want)
		}
	}
}

func TestDetectFromDeterministic(t *testing.T) {
	env := map[string]string{
		"NO_COLOR":            "1",
		"A11Y_REDUCED_MOTION": "true",
		"COLUMNS":             "72",
	}
	getenv := envGetter(env)
	want := DetectFrom(getenv)
	for i := 0; i < 100; i++ {
		if got := DetectFrom(getenv); got != want {
			t.Fatalf("iteração %d: DetectFrom() = %#v, quer %#v", i, got, want)
		}
	}
}
