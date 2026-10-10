/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"go-harbor-cli/cmd/internal/cli/project"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-harbor-cli",
	Short: "A CLI tool to interact with Harbor API v2.0",
	Long: `A CLI tool to interact with Harbor API v2.0
to manage private projects, quotas, members, etc. For example:

go-harbor-cli project create --name myproj`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(project.NewCommand())

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-harbor-cli.yaml)")
	rootCmd.PersistentFlags().String("username", "default", "harbor username for basic authentication")
	rootCmd.PersistentFlags().String("password", "default", "harbor password for basic authentication")
	rootCmd.PersistentFlags().String("baseurl", "https://localhost:30080", "harbor baseurl")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
