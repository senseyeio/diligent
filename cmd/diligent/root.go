package main

import (
	"github.com/senseyeio/diligent"
	"github.com/spf13/cobra"
)

var (
	csvFilePath      string
	licenseWhitelist []string
	npmDevDeps       bool
)

var RootCmd = &cobra.Command{
	Short: "Get the licenses associated with your software dependencies",
	Long:  `Diligent is a CLI tool which determines the licenses associated with your software dependencies`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		licenseWhitelist = diligent.ReplaceCategoriesWithIdentifiers(licenseWhitelist)
		if err := checkWhitelist(); err != nil {
			fatal(70, err.Error())
		}
	},
}

func init() {
	cobra.OnInitialize()
}

func applyCommonFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&npmDevDeps, "npm-dev-deps", "", false, "[NPM] Include developer dependencies")
	cmd.Flags().StringVar(&csvFilePath, "csv", "", "Writes CSV to the provided file path")
}

func applyWhitelistFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(&licenseWhitelist, "whitelist", "w", nil, "Specify licenses compatible with your software. If licenses are found which are not in your whitelist, the command will return with a non zero exit code. Whitelisting license identifiers or categories of licenses is possible, the following categories are supported: 'all', 'permissive', 'copyleft', 'copyleft-limited', 'free-restricted', 'proprietary-free', 'public-domain'. See the readme for more details.")
}
