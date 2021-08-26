package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/benchmark"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/kpt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/krmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
)

// generateNTopics generates a slice of topicCount topic names
func generateNTopics(topicCount int) []string {
	var m []string
	for i := 0; i < topicCount; i++ {
		m = append(m, fmt.Sprintf("topic-%d", i))
	}
	return m
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
			buildDir, cleanup := benchmark.GetBuildDir(b)
			defer cleanup()
			// init empty kpt pkg in the build dir
			kptHelper := kpt.NewCmdConfig(b, kpt.WithDir(buildDir))
			kptHelper.RunCmd("pkg", "init")
			// generate package variants into the build dir
			topicNames := generateNTopics(topicCount)
			for _, name := range topicNames {
				benchmark.CreateVariant(b, blueprintDir, buildDir, name, map[string]string{"topic-name": name, "project-id": projectID})
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
				benchmark.KubectlWaitForDeletion(b, buildDir, 50, 5*time.Second)
				// restart timer
				b.StartTimer()
			}

		})
	}
}
