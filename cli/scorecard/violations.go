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
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	tfconverter "github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	"github.com/forseti-security/config-validator/pkg/multierror"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// attachValidator attaches a Validator to the given config
func attachValidator(stopCh chan struct{}, config *ScoringConfig) error {
	v, err := gcv.NewValidator(
		stopCh,
		filepath.Join(config.PolicyPath, "policies"),
		filepath.Join(config.PolicyPath, "lib"),
	)
	config.validator = v
	return err
}

func getReadersForBucket(bucketName string) ([]io.ReadCloser, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	var readers []io.ReadCloser
	bucket := client.Bucket(bucketName)
	for _, objectName := range destinationObjectNames {
		reader, err := bucket.Object(objectName).NewReader(ctx)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return readers, nil
}

func getReadersForFile(caiDirName string) ([]io.ReadCloser, error) {
	var readers []io.ReadCloser
	for _, objectName := range destinationObjectNames {
		path := filepath.Join(caiDirName, objectName)
		glog.Infof("creating reader for %s", path)
		reader, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return readers, nil
}

// getViolations finds all Config Validator violations for a given Inventory
func getViolations(inventory *InventoryConfig, config *ScoringConfig) (*validator.AuditResponse, error) {
	glog.Info("getting violations")
	var err error
	var readers []io.ReadCloser
	pipeline := NewPipeline(config.validator)
	if inventory.bucketName != "" {
		readers, err = getReadersForBucket(inventory.bucketName)
		if err != nil {
			return nil, errors.Wrap(err, "Fetching inventory from Bucket")
		}
	} else {
		readers, err = getReadersForFile(inventory.dirPath)
		if err != nil {
			return nil, errors.Wrap(err, "Fetching inventory from local directory")
		}
	}

	go func() {
		for _, r := range readers {
			pipeline.AddInput(r)
		}
		pipeline.CloseInput()
	}()

	var count int
	var errs multierror.Errors
	var violations []*validator.Violation
	for result := range pipeline.Results() {
		count++
		for _, err := range result.Errs {
			errs.Add(err)
		}
		violations = append(violations, result.Violations...)
	}
	auditResponse := &validator.AuditResponse{
		Violations: violations,
	}
	return auditResponse, errs.ToError()
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
