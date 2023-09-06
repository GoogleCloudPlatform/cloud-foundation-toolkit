package bptest

import (
	"fmt"
	"os"
	"path"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/util"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/discovery"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/iancoleman/strcase"
)

func initTest(name string) error {
	// check if test already exist
	testDir := path.Join(intTestPath, name)
	exists, err := util.Exists(testDir)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s already exists", testDir)
	}

	// write go mod if not exists
	goModpath := path.Join(intTestPath, goModFilename)
	exists, err = util.Exists(goModpath)
	if err != nil {
		return err
	}
	if !exists {
		goMod, err := getTmplFileContents(goModFilename)
		if err != nil {
			return err
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = writeFile(goModpath, fmt.Sprintf(goMod, path.Base(cwd)))
		if err != nil {
			return fmt.Errorf("error writing go mod file: %w", err)
		}
	}

	// discover test configs
	testCfg, err := discovery.GetConfigDirFromTestDir(testDir)
	if err != nil {
		return fmt.Errorf("unable to discover test configs for %s: %w", testDir, err)
	}

	// Parse config to expose outputs within test
	mod, diags := tfconfig.LoadModule(testCfg)
	if diags.HasErrors() {
		return fmt.Errorf("error parsing outputs: %w", diags)
	}
	outputs := make([]string, 0, len(mod.Outputs))
	for _, op := range mod.Outputs {
		// todo(bharathkkb): make templates type aware
		outputs = append(outputs, op.Name)
	}

	// render and write test
	testFile, err := getBPTestFromTmpl(name, outputs)
	if err != nil {
		return fmt.Errorf("error creating blueprint test: %w", err)
	}
	err = os.MkdirAll(testDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating test dir: %w", err)
	}
	return writeFile(path.Join(testDir, fmt.Sprintf("%s_test.go", strcase.ToSnake(name))), testFile)
}
