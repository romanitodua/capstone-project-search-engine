package cmd

import (
	"bufio"
	"cli-search-engine/strategies"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var pMatchCmd = &cobra.Command{
	Use:   "patternMatch",
	Short: "search for a pattern within a give json file",
	Long: "load json file containing the data of using the parsed command \n" +
		"Usage, : patternMatch pmatchDocs.json -> type in  any pattern",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CompError("please provide a json file")
			return
		}
		filePath = args[0]
		patternMatchStrategy, err := strategies.NewPatternMatchEngine(filePath)
		if err != nil {
			panic(fmt.Sprintf("error: failed to initialize pattern matching engine - %s", err.Error()))
		}

		fmt.Println(startUpMessage)
		for {
			fmt.Println("type in a pattern to search")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("error reading input: %v", err)
			}
			input = strings.TrimSpace(input)
			result := patternMatchStrategy.Search(input)
			fmt.Println(result)
		}
	},
}
