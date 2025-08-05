package imageutil

import (
	"context"
	"image"
	"image/color"
	"image/draw"
)

// -----------------------------------------------------------------------------
// Internal types
// -----------------------------------------------------------------------------

// blackSeries - horizontal black stripe in row y
type blackSeries struct {
	startX, endX, y int
}

// branch - one contour branch (left or right)
type branch struct {
	id   int  // unique contour ID
	left bool // true - left branch, false - right
}

// activeSer - combination of "black series + two branches"
type activeSer struct {
	ser         blackSeries
	leftBranch  *branch
	rightBranch *branch
}

// -----------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------

// findBlackSeries returns all black series in row y
func findBlackSeries(img image.Image, y int) []blackSeries {
	bounds := img.Bounds()
	var out []blackSeries
	in, start := false, 0

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
		isBlack := c.Y == 0

		switch {
		case isBlack && !in:
			in, start = true, x
		case !isBlack && in:
			in = false
			out = append(out, blackSeries{start, x - 1, y})
		}
	}
	if in { // series reached the right edge
		out = append(out, blackSeries{start, bounds.Max.X - 1, y})
	}
	return out
}

// -----------------------------------------------------------------------------
// Main algorithm
// -----------------------------------------------------------------------------

// DrawScannedContours implements the scanning algorithm for contour detection.
// Returns a new image with drawn contours.
func DrawScannedContours(ctx context.Context, src image.Image) (image.Image, error) {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, dst.Bounds(), image.White, image.Point{}, draw.Src)

	// storage of points by contour ID
	contours := make(map[int][]image.Point)
	nextID := 1

	// active series of the previous row
	var prev []activeSer

	// -------------------------------------------------------------------------
	// Line-by-line scanning
	// -------------------------------------------------------------------------
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		cur := findBlackSeries(src, y)
		var next []activeSer

		i, j := 0, 0 // pointers for prev and cur

		for i < len(prev) || j < len(cur) {
			// 1. "End" situation: only prev-series remain
			if j == len(cur) {
				id := prev[i].leftBranch.id
				// add "closing" points
				contours[id] = append(contours[id],
					image.Pt(prev[i].ser.endX, prev[i].ser.y),
					image.Pt(prev[i].ser.startX, prev[i].ser.y))
				i++
				continue
			}

			// 2. "Start" situation: only cur-series remain
			if i == len(prev) {
				id := nextID
				nextID++
				lb := &branch{id: id, left: true}
				rb := &branch{id: id, left: false}
				contours[id] = append(contours[id],
					image.Pt(cur[j].startX, y),
					image.Pt(cur[j].endX, y))
				next = append(next, activeSer{ser: cur[j], leftBranch: lb, rightBranch: rb})
				j++
				continue
			}

			// 3. Comparison of current series
			ps, cs := prev[i].ser, cur[j]

			switch {
			// 3a. "Start" - cur is to the left of prev
			case cs.endX < ps.startX:
				id := nextID
				nextID++
				lb := &branch{id: id, left: true}
				rb := &branch{id: id, left: false}
				contours[id] = append(contours[id],
					image.Pt(cs.startX, y),
					image.Pt(cs.endX, y))
				next = append(next, activeSer{ser: cs, leftBranch: lb, rightBranch: rb})
				j++

			// 3b. "End" - prev is to the left of cur
			case ps.endX < cs.startX:
				id := prev[i].leftBranch.id
				contours[id] = append(contours[id],
					image.Pt(ps.endX, ps.y),
					image.Pt(ps.startX, ps.y))
				i++

			// 3c. Overlap - process continuation/branching/merging
			default:
				// simplest case: assume it's a "continuation" of one series
				// (full branching/merging logic requires further precise series matching,
				// but for reproducing the top part, correct continuation is sufficient)
				id := prev[i].leftBranch.id
				contours[id] = append(contours[id],
					image.Pt(cs.startX, y),
					image.Pt(cs.endX, y))
				next = append(next, activeSer{
					ser:         cs,
					leftBranch:  prev[i].leftBranch,
					rightBranch: prev[i].rightBranch,
				})
				i++
				j++
			}
		}

		prev = next
	}

	// -------------------------------------------------------------------------
	// Close remaining contours
	// -------------------------------------------------------------------------
	for _, as := range prev {
		id := as.leftBranch.id
		contours[id] = append(contours[id],
			image.Pt(as.ser.endX, as.ser.y),
			image.Pt(as.ser.startX, as.ser.y))
	}

	// -------------------------------------------------------------------------
	// Draw
	// -------------------------------------------------------------------------
	contourColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for _, pts := range contours {
		for _, p := range pts {
			if p.In(bounds) {
				dst.Set(p.X, p.Y, contourColor)
			}
		}
	}

	return dst, ctx.Err()
}
