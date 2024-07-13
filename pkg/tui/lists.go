package tui

import (
	"log"
	"net/url"
	"sync"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
)

type focusState int

var (
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	focusedStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("62"))
	unfocusedStyle = lipgloss.NewStyle().Padding(1)
)

const (
	IssuesList   = 0
	ProjectsList = 1
	EpicsList    = 2
)

func (m *model) createAssignedIssuesList(defaults list.DefaultDelegate, wg *sync.WaitGroup) *model {
	defer wg.Done()
	issues := getAssignedIssues(m.client)

	var items []list.Item
	for _, issue := range issues {
		items = append(items, issue)
	}

	m.lists[IssuesList] = list.New(items, defaults, 30, 15)
	m.lists[IssuesList].Title = "Assigned Issues"
	m.lists[IssuesList].SetShowHelp(false)
	return m
}

func (m *model) createProjectsList(defaults list.DefaultDelegate, wg *sync.WaitGroup) *model {
	defer wg.Done()
	projects := getRecentProjects(m.client)

	var items []list.Item
	for _, project := range projects {
		items = append(items, project)
	}

	m.lists[ProjectsList] = list.New(items, defaults, 30, 15)
	m.lists[ProjectsList].Title = "Projects"
	m.lists[ProjectsList].SetShowHelp(false)
	return m
}
func (m *model) createEpicsList(defaults list.DefaultDelegate, wg *sync.WaitGroup) *model {
	defer wg.Done()
	// TODO: Get Epics endpoint
	issues := getAssignedIssues(m.client)

	var items []list.Item
	for _, issue := range issues {
		items = append(items, issue)
	}

	m.lists[EpicsList] = list.New(items, defaults, 30, 15)
	m.lists[EpicsList].Title = "Epics"
	m.lists[EpicsList].SetShowHelp(false)
	return m
}

func getAssignedIssues(client *jira.Client) []jira.Issue {
	builder := jira.NewJQLBuilder().
		Equals("assignee", "currentUser()", true).
		NotIn("status", []string{"Done", "Closed", "Resolved"})

	jqlQuery := builder.Build()

	params := url.Values{}
	params.Add("jql", jqlQuery)
	params.Add("fields", "summary,status,description")
	params.Add("expand", "renderedFields")
	resp, err := jira.SearchIssues(client, params)

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return resp.Issues
}

func getRecentProjects(client *jira.Client) []jira.Project {
	projects, err := jira.GetRecentProjects(client)

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return projects
}
