package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

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
		item1Button.SetText("Press: " + strconv.Itoa(rand.Int()%1000))
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

	prevImage := image.NewRGBA(fb.Bounds())
	draw.Draw(prevImage, prevImage.Bounds(), image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}), image.Point{}, draw.Src)
	a := app.NewWithSoftwareDriver("test", func(img image.Image, rects []image.Rectangle) {
		fmt.Println("Rendering canvas")
		t := time.Now()
		var rect image.Rectangle
		if len(rects) == 0 {
			rect = img.Bounds()
		} else {
			rect = rects[0]
			for _, r := range rects {
				rect = rect.Union(r)
			}
		}
		fmt.Println("Rect:", rect)

		imgR := img.(*image.NRGBA)
		// if imgR.Bounds().Dx() != fb.Bounds().Dx() || imgR.Bounds().Dy() != fb.Bounds().Dy() || len(imgR.Pix)/4 < len(fb.Pixels)/2 {
		// 	fmt.Println("x:", imgR.Bounds().Dx(), fb.Bounds().Dx())
		// 	fmt.Println("y:", imgR.Bounds().Dy(), fb.Bounds().Dy())
		// 	fmt.Println("arr:", len(imgR.Pix)/4, len(fb.Pixels)/2)
		// 	panic("image size does not match framebuffer size")
		// }
		// draw.Draw(fb, fb.Bounds(), img, image.Point{}, draw.Src)
		// Dx := fb.Bounds().Dx()
		// dirty := false
		const ma = 1<<16 - 1
		for i := 0; i < len(imgR.Pix)/4; i++ {
			pixArr := imgR.Pix[i*4 : i*4+4 : i*4+4]
			prevArr := prevImage.Pix[i*4 : i*4+4 : i*4+4]
			a := uint32(pixArr[3]) * 0x101
			if a == 0 {
				continue
			}

			r := uint32(pixArr[0]) * a / 0xff
			g := uint32(pixArr[1]) * a / 0xff
			b := uint32(pixArr[2]) * a / 0xff
			am := (ma - a) * 0x101
			dr := uint32(prevArr[0]) | uint32(prevArr[0])<<8
			dg := uint32(prevArr[1]) | uint32(prevArr[1])<<8
			db := uint32(prevArr[2]) | uint32(prevArr[2])<<8

			rgb := rmkit.ToRGB565(
				dr*am/ma+r,
				dg*am/ma+g,
				db*am/ma+b,
			)

			prevArr[0] = uint8((dr*am/ma + r) >> 8)
			prevArr[1] = uint8((dg*am/ma + g) >> 8)
			prevArr[2] = uint8((db*am/ma + b) >> 8)

			// 	fb.DirtyBounds = fb.DirtyBounds.Union(image.Rect(x, y, x+1, y+1))

			pixelOffset := (i * 2) + ((i * 2 / (fb.Pitch - 8)) * 8)
			p := fb.Pixels[pixelOffset : pixelOffset+2 : pixelOffset+2]
			// if p[0] != byte(rgb&0xFF) || p[1] != byte(rgb>>8) {
			// 	x := i % Dx
			// 	y := i / Dx
			// 	if !dirty {
			// 		dirty = true
			// 		fb.DirtyBounds = image.Rect(x, y, x+1, y+1)
			// 	}
			// 	fb.DirtyBounds = fb.DirtyBounds.Union(image.Rect(x, y, x+1, y+1))
			p[0] = byte(rgb & 0xFF)
			p[1] = byte(rgb >> 8)
			// }
		}
		fb.DirtyBounds = fb.Bounds()
		m, err = rmkit.Redraw(fb, false)
		if err != nil {
			if err.Error() == "nothing to redraw" {
				return
			}
			panic(err)
		}

		fmt.Println("Redrew screen, ", m, time.Since(t))
		t = time.Now()
		err = rmkit.WaitForRedraw(fb, m)
		if err != nil {
			panic(err)
		}
		fmt.Println("Hardware took", time.Since(t))
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
