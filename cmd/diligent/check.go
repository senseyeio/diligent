package main

import (
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [filePath]",
	Short: "Ensures the licenses are compatible with your license whitelist",
	Long:  `Calling check will ensure that the licenses associated with your dependencies are compatible with your license whitelist.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(licenseWhitelist) == 0 {
			warning("your whitelist is empty, consider using the ls command instead")
		}
		run(args)
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
	applyCommonFlags(checkCmd)
	applyWhitelistFlag(checkCmd)
}
