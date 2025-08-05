package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rifux/Go-BasicBorderScanner/internal/app/cli"
	"github.com/rifux/Go-BasicBorderScanner/internal/app/gui"
)

const helpText = `Go-BasicBorderScanner [https://rifux.dev]

A simple tool to find and draw the contour of an object in an image.

USAGE:
  %s <command> [flags]

COMMANDS:
  gui     Run the application in graphical user interface mode.
  cli     Run the application in command-line interface mode.
  help    Show this help message.

Run '%s <command> --help' for more information on a command.
`

func printHelp() {
	p := filepath.Base(os.Args[0])
	fmt.Printf(helpText, p, p)
}

func hasGUI() bool {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "386", "amd64", "arm64":
			return true
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64", "arm64":
			return true
		}
	case "linux":
		switch runtime.GOARCH {
		case "386", "amd64", "arm", "arm64", "loong64", "mips64le", "ppc64le", "riscv64", "s390x":
			return true
		}
	case "freebsd", "openbsd", "netbsd", "dragonfly":
		if runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64" {
			return true
		}
	}
	return false
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var args []string
	var mode string

	if len(os.Args) < 2 {
		if hasGUI() {
			mode = "gui"
			args = os.Args[1:] // Pass all args to gui
		} else {
			printHelp()
			return
		}
	} else {
		mode = os.Args[1]
		args = os.Args[2:] // Args for the command
	}

	var err error
	switch strings.ReplaceAll(mode, "-", "") {
	case "gui":
		// Assuming gui.Run is updated to parse its own flags from its args
		err = gui.Run(ctx, "")
	case "cli":
		err = cli.Run(ctx, args)
	case "help", "h":
		printHelp()
	default:
		// Default to GUI if the first arg isn't a known command but GUI is available
		if hasGUI() {
			err = gui.Run(ctx, "")
		} else {
			fmt.Printf("Error: Unknown command %q\n\n", mode)
			printHelp()
			os.Exit(1)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}
