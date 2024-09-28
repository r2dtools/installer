package utils

import (
	"log"
	"os/exec"
)

type SH struct {
	Logger *log.Logger
}

func (sh *SH) Exec(command string) error {
	c := exec.Command("/bin/sh", "-c", command)
	output, err := c.CombinedOutput()
	outputStr := string(output)

	if outputStr != "" {
		sh.Logger.Println(outputStr)
	}

	return err
}
