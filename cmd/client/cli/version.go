package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	BuildVersion string = "N/A"
	BuildDate    string = "N/A"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gophkeeper",
	Long:  `All software has versions. This is Gophkeeper's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nGophkeeper v%v, %v", BuildVersion, BuildDate)
	},
}
