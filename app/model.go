package app

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/rebelice/lazypsql/postgres"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	id    string
	title string
	desc  string
}

func (i Item) Title() string       { return zone.Mark(i.id, i.title) }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return zone.Mark(i.id, i.title) }

type Model struct {
	Database *postgres.Database

	SchemaList *list.Model
	TableList  *list.Model

	Err error
}

var (
	initItems = []list.Item{
		// an ID field has been added here, however it's not required. You could use
		// any text field as long as it's unique for the zone.
		Item{id: "item_1", title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		Item{id: "item_2", title: "Nutella", desc: "It's good on toast"},
		Item{id: "item_3", title: "Bitter melon", desc: "It cools you down"},
		Item{id: "item_4", title: "Nice socks", desc: "And by that I mean socks without holes"},
		Item{id: "item_5", title: "Eight hours of sleep", desc: "I had this once"},
		Item{id: "item_6", title: "Cats", desc: "Usually"},
		Item{id: "item_7", title: "Plantasia, the album", desc: "My plants love it too"},
		Item{id: "item_8", title: "Pour over coffee", desc: "It takes forever to make though"},
		Item{id: "item_9", title: "VR", desc: "Virtual reality...what is there to say?"},
		Item{id: "item_10", title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		Item{id: "item_11", title: "Linux", desc: "Pretty much the best OS"},
		Item{id: "item_12", title: "Business school", desc: "Just kidding"},
		Item{id: "item_13", title: "Pottery", desc: "Wet clay is a great feeling"},
		Item{id: "item_14", title: "Shampoo", desc: "Nothing like clean hair"},
		Item{id: "item_15", title: "Table tennis", desc: "It’s surprisingly exhausting"},
		Item{id: "item_16", title: "Milk crates", desc: "Great for packing in your extra stuff"},
		Item{id: "item_17", title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		Item{id: "item_18", title: "Stickers", desc: "The thicker the vinyl the better"},
		Item{id: "item_19", title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		Item{id: "item_20", title: "Warm light", desc: "Like around 2700 Kelvin"},
		Item{id: "item_21", title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		Item{id: "item_22", title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		Item{id: "item_23", title: "Terrycloth", desc: "In other words, towel fabric"},
	}
)

func NewModel(database *postgres.Database, f *os.File) tea.Model {
	zone.NewGlobal()

	schemas := []list.Item{
		Item{id: "schema_1", title: "public", desc: "public schema"},
		Item{id: "schema_2", title: "dev_schema", desc: "dev schema"},
		Item{id: "schema_3", title: "prod_schema", desc: "prod schema"},
	}
	// items := initItems
	// items := []list.Item{}
	schemaList := list.New(schemas, list.NewDefaultDelegate(), 0, 0)
	result := Model{
		SchemaList: &schemaList,
		TableList:  nil,
	}
	result.SchemaList.Title = "Left click on an items title to select it"
	// result.TableList.Title = "Left click on an items title to select it"

	result.Database = database
	return &result
}

func (m *Model) Init() tea.Cmd {
	return m.ConnectDatabase
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConnectMsg:
		for len(m.SchemaList.Items()) > 0 {
			m.SchemaList.RemoveItem(0)
		}
		for i, schema := range m.Database.Metadata.Schemas {
			m.SchemaList.InsertItem(i, Item{id: fmt.Sprintf("schema_%d", i), title: schema.Name, desc: "This is desc"})
		}
		m.SchemaList.Title = "Choose the schema"

		// m.TableList = list.New(initItems, list.NewDefaultDelegate(), 0, 0)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.SchemaList.SetSize((msg.Width-h)/2, msg.Height-v)
		// m.TableList.SetSize((msg.Width-h)/2, msg.Height-v)
	case tea.MouseMsg:
		if msg.Type == tea.MouseWheelUp {
			m.SchemaList.CursorUp()
			return m, nil
		}

		if msg.Type == tea.MouseWheelDown {
			m.SchemaList.CursorDown()
			return m, nil
		}

		if msg.Type == tea.MouseLeft {
			for i, listItem := range m.SchemaList.VisibleItems() {
				item, _ := listItem.(Item)
				// Check each item to see if it's in bounds.
				if zone.Get(item.id).InBounds(msg) {
					// If so, select it in the list.
					m.SchemaList.Select(i)
					break
				}
			}
		}

		return m, nil
	case ErrMsg:
		m.Err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	*m.SchemaList, cmd = m.SchemaList.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	// Wrap the main models view in zone.Scan.
	body := docStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.SchemaList.View(), //m.TableList.View(),
		),
	)
	dialog := lipgloss.Place(10, 10, lipgloss.Center, lipgloss.Center, "This is floating window!", lipgloss.WithWhitespaceChars(body))
	// return zone.Scan(body + dialog)
	if m.Err != nil {
		body += "\n" + m.Err.Error()
	}
	// return zone.Scan(body)
	return zone.Scan(dialog)
}

type ErrMsg struct {
	err error
}

func (e ErrMsg) Error() string {
	return e.err.Error()
}
