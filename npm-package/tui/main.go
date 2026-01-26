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
	"github.com/mattn/go-isatty"

	"vibepup-tui/animations"
	"vibepup-tui/config"
	"vibepup-tui/motion"
	"vibepup-tui/persona"
	"vibepup-tui/theme"
	"vibepup-tui/ui"
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

type styleSet struct {
	box    lipgloss.Style
	title  lipgloss.Style
	tip    lipgloss.Style
	status lipgloss.Style
	text   lipgloss.Style
}

func buildStyles(t theme.Theme) styleSet {
	return styleSet{
		box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Padding(1, 2).
			Margin(1, 1).
			Background(t.Background),
		title: lipgloss.NewStyle().
			Foreground(t.Background).
			Background(t.Accent).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1),
		tip: lipgloss.NewStyle().
			Foreground(t.AccentAlt).
			Italic(true).
			MarginTop(1),
		status: lipgloss.NewStyle().
			Foreground(t.Highlight).
			Bold(true),
		text: lipgloss.NewStyle().
			Foreground(t.Foreground),
	}
}

// --- Key Bindings ---

type KeyMap struct {
	Quit      key.Binding
	Help      key.Binding
	NextTheme key.Binding
	NextAnim  key.Binding
	NextSnark key.Binding
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
		NextTheme: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "next theme"),
		),
		NextAnim: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "next anim"),
		),
		NextSnark: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "next snark"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.NextTheme, k.NextAnim, k.NextSnark}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit, k.NextTheme, k.NextAnim, k.NextSnark},
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
	flags      config.Flags
	theme      theme.Theme
	styles     styleSet
	snark      persona.SnarkLevel
	anim       animations.Preset
	animFrame  int
	motion     motion.Engine
	statusBar  ui.StatusBar
	proc       *exec.Cmd
	procCancel context.CancelFunc
}

type processStartedMsg struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
	done   <-chan processDoneMsg
}

type processDoneMsg struct {
	err error
}

func initialModel(flags config.Flags) model {
	args := []string{"--watch"}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	spring := harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.2)
	th := theme.Get(flags.Theme)
	styles := buildStyles(th)
	snark := persona.ParseSnark(flags.Snark)
	motionEngine := motion.New(flags.PerfLow, flags.Quiet)
	animPreset := animations.Get(flags.Anim)

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
		flags:      flags,
		theme:      th,
		styles:     styles,
		snark:      snark,
		motion:     motionEngine,
		anim:       animPreset,
		statusBar:  ui.StatusBar{Theme: th},
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
	return m.motion.Next()
}

func (m *model) stopProc() {
	if m.proc != nil && m.proc.ProcessState == nil {
		if m.procCancel != nil {
			m.procCancel()
		}
		_ = m.proc.Process.Kill()
	}
	m.proc = nil
	m.procCancel = nil
}

func nextTheme(name string) theme.Theme {
	all := theme.All()
	for i, t := range all {
		if t.Name == name {
			return all[(i+1)%len(all)]
		}
	}
	return all[0]
}

func nextAnim(name string) animations.Preset {
	all := animations.All()
	for i, a := range all {
		if a.Name == name {
			return all[(i+1)%len(all)]
		}
	}
	return all[0]
}

