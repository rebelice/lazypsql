package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) ConnectDatabase() tea.Msg {
	if err := m.Database.Connect(); err != nil {
		return ErrMsg{err}
	}

	return nil
}
