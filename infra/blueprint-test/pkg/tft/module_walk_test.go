package tft

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func toSet(args ...string) stringSet {
	result := make(stringSet)
	for _, arg := range args {
		result[arg] = struct{}{}
	}
	return result
}

func TestFindModulesUnderTest(t *testing.T) {
	// Monkey patch findReferencedModules() in module_walk.go
	// to skip parsing real Terraform modules.
	findReferencedModules = func(tfDir string) (stringSet, error) {
		switch tfDir {

		// Group 1 (with external module references omitted):
		//   fixture1 -> [example1, example2]
		//   example1 -> [root, abc]
		//   example2 -> [xyz]
		case "local/test/integration/fixture1":
			return toSet("external1", "../../../examples/example1", "../../../examples/example2"), nil
		case "local/examples/example1":
			return toSet("../..", "../../modules/abc"), nil
		case "local/examples/example2":
			return toSet("../../modules/xyz", "external2"), nil

		// Group 2: fixture2 -> external_example, both point to external modules.
		case "local/test/integration/fixture2":
			return toSet("terraform-google-modules/network/google", "../../../examples/external_example"), nil
		case "local/examples/external_example":
			return toSet("terraform-google-modules/kubernetes-engine/google//modules/private-cluster"), nil

		default:
			return nil, fmt.Errorf("no fake behavior configured for module path %q", tfDir)
		}
	}

	tests := []struct {
		startPath string
		want      stringSet
	}{
		// Group 1.
		{
			startPath: "local/test/integration/fixture1",
			want:      toSet("root", "abc", "xyz"),
		},
		{
			startPath: "local/examples/example1",
			want:      toSet("root", "abc"),
		},
		{
			startPath: "local/examples/example2",
			want:      toSet("xyz"),
		},

		// Group 2: no modules under test found.
		{
			startPath: "local/test/integration/fixture2",
			want:      toSet(),
		},
		{
			startPath: "local/examples/external_example",
			want:      toSet(),
		},
	}
	for _, tt := range tests {
		t.Run(filepath.Base(tt.startPath), func(t *testing.T) {
			got, err := findModulesUnderTest(tt.startPath)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findModulesUnderTest(%q) returned unexpected result (-want +got):\n%s", tt.startPath, diff)
			}
		})
	}
}
