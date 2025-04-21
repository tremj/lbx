package parser

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

const baseYAML string = `
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

func TestNonYAMLFileContentRetreival(t *testing.T) {
	fileContent := `helloLOL`
	cmd, args, buf, tmpFile, err := SetUpTest(fileContent, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	_, err = retrieveFileContent(cmd, args)
	if err == nil {
		t.Fatalf("Expected error when retrieving file: %v", ErrUnmarshalFail)
	}

	if err.Error() != ErrUnmarshalFail.Error() {
		t.Fatalf("Expected error message: %v\nGot: %v", ErrUnmarshalFail, err)
	}

	if buf.String() != "" {
		t.Fatalf("Unexpected output from command: %s", buf.String())
	}
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
	cmd, args, buf, tmpFile, err := SetUpTest(baseYAML, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer os.Remove(tmpFile)

	Parse(cmd, args)

	if !bytes.Contains(buf.Bytes(), []byte("Valid YAML configuration!!")) {
		t.Fatalf("Expecting 'Valid YAML configuration!!' message, got: %s", buf.String())
	}
}
