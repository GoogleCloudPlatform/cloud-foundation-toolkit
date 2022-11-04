package bpmetadata

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
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
		return nil, fmt.Errorf("error parsing blueprint version: %v", err)
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
			gohcl.DecodeExpression(versionAttr.Expr, nil, &modName)
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

// getBlueprintInterfaces gets the variables and outputs associated
// with the blueprint
func getBlueprintInterfaces(configPath string) (*BlueprintInterface, error) {
	//load the configs from the dir path
	mod, diags := tfconfig.LoadModule(configPath)
	err := hasTfconfigErrors(diags)
	if err != nil {
		return nil, err
	}

	var variables []BlueprintVariable
	for _, val := range mod.Variables {
		v := getBlueprintVariable(val)
		variables = append(variables, v)
	}

	// Sort variables
	sort.SliceStable(variables, func(i, j int) bool { return variables[i].Name < variables[j].Name })

	var outputs []BlueprintOutput
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

//build variable
func getBlueprintVariable(modVar *tfconfig.Variable) BlueprintVariable {
	return BlueprintVariable{
		Name:        modVar.Name,
		Description: modVar.Description,
		Default:     modVar.Default,
		Required:    modVar.Required,
		VarType:     modVar.Type,
	}
}

//build output
func getBlueprintOutput(modOut *tfconfig.Output) BlueprintOutput {
	return BlueprintOutput{
		Name:        modOut.Name,
		Description: modOut.Description,
	}
}

// getBlueprintRequirements gets the services and roles associated
// with the blueprint
func getBlueprintRequirements(rolesConfigPath, servicesConfigPath string) (*BlueprintRequirements, error) {
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

	return &BlueprintRequirements{
		Roles:    r,
		Services: s,
	}, nil
}

// parseBlueprintRoles gets the roles required for the blueprint to be provisioned
func parseBlueprintRoles(rolesFile *hcl.File) ([]BlueprintRoles, error) {
	var r []BlueprintRoles
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

		for k, _ := range iamAttrs {
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

			containerRoles := BlueprintRoles{
				// TODO: (b/248123274) no good way to associate granularity yet
				Level: "Project",
				Roles: iamRoles,
			}

			r = append(r, containerRoles)
		}

		// because we're only interested in the top-level locals block
		break
	}

	return r, nil
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

		gohcl.DecodeExpression(apisAttr.Expr, nil, &s)

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
