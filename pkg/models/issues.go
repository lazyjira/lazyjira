package models

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

func (i Issue) FilterValue() string {
	return i.Fields.Summary
}
