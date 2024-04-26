package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"rm-cal/rmkit"
)

func main() {
	fmt.Println("Starting...")
	fb, err := rmkit.InitRM()
	if err != nil {
		panic(err)
	}
	defer fb.Close()

	fmt.Println("Initialized framebuffer")

	m, err := rmkit.ClearScreen(fb)
	if err != nil {
		panic(err)
	}

	fmt.Println("Cleared screen, ", m)

	err = rmkit.WaitForRedraw(fb, m)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waited for redraw")

	// time.Sleep(2 * time.Second)

	t := &rmkit.RMTheme{}

	backButton := widget.NewButton("Back", func() {})
	backButton.SetIcon(t.Icon(theme.IconNameNavigateBack))
	bbCon := container.NewHBox(backButton)

	item1Button := widget.NewButton("Start Tutorial", func() {})
	item1Button.SetIcon(t.Icon(theme.IconNameNavigateNext))
	item1Button.IconPlacement = widget.ButtonIconTrailingText

	item1Image := canvas.NewImageFromResource(theme.FyneLogo())
	item1Image.FillMode = canvas.ImageFillContain
	item1Image.SetMinSize(fyne.NewSize(200, 200))

	item1Header := widget.NewLabel("Fyne Toolkit")
	item1Header.TextStyle = fyne.TextStyle{Bold: true}

	item1Label := widget.NewLabel("An easy-to-use UI toolkit \nand app API written in Go.")

	item1 := container.NewPadded(frame(container.NewPadded(container.NewVBox(
		container.NewBorder(nil, nil, item1Image, nil),
		item1Header,
		item1Label,
		container.NewBorder(nil, nil, nil, container.NewPadded(frame(item1Button))),
	))))

	item2Button := widget.NewButton("Start Tutorial", func() {})
	item2Button.SetIcon(t.Icon(theme.IconNameNavigateNext))
	item2Button.IconPlacement = widget.ButtonIconTrailingText

	item2Image := canvas.NewImageFromResource(theme.FyneLogo())
	item2Image.FillMode = canvas.ImageFillContain
	item2Image.SetMinSize(fyne.NewSize(200, 200))

	item2Header := widget.NewLabel("Fyne Toolkit")
	item2Header.TextStyle = fyne.TextStyle{Bold: true}

	item2Label := widget.NewLabel("An easy-to-use UI toolkit \nand app API written in Go.")

	item2 := container.NewPadded(frame(container.NewPadded(container.NewVBox(
		container.NewBorder(nil, nil, item2Image, nil),
		item2Header,
		item2Label,
		container.NewBorder(nil, nil, nil, container.NewPadded(frame(item2Button))),
	))))

	contentGrid := container.NewGridWithColumns(2, item1, item2)

	label := widget.NewLabel("Guides")
	label.TextStyle = fyne.TextStyle{Bold: true}

	con := container.NewBorder(bbCon, nil, nil, nil,
		container.NewPadded(container.NewVBox(label, contentGrid)))

	s := canvas.NewRectangle(rmkit.White)
	s.SetMinSize(fyne.NewSize(10, 10))

	c := software.NewCanvas()
	c.SetContent(container.NewBorder(s, s, s, s, con))
	c.Resize(fyne.NewSize(float32(fb.Bounds().Dx()), float32(fb.Bounds().Dy())))

	i := software.RenderCanvas(c, t)

	fmt.Println("Rendered canvas")

	draw.Draw(fb, fb.Bounds(), i, image.Point{}, draw.Src)
	m, err = rmkit.Redraw(fb, false)
	if err != nil {
		panic(err)
	}

	fmt.Println("Redrew screen, ", m)

	err = rmkit.WaitForRedraw(fb, m)
	if err != nil {
		panic(err)
	}

	fmt.Println("Waited for redraw")

	fmt.Println("Exiting...")
}

func frame(c fyne.CanvasObject) fyne.CanvasObject {
	s1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	s1.StrokeColor = rmkit.Black
	s1.StrokeWidth = 3
	s1.CornerRadius = -1

	s2 := canvas.NewRectangle(rmkit.White)

	return container.NewStack(s2, container.NewPadded(c), s1)
}
