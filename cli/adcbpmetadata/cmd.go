package adcbpmetadata

import (
	"github.com/spf13/cobra"
)



func init() {

}

var Cmd = &cobra.Command{
	Use:   "adc_validate",
	Short: "Validated terraform module for ingestion into ADC",
	Long:  `Validated terraform module for ingestion into ADC`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("ADC validate command in porgress...")
	},
}
