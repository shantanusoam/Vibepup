package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type viewState int

const (
	stateSplash viewState = iota
	stateSetup
	stateRunning
	stateDone
)

type model struct {
	state     viewState
	start     time.Time
	frame     int
	args      []string
	form      *huh.Form
	selected  string
	spring    harmonica.Spring
	pos       float64
	velocity  float64
	targetPos float64
}

func initialModel() model {
	args := []string{"--watch"}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	spring := harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.2)

	m := model{state: stateSplash, start: time.Now(), args: args, selected: "watch", spring: spring, targetPos: 10}
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("♥ What shall we do today? ♥").
				Options(
					huh.NewOption("Watch (recommended)", "watch"),
					huh.NewOption("Run 5 iterations", "run"),
					huh.NewOption("New project from idea", "new"),
				).
				Value(&m.selected),
		),
	)

	return m
}

func (m model) Init() tea.Cmd {
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)
	log.SetTimeFormat("")
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg { return t })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case time.Time:
		switch m.state {
		case stateSplash:
			m.frame++
			if time.Since(m.start) > time.Second*2 {
				m.state = stateSetup
				return m, m.form.Init()
			}
			return m, tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg { return t })
		case stateRunning:
			m.frame++
			if m.frame%20 == 0 {
				if m.targetPos == 0 {
					m.targetPos = 10
				} else {
					m.targetPos = 0
				}
			}
			m.pos, m.velocity = m.spring.Update(m.pos, m.velocity, m.targetPos)
			return m, tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg { return t })
		}
	case tea.KeyMsg:
		if m.state == stateSetup {
			var cmd tea.Cmd
			fm, cmd := m.form.Update(msg)
			m.form = fm.(*huh.Form)
			if m.form.State == huh.StateCompleted {
				m.state = stateRunning
				return m, launchCmd(m.selected, m.args)
			}
			return m, cmd
		}
	}

	if m.state == stateSetup {
		fm, cmd := m.form.Update(msg)
		m.form = fm.(*huh.Form)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	// Theme Palette
	pink := lipgloss.Color("#FFB7C5")     // Sakura Pink
	hotPink := lipgloss.Color("#FF69B4")  // Hot Pink
	babyBlue := lipgloss.Color("#89CFF0") // Baby Blue
	lavender := lipgloss.Color("#E6E6FA") // Lavender
	gray := lipgloss.Color("245")         // Muted Gray

	// Layout Styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pink).
		Padding(1, 2).
		Margin(1, 1)

	titleStyle := lipgloss.NewStyle().
		Foreground(hotPink).
		Background(lavender).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	// Animation Frames
	frames := []string{
		"૮ ˶ᵔ ᵕ ᵔ˶ ა",        // Happy
		"૮ ˶• ﻌ •˶ ა",        // Alert
		"૮ ≧ ﻌ ≦ ა",          // Blink
		"૮ / ˶ • ﻌ • ˶ \\ ა", // Paws up
	}

	// Sparkle Animation
	sparkles := []string{"｡･ﾟ✧", "✧･ﾟ｡", "｡･ﾟ★", "☆･ﾟ｡"}

	// Frame Calculations
	// Slow down the dog animation (every 4th frame)
	idx := (m.frame / 4) % len(frames)
	// Fast sparkles (every 2nd frame)
	sIdx := (m.frame / 2) % len(sparkles)

	currentDog := frames[idx]
	currentSparkle := sparkles[sIdx]

	// Spring Animation (Horizontal movement)
	pad := ""
	if m.pos > 0 {
		pad = strings.Repeat(" ", int(m.pos))
	}

	// Common Elements
	dogRender := lipgloss.NewStyle().Foreground(pink).Render(pad + currentDog + " " + currentSparkle)

	// Splash Screen Content
	splashContent := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("♥ Vibepup TUI ♥"),
		"",
		dogRender,
		"",
		lipgloss.NewStyle().Foreground(gray).Render("Loading cuteness..."),
	)

	switch m.state {
	case stateSplash:
		return boxStyle.Render(splashContent)

	case stateSetup:
		return boxStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Foreground(babyBlue).Bold(true).Render("♥ Setup Phase"),
				"",
				m.form.View(),
			),
		)

	case stateRunning:
		status := lipgloss.NewStyle().Foreground(hotPink).Bold(true).Render("♥ Vibepup is Working ♥")
		mode := lipgloss.NewStyle().Foreground(babyBlue).Render("MODE: " + strings.ToUpper(m.selected))
		spinner := sparkles[m.frame%len(sparkles)]

		content := lipgloss.JoinVertical(lipgloss.Left,
			status,
			mode,
			"",
			lipgloss.NewStyle().Foreground(gray).Render(spinner+" Running engine... (check output below)"),
			"",
			dogRender,
		)
		return boxStyle.Render(content)

	case stateDone:
		return boxStyle.Render(
			lipgloss.NewStyle().Foreground(hotPink).Bold(true).Render("♥ All Done! Good Pup! ♥"),
		)

	default:
		return boxStyle.Render(splashContent)
	}
}

func launchCmd(choice string, args []string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if choice == "watch" {
			args = append([]string{"--watch"}, args...)
		}
		if choice == "run" {
			args = append([]string{"5"}, args...)
		}
		if choice == "new" {
			args = append([]string{"new", "A vibe-coded project"}, args...)
		}

		log.Info("Starting Vibepup", "args", strings.Join(args, " "))
		cmd := exec.CommandContext(ctx, "bash", "-lc", "vibepup "+strings.Join(args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		_ = cmd.Run()
		return tea.Quit()
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("failed to start vibepup tui")
		os.Exit(1)
	}
}
