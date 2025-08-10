package tui

import (
	"fmt"
	"log"
	"time"

	help2 "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/matthewrobinsondev/lazyjira/pkg/models"
	"github.com/matthewrobinsondev/lazyjira/pkg/services"
)

const (
	Project = iota + 1
	IssuesList
	Details
)

const (
	DetailsTab = iota
	CommentsTab
)

const (
	COLUMN_COUNT         = 12
	LEFT_PANEL_COLUMNS   = 3
	RIGHT_PANEL_COLUMNS  = 9
	PANEL_PADDING_OFFSET = 6
	LIST_WIDTH_OFFSET    = 10
	MIN_LIST_WIDTH       = 20
	PROJECT_LIST_HEIGHT  = 15
	ISSUES_LIST_HEIGHT   = 30
	PANEL_HEIGHT_OFFSET  = 2
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
	issuesService    services.IssuesService
	projectsService  services.ProjectsService
	width            int
	height           int
	columnSize       int
	selectedIssue    int
	focusedTab       int
	rightPanelTab    int
	maxPanels        int
	showProjectsList bool
	projectName      string
	projectItems     []models.Project
	projectsList     list.Model
	issuesItems      []models.Issue
	issuesList       list.Model
	projectKMap      projectKM
	detailsKMap      detailsKM
}

func NewJiraTui(issuesService services.IssuesService, projectsService services.ProjectsService) JiraModel {

	initProjectsList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	initProjectsList.Title = "Projects"

	initIssuesList := list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	initIssuesList.Title = "Assigned Issues"

	return JiraModel{
		issuesService:    issuesService,
		projectsService:  projectsService,
		focusedTab:       Project,
		rightPanelTab:    DetailsTab,
		maxPanels:        2,
		showProjectsList: false,
		projectName:      "Selected Project Name",
		projectsList:     initProjectsList,
		issuesList:       initIssuesList,
		projectKMap:      projectKMKeys,
		detailsKMap:      detailsKMKeys,
	}
}

func (m JiraModel) Init() tea.Cmd {
	return tea.Batch(
		getProjectsList(m.projectsService),
		getIssuesList(m.issuesService),
	)
}

func (m *JiraModel) NextPanel() {
	if m.focusedTab == Project {
		m.focusedTab = IssuesList
	} else {
		m.focusedTab = Project
	}
}

func (m *JiraModel) PrevPanel() {
	if m.focusedTab == Project {
		m.focusedTab = IssuesList
	} else {
		m.focusedTab = Project
	}
}

func (m *JiraModel) ToggleProjectSwitch() {
	m.showProjectsList = !m.showProjectsList
}

func (m *JiraModel) NextRightPanelTab() {
	m.rightPanelTab = (m.rightPanelTab + 1) % 2
}

func (m *JiraModel) PrevRightPanelTab() {
	if m.rightPanelTab == 0 {
		m.rightPanelTab = 1
	} else {
		m.rightPanelTab = 0
	}
}

func (m *JiraModel) setProjectListItems(projects []models.Project) {
	var items []list.Item
	for _, project := range projects {
		items = append(items, project)
	}

	listWidth := 25
	if m.columnSize > 0 {
		listWidth = m.columnSize*LEFT_PANEL_COLUMNS - LIST_WIDTH_OFFSET
	}
	m.projectsList = list.New(items, list.NewDefaultDelegate(), listWidth, PROJECT_LIST_HEIGHT)
	m.projectsList.Title = "📁 Projects"
	m.projectsList.DisableQuitKeybindings()

	m.projectItems = projects
}

func (m *JiraModel) setIssuesListItems(issues []models.Issue) {
	var items []list.Item
	for _, issue := range issues {
		items = append(items, issue)
	}

	listWidth := 25
	if m.columnSize > 0 {
		listWidth = m.columnSize*LEFT_PANEL_COLUMNS - LIST_WIDTH_OFFSET
	}
	m.issuesList = list.New(items, list.NewDefaultDelegate(), listWidth, ISSUES_LIST_HEIGHT)
	m.issuesList.Title = "🎯 Assigned Issues"
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
		case "p":
			if m.focusedTab == Project {
				m.ToggleProjectSwitch()
			}
		case "[":
			m.PrevRightPanelTab()
		case "]":
			m.NextRightPanelTab()
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.columnSize = msg.Width / COLUMN_COUNT

		listWidth := m.columnSize*LEFT_PANEL_COLUMNS - LIST_WIDTH_OFFSET
		if listWidth < MIN_LIST_WIDTH {
			listWidth = MIN_LIST_WIDTH
		}
		m.projectsList.SetSize(listWidth, PROJECT_LIST_HEIGHT)
		m.issuesList.SetSize(listWidth, ISSUES_LIST_HEIGHT)
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
		return panelStyleFocused.Width(m.columnSize*LEFT_PANEL_COLUMNS - PANEL_PADDING_OFFSET).Render(content)
	}

	return panelStyle.Width(m.columnSize*LEFT_PANEL_COLUMNS - PANEL_PADDING_OFFSET).Render(content)
}

