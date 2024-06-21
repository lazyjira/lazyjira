package tui

import (
	"fmt"
	"log"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
)

type model struct {
	cursor int
	issues []jira.Issue
	chosen map[int]struct{}
	ready  bool
}

func NewTuiModel(client *jira.Client) *model {
	builder := jira.NewJQLBuilder().
		Equals("assignee", "currentUser()", true).
		NotIn("status", []string{"Done", "Closed", "Resolved"})

	jqlQuery := builder.Build()

	fmt.Println(jqlQuery)

	params := url.Values{}
	params.Add("jql", jqlQuery)
	params.Add("fields", "summary,status")

	resp, err := jira.SearchIssues(client, params)

	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	return &model{
		issues: resp.Issues,
		chosen: make(map[int]struct{}),
		ready:  true,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.issues)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			_, ok := m.chosen[m.cursor]
			if ok {
				delete(m.chosen, m.cursor)
			} else {
				m.chosen[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m *model) View() string {
	if !m.ready {
		return "Loading..."
	}
	s := "Jira Issues:\n\n"
	for i, issue := range m.issues {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.chosen[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, issue.Fields.Summary)
	}
	s += "\nPress esc to quit.\n"
	return s
}
