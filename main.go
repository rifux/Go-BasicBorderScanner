package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/rifux/Go-BasicBorderScanner/internal/app/cli"
	"github.com/rifux/Go-BasicBorderScanner/internal/app/gui"
)

var (
	forceGUI = flag.Bool("gui", false, "force GUI mode")
	forceCLI = flag.Bool("cli", false, "force CLI mode")
	logMode  = flag.String("log", "auto", "log output mode: auto|json|text")
)

const helpText = `Go-BasicBorderScanner [https://rifux.dev]

USAGE:
  %s [flags]

GLOBAL FLAGS:
  -gui       force GUI mode (default on supported platforms)
  -cli       force CLI mode
  -log mode  log output mode: auto|json|text (default "auto")

CLI MODE FLAGS (only when -cli is given or GUI unavailable):
  -in file      input image (png, jpg, jpeg, gif, tiff, bmp)
  -out file     output file (default "out.png")
  -outfmt fmt   output format: png|jpeg|gif|tiff|bmp (default "png")
  -h, --help    show this help and exit

EXAMPLES:
  GUI:   %s
  CLI:   %s -cli -in scan.jpg -out result.png
`

func printHelp() {
	fmt.Printf(helpText, filepath.Base(os.Args[0]), os.Args[0], os.Args[0])
}

func hasGUI() bool {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "386", "amd64", "arm64":
			return true
		}
		return false
	case "darwin":
		switch runtime.GOARCH {
		case "amd64", "arm64":
			return true
		}
		return false
	case "linux":
		switch runtime.GOARCH {
		case "386", "amd64", "arm", "arm64", "loong64", "mips64le", "ppc64le", "riscv64", "s390x":
			return true
		}
		return false
	case "freebsd", "openbsd", "netbsd", "dragonfly":
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64" {
			return true
		}
		return false
	case "android", "ios":
		return true
	default:
		return false
	}
}

func main() {
	flag.Usage = printHelp
	flag.Parse()

	// show help for common help flags
	for _, a := range os.Args[1:] {
		switch a {
		case "-h", "--help", "help":
			printHelp()
			return
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var mode string
	switch {
	case *forceCLI:
		mode = "cli"
	case *forceGUI:
		mode = "gui"
	case hasGUI():
		mode = "gui"
	default:
		mode = "cli"
	}

	var err error
	switch mode {
	case "gui":
		err = gui.Run(ctx, *logMode)
	case "cli":
		err = cli.Run(ctx, *logMode)
	}
	if err != nil {
		log.Fatal(err.Error())
	}
}
