package query_test

import (
	"github.com/matthewrobinsondev/lazyjira/pkg/query"
	"testing"
)

func TestJQLQuery_Equals(t *testing.T) {
	expected := "assignee = currentUser()"

	builder := query.NewJQLQuery().
		Equals("assignee", "currentUser()", true)

	if result := builder.Build(); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestJQLQuery_In(t *testing.T) {
	expected := `status IN ("Open", "In Progress")`

	builder := query.NewJQLQuery().
		In("status", []string{"Open", "In Progress"})

	if result := builder.Build(); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestJQLQuery_NotIn(t *testing.T) {
	expected := `status NOT IN ("Open", "In Progress")`

	builder := query.NewJQLQuery().
		NotIn("status", []string{"Open", "In Progress"})

	if result := builder.Build(); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestJQLQuery_ChainedQuery(t *testing.T) {
	expected := `assignee = currentUser() AND status IN ("Open")`

	builder := query.NewJQLQuery().
		Equals("assignee", "currentUser()", true).
		In("status", []string{"Open"})

	if result := builder.Build(); result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
