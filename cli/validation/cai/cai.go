package cai

// Asset is the CAI representation of a resource.
type Asset struct {
	// The name, in a peculiar format: `\\<api>.googleapis.com/<self_link>`
	Name     string `json:"name"`
	Ancestry string `json:"ancestry_path"`
	// The type name in `google.<api>.<resourcename>` format.
	Type      string         `json:"asset_type"`
	Resource  *AssetResource `json:"resource,omitempty"`
	IAMPolicy *IAMPolicy     `json:"iam_policy,omitempty"`
}

// AssetResource is the Asset's Resource field.
type AssetResource struct {
	// Api version
	Version string `json:"version"`
	// URI including scheme for the discovery doc - assembled from
	// product name and version.
	DiscoveryDocumentURI string `json:"discovery_document_uri"`
	// Resource name.
	DiscoveryName string `json:"discovery_name"`
	// Actual resource state as per deployment.  Note that this does
	// not necessarily correspond perfectly with the CAI representation
	// as there are occasional deviations between CAI and API responses.
	// This returns the API response values instead.
	Data map[string]interface{} `json:"data,omitempty"`
}

type IAMPolicy struct {
	Bindings []IAMBinding `json:"bindings"`
}

type IAMBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}
