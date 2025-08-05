package imageutil

import (
	"context"
	"image"
	"image/color"
)

// OtsuBinarize applies Otsu's method to binarize an image.
// It automatically determines the optimal threshold to separate pixels into foreground and background.
func OtsuBinarize(ctx context.Context, src image.Image) (image.Image, error) {
	bounds := src.Bounds()
	out := image.NewGray(bounds)

	// Step 1: Compute the grayscale histogram.
	histogram := make([]int, 256)
	totalPixels := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(src.At(x, y)).(color.Gray)
			histogram[gray.Y]++
			totalPixels++
		}
	}

	// Step 2: Calculate the total sum of the histogram.
	var sum float64
	for i, h := range histogram {
		sum += float64(i) * float64(h)
	}

	var sumB float64
	var wB, wF int
	var maxVariance float64
	var threshold int

	// Step 3: Find the optimal threshold.
	for t := 0; t < 256; t++ {
		// Check for context cancellation.
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		wB += histogram[t] // Weight of the background
		if wB == 0 {
			continue
		}

		wF = totalPixels - wB // Weight of the foreground
		if wF == 0 {
			break
		}

		sumB += float64(t) * float64(histogram[t])

		meanB := sumB / float64(wB)
		meanF := (sum - sumB) / float64(wF)

		// Calculate the between-class variance.
		variance := float64(wB) * float64(wF) * (meanB - meanF) * (meanB - meanF)

		if variance > maxVariance {
			maxVariance = variance
			threshold = t
		}
	}

	// Step 4: Apply the threshold to create the binary image.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// Check for context cancellation.
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(src.At(x, y)).(color.Gray)
			if gray.Y > uint8(threshold) {
				out.SetGray(x, y, color.Gray{Y: 255}) // Foreground (white)
			} else {
				out.SetGray(x, y, color.Gray{Y: 0}) // Background (black)
			}
		}
	}

	return out, ctx.Err()
}
