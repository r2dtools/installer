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

		if err := downloadAndUnpackAgent(ARCHIVE_NAME, AGENT_DIR_NAME, "", true); err != nil {
			return err
		}

		if err := updatePermissions(USER, GROUP); err != nil {
			return err
		}

		if err := systemD.StartService(AGENT_BIN_FILE_NAME); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
