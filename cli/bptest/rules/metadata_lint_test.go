package bptest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sample metadata for testing
var sampleMetadata = `
spec:
  interfaces:
    - variables:
      - connections:
        - source: "blueprint://valid-source"
          version: "v1.0.0"
        - source: "blueprint://invalid-source"
          version: "invalid-version"
`

// Mock for the HTTP Get request
func mockHTTPGet(url string) (*http.Response, error) {
	log.Printf("zheng: called the mock http get function.")

	// Simulate different responses based on URL
	if regexp.MustCompile(`https://github.com/valid/source`).MatchString(url) {
		// Mocking a valid source
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(nil),
		}, nil
	} else if regexp.MustCompile(`https://github.com/valid/source/release/tag/1.0.0`).MatchString(url) {
		// Mocking an invalid source
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(nil),
		}, nil
	} else if regexp.MustCompile(`https://github.com/invalid/source`).MatchString(url) {
		// Mocking an invalid source
		return &http.Response{
			StatusCode: 404,
			Body:       ioutil.NopCloser(nil),
		}, nil
	}

	// Default case
	return nil, fmt.Errorf("URL not found: %s", url)
}

// TestValidateBlueprintSource tests the blueprint source validation
func TestValidateBlueprintSource(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/valid/path" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	httpGetFunc = func(url string) (*http.Response, error) {
		return http.Get(mockServer.URL + url[len("https://mockserver"):]) // Adjust path to test server
	}

	// Mock a valid source
	err := validateBlueprintSource("valid/source")
	assert.NoError(t, err, "Expected valid source to pass")

	//// Mock an invalid source
	//err = validateBlueprintSource("invalid/source")
	//assert.Error(t, err, "Expected invalid source to fail")
}

// TestValidateBlueprintVersion tests the blueprint version validation
func TestValidateBlueprintVersion(t *testing.T) {
	// Valid version
	err := validateBlueprintVersion("=1.0.0", "valid/source")
	assert.NoError(t, err, "Expected valid version to pass")

	// Invalid version
	err = validateBlueprintVersion("invalid-version", "valid/source")
	assert.Error(t, err, "Expected invalid version to fail")
}

// TestValidateMetadataFile tests the overall metadata file validation logic
func TestValidateMetadataFile(t *testing.T) {
	// Create a temporary directory for test metadata files
	dir := t.TempDir()
	metadataFilePath := filepath.Join(dir, "metadata.yaml")

	// Write the sample metadata to the file
	err := ioutil.WriteFile(metadataFilePath, []byte(sampleMetadata), 0644)
	assert.NoError(t, err, "Failed to create metadata file")

	// Mock HTTP GET function
	httpGetFunc = mockHTTPGet

	// Validate the metadata file
	err = validateMetadataFile(metadataFilePath)
	assert.Error(t, err, "Expected metadata validation to fail due to invalid source and version")
}

// TestRunMetadataLintCommand tests the main linting command function
func TestRunMetadataLintCommand(t *testing.T) {
	// Create a temporary directory for test metadata files
	dir := t.TempDir()
	metadataFilePath := filepath.Join(dir, "metadata.yaml")

	// Write the sample metadata to the file
	err := ioutil.WriteFile(metadataFilePath, []byte(sampleMetadata), 0644)
	assert.NoError(t, err, "Failed to create metadata file")

	// Mock HTTP GET function
	httpGetFunc = mockHTTPGet

	// Run the lint command (simulate running from the CLI)
	err = lintCmd.Execute()
	assert.NoError(t, err, "Expected runMetadataLintCommand to pass for valid metadata")
}
