package bpmetadata

import (
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
	moduleVersion   string
	requiredVersion string
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

	//parse out the blueprint version from the config
	modName, diags := parseBlueprintVersion(versionsFile, diags)

	//parse out the required version from the config
	var hclModule tfconfig.Module
	hclModuleDiag := tfconfig.LoadModuleFromFile(versionsFile, &hclModule)
	diags = append(diags, hclModuleDiag...)

	for _, diag := range diags {
		if diag.Severity == hcl.DiagError {
			return nil
		}
	}

	requiredCore := ""
	if len(hclModule.RequiredCore) != 0 {
		//always looking for the first element since tf blueprints
		//have one required core version
		requiredCore = hclModule.RequiredCore[0]
	}

	return &blueprintVersion{
		moduleVersion:   modName,
		requiredVersion: requiredCore,
	}
}

// parseBlueprintVersion gets the blueprint version from the provided config
// from the provider_meta block
func parseBlueprintVersion(versionsFile *hcl.File, diags hcl.Diagnostics) (string, hcl.Diagnostics) {
	rootContent, _, rootContentDiags := versionsFile.Body.PartialContent(rootSchema)
	diags = append(diags, rootContentDiags...)
	for _, rootBlock := range rootContent.Blocks {
		if rootBlock.Type == "terraform" {
			tfContent, _, tfContentDiags := rootBlock.Body.PartialContent(metaSchema)
			diags = append(diags, tfContentDiags...)
			for _, tfContentBlock := range tfContent.Blocks {
				if tfContentBlock.Type == "provider_meta" {
					metaContent, _, metaContentDiags := tfContentBlock.Body.PartialContent(metaBlockSchema)
					diags = append(diags, metaContentDiags...)
					if versionAttr, defined := metaContent.Attributes["module_name"]; defined {
						var modName string
						gohcl.DecodeExpression(versionAttr.Expr, nil, &modName)
						re := regexp.MustCompile(versionRegEx)
						m := re.FindStringSubmatch(modName)
						if len(m) > 0 {
							return m[len(m)-1], diags
						}

						return "", diags
					}
				}
			}
			break
		}
	}

	return "", diags
}

// getBlueprintInterfaces gets the variables and outputs associated
// with the blueprint
func getBlueprintInterfaces(configPath string) *BlueprintInterface {
	//load the configs from the dir path
	mod, diags := tfconfig.LoadModule(configPath)
	for _, diag := range diags {
		if diag.Severity == tfconfig.DiagError {
			return nil
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

	i := &BlueprintInterface{
		Variables: variables,
		Outputs:   outputs,
	}

	return i
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

// getBlueprintRequirements gets the variables and outputs associated
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
	iamContent, _, _ := rolesFile.Body.PartialContent(rootSchema)
	for _, block := range iamContent.Blocks {
		if block.Type == "locals" {
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
	}

	return r
}

// parseBlueprintServices gets the gcp api services required for the blueprint
// to be provisioned
func parseBlueprintServices(servicesFile *hcl.File) []string {
	var s []string
	servicesContent, _, _ := servicesFile.Body.PartialContent(rootSchema)
	for _, block := range servicesContent.Blocks {
		if block.Type == "locals" {
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
	}

	return s
}
