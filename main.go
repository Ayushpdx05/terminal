package main

import (
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

func main() {
	a := app.New()
	w := a.NewWindow("germ")

	ui := widget.NewTextGrid()       // Create a new TextGrid
	ui.SetText("I'm on a terminal!") // Set initial text

	c := exec.Command("/bin/bash")
	p, err := pty.Start(c) // Start the pty with the bash command
	if err != nil {
		fyne.LogError("Failed to open pty", err)
		os.Exit(1)
	}
	defer p.Close()
	defer c.Process.Kill()

	// Handle keypresses and input from the terminal user
	onTypedKey := func(e *fyne.KeyEvent) {
		if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
			_, _ = p.Write([]byte{'\r'}) // Send carriage return on Enter
		}
	}

	onTypedRune := func(r rune) {
		_, _ = p.WriteString(string(r)) // Write characters typed to the pty
	}

	// Assign the key and rune handlers to the window's canvas
	w.Canvas().SetOnTypedKey(onTypedKey)
	w.Canvas().SetOnTypedRune(onTypedRune)

	// Goroutine for continuously reading from the pty
	go func() {
		for {
			time.Sleep(100 * time.Millisecond) // Adjust delay to reduce latency
			b := make([]byte, 256)             // Create a buffer for reading pty output
			n, err := p.Read(b)                // Read output from pty
			if err != nil {
				fyne.LogError("Failed to read pty", err)
				return
			}

			// Set the text in the terminal window with the output
			ui.SetText(string(b[:n])) // Update TextGrid with new output
		}
	}()

	// Create a container for the terminal UI
	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewGridWrapLayout(fyne.NewSize(420, 200)), // Fixed terminal size
			ui,
		),
	)

	w.ShowAndRun()
}
