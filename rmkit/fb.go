package rmkit

/*
#include <sys/ioctl.h>
#include <linux/fb.h>

struct fb_fix_screeninfo getFixScreenInfo(int fd) {
	struct fb_fix_screeninfo info;
	ioctl(fd, FBIOGET_FSCREENINFO, &info);
	return info;
}

struct fb_var_screeninfo getVarScreenInfo(int fd) {
	struct fb_var_screeninfo info;
	ioctl(fd, FBIOGET_VSCREENINFO, &info);
	return info;
}
*/
import "C"
import (
	"errors"
	"image"
	"image/color"
	"os"
	"sync"
	"syscall"
)

// Open expects a framebuffer device as its argument (such as "/dev/fb0"). The
// device will be memory-mapped to a local buffer. Writing to the device changes
// the screen output.
// The returned Device implements the draw.Image interface. This means that you
// can use it to copy to and from other images.
// The only supported color model for the specified frame buffer is RGB565.
// After you are done using the Device, call Close on it to unmap the memory and
// close the framebuffer file.
func Open(device string) (*Device, error) {
	file, err := os.OpenFile(device, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	fixInfo := C.getFixScreenInfo(C.int(file.Fd()))
	varInfo := C.getVarScreenInfo(C.int(file.Fd()))

	pixels, err := syscall.Mmap(
		int(file.Fd()),
		0, int(fixInfo.smem_len),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED,
	)
	if err != nil {
		file.Close()
		return nil, err
	}

	var colorModel color.Model
	if varInfo.red.offset == 11 && varInfo.red.length == 5 && varInfo.red.msb_right == 0 &&
		varInfo.green.offset == 5 && varInfo.green.length == 6 && varInfo.green.msb_right == 0 &&
		varInfo.blue.offset == 0 && varInfo.blue.length == 5 && varInfo.blue.msb_right == 0 {
		colorModel = RGB565(0)
	} else {
		return nil, errors.New("unsupported color model")
	}

	return &Device{
		file,
		pixels,
		int(fixInfo.line_length),
		image.Rect(0, 0, int(varInfo.xres), int(varInfo.yres)),
		colorModel,
		image.Rect(0, 0, 0, 0),
		sync.RWMutex{},
	}, nil
}

// Device represents the frame buffer. It implements the draw.Image interface.
type Device struct {
	File        *os.File
	Pixels      []byte
	Pitch       int
	bounds      image.Rectangle
	colorModel  color.Model
	DirtyBounds image.Rectangle

	Lock sync.RWMutex
}

// Close unmaps the framebuffer memory and closes the device file. Call this
// function when you are done using the frame buffer.
func (d *Device) Close() {
	syscall.Munmap(d.Pixels)
	d.File.Close()
}

// Bounds implements the image.Image (and draw.Image) interface.
func (d *Device) Bounds() image.Rectangle {
	return d.bounds
}

// ColorModel implements the image.Image (and draw.Image) interface.
func (d *Device) ColorModel() color.Model {
	return d.colorModel
}

// At implements the image.Image (and draw.Image) interface.
func (d *Device) At(x, y int) color.Color {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	if x < d.bounds.Min.X || x >= d.bounds.Max.X ||
		y < d.bounds.Min.Y || y >= d.bounds.Max.Y {
		return RGB565(0)
	}
	i := y*d.Pitch + 2*x
	return RGB565(d.Pixels[i+1])<<8 | RGB565(d.Pixels[i])
}

// Set implements the draw.Image interface.
func (d *Device) Set(x, y int, c color.Color) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	// the min bounds are at 0,0 (see Open)
	if x >= 0 && x < d.bounds.Max.X &&
		y >= 0 && y < d.bounds.Max.Y {
		var rgb RGB565
		var ok bool
		if rgb, ok = c.(RGB565); !ok {
			r, g, b, a := c.RGBA()
			if a <= 0 {
				return
			}
			rgb = ToRGB565(r, g, b)
		}
		// if rgb != d.At(x, y).(RGB565) {
		// 	d.DirtyBounds = d.DirtyBounds.Union(image.Rect(x, y, x+1, y+1))
		i := y*d.Pitch + 2*x
		// This assumes a little endian system which is the default for
		// Raspbian. The d.Pixels indices have to be swapped if the target
		// system is big endian.
		d.Pixels[i+1] = byte(rgb >> 8)
		d.Pixels[i] = byte(rgb & 0xFF)
		// }
	}
}

// Dirty returns the rectangle that has been modified since the last call to Set.
func (d *Device) Dirty() image.Rectangle {
	return d.DirtyBounds
}

// ResetDirty resets the DirtyBounds rectangle to an empty rectangle.
func (d *Device) ResetDirty() {
	d.DirtyBounds = image.Rect(0, 0, 0, 0)
}

// ToRGB565 helps convert a color.Color to rgb565. In a color.Color each
// channel is represented by the lower 16 bits in a uint32 so the maximum value
// is 0xFFFF. This function simply uses the highest 5 or 6 bits of each channel
// as the RGB values.
func ToRGB565(r, g, b uint32) RGB565 {
	// RRRRRGGGGGGBBBBB
	return RGB565((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}

// RGB565 implements the color.Color and color.Model interfaces.
// The default color model under the Raspberry Pi is RGB 565. Each pixel is
// represented by two bytes, with 5 bits for red, 6 bits for green and 5 bits
// for blue. There is no alpha channel, so alpha is assumed to always be 100%
// opaque.
// This shows the memory layout of a pixel:
//
//	bit 76543210  76543210
//	    RRRRRGGG  GGGBBBBB
//	   high byte  low byte
type RGB565 uint16

func (RGB565) Convert(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return ToRGB565(r, g, b)
}

// RGBA implements the color.Color interface.
func (c RGB565) RGBA() (r, g, b, a uint32) {
	// To convert a color channel from 5 or 6 bits back to 16 bits, the short
	// bit pattern is duplicated to fill all 16 bits.
	// For example the green channel in RGB565 is the middle 6 bits:
	//     00000GGGGGG00000
	//
	// To create a 16 bit channel, these bits are or-ed together starting at the
	// highest bit:
	//     GGGGGG0000000000 shifted << 5
	//     000000GGGGGG0000 shifted >> 1
	//     000000000000GGGG shifted >> 7
	//
	// These patterns map the minimum (all bits 0) and maximum (all bits 1)
	// 5 and 6 bit channel values to the minimum and maximum 16 bit channel
	// values.
	//
	// Alpha is always 100% opaque since this model does not support
	// transparency.
	rBits := uint32(c & 0xF800) // RRRRR00000000000
	gBits := uint32(c & 0x7E0)  // 00000GGGGGG00000
	bBits := uint32(c & 0x1F)   // 00000000000BBBBB
	r = uint32(rBits | rBits>>5 | rBits>>10 | rBits>>15)
	g = uint32(gBits<<5 | gBits>>1 | gBits>>7)
	b = uint32(bBits<<11 | bBits<<6 | bBits<<1 | bBits>>4)
	a = 0xFFFF
	return
}
