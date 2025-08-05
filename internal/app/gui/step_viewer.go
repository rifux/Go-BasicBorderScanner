package gui

import (
	"context"
	"image"
	"image/draw"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ---- step viewer for detailed showcase ----
func ShowStepViewer(binImg, scanImg image.Image) {
	app := fyne.CurrentApp()
	w := app.NewWindow("Step viewer")
	w.Resize(fyne.NewSize(700, 550))

	bounds := binImg.Bounds()
	h := bounds.Max.Y

	base := image.NewRGBA(bounds)
	draw.Draw(base, bounds, binImg, bounds.Min, draw.Src)

	imgCanvas := canvas.NewImageFromImage(base)
	imgCanvas.FillMode = canvas.ImageFillContain
	imgCanvas.SetMinSize(fyne.NewSize(400, 400))

	info := widget.NewLabel("Lines shown: 0")
	slider := widget.NewSlider(0, float64(h))
	slider.Step = 1
	slider.ExtendBaseWidget(slider)

	const bufSize = 2
	workCh := make(chan int, bufSize)

	go func() {
		var ctx context.Context
		var cancel context.CancelFunc

		defer func() {
			if cancel != nil {
				cancel()
			}
		}()

		for n := range workCh {
			if cancel != nil {
				cancel()
			}
			ctx, cancel = context.WithCancel(context.Background())
			go func(target int) {
				defer cancel() // cancel by exit
				select {
				case <-ctx.Done():
					return
				default:
					tmp := image.NewRGBA(bounds)
					draw.Draw(tmp, bounds, binImg, bounds.Min, draw.Src)
					if target > 0 {
						r := image.Rect(bounds.Min.X, bounds.Min.Y,
							bounds.Max.X, bounds.Min.Y+target)
						draw.Draw(tmp, r, scanImg, image.Point{}, draw.Over)
					}

					fyne.Do(func() {
						draw.Draw(base, bounds, tmp, bounds.Min, draw.Src)
						imgCanvas.Refresh()
						info.SetText("Lines shown: " + strconv.Itoa(target))
					})
				}
			}(n)
		}
	}()

	slider.OnChangeEnded = func(v float64) {
		n := int(v)
		if len(workCh) == bufSize {
			<-workCh
		}
		workCh <- n
	}

	exitBtn := widget.NewButton("Exit", func() {
		close(workCh)
		w.Close()
	})

	w.SetContent(
		container.NewBorder(
			container.NewHBox(info, layout.NewSpacer(), exitBtn),
			container.NewVBox(slider),
			nil, nil,
			imgCanvas,
		),
	)
	slider.SetValue(float64(h))
	w.Show()
}
