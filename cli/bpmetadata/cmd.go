package bpmetadata

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

var mdFlags struct {
	path   string
	nested bool
	force  bool
}

const (
	readmeFileName     = "README.md"
	tfVersionsFileName = "versions.tf"
	tfRolesFileName    = "test/setup/iam.tf"
	tfServicesFileName = "test/setup/main.tf"
	iconFilePath       = "assets/icon.png"
	modulesPath        = "modules"
	examplesPath       = "examples"
	metadataFileName   = "metadata.yaml"
	metadataApiVersion = "blueprints.cloud.google.com/v1alpha1"
	metadataKind       = "BlueprintMetadata"
)

func init() {
	viper.AutomaticEnv()

	Cmd.Flags().StringVarP(&mdFlags.path, "path", "p", ".", "Path to the blueprint for generating metadata.")
	Cmd.Flags().BoolVar(&mdFlags.nested, "nested", true, "Flag for generating metadata for nested blueprint, if any.")
	Cmd.Flags().BoolVarP(&mdFlags.force, "force", "f", false, "Force the generation of metadata")
}

var Cmd = &cobra.Command{
	Use:   "metadata",
	Short: "Generates blueprint metatda",
	Long:  `Generates metadata.yaml for specified blueprint`,
	Args:  cobra.NoArgs,
	RunE:  generate,
}

// The top-level command function that generates metadata based on the provided flags
func generate(cmd *cobra.Command, args []string) error {
	wdPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working dir: %w", err)
	}

	bpPath := path.Join(wdPath, mdFlags.path)

	//try to read existing metadata.yaml
	bpObj, err := UnmarshalMetadata(bpPath)
	if err != nil && !mdFlags.force {
		return err
	}

	// create metadata details
	bpMetaObj, err := CreateBlueprintMetadata(bpPath, bpObj)
	if err != nil {
		return fmt.Errorf("error creating metadata for blueprint: %w", err)
	}

	// write metadata to disk
	err = WriteMetadata(bpMetaObj)
	if err != nil {
		return fmt.Errorf("error writing metadata to disk: %w", err)
	}

	return nil
}

func CreateBlueprintMetadata(bpPath string, bpMetadataObj *BlueprintMetadata) (*BlueprintMetadata, error) {
	// verfiy that the blueprint path is valid & get repo details
	repoDetails, err := getRepoDetailsByPath(bpPath)
	if err != nil {
		return nil, err
	}

	// start creating blueprint metadata
	bpMetadataObj.Meta = yaml.ResourceMeta{
		TypeMeta: yaml.TypeMeta{
			APIVersion: metadataApiVersion,
			Kind:       metadataKind,
		},
		ObjectMeta: yaml.ObjectMeta{
			NameMeta: yaml.NameMeta{
				Name:      repoDetails.Name,
				Namespace: "",
			},
			Labels:      bpMetadataObj.Meta.ObjectMeta.Labels,
			Annotations: map[string]string{"config.kubernetes.io/local-config": "true"},
		},
	}

	// start creating the Spec node
	readmeContent, err := os.ReadFile(path.Join(bpPath, readmeFileName))
	if err != nil {
		return nil, fmt.Errorf("error reading blueprint readme markdown: %w", err)
	}

	info, err := createInfo(bpPath, readmeContent)
	if err != nil {
		return nil, fmt.Errorf("error creating blueprint info: %w", err)
	}

	interfaces, err := createInterfaces(bpPath, &bpMetadataObj.Spec.Interfaces)
	if err != nil {
		return nil, fmt.Errorf("error creating blueprint interfaces: %w", err)
	}

	rolesCfgPath := path.Join(repoDetails.Source.RootPath, tfRolesFileName)
	svcsCfgPath := path.Join(repoDetails.Source.RootPath, tfServicesFileName)
	requirements, err := getBlueprintRequirements(rolesCfgPath, svcsCfgPath)
	if err != nil {
		return nil, fmt.Errorf("error creating blueprint requirements: %w", err)
	}

	content := createContent(bpPath, repoDetails.Source.RootPath, readmeContent, &bpMetadataObj.Spec.Content)

	bpMetadataObj.Spec = BlueprintMetadataSpec{
		Info:         *info,
		Content:      *content,
		Interfaces:   *interfaces,
		Requirements: *requirements,
	}

	return bpMetadataObj, nil
}

