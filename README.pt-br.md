# charm-a11y

Perfis de acessibilidade e degradação graciosa para aplicações de terminal Charm: respeite `NO_COLOR`, caia para ASCII, desligue animações e emita texto linear para leitores de tela.

## Instalação

```bash
go get github.com/junhinhow/charm-a11y
```

## Uso

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

Com `TERM=dumb`, o mesmo programa imprime uma borda ASCII, sem cor, e uma linha de status em texto puro. Com `A11Y_REDUCED_MOTION=1`, `p.Frames` devolve um único frame.

## API

| Símbolo | Descrição |
| --- | --- |
| `Detect() Profile` | Lê o ambiente do processo uma vez. |
| `DetectFrom(getenv func(string) string) Profile` | Mesma detecção com ambiente injetável; `nil` significa vazio. |
| `Default() Profile` | Perfil do processo, detectado no primeiro uso. |
| `Profile` | Capacidades: `NoColor`, `Dumb`, `HighContrast`, `ReducedMotion`, `ScreenReader`, `Unicode`, `Width`. |
| `StripANSI(s string) string` | Remove CSI, OSC, escapes de dois bytes e sequências de controle relacionadas. |
| `(Profile) Sanitize(s string) string` | Remove ANSI e, sem Unicode, translitera glifos para ASCII. |
| `(Profile) Symbol(unicode, ascii string) string` | Escolhe a variante Unicode ou ASCII. |
| `(Profile) Style(s lipgloss.Style) lipgloss.Style` | Remove cor, aumenta contraste e converte bordas para ASCII. |
| `(Profile) Border(b lipgloss.Border) lipgloss.Border` | Devolve uma borda segura em ASCII quando necessário. |
| `(Profile) Frames(frames []string) []string` | Devolve um frame sob movimento reduzido. |
| `(Profile) Announce(format string, args ...any) string` | Formata uma mensagem linear e sem ANSI para leitores de tela. |
| `(Profile) Describe() string` | Resumo estável de uma linha para `--version` ou debug. |
| `(Profile) String() string` | Forma compacta estável usada por `Describe`. |

## Variáveis de ambiente

| Variável | Valores aceitos | Efeito |
| --- | --- | --- |
| `NO_COLOR` | Qualquer valor não vazio | `NoColor`; também `Unicode=false` |
| `CLICOLOR` | `0` | `NoColor`; também `Unicode=false` |
| `TERM` | `dumb` (sem diferenciar caixa) | `NoColor`, `Dumb`, `ScreenReader`; `Unicode=false` |
| `ASCII` | `1` ou `true` | `Unicode=false` |
| `A11Y_ASCII` | `1` ou `true` | `Unicode=false` |
| `A11Y_HIGH_CONTRAST` | `1` ou `true` | `HighContrast=true` |
| `A11Y_REDUCED_MOTION` | `1` ou `true` | `ReducedMotion=true` |
| `A11Y_SCREEN_READER` | `1` ou `true` | `ScreenReader=true` |
| `COLUMNS` | Inteiro positivo | Preenche `Width`; caso contrário, `0` (desconhecido) |

Observações:

- Conforme a especificação [no-color.org](https://no-color.org), `NO_COLOR` só vale quando definida com um valor **não vazio**. `NO_COLOR=` não faz nada; `NO_COLOR=0` desliga cor porque é não vazio.
- Os valores de `A11Y_*` e `ASCII` têm espaços removidos e são comparados sem diferenciar maiúsculas de minúsculas.
- `Unicode` é deliberadamente desligado sempre que a cor é desligada, não apenas em `TERM=dumb`. Essa é a escolha conservadora para pipes, emuladores antigos e ambientes mínimos. Atualmente não há override que force Unicode com a cor desligada.

## Transliteração ASCII

`Sanitize` e `Border` usam a tabela abaixo. Glifos não listados são preservados, pois não há equivalente ASCII universal. Substituições com mais de um caractere são intencionais; a largura pode mudar.

### Desenho de caixas

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

### Blocos e sombreados

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

### Setas e triângulos

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

### Marcas, pontuação e símbolos

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

## Checklist de acessibilidade para apps de terminal

- Detecte o perfil antes da primeira renderização e passe-o pela árvore de UI.
- Nunca comunique estado apenas com cor; combine cor com texto ou um glifo distinto.
- Teste com `NO_COLOR=1`, `TERM=dumb`, `ASCII=1` e um terminal de 80 colunas.
- Ofereça caminho ASCII para cada ícone, spinner e borda personalizados.
- Respeite `A11Y_REDUCED_MOTION`: sem spinners, marquees ou redraws rápidos; mostre um status estático.
- Mantenha a saída linear e descritiva; anuncie mudanças de estado relevantes em vez de apenas repintar.
- Suporte navegação só por teclado e mantenha um indicador de foco visível.
- Evite texto piscando, flashes e rolagem automática.
- Evite cinza de baixo contraste sobre preto; teste com `A11Y_HIGH_CONTRAST=1`.
- Não exija mouse ou hover para acessar funcionalidades.
- Não redesenhe a tela inteira sem necessidade; leitores de tela lidam melhor com texto estável e acrescentado.
- Exponha erros e código de saída como texto puro.
- Documente as variáveis de ambiente que o app respeita.

## Notas de design

- **Somente ambiente.** A detecção nunca consulta o terminal. É determinística, funciona em CI e é trivial de testar com `DetectFrom`.
- **Unicode conservador.** `NoColor` implica `Unicode=false`. Isso segue a especificação do projeto e protege ambientes limitados, mas significa que um terminal colorido em um pipeline sem cor também perde o desenho de caixas. Glifos desconhecidos são preservados por `Sanitize`, então conteúdo nunca é descartado silenciosamente.
- **Sem adivinhar cores.** `HighContrast` remove `faint` e adiciona `bold`; nunca escolhe cores de texto ou fundo. Escolher uma paleta de alto contraste exigiria conhecer o fundo do terminal, que este pacote não detecta de forma confiável.
- **Remoção de ANSI** delega a `github.com/charmbracelet/x/ansi`, que cobre CSI, OSC (BEL e ST), DCS, APC, PM, SOS e escapes de dois bytes. Sequências malformadas ou sem terminador são consumidas até o fim da string, sem panic.
- **Announce não é um driver de leitor de tela.** Terminais expõem o texto renderizado à tecnologia assistiva; `Announce` reduz uma mensagem a texto estável, linear e sem ANSI, que essas ferramentas leem de forma mais previsível.
- **Largura é consultiva.** `COLUMNS` é a única fonte; `0` significa desconhecido. Nenhum ioctl é executado.
- **Frames não aloca.** Sob movimento reduzido, devolve um subslice da entrada; o chamador não deve contar com uma cópia.
- **Fora do escopo:** cálculo de contraste de cor, detecção automática do fundo do terminal, localização BCP-47 e normalização Unicode completa.

## Licença

MIT
