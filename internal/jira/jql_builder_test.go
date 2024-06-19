package jira_test

import (
	"github.com/matthewrobinsondev/lazyjira/internal/jira"
	"testing"
)

func TestJQLBuilder_Equals(t *testing.T) {
	builder := jira.NewJQLBuilder().Equals("assignee", "currentUser()", true)
	expected := `assignee = currentUser()`
	if result := builder.Build(); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestJQLBuilder_In(t *testing.T) {
	builder := jira.NewJQLBuilder().In("status", []string{"Open", "In Progress", "Closed"})
	expected := `status IN ("Open", "In Progress", "Closed")`
	if result := builder.Build(); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestJQLBuilder_ComplexQuery(t *testing.T) {
	builder := jira.NewJQLBuilder().
		Equals("assignee", "currentUser()", true).
		In("status", []string{"Open", "Closed"}).
		Equals("project", "MYPROJECT", false)
	expected := `assignee = currentUser() AND status IN ("Open", "Closed") AND project = "MYPROJECT"`
	if result := builder.Build(); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
