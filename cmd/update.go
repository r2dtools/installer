package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unknwon/com"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update R2DTools agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := getConfigContent()
		if err != nil {
			return err
		}

		if err := systemD.StopService(AGENT_BIN_FILE_NAME); err != nil {
			return err
		}

		if err := downloadAndUnpackAgent(ARCHIVE_NAME, AGENT_DIR_NAME, ""); err != nil {
			return err
		}

		if err := putConfigContent(config); err != nil {
			return nil
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

func getConfigContent() (string, error) {
	configPath := getConfigFilePath()
	if !com.IsFile(configPath) {
		return "", nil
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("could not read config file '%s': %v", configPath, err)
	}

	return string(content), nil
}

func putConfigContent(content string) error {
	configPath := getConfigFilePath()
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("could not update config '%s': %v", configPath, err)
	}

	return nil
}

func getConfigFilePath() string {
	return filepath.Join(getAgentDir(), "config", "params.yaml")
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
