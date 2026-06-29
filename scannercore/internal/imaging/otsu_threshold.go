package imaging

import (
	"image"
	"math"
)

// Naive O(N×256) Otsu search: 256 full-image scans, float64 accumulation per pixel.
// 1920×1080 benchmark: ~1050ms on Ryzen 7 (L3-resident). Phone DRAM-eviction: est.3-5×
// slower → ~60-100s for 20 pages. Intentionally unoptimized baseline; histogram in
// otsuThreshold(histStats) collapses this to one O(N) pass + O(256) bin search.

func otsuThreshold(in *image.Gray) uint8 {
	var out uint8
	maxVariance := 0.0

	for t := range 256 {
		//for each possible threshold we need to iterate over the whole image so we determine the optimal t (sort of uneducated brute-force)
		countB, countF := 0.0, 0.0
		sumB, sumF := 0.0, 0.0

		//now we need to separate into two groups the background and foreground
		for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
			for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
				pixelValue := in.Pix[in.PixOffset(x, y)]
				if pixelValue > uint8(t) { //foreground
					countF++
					sumF += float64(pixelValue)
				} else { //background
					countB++
					sumB += float64(pixelValue)
				}
			}
		}

		totalPixels := countF + countB
		if countB == 0 || countF == 0 {
			continue
		}
		wB := countB / totalPixels
		wF := countF / totalPixels
		muB := sumB / countB
		muF := sumF / countF

		variance := wB * wF * math.Pow(muB-muF, 2)
		if variance > maxVariance {
			maxVariance = variance
			out = uint8(t)
		}
	}

	return out
}
