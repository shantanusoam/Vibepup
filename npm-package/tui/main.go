package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"

	"vibepup-tui/config"
	"vibepup-tui/motion"
	"vibepup-tui/persona"
	"vibepup-tui/process"
	"vibepup-tui/theme"
	"vibepup-tui/ui"
)

// --- Key Bindings ---

type KeyMap struct {
	Quit      key.Binding
	Help      key.Binding
	NextTheme key.Binding
	Pet       key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		NextTheme: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "theme")),
		Pet: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "pet dog")),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.NextTheme, k.Pet}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Help, k.Quit, k.NextTheme, k.Pet}}
}

// --- Model ---

type viewState int

const (
	stateSplash viewState = iota
	stateSetup
	stateRunning
	stateDone
)

type model struct {
	state      viewState
	width      int
	height     int
	ready      bool
	
	// Components
	keys       KeyMap
	help       help.Model
	form       *huh.Form
	newForm    *huh.Form
	viewport   ui.LogViewport
	spinner    spinner.Model
	
	// Config & State
	flags      config.Flags
	theme      theme.Theme
	styles     lipgloss.Style // Simplified, use theme package directly where possible
	snark      persona.SnarkLevel
	
	// Process
	runner     *process.Runner
	selected   string
	newIdea    string
	args       []string
	
	// Animation
	motion     motion.Engine
	frame      int
	dogState   string // "sleeping", "running", "barking", "happy"
}

func initialModel(flags config.Flags) model {
	th := theme.Get(flags.Theme)
	snark := persona.ParseSnark(flags.Snark)
	
	s := spinner.New()
	s.Spinner = spinner.Points // More modern spinner
	s.Style = lipgloss.NewStyle().Foreground(th.Highlight)

	m := model{
		state:    stateSplash,
		keys:     DefaultKeyMap(),
		help:     help.New(),
		flags:    flags,
		theme:    th,
		snark:    snark,
		motion:   motion.New(flags.PerfLow, flags.Quiet),
		spinner:  s,
		dogState: "sleeping",
		selected: "watch", // Default
		args:     os.Args[1:],
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
	return tea.Batch(
		m.motion.Next(),
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		
		// Dynamic Layout Calculation
		headerHeight := 10 // Approximation, should be measured
		footerHeight := 3
		vpHeight := msg.Height - headerHeight - footerHeight - 2 // Borders
		
		if !m.ready {
			m.viewport = ui.NewLogViewport(msg.Width-4, vpHeight)
			m.ready = true
		} else {
			m.viewport.SetSize(msg.Width-4, vpHeight)
		}

	case tea.KeyMsg:
		if m.state == stateRunning && m.runner != nil {
			// Pass interactions to viewport if needed
		}
		
		switch {
		case key.Matches(msg, m.keys.Quit):
			if m.runner != nil {
				m.runner.Kill() // ZOMBIE KILLER
			}
			return m, tea.Quit
		case key.Matches(msg, m.keys.Pet):
			m.dogState = "happy"
			cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return "dog_reset"
			}))
		}

	case string:
		if msg == "dog_reset" {
			if m.runner != nil {
				m.dogState = "running"
			} else {
				m.dogState = "sleeping"
			}
		}

	case motion.TickMsg:
		m.frame++
		cmds = append(cmds, m.motion.Next())

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case process.OutputMsg:
		m.viewport.WriteLine(string(msg))
		cmds = append(cmds, m.runner.WaitForOutput())

	case process.DoneMsg:
		m.dogState = "sleeping"
		m.viewport.WriteLine("\n--- Process Finished ---")
		if msg.Err != nil {
			m.viewport.WriteLine(fmt.Sprintf("Error: %v", msg.Err))
			m.dogState = "barking"
		}
		m.runner = nil
	}

	// Handle Forms
	if m.state == stateSetup {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			if m.form.State == huh.StateCompleted {
				if m.selected == "new" {
					m.state = stateRunning
					cmds = append(cmds, m.newForm.Init())
				} else {
					m.state = stateRunning
					cmds = append(cmds, m.startProcess())
				}
			}
		}
		cmds = append(cmds, cmd)
	}

	if m.state == stateRunning && m.selected == "new" && m.runner == nil {
		form, cmd := m.newForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.newForm = f
			if m.newForm.State == huh.StateCompleted {
				m.args = append(m.args, "new", m.newIdea)
				cmds = append(cmds, m.startProcess())
			}
		}
		cmds = append(cmds, cmd)
	}
	
	// Auto-advance splash
	if m.state == stateSplash && m.frame > 100 {
		m.state = stateSetup
	}

	// Update viewport
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) startProcess() tea.Cmd {
	if !m.flags.ForceRun && !isatty.IsTerminal(os.Stdout.Fd()) {
		m.viewport.WriteLine("Error: Not a TTY. Use --force-run.")
		return nil
	}

	args := m.args
	if m.selected == "watch" {
		args = append([]string{"--watch"}, args...)
	} else if m.selected == "run" {
		args = append([]string{"5"}, args...)
	} else if m.selected == "free" {
		args = []string{"free"}
	}
	
	// Actually invoke the CLI (ralph.js -> ralph.sh mechanism, but we call 'vibepup' assuming it's in path or we call the shell script directly)
	// For local dev, we might need to call the script directly if 'vibepup' isn't in PATH.
	// But let's assume 'vibepup' is the command.
	runCmd := "vibepup"
	if m.flags.Runner != "" {
		runCmd = m.flags.Runner
	}
	
	// If running locally from repo, we might want to call the script directly?
	// The user said "run this project from the build".
	// We'll stick to "vibepup" and assume it's linked or we can use absolute path if needed.
	// Let's use the first arg as the command if provided, or default to "vibepup"
	
	m.dogState = "running"
	m.viewport.WriteLine("--- Starting Vibepup ---")
	
	var cmd tea.Cmd
	m.runner, cmd = process.Start(context.Background(), runCmd, args)
	
	return tea.Batch(cmd, m.runner.WaitForOutput())
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// 1. Splash
	if m.state == stateSplash {
		return ui.BoxStyle.Render(
			lipgloss.JoinVertical(lipgloss.Center,
				lipgloss.NewStyle().Foreground(m.theme.Accent).Render("♥ VIBEPUP TUI ♥"),
				"\nLoading Vibes...\n",
				m.spinner.View(),
			),
		)
	}

	// 2. Form
	if m.state == stateSetup || (m.state == stateRunning && m.runner == nil && m.selected == "new" && m.newForm.State != huh.StateCompleted) {
		form := m.form.View()
		if m.selected == "new" {
			form = m.newForm.View()
		}
		return ui.BoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Foreground(m.theme.Accent).Render("SETUP"),
			form,
		))
	}

	// 3. Running
	header := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(m.theme.Highlight).Render("♥ "+persona.GetStatus(m.selected, m.snark)+" ♥ "+m.spinner.View()),
		motion.GetDogFrame(m.dogState, m.frame),
		persona.RandomQuip(m.snark),
	)

	return ui.BoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
		m.help.View(m.keys),
	))
}

func main() {
	flags := config.Parse()
	m := initialModel(flags)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
