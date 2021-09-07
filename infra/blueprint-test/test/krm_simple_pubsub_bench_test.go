package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/benchmark"
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
			// generate setters for each variants
			topicNames := generateNTopics(topicCount)
			variantSetters := make(map[string]map[string]string)
			for _, name := range topicNames {
				variantSetters[name] = map[string]string{"topic-name": name, "project-id": projectID}
			}
			pubsubTest, buildDir, cleanup := benchmark.CreateTestVariant(b, blueprintDir, variantSetters)
			defer cleanup()
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
