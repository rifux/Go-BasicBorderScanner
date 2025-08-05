package imageutil

import (
	"context"
	"image"
	"image/color"
)

// ---- basic image inversion ----
func Invert(ctx context.Context, src image.Image) (image.Image, error) {
	out := image.NewRGBA(src.Bounds())

	for y := out.Bounds().Min.Y; y < out.Bounds().Max.Y; y++ {
		if err := ctx.Err(); err != nil {
			return nil, ctx.Err()
		}

		for x := out.Bounds().Min.X; x <= out.Bounds().Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			out.Set(x, y, color.RGBA{
				R: 255 - uint8(r>>8),
				G: 255 - uint8(g>>8),
				B: 255 - uint8(b>>8),
				A: uint8(a >> 8),
			})
		}
	}

	return out, ctx.Err()
}
