package jira

import (
	"encoding/json"
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
	Expand string `json:"expand"`
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
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

func SearchIssues(client *Client, params url.Values) (*IssueResponse, error) {
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
