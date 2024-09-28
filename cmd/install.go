package cmd

import (
	"fmt"
	"os/user"

	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
	"github.com/unknwon/com"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install R2DTools agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := host.Info()

		if err != nil {
			return err
		}

		logger.Println("checking platform ...")
		supportedPlatforms := []string{"ubuntu", "centos", "debian"}

		if !com.IsSliceContainsStr(supportedPlatforms, info.Platform) {
			return fmt.Errorf("platform %s is not supported", info.Platform)
		}

		if err = addUserGroup(USER, GROUP); err != nil {
			return err
		}

		if err = downloadAndUnpackAgent(ARCHIVE_NAME, AGENT_DIR_NAME, version, false); err != nil {
			return err
		}

		if err = updatePermissions(USER, GROUP); err != nil {
			return err
		}

		if err = systemD.CreateService(getAgentBinPath(), ROOT_USER, GROUP); err != nil {
			return err
		}

		return systemD.StartService(AGENT_BIN_FILE_NAME)
	},
}

func addUserGroup(userName, groupName string) error {
	logger.Printf("adding group %s ...\n", groupName)
	err := sh.Exec(fmt.Sprintf("grep -q %s /etc/group", groupName))

	if err != nil {
		if err := sh.Exec("groupadd " + groupName); err != nil {
			return fmt.Errorf("could not create group '%s': %v", groupName, err)
		}
	} else {
		logger.Printf("the group '%s' already exists\n", groupName)
	}

	logger.Printf("adding user %s ...\n", userName)
	_, err = user.Lookup(userName)

	if err != nil {
		if _, ok := err.(user.UnknownUserError); !ok {
			return err
		} else {
			if err := sh.Exec(fmt.Sprintf("useradd -g %s %s", groupName, userName)); err != nil {
				return fmt.Errorf("could not create user '%s': %v", userName, err)
			}
		}
	} else {
		logger.Printf("the user '%s' already exists\n", userName)
	}

	return nil
}

func init() {
	installCmd.Flags().StringVar(&version, "version", "", "Version to install. If the version is not specified the latest one will be installed.")
	rootCmd.AddCommand(installCmd)
}
