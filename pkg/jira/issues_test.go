package jira_test

import (
	"io"
	"net/url"
	"testing"

	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
)

type MockClient struct {
	Response []byte
	Err      error
}

func (m *MockClient) NewRequest(method, endpoint string, params url.Values, body io.Reader) ([]byte, error) {
	return m.Response, m.Err
}

func TestSearchIssues(t *testing.T) {
	mockData := `{"expand":"names","issues":[{"id":"1","key":"TEST-1","fields":{"summary":"Test issue","status":{"name":"Open"}}}]}`
	mockClient := &MockClient{
		Response: []byte(mockData),
		Err:      nil,
	}

	params := url.Values{}
	params.Add("jql", "assignee = currentUser()")
	params.Add("fields", "summary,status")

	resp, err := jira.SearchIssues(mockClient, params)

	if err != nil {
		t.Fatalf("SearchIssues returned an error: %v", err)
	}

	if len(resp.Issues) != 1 || resp.Issues[0].Key != "TEST-1" {
		t.Errorf("Unexpected response data: %+v", resp.Issues)
	}
}
