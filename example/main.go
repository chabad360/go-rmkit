package main

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go-rmkit/rmkit"
)

func main() {

	// time.Sleep(2 * time.Second)

	t := &rmkit.RMTheme{}

	a, fb, err := rmkit.InitApp("test")
	if err != nil {
		panic(err)
	}

	backButton := widget.NewButton("Back", func() {})
	backButton.SetIcon(t.Icon(theme.IconNameNavigateBack))
	bbCon := container.NewHBox(backButton)

	square := canvas.NewRectangle(color.RGBA{R: 100, G: 100, B: 100, A: 255})
	square.SetMinSize(fyne.NewSize(50, 100))

	item1Button := widget.NewButton("Start Tutorial", func() {})
	item1Button.SetIcon(t.Icon(theme.IconNameNavigateNext))
	item1Button.IconPlacement = widget.ButtonIconTrailingText
	item1Button.OnTapped = func() {
		item1Button.SetText("Press: " + strconv.Itoa(rand.Int()%1000))
	}

	item1Image := canvas.NewImageFromResource(theme.FyneLogo())
	item1Image.FillMode = canvas.ImageFillContain
	item1Image.SetMinSize(fyne.NewSize(200, 200))

	item1Header := widget.NewLabel("Fyne Toolkit")
	item1Header.TextStyle = fyne.TextStyle{Bold: true}

	item1Label := widget.NewLabel("An easy-to-use UI toolkit \nand app API written in Go.")

	item1 := container.NewPadded(rmkit.Frame(container.NewThemeOverride(container.NewPadded(container.NewPadded(container.NewVBox(
		container.NewBorder(nil, nil, item1Image, nil),
		item1Header,
		item1Label,
		container.NewBorder(nil, nil, nil, container.NewPadded(rmkit.Frame(container.NewStack(item1Button, square)))),
	))), t)))

	item2Button := widget.NewButton("Exit", func() {})
	item2Button.SetIcon(t.Icon(theme.IconNameNavigateNext))
	item2Button.IconPlacement = widget.ButtonIconTrailingText
	item2Button.OnTapped = func() {
		panic("Exiting...")
	}

	item2Image := canvas.NewImageFromResource(theme.FyneLogo())
	item2Image.FillMode = canvas.ImageFillContain
	item2Image.SetMinSize(fyne.NewSize(200, 200))

	item2Header := widget.NewLabel("Fyne Toolkit")
	item2Header.TextStyle = fyne.TextStyle{Bold: true}

	item2Label := widget.NewLabel("An easy-to-use UI toolkit \nand app API written in Go.")

	item2 := container.NewPadded(rmkit.Frame(container.NewThemeOverride(container.NewPadded(container.NewPadded(container.NewVBox(
		container.NewBorder(nil, nil, item2Image, nil),
		item2Header,
		item2Label,
		container.NewBorder(nil, nil, nil, container.NewPadded(rmkit.Frame(item2Button))),
	))), t)))

	contentGrid := container.NewGridWithColumns(2, item1, item2)

	label := widget.NewLabel("Guides")
	label.TextStyle = fyne.TextStyle{Bold: true}

	con := container.NewBorder(bbCon, nil, nil, nil,
		container.NewPadded(container.NewVBox(label, contentGrid)))

	w := a.NewWindow("test")
	w.SetContent(container.NewPadded(con))
	rmkit.SetFyneWindowSize(fb, w)

	w.ShowAndRun()

	fmt.Println("Exiting...")
}
