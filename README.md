# charm-a11y

Accessibility profiles and graceful degradation for Charm terminal applications: respect `NO_COLOR`, fall back to ASCII, disable animation, and emit linear text for screen readers.

## Install

```bash
go get github.com/junhinhow/charm-a11y
```

## Usage

```go
package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/junhinhow/charm-a11y"
)

func main() {
	p := a11y.Default()

	title := p.Style(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		Border(lipgloss.RoundedBorder()))

	fmt.Println(title.Render(p.Symbol("✓", "ok") + " pronto"))
	fmt.Println(p.Describe())
}
```

With `TERM=dumb` the same program prints an ASCII border, no color, and a plain-text status line. With `A11Y_REDUCED_MOTION=1`, `p.Frames` returns a single frame.

## API

| Symbol | Description |
| --- | --- |
| `Detect() Profile` | Reads the process environment once. |
| `DetectFrom(getenv func(string) string) Profile` | Same detection with an injectable environment; `nil` means empty. |
| `Default() Profile` | Process-wide profile detected on first use. |
| `Profile` | Capabilities: `NoColor`, `Dumb`, `HighContrast`, `ReducedMotion`, `ScreenReader`, `Unicode`, `Width`. |
| `StripANSI(s string) string` | Removes CSI, OSC, two-byte escapes, and related control sequences. |
| `(Profile) Sanitize(s string) string` | Strips ANSI and, without Unicode, transliterates glyphs to ASCII. |
| `(Profile) Symbol(unicode, ascii string) string` | Picks the Unicode or ASCII variant. |
| `(Profile) Style(s lipgloss.Style) lipgloss.Style` | Removes color, boosts contrast, and converts borders to ASCII. |
| `(Profile) Border(b lipgloss.Border) lipgloss.Border` | Returns an ASCII-safe border when needed. |
| `(Profile) Frames(frames []string) []string` | Returns one frame under reduced motion. |
| `(Profile) Announce(format string, args ...any) string` | Formats a linear, ANSI-free message for screen readers. |
| `(Profile) Describe() string` | One-line, stable profile summary for `--version` or debug. |
| `(Profile) String() string` | Compact stable form used by `Describe`. |

## Environment variables

| Variable | Accepted values | Effect |
| --- | --- | --- |
| `NO_COLOR` | Any non-empty value | `NoColor`; also `Unicode=false` |
| `CLICOLOR` | `0` | `NoColor`; also `Unicode=false` |
| `TERM` | `dumb` (case-insensitive) | `NoColor`, `Dumb`, `ScreenReader`; `Unicode=false` |
| `ASCII` | `1` or `true` | `Unicode=false` |
| `A11Y_ASCII` | `1` or `true` | `Unicode=false` |
| `A11Y_HIGH_CONTRAST` | `1` or `true` | `HighContrast=true` |
| `A11Y_REDUCED_MOTION` | `1` or `true` | `ReducedMotion=true` |
| `A11Y_SCREEN_READER` | `1` or `true` | `ScreenReader=true` |
| `COLUMNS` | Positive integer | Fills `Width`; otherwise `0` (unknown) |

Notes:

