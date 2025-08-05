package gui

import (
	"context"
	"image"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/rifux/Go-BasicBorderScanner/internal/imageutil"
)

// Run starts the GUI.
func Run(_ context.Context, _ string) error {
	a := app.New()
	w := a.NewWindow("Border-scanner viewer")
	w.Resize(fyne.NewSize(900, 600))

	var (
		inImg  image.Image // original
		binImg image.Image // binarized
		outImg image.Image // processed
	)

	inIV := canvas.NewImageFromImage(nil)
	inIV.FillMode = canvas.ImageFillContain
	inIV.SetMinSize(fyne.NewSize(400, 400))

	outIV := canvas.NewImageFromImage(nil)
	outIV.FillMode = canvas.ImageFillContain
	outIV.SetMinSize(fyne.NewSize(400, 400))

	btnUpload := widget.NewButton("Upload", nil)
	btnRun := widget.NewButton("Run", nil)
	btnSave := widget.NewButton("Save", nil)
	btnStep := widget.NewButton("Detailed viewer", nil)

	btnUpload.OnTapped = func() {
		fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			defer uc.Close()
			img, _, err := image.Decode(uc)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			inImg = img
			inIV.Image = img
			inIV.Refresh()
			outIV.Image = nil
			outIV.Refresh()
			binImg = nil
			outImg = nil
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter(openExts))
		fd.Show()
	}

	btnRun.OnTapped = func() {
		if inImg == nil {
			dialog.ShowInformation("No image", "Load an image first", w)
			return
		}
		ctx := context.TODO()
		var err error
		binImg, err = imageutil.OtsuBinarize(ctx, inImg)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		outImg, err = imageutil.DrawScannedContours(ctx, binImg)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		outIV.Image = outImg
		outIV.Refresh()
	}

	btnSave.OnTapped = func() {
		if outImg == nil {
			dialog.ShowInformation("Nothing to save", "Run the pipeline first", w)
			return
		}
		var labels []string
		for ext := range encoders {
			labels = append(labels, strings.ToUpper(ext[1:]))
		}
		selectedExt := ".png"
		fmtSel := widget.NewSelect(labels, func(s string) {
			selectedExt = "." + strings.ToLower(s)
		})
		fmtSel.SetSelected("PNG")

		dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			defer uc.Close()
			fName := withExt(uc.URI().Path(), selectedExt)
			f, err := os.Create(fName)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			defer f.Close()
			if enc, ok := encoders[selectedExt]; ok {
				if err := enc(f, outImg); err != nil {
					dialog.ShowError(err, w)
				}
			}
		}, w)
	}

	btnStep.OnTapped = func() {
		if binImg == nil || outImg == nil {
			dialog.ShowInformation("No data", "Run the pipeline first", w)
			return
		}
		ShowStepViewer(binImg, outImg)
	}

	// --- centered buttons row ---
	btnBox := container.NewHBox(
		btnUpload, btnRun, btnSave, btnStep,
	)

	// --- full-width bottom bar: info left, centered buttons right ---
	bottom := container.NewHBox(
		widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
			ShowInfoWindow(w)
		}),
		layout.NewSpacer(),
		btnBox,
	)

	grid := container.NewGridWithColumns(2,
		container.NewBorder(
			widget.NewLabel("Input"), nil, nil, nil,
			inIV,
		),
		container.NewBorder(
			widget.NewLabel("Output"), nil, nil, nil,
			outIV,
		),
	)

	w.SetContent(container.NewBorder(nil, bottom, nil, nil, grid))
	w.ShowAndRun()
	return nil
}
