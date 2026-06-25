package imaging

import (
	"image"
	"image/color"
)

func ToGray(input image.Image) *image.Gray {
	bounds := input.Bounds()        //obtain the size of original image
	target := image.NewGray(bounds) //we define our return based on original

	//fast path
	if ycbcr, ok := input.(*image.YCbCr); ok {
		for row := 0; row < bounds.Dy(); row++ {
			srcOffset := row * ycbcr.YStride
			dstOffset := row * target.Stride
			copy(target.Pix[dstOffset:dstOffset+bounds.Dx()], ycbcr.Y[srcOffset:srcOffset+bounds.Dx()])
		}
	} else { //non-fast path
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := input.At(x, y).RGBA()
				r8 := r >> 8
				g8 := g >> 8
				b8 := b >> 8
				gray := uint8((77*r8 + 150*g8 + 29*b8) >> 8)

				target.SetGray(x, y, color.Gray{Y: gray})
			}
		}
	}

	return target
}
