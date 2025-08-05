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

	"github.com/rifux/Go-BasicBorderScanner/internal/imageutil"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type encoderFn func(w *os.File, m image.Image) error

// encoderFor gets the correct image encoder function for a given format name.
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
	case "tiff", "tif":
		return func(w *os.File, m image.Image) error {
			return tiff.Encode(w, m, &tiff.Options{Compression: tiff.Deflate})
		}, nil
	case "bmp":
		return func(w *os.File, m image.Image) error { return bmp.Encode(w, m) }, nil
	default:
		return nil, fmt.Errorf("unsupported output format %q (from file extension)", fmtName)
	}
}

// Run executes the command-line interface logic.
func Run(ctx context.Context, args []string) error {
	cliFlags := flag.NewFlagSet("cli", flag.ExitOnError)

	// Define flags for the CLI mode. Note that -outfmt is now gone.
	inPath := cliFlags.String("in", "", "input image (png, jpg, etc) (required)")
	outPath := cliFlags.String("out", "out.png", "output file (format inferred from extension)")
	logMode := cliFlags.String("log", "auto", "log output mode: auto|json|text")

	cliFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s cli [flags]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Flags for cli command:")
		cliFlags.PrintDefaults()
	}

	if err := cliFlags.Parse(args); err != nil {
		return err
	}

	if *inPath == "" {
		cliFlags.Usage()
		return fmt.Errorf("\nError: flag -in is required")
	}

	// --- NEW LOGIC STARTS HERE ---

	outFilename := *outPath
	var outputFormat string

	// Get the file extension from the output filename.
	ext := filepath.Ext(outFilename)

	if ext == "" {
		// If no extension, default to "png".
		outputFormat = "png"
		outFilename += "." + outputFormat // Append ".png" to the filename.
	} else {
		// If an extension exists, use it as the format.
		// Trim the leading dot, e.g., ".jpg" becomes "jpg".
		outputFormat = strings.TrimPrefix(ext, ".")
	}

	fmt.Printf("Input: %s\n", *inPath)
	fmt.Printf("Output: %s (format: %s)\n", outFilename, outputFormat)
	fmt.Printf("Log mode: %s\n", *logMode)

	// Get the encoder function based on the determined format.
	enc, err := encoderFor(outputFormat)
	if err != nil {
		// This will now catch invalid extensions like "result.txt".
		return err
	}

	// --- END OF NEW LOGIC ---

	// Open input file
	src, err := os.Open(*inPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Decode image
	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	// Process image (Otsu binarization and contour drawing)
	binImg, err := imageutil.OtsuBinarize(ctx, img)
	if err != nil {
		return err
	}

	outImg, err := imageutil.DrawScannedContours(ctx, binImg)
	if err != nil {
		return err
	}

	// Create output file
	dst, err := os.Create(outFilename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Encode and save the final image
	return enc(dst, outImg)
}
