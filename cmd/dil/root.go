package main

import (
	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/csv"
	"github.com/senseyeio/diligent/stdout"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var (
	csvFilePath string
)

var RootCmd = &cobra.Command{
	Use:   "dil [filePath]",
	Short: "Get the licenses associated with your software dependencies",
	Long:  `Diligent is a CLI tool which determines the licenses associated with your software dependencies`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		deper, err := getDeper(filePath, fileBytes)
		if err != nil {
			log.Fatal(err.Error())
		}

		runDep(deper, getReporter(), filePath)
	},
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().StringVar(&csvFilePath, "csv", "", "Writes CSV to the provided file path")
}

func getReporter() diligent.Reporter {
	if csvFilePath != "" {
		return csv.NewReporter(csvFilePath)
	}
	return stdout.NewReporter()
}
