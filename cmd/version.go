package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "v0.0.0" // default version
)

// versionCmd represents the network command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version for the build",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
