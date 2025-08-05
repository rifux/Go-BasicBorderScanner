package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func richtext(text string) fyne.Widget {
	richText := widget.NewRichTextFromMarkdown(text)
	richText.Wrapping = fyne.TextWrapWord
	return richText
}

// ---- app info dialogue ----
func ShowInfoWindow(parent fyne.Window) {
	osInfo := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	// Create a markdown strings with hyperlinks
	aboutText := `# Go-BasicBorderScanner 🚀
A simple tool for detecting and drawing borders on images.`
	licensingText := `### License & Author
* **License:** [Apache-2.0](https://github.com/rifux/Go-BasicBorderScanner/blob/release/LICENSE)
* **Author:** [rifux 🌸 Vladimir Blinkov](https://rifux.dev)
* **GitHub:** [https://github.com/rifux/Go-BasicBorderScanner](https://github.com/rifux/Go-BasicBorderScanner)`
	usageText := fmt.Sprintf(`### CLI Usage

	%s --help
`, filepath.Base(os.Args[0]))
	sysinfoText := fmt.Sprintf(`Running as: **%s**`, osInfo)

	// Create a vertical container box
	content := container.NewVBox(
		richtext(aboutText),
		richtext(licensingText),
		richtext(usageText),
		richtext(sysinfoText),
	)

	// Create the custom dialog object
	d := dialog.NewCustom(
		"About",
		"Close",
		content,
		parent,
	)

	// Resize the dialog itself
	d.Resize(fyne.NewSize(450, 220))

	// Show the resized dialog
	d.Show()
}
