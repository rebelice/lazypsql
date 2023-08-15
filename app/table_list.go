package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type tableItem struct {
	id   string
	name string
}

func (i tableItem) Title() string       { return zone.Mark(i.id, i.name) }
func (i tableItem) Description() string { return "" }
func (i tableItem) FilterValue() string { return zone.Mark(i.id, i.name) }

type TableList struct {
	id    string
	model *list.Model
	style lipgloss.Style
}

func NewTableList(id string) *TableList {
	tableList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	tableList.Title = "Tables"
	return &TableList{
		id:    id,
		model: &tableList,
		style: NewUnfocusedModelStyle(0, 0),
	}
}

func (*TableList) Init() tea.Cmd {
	return nil
}

func (t *TableList) Update(msg tea.Msg) (*TableList, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case SyncTablesMsg:
		for len(t.model.Items()) > 0 {
			t.model.RemoveItem(0)
		}
		for i, table := range msg.TableList() {
			cmd := t.model.InsertItem(i, tableItem{
				id:   fmt.Sprintf("table_%d", i),
				name: table,
			})
			cmds = append(cmds, cmd)
		}
		return t, tea.Batch(cmds...)
	}
	return t, nil
}

func (t *TableList) View() string {
	return t.style.Render(t.model.View())
}

func (t *TableList) SetSize(width, height int) {
	h, v := t.style.Width(width).Height(height).GetFrameSize()
	t.model.SetSize(width-h, height-v)
}
