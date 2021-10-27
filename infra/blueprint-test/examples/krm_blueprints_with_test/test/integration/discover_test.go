package test

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/krmt"
)

func TestAll(t *testing.T) {
	krmt.AutoDiscoverAndTest(t)
}
