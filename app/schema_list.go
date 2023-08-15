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

type schemaListState string

const (
	schemaListStateChoosing schemaListState = "choosing"
	schemaListStateChosen   schemaListState = "chosen"
)

type SchemaList struct {
	model *list.Model

	id    string
	style lipgloss.Style

	state schemaListState
}

func NewSchemaList(id string) *SchemaList {
	schemaList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	schemaList.Title = "Choose a schema"
	return &SchemaList{
		id:    id,
		model: &schemaList,
		style: NewFocusedModelStyle(0, 0),
		state: schemaListStateChoosing,
	}
}

func (*SchemaList) Init() tea.Cmd {
	return nil
}

func (s *SchemaList) Update(msg tea.Msg) (*SchemaList, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case ConnectMsg:
		for len(s.model.Items()) > 0 {
			s.model.RemoveItem(0)
		}
		for i, schema := range msg.SchemaList() {
			cmd := s.model.InsertItem(i, schemaItem{id: fmt.Sprintf("schema_%d", i), title: schema})
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "o":
			s.state = schemaListStateChosen
		}
	}

	model, cmd := s.model.Update(msg)
	cmds = append(cmds, cmd)
	s.model = &model
	return s, tea.Batch(cmds...)
}

func (s *SchemaList) View() string {
	return s.style.Render(s.model.View())
}

func (s *SchemaList) SetSize(width, height int) {
	h, v := s.style.Width(width).Height(height).GetFrameSize()
	// panic(fmt.Sprintf("h: %d, v: %d, width: %d, height: %d", h, v, width, height))
	s.model.SetSize(width-h, height-v)
}
