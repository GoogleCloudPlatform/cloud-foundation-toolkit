package simple

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestCFTSimpleModule(t *testing.T) {
	networkBlueprint := tft.NewTFBlueprintTest(t)
	networkBlueprint.DefineVerify(
		func(assert *assert.Assertions) {
			networkBlueprint.DefaultVerify(assert)
			op := gcloud.Run(t, fmt.Sprintf("compute networks subnets describe subnet-01 --project %s --region us-west1", networkBlueprint.GetStringOutput("project_id")))
			assert.Equal("10.10.10.0/24", op.Get("ipCidrRange").String(), "should have the right CIDR")
			assert.Equal("false", op.Get("logConfig.enable").String(), "logConfig should not be enabled")
		})
	networkBlueprint.Test()
}
