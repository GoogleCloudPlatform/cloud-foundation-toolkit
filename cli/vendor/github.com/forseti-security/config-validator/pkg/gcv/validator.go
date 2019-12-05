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

// Package gcv provides a library and a RPC service for Forseti Config Validator.
package gcv

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	asset2 "github.com/forseti-security/config-validator/pkg/asset"
	"github.com/forseti-security/config-validator/pkg/gcv/cf"
	"github.com/forseti-security/config-validator/pkg/gcv/configs"
	"github.com/forseti-security/config-validator/pkg/multierror"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	logRequestsVerboseLevel = 2
	// The JSON object key for ancestry path
	ancestryPathKey = "ancestry_path"
	// The JSON object key for ancestors list
	ancestorsKey = "ancestors"
)

var flags struct {
	workerCount int
}

func init() {
	flag.IntVar(
		&flags.workerCount,
		"workerCount",
		runtime.NumCPU(),
		"Number of workers that Validator will spawn to handle validate calls, this defaults to core count on the host")
}

// Validator checks GCP resource metadata for constraint violation.
//
// Expected usage pattern:
//   - call NewValidator to create a new Validator
//   - call AddData one or more times to add the GCP resource metadata to check
//   - call Audit to validate the GCP resource metadata that has been added so far
//   - call Reset to delete existing data
//   - call AddData to add a new set of GCP resource metadata to check
//   - call Reset to delete existing data
//
// Any data added in AddData stays in the underlying rule evaluation engine's memory.
// To avoid out of memory errors, callers can invoke Reset to delete existing data.
type Validator struct {
	// policyPaths points to a list of directories where the constraints and
	// constraint templates are stored as yaml files.
	policyPaths []string
	// policy dependencies directory points to rego files that provide supporting code for templates.
	// These rego dependencies should be packaged with the GCV deployment.
	// Right now expected to be set to point to "//policies/validator/lib" folder
	policyLibraryDir    string
	constraintFramework *cf.ConstraintFramework
	work                chan func()
}

func loadRegoFiles(dir string) (map[string]string, error) {
	loadedFiles := make(map[string]string)
	files, err := configs.ListRegoFiles(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list rego files from %s", dir)
	}
	for _, filePath := range files {
		glog.V(logRequestsVerboseLevel).Infof("Loading rego file: %s", filePath)
		if _, exists := loadedFiles[filePath]; exists {
			// This shouldn't happen
			return nil, errors.Errorf("unexpected file collision with file %s", filePath)
		}
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read file %s", filePath)
		}
		loadedFiles[filePath] = string(fileBytes)
	}
	return loadedFiles, nil
}

func loadYAMLFiles(dirs []string) ([]*configs.ConstraintTemplate, []*configs.Constraint, error) {
	var templates []*configs.ConstraintTemplate
	var constraints []*configs.Constraint
	var files []string
	for _, dir := range dirs {
		f, err := configs.ListYAMLFiles(dir)
		if err != nil {
			return nil, nil, err
		}
		files = append(files, f...)
	}
	for _, filePath := range files {
		glog.V(logRequestsVerboseLevel).Infof("Loading yaml file: %s", filePath)
		fileContents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read file %s", filePath)
		}
		categorizedData, err := configs.CategorizeYAMLFile(fileContents, filePath)
		if err != nil {
			glog.Infof("Unable to convert file %s, with error %v, assuming this file should be skipped and continuing", filePath, err)
			continue
		}
		switch data := categorizedData.(type) {
		case *configs.ConstraintTemplate:
			templates = append(templates, data)
		case *configs.Constraint:
			constraints = append(constraints, data)
		default:
			// Unexpected: CategorizeYAMLFile shouldn't return any types
			return nil, nil, errors.Errorf("CategorizeYAMLFile returned unexpected data type when converting file %s", filePath)
		}
	}
	return templates, constraints, nil
}

// NewValidator returns a new Validator.
// By default it will initialize the underlying query evaluation engine by loading supporting library, constraints, and constraint templates.
// We may want to make this initialization behavior configurable in the future.
func NewValidator(stopChannel <-chan struct{}, policyPaths []string, policyLibraryPath string) (*Validator, error) {
	if len(policyPaths) == 0 {
		return nil, errors.Errorf("No policy path set, provide an option to set the policy path gcv.PolicyPath")
	}
	if policyLibraryPath == "" {
		return nil, errors.Errorf("No policy library set")
	}

	ret := &Validator{
		work: make(chan func(), flags.workerCount*2),
	}

	glog.V(logRequestsVerboseLevel).Infof("loading policy library dir: %s", ret.policyLibraryDir)
	regoLib, err := loadRegoFiles(policyLibraryPath)
	if err != nil {
		return nil, err
	}

	ret.constraintFramework, err = cf.New(regoLib)
	if err != nil {
		return nil, err
	}
	glog.V(logRequestsVerboseLevel).Infof("loading policy dir: %v", ret.policyPaths)
	templates, constraints, err := loadYAMLFiles(policyPaths)
	if err != nil {
		return nil, err
	}

	if err := ret.constraintFramework.Configure(templates, constraints); err != nil {
		return nil, err
	}

	go func() {
		<-stopChannel
		glog.Infof("validator stopchannel closed, closing work channel")
		close(ret.work)
	}()

	workerCount := flags.workerCount
	glog.Infof("starting %d workers", workerCount)
	for i := 0; i < workerCount; i++ {
		go ret.reviewWorker(i)
	}

	return ret, nil
}