func nextSnark(s persona.SnarkLevel) persona.SnarkLevel {
	levels := []persona.SnarkLevel{persona.Mild, persona.Spicy, persona.Unhinged}
	for i, v := range levels {
		if v == s {
			return levels[(i+1)%len(levels)]
		}
	}
	return persona.Mild
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
			m.stopProc()
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.NextTheme):
			m.theme = nextTheme(m.theme.Name)
			m.styles = buildStyles(m.theme)
			m.statusBar = ui.StatusBar{Theme: m.theme}
		case key.Matches(msg, m.keys.NextAnim):
			m.anim = nextAnim(m.anim.Name)
			m.animFrame = 0
		case key.Matches(msg, m.keys.NextSnark):
			m.snark = nextSnark(m.snark)
		}

	case motion.TickMsg:
		m.frame++
		m.punTimer++
		if n := len(m.anim.Frames); n > 0 {
			m.animFrame = (m.animFrame + 1) % n
		}

		// Rotate tips/puns
		if m.punTimer > 80 { // Every ~4 seconds
			m.punTimer = 0
			m.currentPun = persona.Quip(persona.StateWaiting, m.snark, nil)
			if m.currentPun == "" {
				m.currentPun = puns[rand.Intn(len(puns))]
			}
			m.currentTip = tips[rand.Intn(len(tips))]
		}

		switch m.state {
		case stateSplash:
			if time.Since(m.start) > time.Second*3 {
				m.state = stateSetup
				return m, m.form.Init()
			}
			return m, m.motion.Next()

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
			return m, m.motion.Next()

		default:
			return m, m.motion.Next()
		}

	case processStartedMsg:
		m.proc = msg.cmd
		m.procCancel = msg.cancel
		if msg.done != nil {
			cmds = append(cmds, func() tea.Msg { return <-msg.done })
		}

	case processDoneMsg:
		if msg.err != nil {
			m.currentPun = "Process bailed: " + msg.err.Error()
		}
		m.stopProc()
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
			return m, startProcess(m.selected, m.args, m.flags)
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
				return m, startProcess("new", append([]string{idea}, m.args...), m.flags)
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
	thinkingFrames := []string{"… thinking …", "… scheming …", "… cooking bits …", "… brewing chaos …"}
	if m.flags.NoEmoji || !m.theme.SupportsEmoji {
		frames = []string{"(•ᴗ• )", "(•̀ᴗ• )✧", "(•ᴗ• )ノ", "(ᵕ•ᴗ•ᵕ)"}
		sparkles = []string{"*", "✶", "✷", "✸", "✹"}
		thinkingFrames = []string{"thinking", "scheming", "brewing", "loading"}
	}

	// Animation logic
	idx := (m.frame / 4) % len(frames)
	sIdx := (m.frame / 2) % len(sparkles)
	currentDog := frames[idx]
	currentSparkle := sparkles[sIdx]
	thinking := thinkingFrames[(m.frame/3)%len(thinkingFrames)]

	// Spring movement
	pad := ""
	if m.pos > 0 {
		pad = strings.Repeat(" ", int(m.pos))
	}

	// Composite Dog
	dogRender := lipgloss.NewStyle().Foreground(m.theme.Accent).Render(pad + currentDog + " " + currentSparkle)

	animFrame, _ := animations.Frame(m.anim, m.animFrame)

	// --- Views ---

	// 1. Splash
	splashContent := lipgloss.JoinVertical(lipgloss.Center,
		m.styles.title.Render("♥ VIBEPUP TUI ♥"),
		"",
		dogRender,
		"",
		lipgloss.NewStyle().Foreground(m.theme.AccentAlt).Render("Initializing chaos engine..."),
		lipgloss.NewStyle().Foreground(m.theme.Muted).Render(m.currentPun),
	)

	// 2. Setup
	if m.state == stateSplash {
		return m.styles.box.Render(splashContent)
	}

	if m.state == stateSetup {
		return m.styles.box.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				m.styles.title.Render("♥ VIBE CHECK ♥"),
				lipgloss.NewStyle().Foreground(m.theme.Foreground).Render("Ready to break some code?"),
				"",
				m.form.View(),
				"",
				m.help.View(m.keys),
			),
		)
	}

	// 3. New Project Input
	if m.state == stateRunning && m.selected == "new" && m.newForm.State != huh.StateCompleted {
		return m.styles.box.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				m.styles.title.Render("♥ GENESIS PROTOCOL ♥"),
				lipgloss.NewStyle().Foreground(m.theme.AccentAlt).Render("What are we manifesting today?"),
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
	loader := animFrame
	if m.flags.Quiet {
		loader = ui.ClampWidth(loader, 6)
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		m.styles.status.Render(statusMsg+" "+loader+"  "+thinking),
		lipgloss.NewStyle().Foreground(m.theme.AccentAlt).Render("MODE: "+strings.ToUpper(m.selected)),
		"",
		lipgloss.NewStyle().Foreground(m.theme.Foreground).Render(currentSparkle+" "+m.currentPun),
		"",
		dogRender,
		m.styles.tip.Render(m.currentTip),
		lipgloss.NewStyle().Foreground(m.theme.Muted).Render("─ Matrix Stream ──────────────────────────"),
	)

	if !m.ready {
		return m.styles.box.Render(header + "\nBooting up the matrix...")
	}

	statusLine := m.statusBar.Render("q: quit", "snark:"+m.flags.Snark, m.viewport.Width)

	return m.styles.box.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			m.viewport.View(),
			"",
			statusLine,
			m.help.View(m.keys),
		),
	)
}

func startProcess(choice string, args []string, flags config.Flags) tea.Cmd {
	return func() tea.Msg {
		if !flags.ForceRun {
			if !(isatty.IsTerminal(os.Stdout.Fd()) && isatty.IsTerminal(os.Stdin.Fd())) {
				return processDoneMsg{err: fmt.Errorf("not a tty; use --force-run to override")}
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
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
		cmd := exec.CommandContext(ctx, "vibepup", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Start(); err != nil {
			cancel()
			return processDoneMsg{err: err}
		}

		done := make(chan processDoneMsg, 1)
		go func() {
			err := cmd.Wait()
			done <- processDoneMsg{err: err}
			close(done)
		}()

		return processStartedMsg{cmd: cmd, cancel: cancel, done: done}
	}
}

func main() {
	flags := config.Parse()
	m := initialModel(flags)
	programOpts := []tea.ProgramOption{}
	// default: do not force fullscreen/alt unless explicitly requested off
	if !flags.NoAlt {
		programOpts = append(programOpts, tea.WithAltScreen())
	}
	p := tea.NewProgram(m, programOpts...)
	if _, err := p.Run(); err != nil {
		fmt.Println("Oof, the puppy tripped.", err)
		os.Exit(1)
	}
}
