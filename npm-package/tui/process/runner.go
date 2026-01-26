package process

import (
	"bufio"
	"context"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

// Msg indicates a line of output from the process
type OutputMsg string

// DoneMsg indicates the process finished
type DoneMsg struct {
	Err error
}

// Runner handles the execution of the external process
type Runner struct {
	Cmd        *exec.Cmd
	Cancel     context.CancelFunc
	OutputChan chan string
}

// Start launches the command in a new process group to allow deep killing
func Start(ctx context.Context, name string, args []string) (*Runner, tea.Cmd) {
	ctx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(ctx, name, args...)

	// Create a new process group so we can kill the whole tree later
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Capture stdout and stderr
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	runner := &Runner{
		Cmd:        cmd,
		Cancel:     cancel,
		OutputChan: make(chan string),
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, func() tea.Msg { return DoneMsg{Err: err} }
	}

	// Stream output to channel
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			runner.OutputChan <- scanner.Text()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			runner.OutputChan <- "ERR: " + scanner.Text()
		}
	}()

	// Wait for completion in background
	cmdCmd := func() tea.Msg {
		err := cmd.Wait()
		return DoneMsg{Err: err}
	}

	return runner, cmdCmd
}

// Kill stops the process and its children
func (r *Runner) Kill() {
	if r.Cmd != nil && r.Cmd.Process != nil {
		// Kill the entire process group (negative PID)
		_ = syscall.Kill(-r.Cmd.Process.Pid, syscall.SIGKILL)
	}
	if r.Cancel != nil {
		r.Cancel()
	}
}

// WaitForOutput returns a command that waits for the next line of output
func (r *Runner) WaitForOutput() tea.Cmd {
	return func() tea.Msg {
		line, ok := <-r.OutputChan
		if !ok {
			return nil
		}
		return OutputMsg(line)
	}
}
