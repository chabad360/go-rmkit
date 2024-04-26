package rmkit

/*
#include <linux/fb.h>
*/
import "C"

type FB_FLAGS uint32

const (
	FB_SYNC_OE_LOW_ACT   FB_FLAGS = 0x80000000
	FB_SYNC_CLK_LAT_FALL FB_FLAGS = 0x40000000
	FB_SYNC_DATA_INVERT  FB_FLAGS = 0x20000000
	FB_SYNC_CLK_IDLE_EN  FB_FLAGS = 0x10000000
	FB_SYNC_SHARP_MODE   FB_FLAGS = 0x08000000
	FB_SYNC_SWAP_RGB     FB_FLAGS = 0x04000000
	FB_ACCEL_TRIPLE_FLAG FB_FLAGS = 0x00000000
	FB_ACCEL_DOUBLE_FLAG FB_FLAGS = 0x00000001
)

type MxcfbGblAlpha struct {
	Enable int
	Alpha  int
}

type MxcfbLocAlpha struct {
	Enable        int
	AlphaInPixel  int
	AlphaPhyAddr0 uint32
	AlphaPhyAddr1 uint32
}

type MxcfbColorKey struct {
	Enable   int
	ColorKey uint32
}

type MxcfbPos struct {
	X, Y uint16
}

type MxcfbGamma struct {
	Enable int
	ConstK [16]int
	SlopeK [16]int
}

type MxcfbGPUSplitFmt struct {
	Var    C.struct_fb_var_screeninfo
	Offset uint64
}

type MxcfbRect struct {
	Top, Left, Width, Height uint32
}

type GS_MODE uint32

const (
	GRAYSCALE_8BIT          GS_MODE = 0x1
	GRAYSCALE_8BIT_INVERTED GS_MODE = 0x2
	GRAYSCALE_4BIT          GS_MODE = 0x3
	GRAYSCALE_4BIT_INVERTED GS_MODE = 0x4
)

type AUTO_UPDATE_MODE int

const (
	AUTO_UPDATE_MODE_REGION_MODE    AUTO_UPDATE_MODE = 0
	AUTO_UPDATE_MODE_AUTOMATIC_MODE AUTO_UPDATE_MODE = 1
)

type UPDATE_SCHEME int

const (
	UPDATE_SCHEME_SNAPSHOT        UPDATE_SCHEME = 0
	UPDATE_SCHEME_QUEUE           UPDATE_SCHEME = 1
	UPDATE_SCHEME_QUEUE_AND_MERGE UPDATE_SCHEME = 2
)

type UPDATE_MODE uint32

const (
	UPDATE_MODE_PARTIAL UPDATE_MODE = 0x0
	UPDATE_MODE_FULL    UPDATE_MODE = 0x1
)

type WAVEFORM_MODE uint32

const (
	WAVEFORM_MODE_GLR16 WAVEFORM_MODE = 4
	WAVEFORM_MODE_GLD16 WAVEFORM_MODE = 5

	WAVEFORM_MODE_INIT   WAVEFORM_MODE = 0x0 /* Screen goes to white (clears) */
	WAVEFORM_MODE_DU     WAVEFORM_MODE = 0x1 /* Grey->white/grey->black */
	WAVEFORM_MODE_GC16   WAVEFORM_MODE = 0x2 /* High fidelity (flashing) */
	WAVEFORM_MODE_GC4    WAVEFORM_MODE = 0x3 /* Lower fidelity */
	WAVEFORM_MODE_A2     WAVEFORM_MODE = 0x4 /* Fast black/white animation */
	WAVEFORM_MODE_DU4    WAVEFORM_MODE = 0x7
	WAVEFORM_MODE_REAGLD WAVEFORM_MODE = 0x9

	WAVEFORM_MODE_AUTO WAVEFORM_MODE = 257
)

type TEMP int32

const TEMP_USE_AMBIENT TEMP = 0x1000
const TEMP_USE_REMARKABLE_DRAW TEMP = 0x0018

type EDPC_FLAGS int32

