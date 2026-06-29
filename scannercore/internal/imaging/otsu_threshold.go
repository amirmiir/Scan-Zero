package imaging

import (
	"image"
	"math"
)

type histStats struct {
	counts      [256]int
	n           int
	weightedSum int
}

// O(N+256) Otsu: single histogram pass
// wB/sumB accumulate incrementally; wF/sumF derived from histogram totals — no secondaccumulator.
// float64 means
func otsuThreshold(histogram histStats) uint8 {
	var out uint8
	maxVariance := 0.0
	// histogram := computeHistogram(in)
	countB, countF := 0.0, float64(histogram.n)
	sumB, sumF := 0.0, 0.0

	for t := range 256 {
		//for each possible threshold we need to iterate over the whole image so we determine the optimal t (sort of uneducated brute-force)

		//now we need to separate into two groups the background and foreground
		//BUT NOW WE USE THE HISTOGRAM
		current := histogram.counts[t]
		countB += float64(current)
		countF = float64(histogram.n) - countB

		sumB += float64(t) * float64(current)
		sumF = float64(histogram.weightedSum) - sumB

		if countB == 0 || countF == 0 {
			continue
		}
		wB := countB / float64(histogram.n)
		wF := countF / float64(histogram.n)
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

func computeHistogram(in *image.Gray) histStats {
	//single pass per image
	//okay, for this, we need to have a histogram after we have received the image in grayscale;
	//O(N) N is the WxH?
	var histogram histStats

	for y := in.Bounds().Min.Y; y < in.Bounds().Max.Y; y++ {
		for x := in.Bounds().Min.X; x < in.Bounds().Max.X; x++ {
			histogram.counts[in.Pix[in.PixOffset(x, y)]]++
		}
	}
	histogram.n = (in.Bounds().Max.X - in.Bounds().Min.X) * (in.Bounds().Max.Y - in.Bounds().Min.Y)
	for intensity := range 256 {
		histogram.weightedSum += histogram.counts[intensity] * intensity
	}

	return histogram
}