func (v *Validator) reviewWorker(idx int) {
	glog.Infof("worker %d starting", idx)
	for f := range v.work {
		f()
	}
	glog.Infof("worker %d terminated", idx)
}

// AddData adds GCP resource metadata to be audited later.
func (v *Validator) AddData(request *validator.AddDataRequest) error {
	for i, asset := range request.Assets {
		if err := asset2.ValidateAsset(asset); err != nil {
			return errors.Wrapf(err, "index %d", i)
		}
		f, err := asset2.ConvertResourceViaJSONToInterface(asset)
		if err != nil {
			return errors.Wrapf(err, "index %d", i)
		}
		v.constraintFramework.AddData(f)
	}

	return nil
}

type assetResult struct {
	violations []*validator.Violation
	err        error
}

func (v *Validator) handleReview(ctx context.Context, idx int, asset *validator.Asset, resultChan chan<- *assetResult) func() {
	return func() {
		resultChan <- func() *assetResult {
			if err := asset2.ValidateAsset(asset); err != nil {
				return &assetResult{err: errors.Wrapf(err, "index %d", idx)}
			}
			if asset.AncestryPath == "" && len(asset.Ancestors) != 0 {
				asset.AncestryPath = ancestryPath(asset.Ancestors)
			}

			assetInterface, err := asset2.ConvertResourceViaJSONToInterface(asset)
			if err != nil {
				return &assetResult{err: errors.Wrapf(err, "index %d", idx)}
			}

			violations, err := v.constraintFramework.Review(ctx, assetInterface)
			if err != nil {
				return &assetResult{err: errors.Wrapf(err, "index %d", idx)}
			}

			return &assetResult{violations: violations}
		}()
	}
}

// ancestryPath returns the ancestry path from a given ancestors list
func ancestryPath(ancestors []string) string {
	cnt := len(ancestors)
	revAncestors := make([]string, len(ancestors))
	for idx := 0; idx < cnt; idx++ {
		revAncestors[cnt-idx-1] = ancestors[idx]
	}
	return strings.Join(revAncestors, "/")
}

// fixAncestry will try to use the ancestors array to create the ancestorPath
// value if it is not present.
func (v *Validator) fixAncestry(input map[string]interface{}) error {
	if _, found := input[ancestryPathKey]; found {
		return nil
	}

	ancestorsIface, found := input[ancestorsKey]
	if !found {
		glog.Infof("asset missing ancestry information: %v", input)
		return nil
	}
	ancestorsIfaceSlice, ok := ancestorsIface.([]interface{})
	if !ok {
		return errors.Errorf("ancestors field not array type: %s", input)
	}
	if len(ancestorsIfaceSlice) == 0 {
		return nil
	}
	ancestors := make([]string, len(ancestorsIfaceSlice))
	for idx, v := range ancestorsIfaceSlice {
		val, ok := v.(string)
		if !ok {
			return errors.Errorf("ancestors field idx %d is not string %s, %s", idx, v, input)
		}
		ancestors[idx] = val
	}
	input[ancestryPathKey] = ancestryPath(ancestors)
	return nil
}

// ReviewJSON reviews the content of a JSON string
func (v *Validator) ReviewJSON(ctx context.Context, data string) ([]*validator.Violation, error) {
	asset := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data), &asset); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal json")
	}
	return v.ReviewUnmarshalledJSON(ctx, asset)
}

// ReviewJSON evaluates a single asset without any threading in the background.
func (v *Validator) ReviewUnmarshalledJSON(ctx context.Context, asset map[string]interface{}) ([]*validator.Violation, error) {
	if err := v.fixAncestry(asset); err != nil {
		return nil, err
	}
	return v.constraintFramework.Review(ctx, asset)
}

// Review evaluates each asset in the review request in parallel and returns any
// violations found.
func (v *Validator) Review(ctx context.Context, request *validator.ReviewRequest) (*validator.ReviewResponse, error) {
	assetCount := len(request.Assets)
	resultChan := make(chan *assetResult, flags.workerCount*2)
	defer close(resultChan)

	go func() {
		for idx, asset := range request.Assets {
			v.work <- v.handleReview(ctx, idx, asset, resultChan)
		}
	}()

	response := &validator.ReviewResponse{}
	var errs multierror.Errors
	for i := 0; i < assetCount; i++ {
		result := <-resultChan
		if result.err != nil {
			errs.Add(result.err)
			continue
		}
		response.Violations = append(response.Violations, result.violations...)
	}

	if !errs.Empty() {
		return response, errs.ToError()
	}
	return response, nil
}

// Reset clears previously added data from the underlying query evaluation engine.
func (v *Validator) Reset(ctx context.Context) error {
	return v.constraintFramework.Reset(ctx)
}

// Audit checks the GCP resource metadata that has been added via AddData to determine if any of the constraint is violated.
func (v *Validator) Audit(ctx context.Context) (*validator.AuditResponse, error) {
	response, err := v.constraintFramework.Audit(ctx)
	return response, err
}
