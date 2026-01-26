package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"vibepup-tui/theme"
)

type StatusBar struct {
	Theme theme.Theme
}

func (s StatusBar) Render(left, right string, width int) string {
	style := lipgloss.NewStyle().Foreground(s.Theme.Foreground).Background(s.Theme.Muted)
	spacer := width - lipgloss.Width(left) - lipgloss.Width(right)
	if spacer < 1 {
		spacer = 1
	}
	return style.Render(fmt.Sprintf("%s%s%s", left, spaces(spacer), right))
}

func spaces(n int) string {
	if n < 0 {
		n = 0
	}
	return fmt.Sprintf("%*s", n, "")
}
