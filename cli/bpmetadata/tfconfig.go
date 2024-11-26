package bpmetadata

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata/parser"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/terraform"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	testingiface "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	versionRegEx = "/v([0-9]+[.0-9]*)$"
)

type blueprintVersion struct {
	moduleVersion     string
	requiredTfVersion string
}

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "terraform",
			LabelNames: nil,
		},
		{
			Type:       "locals",
			LabelNames: nil,
		},
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "module",
			LabelNames: []string{"name"},
		},
	},
}

var metaSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "provider_meta",
			LabelNames: []string{"name"},
		},
	},
}

var variableSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
	},
}

var metaBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "module_name",
		},
	},
}

var moduleSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "activate_apis",
		},
	},
}

// Create alias for generateTFStateFile so we can mock it in unit test.
var tfState = generateTFState

// getBlueprintVersion gets both the required core version and the
// version of the blueprint
func getBlueprintVersion(configPath string) (*blueprintVersion, error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	//create hcl file object from the provided tf config
	fileName := filepath.Base(configPath)
	var diags hcl.Diagnostics
	p := hclparse.NewParser()
	versionsFile, fileDiags := p.ParseHCL(bytes, fileName)
	diags = append(diags, fileDiags...)
	err = hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	//parse out the blueprint version from the config
	modName, err := parseBlueprintVersion(versionsFile, diags)
	if err != nil {
		return nil, fmt.Errorf("error parsing blueprint version: %w", err)
	}

	//parse out the required version from the config
	var hclModule tfconfig.Module
	hclModule.RequiredProviders = make(map[string]*tfconfig.ProviderRequirement)
	hclModuleDiag := tfconfig.LoadModuleFromFile(versionsFile, &hclModule)
	diags = append(diags, hclModuleDiag...)
	err = hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	requiredCore := ""
	if len(hclModule.RequiredCore) != 0 {
		//always looking for the first element since tf blueprints
		//have one required core version
		requiredCore = hclModule.RequiredCore[0]
	}

	return &blueprintVersion{
		moduleVersion:     modName,
		requiredTfVersion: requiredCore,
	}, nil
}

// parseBlueprintVersion gets the blueprint version from the provided config
// from the provider_meta block
func parseBlueprintVersion(versionsFile *hcl.File, diags hcl.Diagnostics) (string, error) {
	re := regexp.MustCompile(versionRegEx)
	// PartialContent() returns TF content containing blocks and attributes
	// based on the provided schema
	rootContent, _, rootContentDiags := versionsFile.Body.PartialContent(rootSchema)
	diags = append(diags, rootContentDiags...)
	err := hasHclErrors(diags)
	if err != nil {
		return "", err
	}

	// based on the content returned, iterate through blocks and look for
	// the terraform block specfically
	for _, rootBlock := range rootContent.Blocks {
		if rootBlock.Type != "terraform" {
			continue
		}

		// do a PartialContent() call again but now for the provider_meta block
		// within the terraform block
		tfContent, _, tfContentDiags := rootBlock.Body.PartialContent(metaSchema)
		diags = append(diags, tfContentDiags...)
		err := hasHclErrors(diags)
		if err != nil {
			return "", err
		}

		for _, tfContentBlock := range tfContent.Blocks {
			if tfContentBlock.Type != "provider_meta" {
				continue
			}

			// this PartialContent() call with get the module_name attribute
			// that contains the version info
			metaContent, _, metaContentDiags := tfContentBlock.Body.PartialContent(metaBlockSchema)
			diags = append(diags, metaContentDiags...)
			err := hasHclErrors(diags)
			if err != nil {
				return "", err
			}

			versionAttr, defined := metaContent.Attributes["module_name"]
			if !defined {
				return "", fmt.Errorf("module_name not defined for provider_meta")
			}

			// get the module name from the version attribute and extract the
			// version name only
			var modName string
			diags := gohcl.DecodeExpression(versionAttr.Expr, nil, &modName)
			err = hasHclErrors(diags)
			if err != nil {
				return "", err
			}

			m := re.FindStringSubmatch(modName)
			if len(m) > 0 {
				return m[len(m)-1], nil
			}

			return "", nil
		}

		break
	}

	return "", nil
}

