package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall R2DTools agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := systemD.StopService(AGENT_BIN_FILE_NAME); err != nil {
			return err
		}
		if err := systemD.RemoveService(AGENT_BIN_FILE_NAME); err != nil {
			return err
		}
		removeUserGroup(USER, GROUP)

		return nil
	},
}

func removeUserGroup(userName, groupName string) {
	sh.Exec(fmt.Sprintf("userdel %s", userName))
	sh.Exec(fmt.Sprintf("groupdel %s", userName))
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
