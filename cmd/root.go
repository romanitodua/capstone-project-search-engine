package cmd

import (
	"cli-search-engine/engine"
	"github.com/spf13/cobra"
	"os"
)

var HtmlEngine *engine.HtmlEngine
var filePath string

const (
	startUpMessagePM     = "pattern matching engine initialized"
	startUpMessageSearch = "search engine initialized"
)

const (
	appName = "goseek"
)

var rootCmd = &cobra.Command{
	Use: appName,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&filePath, "file", "", "path to the JSON file for HtmlEngine")
	rootCmd.AddCommand(search)
	rootCmd.AddCommand(pMatchCmd)
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(crawlCmd)
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
