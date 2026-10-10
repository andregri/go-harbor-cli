/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"strconv"

	harbor "go-harbor-cli/cmd/internal/client"

	"github.com/spf13/cobra"
)

func NewProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage Harbor projects",
		Long: `Manage Harbor projects: create and get projects

create: go-harbor-cli project create --name myproj
get: go-harbor-cli project get --name myproj`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("projct called")
			return nil
		},
	}

	createCmd := newCreateProjectSubcommand()
	cmd.AddCommand(createCmd)

	return cmd
}

func newCreateProjectSubcommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Harbor projects",
		Long:  `Create a private Harbor project with the specified name and auto_scan enabled.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			quota_size_gib, _ := cmd.Flags().GetInt64("quota-size-gib")
			quotaSizeBytes := quota_size_gib * 1024 * 1024 * 1024 // Convert GiB to bytes
			// admin_members, _ := cmd.Flags().GetStringSlice("admin-members")

			fmt.Printf("\ncreating project %q...\n", name)

			p := harbor.Project{
				ProjectName:  name,
				Public:       false,
				Metadata:     harbor.ProjectMetadata{AutoScan: "true"},
				StorageLimit: quotaSizeBytes,
			}

			client := newHarborClient(cmd)
			projectId, err := client.Create(p)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}

			fmt.Printf("Created private project %s\n", strconv.FormatInt(projectId, 10))

			return nil
		},
	}

	cmd.Flags().String("name", "", "Project name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int64("quota-size-gib", 1, "Max quota size in GiB for the project")
	cmd.MarkFlagRequired("quota-size-gib")
	// cmd.Flags().StringSlice("admin-members", []string{}, "Admin members of the project")
	// cmd.MarkFlagRequired("admin-members")

	return cmd
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
