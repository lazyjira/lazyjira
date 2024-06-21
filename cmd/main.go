package main

import (
	"fmt"
	"log"

	"os"

	"github.com/matthewrobinsondev/lazyjira/pkg/config"
	"github.com/matthewrobinsondev/lazyjira/pkg/jira"
	"github.com/matthewrobinsondev/lazyjira/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	jiraClient := jira.NewClient(cfg)
	p := tea.NewProgram(tui.NewTuiModel(jiraClient))

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start the program: %v\n", err)
		os.Exit(1)
	}
}
