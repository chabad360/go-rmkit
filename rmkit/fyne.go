package rmkit

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	"github.com/kenshaw/evdev"
	touch2 "golang.org/x/mobile/event/touch"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func InitApp(id string) (fyne.App, *Device, error) {
	fmt.Println("Starting...")
	fb, err := InitRM()
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Initialized framebuffer")

	m, err := ClearScreen(fb)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Cleared screen, ", m)

	err = WaitForRedraw(fb, m)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Waited for redraw")

	eC := make(chan any)

	prevImage := image.NewRGBA(fb.Bounds())
	draw.Draw(prevImage, prevImage.Bounds(), image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: 255}), image.Point{}, draw.Src)
	a := app.NewWithSoftwareDriver("test", func(img image.Image) {
		fb.Lock.Lock()
		defer fb.Lock.Unlock()
		fmt.Println("Rendering canvas")
		t := time.Now()
		imgR := img.(*image.NRGBA)

		Dx := fb.Bounds().Dx()
		dirty := false
		minX, minY, maxX, maxY := 0, 0, 0, 0
		const ma = 1<<16 - 1
		for i := 0; i < len(imgR.Pix)/4; i++ {
			pixArr := imgR.Pix[i*4 : i*4+4 : i*4+4]
			prevArr := prevImage.Pix[i*4 : i*4+4 : i*4+4]
			a := uint32(pixArr[3]) * 0x101
			if a == 0 {
				continue
			}
			r := (uint32(pixArr[0]) * 0x101) * a / 0xff
			g := (uint32(pixArr[1]) * 0x101) * a / 0xff
			b := (uint32(pixArr[2]) * 0x101) * a / 0xff

			dr := uint32(prevArr[0]) * 0x101
			dg := uint32(prevArr[1]) * 0x101
			db := uint32(prevArr[2]) * 0x101
			da := uint32(prevArr[3]) * 0x101

			am := (ma - a) * 0x101

			dr = dr*am/ma + r
			dg = dg*am/ma + g
			db = db*am/ma + b
			da = da*am/ma + a

			prevArr[0] = uint8(dr >> 8)
			prevArr[1] = uint8(dg >> 8)
			prevArr[2] = uint8(db >> 8)
			prevArr[3] = uint8(da >> 8)

			rgb := ToRGB565(dr, dg, db)

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
		m, err = Redraw(fb, false)
		if err != nil {
			if err.Error() == "nothing to redraw" {
				return
			}
			panic(err)
		}

		fmt.Println("Redrew screen, ", m, time.Since(t))
		t = time.Now()
		err = WaitForRedraw(fb, m)
		if err != nil {
			panic(err)
		}
		fmt.Println("Hardware took", time.Since(t))
	}, eC, true)

	a.Settings().SetTheme(&RMTheme{})

	InitTouch(eC)

	return a, fb, nil
}

const xRatio = float32(1404) / float32(767)
const yRatio = float32(1872) / float32(1023)

func InitTouch(eC chan<- any) error {
	touchF, err := os.Open("/dev/input/event2")
	if err != nil {
		return err
	}
	// defer touchF.Close()

	touch := evdev.OpenReader(touchF)
	if touch == nil {
		return err
	}
	// defer touch.Close()

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

	return nil
}

func SetFyneWindowSize(fb *Device, w fyne.Window) {
	w.Resize(fyne.NewSize(float32(fb.Bounds().Dx()), float32(fb.Bounds().Dy())))
	w.SetFullScreen(true)
	w.SetFixedSize(true)
}
