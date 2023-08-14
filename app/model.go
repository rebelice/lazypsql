package app

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/rebelice/lazypsql/postgres"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ModeState string

const (
	ModeStateFocusInfo    ModeState = "state.focus-info"
	ModeStateCommandMode  ModeState = "state.command-mode"
	ModeStateChooseSchema ModeState = "state.choose-schema"
)

type Model struct {
	State    ModeState
	Database *postgres.Database

	Height int
	Width  int

	CommandPanel *CommandPanel
	SchemaList   *SchemaList
	InfoPanel    *InfoPanel

	Err error
}

func NewModel(database *postgres.Database, f *os.File) tea.Model {
	zone.NewGlobal()

	// schemas := []list.Item{
	// 	Item{id: "schema_1", title: "public", desc: "public schema"},
	// 	Item{id: "schema_2", title: "dev_schema", desc: "dev schema"},
	// 	Item{id: "schema_3", title: "prod_schema", desc: "prod schema"},
	// }
	// // items := initItems
	// // items := []list.Item{}
	// schemaList := list.New(schemas, list.NewDefaultDelegate(), 0, 0)
	// result := Model{
	// 	SchemaList: schemaList,
	// }
	// result.SchemaList.Title = "Left click on an items title to select it"
	// // result.TableList.Title = "Left click on an items title to select it"

	result := Model{
		State: ModeStateChooseSchema,
	}
	result.Database = database
	result.SchemaList = NewSchemaList("schema_list")
	result.CommandPanel = NewCommandPanel("command_panel")
	result.InfoPanel = NewInfoPanel("info_panel")
	return &result
}

func (m *Model) Init() tea.Cmd {
	return m.ConnectDatabase
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	// Common Updates
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+p":
			m.State = ModeStateCommandMode
			return m, m.CommandPanel.Focus()
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if zone.Get(m.CommandPanel.id).InBounds(msg) {
				m.State = ModeStateCommandMode
				return m, m.CommandPanel.Focus()
			}
		}
	case ConnectMsg:
		var cmd tea.Cmd
		m.SchemaList, cmd = m.SchemaList.Update(msg)
		cmds = append(cmds, cmd)
		m.InfoPanel, cmd = m.InfoPanel.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	// Mode specific updates
	switch m.State {
	// case ModeStateFocusInfo:
	// 	var cmd tea.Cmd
	// 	m.InfoPanel, cmd = m.InfoPanel.Update(msg)
	// 	return m, cmd
	case ModeStateCommandMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.State = ModeStateChooseSchema
				m.CommandPanel.Blur()
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.CommandPanel, cmd = m.CommandPanel.Update(msg)
		return m, cmd
	case ModeStateChooseSchema:
		var cmd tea.Cmd
		m.SchemaList, cmd = m.SchemaList.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case ConnectMsg:
		// for len(m.SchemaList.Items()) > 0 {
		// 	m.SchemaList.RemoveItem(0)
		// }
		// for i, schema := range m.Database.Metadata.Schemas {
		// 	cmd := m.SchemaList.InsertItem(i, schemaItem{id: fmt.Sprintf("schema_%d", i), title: schema.Name})
		// 	cmds = append(cmds, cmd)
		// }

		// m.TableList = list.New(initItems, list.NewDefaultDelegate(), 0, 0)
	case tea.WindowSizeMsg:
		horizontalFrame, verticalFrame := docStyle.GetFrameSize()
		w, h := msg.Width-horizontalFrame, msg.Height-verticalFrame-1
		// m.Err = errors.New(fmt.Sprintf("w: %d, h: %d, hf: %d, vf: %d", w, h, horizontalFrame, verticalFrame))
		infoPanelHeight := h / 3
		m.InfoPanel.SetSize(w/3, infoPanelHeight-2)
		m.SchemaList.SetSize(w/3, h-infoPanelHeight-2)
		m.CommandPanel.Width = w
		// return m, nil
	// case tea.MouseMsg:
	// 	if msg.Type == tea.MouseWheelUp {
	// 		m.SchemaList.CursorUp()
	// 		return m, nil
	// 	}

	// 	if msg.Type == tea.MouseWheelDown {
	// 		m.SchemaList.CursorDown()
	// 		return m, nil
	// 	}

	// 	if msg.Type == tea.MouseLeft {
	// 		for i, listItem := range m.SchemaList.VisibleItems() {
	// 			item, _ := listItem.(Item)
	// 			// Check each item to see if it's in bounds.
	// 			if zone.Get(item.id).InBounds(msg) {
	// 				// If so, select it in the list.
	// 				m.SchemaList.Select(i)
	// 				break
	// 			}
	// 		}
	// 	}

	// 	return m, nil
	case ErrMsg:
		m.Err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	m.SchemaList, cmd = m.SchemaList.Update(msg)
	cmds = append(cmds, cmd)
	m.InfoPanel, cmd = m.InfoPanel.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// Wrap the main models view in zone.Scan.
	body := docStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.InfoPanel.View(),
				m.SchemaList.View(), //m.TableList.View(),
			),
			zone.Mark(m.CommandPanel.id, m.CommandPanel.View()),
			// m.CommandPanel.View(),
		),
	)
	if m.Err != nil {
		body += "\n" + m.Err.Error()
	}
	return zone.Scan(body)
}

type ErrMsg struct {
	err error
}

func (e ErrMsg) Error() string {
	return e.err.Error()
}
