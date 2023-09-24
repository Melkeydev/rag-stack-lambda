package creator

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type ExitProgram struct {
	Value bool
}

func (e *ExitProgram) CheckExitStatus(program *tea.Program) {
	if e.Value {
		program.ReleaseTerminal()
		os.Exit(1)
	}
}

