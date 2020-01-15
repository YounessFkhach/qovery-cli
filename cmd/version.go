package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information for the Qovery CLI",
	Long: `VERSION allows you to print version information for the qovery-cli. For example:

	qovery version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Qovery version 0.7.2")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