const (
	EDPC_FLAG_ENABLE_INVERSION EDPC_FLAGS = 0x01
	EDPC_FLAG_FORCE_MONOCHROME EDPC_FLAGS = 0x02
	EDPC_FLAG_USE_CMAP         EDPC_FLAGS = 0x04
	EDPC_FLAG_USE_ALT_BUFFER   EDPC_FLAGS = 0x100
	EDPC_FLAG_TEST_COLLISION   EDPC_FLAGS = 0x200
	EDPC_FLAG_GROUP_UPDATE     EDPC_FLAGS = 0x400
	EDPC_FLAG_USE_DITHERING_Y1 EDPC_FLAGS = 0x2000
	EDPC_FLAG_USE_DITHERING_Y4 EDPC_FLAGS = 0x4000
	EDPC_FLAG_USE_REGAL        EDPC_FLAGS = 0x8000
)

type EDPC_DITHERING_FLAGS EDPC_FLAGS

const (
	EDPC_FLAG_USE_DITHERING_PASSTHROUGH EDPC_DITHERING_FLAGS = iota
	EDPC_FLAG_USE_DITHERING_FLOYD_STEINBERG
	EDPC_FLAG_USE_DITHERING_ATKINSON
	EDPC_FLAG_USE_DITHERING_ORDERED
	EDPC_FLAG_USE_DITHERING_QUANT_ONLY
	EDPC_FLAG_USE_DITHERING_MAX

	EDPC_FLAG_EXP1 EDPC_DITHERING_FLAGS = 0x270ce20
)

const (
	FB_POWERDOWN_DISABLE        int = -1
	FB_TEMP_AUTO_UPDATE_DISABLE int = -1
)

type MxcfbAltBufferData struct {
	PhysAddr        uint32
	Width           uint32
	Height          uint32
	AltUpdateRegion MxcfbRect
}

type MxcfbUpdateData struct {
	UpdateRegion  MxcfbRect
	WaveformMode  WAVEFORM_MODE
	UpdateMode    UPDATE_MODE
	UpdateMarker  uint32
	Temp          TEMP
	Flags         EDPC_FLAGS
	DitherMode    EDPC_DITHERING_FLAGS
	QuantBit      int32
	AltBufferData MxcfbAltBufferData
}

type MxcfbUpdateMarkerData struct {
	UpdateMarker  uint32
	CollisionTest uint32
}

type MxcfbWaveformModes struct {
	ModeInit int
	ModeDu   int
	ModeGc4  int
	ModeGc8  int
	ModeGc16 int
	ModeGc32 int
}

type MxcfbCscMatrix struct {
	Param [5][3]int
}

//
// #define MXCFB_WAIT_FOR_VSYNC	_IOW('F', 0x20, u_int32_t)
// #define MXCFB_SET_GBL_ALPHA     _IOW('F', 0x21, struct mxcfb_gbl_alpha)
// #define MXCFB_SET_CLR_KEY       _IOW('F', 0x22, struct mxcfb_color_key)
// #define MXCFB_SET_OVERLAY_POS   _IOWR('F', 0x24, struct mxcfb_pos)
// #define MXCFB_GET_FB_IPU_CHAN 	_IOR('F', 0x25, u_int32_t)
// #define MXCFB_SET_LOC_ALPHA     _IOWR('F', 0x26, struct mxcfb_loc_alpha)
// #define MXCFB_SET_LOC_ALP_BUF    _IOW('F', 0x27, unsigned long)
// #define MXCFB_SET_GAMMA	       _IOW('F', 0x28, struct mxcfb_gamma)
// #define MXCFB_GET_FB_IPU_DI 	_IOR('F', 0x29, u_int32_t)
// #define MXCFB_GET_DIFMT	       _IOR('F', 0x2A, u_int32_t)
// #define MXCFB_GET_FB_BLANK     _IOR('F', 0x2B, u_int32_t)
// #define MXCFB_SET_DIFMT		_IOW('F', 0x2C, u_int32_t)
// #define MXCFB_CSC_UPDATE	_IOW('F', 0x2D, struct mxcfb_csc_matrix)
// #define MXCFB_SET_GPU_SPLIT_FMT	_IOW('F', 0x2F, struct mxcfb_gpu_split_fmt)
// #define MXCFB_SET_PREFETCH	_IOW('F', 0x30, int)
// #define MXCFB_GET_PREFETCH	_IOR('F', 0x31, int)
//
// /* IOCTLs for E-ink panel updates */
// #define MXCFB_SET_WAVEFORM_MODES	_IOW('F', 0x2B, struct mxcfb_waveform_modes)
// #define MXCFB_SET_TEMPERATURE		_IOW('F', 0x2C, int32_t)
// #define MXCFB_SET_AUTO_UPDATE_MODE	_IOW('F', 0x2D, __u32)
// #define MXCFB_SEND_UPDATE		_IOW('F', 0x2E, struct mxcfb_update_data)
// #define MXCFB_WAIT_FOR_UPDATE_COMPLETE	_IOWR('F', 0x2F, struct mxcfb_update_marker_data)
// #define MXCFB_SET_PWRDOWN_DELAY		_IOW('F', 0x30, int32_t)
// #define MXCFB_GET_PWRDOWN_DELAY		_IOR('F', 0x31, int32_t)
// #define MXCFB_SET_UPDATE_SCHEME		_IOW('F', 0x32, __u32)
// #define MXCFB_GET_WORK_BUFFER		_IOWR('F', 0x34, unsigned long)
// #define MXCFB_SET_TEMP_AUTO_UPDATE_PERIOD      _IOW('F', 0x36, int32_t)
// #define MXCFB_DISABLE_EPDC_ACCESS	_IO('F', 0x35)
// #define MXCFB_ENABLE_EPDC_ACCESS	_IO('F', 0x36)

