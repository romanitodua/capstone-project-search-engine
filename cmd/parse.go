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

var textRegex = regexp.MustCompile(`[^\w\s]+`) //REGEX FOR CHARACTERS THAT ARE OUTSIDE [A-Z,0-9]

const (
	htmlDataDirectoryPath     = "./html_data"
	htmlDataParsedDestination = "htmlDocs.json"
	pmDocsDestination         = "pmDocs.json"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse html files into json with stemming and tokenization of words & parse for patternMatching for later use by 'patternMatch command'",
	Long:  "parses json acquired html files from crawl command",
	Run: func(cmd *cobra.Command, args []string) {
		docs, pmDocs, err := parseHtml()
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
		pmdocsToJson, err := json.Marshal(pmDocs)
		if err != nil {
			log.Fatal(err)
		}
		pmDocsFile, err := os.Create(pmDocsDestination)
		if err != nil {
			log.Fatal(err)
		}
		defer pmDocsFile.Close()
		_, err = pmDocsFile.Write(pmdocsToJson)
		if err != nil {
			panic(err)
		}
	},
}

func parseHtml() ([]map[string]map[string]int, map[string]string, error) {
	docs := []map[string]map[string]int{}
	pmDocs := map[string]string{}
	err := filepath.Walk(htmlDataDirectoryPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dtf := make(map[string]map[string]int, 1)
			tf := make(map[string]int)
			file, err := os.OpenFile(path, os.O_RDONLY, 0622)
			if err != nil {
				log.Fatal(err)
			}
			filename := strings.ReplaceAll(strings.ReplaceAll(path, `html_data\`, ""), "_", " ")
			defer file.Close()
			doc, err := goquery.NewDocumentFromReader(file)
			if err != nil {
				return err
			}
			tokens, err := tokenizeFile(doc)
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
			dtf[filename] = tf
			docs = append(docs, dtf)

			text, err := extractTextPM(doc)
			if err != nil {
				return err
			}
			pmDocs[filename] = text
		}
		return nil
	})
	return docs, pmDocs, err
}

func tokenizeFile(doc *goquery.Document) ([]string, error) {

	var tokens []string

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		text = strings.ReplaceAll(text, "\n", " ")
		text = textRegex.ReplaceAllString(text, "")

		parts := strings.Fields(text)
		for _, part := range parts {
			if part != "" {
				tokens = append(tokens, stemmer.Stem(strings.ToLower(part)))
			}
		}
	})
	return tokens, nil
}

func extractTextPM(doc *goquery.Document) (string, error) {
	var finalText strings.Builder

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		text = strings.ReplaceAll(text, "\n", " ")
		text = textRegex.ReplaceAllString(text, "")

		finalText.WriteString(text + " ")
	})
	return finalText.String(), nil
}
