package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anoriqq/jb/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jb",
	Short: "jb",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configHome := filepath.Join(home, ".config/jb")
	configName := "config"
	configType := "yml"
	configPath := filepath.Join(configHome, configName+"."+configType)

	cobra.OnInitialize(func() {
		err := os.MkdirAll(configHome, 0777)
		if err != nil {
			panic(err)
		}

		viper.AddConfigPath(configHome)
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)

		viper.AutomaticEnv()

		_, err = os.Stat(configPath)
		if os.IsNotExist(err) {
			if _, err := os.Create(configPath); err != nil {
				panic(err)
			}
		}

		err = config.LoadConfig()
		if err != nil {
			panic(err)
		}
	})
}
