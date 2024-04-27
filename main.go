package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	touch2 "golang.org/x/mobile/event/touch"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kenshaw/evdev"

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
	item1Button.OnTapped = func() {
		item1Button.SetText("Tutorial Started")
	}

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
	item2Button.OnTapped = func() {
		item2Button.SetText("Tutorial Started")
	}

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
	//
	// c := software.NewCanvas()
	// c.SetContent(container.NewBorder(s, s, s, s, con))
	// c.Resize(fyne.NewSize(float32(fb.Bounds().Dx()), float32(fb.Bounds().Dy())))
	//
	// i := software.RenderCanvas(c, t)
	//
	// fmt.Println("Rendered canvas")
	//
	// draw.Draw(fb, fb.Bounds(), i, image.Point{}, draw.Src)
	// m, err = rmkit.Redraw(fb, false)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println("Redrew screen, ", m)
	//
	// err = rmkit.WaitForRedraw(fb, m)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println("Waited for redraw")

	eC := make(chan any)

	a := app.NewWithSoftwareDriver("test", func(i image.Image) {
		draw.Draw(fb, fb.Bounds(), i, image.Point{}, draw.Src)
		m, err = rmkit.Redraw(fb, false)
		if err != nil {
			if err.Error() == "nothing to redraw" {
				return
			}
			panic(err)
		}

		fmt.Println("Redrew screen, ", m)
		err = rmkit.WaitForRedraw(fb, m)
		if err != nil {
			panic(err)
		}
	}, eC)

	a.Settings().SetTheme(t)

	w := a.NewWindow("test")
	w.SetContent(container.NewBorder(s, s, s, s, con))
	w.Resize(fyne.NewSize(float32(fb.Bounds().Dx()), float32(fb.Bounds().Dy())))
	w.SetFullScreen(true)
	w.SetFixedSize(true)

	touchF, err := os.Open("/dev/input/event2")
	if err != nil {
		panic(err)
	}
	defer touchF.Close()

	touch := evdev.OpenReader(touchF)
	if touch == nil {
		panic("could not open touch")
	}
	defer touch.Close()

	touchC := touch.Poll(context.Background())

	go func() {
		touchEv := touch2.Event{}
		for {
			select {
			case e := <-touchC:
				ev := e.Event
				if ev.Code == uint16(evdev.AbsoluteMTPositionX) {
					// ev.Code = uint16(evdev.AbsoluteMTPositionY)
					ev.Value = 767 - ev.Value
					// ev.Value = ev.Value * 10
					touchEv.X = xRatio * float32(ev.Value)
				} else if ev.Code == uint16(evdev.AbsoluteMTPositionY) {
					// ev.Code = uint16(evdev.AbsoluteMTPositionX)
					// ev.Value = ev.Value * 10
					ev.Value = 1023 - ev.Value
					touchEv.Y = yRatio * float32(ev.Value)
				} else if ev.Code == uint16(evdev.AbsoluteMTTrackingID) {
					if ev.Value == -1 {
						touchEv.Type = touch2.TypeEnd
					} else {
						if touchEv.Type == touch2.TypeEnd {
							touchEv.Type = touch2.TypeBegin
						} else {
							touchEv.Type = touch2.TypeMove
						}
					}
				} else if ev.Type == evdev.EventSync {
					eC <- touchEv
				}
			}
		}
	}()

	w.ShowAndRun()

	fmt.Println("Exiting...")
}

const xRatio = float32(1404) / float32(767)
const yRatio = float32(1872) / float32(1023)

func frame(c fyne.CanvasObject) fyne.CanvasObject {
	s1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	s1.StrokeColor = rmkit.Black
	s1.StrokeWidth = 3
	s1.CornerRadius = -1

	s2 := canvas.NewRectangle(rmkit.White)

	return container.NewStack(s2, container.NewPadded(c), s1)
}
