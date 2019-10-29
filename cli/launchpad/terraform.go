package launchpad

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	tfTerraformVer          = ">= 0.12"
	tfProviderGoogleVer     = "~> 2.1"
	tfProviderGoogleBetaVer = "~> 2.1"
)

// tfConstruct represents a segment of Terraform code.
//
// A tfFile is expected to consist of multiple tfConstructs.
type tfConstruct interface {
	tfTemplate() string
	tfArguments() interface{}
}

// tfFile represents a single Terraform File.
type tfFile struct {
	filename   string // Filename without suffix
	dirname    string // ParentId directory name
	constructs []tfConstruct
}

// path generates the full filepath for a tfFile.
func (f *tfFile) path() string { return filepath.Join(f.dirname, fmt.Sprintf("%s.tf", f.filename)) }

// render produces Terraform code based on tfConstruct's template and arguments.
func (f *tfFile) render() string {
	var buff strings.Builder
	for _, tfc := range f.constructs {
		buff.WriteString(renderTFTemplate(tfc.tfTemplate(), tfc.tfArguments()))
	}
	return buff.String()
}

// newTfFile generates a tfFile with license attached.
func newTfFile(filename string, dirname string, constructs []tfConstruct) *tfFile {
	constructs = append([]tfConstruct{&tfLicense{}}, constructs...) // Auto prepend license
	return &tfFile{filename: filename, dirname: dirname, constructs: constructs}
}

// renderTFTemplate renders a given Terraform template and produces segment of Terraform code.
func renderTFTemplate(tmplfp string, data interface{}) string {
	tmplString, err := loadFile(tmplfp)
	if err != nil {
		log.Fatalln("Unable to load template", tmplfp, err)
	}
	var buff bytes.Buffer
	tmpl, err := template.New(tmplfp).Parse(tmplString)
	if err != nil {
		log.Fatalln("Unable to parse template", tmplfp, err)
	}
	if err := tmpl.Execute(&buff, data); err != nil {
		log.Fatalln("Unable to render template", tmplfp, err)
	}
	return buff.String()
}

// ==== Resources ====

// tfLicense is a tfConstruct license representation.
type tfLicense struct{}

func (t *tfLicense) tfTemplate() string       { return "launchpad/static/tmpl/tf/license.tf.tmpl" }
func (t *tfLicense) tfArguments() interface{} { return nil }

//
type tfTerraform struct {
	RequiredVersion string
}

func (t *tfTerraform) tfTemplate() string       { return "launchpad/static/tmpl/tf/_terraform.tf.tmpl" }
func (t *tfTerraform) tfArguments() interface{} { return t }
func newTfTerraform(requiredVer string) *tfTerraform {
	// TODO Default!
	return &tfTerraform{RequiredVersion: requiredVer}
}

// tfOutput represents a single Terraform "output".
type tfOutput struct {
	Id  string
	Val string
}

func (t *tfOutput) tfTemplate() string       { return "launchpad/static/tmpl/tf/_output.tf.tmpl" }
func (t *tfOutput) tfArguments() interface{} { return t }
func newTfOutput(id string, val string) *tfOutput {
	return &tfOutput{
		Id:  id,
		Val: val,
	}
}

// tfVariable represents a single Terraform "var"/"variable".
type tfVariable struct {
	Id          string
	Description string
	Default     string
}

func (t *tfVariable) tfTemplate() string       { return "launchpad/static/tmpl/tf/_variable.tf.tmpl" }
func (t *tfVariable) tfArguments() interface{} { return t }
func newTfVariable(id string, desc string, def string) *tfVariable {
	// TODO (@wengm) consider optional description & default
	v := &tfVariable{Id: id, Description: desc, Default: def}
	return v
}

// tfGoogleProvider represents any Google related Terraform "provider".
type tfGoogleProvider struct {
	Credentials string
	Version     string
}

func (tf *tfGoogleProvider) tfTemplate() string {
	return "launchpad/static/tmpl/tf/google_provider.tf.tmpl"
}
func (tf *tfGoogleProvider) tfArguments() interface{} { return tf }
func newTfGoogleProvider(options ...func(*tfGoogleProvider) error) *tfGoogleProvider {
	p := &tfGoogleProvider{
		Credentials: "var.credentials_file_path",
		Version:     tfProviderGoogleVer,
	}
	for _, op := range options {
		if err := op(p); err != nil {
			log.Fatalln("unable to process Google provider options")
		}
	}
	return p
}

// tfGoogleFolder represents a single GCP Folder in Terraform resource as "google_folder".
type tfGoogleFolder struct {
	Id          string
	DisplayName string
	Parent      string
}

func (t *tfGoogleFolder) tfTemplate() string       { return "launchpad/static/tmpl/tf/google_folder.tf.tmpl" }
func (t *tfGoogleFolder) tfArguments() interface{} { return *t }
func newTfGoogleFolder(id string, name string, parentPtr *parentRefYAML) *tfGoogleFolder {
	parent := ""
	switch parentPtr.ParentType {
	case KindOrganization:
		parent = fmt.Sprintf("\"organizations/${var.organization_id}\"")
	case KindFolder:
		parent = fmt.Sprintf("google_folder.%s.name", parentPtr.ParentId)
	default:
		log.Fatalln("folder contained in non folder or org")
	}
	return &tfGoogleFolder{
		Id:          id,
		DisplayName: name,
		Parent:      parent,
	}
}
