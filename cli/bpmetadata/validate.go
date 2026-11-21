package bpmetadata

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

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

	if path.Base(m) == metadataDisplayFileName {
		if err := validateDisplayMetadataCompleteness(m); err != nil {
			Log.Error("display metadata completeness validation failed", "err", err)
			return err
		}
	}

	Log.Info("metadata is valid", "path", m)
	return nil
}

// validateDisplayMetadataCompleteness verifies that all variables defined in metadata.yaml
// have a corresponding entry in metadata.display.yaml when metadata.yaml exists in the same directory.
func validateDisplayMetadataCompleteness(dispPath string) error {
	dir := path.Dir(dispPath)
	coreMetaPath := path.Join(dir, metadataFileName)
	if _, err := os.Stat(coreMetaPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	coreMeta, err := UnmarshalMetadata(dir, metadataFileName)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s for display metadata validation: %w", coreMetaPath, err)
	}

	dispMeta, err := UnmarshalMetadata(dir, metadataDisplayFileName)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", dispPath, err)
	}

	if coreMeta.Spec == nil || coreMeta.Spec.Interfaces == nil || len(coreMeta.Spec.Interfaces.Variables) == 0 {
		return nil
	}

	dispVars := make(map[string]bool)
	if dispMeta.Spec != nil && dispMeta.Spec.Ui != nil && dispMeta.Spec.Ui.Input != nil && dispMeta.Spec.Ui.Input.Variables != nil {
		for vName := range dispMeta.Spec.Ui.Input.Variables {
			dispVars[vName] = true
		}
	}

	var missing []string
	for _, v := range coreMeta.Spec.Interfaces.Variables {
		if !dispVars[v.Name] {
			missing = append(missing, v.Name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("metadata.display.yaml at %s is missing variable(s) defined in metadata.yaml: %s", dispPath, strings.Join(missing, ", "))
	}

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
