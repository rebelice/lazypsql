package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SchemaList struct {
	list.Model

	id    string
	style lipgloss.Style
}

func NewSchemaList(id string) *SchemaList {
	schemaList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	schemaList.Title = "Choose a schema"
	return &SchemaList{
		id:    id,
		Model: schemaList,
		style: NewFocusedModelStyle(0, 0),
	}
}

func (*SchemaList) Init() tea.Cmd {
	return nil
}

func (s *SchemaList) Update(msg tea.Msg) (*SchemaList, tea.Cmd) {
	model, cmd := s.Model.Update(msg)
	s.Model = model
	return s, cmd
}

func (s *SchemaList) View() string {
	return s.style.Render(s.Model.View())
}

func (s *SchemaList) SetSize(width, height int) {
	h, v := s.style.Width(width).Height(height).GetFrameSize()
	// panic(fmt.Sprintf("h: %d, v: %d, width: %d, height: %d", h, v, width, height))
	s.Model.SetSize(width-h, height-v)
}
