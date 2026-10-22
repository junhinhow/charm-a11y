package a11y

import "github.com/charmbracelet/x/ansi"

// StripANSI remove sequências ANSI/OSC de s.
//
// São removidos códigos CSI (ESC '[' ...), OSC (ESC ']' ... BEL ou ESC '\'),
// sequências de dois bytes como ESC '(' 'B' e demais famílias de controle
// reconhecidas pelo parser de github.com/charmbracelet/x/ansi. Texto comum,
// incluindo colchetes literais, é preservado.
//
// Sequências malformadas ou sem terminador não causam panic; o parser consome
// o restante da entrada como parte da sequência incompleta. Strings vazias
// devolvem string vazia.
func StripANSI(s string) string {
	if s == "" {
		return ""
	}
	return ansi.Strip(s)
}
