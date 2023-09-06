package bpmetadata

import (
	_ "embed"
	"fmt"
	"os"
	"path"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

//go:embed schema/gcp-blueprint-metadata.json
var s []byte

// validateMetadata validates the metadata files for the provided
// blueprint path. This validation occurs for top-level blueprint
// metadata and blueprints in the modules/ folder, if present
func validateMetadata(bpPath, wdPath string) error {
	// load schema from the binary
	schemaLoader := gojsonschema.NewStringLoader(string(s))

	// check if the provided output path is relative
	if !path.IsAbs(bpPath) {
		bpPath = path.Join(wdPath, bpPath)
	}

	// We don't need to validate metadata under .terraform folders
	skipDirsToValidate := []string{".terraform/"}
	metadataFiles, err := util.FindFilesWithPattern(bpPath, `^metadata(?:.display)?.yaml$`, skipDirsToValidate)
	if err != nil {
		Log.Error("unable to read at: %s", bpPath, "err", err)
	}

	var vErrs []error
	for _, f := range metadataFiles {
		err = validateMetadataYaml(f, schemaLoader)
		if err != nil {
			vErrs = append(vErrs, err)
			Log.Error("core metadata validation failed", "err", err)
		}
	}

	if len(vErrs) > 0 {
		return fmt.Errorf("metadata validation failed for at least one blueprint")
	}

	return nil
}

// validateMetadata validates an individual yaml file present at path "m"
func validateMetadataYaml(m string, schema gojsonschema.JSONLoader) error {
	// prepare metadata for validation by converting it from YAML to JSON
	mBytes, err := convertYamlToJson(m)
	if err != nil {
		return fmt.Errorf("yaml to json conversion failed for metadata at path %s. error: %w", m, err)
	}

	// load metadata from the path
	yamlLoader := gojsonschema.NewStringLoader(string(mBytes))

	// validate metadata against the schema
	result, err := gojsonschema.Validate(schema, yamlLoader)
	if err != nil {
		return fmt.Errorf("metadata validation failed for %s. error: %w", m, err)
	}

	if !result.Valid() {
		for _, e := range result.Errors() {
			Log.Error("validation error", "err", e)
		}

		return fmt.Errorf("metdata validation failed for: %s", m)
	}

	Log.Info("metadata is valid", "path", m)
	return nil
}

// prepares metadata bytes for validation since direct
// validation of YAML is not possible
func convertYamlToJson(m string) ([]byte, error) {
	// read metadata for validation
	b, err := os.ReadFile(m)
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata at path %s. error: %w", m, err)
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("metadata contents can not be empty")
	}

	json, err := yaml.YAMLToJSON(b)
	if err != nil {
		return nil, fmt.Errorf("metadata contents are invalid: %s", err.Error())
	}

	return json, nil
}
