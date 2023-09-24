package spinner

import (
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

var (
	exitColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#E11D48")).Bold(true)
)

type LoadingState struct {
	Loading bool
}

type model struct {
	spinner      spinner.Model
	err          error
	msg          string
	LoadingState *LoadingState
	wg           *sync.WaitGroup
}

func initialModel(loaderMsg string, LoadingState *LoadingState, wg *sync.WaitGroup) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6"))
	return model{spinner: s, msg: loaderMsg, LoadingState: LoadingState, wg: wg}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.LoadingState.Loading = false
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		if !m.LoadingState.Loading {
			m.wg.Done()
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n   %s %s... press %s to quit\n", m.spinner.View(), m.msg, exitColor.Render("q"))
	if !m.LoadingState.Loading {
		return str + "\n"
	}
	return str
}

func SpineMe(loaderMsg string, LoadingState *LoadingState, wg *sync.WaitGroup) *tea.Program {
	p := tea.NewProgram(initialModel(loaderMsg, LoadingState, wg))
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	return p
}