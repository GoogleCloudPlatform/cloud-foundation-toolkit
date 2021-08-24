package vpc_custom_test

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/krmt"
	"github.com/stretchr/testify/assert"
)

func TestVPCCustomBlueprint(t *testing.T) {
	networkBlueprint := krmt.NewKRMBlueprintTest(t,
		krmt.WithUpdateCommit("2b93fd6d4f1a3eabdf4dfce05b93ccb1f9f671c5"),
	)
	networkBlueprint.DefineVerify(
		func(assert *assert.Assertions) {
			networkBlueprint.DefaultVerify(assert)
			op := gcloud.Run(t, fmt.Sprintf("compute networks describe custom-network --project %s", networkBlueprint.ValFromEnv("PROJECT_ID")))
			assert.Equal("GLOBAL", op.Get("routingConfig.routingMode").String(), "should be GLOBAL")
			assert.Equal("false", op.Get("autoCreateSubnetworks").String(), "autoCreateSubnetworks should not be enabled")
		})
	networkBlueprint.Test()
}
