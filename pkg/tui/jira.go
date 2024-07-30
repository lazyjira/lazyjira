package tui

import (
	"fmt"
	help2 "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/matthewrobinsondev/lazyjira/pkg/models"
	"github.com/matthewrobinsondev/lazyjira/pkg/services"
	"log"
	"strings"
)

const (
	Project = iota + 1
	IssuesList
	Details
)

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("241")).
			Padding(0, 1)

	panelStyleFocused = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69")).
				Padding(0, 1)
)

type JiraModel struct {
	issuesService    services.IssuesService
	projectsService  services.ProjectsService
	width            int
	height           int
	columnSize       int
	selectedIssue    models.Issue
	focusedTab       int
	maxPanels        int
	showProjectsList bool
	projectName      string
	projectItems     []models.Project
	projectsList     list.Model
	issuesItems      []models.Issue
	issuesList       list.Model
	projectKMap      projectKM
}

func NewJiraTui(issuesService services.IssuesService, projectsService services.ProjectsService) JiraModel {

	initProjectsList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	initProjectsList.Title = "Projects"

	initIssuesList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	initIssuesList.Title = "Assigned Issues"

	return JiraModel{
		issuesService:    issuesService,
		projectsService:  projectsService,
		focusedTab:       1,
		maxPanels:        3,
		showProjectsList: false,
		projectName:      "Selected Project Name",
		projectsList:     initProjectsList,
		issuesList:       initIssuesList,
		projectKMap:      projectKMKeys,
	}
}

func (m JiraModel) Init() tea.Cmd {
	return tea.Batch(
		getProjectsList(m.projectsService),
		getIssuesList(m.issuesService),
	)
}

func (m *JiraModel) NextPanel() {
	m.focusedTab = (m.focusedTab + 1) % 4
}

func (m *JiraModel) PrevPanel() {
	m.focusedTab = (m.focusedTab + 2) % 4
}

func (m *JiraModel) ToggleProjectSwitch() {
	m.showProjectsList = !m.showProjectsList
}

func (m *JiraModel) setProjectListItems(projects []models.Project) {
	var items []list.Item
	for _, project := range projects {
		items = append(items, project)
	}

	m.projectsList = list.New(items, list.NewDefaultDelegate(), 30, 15)
	m.projectsList.Title = "Projects"
	m.projectsList.DisableQuitKeybindings()

	m.projectItems = projects
}

func (m *JiraModel) setIssuesListItems(issues []models.Issue) {
	var items []list.Item
	for _, issue := range issues {
		items = append(items, issue)
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
		case "ctrl+s":
			if m.focusedTab == Project {
				m.ToggleProjectSwitch()
			}
		case "enter":
			if m.focusedTab == IssuesList {
				m.selectedIssue = m.issuesList.SelectedItem().(models.Issue)
			}
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

	if m.focusedTab == Project {
		m.projectsList, cmd = m.projectsList.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusedTab == IssuesList {
		m.issuesList, cmd = m.issuesList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m JiraModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarView(m), contentView(m, 3))
}

func projectSummaryView(m JiraModel, activeIndex int) string {
	help := help2.New()
	summaryContainer := lipgloss.NewStyle().Padding(0, 1)
	content := summaryContainer.Render(lipgloss.JoinVertical(lipgloss.Left, m.projectName, help.View(m.projectKMap)))

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
	if m.showProjectsList {
		return ""
	}

	content := m.issuesList.View()

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize * 3).Render(content)
	}

	return panelStyle.Width(m.columnSize * 3).Render(content)
}

func sidebarView(m JiraModel) string {

	projectView := projectSummaryView(m, Project)

	if m.showProjectsList {
		projectView = projectListView(m, Project)
	}

	return lipgloss.JoinVertical(
		lipgloss.Right,
		projectView,
		issuesListView(m, IssuesList),
	)
}

func contentView(m JiraModel, activeIndex int) string {
	var output strings.Builder

	if (models.Issue{} != m.selectedIssue) {
		output.Write([]byte(renderIssueDescription(m.selectedIssue)))
	}

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize * 9).Height(m.height - 2).Render(output.String())
	}

	return panelStyle.Width(m.columnSize * 9).Height(m.height - 2).Render(output.String())
}

func renderIssueDescription(i models.Issue) string {
	var output strings.Builder

	output.Write([]byte(fmt.Sprintf("# %s", i.Fields.Summary)))
	output.Write([]byte(fmt.Sprintf("\n## %s", i.Fields.Status.Name)))
	output.Write([]byte(fmt.Sprintf("\n## %s", i.Key)))
	output.Write([]byte(fmt.Sprintf("\n %s", i.RenderedFields.Description)))

	out, _ := glamour.Render(output.String(), "dark")
	return out
}

// Fetch Data
type projectsListResponse struct {
	projectsItems []models.Project
	err           error
}

type issuesListResponse struct {
	issuesItems []models.Issue
	err         error
}

func getIssuesList(i services.IssuesService) tea.Cmd {
	return func() tea.Msg {
		issues, err := i.GetAssignedIssues()

		return issuesListResponse{
			issuesItems: issues,
			err:         err,
		}
	}
}

func getProjectsList(p services.ProjectsService) tea.Cmd {
	return func() tea.Msg {
		projects, err := p.GetRecentProjects()

		return projectsListResponse{
			projectsItems: projects,
			err:           err,
		}
	}
}
