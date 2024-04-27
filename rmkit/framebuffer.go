package rmkit

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/chabad360/framebuffer"
)

var umCount uint32 = 10000

func _um() uint32 {
	umCount++
	return umCount
}

func InitRM() (*framebuffer.Device, error) {
	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		return nil, err
	}

	// updateMode := AUTO_UPDATE_MODE_AUTOMATIC_MODE

	req := IOW('F', MXCFB_SET_AUTO_UPDATE_MODE, unsafe.Sizeof(AUTO_UPDATE_MODE_AUTOMATIC_MODE))
	if err := unix.IoctlSetPointerInt(int(fb.File.Fd()), uint(req), int(AUTO_UPDATE_MODE_AUTOMATIC_MODE)); err != nil {
		return nil, err
	}

	return fb, nil
}

func Redraw(fb *framebuffer.Device, fullScreen bool) (uint32, error) {
	if fb.Dirty() == image.Rect(0, 0, 0, 0) && !fullScreen {
		return 0, errors.New("nothing to redraw")
	}

	var updateRect MxcfbRect
	if fullScreen {
		updateRect = MxcfbRect{
			Top:    0,
			Left:   0,
			Width:  uint32(fb.Bounds().Max.X),
			Height: uint32(fb.Bounds().Max.Y),
		}
	} else {
		updateRect = MxcfbRect{
			Top:    uint32(fb.Dirty().Min.Y),
			Left:   uint32(fb.Dirty().Min.X),
			Width:  uint32(fb.Dirty().Dx()),
			Height: uint32(fb.Dirty().Dy()),
		}
		fb.ResetDirty()
	}

	fmt.Println(updateRect)
	fb.ResetDirty()

	updateData := MxcfbUpdateData{
		UpdateRegion: updateRect,
		WaveformMode: WaveformModeInit,
		UpdateMode:   UPDATE_MODE_FULL,
		DitherMode:   EDPC_FLAG_EXP1,
		Temp:         TEMP_USE_REMARKABLE_DRAW,
		UpdateMarker: _um(),
	}

	if !fullScreen {
		updateData.WaveformMode = WaveformModeGC16
		updateData.UpdateMode = UPDATE_MODE_PARTIAL
	}

	req := IOW('F', MXCFB_SEND_UPDATE, unsafe.Sizeof(updateData))
	_, _, err := ioctl(fb, req, &updateData)
	if err != nil {
		return 0, err
	}

	return updateData.UpdateMarker, nil
}

func WaitForRedraw(fb *framebuffer.Device, marker uint32) error {
	if marker == 0 {
		return errors.New("invalid marker")
	}

	ud := MxcfbUpdateMarkerData{
		UpdateMarker: marker,
	}
	req := IOWR('F', MXCFB_WAIT_FOR_UPDATE_COMPLETE, unsafe.Sizeof(ud))
	if _, _, err := ioctl(fb, req, &ud); err != nil {
		return err
	}

	return nil
}

func ClearScreen(fb *framebuffer.Device) (uint32, error) {
	draw.Draw(fb, fb.Bounds(), image.NewUniform(White), image.Point{}, draw.Src)
	return Redraw(fb, true)
}

func ioctl[T any](fb *framebuffer.Device, cmd uintptr, arg *T) (uintptr, uintptr, error) {
	r1, r2, errno := syscall.Syscall(syscall.SYS_IOCTL, fb.File.Fd(), cmd, uintptr(unsafe.Pointer(arg)))
	if errno != 0 {
		return r1, r2, &os.SyscallError{"SYS_IOCTL", errno}
	}
	return r1, r2, nil
}
