/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func NewQuotaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Manage quotas of Harbor projects",
		Long: `Manage hard storaga quotas of Harbor projects

go-harbor-cli quota get --project-id 1
go-harbor-cli quota put --id 1 --storage-gib 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("quota called")
			return nil
		},
	}

	getCmd := newGetQuotaSubcommand()
	cmd.AddCommand(getCmd)

	putCmd := newSetQuotaSubcommand()
	cmd.AddCommand(putCmd)

	cmd.PersistentFlags().String("name", "", "Project name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newGetQuotaSubcommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get quotas of Harbor projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectId, _ := cmd.Flags().GetInt64("project-id")

			fmt.Printf("\ngetting quotas of project id %s...\n", strconv.FormatInt(projectId, 10))

			client := newHarborClient(cmd)
			quota, err := client.GetQuota(projectId)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}

			quotaBytes, err := json.MarshalIndent(quota, "", "  ")
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("json.MarshalIndent(quota): %w", err)
			}
			fmt.Printf("Quota:\n%s\n", string(quotaBytes))

			return nil
		},
	}

	cmd.Flags().Int64("project-id", 1, "project id")
	cmd.MarkFlagRequired("project-id")

	return cmd
}

func newSetQuotaSubcommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set storage quotas of Harbor projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectId, _ := cmd.Flags().GetInt64("project-id")
			quotaSizeGiB, _ := cmd.Flags().GetInt64("quota-size-gib")
			quotaSizeBytes := quotaSizeGiB * 1024 * 1024 * 1024

			fmt.Printf("\nsetting quotas hard storage limit to prohject %s...\n", strconv.FormatInt(projectId, 10))

			client := newHarborClient(cmd)
			quota, err := client.GetQuota(projectId)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}

			if err := client.PutQuota(quota.Id, map[string]int64{"storage": quotaSizeBytes}); err != nil {
				cmd.SilenceUsage = true
				return err
			}

			return nil
		},
	}

	cmd.Flags().Int64("project-id", 1, "project id")
	cmd.MarkFlagRequired("project-id")
	cmd.Flags().Int64("quota-size-gib", 1, "Max quota size in GiB for the project")
	cmd.MarkFlagRequired("quota-size-gib")

	return cmd
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quotaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quotaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
