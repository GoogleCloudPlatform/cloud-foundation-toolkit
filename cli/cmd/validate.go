package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/deployment"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/validation"
)

func init() {
	initProjectFlag(validateCmd)
	initPolicyPathFlag(validateCmd)
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate deployment(s)",
	Long:  `Validate deployment(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		setDefaultProjectID()

		if len(args) < 1 {
			log.Fatalf("At least one deployment name is expected")
		}

		for _, name := range args {
			_, err := validation.ValidateDeployment(name, policyPathFlag, deployment.DefaultProjectID)
			if err != nil {
				log.Fatalf("Validation error")
			}
		}
	},
}
