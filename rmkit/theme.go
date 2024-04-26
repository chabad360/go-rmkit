package rmkit

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = (*RMTheme)(nil)

type RMTheme struct{}

func (t *RMTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return White
	case theme.ColorNameForeground:
		return Black
	case theme.ColorNameButton:
		return White
	case theme.ColorNameInputBorder:
		return Black
	case theme.ColorNameShadow:
		return Black
	default:
		d := theme.DefaultTheme().Color(name, variant)
		fmt.Printf("Color %s, %d: %v\n", name, variant, d)
		return d
	}
}

func (t *RMTheme) Font(style fyne.TextStyle) fyne.Resource {
	switch style {
	default:
		d := theme.DefaultTheme().Font(style)
		fmt.Printf("Font %#v: %s\n", style, d.Name())
		return d
	}
}

func (t *RMTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	switch name {
	default:
		d := theme.DefaultTheme().Icon(name)
		fmt.Printf("Icon %s: %v\n", name, d.Name())
		return d
	}
}

func (t *RMTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 36
	case theme.SizeNameInlineIcon:
		return 50
	case theme.SizeNameLineSpacing:
		return 0
	case theme.SizeNameInnerPadding:
		return 5
	case theme.SizeNamePadding:
		return 20
	case theme.SizeNameInputRadius:
		return 0
	default:
		d := theme.DefaultTheme().Size(name)
		fmt.Printf("Size %s: %v\n", name, d)
		return d
	}
}
