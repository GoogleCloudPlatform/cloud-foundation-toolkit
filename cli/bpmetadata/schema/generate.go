package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"github.com/invopop/jsonschema"
)

const schemaFileName = "gcp-blueprint-metadata.json"

// generateSchema creates a JSON Schema based on the types
// defined in the type BlueprintMetadata and it's recursive
// children. The generated schema will be used to validate
// all metadata files for consistency and will be uploaded
// to https://www.schemastore.org/ to provide IntelliSense
// VSCode for authors manually authoring the metadata.
func generateSchemaFile(o, wdPath string) error {
	sData, err := GenerateSchema()
	if err != nil {
		return err
	}
	sData = append(sData, []byte("\n")...)

	// check if the provided output path is relative
	if !path.IsAbs(o) {
		o = path.Join(wdPath, o)
	}

	err = os.WriteFile(path.Join(o, schemaFileName), sData, 0644)
	if err != nil {
		return err
	}

	Log.Info("generated JSON schema for BlueprintMetadata", "path", path.Join(o, schemaFileName))
	return nil
}

func GenerateSchema() ([]byte, error) {
	r := &jsonschema.Reflector{}
	s := r.Reflect(&bpmetadata.BlueprintMetadata{})
	s.Version = "http://json-schema.org/draft-07/schema#"

	// defaultValue was defined as interface{} and has changed to
	// Value type with proto definitions. To keep backwards
	// compatibility for schema validation, this is being set to
	// true i.e. it's presence is validated regardless of type.
	vDef, defExists := s.Definitions["BlueprintVariable"]
	if defExists {
		vDef.Properties.Set("defaultValue", jsonschema.TrueSchema)
	}
	// JSON schema seems to infer google.protobuf.Value as object type
	// so we use the same workaround as above.
	oDef, defExists := s.Definitions["BlueprintOutput"]
	if defExists {
		oDef.Properties.Set("type", jsonschema.TrueSchema)
	}
	altDefaultDef, defExists := s.Definitions["DisplayVariable_AlternateDefault"]
	if defExists {
		altDefaultDef.Properties.Set("value", jsonschema.TrueSchema)
	}

	sData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}

	return sData, nil
}
