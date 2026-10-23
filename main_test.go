package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"golang.design/x/clipboard"
)

func TestGenerateUuidAndCopyToClipboard(t *testing.T) {
	for _, tt := range []struct {
		name    string
		enabled bool
		changed <-chan struct{}
		wantMsg string
	}{
		{"clipboard disabled", false, nil, "You need to manually copy the uuid"},
		{"copy succeeds", true, make(chan struct{}), "auto save to clipboard"},
		{"copy fails", true, nil, "Failed to copy to clipboard; manually copy the uuid"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var calls int
			var copied string
			write := func(format clipboard.Format, data []byte) <-chan struct{} {
				calls++
				if format != clipboard.FmtText {
					t.Errorf("clipboard format = %v, want text", format)
				}
				copied = string(data)
				return tt.changed
			}

			u, msg := generateUuidAndCopyToClipboard(tt.enabled, write)
			parsed, err := uuid.Parse(u)
			if err != nil {
				t.Fatalf("invalid UUID %q: %v", u, err)
			}
			if parsed.Version() != 4 || parsed.Variant() != uuid.RFC4122 {
				t.Errorf("expected an RFC 4122 version 4 UUID, got %s", u)
			}
			if msg != tt.wantMsg {
				t.Errorf("status = %q, want %q", msg, tt.wantMsg)
			}
			if tt.enabled {
				if calls != 1 || copied != u {
					t.Errorf("clipboard writes = %d, text = %q; want one write of %q", calls, copied, u)
				}
			} else if calls != 0 {
				t.Errorf("clipboard disabled but write called %d times", calls)
			}
		})
	}
}

func TestGeneratedUUIDsAreUnique(t *testing.T) {
	u1, _ := generateUuidAndCopyToClipboard(false, nil)
	u2, _ := generateUuidAndCopyToClipboard(false, nil)
	if u1 == u2 {
		t.Fatalf("expected different UUIDs, got identical: %s", u1)
	}
}

func TestPrintMessage(t *testing.T) {
	for _, tt := range []struct {
		name    string
		message string
		want    []string
	}{
		{"empty message", "", []string{"", "", "", ""}},
		{"single line", "Here: UUID", []string{"Here: UUID", "", "", ""}},
		{"instructions on separate lines", "Here: UUID\nPress Enter\nOther keys exit", []string{"Here: UUID", "Press Enter", "Other keys exit", ""}},
		{"blank line", "First\n\nThird\n", []string{"First", "", "Third", ""}},
		{"unicode text", "café", []string{"café", "", "", ""}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			screen := tcell.NewSimulationScreen("")
			if err := screen.Init(); err != nil {
				t.Fatal(err)
			}
			defer screen.Fini()
			screen.SetSize(40, len(tt.want))
			printMessage(screen, tt.message)
			screen.Show()

			cells, width, _ := screen.GetContents()
			for y, want := range tt.want {
				var line strings.Builder
				for _, cell := range cells[y*width : (y+1)*width] {
					line.WriteString(string(cell.Runes))
				}
				if got := strings.TrimRight(line.String(), " "); got != want {
					t.Errorf("row %d = %q, want %q", y, got, want)
				}
			}
		})
	}
}

func TestRunWithClipboardDisabled(t *testing.T) {
	sim := tcell.NewSimulationScreen("")
	if err := sim.Init(); err != nil {
		t.Fatal(err)
	}
	defer sim.Fini()
	sim.SetSize(160, 5)
	screen := &recordingScreen{Screen: sim}
	for _, event := range []tcell.Event{
		tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone),
	} {
		if err := sim.PostEvent(event); err != nil {
			t.Fatal(err)
		}
	}

	run(screen, false, func(clipboard.Format, []byte) <-chan struct{} {
		t.Fatal("must not write to an unavailable clipboard")
		return nil
	})

	if len(screen.frames) != 2 {
		t.Fatalf("rendered frames = %d, want 2 (initial and Enter, then exit on q)", len(screen.frames))
	}
	var previous string
	for i, frame := range screen.frames {
		value, hasPrefix := strings.CutPrefix(frame, "Here: ")
		value, hasSuffix := strings.CutSuffix(value, " - You need to manually copy the uuid")
		if !hasPrefix || !hasSuffix {
			t.Fatalf("frame %d is missing UUID display or manual copy instructions: %q", i, frame)
		}
		if _, err := uuid.Parse(value); err != nil {
			t.Fatalf("frame %d has invalid UUID %q: %v", i, value, err)
		}
		if value == previous {
			t.Error("Enter must generate a new UUID even when the clipboard is unavailable")
		}
		previous = value
	}
}

func TestRunClipboardStatusSurvivesRedrawAndRecovers(t *testing.T) {
	sim := tcell.NewSimulationScreen("")
	if err := sim.Init(); err != nil {
		t.Fatal(err)
	}
	defer sim.Fini()
	sim.SetSize(160, 5)
	screen := &recordingScreen{Screen: sim}

	for _, event := range []tcell.Event{
		tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
		tcell.NewEventResize(160, 5),
		tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
		tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone),
	} {
		if err := sim.PostEvent(event); err != nil {
			t.Fatal(err)
		}
	}

	var copied []string
	write := func(_ clipboard.Format, data []byte) <-chan struct{} {
		copied = append(copied, string(data))
		if len(copied) == 2 {
			return nil
		}
		return make(chan struct{})
	}
	run(screen, true, write)

	if len(copied) != 3 {
		t.Fatalf("clipboard writes = %d, want 3 (initial and two Enter presses)", len(copied))
	}
	if copied[0] == copied[1] || copied[1] == copied[2] || copied[0] == copied[2] {
		t.Fatalf("Enter did not generate new UUIDs: %v", copied)
	}
	wantFrames := []string{
		fmt.Sprintf("Here: %s - auto save to clipboard", copied[0]),
		fmt.Sprintf("Here: %s - Failed to copy to clipboard; manually copy the uuid", copied[1]),
		fmt.Sprintf("Here: %s - Failed to copy to clipboard; manually copy the uuid", copied[1]),
		fmt.Sprintf("Here: %s - auto save to clipboard", copied[2]),
	}
	if len(screen.frames) != len(wantFrames) {
		t.Fatalf("rendered frames = %d, want %d", len(screen.frames), len(wantFrames))
	}
	for i, want := range wantFrames {
		if screen.frames[i] != want {
			t.Errorf("frame %d = %q, want %q", i, screen.frames[i], want)
		}
	}
}

type recordingScreen struct {
	tcell.Screen
	frames []string
}

func (s *recordingScreen) Show() {
	s.Screen.Show()
	width, _ := s.Size()
	var line strings.Builder
	for x := 0; x < width; x++ {
		r, _, _, _ := s.GetContent(x, 0)
		line.WriteRune(r)
	}
	s.frames = append(s.frames, strings.TrimRight(line.String(), " "))
}
