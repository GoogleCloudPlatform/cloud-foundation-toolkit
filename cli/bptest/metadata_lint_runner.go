package bptest

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RunMetadataLintCommand is the function that will run the lint checks.
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
