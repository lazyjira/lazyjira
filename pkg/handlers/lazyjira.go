package handlers

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matthewrobinsondev/lazyjira/pkg/clients"
	"github.com/matthewrobinsondev/lazyjira/pkg/config"
	"github.com/matthewrobinsondev/lazyjira/pkg/services"
	"github.com/matthewrobinsondev/lazyjira/pkg/tui"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func NewLazyJiraHandler(cmd *cobra.Command, args []string) error {
	confServ := config.NewConfigService()
	cfg, err := confServ.Load()
	if err != nil {
		return errors.New(fmt.Sprintf("Error loading configuration: %v", err))
	}

	jiraClient := clients.NewJiraClient(cfg, &http.Client{})

	issuesService := services.NewIssuesJiraService(jiraClient)
	projectsService := services.NewProjectsJiraService(jiraClient)

	p := tea.NewProgram(tui.NewJiraTui(issuesService, projectsService), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	return nil
}
