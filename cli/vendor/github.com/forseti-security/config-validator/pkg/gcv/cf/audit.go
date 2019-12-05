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

package cf

// AuditRego contains the Rego snippet that implements matching logic
// and applies constraints to assets.
const AuditRego = `
package validator.gcp.lib

# Audit endpoint to grab all violations found in data.constraints

audit[result] {
	inventory := data.inventory

	asset := inventory[_]
	handle_asset[result] with input.asset as asset 
}

handle_asset[result] {
	asset := input.asset
	constraints := data.constraints
	constraint := constraints[_]
	spec := _get_default(constraint, "spec", {})
	match := _get_default(spec, "match", {})
	# Default matcher behavior is to match everything.
	target := _get_default(match, "target", ["organization/*"])
	# TODO: Retire the "gcp" wrapper.
	# See https://github.com/forseti-security/config-validator/issues/42
	gcp := _get_default(match, "gcp", {})
	gcp_target := _get_default(gcp, "target", target)
	re_match(gcp_target[_], asset.ancestry_path)
	exclude := _get_default(match, "exclude", [])
	gcp_exclude := _get_default(gcp, "exclude", exclude)
	exclusion_match := {asset.ancestry_path | re_match(gcp_exclude[_], asset.ancestry_path)}
	count(exclusion_match) == 0

	violations := data.templates.gcp[constraint.kind].deny with input.asset as asset
		 with input.constraint as constraint

	violation := violations[_]

	result := {
		"asset": asset.name,
		"constraint": constraint.metadata.name,
		"constraint_config": constraint,
		"violation": violation,
	}
}

# The following functions are prefixed with underscores, because their names
# conflict with the existing functions in policy library. We want to separate
# them here to ensure that there is no dependency.

# has_field returns whether an object has a field
_has_field(object, field) {
	object[field]
}

# False is a tricky special case, as false responses would create an undefined document unless
# they are explicitly tested for
_has_field(object, field) {
	object[field] == false
}

_has_field(object, field) = false {
	not object[field]
	not object[field] == false
}

# get_default returns the value of an object's field or the provided default value.
# It avoids creating an undefined state when trying to access an object attribute that does
# not exist
_get_default(object, field, _default) = output {
	_has_field(object, field)
	output = object[field]
}

_get_default(object, field, _default) = output {
	_has_field(object, field) == false
	output = _default
}
`
