package bptest

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// RunMetadataLintCommand is the entry function that will run the metadata.yml lint checks.
func RunMetadataLintCommand(cmd *cobra.Command, args []string) {
	metadataFile := findMetadataFile()
	if metadataFile == "" {
		fmt.Println("Error: metadata.yaml file not found")
		os.Exit(1)
	}

	metadata, err := ParseMetadataFile(metadataFile)
	if err != nil {
		fmt.Printf("Error parsing metadata file: %v\n", err)
		os.Exit(1)
	}

	ctx := LintContext{
		Metadata: metadata,
		FilePath: metadataFile,
	}

	// Create a lint runner and register all lint rules
	runner := &LintRunner{}
	runner.RegisterRule(&BlueprintSourceRule{})
	runner.RegisterRule(&BlueprintVersionRule{})
	runner.RegisterRule(&EmptyRule{}) // Example empty rule

	// Run lint checks
	errs := runner.Run(ctx)
	if len(errs) > 0 {
		fmt.Println("Linting failed with the following errors:")
		for _, err := range errs {
			fmt.Println("- ", err)
		}
		os.Exit(1)
	} else {
		fmt.Println("All lint checks passed!")
	}
}

// findMetadataFile searches for 'metadata.yaml' in the current directory or parent directories.
func findMetadataFile() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	for {
		metadataFilePath := dir + "/metadata.yaml"
		if _, err := os.Stat(metadataFilePath); err == nil {
			return metadataFilePath
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return ""
}

// ParseMetadataFile reads a YAML file, converts it to JSON, and unmarshals it into a proto message.
func ParseMetadataFile(filePath string) (*BlueprintMetadata, error) {
	// Read the YAML file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	// Unmarshal YAML into a map (intermediate structure)
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Convert the YAML map to JSON
	jsonData, err := yamlToJSON(yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}

	// Unmarshal JSON into the proto message
	var blueprintMetadata BlueprintMetadata
	if err := protojson.Unmarshal(jsonData, &blueprintMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON into proto: %w", err)
	}

	return &blueprintMetadata, nil
}

// yamlToJSON is a helper function that converts a YAML map to JSON.
func yamlToJSON(yamlData map[string]interface{}) ([]byte, error) {
	jsonData, err := protojson.Marshal(yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return jsonData, nil
}
