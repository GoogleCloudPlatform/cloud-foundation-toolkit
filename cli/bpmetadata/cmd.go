package bpmetadata

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/itchyny/json2yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"
)

var mdFlags struct {
	path          string
	nested        bool
	force         bool
	display       bool
	validate      bool
	quiet         bool
	genOutputType bool
}

const (
	readmeFileName          = "README.md"
	tfVersionsFileName      = "versions.tf"
	tfRolesFileName         = "test/setup/iam.tf"
	tfServicesFileName      = "test/setup/main.tf"
	iconFilePath            = "assets/icon.png"
	modulesPath             = "modules/"
	examplesPath            = "examples"
	metadataFileName        = "metadata.yaml"
	metadataDisplayFileName = "metadata.display.yaml"
	metadataApiVersion      = "blueprints.cloud.google.com/v1alpha1"
	metadataKind            = "BlueprintMetadata"
	localConfigAnnotation   = "config.kubernetes.io/local-config"
)

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().BoolVarP(&mdFlags.display, "display", "d", false, "Generate the display metadata used for UI rendering.")
	Cmd.Flags().BoolVarP(&mdFlags.force, "force", "f", false, "Force the generation of fresh metadata.")
	Cmd.Flags().StringVarP(&mdFlags.path, "path", "p", ".", "Path to the blueprint for generating metadata.")
	Cmd.Flags().BoolVar(&mdFlags.nested, "nested", true, "Flag for generating metadata for nested blueprint, if any.")
	Cmd.Flags().BoolVarP(&mdFlags.validate, "validate", "v", false, "Validate metadata against the schema definition.")
	Cmd.Flags().BoolVarP(&mdFlags.quiet, "quiet", "q", false, "Run in quiet mode suppressing all prompts.")
	Cmd.Flags().BoolVarP(&mdFlags.genOutputType, "generate-output-type", "g", false, "Automatically generate type field for outputs.")
}

var Cmd = &cobra.Command{
	Use:   "metadata",
	Short: "Generates blueprint metadata",
	Long:  `Generates metadata.yaml for specified blueprint`,
	Args:  cobra.NoArgs,
	RunE:  generate,
}

var repoDetails repoDetail