// convert the above code to Go

const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS

	_IOC_WRITE = 1
	_IOC_READ  = 2
)

func IOWR(t, nr, size uintptr) uintptr {
	return (_IOC_READ|_IOC_WRITE)<<_IOC_DIRSHIFT | t<<_IOC_TYPESHIFT | nr<<_IOC_NRSHIFT | size<<_IOC_SIZESHIFT
}
func IOW(t, nr, size uintptr) uintptr {
	return (_IOC_WRITE)<<_IOC_DIRSHIFT | t<<_IOC_TYPESHIFT | nr<<_IOC_NRSHIFT | size<<_IOC_SIZESHIFT
}
func IOR(t, nr, size uintptr) uintptr {
	return (_IOC_READ)<<_IOC_DIRSHIFT | t<<_IOC_TYPESHIFT | nr<<_IOC_NRSHIFT | size<<_IOC_SIZESHIFT
}

const (
	MXCFB_WAIT_FOR_VSYNC         uintptr = 0x20
	MXCFB_SET_GBL_ALPHA          uintptr = 0x21
	MXCFB_SET_CLR_KEY            uintptr = 0x22
	MXCFB_SET_OVERLAY_POSuintptr         = 0x24
	MXCFB_GET_FB_IPU_CHAN        uintptr = 0x25
	MXCFB_SET_LOC_ALPHA          uintptr = 0x26
	MXCFB_SET_LOC_ALP_BUF        uintptr = 0x27
	MXCFB_SET_GAMMA                      = 0x28
	MXCFB_GET_FB_IPU_DI                  = 0x29
	MXCFB_GET_DIFMT                      = 0x2A
	MXCFB_GET_FB_BLANK                   = 0x2B
	MXCFB_SET_DIFMT                      = 0x2C
	MXCFB_CSC_UPDATE                     = 0x2D
	MXCFB_SET_GPU_SPLIT_FMT              = 0x2F
	MXCFB_SET_PREFETCH                   = 0x30
	MXCFB_GET_PREFETCH                   = 0x31

	MXCFB_SET_WAVEFORM_MODES                  = 0x2B
	MXCFB_SET_TEMPERATURE                     = 0x2C
	MXCFB_SET_AUTO_UPDATE_MODE                = 0x2D
	MXCFB_SEND_UPDATE                 uintptr = 0x2E
	MXCFB_WAIT_FOR_UPDATE_COMPLETE    uintptr = 0x2F
	MXCFB_SET_PWRDOWN_DELAY                   = 0x30
	MXCFB_GET_PWRDOWN_DELAY                   = 0x31
	MXCFB_SET_UPDATE_SCHEME                   = 0x32
	MXCFB_GET_WORK_BUFFER                     = 0x34
	MXCFB_SET_TEMP_AUTO_UPDATE_PERIOD         = 0x36
	MXCFB_DISABLE_EPDC_ACCESS                 = 0x35
	MXCFB_ENABLE_EPDC_ACCESS                  = 0x36
)
