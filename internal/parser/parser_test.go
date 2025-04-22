package parser

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func SetUpTest(fileContent string, setupFilepath bool) (*cobra.Command, []string, *bytes.Buffer, string, error) {
	tmpFile, err := os.CreateTemp("", "test-successful-retrieval-*.yaml")
	if err != nil {
		return nil, nil, nil, "", fmt.Errorf("failed to create tmp file: %v", err)
	}

	if _, err = tmpFile.WriteString(fileContent); err != nil {
		return nil, nil, nil, "", fmt.Errorf("failed to create dummy file content to tmp file")
	}
	tmpFile.Close()

	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)

	cmd.SetOut(buf)
	cmd.Flags().StringP("filepath", "f", "", "Path to file")
	if setupFilepath {
		cmd.Flags().Set("filepath", tmpFile.Name())
	}

	args := []string{}

	return cmd, args, buf, tmpFile.Name(), nil
}

func TestNoFileExists(t *testing.T) {
	fileContent := `g2`
	cmd, args, buf, tmpFile, err := SetUpTest(fileContent, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = os.Remove(tmpFile) // remove file before running command
	if err != nil {
		t.Fatalf("Error removing file: %v", err)
	}

	expectedMsg := fmt.Sprintf("failed to read file: %s", tmpFile)
	_, err = retrieveFileContent(cmd, args)
	if err == nil {
		t.Fatalf("Expected error when retrieving file: %s", expectedMsg)
	}

	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s\nGot: %v", expectedMsg, err)
	}

	if buf.String() != "" {
		t.Fatalf("Unexpected output from command: %s", buf.String())
	}

}

func TestNoFilepathFlag(t *testing.T) {
	fileContent := `yogurt`
	cmd, args, buf, tmpFile, err := SetUpTest(fileContent, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	_, err = retrieveFileContent(cmd, args)
	if err == nil {
		t.Fatalf("Expected error when retreiving file: %v", ErrNoFilepathFlag)
	}

	if err.Error() != ErrNoFilepathFlag.Error() {
		t.Fatalf("Expected error message: %v\nGot: %v", ErrNoFilepathFlag, err)
	}

	if buf.String() != "" {
		t.Fatalf("Unexpected output from command: %s", buf.String())
	}
}

func TestValidYAML(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, buf, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	Parse(cmd, args)

	if !bytes.Contains(buf.Bytes(), []byte("Valid YAML configuration!!")) {
		t.Fatalf("Expecting 'Valid YAML configuration!!' message, got: %s", buf.String())
	}
}

func TestUnknownYAMLField(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

hello: hiiii

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	_, err = retrieveFileContent(cmd, args)
	expectedSection := "field hello not found in type parser.Config"
	if !strings.Contains(err.Error(), expectedSection) {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestMissingNameDescription(t *testing.T) {
	config := `
listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"missing name of your LB configuration", "missing description of your LB configuration"}
	if len(msgArr) != len(expectedMsg) {
		t.Fatalf("Expected %d errors, got %d errors", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if msgArr[i] != expectedMsg[i] {
			t.Fatalf("Expected '%s'\nGot '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestMissingListener(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := "missing or empty listener config"

	if len(msgArr) != 1 {
		t.Fatalf("Expected 1 error, got %d errors", len(msgArr))
	}

	if expectedMsg != msgArr[0] {
		t.Fatalf("Expected '%s'\nGot '%s'", expectedMsg, msgArr[0])
	}
}

func TestEmptyListener(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"missing or empty listener config"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestListenerMissingNamePort(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - protocol: http

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Error retrieving file content: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"listener 1: Missing name", "listener 1: Missing port number"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestListenerMissingProtocol(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Error retrieving file content: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"listener 1: Missing protocol"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestListenerWrongHTTPPort(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 81

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Error retrieving file content: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"listener 1: Invalid port 81 for protocol \"http\""}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestListenerWrongHTTPSPortAndMissingTLS(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-https
    protocol: https
    port: 81

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Error retrieving file content: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"listener 1: Invalid port 81 for protocol \"https\"", "listener 1: Missing certificate information"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestListenerInvalidProtocol(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-ssh
    protocol: ssh
    port: 81

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Error retrieving file content: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"listener 1: Invalid protocol \"ssh\""}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))

	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestListenerValidTLS(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-https
    protocol: https
    port: 443
    tls_cert: "/path/to/cert.pem"
    tls_key: "/path/to/key.pem"

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Error retrieving file content: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))

	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestMissingBackend(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"missing or empty backend config"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestEmptyBackend(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"missing or empty backend config"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestBackendMissingNamePort(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - port: 8080
  - name: backend2
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"backend 1: Missing name", "backend 2: Missing port"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}

func TestBackendInvalidPort(t *testing.T) {
	config := `
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend9
    port: 100000
`
	cmd, args, _, tmpFile, err := SetUpTest(config, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	lbConfig, err := retrieveFileContent(cmd, args)
	if err != nil {
		t.Fatalf("Error retreiving file content: %v", err)
	}

	msgArr := validateYAML(lbConfig)
	expectedMsg := []string{"backend 1: Invalid port 100000"}

	if len(expectedMsg) != len(msgArr) {
		t.Fatalf("Expected %d errors, got %d", len(expectedMsg), len(msgArr))
	}

	for i := range expectedMsg {
		if expectedMsg[i] != msgArr[i] {
			t.Fatalf("Expected '%s'\nGot: '%s'", expectedMsg[i], msgArr[i])
		}
	}
}
