package main

import (
	"fmt"

	"github.com/google/uuid"
	term "github.com/nsf/termbox-go"
	"golang.design/x/clipboard"
)

func main() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()

	isClipboardInit := true
	err = clipboard.Init()
	if err != nil {
		isClipboardInit = false
	}

	msg := "You need to manually copy the uuid"
	isExit := false

	if isClipboardInit {
		msg = "auto save to clipboard"
	}

	for {
		uuid := uuid.New().String()

		if isClipboardInit {
			clipboard.Write(clipboard.FmtText, []byte(uuid))
		}

		fmt.Printf("Here: %s - %s\n", uuid, msg)
		fmt.Println("Press Enter to generate new uuid")
		fmt.Println("Press any button to exit")

		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEnter:
				term.Sync()
				isExit = false

			default:
				term.Sync()
				isExit = true
			}
		case term.EventError:
			panic(ev.Err)
		}

		if isExit {
			break
		}
	}
}
