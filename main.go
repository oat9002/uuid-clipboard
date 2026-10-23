package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"golang.design/x/clipboard"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	run(screen, clipboard.Init() == nil, clipboard.Write)
}

func run(screen tcell.Screen, isClipboardInit bool, write func(clipboard.Format, []byte) <-chan struct{}) {
	uuid, msg := generateUuidAndCopyToClipboard(isClipboardInit, write)

	for {
		screen.Clear()

		message := fmt.Sprintf("Here: %s - %s\nPress Enter to generate new UUID\nPress any other key to exit", uuid, msg)
		printMessage(screen, message)

		screen.Show()

		switch ev := screen.PollEvent(); ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.(*tcell.EventKey).Key() {
			case tcell.KeyEnter:
				uuid, msg = generateUuidAndCopyToClipboard(isClipboardInit, write)
			default:
				return
			}
		}
	}
}

func printMessage(screen tcell.Screen, message string) {
	lines := []rune(message)
	x, y := 0, 0
	for _, r := range lines {
		if r == '\n' {
			x = 0
			y++
			continue
		}
		screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		x++
	}
}

func generateUuidAndCopyToClipboard(isClipboardInit bool, write func(clipboard.Format, []byte) <-chan struct{}) (string, string) {
	uuid := uuid.New().String()

	if !isClipboardInit {
		return uuid, "You need to manually copy the uuid"
	}

	if write(clipboard.FmtText, []byte(uuid)) == nil {
		return uuid, "Failed to copy to clipboard; manually copy the uuid"
	}

	return uuid, "auto save to clipboard"
}
