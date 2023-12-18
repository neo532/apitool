package main

import (
	"log"

	"github.com/neo532/apitool/internal/proto/add"
	"github.com/neo532/apitool/internal/proto/httpclient"
	"github.com/neo532/apitool/internal/proto/pbstruct"
	"github.com/neo532/apitool/internal/proto/service"
	"github.com/neo532/apitool/internal/upgrade"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "apitool",
		Short:   "Apitool: An elegant toolkit for api.",
		Long:    `Apitool: An elegant toolkit for api.`,
		Version: "v1.0.39",
	}
)

func init() {
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(add.CmdAdd)
	rootCmd.AddCommand(pbstruct.CmdStruct)

	rootCmd.AddCommand(httpclient.CmdClient)

	rootCmd.AddCommand(service.CmdService)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
