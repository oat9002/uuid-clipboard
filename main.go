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

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	isExit := false

	for {
		uuid := uuid.New().String()
		clipboard.Write(clipboard.FmtText, []byte(uuid))
		fmt.Printf("Here: %s - auto save to clipboard\n", uuid)
		fmt.Println("Press Enter to generate new uuid")
		fmt.Println("Press any button to exit\n")

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
