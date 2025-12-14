package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ----------------------
// Embedded font resources
// ----------------------
var FontJP fyne.Resource // Japanese (includes Latin characters)
var FontKR fyne.Resource // Korean (includes Latin characters)
var FontSC fyne.Resource // Simplified Chinese (includes Latin characters)

// ----------------------
// Custom CJK theme
// ----------------------
type CustomCJKTheme struct {
	Base fyne.Theme
}

func (t *CustomCJKTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.Base.Color(name, variant)
}

func (t *CustomCJKTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.Base.Size(name)
}

func (t *CustomCJKTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.Base.Icon(name)
}

// Font returns the CJK font for all text
// Noto Sans CJK fonts include Latin/English characters, so this works for both
func (t *CustomCJKTheme) Font(style fyne.TextStyle) fyne.Resource {
	// For monospace, use default
	if style.Monospace {
		return t.Base.Font(style)
	}

	// For bold
	if style.Bold {
		// You can load bold variants if you have them
		// For now, fall back to regular (bold rendering will be synthesized)
		return FontJP
	}

	// For italic
	if style.Italic {
		// CJK fonts typically don't have italic variants
		// Fyne will synthesize italic rendering
		return FontJP
	}

	// Default: return Japanese font (includes Latin + Japanese)
	// This will work for English text too!
	return FontJP
}

// ----------------------
// Initialize theme globally
// ----------------------
func InitCJKTheme(app fyne.App, jp, kr, sc fyne.Resource) {
	FontJP = jp
	FontKR = kr
	FontSC = sc

	app.Settings().SetTheme(&CustomCJKTheme{
		Base: theme.DefaultTheme(),
	})
}
