package main

// Gonna try using ast-grep to prototype the idea
// package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ngynkvn/crepe/cmd/simple/app"
)

const (
	initialInputs = 1
	maxInputs     = 6
	minInputs     = 1
	helpHeight    = 5
)

type keymap = struct {
	next, quit key.Binding
}

type AstGrepCommand struct{}

type Size struct {
	width  int
	height int
}

type model struct {
	help   help.Model
	keymap keymap
	input  textinput.Model
	panel  Panel
	Size
}

type Panel struct {
	buffer string
	Size
}

func (p *Panel) SetValue(s string) {
	p.buffer = s
}

func (m *model) Resize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
	m.input.Width = msg.Width
	m.panel.Resize(msg)
}

func (p *Panel) Resize(msg tea.WindowSizeMsg) {
	p.height = msg.Height
	p.width = msg.Width
}

func (p Panel) View() string {
	return lipgloss.NewStyle().
		Padding(1).
		Width(p.width-10).
		Height(p.height-10).
		Border(lipgloss.RoundedBorder(), true).
		Render(p.buffer)
}

func (agc AstGrepCommand) Grep(input string) *exec.Cmd {
	cmd := exec.Command(
		"ast-grep",
		"run",
		"--color",
		"always",
		"--pattern",
		input,
	)
	cmd.Env = os.Environ()
	return cmd
}

func newModel() model {
	m := model{
		input: textinput.New(),
		panel: Panel{},
		help:  help.New(),
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.input.Focus()
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		}
		agc := AstGrepCommand{}
		oscmd := agc.Grep(m.input.Value())
		output, err := oscmd.Output()
		if err != nil {
			m.panel.SetValue(err.Error())
		} else {
			m.panel.SetValue(string(output))
		}
	case tea.WindowSizeMsg:
		m.Resize(msg)
	}

	// Update all textareas
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.quit,
	})

	var views []string
	views = append(views, app.DefaultStyles.InputStyle.Render(m.input.View()))
	views = append(views, m.panel.View())

	return lipgloss.JoinVertical(lipgloss.Center, views...) + "\n\n" + help
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
