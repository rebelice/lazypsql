package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CommandPanel struct {
	textinput.Model

	id string
}

func NewCommandPanel(id string) *CommandPanel {
	input := textinput.New()
	return &CommandPanel{
		id:    id,
		Model: input,
	}
}

func (*CommandPanel) Init() tea.Cmd {
	return nil
}

func (c *CommandPanel) Update(msg tea.Msg) (*CommandPanel, tea.Cmd) {
	model, cmd := c.Model.Update(msg)
	c.Model = model
	return c, cmd
}

func (c *CommandPanel) View() string {
	return c.Model.View()
}
