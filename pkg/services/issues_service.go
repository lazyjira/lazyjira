package services

import (
	"encoding/json"
	"github.com/matthewrobinsondev/lazyjira/pkg/clients"
	"github.com/matthewrobinsondev/lazyjira/pkg/models"
	"github.com/matthewrobinsondev/lazyjira/pkg/query"
	"net/http"
	"net/url"
)

type IssuesService interface {
	GetAssignedIssues() ([]models.Issue, error)
}

type IssuesJiraService struct {
	jiraClient clients.JiraClient
}

func NewIssuesJiraService(client clients.JiraClient) *IssuesJiraService {
	return &IssuesJiraService{
		jiraClient: client,
	}
}

func (s IssuesJiraService) GetAssignedIssues() ([]models.Issue, error) {
	builder := query.NewJQLQuery().
		Equals("assignee", "currentUser()", true).
		NotIn("status", []string{"Done", "Closed", "Resolved"})

	params := url.Values{}
	params.Add("jql", builder.Build())
	params.Add("fields", "summary,status,description")
	params.Add("expand", "renderedFields")

	resp, err := s.jiraClient.NewRequest(http.MethodGet, "/search", params, nil, clients.VERSION_3)
	if err != nil {
		return []models.Issue{}, err
	}

	var issues models.IssueResponse
	err = json.Unmarshal(resp, &issues)
	if err != nil {
		return []models.Issue{}, err
	}

	return issues.Issues, nil
}
