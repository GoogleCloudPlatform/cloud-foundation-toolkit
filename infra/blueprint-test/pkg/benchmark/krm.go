package benchmark

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/kpt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/krmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/mitchellh/go-testing-interface"
	"github.com/otiai10/copy"
)

// KubectlWaitForDeletion waits for resources in dir to be deleted.
// Workaround until https://github.com/GoogleContainerTools/kpt/issues/2374
func KubectlWaitForDeletion(b testing.TB, dir string, retries int, retryInterval time.Duration) {
	waitArgs := []string{"get", "-R", "-f", ".", "--ignore-not-found"}
	kptCmd := shell.Command{
		Command:    "kubectl",
		Args:       waitArgs,
		Logger:     utils.GetLoggerFromT(),
		WorkingDir: dir,
	}
	waitFunction := func() (bool, error) {
		op, err := shell.RunCommandAndGetStdOutE(b, kptCmd)
		if err != nil {
			return false, err
		}
		// retry if output is not empty
		retry := op != ""
		return retry, nil
	}
	utils.Poll(b, waitFunction, retries, retryInterval)
}

// CreateVariant creates a variant of baseDir blueprint in the buildDir/variantName and upserts any given setters for that variant.
func CreateVariant(b testing.TB, baseDir string, buildDir string, variantName string, setters map[string]string) {
	for _, p := range []string{baseDir, buildDir} {
		_, err := os.Stat(p)
		if err != nil {
			b.Fatalf("%s does not exist", p)
		}
	}
	variantPath := path.Join(buildDir, variantName)
	err := copy.Copy(baseDir, variantPath)
	if err != nil {
		b.Fatalf("Error copying resource from %s to %s", baseDir, variantPath)
	}
	rs, err := kpt.ReadPkgResources(variantPath)
	if err != nil {
		b.Fatalf("unable to read resources in %s :%v", variantPath, err)
	}
	if err := kpt.UpsertSetters(rs, setters); err != nil {
		b.Fatalf("unable to upsert setters in %s :%v", variantPath, err)
	}
	err = kpt.WritePkgResources(variantPath, rs)
	if err != nil {
		b.Fatalf("unable to write resources in %s :%v", variantPath, err)
	}
}

// GetBuildDir creates a directory to store generated variants and cleanup fn.
func GetBuildDir(b testing.TB) (string, func()) {
	buildDir := path.Join(".build", b.Name())
	err := os.MkdirAll(buildDir, 0755)
	if err != nil {
		b.Fatalf("unable to create %s :%v", buildDir, err)
	}
	abs, err := filepath.Abs(buildDir)
	if err != nil {
		b.Fatalf("unable to get absolute path for %s :%v", buildDir, err)
	}
	rmBuildDir := func() {
		os.RemoveAll(buildDir)
	}
	return abs, rmBuildDir
}

// CreateTestVariant creates variants of a baseBlueprint and renders them with variantSetters
func CreateTestVariant(b testing.TB, baseBlueprintDir string, variantSetters map[string]map[string]string) (*krmt.KRMBlueprintTest, string, func()) {
	// precreate a custom build directory to generate variants for a given resource blueprint
	buildDir, cleanup := GetBuildDir(b)
	// init empty kpt pkg in the build dir
	kptHelper := kpt.NewCmdConfig(b, kpt.WithDir(buildDir))
	// generate package variants into the build dir
	kptHelper.RunCmd("pkg", "init")
	// generate variants of the base blueprint
	for name, setters := range variantSetters {
		CreateVariant(b, baseBlueprintDir, buildDir, name, setters)
	}
	// render variants
	// TODO(bharathkkb): this is currently done in serial by kpt and can be slow for bigger topicCounts
	// We should look into doing this in parallel possibly bundling variant creation with rendering
	kptHelper.RunCmd("fn", "render")
	kptHelper.RunCmd("live", "install-resource-group")
	kptHelper.RunCmd("live", "init")
	return krmt.NewKRMBlueprintTest(b, krmt.WithDir(buildDir), krmt.WithBuildDir(buildDir), krmt.WithUpdatePkgs(false)), buildDir, cleanup
}
