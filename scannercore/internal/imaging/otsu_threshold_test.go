package imaging

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Bimodal: half pixels at 50 (ink), half at 200 (paper).
// Threshold must land between the two clusters.
func TestOtsuThreshold_Bimodal(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 2, 2))
	img.Pix[0] = 50
	img.Pix[1] = 50
	img.Pix[2] = 200
	img.Pix[3] = 200

	result := otsuThreshold(computeHistogram(img))

	assert.GreaterOrEqual(t, result, uint8(50))
	assert.Less(t, result, uint8(200))
}

// Degenerate: all pixels identical — no valid split exists.
// Must not panic; empty-class guard holds throughout.
func TestOtsuThreshold_Degenerate(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := range img.Pix {
		img.Pix[i] = 128
	}

	assert.NotPanics(t, func() { otsuThreshold(computeHistogram(img)) })
}

// Fixture {0,0,128,255}: counts, total pixel count, and weighted sum must be exact.
func TestComputeHistogram_Fixture(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 2, 2))
	img.Pix[0] = 0
	img.Pix[1] = 0
	img.Pix[2] = 128
	img.Pix[3] = 255

	h := computeHistogram(img)

	assert.Equal(t, 2, h.counts[0])
	assert.Equal(t, 1, h.counts[128])
	assert.Equal(t, 1, h.counts[255])
	assert.Equal(t, 4, h.n)
	assert.Equal(t, 383, h.weightedSum) // 0×2 + 128×1 + 255×1
}

// Stride guard: Stride > Width means padding bytes exist between rows.
// computeHistogram must count only real pixels, not padding.
func TestComputeHistogram_StrideGuard(t *testing.T) {
	img := &image.Gray{
		Pix:    make([]uint8, 8), // 2 rows × Stride 4, but Width=2
		Stride: 4,
		Rect:   image.Rect(0, 0, 2, 2),
	}
	// real pixels: offsets 0,1 (row 0) and 4,5 (row 1)
	img.Pix[0] = 10
	img.Pix[1] = 20
	img.Pix[4] = 30
	img.Pix[5] = 40
	// padding bytes at offsets 2,3,6,7 — set nonzero to prove they are excluded
	img.Pix[2] = 255
	img.Pix[3] = 255
	img.Pix[6] = 255
	img.Pix[7] = 255

	h := computeHistogram(img)

	assert.Equal(t, 4, h.n)
	assert.Equal(t, 0, h.counts[255]) // padding must not appear in any bin
}

func BenchmarkOtsuThreshold(b *testing.B) {
	img := image.NewGray(image.Rect(0, 0, 1920, 1080))
	// fill with some non-uniform data so classes are non-empty
	for i := range img.Pix {
		img.Pix[i] = uint8(i % 256)
	}
	b.ResetTimer()
	for range b.N {
		otsuThreshold(computeHistogram(img))
	}
}