// The top-level command function that generates metadata based on the provided flags
func generate(cmd *cobra.Command, args []string) error {
	wdPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %w", err)
	}

	// validate metadata if there is an argument passed into the command
	if mdFlags.validate {
		if err := validateMetadata(mdFlags.path, wdPath); err != nil {
			return err
		}

		return nil
	}

	currBpPath := mdFlags.path
	if !path.IsAbs(mdFlags.path) {
		currBpPath = path.Join(wdPath, mdFlags.path)
	}

	var allBpPaths []string
	_, err = os.Stat(path.Join(currBpPath, readmeFileName))

	// throw an error and exit if root level readme.md doesn't exist
	if err != nil {
		return fmt.Errorf("top-level module does not have a readme: %w", err)
	}

	allBpPaths = append(allBpPaths, currBpPath)
	var errors []string

	// if nested, check if modules/ exists and create paths
	// for submodules
	if mdFlags.nested {
		modulesPathforBp := path.Join(currBpPath, modulesPath)
		_, err = os.Stat(modulesPathforBp)
		if os.IsNotExist(err) {
			Log.Info("sub-modules do not exist for this blueprint")
		} else {
			moduleDirs, err := util.WalkTerraformDirs(modulesPathforBp)
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				allBpPaths = append(allBpPaths, moduleDirs...)
			}
		}
	}

	for _, modPath := range allBpPaths {
		// check if module path has readme.md
		_, err := os.Stat(path.Join(modPath, readmeFileName))

		// log info if a sub-module doesn't have a readme.md and continue
		if err != nil {
			Log.Info("skipping metadata for sub-module identified as an internal module", "Path:", modPath)
			continue
		}

		err = generateMetadataForBpPath(modPath)
		if err != nil {
			e := fmt.Sprintf("path: %s\n %s", modPath, err.Error())
			errors = append(errors, e)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	Log.Info("metadata generated successfully")
	return nil
}

func generateMetadataForBpPath(bpPath string) error {
	//try to read existing metadata.yaml
	bpObj, err := UnmarshalMetadata(bpPath, metadataFileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) && !mdFlags.force {
		return err
	}

	// create core metadata
	bpMetaObj, err := CreateBlueprintMetadata(bpPath, bpObj)
	if err != nil {
		return fmt.Errorf("error creating metadata for blueprint at path: %s. Details: %w", bpPath, err)
	}

	// If the flag is set, update output types
	if mdFlags.genOutputType {
		err = updateOutputTypes(bpPath, bpMetaObj.Spec.Interfaces)
		if err != nil {
			return fmt.Errorf("error updating output types: %w", err)
		}
	}

	// write core metadata to disk
	err = WriteMetadata(bpMetaObj, bpPath, metadataFileName)
	if err != nil {
		return fmt.Errorf("error writing metadata to disk for blueprint at path: %s. Details: %w", bpPath, err)
	}

	// continue with creating display metadata if the flag is set,
	// else let the command exit
	if !mdFlags.display {
		return nil
	}

	bpDpObj, err := UnmarshalMetadata(bpPath, metadataDisplayFileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) && !mdFlags.force {
		return err
	}

	// create display metadata
	bpMetaDpObj, err := CreateBlueprintDisplayMetadata(bpPath, bpDpObj, bpMetaObj)
	if err != nil {
		return fmt.Errorf("error creating display metadata for blueprint at path: %s. Details: %w", bpPath, err)
	}

	// write display metadata to disk
	err = WriteMetadata(bpMetaDpObj, bpPath, metadataDisplayFileName)
	if err != nil {
		return fmt.Errorf("error writing display metadata to disk for blueprint at path: %s. Details: %w", bpPath, err)
	}

	return nil
}

func CreateBlueprintMetadata(bpPath string, bpMetadataObj *BlueprintMetadata) (*BlueprintMetadata, error) {
	// Verify that readme is present.
	readmeContent, err := os.ReadFile(path.Join(bpPath, readmeFileName))
	if err != nil {
		return nil, fmt.Errorf("blueprint readme markdown is missing, create one using https://tinyurl.com/tf-mod-readme | error: %w", err)
	}

	// verify that the blueprint path is valid & get repo details
	getRepoDetailsByPath(bpPath, &repoDetails, readmeContent)
	if repoDetails.ModuleName == "" && !mdFlags.quiet {
		fmt.Printf("Provide a name for the blueprint at path [%s]: ", bpPath)
		_, err := fmt.Scan(&repoDetails.ModuleName)
		if err != nil {
			fmt.Println("Unable to scan the name for the blueprint.")
		}
	}

	if repoDetails.Source.URL == "" && !mdFlags.quiet {
		fmt.Printf("Provide a URL for the blueprint source at path [%s]: ", bpPath)
		_, err := fmt.Scan(&repoDetails.Source.URL)
		if err != nil {
			fmt.Println("Unable to scan the URL for the blueprint.")
		}
	}

	// start creating blueprint metadata
	bpMetadataObj.ApiVersion = metadataApiVersion
	bpMetadataObj.Kind = metadataKind

	if bpMetadataObj.Metadata == nil {
		bpMetadataObj.Metadata = &ResourceTypeMeta{
			Name:        repoDetails.ModuleName,
			Annotations: map[string]string{localConfigAnnotation: "true"},
		}
	}

	if bpMetadataObj.Spec == nil {
		bpMetadataObj.Spec = &BlueprintMetadataSpec{}
	}

	if bpMetadataObj.Spec.Info == nil {
		bpMetadataObj.Spec.Info = &BlueprintInfo{}
	}

	// create blueprint info
	err = bpMetadataObj.Spec.Info.create(bpPath, repoDetails, readmeContent)
	if err != nil {
		return nil, fmt.Errorf("error creating blueprint info: %w", err)
	}

	var existingInterfaces *BlueprintInterface
	if bpMetadataObj.Spec.Interfaces == nil {
		bpMetadataObj.Spec.Interfaces = &BlueprintInterface{}
	} else {
		existingInterfaces = proto.Clone(bpMetadataObj.Spec.Interfaces).(*BlueprintInterface)
	}

	// create blueprint interfaces i.e. variables & outputs
	err = bpMetadataObj.Spec.Interfaces.create(bpPath)
	if err != nil {
		return nil, fmt.Errorf("error creating blueprint interfaces: %w", err)
	}

	// Merge existing connections (if any) into the newly generated interfaces
	mergeExistingConnections(bpMetadataObj.Spec.Interfaces, existingInterfaces)

	// Merge existing output types (if any) into the newly generated interfaces
	mergeExistingOutputTypes(bpMetadataObj.Spec.Interfaces, existingInterfaces)

	// get blueprint requirements
	rolesCfgPath := path.Join(repoDetails.Source.BlueprintRootPath, tfRolesFileName)
	svcsCfgPath := path.Join(repoDetails.Source.BlueprintRootPath, tfServicesFileName)
	versionsCfgPath := path.Join(bpPath, tfVersionsFileName)
	requirements, err := getBlueprintRequirements(rolesCfgPath, svcsCfgPath, versionsCfgPath)
	if err != nil {
		Log.Info("skipping blueprint requirements since roles and/or services configurations were not found as per https://tinyurl.com/tf-iam and https://tinyurl.com/tf-services")
	} else {
		bpMetadataObj.Spec.Requirements = requirements
	}

	if bpMetadataObj.Spec.Content == nil {
		bpMetadataObj.Spec.Content = &BlueprintContent{}
	}

	// create blueprint content i.e. documentation, icons, etc.
	bpMetadataObj.Spec.Content.create(bpPath, repoDetails.Source.BlueprintRootPath, readmeContent)
	return bpMetadataObj, nil
}

func CreateBlueprintDisplayMetadata(bpPath string, bpDisp, bpCore *BlueprintMetadata) (*BlueprintMetadata, error) {
	// start creating blueprint metadata
	bpDisp.ApiVersion = bpCore.ApiVersion
	bpDisp.Kind = bpCore.Kind

	if bpDisp.Metadata == nil {
		bpDisp.Metadata = &ResourceTypeMeta{
			Name:        bpCore.Metadata.Name + "-display",
			Annotations: map[string]string{localConfigAnnotation: "true"},
		}
	}

	if bpDisp.Spec == nil {
		bpDisp.Spec = &BlueprintMetadataSpec{}
	}

	if bpDisp.Spec.Info == nil {
		bpDisp.Spec.Info = &BlueprintInfo{}
	}

	if bpDisp.Spec.Ui == nil {
		bpDisp.Spec.Ui = &BlueprintUI{}
		bpDisp.Spec.Ui.Input = &BlueprintUIInput{}
	}

	bpDisp.Spec.Info.Title = bpCore.Spec.Info.Title
	bpDisp.Spec.Info.Source = bpCore.Spec.Info.Source
	buildUIInputFromVariables(bpCore.Spec.Interfaces.Variables, bpDisp.Spec.Ui.Input)

	existingInput := func() *BlueprintUIInput {
		if bpCore.Spec.Ui != nil && bpCore.Spec.Ui.Input != nil {
			return proto.Clone(bpCore.Spec.Ui.Input).(*BlueprintUIInput)
		}
		return &BlueprintUIInput{}
	}()
	// Merge existing data (if any) into the newly generated UI Input
	mergeExistingAltDefaults(bpDisp.Spec.Ui.Input, existingInput)

	return bpDisp, nil
}

func (i *BlueprintInfo) create(bpPath string, r repoDetail, readmeContent []byte) error {
	title, err := getMdContent(readmeContent, 1, 1, "", false)
	if err != nil {
		return fmt.Errorf("title tag missing in markdown, err: %w", err)
	}

	i.Title = title.literal
	rootPath := r.Source.RepoRootPath
	if rootPath == "" {
		rootPath = r.Source.BlueprintRootPath
	}

	bpDir := strings.ReplaceAll(bpPath, rootPath, "")
	i.Source = &BlueprintRepoDetail{
		Repo:       r.Source.URL,
		SourceType: r.Source.SourceType,
		Dir:        bpDir,
	}

	versionInfo, err := getBlueprintVersion(path.Join(bpPath, tfVersionsFileName))
	if err == nil {
		i.Version = versionInfo.moduleVersion
		i.ActuationTool = &BlueprintActuationTool{
			Version: versionInfo.requiredTfVersion,
			Flavor:  "Terraform",
		}
	}

	// create descriptions
	i.Description = &BlueprintDescription{}
	tagline, err := getMdContent(readmeContent, -1, -1, "Tagline", true)
	if err == nil {
		i.Description.Tagline = tagline.literal
	}

	detailed, err := getMdContent(readmeContent, -1, -1, "Detailed", true)
	if err == nil {
		i.Description.Detailed = detailed.literal
	}

	preDeploy, err := getMdContent(readmeContent, -1, -1, "PreDeploy", true)
	if err == nil {
		i.Description.PreDeploy = preDeploy.literal
	}

	var archListToSet []string
	architecture, err := getMdContent(readmeContent, -1, -1, "Architecture", true)
	if err == nil {
		for _, li := range architecture.listItems {
			archListToSet = append(archListToSet, li.text)
		}

		i.Description.Architecture = archListToSet
	}

	// create icon
	iPath := path.Join(r.Source.BlueprintRootPath, iconFilePath)
	exists, _ := fileExists(iPath)
	if exists {
		i.Icon = iconFilePath
	}

	d, err := getDeploymentDuration(readmeContent, "Deployment Duration")
	if err == nil {
		i.DeploymentDuration = d
	}

	c, err := getCostEstimate(readmeContent, "Cost")
	if err == nil {
		i.CostEstimate = c
	}

	return nil
}

func (i *BlueprintInterface) create(bpPath string) error {
	interfaces, err := getBlueprintInterfaces(bpPath)
	if err != nil {
		return err
	}

	i.Variables = interfaces.Variables
	i.Outputs = interfaces.Outputs

	return nil
}

func (c *BlueprintContent) create(bpPath string, rootPath string, readmeContent []byte) {
	var docListToSet []*BlueprintListContent
	documentation, err := getMdContent(readmeContent, -1, -1, "Documentation", true)
	if err == nil {
		for _, li := range documentation.listItems {
			doc := &BlueprintListContent{
				Title: li.text,
				Url:   li.url,
			}

			docListToSet = append(docListToSet, doc)
		}

		c.Documentation = docListToSet
	}

	// create architecture
	a, err := getArchitctureInfo(readmeContent, "Architecture")
	if err == nil {
		c.Architecture = a
	}

	// create sub-blueprints
	modPath := path.Join(bpPath, modulesPath)
	modContent, err := getModules(modPath)
	if err == nil {
		c.SubBlueprints = modContent
	}

	// create examples
	exPath := path.Join(rootPath, examplesPath)
	exContent, err := getExamples(exPath)
	if err == nil {
		c.Examples = exContent
	}
}

func WriteMetadata(obj *BlueprintMetadata, bpPath, fileName string) error {
	jBytes, err := protojson.Marshal(obj)
	if err != nil {
		return err
	}

	input := strings.NewReader(string(jBytes))
	var output strings.Builder
	if err := json2yaml.Convert(&output, input); err != nil {
		return err
	}

	return os.WriteFile(path.Join(bpPath, fileName), []byte(output.String()), 0644)
}

func UnmarshalMetadata(bpPath, fileName string) (*BlueprintMetadata, error) {
	bpObj := BlueprintMetadata{}
	metaFilePath := path.Join(bpPath, fileName)

	// return empty metadata if file does not exist or if the file is not read
	if _, err := os.Stat(metaFilePath); errors.Is(err, os.ErrNotExist) {
		return &bpObj, err
	}

	f, err := os.ReadFile(metaFilePath)
	if err != nil {
		return &bpObj, fmt.Errorf("unable to read metadata from the existing file: %w", err)
	}

	// convert yaml bytes to json bytes for unmarshaling metadata
	// content to proto definition
	j, err := yaml.YAMLToJSON(f)
	if err != nil {
		return nil, err
	}

	if err := protojson.Unmarshal(j, &bpObj); err != nil {
		return &bpObj, err
	}

	currVersion := bpObj.ApiVersion
	currKind := bpObj.Kind

	//validate GVK for current metadata
	if currVersion != metadataApiVersion {
		return &bpObj, fmt.Errorf("found incorrect version for the metadata: %s. Supported version is: %s", currVersion, metadataApiVersion)
	}

	if currKind != metadataKind {
		return &bpObj, fmt.Errorf("found incorrect kind for the metadata: %s. Supported kind is %s", currKind, metadataKind)
	}

	return &bpObj, nil
}
