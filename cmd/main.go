package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/matthewrobinsondev/lazyjira/internal/config"
	"github.com/matthewrobinsondev/lazyjira/internal/jira"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	jiraClient := jira.NewClient(cfg)
	builder := jira.NewJQLBuilder().
		Equals("assignee", "currentUser()", true).
		In("status", []string{"Done", "Closed", "Resolved"})

	jqlQuery := builder.Build()

	fmt.Println(jqlQuery)

	params := url.Values{}
	params.Add("jql", jqlQuery)
	params.Add("fields", "summary,status")

	fmt.Println(jiraClient.NewRequest(http.MethodGet, "/search", params, nil))
}
