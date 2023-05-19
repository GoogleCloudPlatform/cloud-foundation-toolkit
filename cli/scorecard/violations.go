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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/config-validator/pkg/api/validator"
	cvasset "github.com/GoogleCloudPlatform/config-validator/pkg/asset"
	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
)

func getDataFromReader(reader io.Reader) ([]*validator.Asset, error) {
	const maxCapacity = 1024 * 1024
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	var pbAssets []*validator.Asset
	for scanner.Scan() {
		pbAsset, err := getAssetFromJSON(scanner.Bytes())
		if err != nil {
			return nil, err
		}
		pbAssets = append(pbAssets, pbAsset)
	}
	return pbAssets, nil
}

func getDataFromBucket(bucketName string) ([]*validator.Asset, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)
	var pbAssets []*validator.Asset
	for _, objectName := range destinationObjectNames {
		reader, err := bucket.Object(objectName).NewReader(ctx)
		if err != nil {
			fmt.Println("WARNING: Unable to read inventory file :", objectName, err)
			continue
		}
		defer reader.Close()
		assets, err := getDataFromReader(reader)
		if err != nil {
			return nil, err
		}

		pbAssets = append(pbAssets, assets...)
	}
	if len(pbAssets) == 0 {
		return nil, fmt.Errorf("No inventory found")
	}
	return pbAssets, nil
}

func getDataFromFile(caiDirName string) ([]*validator.Asset, error) {
	var pbAssets []*validator.Asset
	for _, objectName := range destinationObjectNames {
		reader, err := os.Open(filepath.Join(caiDirName, objectName))
		if err != nil {
			fmt.Println("WARNING: Unable to read inventory file :", objectName, err)
			continue
		}
		defer reader.Close()
		assets, err := getDataFromReader(reader)
		if err != nil {
			return nil, err
		}

		pbAssets = append(pbAssets, assets...)
	}
	if len(pbAssets) == 0 {
		return nil, fmt.Errorf("No inventory found")
	}
	return pbAssets, nil
}

func getDataFromStdin() ([]*validator.Asset, error) {
	return getDataFromReader(os.Stdin)
}

// getViolations finds all Config Validator violations for a given Inventory
func getViolations(inventory *InventoryConfig, config *ScoringConfig) ([]*RichViolation, error) {
	var err error
	var pbAssets []*validator.Asset

	if inventory.bucketName != "" {
		pbAssets, err = getDataFromBucket(inventory.bucketName)
		if err != nil {
			return nil, errors.Wrap(err, "Fetching inventory from Bucket")
		}
	} else if inventory.dirPath != "" {
		pbAssets, err = getDataFromFile(inventory.dirPath)
		if err != nil {
			return nil, errors.Wrap(err, "Fetching inventory from local directory")
		}
	} else if inventory.readFromStdin {
		pbAssets, err = getDataFromStdin()
		if err != nil {
			return nil, errors.Wrap(err, "Reading from stdin")
		}
	}

	richViolations := make([]*RichViolation, 0)
	wp := workerpool.New(inventory.workers)
	var badAsset *validator.Asset
	var mu sync.Mutex
	for _, asset := range pbAssets {
		asset := asset
		wp.Submit(func() {
			violations, errAsset := config.validator.ReviewAsset(context.Background(), asset)
			if errAsset != nil {
				err = errAsset
				badAsset = asset
				wp.Stop()
			}
			for _, violation := range violations {
				richViolation := RichViolation{violation, "", violation.Resource, violation.Message, violation.Metadata, asset}
				mu.Lock()
				richViolations = append(richViolations, &richViolation)
				mu.Unlock()
			}
		})
	}
	wp.StopWait()

	if err != nil {
		return nil, errors.Wrapf(err, "reviewing asset %s", badAsset)
	} else {
		return richViolations, nil
	}
}

// converts raw JSON into Asset proto
func getAssetFromJSON(input []byte) (*validator.Asset, error) {
	pbAsset := &validator.Asset{}
	err := unMarshallAsset(input, pbAsset)
	if err != nil {
		return nil, errors.Wrapf(err, "converting asset %s to proto", string(input))
	}

	err = cvasset.SanitizeAncestryPath(pbAsset)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching ancestry path for %s", pbAsset)
	}

	Log.Debug("Asset converted", "name", pbAsset.GetName(), "ancestry", pbAsset.GetAncestryPath())
	return pbAsset, nil
}
