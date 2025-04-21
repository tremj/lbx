package parser

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ErrNoFilepathFlag = errors.New("failed to get 'filepath' flag")
	ErrUnmarshalFail  = errors.New("failed to unmarshal file contents")
)

func retrieveFileContent(cmd *cobra.Command, _ []string) (Config, error) {
	filepath, err := cmd.Flags().GetString("filepath")
	if err != nil || filepath == "" {
		return Config{}, ErrNoFilepathFlag
	}

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %s", filepath)
	}

	var cfg Config
	if err = yaml.Unmarshal(fileContent, &cfg); err != nil {
		return Config{}, ErrUnmarshalFail
	}

	return cfg, nil
}

func validateYAML(config Config) []string {
	var errMsg []string
	if config.Name == "" {
		errMsg = append(errMsg, "missing name of your LB configuration")
	}

	if config.Description == "" {
		errMsg = append(errMsg, "missing description of your LB configuration")
	}

	if len(config.Listeners) == 0 {
		errMsg = append(errMsg, "missing or empty listener config")
	} else {
		for i, l := range config.Listeners {
			if l.Name == "" {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Missing name", i+1))
			}

			if l.Protocol == "" {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Missing protocol", i+1))
			} else if l.Protocol != "http" && l.Protocol != "https" {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Invalid protocol %s", i+1, l.Protocol))
			}

			if l.Port == 0 {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Missing port number", i+1))
			} else if (l.Protocol == "http" && l.Port != 80) || (l.Protocol == "https" && l.Port != 443) {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Invalid port %d for protocol \"%s\"", i+1, l.Port, l.Protocol))
			}

			if l.Protocol == "https" && (l.TLSCert == "" || l.TLSKey == "") {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Missing certificate information", i+1))
			}
		}
	}

	if len(config.Backends) == 0 {
		errMsg = append(errMsg, "missing or empty backend config")
	} else {
		for i, b := range config.Backends {
			if b.Name == "" {
				errMsg = append(errMsg, fmt.Sprintf("backend %d: Missing name", i+1))
			}

			if b.Port < 0 || b.Port > 65535 {
				errMsg = append(errMsg, fmt.Sprintf("backend %d: Invalid port %d", i+1, b.Port))
			}
		}
	}

	return errMsg
}

func Parse(cmd *cobra.Command, args []string) {
	if err := parse(cmd, args); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Valid YAML configuration!!")
}

func parse(cmd *cobra.Command, args []string) error {
	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		return fmt.Errorf("error retreiving file content: %v", err)
	}

	if errArr := validateYAML(lbConfig); len(errArr) > 0 {
		errArr[0] = fmt.Sprintf(" - %s", errArr[0])
		return fmt.Errorf("error(s) parsing YAML:\n%v", strings.Join(errArr, "\n - "))
	}

	return nil
}
