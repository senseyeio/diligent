package main

import (
	"fmt"

	"github.com/senseyeio/diligent"
	"github.com/spf13/cobra"
)

// whitelistCmd represents the whitelist command
var whitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Details the licenses allowed by the current whitelist settings",
	Long: `Calling whitelist will detail what licenses are permitted with the provided flags. This can be used to
validate your are whitelisting just the licenses you are interested in`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		all := diligent.GetLicenses()
		for _, l := range all {
			if isInWhitelist(l) {
				fmt.Println(l.Identifier)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(whitelistCmd)
	applyWhitelistFlag(whitelistCmd)
}
