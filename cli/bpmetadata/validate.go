package bpmetadata

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

//go:embed schema/bpmetadataschema.json
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

	moduleDirs := []string{bpPath}
	modulesPath := path.Join(bpPath, modulesPath)
	_, err := os.Stat(modulesPath)
	if err == nil {
		subModuleDirs, err := util.WalkTerraformDirs(modulesPath)
		if err != nil {
			Log.Warn("unable to read the submodules i.e. modules/ folder", "err", err)
		}

		moduleDirs = append(moduleDirs, subModuleDirs...)
	}

	var vErrs []error
	for _, d := range moduleDirs {
		// validate core metadata
		core := path.Join(d, metadataFileName)
		_, err := os.Stat(core)

		// log info msg and continue if the file does not exist
		if err != nil {
			Log.Info("core metadata for module does not exist", "path", core)
			continue
		}

		err = validateMetadataYaml(core, schemaLoader)
		if err != nil {
			vErrs = append(vErrs, err)
			Log.Error("core metadata validation failed", "err", err)
		}

		// validate display metadata
		disp := path.Join(d, metadataDisplayFileName)
		_, err = os.Stat(disp)

		// log info msg and continue if the file does not exist
		if err != nil {
			Log.Info("display metadata for module does not exist", "path", disp)
			continue
		}

		err = validateMetadataYaml(disp, schemaLoader)
		if err != nil {
			vErrs = append(vErrs, err)
			Log.Error("display metadata validation failed", "err", err)
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
		return fmt.Errorf("yaml to json conversion failed for metadata at path %s. error: %s", m, err)
	}

	// load metadata from the path
	yamlLoader := gojsonschema.NewStringLoader(string(mBytes))

	// validate metadata against the schema
	result, err := gojsonschema.Validate(schema, yamlLoader)
	if err != nil {
		return fmt.Errorf("metadata validation failed for %s. error: %s", m, err)
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
	b, err := ioutil.ReadFile(m)
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata at path %s. error: %s", m, err)
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
