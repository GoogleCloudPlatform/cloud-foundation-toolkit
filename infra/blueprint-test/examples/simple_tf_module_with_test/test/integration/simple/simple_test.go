package simple

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/bpt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/modules/tft"
	"github.com/stretchr/testify/assert"
)

func TestCFTSimpleModule(t *testing.T) {
	nt := tft.Init(t)
	bpt.TestBlueprint(t, nt, func(tft *bpt.BlueprintTest) {
		tft.DefineVerify(func(assert *assert.Assertions) {
			op := gcloud.Run(t, fmt.Sprintf("compute networks subnets describe subnet-01 --project %s --region us-west1", nt.GetStringOutput("project_id")))
			assert.Equal(op.Get("ipCidrRange").String(), "10.10.10.0/24", "should have the right CIDR")
			assert.Equal(op.Get("logConfig.enable").String(), "false", "logConfig should not be enabled")
		})
	})
}
