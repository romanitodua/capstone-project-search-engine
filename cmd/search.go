package cmd

import (
	"bufio"
	"cli-search-engine/engine"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var search = &cobra.Command{
	Use:   "search",
	Short: "search by keywords,sentences in documents present in the json file. Accepted sort strategies quickSort,bitonicSort",
	Long: "load json file containing the contents parsed by parse command \n" +
		"Usage: load docs.json -> boston tea party --quickSort --quickSort is a flag \n" +
		"representing the sort strategy to with which to sort documents tf-idf values",
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
