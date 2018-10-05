package main

import (
	"github.com/senseyeio/diligent"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [path]",
	Short: "Lists the licenses associated with your dependencies",
	Long:  `Calling ls will list the licenses associated with your dependencies.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		licenseWhitelist = diligent.GetLicenseIdentifiers()
		run(args)
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
	applyCommonFlags(lsCmd)
}
