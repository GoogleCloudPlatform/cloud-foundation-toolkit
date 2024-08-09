package rules

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	testdataDir    = "testdata"
	expectedSuffix = ".expected"
	updateEnvVar   = "UPDATE_EXPECTED"
	issueFile      = "issues.json"
)

var validExtensions = []string{".tf", ".hcl"}

// ruleTC is a single rule test case.
type ruleTC struct {
	// Dir with root module to be tested.
	dir string
}

// ruleTest tests rule r with test case rt.
func ruleTest(t *testing.T, r tflint.Rule, rt ruleTC) {
	t.Helper()

	config := configForTest(t, rt.dir)
	runner := helper.TestRunner(t, config)

	if err := r.Check(runner); err != nil {
		t.Fatalf("rule check: %s", err)
	}

	expected := path.Join(testdataDir, rt.dir, expectedSuffix, issueFile)
	updateExpected(t, expected, issuesToJSON(t, runner.Issues))
	wantIssues := issuesFromJSON(t, expected)
	helper.AssertIssues(t, wantIssues, runner.Issues)
}

// configForTest returns a map of TF configs paths to data stored in subdir.
// Paths are relative to subdir.
func configForTest(t *testing.T, subdir string) map[string]string {
	t.Helper()

	modDir := path.Join(testdataDir, subdir)
	configs := map[string]string{}
	err := filepath.WalkDir(modDir, func(fp string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// ignore hidden dirs
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !d.IsDir() && slices.Contains(validExtensions, path.Ext(fp)) {
			relPath, err := filepath.Rel(modDir, fp)
			if err != nil {
				return err
			}
			cfg, err := os.ReadFile(fp)
			if err != nil {
				return err
			}
			configs[relPath] = string(cfg)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("fetching testdata: %v", err)
	}
	return configs
}

// issuesFromJSON converts file at fp to a helper issues.
func issuesFromJSON(t *testing.T, fp string) helper.Issues {
	t.Helper()

	issues := helper.Issues{}
	data, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("reading issues: %v", err)
	}
	err = json.Unmarshal(data, &issues)
	if err != nil {
		t.Fatalf("unmarshalling issues: %v", err)
	}
	return issues
}

// issuesToJSON marshals issues to json bytes ignoring issue rule.
func issuesToJSON(t *testing.T, issues helper.Issues) []byte {
	t.Helper()

	// Workaround for unmarshal error of rule interface.
	for _, i := range issues {
		i.Rule = nil
	}
	data, err := json.MarshalIndent(issues, "", " ")
	if err != nil {
		t.Fatalf("marshalling issues: %v", err)
	}
	data = append(data, "\n"...)
	return data
}

// UpdateExpected updates expected file at fp with data with update env var is set.
func updateExpected(t *testing.T, fp string, data []byte) {
	t.Helper()

	if strings.ToLower(os.Getenv(updateEnvVar)) != "true" {
		return
	}
	err := os.MkdirAll(path.Dir(fp), os.ModePerm)
	if err != nil {
		t.Fatalf("updating result: %v", err)
	}

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		_, err := os.Create(fp)
		if err != nil {
			t.Fatalf("creating %s: %v", fp, err)
		}
	}

	err = os.WriteFile(fp, data, os.ModePerm)
	if err != nil {
		t.Fatalf("updating result: %v", err)
	}
}
