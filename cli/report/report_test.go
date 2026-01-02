package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindReports_EmptyResults(t *testing.T) {
	tempDir := t.TempDir()

	// Create a dummy rego file that produces no results for data.reports
	regoContent := `package not_reports
# No rules defined
`
	regoFile := filepath.Join(tempDir, "empty.rego")
	require.NoError(t, os.WriteFile(regoFile, []byte(regoContent), 0644))

	// Create a dummy data file (can be empty json object)
	dataContent := `{}`
	dataFile := filepath.Join(tempDir, "data.json")
	require.NoError(t, os.WriteFile(dataFile, []byte(dataContent), 0644))

	// Call findReports
	results, err := findReports([]string{regoFile, dataFile})
	require.NoError(t, err)

	// Check if results is a map (expected empty map)
	require.IsType(t, map[string]interface{}{}, results)
	require.Empty(t, results)
}
