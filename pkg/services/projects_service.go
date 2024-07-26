package services

import (
	"encoding/json"
	"github.com/matthewrobinsondev/lazyjira/pkg/clients"
	"github.com/matthewrobinsondev/lazyjira/pkg/models"
	"net/http"
)

type ProjectsService interface {
	GetRecentProjects() ([]models.Project, error)
}

type ProjectsJiraService struct {
	jiraClient clients.JiraClient
}

func NewProjectsJiraService(client clients.JiraClient) ProjectsService {
	return &ProjectsJiraService{
		jiraClient: client,
	}
}

func (p ProjectsJiraService) GetRecentProjects() ([]models.Project, error) {
	resp, err := p.jiraClient.NewRequest(
		http.MethodGet,
		"/project/recent",
		nil,
		nil,
		clients.VERSION_3,
	)
	if err != nil {
		return []models.Project{}, err
	}

	var projects []models.Project
	err = json.Unmarshal(resp, &projects)
	if err != nil {
		return []models.Project{}, err
	}

	return projects, nil
}
