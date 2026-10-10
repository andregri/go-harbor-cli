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
			projectName, _ := cmd.Flags().GetString("name")
			quota_size_gib, _ := cmd.Flags().GetInt64("quota-size-gib")
			quotaSizeBytes := quota_size_gib * 1024 * 1024 * 1024 // Convert GiB to bytes
			admin_users, _ := cmd.Flags().GetStringSlice("admin-users")

			fmt.Printf("\ncreating project %q...\n", projectName)

			p := harbor.Project{
				ProjectName:  projectName,
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

			for _, username := range admin_users {
				// m := &harbor.Member{
				// 	RoleId:
				// 	User: harbor.MemberUser{
				// 		Username: username,
				// 	},
				// }
				m := map[string]any{
					"role_id": harbor.USER_MEMBER_ADMIN,
					"member_user": map[string]any{
						"username": username,
					},
				}
				err = client.AddProjectMember(projectName, m)
				if err != nil {
					cmd.SilenceUsage = true
					return err
				}

				fmt.Printf("Added admin user %q to project %q\n", username, projectName)
			}

			return nil
		},
	}

	cmd.Flags().String("name", "", "Project name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int64("quota-size-gib", 1, "Max quota size in GiB for the project")
	cmd.MarkFlagRequired("quota-size-gib")
	cmd.Flags().StringSlice("admin-users", []string{}, "Admin users of the project")
	cmd.MarkFlagRequired("admin-users")

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
