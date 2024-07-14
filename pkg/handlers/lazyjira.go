package handlers

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matthewrobinsondev/lazyjira/pkg/config"
	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
	"github.com/matthewrobinsondev/lazyjira/pkg/tui"
	"github.com/spf13/cobra"
	"log"
)

func NewLazyJiraHandler(cmd *cobra.Command, args []string) error {
	confServ := config.NewConfigService()
	cfg, err := confServ.Load()
	if err != nil {
		return errors.New(fmt.Sprintf("Error loading configuration: %v", err))
	}

	jiraClient := jira.NewClient(cfg)

	p := tea.NewProgram(tui.NewJiraTui(jiraClient), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	return nil
}
