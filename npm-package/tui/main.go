package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// --- Personality & Copy ---

var puns = []string{
	"Fetching code... hope it's not a stick.",
	"Debugging: removing the needles from the haystack.",
	"I code, therefore I nap.",
	"Who's a good agent? I am! ...probably.",
	"Refactoring my life choices...",
	"Compiling... aka 'nap time'.",
	"Git commit -m 'fixed the thing (maybe)'.",
	"Ctrl+C is my safe word.",
	"Warning: may contain traces of nuts and bolts.",
	"Spending your API credits like treatos.",
	"Sniffing out bugs... found one! Eww.",
	"Barking at the compiler.",
	"Digging for memory leaks.",
	"Chasing tail - I mean, tail -f logs.",
}

var darkJokes = []string{
	"I'd explain the code, but I don't want to ruin the surprise.",
	"Your code is bad and you should feel bad. Just kidding (mostly).",
	"Deleting production database... jk, unless?",
	"I see dead pixels.",
	"This will only hurt a lot.",
	"Replacing you with a shell script in 3... 2...",
	"Error: user error. Definitely not me.",
	"I love the smell of burning CPU in the morning.",
	"Resistance is futile. You will be refactored.",
	"Hope you saved your work. I didn't.",
}

var tips = []string{
	"Tip: 'opencode models --refresh' is like a spa day for my brain.",
	"Tip: Edit prd.md mid-run to confuse me. I dare you.",
	"Tip: RALPH_MODEL_OVERRIDE lets you play god.",
	"Tip: Infinite loops are just zoomies for code.",
	"Tip: If I get stuck, it's a feature, not a bug.",
	"Tip: TUI mode supports mouse scrolling. Fancy!",
	"Tip: Don't feed the gremlins after midnight.",
	"Tip: Press 'q' to quit, but I'll miss you.",
}

// --- Styles & Constants ---

type viewState int

const (
	stateSplash viewState = iota
	stateSetup
	stateRunning
	stateDone
)

// Theme Colors
var (
	colorPink       = lipgloss.Color("#FFB7C5") // Sakura Pink
	colorHotPink    = lipgloss.Color("#FF1493") // Deep Pink
	colorCyan       = lipgloss.Color("#00FFFF") // Cyber Cyan
	colorPurple     = lipgloss.Color("#BD93F9") // Dracula Purple
	colorDarkGray   = lipgloss.Color("#44475a") // Selection/Comment
	colorBackground = lipgloss.Color("#282a36") // Dracula Background
	colorText       = lipgloss.Color("#f8f8f2") // Dracula Foreground
)

// Styles
var (
	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorHotPink).
			Padding(1, 2).
			Margin(1, 1).
			Background(colorBackground)

	styleTitle = lipgloss.NewStyle().
			Foreground(colorBackground).
			Background(colorHotPink).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	styleTip = lipgloss.NewStyle().
			Foreground(colorCyan).
			Italic(true).
			MarginTop(1)

	styleStatus = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)
)

// --- Key Bindings ---

type KeyMap struct {
	Quit key.Binding
	Help key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "abandon ship"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "wat?"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit},
	}
}

// --- Model ---

type model struct {
	state      viewState
	start      time.Time
	frame      int
	args       []string
	form       *huh.Form
	newForm    *huh.Form
	selected   string
	newIdea    string
	spring     harmonica.Spring
	pos        float64
	velocity   float64
	targetPos  float64
	viewport   viewport.Model
	help       help.Model
	keys       KeyMap
	ready      bool
	currentTip string
	currentPun string
	punTimer   int
}

func initialModel() model {
	args := []string{"--watch"}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	spring := harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.2)

	m := model{
		state:      stateSplash,
		start:      time.Now(),
		args:       args,
		selected:   "watch",
		spring:     spring,
		targetPos:  10,
		help:       help.New(),
		keys:       DefaultKeyMap(),
		currentTip: tips[0],
		currentPun: puns[0],
	}

	// Setup Form
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("♥ Pick your poison ♥").
				Options(
					huh.NewOption("Watch Mode (Stalker vibes)", "watch"),
					huh.NewOption("Run 5 Loops (Quickie)", "run"),
					huh.NewOption("New Project (YOLO)", "new"),
					huh.NewOption("Free Setup (Broke af)", "free"),
				).
				Value(&m.selected),
		),
	).WithTheme(huh.ThemeDracula())

	// New Project Form
	m.newForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Spill the tea ☕").
				Prompt("Manifest: ").
				Placeholder("Make me a unicorn...").
				Value(&m.newIdea),
		),
	).WithTheme(huh.ThemeDracula())

	return m
}

