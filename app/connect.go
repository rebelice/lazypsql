package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rebelice/lazypsql/postgres"
)

type ConnectMsg struct {
	Database *postgres.Database
}

func (m Model) ConnectDatabase() tea.Msg {
	if err := m.Database.Connect(); err != nil {
		return ErrMsg{err}
	}

	if err := m.Database.Ping(); err != nil {
		return ErrMsg{err}
	}

	if err := m.Database.FetchSchemas(); err != nil {
		return ErrMsg{err}
	}

	return ConnectMsg{m.Database}
}
