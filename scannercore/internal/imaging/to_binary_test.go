package imaging

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Pixels at or below threshold → 0 (ink). Pixels above → 255 (paper).
// Boundary: pixel == t is ink, not paper.
func TestApplyBinary(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 3, 1))
	img.Pix[0] = 50  // below t
	img.Pix[1] = 128 // at t — must be ink (0)
	img.Pix[2] = 200 // above t

	out := applyBinary(img, 128)

	assert.Equal(t, uint8(0), out.Pix[0])
	assert.Equal(t, uint8(0), out.Pix[1])
	assert.Equal(t, uint8(255), out.Pix[2])
}

// Input image must not be mutated.
func TestApplyBinary_NoMutate(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 2, 1))
	img.Pix[0] = 10
	img.Pix[1] = 200

	applyBinary(img, 128)

	assert.Equal(t, uint8(10), img.Pix[0])
	assert.Equal(t, uint8(200), img.Pix[1])
}

// Three zones: below vLow → 0, above vHigh → 255, band → original value.
func TestApplyContrastAnchor_ThreeZones(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 3, 1))
	img.Pix[0] = 50  // below vLow (100-20=80) → black
	img.Pix[1] = 90  // inside band [80,120] → preserved
	img.Pix[2] = 150 // above vHigh (100+20=120) → white

	out := applyContrastAnchor(img, 100, 20, 20)

	assert.Equal(t, uint8(0), out.Pix[0])
	assert.Equal(t, uint8(90), out.Pix[1])
	assert.Equal(t, uint8(255), out.Pix[2])
}

// vLow clamp: t-l underflows. Without clamp, vLow wraps to uint8(10-50)=216 and
// snaps pixel=100 to black (100 ≤ 216). With clamp vLow=0, vHigh=20, pixel=100
// is above vHigh → white (255).
func TestApplyContrastAnchor_LowClamp(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 1, 1))
	img.Pix[0] = 100

	out := applyContrastAnchor(img, 10, 50, 10) // vLow=max(0,10-50)=0, vHigh=20

	assert.Equal(t, uint8(255), out.Pix[0])
}

// vHigh clamp: t+h overflows. Without clamp, vHigh wraps to uint8(300)=44 and
// snaps pixel=100 to white (100 ≥ 44). With clamp vHigh=255, vLow=240, pixel=100
// is below vLow → black (0).
func TestApplyContrastAnchor_HighClamp(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 1, 1))
	img.Pix[0] = 100

	out := applyContrastAnchor(img, 250, 10, 50) // vLow=240, vHigh=min(255,300)=255

	assert.Equal(t, uint8(0), out.Pix[0])
}
