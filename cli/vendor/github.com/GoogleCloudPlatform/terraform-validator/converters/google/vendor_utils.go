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

package google

import (
	"fmt"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
)

// NOTE: These functions were pulled from github.com/terraform-providers/terraform-provider-google. They can go away when the functionality they are providing is implemented in the future github.com/GoogleCloudPlatform/terraform-converters package.

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d converter.TerraformResourceData, config *converter.Config) (string, error) {
	return getProjectFromSchema("project", d, config)
}

func getProjectFromSchema(projectSchemaField string, d converter.TerraformResourceData, config *converter.Config) (string, error) {
	res, ok := d.GetOk(projectSchemaField)
	if ok && projectSchemaField != "" {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%s: required field is not set", projectSchemaField)
}
