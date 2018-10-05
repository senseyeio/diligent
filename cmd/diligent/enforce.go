package main

import (
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var enforceCmd = &cobra.Command{
	Use:   "enforce [filePath]",
	Short: "Ensures the licenses are compatible with your license whitelist",
	Long:  `Calling enforce will check that the licenses associated with your dependencies are compatible with your license whitelist.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(licenseWhitelist) == 0 {
			warning("your whitelist is empty, consider using the ls command instead")
		}
		run(args)
	},
}

func init() {
	RootCmd.AddCommand(enforceCmd)
	applyCommonFlags(enforceCmd)
	applyWhitelistFlag(enforceCmd)
}
