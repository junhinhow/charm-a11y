package a11y

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Profile descreve as capacidades do terminal/ambiente.
//
// O valor zero representa um ambiente sem restrições conhecidas, exceto que
// Unicode é falso: campos booleanos zero não afirmam suporte. Prefira [Detect]
// para obter um perfil preenchido a partir do ambiente.
type Profile struct {
	NoColor       bool // NO_COLOR ou CLICOLOR=0 ou TERM=dumb
	Dumb          bool // TERM=dumb
	HighContrast  bool // A11Y_HIGH_CONTRAST=1
	ReducedMotion bool // A11Y_REDUCED_MOTION=1
	ScreenReader  bool // A11Y_SCREEN_READER=1 ou TERM=dumb
	Unicode       bool // false quando NoColor/Dumb/ASCII=1
	Width         int  // 0 = desconhecido
}

// Detect monta o perfil a partir das variáveis de ambiente do processo.
func Detect() Profile {
	return DetectFrom(os.Getenv)
}

// DetectFrom monta o perfil usando getenv, o que permite testar a detecção
// sem alterar o ambiente global. Um getenv nil é tratado como ambiente vazio.
//
// Regras aplicadas, nesta ordem de leitura:
//
//   - NO_COLOR definido e não vazio desliga cores, conforme no-color.org.
//   - CLICOLOR=0 desliga cores.
//   - TERM=dumb desliga cores, marca Dumb e ScreenReader e desliga Unicode.
//   - ASCII=1 ou A11Y_ASCII=1 (ou "true") desliga Unicode.
//   - A11Y_HIGH_CONTRAST, A11Y_REDUCED_MOTION e A11Y_SCREEN_READER aceitam
//     "1" ou "true", sem diferenciar maiúsculas de minúsculas.
//   - COLUMNS, quando inteiro positivo, preenche Width; caso contrário fica 0.
//
// Qualquer condição que desligue cores também desliga Unicode. Essa escolha é
// deliberadamente conservadora: ambientes sem cor costumam ser terminais
// limitados, pipes ou emuladores antigos.
func DetectFrom(getenv func(string) string) Profile {
	if getenv == nil {
		getenv = func(string) string { return "" }
	}

	term := strings.TrimSpace(getenv("TERM"))
	dumb := strings.EqualFold(term, "dumb")

	noColor := getenv("NO_COLOR") != "" ||
		strings.TrimSpace(getenv("CLICOLOR")) == "0" ||
		dumb
	ascii := truthy(getenv("ASCII")) || truthy(getenv("A11Y_ASCII"))

	return Profile{
		NoColor:       noColor,
		Dumb:          dumb,
		HighContrast:  truthy(getenv("A11Y_HIGH_CONTRAST")),
		ReducedMotion: truthy(getenv("A11Y_REDUCED_MOTION")),
		ScreenReader:  truthy(getenv("A11Y_SCREEN_READER")) || dumb,
		Unicode:       !noColor && !dumb && !ascii,
		Width:         columns(getenv("COLUMNS")),
	}
}

// truthy informa se v é "1" ou "true", ignorando espaços e caixa.
func truthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true":
		return true
	default:
		return false
	}
}

// columns converte COLUMNS em largura, devolvendo 0 para valores inválidos ou
// não positivos.
func columns(v string) int {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

// defaultProfile detecta o perfil uma única vez, na primeira chamada.
var defaultProfile = sync.OnceValue(Detect)

// Default devolve o perfil detectado uma vez no início do processo.
//
// O resultado é estável durante toda a vida do processo: mudanças posteriores
// no ambiente não o alteram.
func Default() Profile {
	return defaultProfile()
}

// String devolve uma representação compacta e estável do perfil, adequada para
// logs e testes.
//
// O formato é uma sequência de campos separados por espaço, sempre na mesma
// ordem: color=on|off, unicode=on|off, width=<n>|unknown e, por fim, os
// marcadores dumb, high-contrast, reduced-motion e screen-reader quando ativos.
func (p Profile) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "color=%s unicode=%s width=", onOff(!p.NoColor), onOff(p.Unicode))
	if p.Width > 0 {
		b.WriteString(strconv.Itoa(p.Width))
	} else {
		b.WriteString("unknown")
	}
	for _, flag := range [...]struct {
		name string
		on   bool
	}{
		{"dumb", p.Dumb},
		{"high-contrast", p.HighContrast},
		{"reduced-motion", p.ReducedMotion},
		{"screen-reader", p.ScreenReader},
	} {
		if flag.on {
			b.WriteByte(' ')
			b.WriteString(flag.name)
		}
	}
	return b.String()
}

// Describe devolve o perfil em uma linha legível, para `--version` ou debug.
// O resultado é prefixado com "a11y: " e é estável entre execuções.
func (p Profile) Describe() string {
	return "a11y: " + p.String()
}

// onOff converte um booleano em "on" ou "off".
func onOff(v bool) string {
	if v {
		return "on"
	}
	return "off"
}

// Symbol escolhe entre a versão Unicode e a ASCII de um símbolo.
//
// Quando o perfil não suporta Unicode, ascii é devolvido; caso contrário,
// unicode. Argumentos vazios são devolvidos como estão.
func (p Profile) Symbol(unicode, ascii string) string {
	if p.Unicode {
		return unicode
	}
	return ascii
}

// Frames devolve os frames de animação compatíveis com o perfil.
//
// Com ReducedMotion, apenas o primeiro frame é devolvido. Um slice nil
// permanece nil e um slice vazio permanece vazio; nenhum dos dois entra em
// panic.
func (p Profile) Frames(frames []string) []string {
	if frames == nil {
		return nil
	}
	if len(frames) == 0 {
		return frames
	}
	if p.ReducedMotion {
		return frames[:1]
	}
	return frames
}

// Announce formata uma mensagem para leitores de tela.
//
// A saída é texto linear: sequências ANSI são removidas, quebras de linha,
// tabs e demais caracteres de controle viram espaço, e espaços repetidos são
// reduzidos. O conteúdo Unicode legítimo é preservado.
//
// Quando não há argumentos, format é usado literalmente, sem interpretar '%'.
func (p Profile) Announce(format string, args ...any) string {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	msg = StripANSI(msg)
	msg = strings.Map(func(r rune) rune {
		switch {
		case r == '\n', r == '\r', r == '\t', r == '\v', r == '\f':
			return ' '
		case r < 0x20, r == 0x7f:
			return ' '
		default:
			return r
		}
	}, msg)
	return strings.Join(strings.Fields(msg), " ")
}
