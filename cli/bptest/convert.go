package bptest

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v2"
)

const (
	intTestPath      = "test/integration"
	inspecInputsFile = "inspec.yml"
	discoverTest     = `package test

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
)

func TestAll(t *testing.T) {
	tft.AutoDiscoverAndTest(t)
}`
)

func getCurrentTestDirs() ([]string, error) {
	files, err := ioutil.ReadDir(intTestPath)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	return dirs, nil
}

type inspecInputs struct {
	Name       string `yaml:"name"`
	Attributes []struct {
		Name string `yaml:"name"`
	} `yaml:"attributes"`
}

func convertTest(dir string) error {
	// read inspec.yaml
	f, err := ioutil.ReadFile(path.Join(dir, inspecInputsFile))
	if err != nil {
		return err
	}
	var inspec inspecInputs
	err = yaml.Unmarshal(f, &inspec)
	if err != nil {
		return err
	}

	// get inputs
	var inputs []string
	for _, i := range inspec.Attributes {
		inputs = append(inputs, i.Name)
	}

	// get bpt skeleton
	testName := path.Base(dir)
	bpTest, err := getBPTestTmpl(testName, inputs)
	if err != nil {
		return err
	}
	// remove old test
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	// write bpt
	err = ioutil.WriteFile(path.Join(dir, fmt.Sprintf("%s_test.go", strcase.ToSnake(testName))), []byte(bpTest), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func getBPTestTmpl(testName string, inputs []string) (string, error) {
	pkgName := strcase.ToSnake(testName)
	fnName := fmt.Sprintf("Test%s", strcase.ToCamel(testName))
	tmpl := `package {{.PkgName}}

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func {{.FnName}}(t *testing.T) {
	bpt := tft.NewTFBlueprintTest(t)

	bpt.DefineVerify(func(assert *assert.Assertions) {
		bpt.DefaultVerify(assert)
		{{range .Inputs}}
		{{toLowerCamel .}} := bpt.GetStringOutput("{{.}}"){{end}}

		op := gcloud.Run(t,"")
		assert.Contains(op.Get("result").String(), randomFileString, "contains random string")
	})

	bpt.Test()
}`
	t, err := template.New("test").Funcs(template.FuncMap{"toLowerCamel": strcase.ToLowerCamel}).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, struct {
		PkgName string
		FnName  string
		Inputs  []string
	}{
		PkgName: pkgName,
		FnName:  fnName,
		Inputs:  inputs,
	},
	)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func getGoMod(dir string) string {
	return fmt.Sprintf(`module github.com/terraform-google-modules/%s/test/integration

go 1.16

require (
	github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test v0.0.0-20220203183701-08c972c3768c
	github.com/gruntwork-io/terratest v0.35.6 // indirect
	github.com/stretchr/testify v1.7.0
)
`, path.Base(dir))
}

func writeFile(p string, content string) error {
	return ioutil.WriteFile(path.Join(intTestPath, p), []byte(content), os.ModePerm)
}

func convert() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	currentMod := path.Base(cwd)
	// write go mod
	err = writeFile("go.mod", getGoMod(currentMod))
	if err != nil {
		return err
	}
	// write discover test
	err = writeFile("discover_test.go", discoverTest)
	if err != nil {
		return err
	}
	testDirs, err := getCurrentTestDirs()
	if err != nil {
		return err
	}
	for _, dir := range testDirs {
		err = convertTest(path.Join(intTestPath, dir))
		if err != nil {
			return err
		}
	}
	return nil
}
