package bptest

import (
	"fmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"os"
)

const metadataFile = "metadata.yaml"

// RunLintCommand is the entry function that will run the metadata.yml lint checks.
func RunLintCommand() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Parse medata.yaml to proto
	metadata, err := bpmetadata.UnmarshalMetadata(dir, "/"+metadataFile)
	metadataFile := dir + "/" + metadataFile

	if err != nil {
		fmt.Printf("Error parsing metadata file: %v\n", err)
		os.Exit(1)
	}

	ctx := lintContext{
		metadata: metadata,
		filePath: metadataFile,
	}

	runner := &lintRunner{}
	runner.RegisterRule(&BlueprintConnectionSourceVersionRule{})

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
