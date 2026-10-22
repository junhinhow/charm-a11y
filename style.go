package a11y

import "charm.land/lipgloss/v2"

// Style degrada um lipgloss.Style conforme o perfil.
//
// Com NoColor, remove cores de texto, fundo, borda e margem. Com
// HighContrast, remove o atributo faint e aplica negrito, sem introduzir
// nenhuma cor. Quando o perfil não suporta Unicode, a borda do estilo é
// trocada pela versão ASCII de [Profile.Border], preservando os lados
// habilitados.
//
// A função é pura: o estilo recebido não é modificado.
func (p Profile) Style(s lipgloss.Style) lipgloss.Style {
	if p.NoColor {
		s = s.UnsetForeground().UnsetBackground()
		s = s.UnsetBorderForeground().UnsetBorderForegroundBlend()
		s = s.UnsetBorderBackground()
		s = s.UnsetMarginBackground()
		if s.GetUnderlineColor() != (lipgloss.NoColor{}) {
			s = s.UnderlineColor(lipgloss.NoColor{})
		}
	}
	if p.HighContrast {
		s = s.UnsetFaint().Bold(true)
	}
	if !p.Unicode {
		if b := s.GetBorderStyle(); b != (lipgloss.Border{}) {
			s = s.BorderStyle(p.Border(b))
		}
	}
	return s
}

// Border devolve uma borda compatível com o perfil.
//
// Com Unicode habilitado, a borda é devolvida sem alterações. Caso contrário,
// cada parte é transliterada para ASCII usando a mesma tabela de
// [Profile.Sanitize]. Bordas personalizadas com glifos desconhecidos mantêm
// esses glifos.
func (p Profile) Border(b lipgloss.Border) lipgloss.Border {
	if p.Unicode {
		return b
	}
	return lipgloss.Border{
		Top:          transliterate(b.Top),
		Bottom:       transliterate(b.Bottom),
		Left:         transliterate(b.Left),
		Right:        transliterate(b.Right),
		TopLeft:      transliterate(b.TopLeft),
		TopRight:     transliterate(b.TopRight),
		BottomLeft:   transliterate(b.BottomLeft),
		BottomRight:  transliterate(b.BottomRight),
		MiddleLeft:   transliterate(b.MiddleLeft),
		MiddleRight:  transliterate(b.MiddleRight),
		Middle:       transliterate(b.Middle),
		MiddleTop:    transliterate(b.MiddleTop),
		MiddleBottom: transliterate(b.MiddleBottom),
	}
}
