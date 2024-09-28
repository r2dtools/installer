package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unknwon/com"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall R2DTools agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		systemD.StopService(AGENT_BIN_FILE_NAME)

		if err := systemD.RemoveService(AGENT_BIN_FILE_NAME); err != nil {
			return err
		}

		agentDir := getAgentDir()

		if com.IsDir(agentDir) {
			if err := os.RemoveAll(agentDir); err != nil {
				return err
			}
		}

		removeUserGroup(USER, GROUP)
		logger.Println("the agent is successfully uninstalled")

		return nil
	},
}

func removeUserGroup(userName, groupName string) {
	sh.Exec(fmt.Sprintf("userdel %s", userName))
	sh.Exec(fmt.Sprintf("groupdel %s", groupName))
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