// parseBlueprintProviderVersions gets the blueprint provider_versions from the provided config
// from the required_providers block.
func parseBlueprintProviderVersions(versionsFile *hcl.File) ([]*ProviderVersion, error) {
	var v []*ProviderVersion
	// parse out the required providers from the config
	var hclModule tfconfig.Module
	hclModule.RequiredProviders = make(map[string]*tfconfig.ProviderRequirement)
	diags := tfconfig.LoadModuleFromFile(versionsFile, &hclModule)
	err := hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	for _, providerData := range hclModule.RequiredProviders {
		if providerData.Source == "" {
			Log.Info("Not found source in provider settings\n")
			continue
		}
		if len(providerData.VersionConstraints) == 0 {
			Log.Info("Not found version in provider settings\n")
			continue
		}
		v = append(v, &ProviderVersion{
			Source:  providerData.Source,
			Version: strings.Join(providerData.VersionConstraints, ", "),
		})
	}
	// Sort provider_versions
	sort.SliceStable(v, func(i, j int) bool { return v[i].Source < v[j].Source })
	return v, nil
}

// getBlueprintInterfaces gets the variables and outputs associated
// with the blueprint
func getBlueprintInterfaces(configPath string) (*BlueprintInterface, error) {
	//load the configs from the dir path
	mod, diags := tfconfig.LoadModule(configPath)
	err := hasTfconfigErrors(diags)
	if err != nil {
		return nil, err
	}

	var variables []*BlueprintVariable
	for _, val := range mod.Variables {
		v := getBlueprintVariable(val)
		variables = append(variables, v)
	}

	// Get the varible orders from tf file.
	variableOrders, sortErr := getBlueprintVariableOrders(configPath)
	if sortErr != nil {
		Log.Info("Failed to get variables orders. Fallback to sort by variable names.", sortErr)
		sort.SliceStable(variables, func(i, j int) bool { return variables[i].Name < variables[j].Name })
	} else {
		Log.Info("Sort variables by the original input order.")
		sort.SliceStable(variables, func(i, j int) bool {
			return variableOrders[variables[i].Name] < variableOrders[variables[j].Name]
		})
	}

	var outputs []*BlueprintOutput
	for _, val := range mod.Outputs {
		o := getBlueprintOutput(val)
		outputs = append(outputs, o)
	}

	// Sort outputs
	sort.SliceStable(outputs, func(i, j int) bool { return outputs[i].Name < outputs[j].Name })

	return &BlueprintInterface{
		Variables: variables,
		Outputs:   outputs,
	}, nil
}

func getBlueprintVariableOrders(configPath string) (map[string]int, error) {
	p := hclparse.NewParser()
	variableFile, hclDiags := p.ParseHCLFile(filepath.Join(configPath, "variables.tf"))
	err := hasHclErrors(hclDiags)
	if hclDiags.HasErrors() {
		return nil, err
	}
	variableContent, _, hclDiags := variableFile.Body.PartialContent(variableSchema)
	err = hasHclErrors(hclDiags)
	if hclDiags.HasErrors() {
		return nil, err
	}
	variableOrderKeys := make(map[string]int)
	for i, block := range variableContent.Blocks {
		// We only care about variable blocks.
		if block.Type != "variable" {
			continue
		}
		// We expect a single label which is the variable name.
		if len(block.Labels) != 1 || len(block.Labels[0]) == 0 {
			return nil, fmt.Errorf("Vaiable block has no name.")
		}

		variableOrderKeys[block.Labels[0]] = i
	}
	return variableOrderKeys, nil
}

// build variable
func getBlueprintVariable(modVar *tfconfig.Variable) *BlueprintVariable {
	v := &BlueprintVariable{
		Name:        modVar.Name,
		Description: modVar.Description,
		Required:    modVar.Required,
		VarType:     modVar.Type,
	}
	if modVar.Default == nil {
		return v
	}

	vl, err := structpb.NewValue(modVar.Default)
	if err == nil {
		v.DefaultValue = vl
	}

	return v
}

