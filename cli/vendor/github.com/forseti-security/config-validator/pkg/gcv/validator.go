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
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv/cf"
	"github.com/forseti-security/config-validator/pkg/gcv/configs"
	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const logRequestsVerboseLevel = 2

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
	// policyPath points to a directory where the constraints and constraint templates are stored as yaml files.
	policyPath string
	// policy dependencies directory points to rego files that provide supporting code for templates.
	// These rego dependencies should be packaged with the GCV deployment.
	// Right now expected to be set to point to "//policies/validator/lib" folder
	policyLibraryDir    string
	constraintFramework *cf.ConstraintFramework
}

// Option is a function for configuring Validator.
// See https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis for background.
type Option func(*Validator) error

// PolicyPath returns an Option that sets the root directory of constraints and constraint templates.
func PolicyPath(p string) Option {
	return func(v *Validator) error {
		v.policyPath = p
		return nil
	}
}

// PolicyLibraryDir returns an Option that sets the policy library directory with rego files.
// This function is expected to be removed in the future when all assumed dependant rego code is inlined in template files,
// and this validator includes the audit.rego files
func PolicyLibraryDir(dir string) Option {
	return func(v *Validator) error {
		v.policyLibraryDir = dir
		return nil
	}
}

func loadRegoFiles(dir string) (map[string]string, error) {
	loadedFiles := make(map[string]string)
	files, err := configs.ListRegoFiles(dir)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	for _, filePath := range files {
		glog.V(logRequestsVerboseLevel).Infof("Loading rego file: %s", filePath)
		if _, exists := loadedFiles[filePath]; exists {
			// This shouldn't happen
			return nil, status.Errorf(codes.Internal, "Unexpected file collision with file %s", filePath)
		}
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, errors.Wrapf(err, "unable to read file %s", filePath).Error())
		}
		loadedFiles[filePath] = string(fileBytes)
	}
	return loadedFiles, nil
}

func loadYAMLFiles(dir string) ([]*configs.ConstraintTemplate, []*configs.Constraint, error) {
	var templates []*configs.ConstraintTemplate
	var constraints []*configs.Constraint
	files, err := configs.ListYAMLFiles(dir)
	if err != nil {
		return nil, nil, err
	}
	for _, filePath := range files {
		glog.V(logRequestsVerboseLevel).Infof("Loading yaml file: %s", filePath)
		fileContents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, nil, status.Error(codes.InvalidArgument, errors.Wrapf(err, "unable to read file %s", filePath).Error())
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
			return nil, nil, status.Errorf(codes.Internal, "CategorizeYAMLFile returned unexpected data type when converting file %s", filePath)
		}
	}
	return templates, constraints, nil
}

// NewValidator returns a new Validator.
// By default it will initialize the underlying query evaluation engine by loading supporting library, constraints, and constraint templates.
// We may want to make this initialization behavior configurable in the future.
func NewValidator(options ...Option) (*Validator, error) {
	ret := &Validator{}
	for _, option := range options {
		if err := option(ret); err != nil {
			return nil, err
		}
	}
	if ret.policyPath == "" {
		return nil, status.Errorf(codes.InvalidArgument, "No policy path set, provide an option to set the policy path gcv.PolicyPath")
	}
	if ret.policyLibraryDir == "" {
		return nil, status.Errorf(codes.InvalidArgument, "No policy library set")
	}

	glog.V(logRequestsVerboseLevel).Infof("loading policy library dir: %s", ret.policyLibraryDir)
	regoLib, err := loadRegoFiles(ret.policyLibraryDir)
	if err != nil {
		return nil, err
	}

	ret.constraintFramework, err = cf.New(regoLib)
	if err != nil {
		return nil, err
	}
	glog.V(logRequestsVerboseLevel).Infof("loading policy dir: %s", ret.policyPath)
	templates, constraints, err := loadYAMLFiles(ret.policyPath)
	if err != nil {
		return nil, err
	}
	for _, template := range templates {
		if err := ret.constraintFramework.AddTemplate(template); err != nil {
			return nil, err
		}
	}
	for _, constraint := range constraints {
		if err := ret.constraintFramework.AddConstraint(constraint); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

// AddData adds GCP resource metadata to be audited later.
func (v *Validator) AddData(request *validator.AddDataRequest) error {

	for i, asset := range request.Assets {
		f, err := convertResourceViaJSONToInterface(asset)
		if err != nil {
			return status.Error(codes.Internal, errors.Wrapf(err, "index %d", i).Error())
		}
		v.constraintFramework.AddData(f)
	}

	return nil
}

func convertResourceViaJSONToInterface(asset *validator.Asset) (interface{}, error) {
	if asset == nil {
		return nil, nil
	}
	m := &jsonpb.Marshaler{
		OrigName: true,
	}
	if asset.Resource != nil {
		cleanStructValue(asset.Resource.Data)
	}
	glog.V(logRequestsVerboseLevel).Infof("converting asset to golang interface: %v", asset)
	var buf bytes.Buffer
	if err := m.Marshal(&buf, asset); err != nil {
		return nil, errors.Wrapf(err, "marshalling to json with asset %s: %v", asset.Name, asset)
	}
	var f interface{}
	err := json.Unmarshal(buf.Bytes(), &f)
	if err != nil {
		return nil, errors.Wrapf(err, "marshalling from json with asset %s: %v", asset.Name, asset)
	}
	return f, nil
}

// Reset clears previously added data from the underlying query evaluation engine.
func (v *Validator) Reset() error {
	v.constraintFramework.Reset()
	return nil
}

// Audit checks the GCP resource metadata that has been added via AddData to determine if any of the constraint is violated.
func (v *Validator) Audit(ctx context.Context) (*validator.AuditResponse, error) {
	response, err := v.constraintFramework.Audit(ctx)
	return response, err
}
