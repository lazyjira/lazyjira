package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
	"log"
)

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("241"))

	panelStyleFocused = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type JiraModel struct {
	width            int
	height           int
	columnSize       int
	selectedIssue    int
	focusedTab       int
	maxPanels        int
	showProjectsList bool
	projectName      string
	client           *jira.Client
	projectItems     []jira.Project
	projectsList     list.Model
	issuesItems      []jira.Issue
	issuesList       list.Model
}

func NewJiraTui(client *jira.Client) JiraModel {

	initProjectsList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	initProjectsList.Title = "Projects"

	initIssuesList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	initIssuesList.Title = "Assigned Issues"

	return JiraModel{
		client:           client,
		focusedTab:       1,
		maxPanels:        3,
		showProjectsList: false,
		projectName:      "Selected Project Name",
		projectsList:     initProjectsList,
		issuesList:       initIssuesList,
	}
}

func (m JiraModel) Init() tea.Cmd {
	return tea.Batch(
		getProjectsList(m.client),
		getIssuesList(m.client),
	)
}

func (m *JiraModel) NextPanel() {
	m.focusedTab = (m.focusedTab + 1) % 4
}

func (m *JiraModel) PrevPanel() {
	m.focusedTab = (m.focusedTab + 2) % 4
}

func (m *JiraModel) setProjectListItems(projects []jira.Project) {
	var items []list.Item
	for _, project := range projects {
		items = append(items, project)
	}

	m.projectsList = list.New(items, list.NewDefaultDelegate(), 30, 15)
	m.projectsList.Title = "Projects"
	m.projectsList.DisableQuitKeybindings()

	m.projectItems = projects
}

func (m *JiraModel) setIssuesListItems(issues []jira.Issue) {
	var items []list.Item
	for _, project := range issues {
		items = append(items, project)
	}

	m.issuesList = list.New(items, list.NewDefaultDelegate(), 30, 15)
	m.issuesList.Title = "Assigned Issues"
	m.issuesList.DisableQuitKeybindings()

	m.issuesItems = issues
}

func (m JiraModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.NextPanel()
		case "shift+tab":
			m.PrevPanel()
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.columnSize = msg.Width / 12
	case projectsListResponse:
		if msg.err != nil {
			// TODO: Handle this better
			log.Fatalf("Error with projects, err: %v", msg.err)
		}
		m.setProjectListItems(msg.projectsItems)
	case issuesListResponse:
		if msg.err != nil {
			// TODO: Handle this better
			log.Fatalf("Error with projects, err: %v", msg.err)
		}
		m.setIssuesListItems(msg.issuesItems)
	}

	if m.focusedTab == 1 {
		m.projectsList, cmd = m.projectsList.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusedTab == 2 {
		m.issuesList, cmd = m.issuesList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m JiraModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarView(m), contentView(m, 3))
}

func projectSummaryView(m JiraModel, activeIndex int) string {
	content := m.projectName

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize * 3).Render(content)
	}

	return panelStyle.Width(m.columnSize * 3).Render(content)
}

func projectListView(m JiraModel, activeIndex int) string {
	content := m.projectsList.View()

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize * 3).Render(content)
	}

	return panelStyle.Width(m.columnSize * 3).Render(content)
}

func issuesListView(m JiraModel, activeIndex int) string {
	content := m.issuesList.View()

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize * 3).Render(content)
	}

	return panelStyle.Width(m.columnSize * 3).Render(content)
}

func sidebarView(m JiraModel) string {

	projectView := projectSummaryView(m, 1)

	if m.showProjectsList {
		projectView = projectListView(m, 1)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		projectView,
		issuesListView(m, 2),
	)
}

func contentView(m JiraModel, activeIndex int) string {
	content := "Details"

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize * 9).Height(m.height - 2).Render("Details")
	}

	return panelStyle.Width(m.columnSize * 9).Height(m.height - 2).Render(content)
}

// Fetch Data
type projectsListResponse struct {
	projectsItems []jira.Project
	err           error
}

type issuesListResponse struct {
	issuesItems []jira.Issue
	err         error
}

func getIssuesList(c *jira.Client) tea.Cmd {
	return func() tea.Msg {
		issues, err := jira.GetAssignedIssues(c)

		return issuesListResponse{
			issuesItems: issues,
			err:         err,
		}
	}
}

func getProjectsList(c *jira.Client) tea.Cmd {
	return func() tea.Msg {
		projects, err := jira.GetRecentProjects(c)

		return projectsListResponse{
			projectsItems: projects,
			err:           err,
		}
	}
}
