/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package test

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/tft"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/utils"
)

func generateNTopicsPerProject(projects []string, topicCount int) map[string]interface{} {
	// split topics across available projects
	n := topicCount / len(projects)
	m := make(map[string]interface{})
	for _, p := range projects {
		for i := 0; i < n; i++ {
			topic := fmt.Sprintf("topic-%d", i)
			m[fmt.Sprintf("%s/%s", p, topic)] = map[string]string{"project": p, "topic": topic}
		}
	}
	return m
}

func BenchmarkTFPubSub(b *testing.B) {
	// benchmarks to run
	// topicCounts := []int{10, 50, 100, 500, 1000}
	topicCounts := []int{10}
	for _, topicCount := range topicCounts {
		b.Run(fmt.Sprintf("PubSub Bench mark with %d topics", topicCount), func(b *testing.B) {
			pubSubTest := tft.Init(b,
				tft.WithSetupPath("setup/simple_tf_bench"),
				tft.WithTFDir("benchmark_fixtures/simple_pubsub_tf"),
				// tft.WithLogger(logger.Discard),
			)
			// get list of available projects that have been setup
			project_ids := pubSubTest.GetTFSetupOPListVal("project_ids")
			// create input as vars for TF config with topics split across available projects
			tfVars := map[string]interface{}{"project_topic_map": generateNTopicsPerProject(project_ids, topicCount)}
			tft.WithVars(tfVars)(pubSubTest)
			// run tf init to download provider(s)
			utils.RunStage("setup", func() { pubSubTest.Setup() })
			// reset benchmark timer to ignore previous time
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				// start apply benchmark
				utils.RunStage("apply", func() { pubSubTest.Apply() })
				//stop timer for cleanup
				b.StopTimer()
				utils.RunStage("destroy", func() { pubSubTest.Teardown() })
				// restart timer
				b.StartTimer()
			}
		})
	}
}
