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
	"io"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/logger"
	tt "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
)

type customLogger struct {
	w io.Writer
}

func (c *customLogger) Logf(t tt.TestingT, format string, args ...interface{}) {
	_, err := fmt.Fprintf(c.w, format, args...)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fmt.Fprintln(c.w)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSimpleTFModule(t *testing.T) {
	path, _ := os.Getwd()

	// Regular logger that also writes to stdout.
	var regularLogs strings.Builder
	regularWriter := io.MultiWriter(&regularLogs, os.Stdout)
	fakeRegularWriter := &customLogger{
		w: regularWriter,
	}
	regularLogger := logger.New(fakeRegularWriter)

	// Sensitive logger to capture sensitive output.
	var sensitiveLogs strings.Builder
	fakeSensitiveWriter :=
		&customLogger{
			w: &sensitiveLogs,
		}

	sensitiveLogger := logger.New(fakeSensitiveWriter)
	statePath := fmt.Sprintf("%s/../examples/simple_tf_module/local_backend.tfstate", path)
	nt := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../examples/simple_tf_module"),
		tft.WithBackendConfig(map[string]interface{}{
			"path": statePath,
		}),
		tft.WithSetupPath("setup"),
		tft.WithEnvVars(map[string]string{"network_name": fmt.Sprintf("foo-%s", utils.RandStr(5))}),
		tft.WithLogger(regularLogger),
		tft.WithSensitiveLogger(sensitiveLogger),
	)

	utils.RunStage("init", func() { nt.Init(nil) })
	defer utils.RunStage("teardown", func() { nt.Teardown(nil) })

	utils.RunStage("plan", func() { nt.Plan(nil) })
	utils.RunStage("apply", func() { nt.Apply(nil) })

	utils.RunStage("verify", func() {
		assert := assert.New(t)
		nt.Verify(assert)
		op := gcloud.Run(t, fmt.Sprintf("compute networks subnets describe subnet-01 --project %s --region us-west1", nt.GetStringOutput("project_id")))
		assert.Equal("10.10.10.0/24", op.Get("ipCidrRange").String(), "should have the right CIDR")
		assert.Equal("false", op.Get("logConfig.enable").String(), "logConfig should not be enabled")
		assert.FileExists(statePath)
	})

	// sa_key is a sensitive output key from setup.
	sensitiveOP := "sa_key"
	if strings.Contains(regularLogs.String(), sensitiveOP) {
		t.Errorf("regular logs should not contain sensitive output")
	}
	if !strings.Contains(sensitiveLogs.String(), sensitiveOP) {
		t.Errorf("sensitive logs should contain sensitive output")
	}

	// Custom plan function not defined, plan should be skipped.
	if !strings.Contains(regularLogs.String(), "skipping plan as no function defined") {
		t.Errorf("plan should be skipped")
	}
}
