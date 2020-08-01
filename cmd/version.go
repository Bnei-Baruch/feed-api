package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/feed-api/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of feed-api",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Feed and recommendations API for new archive site version %s\n", version.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
