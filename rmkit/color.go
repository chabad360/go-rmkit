package rmkit

const (
	White RGB565 = 0xFFFF
	Gray  RGB565 = 0x4444
	Black RGB565 = 0
)

type RGBColor struct {
	R, G, B uint8
}

func Gray32(n int) RGB565 {
	if n == 0 {
		return RGB565((n << 11) | (((2 * n) | 0) << 5) | n)
	}
	return RGB565((n << 11) | (((2 * n) | 1) << 5) | n)
}

func FromFloat(n float32) RGB565 {
	return Gray32(int(31.0 * n))
}

func ToRGB8(s RGB565) RGBColor {
	return RGBColor{
		uint8(((uint32(s) & 0xf800) >> 11) / 31. * 255.),
		uint8(((uint32(s) & 0x7e0) >> 5) / 63. * 255.),
		uint8((uint32(s) & 0x1f) / 31. * 255.),
	}
}

func ToFloat(c RGB565) float32 {
	return float32((uint32(c)>>11)&31)*(0.21/31) +
		float32((uint32(c)>>5)&63)*(0.72/63) +
		float32(uint32(c)&31)*(0.07/31)
}

var (
	Gray1  = Gray32(2)
	Gray2  = Gray32(4)
	Gray3  = Gray32(6)
	Gray4  = Gray32(8)
	Gray5  = Gray32(10)
	Gray6  = Gray32(12)
	Gray7  = Gray32(14)
	Gray8  = Gray32(16)
	Gray9  = Gray32(18)
	Gray10 = Gray32(20)
	Gray11 = Gray32(22)
	Gray12 = Gray32(24)
	Gray13 = Gray32(26)
	Gray14 = Gray32(28)

	Scale16 = []RGB565{
		Black, Gray1, Gray2, Gray3, Gray4, Gray5, Gray6, Gray7,
		Gray8, Gray9, Gray10, Gray11, Gray12, Gray13, Gray14, White,
	}
)

func Quantize2(c float32) RGB565 {
	if c >= 0.5 {
		return White
	}
	return Black
}

func Quantize4(c float32) RGB565 {
	if c < 0.25 {
		return Black
	}
	if c >= 0.75 {
		return White
	}
	if c >= 0.5 {
		return Gray10
	}
	return Gray5
}

func Quantize16(c float32) RGB565 {
	if c < (1 / 16.0) {
		return Black
	}
	if c >= (15 / 16.0) {
		return White
	}
	return Scale16[int(c*15+0.5)]
}