func projectListView(m JiraModel, activeIndex int) string {
	content := m.projectsList.View()

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize*LEFT_PANEL_COLUMNS - PANEL_PADDING_OFFSET).Render(content)
	}

	return panelStyle.Width(m.columnSize*LEFT_PANEL_COLUMNS - PANEL_PADDING_OFFSET).Render(content)
}

func issuesListView(m JiraModel, activeIndex int) string {
	if m.showProjectsList {
		return ""
	}

	content := m.issuesList.View()

	if m.focusedTab == activeIndex {
		return panelStyleFocused.Width(m.columnSize*LEFT_PANEL_COLUMNS - PANEL_PADDING_OFFSET).Render(content)
	}

	return panelStyle.Width(m.columnSize*LEFT_PANEL_COLUMNS - PANEL_PADDING_OFFSET).Render(content)
}

func sidebarView(m JiraModel) string {

	projectView := projectSummaryView(m, Project)

	if m.showProjectsList {
		projectView = projectListView(m, Project)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		projectView,
		issuesListView(m, IssuesList),
	)
}

func contentView(m JiraModel, activeIndex int) string {
	var content string
	var tabsHeader string

	if len(m.issuesItems) > 0 {
		selectedIndex := m.issuesList.Index()
		if selectedIndex >= 0 && selectedIndex < len(m.issuesItems) {
			issue := m.issuesItems[selectedIndex]

			detailsTabStyle := lipgloss.NewStyle().
				Padding(0, 2).
				Background(lipgloss.Color("238")).
				Foreground(lipgloss.Color("241")).
				Border(lipgloss.RoundedBorder(), true, true, false, false)

			commentsTabStyle := lipgloss.NewStyle().
				Padding(0, 2).
				Background(lipgloss.Color("238")).
				Foreground(lipgloss.Color("241")).
				Border(lipgloss.RoundedBorder(), true, true, false, false)

			if m.rightPanelTab == DetailsTab {
				detailsTabStyle = detailsTabStyle.
					Background(lipgloss.Color("69")).
					Foreground(lipgloss.Color("15")).
					Bold(true)
			} else {
				commentsTabStyle = commentsTabStyle.
					Background(lipgloss.Color("69")).
					Foreground(lipgloss.Color("15")).
					Bold(true)
			}

			tabsHeader = lipgloss.JoinHorizontal(lipgloss.Top,
				detailsTabStyle.Render("📋 Details"),
				commentsTabStyle.Render("💬 Comments"),
			)

			var mainContent string
			if m.rightPanelTab == DetailsTab {
				titleSection := lipgloss.JoinVertical(lipgloss.Left,
					lipgloss.NewStyle().
						Foreground(lipgloss.Color("69")).
						Bold(true).
						Render("📋 "+issue.Key),
					lipgloss.NewStyle().
						Foreground(lipgloss.Color("15")).
						Bold(true).
						MarginBottom(1).
						Render(issue.Fields.Summary),
				)

				statusBadge := getStatusBadge(issue.Fields.Status.Name)

				metaSection := lipgloss.JoinHorizontal(lipgloss.Left,
					statusBadge,
					"  ",
				)

				descriptionTitle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("12")).
					Bold(true).
					MarginTop(1).
					MarginBottom(1).
					Render("📝 Description")

				description := issue.GetRenderedDescription()
				if description == "" {
					description = "_No description provided_"
				}
				renderedDescription := formatText(description, m.columnSize*RIGHT_PANEL_COLUMNS)

				mainContent = lipgloss.JoinVertical(lipgloss.Left,
					titleSection,
					"",
					metaSection,
					"",
					descriptionTitle,
					renderedDescription,
				)
			} else {
				if len(issue.Fields.Comment.Comments) == 0 {
					mainContent = lipgloss.JoinVertical(lipgloss.Left,
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("12")).
							Bold(true).
							MarginBottom(1).
							Render("💬 Comments"),
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("241")).
							Italic(true).
							Render("No comments available"),
					)
				} else {
					var commentStrings []string

					commentStrings = append(commentStrings,
						lipgloss.NewStyle().
							Foreground(lipgloss.Color("12")).
							Bold(true).
							MarginBottom(1).
							Render(fmt.Sprintf("💬 Comments (%d)", len(issue.Fields.Comment.Comments))),
						"",
					)

					for i, comment := range issue.Fields.Comment.Comments {
						authorStyle := lipgloss.NewStyle().
							Background(lipgloss.Color("69")).
							Foreground(lipgloss.Color("15")).
							Padding(0, 1).
							Bold(true)

						dateStyle := lipgloss.NewStyle().
							Foreground(lipgloss.Color("241")).
							Italic(true)

						commentHeader := lipgloss.JoinHorizontal(lipgloss.Left,
							authorStyle.Render("👤 "+comment.Author.DisplayName),
							"  ",
							dateStyle.Render(formatDate(comment.Created)),
						)

						commentText := extractCommentText(comment)
						if commentText == "" {
							commentText = "_No content_"
						}
						renderedComment := formatText(commentText, m.columnSize*RIGHT_PANEL_COLUMNS-4)

						commentContainer := lipgloss.NewStyle().
							Border(lipgloss.NormalBorder(), false, false, false, true).
							BorderForeground(lipgloss.Color("241")).
							PaddingLeft(2).
							MarginBottom(1)

						commentStrings = append(commentStrings,
							commentHeader,
							commentContainer.Render(renderedComment),
						)

						if i < len(issue.Fields.Comment.Comments)-1 {
							commentStrings = append(commentStrings, "")
						}
					}

					mainContent = lipgloss.JoinVertical(lipgloss.Left, commentStrings...)
				}
			}

			content = lipgloss.JoinVertical(lipgloss.Left,
				tabsHeader,
				"",
				mainContent,
			)
		}
	} else {
		detailsTabStyle := lipgloss.NewStyle().
			Padding(0, 2).
			Background(lipgloss.Color("69")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Border(lipgloss.RoundedBorder(), true, true, false, false)

		commentsTabStyle := lipgloss.NewStyle().
			Padding(0, 2).
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("241")).
			Border(lipgloss.RoundedBorder(), true, true, false, false)

		tabsHeader = lipgloss.JoinHorizontal(lipgloss.Top,
			detailsTabStyle.Render("📋 Details"),
			commentsTabStyle.Render("💬 Comments"),
		)

		emptyState := lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Bold(true).
				MarginTop(3).
				Render("🔍 No Issue Selected"),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Italic(true).
				Render("Select an issue from the list to view details"),
		)

		content = lipgloss.JoinVertical(lipgloss.Left, tabsHeader, "", emptyState)
	}

	return panelStyle.Width(m.columnSize*RIGHT_PANEL_COLUMNS + 2).Height(m.height - PANEL_HEIGHT_OFFSET).Render(content)
}

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

func formatText(content string, width int) string {
	if content == "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true).Render("No content available")
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(content)
}

func getStatusBadge(status string) string {
	var style lipgloss.Style
	var icon string

	switch status {
	case "Done", "Closed", "Resolved":
		style = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)
		icon = "✅"
	case "In Progress", "In Review", "Testing":
		style = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)
		icon = "🔄"
	case "To Do", "Open", "Backlog":
		style = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)
		icon = "📋"
	default:
		style = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)
		icon = "🔵"
	}

	return style.Render(icon + " " + status)
}

func formatDate(dateStr string) string {
	t, err := time.Parse("2006-01-02T15:04:05.000-0700", dateStr)
	if err != nil {
		return dateStr
	}

	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes < 1 {
			return "just now"
		}
		return fmt.Sprintf("%dm ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	} else {
		return t.Format("Jan 2, 2006")
	}
}

func extractCommentText(comment models.Comment) string {
	var text string
	for _, content := range comment.Body.Content {
		for _, element := range content.Content {
			if element.Type == "text" {
				text += element.Text
			}
		}
	}
	if text == "" {
		return "Comment content not available"
	}
	return text
}
