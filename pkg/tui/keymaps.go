package tui

import "github.com/charmbracelet/bubbles/key"

type projectKM struct {
	toggle key.Binding
}

func (m projectKM) ShortHelp() []key.Binding {
	return []key.Binding{m.toggle}
}

func (m projectKM) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.toggle},
	}
}

var projectKMKeys = projectKM{
	toggle: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "switch project"),
	),
}
