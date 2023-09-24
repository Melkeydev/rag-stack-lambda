package textinput

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	exitColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E11D48")).Bold(true)
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
)

type (
	errorMsg error
)

type Output struct {
	Output string
}

func (o *Output) update(val string) {
	o.Output = val
}

type model struct {
	textinput textinput.Model
	err       error
	header    string
	output    *Output
}

func initialModel(placeholder string, header string, output *Output) model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 20
	return model{
		textinput: ti,
		err:       nil,
		header:    titleStyle.Render(header),
		output:    output,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC:
			if len(m.textinput.Value()) > 1 {
				m.output.update(m.textinput.Value())

				return m, tea.Quit
			}
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			os.Exit(1)
		}
	case errorMsg:
		m.err = msg
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		m.header,
		m.textinput.View(),
		fmt.Sprintf("press %s to quit", exitColor.Render("q")))
}

func TextInputRun(placeholder string, header string, output *Output) {
	p := tea.NewProgram(initialModel(placeholder, header, output))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
