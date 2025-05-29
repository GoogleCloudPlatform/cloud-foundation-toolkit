package adcbpmetadata

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
)



const mdTemplate = `
# {{.Title}}

## Author:
## Description
	### Tagline
	### Detailed
	### Predeploy
	### Architecture

## DeploymentDuration
## CostEstimate

## Inputs
## Outputs
## Requirements
	### APIs
	### Roles
`


// PageData is a struct that holds the data to be njected into the mdTemplate.
// Each field corresponds to a section or piece of information in the README.md file.
type PageData struct {
	Title    string
	TitleDescription string
	Data string
}

// markdownHeading struct
type markdownHeading struct {
	HeadLevel   int
	HeadOrder   int
	HeadTitle   string
	Content     bool
	Description string
}


// createReadme generates or overwrites a README.md file in bpCurrPath.
// It uses fileTemplate if provided; otherwise, it defaults to the global mdTemplate.
// The content is populated using data from pageData.
func createReadme(bpCurrPath string, pageData PageData, fileTemplate string) error {
	currentTemplate := fileTemplate
	if len(currentTemplate) == 0 {
		currentTemplate = mdTemplate
	}

	tmpl, err := template.New("markdown").Parse(currentTemplate)
	if err != nil {
		return fmt.Errorf("Error parsing README template:"+err.Error())

	}

	filePath := filepath.Join(bpCurrPath, readmeFileName)
	file, err := os.Create(filePath) // os.Create truncates the file if it exists or creates it
	if err != nil {
		return fmt.Errorf("Error creating/updating README file '"+ filePath+ "': "+ err.Error())

	}
	defer file.Close()

	err = tmpl.Execute(file, pageData)
	if err != nil {
		return fmt.Errorf("Error executing README template for file '"+ filePath+ "': "+ err.Error())

	}
	Log.Info("README file created/updated successfully: ", filePath)
	return nil
}

// extractReadmeData reads the content of the file at the given path.
func extractReadmeData(path string) ([]byte, error) {
	readmeContent, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("blueprint README.md does not exist at %s (see https://tinyurl.com/tf-mod-readme): %w", path, err)
		}
		return nil, fmt.Errorf("error reading README.md at %s: %w", path, err)
	}
	return readmeContent, nil
}


// validateReadme checks an existing README.md file for required and recommended sections.
func validateReadme(readmeFilePath string) error {
	Log.Info("\n\n**************** Validating README.md ****************")

	readmeContent, err:= extractReadmeData(readmeFilePath)
	if err!=nil {
		return err;
	}
	// Repo details validation (comment suggests further validation, not implemented here)

	// Must have headings
	mustHaveHeadings := []markdownHeading{
		{
			HeadLevel:   1,
			HeadOrder:   1,
			HeadTitle:   "",
			Content:     false,
			Description: "Title of the module (first H1 heading)",
		},
	}

	missingMustHaveHeadings := missingHeadings(readmeContent, mustHaveHeadings)
	var errorMessages []string // Renamed from 'errors' to avoid shadowing package
	if len(missingMustHaveHeadings) > 0 {
		for _, heading := range missingMustHaveHeadings {
			title := heading.HeadTitle
			if title == "" { // Provide a more descriptive name if HeadTitle was for matching any
				title = fmt.Sprintf("H%d (order %d, expecting: %s)", heading.HeadLevel, heading.HeadOrder, heading.Description)
			}
			e := fmt.Sprintf("Missing required heading '%s': %s", title, heading.Description)
			errorMessages = append(errorMessages, e)
		}
		// This inner 'if len(errors) > 0' is redundant due to the outer 'if len(missingMustHaveHeadings) > 0'
		// but kept as it doesn't harm.
		if len(errorMessages) > 0 {
			return fmt.Errorf(strings.Join(errorMessages, "\n"))
		}
	}

	goodToHaveMarkdownHeadings := []markdownHeading{

		{
			HeadLevel:   -1, // -1 can mean any level if HeadTitle is specific
			HeadOrder:   -1, // -1 can mean any order
			HeadTitle:   "Tagline",
			Content:     true,
			Description: "corresponds to `BlueprintMetadataSpec.BlueprintInfo.description.Tagline` field in metadata.yaml",
		},
	}
	missingGoodToHaveHeadings := missingHeadings(readmeContent, goodToHaveMarkdownHeadings)

	if len(missingGoodToHaveHeadings) > 0 {
		var warningMsg []string
		for _, heading := range missingGoodToHaveHeadings {
			warningMsg = append(warningMsg, fmt.Sprintf("# %s: \t%s", heading.HeadTitle, heading.Description))
		}
		Log.Warn("\nGood to have headings missing in README.md:\n" + strings.Join(warningMsg, "\n") + "\n")
	}

	// Define otherHeadings (similar to goodToHaveMarkdownHeadings but for informational logging)
	otherHeadings := []markdownHeading{
		{HeadLevel: -1, HeadOrder: -1, HeadTitle: "Detailed", Content: true, Description: "corresponds to `BlueprintMetadataSpec.BlueprintInfo.Description.Detailed`"},
		{HeadLevel: -1, HeadOrder: -1, HeadTitle: "PreDeploy", Content: true, Description: "corresponds to `BlueprintMetadataSpec.BlueprintInfo.Description.PreDeploy`"},
		{HeadLevel: -1, HeadOrder: -1, HeadTitle: "Architecture", Content: true, Description: "corresponds to `BlueprintMetadataSpec.BlueprintInfo.Description.Architecture`"},
		{HeadLevel: -1, HeadOrder: -1, HeadTitle: "Deployment Duration", Content: true, Description: "corresponds to `BlueprintMetadataSpec.BlueprintInfo.DeploymentDuration`"},
		{HeadLevel: -1, HeadOrder: -1, HeadTitle: "Cost", Content: true, Description: "corresponds to `BlueprintMetadataSpec.BlueprintInfo.CostEstimate`"},
		{HeadLevel: -1, HeadOrder: -1, HeadTitle: "Documentation", Content: true, Description: "corresponds to `BlueprintMetadataSpec.BlueprintContent.Documentation` (diagram and description)"},
	}
	otherMissingHeadings := missingHeadings(readmeContent, otherHeadings)

	if len(otherMissingHeadings) > 0 {
		var infoMsg []string
		for _, heading := range otherMissingHeadings {
			infoMsg = append(infoMsg, fmt.Sprintf("# %s:\t%s", heading.HeadTitle, heading.Description))
		}
		Log.Info("\nOther informational headings that could be added to README.md:\n" + strings.Join(infoMsg, "\n") + "\n")
	}
	Log.Info("\n**************** Done Validating README.md ****************\n")
	return nil
}


// missingHeadings checks for the presence of specified markdown headings in the content.
func missingHeadings(content []byte, markdownHeadings []markdownHeading) []markdownHeading {
	var notFoundHeadings []markdownHeading
	for _, heading := range markdownHeadings {
		_, err := bpmetadata.GetMdContent(content, heading.HeadLevel, heading.HeadOrder, heading.HeadTitle, heading.Content)
		if err != nil {
			notFoundHeadings = append(notFoundHeadings, heading)
		}
	}
	return notFoundHeadings
}
