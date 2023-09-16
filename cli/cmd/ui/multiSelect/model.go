package multiselect

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type Selection struct {
	Choice string
}

func (s *Selection) Update(value string) {
	s.Choice = value
}

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
	choice   *Selection
	header   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel(choices []string, selection *Selection, header string) model {
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
		choice:   selection,
		header:   header,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			os.Exit(1)
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.selected) == 1 {
				m.selected = make(map[int]struct{})
			}
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "y":
			if len(m.selected) == 1 {
				m.choice.Update(m.choices[m.cursor])
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := m.header + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress y to confirm choice.\n"
	s += "\nPress q to quit.\n\n"

	return s
}

func MultiBoxSelectRun(options []string, selection *Selection, header string) {
	p := tea.NewProgram(InitialModel(options, selection, header))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
