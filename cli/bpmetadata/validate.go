package bpmetadata

import (
	_ "embed"
	"encoding/json"
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

	var vErrs []error
	err := validateMetadataYaml(bpPath, schemaLoader)
	if err != nil {
		vErrs = append(vErrs, err)
		Log.Error("metadata validation failed", "err", err)
	}

	modulesPath := path.Join(bpPath, modulesPath)
	_, err = os.Stat(modulesPath)

	if err != nil {
		moduleDirs, _ := util.WalkTerraformDirs(modulesPath)

		for _, d := range moduleDirs {

			m := path.Join(d, metadataFileName)
			_, err := os.Stat(m)

			// log info msg and continue if the file does not exist
			if err != nil {
				Log.Info("metadata for module does not exist", "path", d)
				continue
			}

			err = validateMetadataYaml(m, schemaLoader)
			if err != nil {
				vErrs = append(vErrs, err)
				Log.Error("metadata validation failed", "err", err)
			}
		}
	}

	if len(vErrs) > 0 {
		return fmt.Errorf("metadata validation failed for at least one blueprint")
	}

	return nil
}

// validateMetadata validates an individual yaml file present at path "p"
func validateMetadataYaml(m string, schema gojsonschema.JSONLoader) error {

	// prepare metadata for validation
	mBytes, _ := ioutil.ReadFile(m)
	err := prepareMetadataBytes(&mBytes)
	if err != nil {
		Log.Error("unable to read metadata", "path", m, "error", err)
		return err
	}

	// we can use a generic interface here since the validation
	// will be done against the schema built from BlueprintMetadata
	var obj interface{}
	if err = json.Unmarshal(mBytes, &obj); err != nil {
		Log.Error("unable to unmarshal metadata", "error", err)
		return err
	}

	// load metadata from the path
	yamlLoader := gojsonschema.NewStringLoader(string(mBytes))

	// validate metadata against the schema
	result, err := gojsonschema.Validate(schema, yamlLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		//var vErrors string
		for _, e := range result.Errors() {
			Log.Error("validation error", "err", e)
			//vErrors += fmt.Sprintf("- %s\n", e)
		}

		return fmt.Errorf("metdata validation failed for file: %s", m)
	}

	Log.Info("metadata is valid", "path", m)
	return nil
}

// prepares metadata bytes for validation since direct
// validation of YAML is not possible
func prepareMetadataBytes(b *[]byte) error {
	if len(*b) == 0 {
		return fmt.Errorf("metadata contents can not be empty")
	}

	json, err := yaml.YAMLToJSON(*b)
	if err != nil {
		return fmt.Errorf("metadata contents are invalid: %s", err.Error())
	}

	// successful conversion
	*b = json

	return nil
}
