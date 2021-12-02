package cmd

import (
	"fmt"
	"os/user"

	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
	"github.com/unknwon/com"
)

var version string

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install R2DTools agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := host.Info()
		if err != nil {
			return err
		}
		if info.KernelArch != "x86_64" {
			return fmt.Errorf("architecture %s is not supported", info.KernelArch)
		}

		logger.Println("checking platform ...")
		supportedPlatforms := []string{"ubuntu", "centos", "debian"}
		if !com.IsSliceContainsStr(supportedPlatforms, info.Platform) {
			return fmt.Errorf("platform %s is not supported", info.Platform)
		}

		if err = installPackages(info); err != nil {
			return err
		}

		if err = addUserGroup(USER, GROUP); err != nil {
			return err
		}

		if err = downloadAndUnpackAgent(ARCHIVE_NAME, AGENT_DIR_NAME, version); err != nil {
			return err
		}

		if err = updatePermissions(USER, GROUP); err != nil {
			return err
		}

		if err = systemD.CreateService(getAgentBinPath(), ROOT_USER, GROUP); err != nil {
			return err
		}

		if err = systemD.StartService(AGENT_BIN_FILE_NAME); err != nil {
			return err
		}

		return nil
	},
}

func installPackages(info *host.InfoStat) error {
	logger.Println("installing augeas package ...")
	packageCmd := "sudo apt-get install libaugeas0"
	if info.PlatformFamily == "rhel" {
		packageCmd = "yum -y install augeas"
	}

	if err := sh.Exec(packageCmd); err != nil {
		return fmt.Errorf("could not install augeas package: %v", err)
	}

	return nil
}

func addUserGroup(userName, groupName string) error {
	logger.Println(fmt.Sprintf("adding group %s ...", groupName))
	err := sh.Exec(fmt.Sprintf("grep -q %s /etc/group", groupName))
	if err != nil {
		if err := sh.Exec("groupadd " + groupName); err != nil {
			return fmt.Errorf("could not create group '%s': %v", groupName, err)
		}
	} else {
		logger.Println(fmt.Sprintf("the group '%s' already exists", groupName))
	}

	logger.Println(fmt.Sprintf("adding user %s ...", userName))
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
		logger.Println(fmt.Sprintf("the user '%s' already exists", userName))
	}

	return nil
}

func init() {
	installCmd.Flags().StringVar(&version, "version", "", "Version to install. If the version is not specified the latest one will be installed.")
	rootCmd.AddCommand(installCmd)
}
