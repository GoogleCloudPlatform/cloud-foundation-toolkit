// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scorecard

import (
	"context"
	"time"

	"github.com/briandowns/spinner"

	asset "cloud.google.com/go/asset/apiv1"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

// Inventory manages a CAI inventory
type Inventory struct {
	ProjectID string
	Parent    string
	GcsPath   string
}

// NewInventory creates a new CAI inventory manager
func NewInventory(projectID string) *Inventory {
	inventory := new(Inventory)
	inventory.ProjectID = projectID
	inventory.Parent = "organizations/816421441114"
	inventory.GcsPath = "gs://clf-gcp-inventory/inventory.json"
	return inventory
}

// ExportInventory creates a new inventory export
func ExportInventory(inventory *Inventory) error {
	ctx := context.Background()
	c, err := asset.NewClient(ctx)
	if err != nil {
		return err
	}

	req := &assetpb.ExportAssetsRequest{
		Parent:      inventory.Parent,
		ContentType: assetpb.ContentType_RESOURCE,
		OutputConfig: &assetpb.OutputConfig{
			Destination: &assetpb.OutputConfig_GcsDestination{
				GcsDestination: &assetpb.GcsDestination{
					ObjectUri: &assetpb.GcsDestination_Uri{
						Uri: inventory.GcsPath,
					},
				},
			},
		},
	}

	op, err := c.ExportAssets(ctx, req)
	if err != nil {
		return err
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Exporting Cloud Asset Inventory to GCS bucket... "
	s.Start()
	_, err = op.Wait(ctx)
	s.Stop()

	if err != nil {
		return err
	}

	return nil
}
