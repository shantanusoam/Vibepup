package motion

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is emitted on each animation frame.
type TickMsg time.Time

// Engine manages frame ticks respecting perf flags.
type Engine struct {
	interval time.Duration
}

func New(perfLow, quiet bool) Engine {
	interval := time.Millisecond * 50
	if perfLow || quiet {
		interval = time.Millisecond * 80
	}
	return Engine{interval: interval}
}

func (e Engine) Next() tea.Cmd {
	return tea.Tick(e.interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
