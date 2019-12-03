package policies

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"github.com/forseti-security/config-validator/pkg/bundlemanager"
)

func add(cmd *cobra.Command, args []string) error {
	source := path.Join(flags.libraryPath, flags.sourcePath)

	bundleManager := bundlemanager.New()
	if err := bundleManager.Load(source); err != nil {
		return err
	}

	for _, name := range args {
		Log.Info("Attempting to add object", "name", name)

		info, err := bundleManager.Inspect(name)
		if err != nil {
			return err
		}

		Log.Info("Found object", "name", name)

		d, err := info.GetYaml()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n\n", string(d))
	}

	return nil
}
