package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func retrieveFileContent(cmd *cobra.Command, args []string) (Config, error) {
	filepath, err := cmd.Flags().GetString("filepath")
	if err != nil {
		return Config{}, fmt.Errorf("failed to get 'filepath' flag")
	}

	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %s", filepath)
	}

	var cfg Config
	if err = yaml.Unmarshal(fileContent, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal file contents")
	}

	return cfg, nil
}

func validateYAML(config Config) error {
	var errMsg []string
	if config.Name == "" {
		errMsg = append(errMsg, "missing name of your LB configuration")
	}

	if config.Description == "" {
		errMsg = append(errMsg, "missing description of your LB configuration")
	}

	if config.Listeners == nil || len(config.Listeners) == 0 {
		errMsg = append(errMsg, "missing or empty listener config")
	} else {
		for i, l := range config.Listeners {
			if l.Name == "" {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Missing name", i+1))
			}

			if l.Protocol == "" {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Missing protocol", i+1))
			} else if l.Protocol != "http" && l.Protocol != "https" {
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Invalid protocl %s", i+1, l.Protocol))
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

	if config.Backends == nil || len(config.Backends) == 0 {
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

	if len(errMsg) == 0 {
		return nil
	}
	errMsg[0] = fmt.Sprintf(" - %s", errMsg[0])
	return fmt.Errorf(strings.Join(errMsg, "\n - "))
}

func Parse(cmd *cobra.Command, args []string) {
	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		fmt.Printf("Error retreiving file content: %v\n", err)
		os.Exit(1)
	}

	if err = validateYAML(lbConfig); err != nil {
		fmt.Printf("Error(s) parsing YAML:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Println("Valid YAML configuration!!")
}
