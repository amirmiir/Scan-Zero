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

	result := otsuThreshold(img)

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

	assert.NotPanics(t, func() { otsuThreshold(img) })
}

func BenchmarkOtsuThreshold(b *testing.B) {
	img := image.NewGray(image.Rect(0, 0, 1920, 1080))
	// fill with some non-uniform data so classes are non-empty
	for i := range img.Pix {
		img.Pix[i] = uint8(i % 256)
	}
	b.ResetTimer()
	for range b.N {
		otsuThreshold(img)
	}
}
