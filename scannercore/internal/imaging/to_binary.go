package imaging

import "image"

func applyBinary(in *image.Gray, t uint8) *image.Gray {
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
