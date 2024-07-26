package models

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

func (p Project) FilterValue() string {
	return p.Name
}
