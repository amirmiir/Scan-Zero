package imaging

import "image"

func ApplyBinary(in *image.Gray, t uint8) *image.Gray {
	//we copy the size
	out := image.NewGray(in.Bounds())

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			if in.Pix[in.PixOffset(x, y)] <= t {
				out.Pix[out.PixOffset(x, y)] = 0
			} else {
				out.Pix[out.PixOffset(x, y)] = 255
			}
		}
	}

	return out
}

func ApplyContrastAnchor(in *image.Gray, t, l, h uint8) *image.Gray {
	out := image.NewGray(in.Bounds())

	vLowInt := int(t) - int(l)
	vHighInt := int(t) + int(h)

	if vLowInt < 0 {
		vLowInt = 0
	}
	if vHighInt > 255 {
		vHighInt = 255
	}

	vLow := uint8(vLowInt)
	vHigh := uint8(vHighInt)

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			if in.Pix[in.PixOffset(x, y)] <= vLow {
				out.Pix[out.PixOffset(x, y)] = 0
			} else if in.Pix[in.PixOffset(x, y)] >= vHigh {
				out.Pix[out.PixOffset(x, y)] = 255
			} else {
				out.Pix[out.PixOffset(x, y)] = in.Pix[in.PixOffset(x, y)]
			}
		}
	}

	return out
}
