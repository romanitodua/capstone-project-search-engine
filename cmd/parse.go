package cmd

import (
	"cli-search-engine/stemmer"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const htmlDataDirectoryPath = "./html_data"
const htmlDataParsedDestination = "htmlDocs.json"

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse html files into json with stemming and tokenization of words & parse for patternMatching for later use by 'patternMatch command'",
	Long:  "parses json acquired html files from crawl command",
	Run: func(cmd *cobra.Command, args []string) {
		docs, err := parseHtml()
		if err != nil {
			log.Fatal(err)
		}
		docsToJson, err := json.Marshal(docs)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(htmlDataParsedDestination)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write(docsToJson)
		if err != nil {
			panic(err)
		}
		//TODO  also parse for pattern matching
	},
}

func parseHtml() ([]map[string]map[string]int, error) {
	docs := []map[string]map[string]int{}
	err := filepath.Walk(htmlDataDirectoryPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dtf := make(map[string]map[string]int, 1)
			tf := make(map[string]int)
			file, err := os.OpenFile(path, os.O_RDONLY, 0622)
			if err != nil {
				log.Fatal(err)
			}
			tokens, err := tokenizeFile(file)
			if err != nil {
				return err
			}
			for _, token := range tokens {
				if _, ok := tf[token]; !ok {
					tf[token] = 1
				} else {
					tf[token] += 1
				}
			}
			dtf[path] = tf
			docs = append(docs, dtf)
		}
		return nil
	})
	return docs, err
}

func tokenizeFile(file *os.File) ([]string, error) {
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	var tokens []string
	re := regexp.MustCompile(`[^\w\s]+`)

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		text = strings.ReplaceAll(text, "\n", " ")
		text = re.ReplaceAllString(text, "")

		parts := strings.Fields(text)
		for _, part := range parts {
			if part != "" {
				tokens = append(tokens, stemmer.Stem(strings.ToLower(part)))
			}
		}
	})
	return tokens, nil
}
