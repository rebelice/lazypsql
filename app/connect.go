package app

import (
	"sort"

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

func (msg ConnectMsg) SchemaList() []string {
	var schemas []string
	for schemaName := range msg.Database.Metadata.Schemas {
		schemas = append(schemas, schemaName)
	}
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i] < schemas[j]
	})
	return schemas
}

type SyncTablesMsg struct {
	Schema   string
	Database *postgres.Database
}

func (m SyncTablesMsg) TableList() []string {
	var result []string
	for _, table := range m.Database.Metadata.Schemas[m.Schema].Tables {
		result = append(result, table.Name)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	return result
}

func (m Model) SyncTables() tea.Msg {
	if err := m.Database.FetchTables(m.CurrentSchema()); err != nil {
		return ErrMsg{err}
	}

	return SyncTablesMsg{m.CurrentSchema(), m.Database}
}
