/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package project

import (
	"fmt"

	client_project "go-harbor-cli/cmd/internal/client/project"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage Harbor projects",
		Long: `Manage Harbor projects: create and get projects

create: go-harbor-cli project create --name myproj
get: go-harbor-cli project get --name myproj`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("project called")
			return nil
		},
	}
	createCmd := newCreateProjectCommand()
	cmd.AddCommand(createCmd)

	cmd.PersistentFlags().String("name", "", "Project name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newCreateProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Harbor projects",
		Long: `Create a private Harbor project with the specified name and quota size in GiB.
The project will be created with auto_scan enabled.
Then you can add admin members to the project using the "admin-members" flag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			baseurl, _ := cmd.Flags().GetString("baseurl")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			name, _ := cmd.Flags().GetString("name")
			quota_size_gib, _ := cmd.Flags().GetInt("quota-size-gib")
			// admin_members, _ := cmd.Flags().GetStringSlice("admin-members")

			fmt.Printf("creating project %q...\n", name)

			cProject := &client_project.Client{
				BaseUrl:  baseurl,
				Username: username,
				Password: password,
			}

			p := client_project.Project{
				ProjectName:  name,
				Public:       false,
				Metadata:     client_project.ProjectMetadata{AutoScan: "true"},
				StorageLimit: quota_size_gib * 1024 * 1024 * 1024, // Convert GiB to bytes
			}
			_, err := cProject.Create(p)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}

			return err
		},
	}

	cmd.Flags().Int("quota-size-gib", 1, "Max quota size in GiB for the project")
	cmd.MarkFlagRequired("quota-size-gib")
	cmd.Flags().StringSlice("admin-members", []string{}, "Admin members of the project")
	cmd.MarkFlagRequired("admin-members")

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
