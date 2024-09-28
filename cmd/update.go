package cmd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update R2DTools agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := systemD.StopService(AGENT_BIN_FILE_NAME); err != nil {
			return err
		}

		if err := downloadAndUnpackAgent(ARCHIVE_NAME, AGENT_DIR_NAME, version, true); err != nil {
			return err
		}

		if err := updatePermissions(USER, GROUP); err != nil {
			return err
		}

		return systemD.StartService(AGENT_BIN_FILE_NAME)
	},
}

func init() {
	updateCmd.Flags().StringVar(&version, "version", "", "Version to upgrade to. If the version is not specified the agent will be upgraded to the latest one.")
	rootCmd.AddCommand(updateCmd)
}
