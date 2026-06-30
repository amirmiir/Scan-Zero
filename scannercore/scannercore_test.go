package scannercore

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	//1.1. define a input
	//1.2. send it to scannercore
	//1.3. validate if response is "pong: "+input
}

func TestUnpackPaths(t *testing.T) {
	//1.1 sending a valid json array of two paths
	//2.1 malformed json
	//3.1 empty array
}

// writeSyntheticJPEG creates a 100×100 gray gradient JPEG at the given path.
func writeSyntheticJPEG(t *testing.T, path string) {
	t.Helper()
	img := image.NewGray(image.Rect(0, 0, 100, 100))
	for i := range img.Pix {
		img.Pix[i] = uint8(i % 256)
	}
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, jpeg.Encode(f, img, &jpeg.Options{Quality: 90}))
}

// ProcessImage must decode, transform, and write a valid JPEG — not copy the input.
func TestProcessImage_Success(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.jpg")
	outputPath := filepath.Join(dir, "output.jpg")

	writeSyntheticJPEG(t, inputPath)
	inputBytes, err := os.ReadFile(inputPath)
	require.NoError(t, err)

	gotPath, err := ProcessImage(inputPath, outputPath)

	require.NoError(t, err)
	assert.Equal(t, outputPath, gotPath)

	// output must exist and be a valid JPEG with same dimensions as input
	outBytes, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	outImg, _, err := image.Decode(bytes.NewReader(outBytes))
	require.NoError(t, err)
	assert.Equal(t, 100, outImg.Bounds().Dx())
	assert.Equal(t, 100, outImg.Bounds().Dy())

	// pipeline must have transformed pixels — output must not be byte-identical to input
	assert.NotEqual(t, inputBytes, outBytes)
}

// Non-existent input must return non-nil error and empty path.
func TestProcessImage_InvalidPath(t *testing.T) {
	gotPath, err := ProcessImage("/nonexistent/input.jpg", "/tmp/out.jpg")

	assert.Error(t, err)
	assert.Empty(t, gotPath)
}
