package jira

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type IssueResponse struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

type Issue struct {
	Expand         string         `json:"expand"`
	ID             string         `json:"id"`
	Key            string         `json:"key"`
	Fields         Fields         `json:"fields"`
	RenderedFields RenderedFields `json:"renderedFields"`
}

type RenderedFields struct {
	Description string `json:"description"`
}

type Fields struct {
	Summary string `json:"summary"`
	Status  Status `json:"status"`
}

type Status struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          string `json:"id"`
}

func (i Issue) Title() string {
	return i.Key + ": " + i.Fields.Summary
}

func (i Issue) Description() string {
	return ""
}

func (i Issue) GetRenderedDescription() string {
	return i.RenderedFields.Description
}

// TODO: Broken
func (i Issue) FilterValue() string {
	return i.Key + " " + i.Fields.Summary
}

func SearchIssues(client ClientInterface, params url.Values) (*IssueResponse, error) {
	apiResp, err := client.NewRequest(http.MethodGet, "/search", params, nil)

	if err != nil {
		return nil, err
	}

	var resp IssueResponse
	err = json.Unmarshal(apiResp, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func GetAssignedIssues(client ClientInterface) ([]Issue, error) {
	// TODO: Add a test for this
	builder := NewJQLBuilder().
		Equals("assignee", "currentUser()", true).
		NotIn("status", []string{"Done", "Closed", "Resolved"})

	jqlQuery := builder.Build()

	params := url.Values{}
	params.Add("jql", jqlQuery)
	params.Add("fields", "summary,status,description")
	params.Add("expand", "renderedFields")
	resp, err := SearchIssues(client, params)

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return resp.Issues, err
}
