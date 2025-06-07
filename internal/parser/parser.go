package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tremj/lbx/internal/types"
	"gopkg.in/yaml.v3"
)

func ValidateConfig(data []byte) error {
	var errMsg []string
	var config types.Config

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		errMsg = append(errMsg, err.Error())
	}

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
				errMsg = append(errMsg, fmt.Sprintf("listener %d: Invalid protocol \"%s\"", i+1, l.Protocol))
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

			if b.Port == 0 {
				errMsg = append(errMsg, fmt.Sprintf("backend %d: Missing port", i+1))
			} else if b.Port < 1 || b.Port > 65535 {
				errMsg = append(errMsg, fmt.Sprintf("backend %d: Invalid port %d", i+1, b.Port))
			}
		}
	}

	return concatErrorMessages(errMsg)
}

func concatErrorMessages(messages []string) error {
	return errors.New(strings.Join(messages, "\n"))
}
