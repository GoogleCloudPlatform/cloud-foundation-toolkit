package vpc_custom_test

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/krmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestVPCCustomBlueprint(t *testing.T) {
	networkName := "test-network"
	networkBlueprint := krmt.NewKRMBlueprintTest(t,
		krmt.WithUpdateCommit("2b93fd6d4f1a3eabdf4dfce05b93ccb1f9f671c5"),
		krmt.WithSetters(map[string]string{"network-name": networkName}),
	)
	networkBlueprint.DefineVerify(
		func(assert *assert.Assertions) {
			networkBlueprint.DefaultVerify(assert)
			op := gcloud.Run(t, fmt.Sprintf("compute networks describe %s --project %s", networkName, utils.ValFromEnv(t, "PROJECT_ID")))
			assert.Equal("GLOBAL", op.Get("routingConfig.routingMode").String(), "should be GLOBAL")
			assert.Equal("false", op.Get("autoCreateSubnetworks").String(), "autoCreateSubnetworks should not be enabled")
		})
	networkBlueprint.Test()
}
