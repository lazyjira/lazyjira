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
	Summary   string       `json:"summary"`
	Status    Status       `json:"status"`
	IssueType IssueType    `json:"issuetype"`
	Comment   CommentField `json:"comment"`
}

type CommentField struct {
	Comments   []Comment `json:"comments"`
	MaxResults int       `json:"maxResults"`
	Total      int       `json:"total"`
	StartAt    int       `json:"startAt"`
}

type IssueType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Comment struct {
	ID      string        `json:"id"`
	Author  CommentAuthor `json:"author"`
	Body    CommentBody   `json:"body"`
	Created string        `json:"created"`
	Updated string        `json:"updated"`
}

type CommentAuthor struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
}

type CommentBody struct {
	Type    string           `json:"type"`
	Version int              `json:"version"`
	Content []CommentContent `json:"content"`
}

type CommentContent struct {
	Type    string                  `json:"type"`
	Content []CommentContentElement `json:"content"`
}

type CommentContentElement struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Status struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          string `json:"id"`
}

func (i Issue) Title() string {
	statusIcon := i.GetStatusIcon()
	return statusIcon + " " + i.Key + ": " + i.Fields.Summary
}

func (i Issue) Description() string {
	return "Type: " + i.Fields.IssueType.Name + " | Status: " + i.Fields.Status.Name
}

func (i Issue) GetStatusIcon() string {
	switch i.Fields.Status.Name {
	case "In Progress", "In Review", "Testing":
		return "🔄"
	case "To Do", "Open", "Backlog":
		return "📋"
	case "Done", "Closed", "Resolved":
		return "✅"
	default:
		return "🔵"
	}
}

func (i Issue) GetRenderedDescription() string {
	return i.RenderedFields.Description
}

func (i Issue) FilterValue() string {
	return i.Fields.Summary
}
