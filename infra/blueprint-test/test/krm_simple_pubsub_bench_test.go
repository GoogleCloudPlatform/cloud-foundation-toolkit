package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/kpt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/krmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/otiai10/copy"
)

// generateNTopics generates a slice of topicCount topic names
func generateNTopics(topicCount int) []string {
	var m []string
	for i := 0; i < topicCount; i++ {
		m = append(m, fmt.Sprintf("topic-%d", i))
	}
	return m
}

// kubectlWaitForDeletion waits for resources in dir to be deleted
// workaround until https://github.com/GoogleContainerTools/kpt/issues/2374
func kubectlWaitForDeletion(b *testing.B, dir string, retries int, retryInterval time.Duration) {
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

// createVariant creates a variant of baseDir blueprint in the buildDir/variantName and upserts any given setters for that variant
func createVariant(b *testing.B, baseDir string, buildDir string, variantName string, setters map[string]string) {
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
	kpt.UpsertSetters(rs, setters)
	err = kpt.WritePkgResources(variantPath, rs)
	if err != nil {
		b.Fatalf("unable to write resources in %s :%v", variantPath, err)
	}
}

// getBuildDir creates a directory to store generated variants and cleanup fn
func getBuildDir(b *testing.B) (string, func()) {
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

func BenchmarkKRMPubSubBlueprint(b *testing.B) {
	projectID := utils.ValFromEnv(b, "PROJECT_ID")
	// base blueprint dir
	blueprintDir := "benchmark_fixtures/simple_pubsub_krm"
	// topicCounts := []int{10, 50, 100, 500, 1000}
	topicCounts := []int{10}
	for _, topicCount := range topicCounts {
		b.Run(fmt.Sprintf("PubSub Bench mark with %d topics", topicCount), func(b *testing.B) {
			// we precreate a custom build directory to generate variants for a given resource blueprint
			buildDir, cleanup := getBuildDir(b)
			defer cleanup()
			// init empty kpt pkg in the build dir
			kptHelper := kpt.NewCmdConfig(b, kpt.WithDir(buildDir))
			kptHelper.RunCmd("pkg", "init")
			// generate package variants into the build dir
			topicNames := generateNTopics(topicCount)
			for _, name := range topicNames {
				createVariant(b, blueprintDir, buildDir, name, map[string]string{"topic-name": name, "project-id": projectID})
			}
			// render variants
			// TODO(bharathkkb): this is currently done in serial by kpt and can be slow for bigger topicCounts
			// We should look into doing this in parallel possibly bundling variant creation with rendering
			kptHelper.RunCmd("fn", "render")
			kptHelper.RunCmd("live", "install-resource-group")
			kptHelper.RunCmd("live", "init")
			pubsubTest := krmt.NewKRMBlueprintTest(b, krmt.WithDir(buildDir), krmt.WithBuildDir(buildDir), krmt.WithUpdatePkgs(false))
			b.ResetTimer()
			// start benchmark
			for n := 0; n < b.N; n++ {
				pubsubTest.Apply(nil)
				b.StopTimer()
				// stop benchmark
				pubsubTest.Teardown(nil)
				// confirm resources are deleted
				kubectlWaitForDeletion(b, buildDir, 50, 5*time.Second)
				// restart timer
				b.StartTimer()
			}

		})
	}
}
