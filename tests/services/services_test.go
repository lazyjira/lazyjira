package services_test

import (
	"errors"
	"github.com/matthewrobinsondev/lazyjira/pkg/services"
	"github.com/matthewrobinsondev/lazyjira/tests"
	"io"
	"net/url"
	"os"
	"testing"
)

const (
	ISSUES_SEARCH_STUB_NAME   = "issues_search.json"
	RECENT_PROJECTS_STUB_NAME = "recent_projects.json"
)

type FakeJiraClient struct {
	MockNewRequest func(method, endpoint string, params url.Values, body io.Reader, version string) ([]byte, error)
}

func (f FakeJiraClient) NewRequest(method, endpoint string, params url.Values, body io.Reader, version string) ([]byte, error) {
	return f.MockNewRequest(method, endpoint, params, body, version)
}

func Test_CanGetAssignedIssues(t *testing.T) {
	mockClient := &FakeJiraClient{
		MockNewRequest: func(method, endpoint string, params url.Values, body io.Reader, version string) ([]byte, error) {
			return os.ReadFile(tests.GetStubPath(ISSUES_SEARCH_STUB_NAME))
		},
	}

	service := services.NewIssuesJiraService(mockClient)

	issues, err := service.GetAssignedIssues()
	if err != nil {
		t.Errorf("Unexpected error getting issues, err: %v", err)
	}

	if len(issues) < 1 {
		t.Errorf("Expected 1 issue, got: %d", len(issues))
	}
}

func Test_CanHandleErrorsFetchingIssues(t *testing.T) {
	mockClient := &FakeJiraClient{
		MockNewRequest: func(method, endpoint string, params url.Values, body io.Reader, version string) ([]byte, error) {
			return []byte(""), errors.New("there was no data")
		},
	}

	service := services.NewIssuesJiraService(mockClient)

	_, err := service.GetAssignedIssues()
	if err == nil {
		t.Error("Unexpected error, got nil")
	}

	if err.Error() != "there was no data" {
		t.Errorf("Unexpected err: %v", err)
	}
}

func Test_GetRecentProjects(t *testing.T) {
	mockClient := &FakeJiraClient{
		MockNewRequest: func(method, endpoint string, params url.Values, body io.Reader, version string) ([]byte, error) {
			return os.ReadFile(tests.GetStubPath(RECENT_PROJECTS_STUB_NAME))
		},
	}

	projectService := services.NewProjectsJiraService(mockClient)

	projects, err := projectService.GetRecentProjects()
	if err != nil {
		t.Errorf("Unexpected error fetching projects, err: %v", err)
	}

	if len(projects) < 1 {
		t.Errorf("expected 1 project to be returned, got: %d", len(projects))
	}
}
