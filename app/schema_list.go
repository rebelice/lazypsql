package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type schemaItem struct {
	id    string
	title string
	desc  string
}

func (i schemaItem) Title() string       { return zone.Mark(i.id, i.title) }
func (i schemaItem) Description() string { return i.desc }
func (i schemaItem) FilterValue() string { return zone.Mark(i.id, i.title) }

type SchemaList struct {
	*list.Model

	id    string
	style lipgloss.Style
}

func NewSchemaList(id string) *SchemaList {
	schemaList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	schemaList.Title = "Choose a schema"
	return &SchemaList{
		id:    id,
		Model: &schemaList,
		style: NewFocusedModelStyle(0, 0),
	}
}

func (*SchemaList) Init() tea.Cmd {
	return nil
}

func (s *SchemaList) Update(msg tea.Msg) (*SchemaList, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case ConnectMsg:
		for len(s.Model.Items()) > 0 {
			s.Model.RemoveItem(0)
		}
		for i, schema := range msg.Database.Metadata.Schemas {
			cmd := s.Model.InsertItem(i, schemaItem{id: fmt.Sprintf("schema_%d", i), title: schema.Name})
			cmds = append(cmds, cmd)
		}
		// var items []list.Item
		// for i, schema := range msg.Database.Metadata.Schemas {
		// 	items = append(items, schemaItem{id: fmt.Sprintf("schema_%d", i), title: schema.Name})
		// }
		// l := list.New(items, list.NewDefaultDelegate(), 0, 0)
		// s.Model = &l
		// return s, nil
	}

	model, cmd := s.Model.Update(msg)
	cmds = append(cmds, cmd)
	s.Model = &model
	return s, tea.Batch(cmds...)
}

func (s *SchemaList) View() string {
	return s.style.Render(s.Model.View())
}

func (s *SchemaList) SetSize(width, height int) {
	h, v := s.style.Width(width).Height(height).GetFrameSize()
	// panic(fmt.Sprintf("h: %d, v: %d, width: %d, height: %d", h, v, width, height))
	s.Model.SetSize(width-h, height-v)
}
