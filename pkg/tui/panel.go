package tui

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle   = inactiveTabStyle.Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("21"))
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.RoundedBorder())
)

type panel struct {
	tabs       []string
	tabContent []string
	activeTab  int
	width      int
}

func (p *panel) View() string {
	var renderedTabs []string
	for i, t := range p.tabs {
		style := inactiveTabStyle
		if i == p.activeTab {
			style = activeTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	markdownOutput, err := renderMarkdown(p.tabContent[p.activeTab])
	if err != nil {
		markdownOutput = "Failed to render markdown content: " + err.Error()
	}

	content := windowStyle.Width((p.width / 2)).Render(markdownOutput)

	return lipgloss.JoinVertical(lipgloss.Top, tabsRow, content)
}

func renderMarkdown(md string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
	)
	if err != nil {
		return "", err
	}

	return renderer.Render(md)
}
