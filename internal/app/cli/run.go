package cli

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/rifux/Go-BasicBorderScanner/internal/imageutil"
)

type encoderFn func(w *os.File, m image.Image) error

func encoderFor(fmtName string) (encoderFn, error) {
	switch strings.ToLower(fmtName) {
	case "png":
		return func(w *os.File, m image.Image) error { return png.Encode(w, m) }, nil
	case "jpeg", "jpg":
		return func(w *os.File, m image.Image) error {
			return jpeg.Encode(w, m, &jpeg.Options{Quality: 90})
		}, nil
	case "gif":
		return func(w *os.File, m image.Image) error { return gif.Encode(w, m, nil) }, nil
	case "tiff":
		return func(w *os.File, m image.Image) error {
			return tiff.Encode(w, m, &tiff.Options{Compression: tiff.Deflate})
		}, nil
	case "bmp":
		return func(w *os.File, m image.Image) error { return bmp.Encode(w, m) }, nil
	default:
		return nil, fmt.Errorf("unsupported output format %q", fmtName)
	}
}

func Run(ctx context.Context, _ string) error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	inPath := fs.String("in", "", "input image")
	outPath := fs.String("out", "out", "output file")
	outFmt := fs.String("outfmt", "png", "output format")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if *inPath == "" {
		return fmt.Errorf("flag -in is required")
	}

	src, err := os.Open(*inPath)
	if err != nil {
		return err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	binImg, err := imageutil.OtsuBinarize(ctx, img)
	if err != nil {
		return err
	}

	outImg, err := imageutil.DrawScannedContours(ctx, binImg)
	if err != nil {
		return err
	}

	enc, err := encoderFor(*outFmt)
	if err != nil {
		return err
	}

	outFile := *outPath
	if ext := filepath.Ext(outFile); ext == "" || !strings.EqualFold(ext, "."+*outFmt) {
		outFile += "." + *outFmt
	}

	dst, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	return enc(dst, outImg)
}
