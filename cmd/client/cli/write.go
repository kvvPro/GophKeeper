package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(writeCmd)
}

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "send info to server",
	Long:  "",
	Run:   writeToServer,
}

func writeToServer(cmd *cobra.Command, args []string) {
	fmt.Printf("\nGophkeeper v%v, %v", BuildVersion, BuildDate)
}
