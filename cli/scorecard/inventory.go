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
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pkg/errors"

	asset "cloud.google.com/go/asset/apiv1"
	assetpb "cloud.google.com/go/asset/apiv1/assetpb"
)

// InventoryConfig manages a CAI inventory
type InventoryConfig struct {
	projectID      string
	folderID       string
	organizationID string
	bucketName     string
	dirPath        string
	readFromStdin  bool
	workers        int
}

// Option for NewInventory
type Option func(*InventoryConfig)

// TargetProject sets the project for collecting inventory data
func TargetProject(projectID string) Option {
	return func(inventory *InventoryConfig) {
		inventory.projectID = projectID
	}
}

// TargetFolder sets the folder for collecting inventory data
func TargetFolder(folderID string) Option {
	return func(inventory *InventoryConfig) {
		inventory.folderID = folderID
	}
}

// TargetOrg sets the organzation for collecting inventory data
func TargetOrg(organizationID string) Option {
	return func(inventory *InventoryConfig) {
		inventory.organizationID = organizationID
	}
}

// WorkerSize sets the number of workers for running violations review concurrently
func WorkerSize(workers int) Option {
	return func(inventory *InventoryConfig) {
		inventory.workers = workers
	}
}

// NewInventory creates a new CAI inventory manager
func NewInventory(bucketName, dirPath string, readFromStdin bool, refresh bool, options ...Option) (*InventoryConfig, error) {
	inventory := new(InventoryConfig)
	inventory.bucketName = bucketName
	inventory.dirPath = dirPath
	inventory.readFromStdin = readFromStdin

	for _, option := range options {
		option(inventory)
	}

	Log.Debug("Initializing inventory", "target", inventory.getParent())
	if refresh {
		err := inventory.Export()
		if err != nil {
			return nil, err
		}
	}
	return inventory, nil
}

func (inventory InventoryConfig) getParent() string {
	if inventory.organizationID != "" {
		return fmt.Sprintf("organizations/%v", inventory.organizationID)
	} else if inventory.folderID != "" {
		return fmt.Sprintf("folders/%v", inventory.folderID)
	}
	return fmt.Sprintf("projects/%v", inventory.projectID)
}

// destinationObjectNames maps the different export types to their expected file location
var destinationObjectNames = map[assetpb.ContentType]string{
	assetpb.ContentType_RESOURCE:      "resource_inventory.json",
	assetpb.ContentType_IAM_POLICY:    "iam_inventory.json",
	assetpb.ContentType_ORG_POLICY:    "org_policy_inventory.json",
	assetpb.ContentType_ACCESS_POLICY: "access_policy_inventory.json",
}

func (inventory InventoryConfig) getGcsDestination(contentType assetpb.ContentType) *assetpb.GcsDestination_Uri {
	objectName := destinationObjectNames[contentType]
	return &assetpb.GcsDestination_Uri{
		Uri: fmt.Sprintf("gs://%v/%v", inventory.bucketName, objectName),
	}
}

// exportToGcs exports an inventory of the given resource type to GCS
func (inventory InventoryConfig) exportToGcs(contentType assetpb.ContentType) error {
	ctx := context.Background()
	c, err := asset.NewClient(ctx)
	if err != nil {
		return err
	}

	destination := inventory.getGcsDestination(contentType)
	req := &assetpb.ExportAssetsRequest{
		Parent:      inventory.getParent(),
		ContentType: contentType,
		OutputConfig: &assetpb.OutputConfig{
			Destination: &assetpb.OutputConfig_GcsDestination{
				GcsDestination: &assetpb.GcsDestination{
					ObjectUri: destination,
				},
			},
		},
	}
	Log.Debug("Exporting Asset ", "contentType", contentType, "parent", inventory.getParent())
	op, err := c.ExportAssets(ctx, req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("destination = %v", destination))
	}

	_, err = op.Wait(ctx)
	return err
}

// Export creates a new inventory export
func (inventory *InventoryConfig) Export() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Exporting Cloud Asset Inventory to GCS bucket... "
	s.Start()
	err := inventory.exportToGcs(assetpb.ContentType_RESOURCE)
	if err != nil {
		s.Stop()
		return err
	}
	err = inventory.exportToGcs(assetpb.ContentType_IAM_POLICY)
	if err != nil {
		s.Stop()
		return err
	}
	err = inventory.exportToGcs(assetpb.ContentType_ORG_POLICY)
	if err != nil {
		s.Stop()
		return err
	}
	err = inventory.exportToGcs(assetpb.ContentType_ACCESS_POLICY)
	if err != nil {
		s.Stop()
		return err
	}
	s.Stop()

	return nil
}
