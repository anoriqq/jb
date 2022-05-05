package cmd

import (
	"fmt"
	"strings"

	"github.com/anoriqq/jb/internal/client"
	"github.com/anoriqq/jb/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// touchCmd represents the touch command
var touchCmd = &cobra.Command{
	Use:   "touch",
	Short: "touch",
	Long:  `Arguments are sent as a message`,
	RunE:  touchRun,
}

func init() {
	rootCmd.AddCommand(touchCmd)
}

func touchRun(_ *cobra.Command, args []string) error {
	c := client.New(config.Cfg)

	resp, err := c.PostChatCommand(config.Cfg.TouchChannel, strings.Join(args, " "))
	if err != nil {
		return err
	}
	if !resp.OK {
		return errors.New(resp.Error)
	}

	fmt.Println("Successfully touch")

	return nil
}
