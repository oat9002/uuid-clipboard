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

	isClipboardInit := true
	err = clipboard.Init()
	if err != nil {
		isClipboardInit = false
	}

	msg := "You need to manually copy the uuid"

	if isClipboardInit {
		msg = "auto save to clipboard"
	}

	for {
		screen.Clear()
		uuid := uuid.New().String()

		if isClipboardInit {
			clipboard.Write(clipboard.FmtText, []byte(uuid))
		}

		message := fmt.Sprintf("Here: %s - %s\nPress Enter to generate new UUID\nPress any other key to exit", uuid, msg)
		printMessage(screen, message)

		screen.Show()

		switch ev := screen.PollEvent(); ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch ev.(*tcell.EventKey).Key() {
			case tcell.KeyEnter:
				fmt.Println("Generating new UUID...")
				screen.Sync()
			default:
				break
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
