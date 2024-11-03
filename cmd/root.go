package cmd

import (
	"cli-search-engine/engine"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var HtmlEngine *engine.HtmlEngine
var filePath string

var rootCmd = &cobra.Command{
	Use: "cli-search-engine",
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "creates a new html engine based on the json file",
	Long:  "load json file containing the data of html documents",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CompError("please provide a json file")
			return
		}
		filePath = args[0]
		HtmlEngine = engine.NewHtmlEngine(filePath)

		for {
			fmt.Println("enter")
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return
			}
			result := HtmlEngine.Search([]string{"roma"}, input)
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&filePath, "file", "", "path to the JSON file for HtmlEngine")
	rootCmd.AddCommand(loadCmd)
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
