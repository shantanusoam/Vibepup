package main

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"vibepup-tui/config"
)

func waitForOutput(t *testing.T, tm *teatest.TestModel, needle []byte) {
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if bytes.Contains(tm.Output(), needle) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("expected output to contain %q, got: %s", needle, string(tm.Output()))
}

func TestTUITransitionsToSetup(t *testing.T) {
	flags := config.Flags{ForceRun: true}
	m := initialModel(flags)

	tm := teatest.NewTestModel(
		t,
		m,
		teatest.WithInitialTermSize(80, 24),
	)
	defer tm.Quit()

	tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})

	waitForOutput(t, tm, []byte("SETUP"))
}

func TestTUIHelpToggle(t *testing.T) {
	flags := config.Flags{ForceRun: true}
	m := initialModel(flags)

	tm := teatest.NewTestModel(
		t,
		m,
		teatest.WithInitialTermSize(80, 24),
	)
	defer tm.Quit()

	tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})
	waitForOutput(t, tm, []byte("SETUP"))

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

	waitForOutput(t, tm, []byte("help"))
}
