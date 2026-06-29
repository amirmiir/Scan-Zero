package imaging

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Test A — generic fallback luma path.

Feeds *image.RGBA to force the non-YCbCr branch.
Verifies integer luma weights (77r+150g+29b) >> 8
*/
func TestToGray_GenericLuma(t *testing.T) {
	cases := []struct {
		name     string
		in       color.RGBA
		wantGray uint8
	}{
		{"red", color.RGBA{R: 255, G: 0, B: 0, A: 255}, 76},
		{"green", color.RGBA{R: 0, G: 255, B: 0, A: 255}, 149},
		{"blue", color.RGBA{R: 0, G: 0, B: 255, A: 255}, 28},
		{"white", color.RGBA{R: 255, G: 255, B: 255, A: 255}, 255},
		{"black", color.RGBA{R: 0, G: 0, B: 0, A: 255}, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := image.NewRGBA(image.Rect(0, 0, 1, 1))
			src.Set(0, 0, tc.in)

			got := ToGray(src).GrayAt(0, 0).Y

			assert.Equal(t, tc.wantGray, got)
		})
	}
}

/*
Test B — YCbCr fast path + stride divergence.

Plain NewYCbCr gives YStride == Dx, which HIDES stride bugs.
*/
func TestToGray_YCbCrStride(t *testing.T) {
	parent := image.NewYCbCr(image.Rect(0, 0, 8, 8), image.YCbCrSubsampleRatio444)
	for i := range parent.Y {
		parent.Y[i] = byte(i) // known ramp
	}

	sub := parent.SubImage(image.Rect(2, 0, 6, 4)).(*image.YCbCr)
	assert.Equal(t, 4, sub.Bounds().Dx())
	assert.Equal(t, 8, sub.YStride)

	out := ToGray(sub)
	for y := sub.Bounds().Min.Y; y < sub.Bounds().Max.Y; y++ {
		for x := sub.Bounds().Min.X; x < sub.Bounds().Max.X; x++ {
			want := parent.Y[parent.YOffset(x, y)]
			got := out.GrayAt(x, y).Y

			assert.Equal(t, want, got)
		}

	}
}

/*
Test C — SubImage with non-zero bounds.Min.Y

Verifies Y-plane copy starts at the correct parent row when Min.Y > 0
*/

func TestToGray_YCbCrBoundsOffset(t *testing.T) {
	parent := image.NewYCbCr(image.Rect(0, 0, 8, 8), image.YCbCrSubsampleRatio444)
	for i := range parent.Y {
		parent.Y[i] = byte(i)
	}

	sub := parent.SubImage(image.Rect(0, 4, 8, 8)).(*image.YCbCr)
	out := ToGray(sub)
	for y := sub.Bounds().Min.Y; y < sub.Bounds().Max.Y; y++ {
		for x := sub.Bounds().Min.X; x < sub.Bounds().Max.X; x++ {
			want := parent.Y[parent.YOffset(x, y)]
			got := out.GrayAt(x, y).Y

			assert.Equal(t, want, got)
		}

	}

}

// sink defeats dead-code elimination in the benchmark.
var sink *image.Gray

// BenchmarkToGray — fast path, realistic frame. allocs/op must be CONSTANT
func BenchmarkToGray(b *testing.B) {
	img := image.NewYCbCr(image.Rect(0, 0, 1920, 1080), image.YCbCrSubsampleRatio420)

	for b.Loop() {
		sink = ToGray(img)
	}
}
