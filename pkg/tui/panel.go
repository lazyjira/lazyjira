package tui

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type panel struct {
	tabs       []string
	tabContent []string
	activeTab  int
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

	content := windowStyle.Render(markdownOutput)

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
