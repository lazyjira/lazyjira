package jira

import (
	"encoding/json"
	"net/http"
)

type Project struct {
	Self            string          `json:"self"`
	ID              string          `json:"id"`
	Key             string          `json:"key"`
	Name            string          `json:"name"`
	ProjectCategory ProjectCategory `json:"projectCategory"`
	ProjectTypeKey  string          `json:"projectTypeKey"`
}

type ProjectCategory struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p Project) Title() string {
	return p.Key + ": " + p.ProjectCategory.Name
}

func (p Project) Description() string {
	return ""
}

// TODO: Broken
func (p Project) FilterValue() string {
	return p.Key + ": " + p.ProjectCategory.Name
}

func GetRecentProjects(client ClientInterface) ([]Project, error) {
	apiResp, err := client.NewRequest(http.MethodGet, "/project/recent", nil, nil)

	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(apiResp, &projects)

	if err != nil {
		return nil, err
	}

	return projects, nil
}
