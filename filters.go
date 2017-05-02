package imagexp

import (
	"image/color"
	"math"
)

func BasicGrayscale(r, g, b, _ uint32) color.Gray16 {
	avg := float64((r + g + b) / 3)

	return color.Gray16{uint16(math.Ceil(avg))}
}

func ImprovedGrayscale(r, g, b, _ uint32) color.Gray16 {
	avg := float64(0.3)*float64(r) + float64(0.59)*float64(g) + float64(0.11)*float64(b)

	return color.Gray16{uint16(math.Ceil(avg))}
}

func Desaturation(r, g, b, a uint32) color.Gray16 {
	avg := float64(maxOfThree(r, g, b, a)+minOfThree(r, g, b, a)) / 2
	return color.Gray16{uint16(math.Ceil(avg))}
}

func DecompositionMax(r, g, b, a uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(maxOfThree(r, g, b, a))))}
}

func DecompositionMin(r, g, b, a uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(minOfThree(r, g, b, a))))}
}

func maxOfThree(r, g, b, _ uint32) uint32 {
	return max(max(r, g), b)
}

func minOfThree(r, g, b, _ uint32) uint32 {
	return min(min(r, g), b)
}

// This is how, I'll do it, Until I figure out a better way
func SingleChannelRed(r, _, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(r)))}
}

func SingleChannelGreen(_, g, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(g)))}
}

func SingleChannelBlue(_, _, b, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(b)))}
}

func RedFilter(r, g, b, a uint32) color.RGBA64 {

	if !(r > b) || !(r > g) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(255)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func GreenFilter(r, g, b, a uint32) color.RGBA64 {
	if !(g > r) || !(g > b) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(255)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func BlueFilter(r, g, b, a uint32) color.RGBA64 {
	if !(b > g) || !(b > r) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(255)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
