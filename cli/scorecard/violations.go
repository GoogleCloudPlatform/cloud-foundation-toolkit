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
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"cloud.google.com/go/storage"

	tfconverter "github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
)

// attachValidator attaches a Validator to the given config
func attachValidator(config *ScoringConfig) error {
	v, err := gcv.NewValidator(
		gcv.PolicyPath(filepath.Join(config.PolicyPath, "policies")),
		gcv.PolicyLibraryDir(filepath.Join(config.PolicyPath, "lib")),
	)
	config.validator = v
	return err
}

func addDataFromScanner(config *ScoringConfig, scanner *bufio.Scanner) error {

	for scanner.Scan() {
		pbAsset, err := getAssetFromJSON(scanner.Bytes())

		pbAssets := []*validator.Asset{pbAsset}

		err = config.validator.AddData(&validator.AddDataRequest{
			Assets: pbAssets,
		})
		if err != nil {
			return errors.Wrap(err, "adding data to validator")
		}
	}
	return nil
}

// getViolations finds all Config Validator violations for a given Inventory
func getViolations(inventory *Inventory, config *ScoringConfig) (*validator.AuditResponse, error) {
	v := config.validator

	if inventory.GcsBucket != "" {
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return nil, err
		}
		bucket := client.Bucket(inventory.GcsBucket)

		for _, objectName := range destinationObjectNames {
			reader, err := bucket.Object(objectName).NewReader(ctx)
			if err != nil {
				return nil, err
			}
			defer reader.Close()
			scanner := bufio.NewScanner(reader)
			err = addDataFromScanner(config, scanner)
			if err != nil {
				return nil, errors.Wrap(err, "Fetching inventory")
			}
		}
	} else {
		for _, objectName := range destinationObjectNames {
			reader, err := os.Open(filepath.Join(inventory.DirName, objectName))
			if err != nil {
				return nil, err
			}
			defer reader.Close()
			scanner := bufio.NewScanner(reader)
			err = addDataFromScanner(config, scanner)
			if err != nil {
				return nil, errors.Wrap(err, "Fetching inventory")
			}
		}
	}

	auditResponse, err := v.Audit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "auditing")
	}

	return auditResponse, nil
}

// converts raw JSON into Asset proto
func getAssetFromJSON(input []byte) (*validator.Asset, error) {
	asset := tfconverter.Asset{}
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
	// TODO(morgantep): make this fetch the actual asset path
	// fmt.Printf("Asset parent: %v\n", pbAsset.GetResource().GetParent())
	return "organization/0/project/test", nil
}
