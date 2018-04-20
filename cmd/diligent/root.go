package main

import (
	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/csv"
	"github.com/senseyeio/diligent/stdout"
	"github.com/spf13/cobra"
)

var (
	csvFilePath      string
	licenseWhitelist []string
)

var RootCmd = &cobra.Command{
	Use:   "diligent [filePath]",
	Short: "Get the licenses associated with your software dependencies",
	Long:  `Diligent is a CLI tool which determines the licenses associated with your software dependencies`,
	Args:  cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		licenseWhitelist = diligent.ReplaceCategoriesWithIdentifiers(licenseWhitelist)
		if err := checkWhitelist(); err != nil {
			fatal(70, err.Error())
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		fileBytes := mustReadFile(filePath)
		deper, err := getDeper(filePath, fileBytes)
		if err != nil {
			fatal(69, err.Error())
		}

		runDep(deper, getReporter(), filePath)
	},
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().StringVar(&csvFilePath, "csv", "", "Writes CSV to the provided file path")
	RootCmd.PersistentFlags().StringSliceVarP(&licenseWhitelist, "whitelist", "w", nil, "Specify licenses compatible with your software. If licenses are found which are not in your whitelist, the command will return with a non zero exit code. Whitelisting license identifiers or categories of licenses is possible, the following categories are supported: 'permissive', 'copyleft', 'copyleft-limited', 'free-restricted', 'proprietary-free', 'public-domain'. See the readme for more details.")
}

func getReporter() diligent.Reporter {
	if csvFilePath != "" {
		return csv.NewReporter(csvFilePath)
	}
	return stdout.NewReporter()
}
