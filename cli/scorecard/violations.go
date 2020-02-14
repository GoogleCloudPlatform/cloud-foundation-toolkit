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
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/pkg/errors"
)

func getDataFromReader(config *ScoringConfig, reader io.Reader) ([]*validator.Asset, error) {
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

func getDataFromBucket(config *ScoringConfig, bucketName string) ([]*validator.Asset, error) {
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
			return nil, err
		}
		defer reader.Close()
		assets, err := getDataFromReader(config, reader)
		if err != nil {
			return nil, err
		}

		pbAssets = append(pbAssets, assets...)
	}
	return pbAssets, nil
}

func getDataFromFile(config *ScoringConfig, caiDirName string) ([]*validator.Asset, error) {
	var pbAssets []*validator.Asset
	for _, objectName := range destinationObjectNames {
		reader, err := os.Open(filepath.Join(caiDirName, objectName))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		assets, err := getDataFromReader(config, reader)
		if err != nil {
			return nil, err
		}

		pbAssets = append(pbAssets, assets...)
	}
	return pbAssets, nil
}

func getDataFromStdin(config *ScoringConfig) ([]*validator.Asset, error) {
	return getDataFromReader(config, os.Stdin)
}

// getViolations finds all Config Validator violations for a given Inventory
func getViolations(inventory *InventoryConfig, config *ScoringConfig) (*validator.AuditResponse, error) {
	var err error
	var pbAssets []*validator.Asset
	if inventory.bucketName != "" {
		pbAssets, err = getDataFromBucket(config, inventory.bucketName)
		if err != nil {
			return nil, errors.Wrap(err, "Fetching inventory from Bucket")
		}
	} else if inventory.dirPath != "" {
		pbAssets, err = getDataFromFile(config, inventory.dirPath)
		if err != nil {
			return nil, errors.Wrap(err, "Fetching inventory from local directory")
		}
	} else if inventory.readFromStdin {
		pbAssets, err = getDataFromStdin(config)
		if err != nil {
			return nil, errors.Wrap(err, "Reading from stdin")
		}
	}

	auditResult := &validator.AuditResponse{}
	for _, asset := range pbAssets {
		violations, err := config.validator.ReviewAsset(context.Background(), asset)
		if err != nil {
			return nil, errors.Wrapf(err, "reviewing asset %s", asset)
		}
		auditResult.Violations = append(auditResult.Violations, violations...)
	}
	return auditResult, nil
}

// converts raw JSON into Asset proto
func getAssetFromJSON(input []byte) (*validator.Asset, error) {
	var asset map[string]interface{}
	err := json.Unmarshal(input, &asset)
	if err != nil {
		return nil, err
	}
	pbAsset := &validator.Asset{}
	err = protoViaJSON(asset, pbAsset)
	if err != nil {
		return nil, errors.Wrapf(err, "converting asset %s to proto", asset.Name)
	}

	pbAsset.AncestryPath, err = getAncestryPath(pbAsset)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching ancestry path for %s", asset.Name)
	}
	Log.Debug("Asset converted", "name", asset.Name, "ancestry", pbAsset.GetAncestryPath())
	return pbAsset, nil
}

// looks up the ancestry path for a given asset
func getAncestryPath(pbAsset *validator.Asset) (string, error) {
	ancestors := pbAsset.Ancestors
	cnt := len(ancestors)
	revAncestors := make([]string, len(ancestors))
	for idx := 0; idx < cnt; idx++ {
		revAncestors[cnt-idx-1] = ancestors[idx]
	}
	ancestorPath := strings.Join(revAncestors, "/")
	if ancestorPath == ""{
		ancestorPath = "organization/0/project/test"
	}
	return ancestorPath, nil
}


// listFiles returns a list of files under a dir. Errors will be grpc errors.
func listFiles(dir string) ([]string, error) {
	files := []string{}
	visit := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "error visiting path %s", path)
		}
		if !f.IsDir() {
			files = append(files, path)
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		return nil, err
	}
	return files, nil
}
