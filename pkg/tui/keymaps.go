package tui

import "github.com/charmbracelet/bubbles/key"

type projectKM struct {
	toggle key.Binding
}

type detailsKM struct {
	prevTab key.Binding
	nextTab key.Binding
}

func (m projectKM) ShortHelp() []key.Binding {
	return []key.Binding{m.toggle}
}

func (m projectKM) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.toggle},
	}
}

func (m detailsKM) ShortHelp() []key.Binding {
	return []key.Binding{m.prevTab, m.nextTab}
}

func (m detailsKM) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.prevTab, m.nextTab},
	}
}

var projectKMKeys = projectKM{
	toggle: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "switch project"),
	),
}

var detailsKMKeys = detailsKM{
	prevTab: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "previous tab"),
	),
	nextTab: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next tab"),
	),
}
