package bpcatalog

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var catalogListFlags struct {
	format renderFormat
	sort   sortOption
}

const (
	tfModulesOrg = "terraform-google-modules"
	gcpOrg       = "GoogleCloudPlatform"
)

var (
	// any repos that match terraform-google-* but should not be included
	repoIgnoreList = map[string]bool{
		"terraform-google-conversion": true,
		"terraform-google-examples":   true,
	}
	// any repos that do not match terraform-google-* but should be included
	repoAllowList = map[string]bool{
		"terraform-example-foundation": true,
	}
)

func init() {
	viper.AutomaticEnv()
	Cmd.AddCommand(listCmd)

	listCmd.Flags().Var(&catalogListFlags.format, "format", fmt.Sprintf("Format to display catalog. Defaults to table. Options are %+v.", renderFormats))
	listCmd.Flags().Var(&catalogListFlags.sort, "sort", fmt.Sprintf("Sort results. Defaults to created date. Options are %+v.", sortOptions))
}

var Cmd = &cobra.Command{
	Use:   "catalog",
	Short: "Blueprint catalog",
	Long:  `Blueprint catalog is used to get information about blueprints catalog.`,
	Args:  cobra.NoArgs,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists blueprints",
	Long:  `Lists blueprints in catalog`,
	Args:  cobra.NoArgs,
	RunE:  listCatalog,
}

func listCatalog(cmd *cobra.Command, args []string) error {
	// defaults
	if catalogListFlags.format.Empty() {
		catalogListFlags.format = renderTable
	}
	if catalogListFlags.sort.Empty() {
		catalogListFlags.sort = sortCreated
	}
	gh := newGHService(withTokenClient(), withOrgs([]string{tfModulesOrg, gcpOrg}))
	repos, err := fetchSortedTFRepos(gh, catalogListFlags.sort)
	if err != nil {
		return err
	}
	return render(repos, os.Stdout, catalogListFlags.format, viper.GetBool("verbose"))
}
