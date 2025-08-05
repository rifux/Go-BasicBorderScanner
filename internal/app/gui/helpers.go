package gui

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

// ---- helpers for extra formats ----
var (
	openExts = []string{".png", ".jpg", ".jpeg", ".gif", ".tiff", ".tif", ".bmp"}

	encoders = map[string]func(w *os.File, m image.Image) error{
		".png":  func(w *os.File, m image.Image) error { return png.Encode(w, m) },
		".jpg":  func(w *os.File, m image.Image) error { return jpeg.Encode(w, m, &jpeg.Options{Quality: 90}) },
		".jpeg": func(w *os.File, m image.Image) error { return jpeg.Encode(w, m, &jpeg.Options{Quality: 90}) },
		".gif":  func(w *os.File, m image.Image) error { return gif.Encode(w, m, nil) },
		".tiff": func(w *os.File, m image.Image) error {
			return tiff.Encode(w, m, &tiff.Options{Compression: tiff.Deflate})
		},
		".tif": func(w *os.File, m image.Image) error {
			return tiff.Encode(w, m, &tiff.Options{Compression: tiff.Deflate})
		},
		".bmp": func(w *os.File, m image.Image) error { return bmp.Encode(w, m) },
	}
)

func withExt(name, ext string) string {
	if filepath.Ext(name) == "" {
		return name + ext
	}
	return name
}
