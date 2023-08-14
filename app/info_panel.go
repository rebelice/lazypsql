package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type infoType string

type InfoPanel struct {
	list.Model

	style lipgloss.Style

	id string

	database string
	user     string
	schema   string
}

type infoItem struct {
	id    string
	key   string
	value string
}

func (i infoItem) Title() string {
	return zone.Mark(i.id, fmt.Sprintf("%s: %s", i.key, i.value))
}

func (i infoItem) Description() string {
	return ""
}

func (i infoItem) FilterValue() string {
	return zone.Mark(i.id, fmt.Sprintf("%s: %s", i.key, i.value))
}

func NewInfoPanel(id string) *InfoPanel {
	infoPanel := list.New([]list.Item{
		infoItem{id: "info_panel_database", key: "Database", value: ""},
		infoItem{id: "info_panel_user", key: "User", value: ""},
		infoItem{id: "info_panel_schema", key: "Schema", value: ""},
	}, list.NewDefaultDelegate(), 0, 0)
	infoPanel.Title = "Database Information"
	infoPanel.SetShowStatusBar(false)
	infoPanel.SetShowHelp(false)
	return &InfoPanel{
		id:    id,
		Model: infoPanel,
		style: NewUnfocusedModelStyle(0, 0),
	}
}

func (*InfoPanel) Init() tea.Cmd {
	return nil
}

func (s *InfoPanel) Update(msg tea.Msg) (*InfoPanel, tea.Cmd) {
	model, cmd := s.Model.Update(msg)
	s.Model = model
	return s, cmd
}

func (s *InfoPanel) View() string {
	return s.style.Render(s.Model.View())
}

func (s *InfoPanel) SetSize(width, height int) {
	h, v := s.style.Width(width).Height(height).GetFrameSize()
	s.Model.SetSize(width-h, height-v)
}
