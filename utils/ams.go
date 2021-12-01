package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const AMS_URL = "https://ams.r2dtools.com/v1"

func GetAgentLatestVersion() (string, error) {
	response, err := http.Get(AMS_URL + "/agent/latest-version")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", response.Status)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	data := struct{ Version string }{}
	if err = json.Unmarshal(responseData, &data); err != nil {
		return "", err
	}

	return data.Version, nil
}
