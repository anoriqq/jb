package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/anoriqq/jb/internal/client"
	"github.com/anoriqq/jb/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login",
	RunE:  loginRun,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func loginRun(_ *cobra.Command, _ []string) error {
	if len(config.Cfg.D) > 0 && len(config.Cfg.Token) > 0 {
		c := client.New(config.Cfg)

		resp, err := c.PostAuthTest()
		if err != nil {
			return err
		}

		if resp.OK {
			var confirm bool
			p := &survey.Confirm{
				Message: fmt.Sprintf("You are already logged in to %s as %v. Continue re-login?", resp.Team, resp.User),
			}
			err := survey.AskOne(p, &confirm)
			if err != nil {
				return errors.WithStack(err)
			}

			if !confirm {
				return nil
			}
		}
	}

	d, err := surveyD()
	if err != nil {
		return err
	}

	viper.Set("d", d)

	token, err := surveyToken()
	if err != nil {
		return err
	}

	viper.Set("token", token)

	err = config.LoadConfig()
	if err != nil {
		return err
	}

	c := client.New(config.Cfg)

	resp, err := c.PostAuthTest()
	if err != nil {
		return err
	}

	if !resp.OK {
		return errors.New(resp.Error)
	}

	touchChannel, err := surveyTouchChannel()
	if err != nil {
		return err
	}

	viper.Set("touchChannel", touchChannel)

	err = viper.WriteConfig()
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("Successfully login")

	return nil
}

func surveyToken() (string, error) {
	var token string
	p := &survey.Input{Message: "Please enter token. Token is in `JSON.parse(localStorage.localConfig_v2).teams`"}
	err := survey.AskOne(p, &token, survey.WithValidator(survey.Required))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

func surveyD() (string, error) {
	var d string
	p := &survey.Input{Message: "Please enter d. D is in Application > Cookie > https://app.slack.com > d"}
	err := survey.AskOne(p, &d, survey.WithValidator(survey.Required))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return d, nil
}

func surveyTouchChannel() (string, error) {
	var touchChannel string
	p := &survey.Input{Message: "Please enter channel ID for touch"}
	err := survey.AskOne(p, &touchChannel, survey.WithValidator(survey.Required))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return touchChannel, nil
}
