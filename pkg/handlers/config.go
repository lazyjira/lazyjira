package handlers

import (
	"github.com/charmbracelet/huh"
	"github.com/matthewrobinsondev/lazyjira/pkg/config"
	"github.com/spf13/cobra"
)

func NewConfigAuthHandlerNewLazyJiraHandler(cmd *cobra.Command, args []string) error {
	configBuf := &config.Config{}

	configForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Whats you jira instance URL?").
				Value(&configBuf.JiraURL),
			huh.NewInput().
				Title("What email do your use with your Access Token?").
				Value(&configBuf.Email),
			huh.NewInput().
				Title("What is your access token?").
				Value(&configBuf.AccessToken),
		),
	).WithTheme(huh.ThemeBase())

	err := configForm.Run()
	if err != nil {
		return err
	}

	confSer := config.NewConfigService()

	if err = confSer.Save(*configBuf); err != nil {
		cobra.CheckErr(err)
	}

	return nil
}