func createInfo(bpPath string, readmeContent []byte) (*BlueprintInfo, error) {
	i := &BlueprintInfo{}
	title, err := getMdContent(readmeContent, 1, 1, "", false)
	if err != nil {
		return nil, err
	}

	i.Title = title.literal

	repoDetails, err := getRepoDetailsByPath(bpPath)
	if err != nil {
		return nil, err
	}

	i.Source = &BlueprintRepoDetail{
		Repo:       repoDetails.Source.Path,
		SourceType: "git",
	}

	versionInfo, err := getBlueprintVersion(path.Join(bpPath, tfVersionsFileName))
	if err != nil {
		return nil, err
	}

	i.Version = versionInfo.moduleVersion

	// actuation tool
	i.ActuationTool = BlueprintActuationTool{
		Version: versionInfo.requiredTfVersion,
		Flavor:  "Terraform",
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

	// create icon
	iPath := path.Join(repoDetails.Source.RootPath, iconFilePath)
	exists, _ := fileExists(iPath)
	if exists {
		i.Icon = iconFilePath
	}

	return i, nil
}

func createInterfaces(bpPath string, interfaces *BlueprintInterface) (*BlueprintInterface, error) {
	i, err := getBlueprintInterfaces(bpPath)
	if err != nil {
		return nil, err
	}

	if interfaces.VariableGroups != nil {
		i.VariableGroups = interfaces.VariableGroups
	}

	return i, nil
}

func createContent(bpPath string, rootPath string, readmeContent []byte, content *BlueprintContent) *BlueprintContent {
	//var content BlueprintContent
	var docListToSet []BlueprintListContent
	documentation, err := getMdContent(readmeContent, -1, -1, "Documentation", true)
	if err == nil {
		for _, li := range documentation.listItems {
			doc := BlueprintListContent{
				Title: li.text,
				Url:   li.url,
			}

			docListToSet = append(docListToSet, doc)
		}

		content.Documentation = docListToSet
	}

	// create sub-blueprints
	modPath := path.Join(bpPath, modulesPath)
	modContent, err := getModules(modPath)
	if err == nil {
		content.SubBlueprints = modContent
	}

	// create examples
	exPath := path.Join(rootPath, examplesPath)
	exContent, err := getExamples(exPath)
	if err == nil {
		content.Examples = exContent
	}

	return content
}

func WriteMetadata(obj *BlueprintMetadata) error {
	// marshal and write the file
	yFile, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	return os.WriteFile(metadataFileName, yFile, 0644)
}

func UnmarshalMetadata(bpPath string) (*BlueprintMetadata, error) {
	bpObj := BlueprintMetadata{}
	metaFilePath := path.Join(bpPath, metadataFileName)

	// return empty metadata if file does not exist or if the file is not read
	if _, err := os.Stat(metaFilePath); errors.Is(err, os.ErrNotExist) {
		return &bpObj, nil
	}

	f, err := os.ReadFile(metaFilePath)
	if err != nil {
		return &bpObj, fmt.Errorf("unable to read metadata from the existing file: %w", err)
	}

	err = yaml.Unmarshal(f, &bpObj)
	if err != nil {
		return &bpObj, err
	}

	currVersion := bpObj.Meta.TypeMeta.APIVersion
	currKind := bpObj.Meta.TypeMeta.Kind

	//validate GVK for current metadata
	if currVersion != metadataApiVersion {
		return &bpObj, fmt.Errorf("found incorrect version for the metadata: %s. Supported version is: %s", currVersion, metadataApiVersion)
	}

	if currKind != metadataKind {
		return &bpObj, fmt.Errorf("found incorrect kind for the metadata: %s. Supported kind is %s", currKind, metadataKind)
	}

	return &bpObj, nil
}