// build output
func getBlueprintOutput(modOut *tfconfig.Output) *BlueprintOutput {
	return &BlueprintOutput{
		Name:        modOut.Name,
		Description: modOut.Description,
	}
}

// getBlueprintRequirements gets the services and roles associated
// with the blueprint
func getBlueprintRequirements(rolesConfigPath, servicesConfigPath, versionsConfigPath string) (*BlueprintRequirements, error) {
	//parse blueprint roles
	p := hclparse.NewParser()
	rolesFile, diags := p.ParseHCLFile(rolesConfigPath)
	err := hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	r, err := parseBlueprintRoles(rolesFile)
	if err != nil {
		return nil, err
	}

	//parse blueprint services
	servicesFile, diags := p.ParseHCLFile(servicesConfigPath)
	err = hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	s, err := parseBlueprintServices(servicesFile)
	if err != nil {
		return nil, err
	}

	versionCfgFileExists, _ := fileExists(versionsConfigPath)

	if !versionCfgFileExists {
		return &BlueprintRequirements{
			Roles:    r,
			Services: s,
		}, nil
	}

	//parse blueprint provider versions
	versionsFile, diags := p.ParseHCLFile(versionsConfigPath)
	err = hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	v, err := parseBlueprintProviderVersions(versionsFile)
	if err != nil {
		return nil, err
	}

	return &BlueprintRequirements{
		Roles:            r,
		Services:         s,
		ProviderVersions: v,
	}, nil

}

// parseBlueprintRoles gets the roles required for the blueprint to be provisioned
func parseBlueprintRoles(rolesFile *hcl.File) ([]*BlueprintRoles, error) {
	var r []*BlueprintRoles
	iamContent, _, diags := rolesFile.Body.PartialContent(rootSchema)
	err := hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	for _, block := range iamContent.Blocks {
		if block.Type != "locals" {
			continue
		}

		iamAttrs, diags := block.Body.JustAttributes()
		err := hasHclErrors(diags)
		if err != nil {
			return nil, err
		}

		for k := range iamAttrs {
			var iamRoles []string
			attrValue, _ := iamAttrs[k].Expr.Value(nil)
			if !attrValue.Type().IsTupleType() {
				continue
			}

			ie := attrValue.ElementIterator()
			for ie.Next() {
				_, v := ie.Element()
				iamRoles = append(iamRoles, v.AsString())
			}

			containerRoles := &BlueprintRoles{
				// TODO: (b/248123274) no good way to associate granularity yet
				Level: "Project",
				Roles: iamRoles,
			}

			r = append(r, containerRoles)
		}

		// because we're only interested in the top-level locals block
		break
	}

	sortBlueprintRoles(r)
	return r, nil
}

// Sort blueprint roles.
func sortBlueprintRoles(r []*BlueprintRoles) {
	sort.SliceStable(r, func(i, j int) bool {
		// 1. Sort by Level
		if r[i].Level != r[j].Level {
			return r[i].Level < r[j].Level
		}

		// 2. Sort by the len of roles
		if len(r[i].Roles) != len(r[j].Roles) {
			return len(r[i].Roles) < len(r[j].Roles)
		}

		// 3. Sort by the first role (if available)
		if len(r[i].Roles) > 0 && len(r[j].Roles) > 0 {
			return r[i].Roles[0] < r[j].Roles[0]
		}

		return false
	})
}

// parseBlueprintServices gets the gcp api services required for the blueprint
// to be provisioned
func parseBlueprintServices(servicesFile *hcl.File) ([]string, error) {
	var s []string
	servicesContent, _, diags := servicesFile.Body.PartialContent(rootSchema)
	err := hasHclErrors(diags)
	if err != nil {
		return nil, err
	}

	for _, block := range servicesContent.Blocks {
		if block.Type != "module" {
			continue
		}

		moduleContent, _, moduleContentDiags := block.Body.PartialContent(moduleSchema)
		diags = append(diags, moduleContentDiags...)
		err := hasHclErrors(diags)
		if err != nil {
			return nil, err
		}

		apisAttr, defined := moduleContent.Attributes["activate_apis"]
		if !defined {
			return nil, fmt.Errorf("activate_apis not defined for project module")
		}

		diags = gohcl.DecodeExpression(apisAttr.Expr, nil, &s)
		err = hasHclErrors(diags)
		if err != nil {
			return nil, err
		}

		// because we're only interested in the top-level modules block
		break
	}

	return s, nil
}

