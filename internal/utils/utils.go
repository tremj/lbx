package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tremj/lbx/internal/storage"
)

func RetrieveSaveCmdInfo(cmd *cobra.Command) ([]byte, string, error) {
	filepath, err := cmd.Flags().GetString("filepath")
	if err != nil || filepath == "" {
		return nil, "", errors.New("failed to get 'filepath' flag")
	}

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %s", filepath)
	}

	configName, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, "", fmt.Errorf("error getting name flag: %v", err)
	}

	val, err := storage.Get(cmd.Context(), configName)
	if len(val) == 0 {
		return nil, "", fmt.Errorf("did not find config for key: %s", configName)
	} else if err != nil {
		return nil, "", err
	}

	return fileContent, configName, nil
}

func RetrieveDeleteCmdInfo(cmd *cobra.Command) (string, error) {
	configName, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", fmt.Errorf("error getting name flag: %v", err)
	}

	val, err := storage.Get(cmd.Context(), configName)
	if len(val) == 0 {
		return "", fmt.Errorf("did not find config for key: %s", configName)
	} else if err != nil {
		return "", err
	}

	return configName, nil
}
