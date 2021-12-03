package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/r2dtools/installer/utils"
	"github.com/spf13/cobra"
	"github.com/unknwon/com"
)

const (
	ROOT_USER           = "root"
	USER                = "r2dtools"
	GROUP               = "r2dtools"
	AGENT_BIN_FILE_NAME = "r2dtools"
	LEGO_BIN_FILE_NAME  = "lego"
	AGENT_DIR_NAME      = "r2dtools"
	ARCHIVE_NAME        = "r2dtools-agent.tar.gz"
	URL                 = "https://github.com/r2dtools/agent/releases/download"
	AGENT_PARENT_DIR    = "/opt"
)

var rootCmd = &cobra.Command{
	Use:   "installer",
	Short: "R2DTools agent installer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
var logger *log.Logger
var systemD *utils.SystemD
var sh *utils.SH

// Execute entry point for cli commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func downloadAndUnpackAgent(archiveName, agentParentDir, version string, update bool) error {
	var err error
	if version == "" {
		version, err = utils.GetAgentLatestVersion()
		if err != nil {
			return fmt.Errorf("could not get the latest version of the agent: %v", err)
		}
		if version == "" {
			return errors.New("could not get the latest version of the agent")
		}
		logger.Println(fmt.Sprintf("the latest version of agent: %s", version))
	}

	tmp := os.TempDir()
	filePath := filepath.Join(tmp, archiveName)
	dirPath := filepath.Join(tmp, agentParentDir)

	if com.IsFile(filePath) {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	if com.IsDir(dirPath) {
		if err := os.RemoveAll(dirPath); err != nil {
			return err
		}
	}

	logger.Println("downloading the agent archive ...")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create agent archive file %s: %v", archiveName, err)
	}
	defer file.Close()
	response, err := http.Get(URL + "/v" + version + "/" + archiveName)
	if err != nil {
		return fmt.Errorf("could not download the agent archive: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("could not download the agent archive: bad status: %s", response.Status)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("could not download the agent archive: %v", err)
	}

	logger.Println("unpacking the agent archive ...")
	agentDir := getAgentDir()

	if com.IsDir(agentDir) {
		if !update {
			if err := os.RemoveAll(agentDir); err != nil {
				return fmt.Errorf("could not remove the existed directory '%s': %v", agentDir, err)
			}
		} else {
			entries, err := os.ReadDir(agentDir)
			if err != nil {
				return fmt.Errorf("could not read the agent directory '%s': %v", agentDir, err)
			}
			for _, entry := range entries {
				entryName := entry.Name()
				if !com.IsSliceContainsStr(getDirsToExclude(), entryName) {
					entryPath := filepath.Join(agentDir, entryName)
					if err = os.RemoveAll(entryPath); err != nil {
						return err
					}
				}
			}
		}
	}

	if !com.IsDir(agentDir) {
		if err = os.Mkdir(agentDir, 0755); err != nil {
			return fmt.Errorf("could not create agent directory: %v", err)
		}
	}

	if err := utils.ExtractTarGz(filePath, agentDir); err != nil {
		return fmt.Errorf("could not unpack archive: %v", err)
	}

	return nil
}

func updatePermissions(userName, groupName string) error {
	agentDir := getAgentDir()
	logger.Println(fmt.Sprintf("setting owner '%s:%s' for the agent directory '%s' ...", userName, groupName, agentDir))
	if err := sh.Exec(fmt.Sprintf("chown -R %s:%s %s", userName, groupName, agentDir)); err != nil {
		return fmt.Errorf("could not set agent directory owner: %v", err)
	}

	logger.Println("changing SUID for the agent bin file ...")
	if err := os.Chmod(getAgentBinPath(), 0744|os.ModeSetuid); err != nil {
		return fmt.Errorf("could not set SUID for the agent bin file: %v", err)
	}

	logger.Println("making lego bin file executable...")
	if err := os.Chmod(getLegoBinPath(), 0744); err != nil {
		return fmt.Errorf("could not make lego bin file executable: %v", err)
	}

	return nil
}

func getAgentDir() string {
	return filepath.Join(AGENT_PARENT_DIR, AGENT_DIR_NAME)
}

func getAgentBinPath() string {
	return filepath.Join(getAgentDir(), AGENT_BIN_FILE_NAME)
}

func getLegoBinPath() string {
	return filepath.Join(getAgentDir(), LEGO_BIN_FILE_NAME)
}

func getDirsToExclude() []string {
	return []string{"var", "config"}
}

func init() {
	logger = log.Default()
	sh = &utils.SH{Logger: logger}
	systemD = &utils.SystemD{Logger: logger, Sh: sh}
}
