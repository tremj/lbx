package parser

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func SetUpTest(fileContent string) (*cobra.Command, []string, *bytes.Buffer, error) {
	tmpFile, err := os.CreateTemp("", "test-successful-retrieval-*.yaml")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create tmp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(fileContent); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create dummy file content to tmp file")
	}
	tmpFile.Close()

	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.Flags().StringP("filepath", "f", "", "Path to file")

	cmd.Flags().Set("filepath", tmpFile.Name())

	args := []string{}

	return cmd, args, buf, nil
}

func TestNonYAMLFileContentRetreival(t *testing.T) {
	fileContent := `helloLOL`
	cmd, args, buf, err := SetUpTest(fileContent)
	if err != nil {
		t.Fatalf("%v", err)
	}

	_, err = retrieveFileContent(cmd, args)
	if err == nil {
		t.Fatalf("Expected error when retrieving file: %v", fmt.Errorf("failed to unmarshal file contents"))
	}

	if buf.String() != "" {
		t.Errorf("Unexpected output from command")
	}
}