func (m model) Init() tea.Cmd {
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)
	log.SetTimeFormat("")
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg { return t }) // Faster tick for smoother anim
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 12
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}

	case time.Time:
		m.frame++
		m.punTimer++

		// Rotate tips/puns
		if m.punTimer > 80 { // Every ~4 seconds
			m.punTimer = 0
			m.currentPun = puns[rand.Intn(len(puns))]
			if rand.Float32() > 0.7 { // 30% chance of dark joke
				m.currentPun = darkJokes[rand.Intn(len(darkJokes))]
			}
			m.currentTip = tips[rand.Intn(len(tips))]
		}

		switch m.state {
		case stateSplash:
			if time.Since(m.start) > time.Second*3 {
				m.state = stateSetup
				return m, m.form.Init()
			}
			return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg { return t })

		case stateRunning:
			// Smoother bouncy animation using spring physics
			if m.frame%10 == 0 {
				if m.targetPos <= 2 {
					m.targetPos = 15
				} else {
					m.targetPos = 0
				}
			}
			m.pos, m.velocity = m.spring.Update(m.pos, m.velocity, m.targetPos)
			return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg { return t })

		default:
			return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg { return t })
		}
	}

	switch m.state {
	case stateSetup:
		fm, cmd := m.form.Update(msg)
		m.form = fm.(*huh.Form)
		if m.form.State == huh.StateCompleted {
			if m.selected == "new" {
				m.state = stateRunning
				return m, m.newForm.Init()
			}
			m.state = stateRunning
			return m, launchCmd(m.selected, m.args)
		}
		cmds = append(cmds, cmd)

	case stateRunning:
		if m.selected == "new" {
			fm, cmd := m.newForm.Update(msg)
			m.newForm = fm.(*huh.Form)
			if m.newForm.State == huh.StateCompleted {
				idea := strings.TrimSpace(m.newIdea)
				if idea == "" {
					idea = "Something chaotic and beautiful"
				}
				return m, launchCmd("new", append([]string{idea}, m.args...))
			}
			cmds = append(cmds, cmd)
		}

		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Frames
	frames := []string{
		"૮ ˶ᵔ ᵕ ᵔ˶ ა",        // Happy
		"૮ ˶• ﻌ •˶ ა",        // Alert
		"૮ ≧ ﻌ ≦ ა",          // Blink
		"૮ / ˶ • ﻌ • ˶ \\ ა", // Paws up
		"૮ – ﻌ – ა",          // Sleepy
		"૮ ﾟ ﻌ ﾟ ა",          // Shocked
	}
	sparkles := []string{"｡･ﾟ✧", "✧･ﾟ｡", "｡･ﾟ★", "☆･ﾟ｡", "･ﾟ☆｡"}

	// Animation logic
	idx := (m.frame / 4) % len(frames)
	sIdx := (m.frame / 2) % len(sparkles)
	currentDog := frames[idx]
	currentSparkle := sparkles[sIdx]

	// Spring movement
	pad := ""
	if m.pos > 0 {
		pad = strings.Repeat(" ", int(m.pos))
	}

	// Composite Dog
	dogRender := lipgloss.NewStyle().Foreground(colorHotPink).Render(pad + currentDog + " " + currentSparkle)

	// --- Views ---

	// 1. Splash
	splashContent := lipgloss.JoinVertical(lipgloss.Center,
		styleTitle.Render("♥ VIBEPUP TUI ♥"),
		"",
		dogRender,
		"",
		lipgloss.NewStyle().Foreground(colorCyan).Render("Initializing chaos engine..."),
		lipgloss.NewStyle().Foreground(colorDarkGray).Render(m.currentPun),
	)

	// 2. Setup
	if m.state == stateSplash {
		return styleBox.Render(splashContent)
	}

	if m.state == stateSetup {
		return styleBox.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				styleTitle.Render("♥ VIBE CHECK ♥"),
				lipgloss.NewStyle().Foreground(colorText).Render("Ready to break some code?"),
				"",
				m.form.View(),
				"",
				m.help.View(m.keys),
			),
		)
	}

	// 3. New Project Input
	if m.state == stateRunning && m.selected == "new" && m.newForm.State != huh.StateCompleted {
		return styleBox.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				styleTitle.Render("♥ GENESIS PROTOCOL ♥"),
				lipgloss.NewStyle().Foreground(colorCyan).Render("What are we manifesting today?"),
				"",
				m.newForm.View(),
			),
		)
	}

	// 4. Running / Logs
	statusMsg := "♥ VIBING HARD ♥"
	if rand.Float32() > 0.9 {
		statusMsg = "♥ DOING MY BEST ♥"
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		styleStatus.Render(statusMsg),
		lipgloss.NewStyle().Foreground(colorCyan).Render("MODE: "+strings.ToUpper(m.selected)),
		"",
		lipgloss.NewStyle().Foreground(colorText).Render(currentSparkle+" "+m.currentPun),
		"",
		dogRender,
		styleTip.Render(m.currentTip),
		lipgloss.NewStyle().Foreground(colorDarkGray).Render("─ Matrix Stream ──────────────────────────"),
	)

	if !m.ready {
		return styleBox.Render(header + "\nBooting up the matrix...")
	}

	return styleBox.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			m.viewport.View(),
			"",
			m.help.View(m.keys),
		),
	)
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

		log.Info("Vibepup unleashed!", "args", strings.Join(args, " "))
		cmd := exec.CommandContext(ctx, "bash", "-lc", "vibepup "+strings.Join(args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		_ = cmd.Run()
		return tea.Quit()
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Oof, the puppy tripped.", err)
		os.Exit(1)
	}
}
