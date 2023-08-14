package app

import "github.com/charmbracelet/lipgloss"

func NewFocusedModelStyle(width int, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Left, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))
}

func NewUnfocusedModelStyle(width int, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Left, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("241"))
}
