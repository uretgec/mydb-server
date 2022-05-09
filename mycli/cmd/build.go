package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build services and run",
	Long:  `build services and run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("build called")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
