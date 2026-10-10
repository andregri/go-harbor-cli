package cli

import (
	harbor "go-harbor-cli/cmd/internal/client"

	"github.com/spf13/cobra"
)

func newHarborClient(cmd *cobra.Command) *harbor.Client {
	baseurl, _ := cmd.Flags().GetString("baseurl")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	return &harbor.Client{
		BaseUrl:  baseurl,
		Username: username,
		Password: password,
	}
}
