package tui

import (
	"log"
	"net/url"

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
	focusOnList1 = 0
	focusOnList2 = 1
	focusOnList3 = 2
)

func (m *model) createAssignedIssuesList(client *jira.Client, defaults list.DefaultDelegate) *model {
	issues, _ := getAssignedIssues(client)

	var items []list.Item
	for _, issue := range issues {
		items = append(items, issue)
	}

	m.lists[0] = list.New(items, defaults, 30, 15)
	m.lists[0].Title = "Assigned Issues"
	m.lists[0].SetShowHelp(false)
	return m
}

func (m *model) createProjectsList(client *jira.Client, defaults list.DefaultDelegate) *model {
	projects, _ := getRecentProjects(client)

	var items []list.Item
	for _, project := range projects {
		items = append(items, project)
	}

	m.lists[1] = list.New(items, defaults, 30, 15)
	m.lists[1].Title = "Projects"
	m.lists[1].SetShowHelp(false)
	return m
}
func (m *model) createEpicsList(client *jira.Client, defaults list.DefaultDelegate) *model {
	// TODO: Get Epics endpoint
	issues, _ := getAssignedIssues(client)

	var items []list.Item
	for _, issue := range issues {
		items = append(items, issue)
	}

	m.lists[2] = list.New(items, defaults, 30, 15)
	m.lists[2].Title = "Epics"
	m.lists[2].SetShowHelp(false)
	return m
}

func getAssignedIssues(client *jira.Client) ([]jira.Issue, error) {
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

	return resp.Issues, nil
}

func getRecentProjects(client *jira.Client) ([]jira.Project, error) {
	projects, err := jira.GetRecentProjects(client)

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return projects, nil
}
