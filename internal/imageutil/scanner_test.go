package imageutil

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"testing"
)

//	[ ! ] THESE TESTS SHOULD BE REWRITTEN TO MEET BASIC DEVELOPER AGREEMENTS

func TestScanner(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Debug("opening file 'input.png'")
	file, errOpen := os.Open("outputOtsu1.png")
	if errOpen != nil {
		t.Fatalf("got error while trying to open input.png file:\n%s", errOpen.Error())
	}
	defer file.Close()
	logger.Debug("file 'input.png' opened")

	logger.Debug("decoding input image")
	img, _, errDecode := image.Decode(file)
	if errDecode != nil {
		t.Fatalf("got error while trying to decode image:\n%s", errDecode.Error())
	}
	logger.Debug("input image decoded successfully")

	logger.Debug(fmt.Sprintf("first pixel color: %v, image bounds: %v", img.At(0, 0), img.Bounds()))

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	logger.Debug("starting process of inverting image")
	type result struct {
		img image.Image
		err error
	}
	ch := make(chan result, 1) // buffer so the goroutine never blocks
	go func() {
		img, err := DrawScannedContours(ctx, img)
		ch <- result{img, err}
	}()
	//time.AfterFunc(time.Millisecond*10, stop)
	r := <-ch
	imgInvert, errInv := r.img, r.err
	if errInv != nil {
		t.Fatalf("got error while trying to invert image:\n%s", errInv.Error())
	}
	logger.Debug("image has been inverted")

	logger.Debug(fmt.Sprintf("first pixel color: %v, image bounds: %v", imgInvert.At(0, 0), imgInvert.Bounds()))

	logger.Debug("saving 'outputScan.png' image")
	fileSave, errSave := os.Create("output.png")
	if errSave != nil {
		t.Fatalf("got error while trying to save 'output.png' file:\n%s", errSave.Error())
	}
	defer fileSave.Close()
	png.Encode(fileSave, imgInvert)
	logger.Debug("output image saved successfully")
}
