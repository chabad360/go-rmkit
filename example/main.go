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
	"go-rmkit/rmkit"

	"github.com/kenshaw/evdev"
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

	eC := make(chan any)

	prevImage := image.NewRGBA(fb.Bounds())
	draw.Draw(prevImage, prevImage.Bounds(), image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}), image.Point{}, draw.Src)
	a := app.NewWithSoftwareDriver("test", func(img image.Image) {
		fmt.Println("Rendering canvas")
		t := time.Now()
		imgR := img.(*image.NRGBA)
		// if imgR.Bounds().Dx() != fb.Bounds().Dx() || imgR.Bounds().Dy() != fb.Bounds().Dy() || len(imgR.Pix)/4 < len(fb.Pixels)/2 {
		// 	fmt.Println("x:", imgR.Bounds().Dx(), fb.Bounds().Dx())
		// 	fmt.Println("y:", imgR.Bounds().Dy(), fb.Bounds().Dy())
		// 	fmt.Println("arr:", len(imgR.Pix)/4, len(fb.Pixels)/2)
		// 	panic("image size does not match framebuffer size")
		// }
		// draw.Draw(fb, fb.Bounds(), img, image.Point{}, draw.Src)
		Dx := fb.Bounds().Dx()
		dirty := false
		minX, minY, maxX, maxY := 0, 0, 0, 0
		const ma = 1<<16 - 1
		for i := 0; i < len(imgR.Pix)/4; i++ {
			pixArr := imgR.Pix[i*4 : i*4+4 : i*4+4]
			prevArr := prevImage.Pix[i*4 : i*4+3 : i*4+3]
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

			dr = dr*am/ma + r
			dg = dg*am/ma + g
			db = db*am/ma + b

			prevArr[0] = uint8(dr >> 8)
			prevArr[1] = uint8(dg >> 8)
			prevArr[2] = uint8(db >> 8)

			rgb := rmkit.ToRGB565(dr, dg, db)

			// 	fb.DirtyBounds = fb.DirtyBounds.Union(image.Rect(x, y, x+1, y+1))

			pixelOffset := (i * 2) + ((i * 2 / (fb.Pitch - 8)) * 8)
			p := fb.Pixels[pixelOffset : pixelOffset+2 : pixelOffset+2]
			// if p[0] != byte(rgb&0xFF) || p[1] != byte(rgb>>8) {
			x := i % Dx
			y := i / Dx
			if !dirty {
				dirty = true
				minX, minY, maxX, maxY = x, y, x, y
			}
			// fb.DirtyBounds = fb.DirtyBounds.Union(image.Rect(x, y, x+1, y+1))
			minX, minY, maxX, maxY = min(minX, x), min(minY, y), max(maxX, x), max(maxY, y)
			p[0] = byte(rgb & 0xFF)
			p[1] = byte(rgb >> 8)
			// }
		}
		// fb.DirtyBounds = fb.Bounds()
		fb.DirtyBounds = image.Rect(minX, minY, maxX+1, maxY+1)
		fmt.Println(fb.Dirty())
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
	}, eC, true)

	a.Settings().SetTheme(t)

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
		begun := false
		for {
			select {
			case e, ok := <-touchC:
				if !ok {
					panic("touch closed")
				}
				if e == nil {
					panic("nil touch event")
				}
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
						begun = false
					} else if touchEv.Type == touch2.TypeEnd {
						touchEv.Type = touch2.TypeBegin
					}
				} else if ev.Type == evdev.EventSync {
					if touchEv.Type == touch2.TypeBegin && !begun {
						begun = true
					} else if begun {
						touchEv.Type = touch2.TypeMove
					}
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
