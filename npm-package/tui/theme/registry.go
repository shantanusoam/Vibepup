package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name          string
	Background    lipgloss.Color
	Foreground    lipgloss.Color
	Accent        lipgloss.Color
	AccentAlt     lipgloss.Color
	Muted         lipgloss.Color
	Border        lipgloss.Color
	Highlight     lipgloss.Color
	SupportsEmoji bool
}

var themes = map[string]Theme{}

func init() {
	RegisterTheme(Theme{
		Name:          "dracula-vibe",
		Background:    lipgloss.Color("#282a36"),
		Foreground:    lipgloss.Color("#f8f8f2"),
		Accent:        lipgloss.Color("#FF1493"),
		AccentAlt:     lipgloss.Color("#00FFFF"),
		Muted:         lipgloss.Color("#44475a"),
		Border:        lipgloss.Color("#FF1493"),
		Highlight:     lipgloss.Color("#BD93F9"),
		SupportsEmoji: true,
	})

	RegisterTheme(Theme{
		Name:          "halloween-glitch",
		Background:    lipgloss.Color("#0d0b1a"),
		Foreground:    lipgloss.Color("#f8e7cf"),
		Accent:        lipgloss.Color("#ff6b00"),
		AccentAlt:     lipgloss.Color("#6df3ff"),
		Muted:         lipgloss.Color("#2b203d"),
		Border:        lipgloss.Color("#ff00aa"),
		Highlight:     lipgloss.Color("#aaff00"),
		SupportsEmoji: true,
	})

	RegisterTheme(Theme{
		Name:          "mono-chill",
		Background:    lipgloss.Color("#101010"),
		Foreground:    lipgloss.Color("#e6e6e6"),
		Accent:        lipgloss.Color("#8ec07c"),
		AccentAlt:     lipgloss.Color("#83a598"),
		Muted:         lipgloss.Color("#3c3836"),
		Border:        lipgloss.Color("#928374"),
		Highlight:     lipgloss.Color("#fabd2f"),
		SupportsEmoji: false,
	})
}

func RegisterTheme(t Theme) {
	themes[t.Name] = t
}

func Get(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["dracula-vibe"]
}

func All() []Theme {
	items := make([]Theme, 0, len(themes))
	for _, t := range themes {
		items = append(items, t)
	}
	return items
}
