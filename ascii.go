package a11y

import (
	"sort"
	"strings"
)

// Sanitize adapta o texto ao perfil.
//
// Sequências ANSI/OSC são sempre removidas. Quando o perfil não suporta
// Unicode, os glifos de caixa, blocos, setas, marcadores e sinais mais comuns
// são transliterados para ASCII. Glifos fora da tabela são preservados, pois
// não há um equivalente ASCII universal para eles.
func (p Profile) Sanitize(s string) string {
	out := StripANSI(s)
	if p.Unicode {
		return out
	}
	return transliterate(out)
}

// asciiFallbacks mapeia glifos Unicode comuns em interfaces de terminal para
// equivalentes ASCII. A preferência é por um caractere por glifo, para manter
// o alinhamento; símbolos compostos usam mais de um caractere quando isso é
// mais legível ("…" vira "..." e "→" vira "->").
//
// A tabela completa, por categoria, está documentada no README.
var asciiFallbacks = map[string]string{
	// Caixas de linha simples e grossa.
	"─": "-", "━": "-", "│": "|", "┃": "|",
	"┌": "+", "┍": "+", "┎": "+", "┏": "+",
	"┐": "+", "┑": "+", "┒": "+", "┓": "+",
	"└": "+", "┕": "+", "┖": "+", "┗": "+",
	"┘": "+", "┙": "+", "┚": "+", "┛": "+",
	"├": "+", "┝": "+", "┞": "+", "┟": "+", "┠": "+", "┡": "+", "┢": "+", "┣": "+",
	"┤": "+", "┥": "+", "┦": "+", "┧": "+", "┨": "+", "┩": "+", "┪": "+", "┫": "+",
	"┬": "+", "┭": "+", "┮": "+", "┯": "+", "┰": "+", "┱": "+", "┲": "+", "┳": "+",
	"┴": "+", "┵": "+", "┶": "+", "┷": "+", "┸": "+", "┹": "+", "┺": "+", "┻": "+",
	"┼": "+", "┽": "+", "┾": "+", "┿": "+", "╀": "+", "╁": "+", "╂": "+", "╃": "+",
	"╄": "+", "╅": "+", "╆": "+", "╇": "+", "╈": "+", "╉": "+", "╊": "+", "╋": "+",

	// Caixas de linha dupla.
	"═": "=", "║": "|",
	"╒": "+", "╓": "+", "╔": "+",
	"╕": "+", "╖": "+", "╗": "+",
	"╘": "+", "╙": "+", "╚": "+",
	"╛": "+", "╜": "+", "╝": "+",
	"╞": "+", "╟": "+", "╠": "+",
	"╡": "+", "╢": "+", "╣": "+",
	"╤": "+", "╥": "+", "╦": "+",
	"╧": "+", "╨": "+", "╩": "+",
	"╪": "+", "╫": "+", "╬": "+",

	// Cantos arredondados.
	"╭": "+", "╮": "+", "╰": "+", "╯": "+",

	// Junções e traços avulsos.
	"╴": "-", "╵": "|", "╶": "-", "╷": "|", "╸": "-", "╹": "|",
	"╺": "-", "╻": "|", "╼": "-", "╽": "|", "╾": "-", "╿": "|",
	"╳": "x",

	// Blocos, meios-blocos e sombreados.
	"█": "#", "▓": "#", "▒": ":", "░": ".",
	"▀": "^", "▄": "_", "▌": "[", "▐": "]",
	"▁": "_", "▂": "_", "▃": "_", "▅": "_", "▆": "_", "▇": "_",
	"▉": "#", "▊": "#", "▋": "#", "▍": "|", "▎": "|", "▏": "|", "▕": "|",
	"▔": "-", "▖": "+", "▗": "+", "▘": "+", "▙": "+", "▚": "+", "▛": "+",
	"▜": "+", "▝": "+", "▞": "+", "▟": "+",

	// Setas.
	"←": "<-", "→": "->", "↑": "^", "↓": "v", "↔": "<->", "↕": "^v",
	"↖": "\\", "↗": "/", "↘": "\\", "↙": "/",
	"⇐": "<=", "⇒": "=>", "⇔": "<=>",

	// Triângulos e ponteiros.
	"▴": "^", "▲": "^", "△": "^",
	"▾": "v", "▼": "v", "▽": "v", "▿": "v",
	"▸": ">", "▹": ">", "▶": ">", "▷": ">", "►": ">",
	"◂": "<", "◀": "<", "◄": "<", "◃": "<",

	// Marcadores e sinais.
	"●": "o", "○": "o", "◦": "o", "◉": "(o)", "◎": "(o)", "◯": "o",
	"✓": "v", "✔": "v", "✗": "x", "✘": "x", "✕": "x", "✖": "x",
	"☑": "[x]", "☒": "[x]",
	"•": "*", "·": ".", "∙": ".", "‣": ">", "⁃": "-",
	"■": "#", "□": "[]", "▪": "*", "▫": "*", "◆": "<>", "◇": "<>",
	"★": "*", "☆": "*",

	// Pontuação, sinais e símbolos de status.
	"…": "...", "–": "-", "—": "-", "―": "-", "‑": "-",
	"“": "\"", "”": "\"", "‘": "'", "’": "'",
	"«": "<<", "»": ">>",
	"≤": "<=", "≥": ">=", "≠": "!=", "≈": "~", "×": "x", "÷": "/", "±": "+/-",
	"⚠": "(!)", "ℹ": "(i)",
}

// asciiReplacer aplica asciiFallbacks. As chaves são ordenadas para que a
// construção seja determinística, ainda que o resultado não dependa da ordem
// para chaves distintas.
var asciiReplacer = newASCIIReplacer(asciiFallbacks)

// newASCIIReplacer monta um strings.Replacer a partir de um mapa, ordenando as
// chaves antes de inserir os pares.
func newASCIIReplacer(table map[string]string) *strings.Replacer {
	keys := make([]string, 0, len(table))
	for k := range table {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys)*2)
	for _, k := range keys {
		pairs = append(pairs, k, table[k])
	}
	return strings.NewReplacer(pairs...)
}

// transliterate troca os glifos conhecidos pelos equivalentes ASCII.
func transliterate(s string) string {
	if s == "" {
		return ""
	}
	return asciiReplacer.Replace(s)
}
