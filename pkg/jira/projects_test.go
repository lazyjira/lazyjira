package jira_test

import (
	"testing"

	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
)

func TestGetRecentProjects(t *testing.T) {
	mockData := `[
        {
            "expand": "description,lead,issueTypes,url,projectKeys,permissions,insight",
            "self": "https://example.com/rest/api/3/project/10001",
            "id": "10001",
            "key": "PROJ1",
            "name": "Project Alpha",
            "projectCategory": {
                "self": "https://example.com/rest/api/3/projectCategory/1000",
                "id": "1000",
                "name": "Software Development",
                "description": "Projects related to software development."
            }
        },
        {
            "expand": "description,lead,issueTypes,url,projectKeys,permissions,insight",
            "self": "https://example.com/rest/api/3/project/10002",
            "id": "10002",
            "key": "PROJ2",
            "name": "Project Beta",
            "projectCategory": {
                "self": "https://example.com/rest/api/3/projectCategory/1001",
                "id": "1001",
                "name": "Marketing Campaigns",
                "description": "Projects related to marketing and campaign development."
            }
        }
    ]`
	mockClient := &MockClient{
		Response: []byte(mockData),
		Err:      nil,
	}

	projects, err := jira.GetRecentProjects(mockClient)

	if err != nil {
		t.Fatalf("GetRecentProjects returned an error: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}
}
