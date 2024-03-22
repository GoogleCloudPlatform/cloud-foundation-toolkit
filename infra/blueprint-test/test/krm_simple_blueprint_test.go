package test

import (
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/krmt"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/assert"
)

func TestKRMSimpleBlueprint(t *testing.T) {
	tfBlueprint := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("setup"),
	)
	gcloud.Runf(t, "container clusters get-credentials %s --region=%s --project %s -q", tfBlueprint.GetStringOutput("cluster_name"), tfBlueprint.GetStringOutput("cluster_region"), tfBlueprint.GetStringOutput("project_id"))

	networkBlueprint := krmt.NewKRMBlueprintTest(t,
		krmt.WithDir("../examples/simple_krm_blueprint"),
		krmt.WithUpdatePkgs(false),
	)
	networkBlueprint.DefineVerify(
		func(assert *assert.Assertions) {
			networkBlueprint.DefaultVerify(assert)
			k8sOpts := k8s.KubectlOptions{}
			op, err := k8s.RunKubectlAndGetOutputE(t, &k8sOpts, "get", "pod", "simple-krm-blueprint", "--no-headers", "-o", "custom-columns=:metadata.name")
			assert.NoError(err)
			result := strings.Split(op, "\n")
			assert.Equal("simple-krm-blueprint", result[len(result)-1])
		})
	networkBlueprint.Test()
}
