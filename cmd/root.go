package cmd

import (
	"bufio"
	"cli-search-engine/engine"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var HtmlEngine *engine.HtmlEngine
var filePath string

const (
	startUpMessage = "start up message"
)

var rootCmd = &cobra.Command{
	Use: "cli-search-engine",
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "search for a patterns in from a provided json file0",
	Long:  "load json file containing the data of html documents",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CompError("please provide a json file")
			return
		}
		filePath = args[0]
		HtmlEngine = engine.NewHtmlEngine(filePath)

		fmt.Println(startUpMessage)
		for {
			fmt.Println("type in search keywords")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("error reading input: %v", err)
			}
			searchTerms, flags := handleUserInput(input)

			result := HtmlEngine.Search(searchTerms, flags)
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&filePath, "file", "", "path to the JSON file for HtmlEngine")
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(pMatchCmd)
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func handleUserInput(input string) ([]string, []string) {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")

	var searchTerms []string
	var flags []string

	for _, part := range parts {
		if strings.HasPrefix(part, "--") {
			flags = append(flags, strings.Replace(part, "--", "", -1))
		} else {
			searchTerms = append(searchTerms, part)
		}
	}
	//TODO STEM SEARCH TERMS
	return searchTerms, flags
}
