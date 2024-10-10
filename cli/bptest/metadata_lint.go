package bptest

import (
	"fmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

// RunMetadataLintCommand is the entry function that will run the metadata.yml lint checks.
func RunMetadataLintCommand() {
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

	runner := &LintRunner{}
	runner.RegisterRule(&BlueprintVersionRule{})

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
func ParseMetadataFile(filePath string) (*bpmetadata.BlueprintMetadata, error) {
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

	// Convert YAML to Protobuf
	blueprintMetadata, err := ConvertYAMLToProto(yamlData)
	if err != nil {
		fmt.Printf("Error converting YAML to Proto: %v\n", err)
		return nil, err
	}

	return blueprintMetadata, nil
}

// Function to convert YAML to Protobuf
func ConvertYAMLToProto(yamlData map[string]interface{}) (*bpmetadata.BlueprintMetadata, error) {
	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %v", err)
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %v", err)
	}

	var protoMessage bpmetadata.BlueprintMetadata
	if err := protojson.Unmarshal(jsonBytes, &protoMessage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to Proto: %v", err)
	}

	return &protoMessage, nil
}