func hasHclErrors(diags hcl.Diagnostics) error {
	for _, diag := range diags {
		if diag.Severity == hcl.DiagError {
			return fmt.Errorf("hcl error: %s | detail: %s", diag.Summary, diag.Detail)
		}
	}

	return nil
}

// this is almost a dup of hasHclErrors because the TF api has two
// different structs for diagnostics...
func hasTfconfigErrors(diags tfconfig.Diagnostics) error {
	for _, diag := range diags {
		if diag.Severity == tfconfig.DiagError {
			return fmt.Errorf("hcl error: %s | detail: %s", diag.Summary, diag.Detail)
		}
	}

	return nil
}

// MergeExistingConnections merges existing connections from an old BlueprintInterface into a new one,
// preserving manually authored connections.
func mergeExistingConnections(newInterfaces, existingInterfaces *BlueprintInterface) {
	if existingInterfaces == nil {
		return // Nothing to merge if existingInterfaces is nil
	}

	for i, variable := range newInterfaces.Variables {
		for _, existingVariable := range existingInterfaces.Variables {
			if variable.Name == existingVariable.Name && existingVariable.Connections != nil {
				newInterfaces.Variables[i].Connections = existingVariable.Connections
			}
		}
	}
}

// mergeExistingOutputTypes merges existing output types from an old BlueprintInterface into a new one,
// preserving manually authored types.
func mergeExistingOutputTypes(newInterfaces, existingInterfaces *BlueprintInterface) {
	if existingInterfaces == nil {
		return // Nothing to merge if existingInterfaces is nil
	}

	existingOutputs := make(map[string]*BlueprintOutput)
	for _, output := range existingInterfaces.Outputs {
		existingOutputs[output.Name] = output
	}

	for i, output := range newInterfaces.Outputs {
		if output.Type != nil {
			continue
		}
		if existingOutput, ok := existingOutputs[output.Name]; ok && existingOutput.Type != nil {
			newInterfaces.Outputs[i].Type = existingOutput.Type
		}
	}
}

// UpdateOutputTypes generates the terraform.tfstate file, extracts output types from it,
// and updates the output types in the provided BlueprintInterface.
func updateOutputTypes(bpPath string, bpInterfaces *BlueprintInterface) error {
	// Generate the terraform.tfstate file
	stateData, err := tfState(bpPath)
	if err != nil {
		return fmt.Errorf("error generating terraform.tfstate file: %w", err)
	}

	// Parse the state file and extract output types
	outputTypes, err := parser.ParseOutputTypesFromState(stateData)
	if err != nil {
		return fmt.Errorf("error parsing output types: %w", err)
	}

	// Update the output types in the BlueprintInterface
	for i, output := range bpInterfaces.Outputs {
		if outputType, ok := outputTypes[output.Name]; ok {
			bpInterfaces.Outputs[i].Type = outputType
		}
	}
	return nil
}

// generateTFState generates the terraform.tfstate by running terraform init and apply, and terraform show to capture the state.
func generateTFState(bpPath string) ([]byte, error) {
	var stateData []byte
	// Construct the path to the test/setup directory
	tfDir := filepath.Join(bpPath)

	// testing.T checks verbose flag to determine its mode. Add this line as a flags initializer
	// so the program doesn't panic
	flag.Parse()
	runtimeT := testingiface.RuntimeT{}

	root := tft.NewTFBlueprintTest(
		&runtimeT,
		tft.WithTFDir(tfDir), // Setup test at the blueprint path,
	)

	root.DefineVerify(func(assert *assert.Assertions) {
		stateStr, err := terraform.ShowE(&runtimeT, root.GetTFOptions())
		if err != nil {
			assert.FailNowf("Failed to generate terraform.tfstate", "Error calling `terraform show`: %v", err)
		}

		stateData = []byte(stateStr)
	})

	root.Test() // This will run terraform init and apply, and then destroy

	return stateData, nil
}
