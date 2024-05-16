package rmkit

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = (*RMTheme)(nil)

type RMTheme struct{}

var fyneWhite = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
var fyneBlack = color.NRGBA{A: 255}

func (t *RMTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return fyneWhite
	case theme.ColorNameForeground:
		return fyneBlack
	case theme.ColorNameButton:
		return fyneWhite
	case theme.ColorNameInputBorder:
		return fyneBlack
	case theme.ColorNameShadow:
		return fyneBlack
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

type RMThemeVariant struct {
	RMTheme
}

func (t *RMThemeVariant) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInnerPadding:
		return 20
	default:
		return t.RMTheme.Size(name)
	}
}

func Frame(c fyne.CanvasObject) fyne.CanvasObject {
	s1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	s1.StrokeColor = Black
	s1.StrokeWidth = 3
	s1.CornerRadius = -1

	s2 := canvas.NewRectangle(White)

	return container.NewStack(s2, container.NewThemeOverride(c, &RMThemeVariant{}), s1)
}
