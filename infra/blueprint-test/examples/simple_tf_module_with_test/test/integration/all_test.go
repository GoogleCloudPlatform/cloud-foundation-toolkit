package test

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
)

func TestAll(t *testing.T) {
	tft.AutoDiscoverAndTest(t)
}
