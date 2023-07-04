package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return "Hello, World!"
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ctrl+c or q exits. Even with short running programs it's good to have
		// a quit key, just in case your logic is off. Users will be very
		// annoyed if they can't exit.
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	// If we happen to get any other messages, don't do anything.
	return m, nil
}
