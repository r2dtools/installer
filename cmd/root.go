package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "installer",
	Short: "R2DTools agent installer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute entry point for cli commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
