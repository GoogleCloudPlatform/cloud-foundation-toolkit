package launchpad

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

type tfFile struct {
	filename   string // Filename without suffix
	dirname    string // ParentId directory name
	constructs []tfConstruct
}

type tfConstruct interface {
	tfTemplate() string
	tfArguments() interface{}
}

func newTfFile(filename string, dirname string, constructs []tfConstruct) *tfFile {
	// All new files should have license attached
	constructs = append([]tfConstruct{&tfLicense{}}, constructs...)
	return &tfFile{filename: filename, dirname: dirname, constructs: constructs}
}

func (f *tfFile) path() string {
	return filepath.Join(f.dirname, fmt.Sprintf("%s.tf", f.filename))
}

func (f *tfFile) render() string {
	var buff strings.Builder
	for _, tfc := range f.constructs {
		buff.WriteString(renderTFTemplate(tfc.tfTemplate(), tfc.tfArguments()))
	}
	return buff.String()
}

func renderTFTemplate(tmplfp string, data interface{}) string {
	var buff bytes.Buffer
	tmpl := template.New(tmplfp)
	tmpl, err := tmpl.Parse(loadFile(tmplfp))
	if err != nil {
		log.Fatalln("Unable to load template", tmplfp, err)
	}
	err = tmpl.Execute(&buff, data)
	if err != nil {
		log.Fatalln("Unable to render template", tmplfp, err)
	}
	return buff.String()
}

// ==== Resources ====
type tfLicense struct{}

func (t *tfLicense) tfTemplate() string       { return "launchpad/static/tmpl/tf/license.tf.tmpl" }
func (t *tfLicense) tfArguments() interface{} { return nil }

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
		Credentials: "${file(var.credentials_file_path)}",
		Version:     "~> 1.19",
	}
	for _, op := range options {
		err := op(p)
		if err != nil {
			log.Fatalln("unable to process Google provider options")
		}
	}
	return p
}

type tfGoogleFolder struct {
	Id          string
	DisplayName string
	Parent      string
}

func (t *tfGoogleFolder) tfTemplate() string       { return "launchpad/static/tmpl/tf/google_folder.tf.tmpl" }
func (t *tfGoogleFolder) tfArguments() interface{} { return *t }
func newTfGoogleFolder(id string, name string, parent string) *tfGoogleFolder {
	return &tfGoogleFolder{
		Id:          id,
		DisplayName: name,
		Parent:      parent,
	}
}
