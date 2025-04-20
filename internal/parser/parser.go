package parser

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func retrieveFileContent(cmd *cobra.Command, args []string) (Config, error) {
	filepath, err := cmd.Flags().GetString("filepath")
	if err != nil {
		return Config{}, err
	}

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = yaml.Unmarshal(fileContent, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func Parse(cmd *cobra.Command, args []string) {
	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
