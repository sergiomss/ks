package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var promptStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3DE063"))

var boldStyle = lipgloss.NewStyle().
	Bold(true)

var instructionsStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#56c7c0"))

var hooverStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#18F1E3"))

type model struct {
	options    []string
	choice     string
	selected   string
	cursor     int
	choiceType string
}

func (m model) Init() tea.Cmd {
	m.options = []string{"Buy carrots", "But celery", "Buy kohlrabi"}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else if m.cursor == 0 {
				m.cursor = len(m.options) - 1
			}
		case "down":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			} else if m.cursor == len(m.options)-1 {
				m.cursor = 0
			}
		case "enter":
			m.selected = m.options[m.cursor]
			return m, tea.Sequence(tea.Println(fmt.Sprintf("Successfully switched to context: %v", m.selected), tea.Quit()))
		}
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("%v %v   %v\n",
		promptStyle.Render("?"),
		boldStyle.Render(fmt.Sprintf("Choose a %v:", m.choiceType)),
		instructionsStyle.Render("[Use arrows to move, enter to select, type to filter]"))

	if m.selected == "" {
		for i, choice := range m.options {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s += fmt.Sprintf("%s\n", hooverStyle.Render(cursor+" "+choice))
			} else {
				s += fmt.Sprintf("%s %s\n", cursor, choice)
			}
		}
	} else {
		s += fmt.Sprintf("Successfully switched to context: %v", m.selected)
	}

	return s
}
