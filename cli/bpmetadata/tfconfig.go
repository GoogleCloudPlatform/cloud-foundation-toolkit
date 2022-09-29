package bpmetadata

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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

// getBlueprintVersion gets both the required core version and the
// version of the blueprint
func getBlueprintVersion(configPath string) *blueprintVersion {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	//create hcl file object from the provided tf config
	fileName := filepath.Base(configPath)
	var diags hcl.Diagnostics
	p := hclparse.NewParser()
	versionsFile, fileDiags := p.ParseHCL(bytes, fileName)
	diags = append(diags, fileDiags...)
	if hasHclErrors(diags) {
		return nil
	}

	//parse out the blueprint version from the config
	modName, diags := parseBlueprintVersion(versionsFile, diags)

	//parse out the required version from the config
	var hclModule tfconfig.Module
	hclModule.RequiredProviders = make(map[string]*tfconfig.ProviderRequirement)
	hclModuleDiag := tfconfig.LoadModuleFromFile(versionsFile, &hclModule)
	diags = append(diags, hclModuleDiag...)
	if hasHclErrors(diags) {
		return nil
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
	}
}

// parseBlueprintVersion gets the blueprint version from the provided config
// from the provider_meta block
func parseBlueprintVersion(versionsFile *hcl.File, diags hcl.Diagnostics) (string, hcl.Diagnostics) {
	re := regexp.MustCompile(versionRegEx)
	// PartialContent() returns TF content containing blocks and attributes
	// based on the provided schema
	rootContent, _, rootContentDiags := versionsFile.Body.PartialContent(rootSchema)
	diags = append(diags, rootContentDiags...)
	if hasHclErrors(diags) {
		return "", diags
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
		for _, tfContentBlock := range tfContent.Blocks {
			if tfContentBlock.Type != "provider_meta" {
				continue
			}

			// this PartialContent() call with get the module_name attribute
			// that contains the version info
			metaContent, _, metaContentDiags := tfContentBlock.Body.PartialContent(metaBlockSchema)
			diags = append(diags, metaContentDiags...)
			versionAttr, defined := metaContent.Attributes["module_name"]
			if !defined {
				return "", diags
			}

			// get the module name from the version attribute and extract the
			// version name only
			var modName string
			gohcl.DecodeExpression(versionAttr.Expr, nil, &modName)
			m := re.FindStringSubmatch(modName)
			if len(m) > 0 {
				return m[len(m)-1], diags
			}

			return "", diags
		}

		break
	}

	return "", diags
}

// getBlueprintInterfaces gets the variables and outputs associated
// with the blueprint
func getBlueprintInterfaces(configPath string) (*BlueprintInterface, error) {
	//load the configs from the dir path
	mod, diags := tfconfig.LoadModule(configPath)
	for _, diag := range diags {
		if diag.Severity == tfconfig.DiagError {
			return nil, fmt.Errorf("unable to load module: %v", diag)
		}
	}

	var variables []BlueprintVariable
	for _, val := range mod.Variables {
		v := getBlueprintVariable(val)
		variables = append(variables, v)
	}

	var outputs []BlueprintOutput
	for _, val := range mod.Outputs {
		o := getBlueprintOutput(val)

		outputs = append(outputs, o)
	}

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
func getBlueprintRequirements(rolesConfigPath, servicesConfigPath string) BlueprintRequirements {
	//parse blueprint roles
	p := hclparse.NewParser()
	rolesFile, _ := p.ParseHCLFile(rolesConfigPath)
	r := parseBlueprintRoles(rolesFile)

	//parse blueprint services
	servicesFile, _ := p.ParseHCLFile(servicesConfigPath)
	s := parseBlueprintServices(servicesFile)

	return BlueprintRequirements{
		Roles:    r,
		Services: s,
	}
}

// parseBlueprintRoles gets the roles required for the blueprint to be provisioned
func parseBlueprintRoles(rolesFile *hcl.File) []BlueprintRoles {
	var r []BlueprintRoles
	iamContent, _, diags := rolesFile.Body.PartialContent(rootSchema)
	if hasHclErrors(diags) {
		return r
	}

	for _, block := range iamContent.Blocks {
		if block.Type != "locals" {
			continue
		}

		iamAttrs, _ := block.Body.JustAttributes()
		for k, _ := range iamAttrs {
			var iamRoles []string
			attrValue, _ := iamAttrs[k].Expr.Value(nil)
			ie := attrValue.ElementIterator()
			for ie.Next() {
				_, v := ie.Element()
				iamRoles = append(iamRoles, v.AsString())
			}

			containerRoles := BlueprintRoles{
				// TODO: (b/248123274) no good way to associate granularity yet
				Granularity: "",
				Roles:       iamRoles,
			}

			r = append(r, containerRoles)
		}

		// because we're only interested in the top-level locals block
		break
	}

	return r
}

// parseBlueprintServices gets the gcp api services required for the blueprint
// to be provisioned
func parseBlueprintServices(servicesFile *hcl.File) []string {
	var s []string
	servicesContent, _, diags := servicesFile.Body.PartialContent(rootSchema)
	if hasHclErrors(diags) {
		return s
	}

	for _, block := range servicesContent.Blocks {
		if block.Type != "locals" {
			continue
		}

		serivceAttrs, _ := block.Body.JustAttributes()
		for k, _ := range serivceAttrs {
			attrValue, _ := serivceAttrs[k].Expr.Value(nil)
			ie := attrValue.ElementIterator()
			for ie.Next() {
				_, v := ie.Element()
				s = append(s, v.AsString())
			}
		}

		// because we're only interested in the top-level locals block
		break
	}

	return s
}

func hasHclErrors(diags hcl.Diagnostics) bool {
	for _, diag := range diags {
		if diag.Severity == hcl.DiagError {
			return true
		}
	}

	return false
}
