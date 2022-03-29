package bptest

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	cb "google.golang.org/api/cloudbuild/v1"
	"sigs.k8s.io/yaml"
)

const (
	intTestPath          = "test/integration"
	intTestBuildFilePath = "build/int.cloudbuild.yaml"
	inspecInputsFile     = "inspec.yml"
	tmplSuffix           = ".tmpl"
	goModFilename        = "go.mod"
	bptTestFilename      = "blueprint_test.go"
)

var (
	//go:embed templates
	templateFiles          embed.FS
	kitchenCFTStageMapping = map[string]string{
		"create":   stages[0],
		"converge": stages[1],
		"verify":   stages[2],
		"destroy":  stages[3],
	}
)

type inspecInputs struct {
	Name       string `yaml:"name"`
	Attributes []struct {
		Name string `yaml:"name"`
	} `yaml:"attributes"`
}

// convertKitchenTests converts all kitchen tests to blueprint tests and updates build files
func convertKitchenTests() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// write go mod
	goMod, err := getTmplFileContents(goModFilename)
	if err != nil {
		return err
	}
	err = writeFile(path.Join(intTestPath, goModFilename), fmt.Sprintf(goMod, path.Base(cwd)))
	if err != nil {
		return fmt.Errorf("error writing go mod file: %v", err)
	}
	// write discover test
	discoverTest, err := getTmplFileContents(discoverTestFilename)
	if err != nil {
		return err
	}
	err = writeFile(path.Join(intTestPath, discoverTestFilename), string(discoverTest))
	if err != nil {
		return fmt.Errorf("error writing discover_test.go: %v", err)
	}
	testDirs, err := getCurrentTestDirs()
	if err != nil {
		return fmt.Errorf("error getting current test dirs: %v", err)
	}
	for _, dir := range testDirs {
		err = convertTest(path.Join(intTestPath, dir))
		if err != nil {
			return fmt.Errorf("error converting %s: %v", dir, err)
		}
	}
	// remove kitchen
	err = os.Remove(".kitchen.yml")
	if err != nil {
		return fmt.Errorf("error removing .kitchen.yml: %v", err)
	}
	// convert build file
	// We use build to identify commands to update and update the commands in the buildFile.
	// This minimizes unnecessary diffs in build yaml due to round tripping.
	build, buildFile, err := getBuildFromFile(intTestBuildFilePath)
	if err != nil {
		return fmt.Errorf("error unmarshalling %s: %v", intTestBuildFilePath, err)
	}
	newBuildFile, err := transformBuild(build, buildFile)
	if err != nil {
		return fmt.Errorf("error transforming buildfile: %v", err)
	}
	return writeFile(intTestBuildFilePath, newBuildFile)
}

// getCurrentTestDirs returns current test dirs in intTestPath
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

// convertTest converts a kitchen test in dir to blueprint test
func convertTest(dir string) error {
	// read inspec.yaml
	f, err := ioutil.ReadFile(path.Join(dir, inspecInputsFile))
	if err != nil {
		return fmt.Errorf("error reading inspec file: %s", err)
	}
	var inspec inspecInputs
	err = yaml.Unmarshal(f, &inspec)
	if err != nil {
		return fmt.Errorf("error unmarshalling inspec file: %s", err)
	}
	// get inspec input attributes
	var inputs []string
	for _, i := range inspec.Attributes {
		inputs = append(inputs, i.Name)
	}
	// get bpt skeleton
	testName := path.Base(dir)
	bpTest, err := getBPTestFromTmpl(testName, inputs)
	if err != nil {
		return fmt.Errorf("error creating blueprint test: %s", err)
	}
	// remove old test
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("error removing old test dir: %s", err)
	}
	// write bpt
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating test dir: %s", err)
	}
	return writeFile(path.Join(dir, fmt.Sprintf("%s_test.go", strcase.ToSnake(testName))), bpTest)
}

// getTmplFileContents returns contents of embedded file f
func getTmplFileContents(f string) (string, error) {
	tmplF := path.Join("templates", fmt.Sprintf("%s%s", f, tmplSuffix))
	contents, err := templateFiles.ReadFile(tmplF)
	if err != nil {
		return "", fmt.Errorf("error reading %s : %v", tmplF, err)
	}
	return string(contents), nil
}

// getTestFnName returns the go test function name
func getTestFnName(name string) string {
	return fmt.Sprintf("Test%s", strcase.ToCamel(name))
}

// getBPTestFromTmpl returns a skeleton blueprint test
func getBPTestFromTmpl(testName string, inputs []string) (string, error) {
	pkgName := strcase.ToSnake(testName)
	fnName := getTestFnName(testName)
	tmpl, err := getTmplFileContents(bptTestFilename)
	if err != nil {
		return "", err
	}
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

// writeFile writes content to file path
func writeFile(p string, content string) error {
	return ioutil.WriteFile(p, []byte(content), os.ModePerm)
}

// transformBuild transforms cloudbuild file contents with kitchen commands to CFT cli commands
func transformBuild(b *cb.Build, f string) (string, error) {
	for _, step := range b.Steps {
		// test commands have at least two args
		if len(step.Args) < 2 {
			continue
		}
		cmd := step.Args[len(step.Args)-1]
		// skip if not a kitchen command
		kitchenCmdIndex := strings.Index(cmd, "kitchen_do")
		if kitchenCmdIndex == -1 {
			continue
		}
		kitchenCmd := cmd[kitchenCmdIndex:]
		newCmd, err := getCFTCmd(kitchenCmd)
		if err != nil {
			return "", err
		}
		f = strings.ReplaceAll(f, cmd, newCmd)
	}
	return f, nil
}

// getCFTCmd returns an equivalent CFT command for a kitchen command
func getCFTCmd(kitchenCmd string) (string, error) {
	if !strings.Contains(kitchenCmd, "kitchen_do") {
		return "", fmt.Errorf("invalid kitchen command: %s", kitchenCmd)
	}
	cmdArr := strings.Split(kitchenCmd, " ")
	cftCmd := []string{"cft", "test", "run"}
	// cmd of form kitchen_do verb
	if len(cmdArr) == 2 {
		kitchenStage := cmdArr[len(cmdArr)-1]
		cftCmd = append(cftCmd, []string{"all", "--stage", kitchenCFTStageMapping[kitchenStage]}...)
	} else if len(cmdArr) == 3 {
		// cmd of form kitchen_do verb test-name
		kitchenTestName := cmdArr[len(cmdArr)-1]
		kitchenStage := cmdArr[len(cmdArr)-2]
		cftTestName := getTestFnName(strings.TrimSuffix(kitchenTestName, "-local"))
		cftCmd = append(cftCmd, []string{cftTestName, "--stage", kitchenCFTStageMapping[kitchenStage]}...)
	} else {
		return "", fmt.Errorf("unknown kitchen command: %s", kitchenCmd)
	}
	cftCmd = append(cftCmd, "--verbose")
	return strings.Join(cftCmd, " "), nil
}

// getBuildFromFile unmarshalls a cloudbuild file
func getBuildFromFile(fp string) (*cb.Build, string, error) {
	f, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, "", err
	}
	j, err := yaml.YAMLToJSON(f)
	if err != nil {
		return nil, "", err
	}
	var b cb.Build
	if err = json.Unmarshal(j, &b); err != nil {
		fmt.Println(err.Error())
	}
	return &b, string(f), nil
}