- Per the [no-color.org](https://no-color.org) specification, `NO_COLOR` takes effect only when set to a **non-empty** value. `NO_COLOR=` does nothing; `NO_COLOR=0` still disables color because it is non-empty.
- The `A11Y_*` and `ASCII` values are trimmed and matched case-insensitively.
- `Unicode` is deliberately disabled whenever color is disabled, not only for `TERM=dumb`. This is the conservative choice for pipes, old emulators, and minimal environments. There is currently no override that forces Unicode while color is off.

## ASCII transliteration

`Sanitize` and `Border` use the following table. Unlisted glyphs are preserved as-is because there is no universal ASCII equivalent. Multi-character replacements are intentional; width changes are possible.

### Box drawing

| Unicode | ASCII |
| --- | --- |
| `─` `━` | `-` |
| `│` `┃` | `|` |
| `┌` `┍` `┎` `┏` `┐` `┑` `┒` `┓` `└` `┕` `┖` `┗` `┘` `┙` `┚` `┛` `├` `┝` `┞` `┟` `┠` `┡` `┢` `┣` `┤` `┥` `┦` `┧` `┨` `┩` `┪` `┫` `┬` `┭` `┮` `┯` `┰` `┱` `┲` `┳` `┴` `┵` `┶` `┷` `┸` `┹` `┺` `┻` `┼` `┽` `┾` `┿` `╀` `╁` `╂` `╃` `╄` `╅` `╆` `╇` `╈` `╉` `╊` `╋` `╭` `╮` `╰` `╯` | `+` |
| `═` | `=` |
| `║` | `|` |
| `╒` `╓` `╔` `╕` `╖` `╗` `╘` `╙` `╚` `╛` `╜` `╝` `╞` `╟` `╠` `╡` `╢` `╣` `╤` `╥` `╦` `╧` `╨` `╩` `╪` `╫` `╬` | `+` |
| `╴` `╶` `╸` `╺` `╼` `╾` | `-` |
| `╵` `╷` `╹` `╻` `╽` `╿` | `|` |
| `╳` | `x` |

### Blocks and shades

| Unicode | ASCII |
| --- | --- |
| `█` `▓` `▉` `▊` `▋` | `#` |
| `▒` | `:` |
| `░` | `.` |
| `▀` | `^` |
| `▄` `▁` `▂` `▃` `▅` `▆` `▇` | `_` |
| `▌` | `[` |
| `▐` | `]` |
| `▍` `▎` `▏` `▕` | `|` |
| `▔` | `-` |
| `▖` `▗` `▘` `▙` `▚` `▛` `▜` `▝` `▞` `▟` | `+` |

### Arrows and triangles

| Unicode | ASCII |
| --- | --- |
| `←` | `<-` |
| `→` | `->` |
| `↑` | `^` |
| `↓` | `v` |
| `↔` | `<->` |
| `↕` | `^v` |
| `↖` `↘` | `\` |
| `↗` `↙` | `/` |
| `⇐` | `<=` |
| `⇒` | `=>` |
| `⇔` | `<=>` |
| `▴` `▲` `△` | `^` |
| `▾` `▼` `▽` `▿` | `v` |
| `▸` `▹` `▶` `▷` `►` | `>` |
| `◂` `◀` `◄` `◃` | `<` |

### Marks, punctuation, and symbols

| Unicode | ASCII |
| --- | --- |
| `●` `○` `◦` `◯` | `o` |
| `◉` `◎` | `(o)` |
| `✓` `✔` | `v` |
| `✗` `✘` `✕` `✖` | `x` |
| `☑` `☒` | `[x]` |
| `•` | `*` |
| `·` `∙` | `.` |
| `‣` | `>` |
| `⁃` | `-` |
| `■` | `#` |
| `□` | `[]` |
| `▪` `▫` `★` `☆` | `*` |
| `◆` `◇` | `<>` |
| `…` | `...` |
| `–` `—` `―` `‑` | `-` |
| `“` `”` | `"` |
| `‘` `’` | `'` |
| `«` | `<<` |
| `»` | `>>` |
| `≤` | `<=` |
| `≥` | `>=` |
| `≠` | `!=` |
| `≈` | `~` |
| `×` | `x` |
| `÷` | `/` |
| `±` | `+/-` |
| `⚠` | `(!)` |
| `ℹ` | `(i)` |

## Accessibility checklist for terminal apps

- Detect the profile before the first render and pass it through the UI tree.
- Never convey state with color alone; pair color with text or a distinct glyph.
- Test with `NO_COLOR=1`, `TERM=dumb`, `ASCII=1`, and an 80-column terminal.
- Provide an ASCII path for every custom icon, spinner, and border.
- Honor `A11Y_REDUCED_MOTION`: no spinners, marquees, or rapid redraws; show a static status.
- Keep output linear and descriptive; announce meaningful state changes instead of only repainting.
- Support keyboard-only navigation and keep a visible focus indicator.
- Avoid blinking, flashing, and scrolling text.
- Avoid low-contrast gray-on-black text; test with `A11Y_HIGH_CONTRAST=1`.
- Do not require a mouse or hover to reach functionality.
- Do not redraw the full screen unnecessarily; screen readers cope better with appended, stable text.
- Make errors and exit status available as plain text.
- Document the environment variables your app honors.

## Design notes

- **Environment only.** Detection never queries the terminal. It is deterministic, works in CI, and is trivially testable with `DetectFrom`.
- **Conservative Unicode.** `NoColor` implies `Unicode=false`. This follows the project specification and protects limited environments, but it means a color-capable terminal in a color-disabled pipeline also loses box drawing. Unknown glyphs are preserved by `Sanitize`, so content is never silently dropped.
- **No color guessing.** `HighContrast` removes `faint` and adds `bold`; it never selects foreground or background colors. Choosing a high-contrast palette requires knowing the terminal background, which this package cannot detect reliably.
- **ANSI stripping** delegates to `github.com/charmbracelet/x/ansi`, which covers CSI, OSC (BEL and ST), DCS, APC, PM, SOS, and two-byte escapes. Malformed or unterminated sequences are consumed to the end of the string without panicking.
- **Announce is not a screen-reader driver.** Terminals expose rendered text to assistive technology; `Announce` reduces a message to stable, linear, ANSI-free text that such tools read more predictably.
- **Width is advisory.** `COLUMNS` is the only source; `0` means unknown. No ioctl is performed.
- **Frames is zero-allocation.** Under reduced motion it returns a subslice of the input; callers must not rely on a copy.
- **Not included:** color-contrast calculation, automatic terminal background detection, BCP-47 localization, and full Unicode normalization.

## License

MIT
