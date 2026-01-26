package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Layout struct {
	HeaderHeight int
	FooterHeight int
	Width        int
	Height       int
}

type LogViewport struct {
	Model      viewport.Model
	Content    *strings.Builder
	AutoScroll bool
}

func NewLogViewport(width, height int) LogViewport {
	vp := viewport.New(width, height)
	vp.YPosition = 0
	vp.HighPerformanceRendering = false
	return LogViewport{
		Model:      vp,
		Content:    &strings.Builder{},
		AutoScroll: true,
	}
}

func (l *LogViewport) SetSize(width, height int) {
	l.Model.Width = width
	l.Model.Height = height
}

func (l *LogViewport) WriteLine(line string) {
	l.Content.WriteString(line + "\n")
	l.Model.SetContent(l.Content.String())
	if l.AutoScroll {
		l.Model.GotoBottom()
	}
}

func (l *LogViewport) Update(msg tea.Msg) (LogViewport, tea.Cmd) {
	var cmd tea.Cmd
	l.Model, cmd = l.Model.Update(msg)
	return *l, cmd
}

func (l *LogViewport) View() string {
	return l.Model.View()
}

var (
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Margin(0, 1) // Reduced margin to save space
)
