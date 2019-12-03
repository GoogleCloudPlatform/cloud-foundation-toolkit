package policies

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/forseti-security/config-validator/pkg/bundlemanager"
)

func list(cmd *cobra.Command, args []string) error {
	bundleManager := bundlemanager.New()
	if err := bundleManager.Load(flags.libraryPath); err != nil {
		return err
	}

	bundleKey := fmt.Sprintf("bundles.validator.forsetisecurity.org/%s", viper.GetString("bundle"))

	bundle := bundleManager.Bundle(bundleKey)

	Log.Info("Found bundle", "bundle", bundleKey, "objects", len(bundle.All()))

	if len(bundle.All()) < 1 {
		fmt.Println("No constraints found, have you initialized the library?")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Name", "Description", "Controls"})

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, obj := range bundle.All() {
		info := bundlemanager.GetInfo(obj)
		Log.Info("Printing object", "name", obj.GetName())
		controls := info.Controls()
		fmt.Println(controls)

		table.Append([]string{obj.GetName(), info.GetDescription(), info.GetControl(bundleKey)})
	}

	table.Render()

	return nil
}
