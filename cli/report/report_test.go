package report

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindReports_EmptyResults(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "report_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy rego file that produces no results for data.reports
	regoContent := `package not_reports
# No rules defined
`
	regoFile := filepath.Join(tempDir, "empty.rego")
	if err := os.WriteFile(regoFile, []byte(regoContent), 0644); err != nil {
		t.Fatalf("Failed to write rego file: %v", err)
	}

	// Create a dummy data file (can be empty json object)
	dataContent := `{}`
	dataFile := filepath.Join(tempDir, "data.json")
	if err := os.WriteFile(dataFile, []byte(dataContent), 0644); err != nil {
		t.Fatalf("Failed to write data file: %v", err)
	}

	// Call findReports
	results, err := findReports([]string{regoFile, dataFile})
	if err != nil {
		t.Fatalf("findReports returned error: %v", err)
	}

	// Check if results is a map (expected empty map)
	resultsMap, ok := results.(map[string]interface{})
	if !ok {
		t.Errorf("Expected results to be map[string]interface{}, got %T", results)
	}
	if len(resultsMap) != 0 {
		t.Errorf("Expected empty results map, got length %d", len(resultsMap))
	}
}
