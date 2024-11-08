package cmd

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const wikipediaRandomArticleUrl = "https://en.wikipedia.org/wiki/Special:Random"
const saveDir = "html_data"

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "crawls the wikipedia website for random for a specified amount of documents and saves them as html for later parsing ",
	Long: "load json file containing the contents parsed by parse command \n" +
		"Usage: load docs.json -> boston tea party --quickSort --quickSort is a flag \n" +
		"representing the sort strategy to with which to sort documents tf-idf values",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CompError("please provide integer number for retrieving documents")
			return
		}
		n, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(fmt.Sprintf("error acquiring a number %v", err))
		}
		crawl(n)
		fmt.Println("crawling finished")
	},
}

func crawl(docsToDownload int) {
	wg := sync.WaitGroup{}
	documentsToDownload := 1
	for i := 0; i < documentsToDownload; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := rand.IntN(10)
			time.Sleep(time.Second * time.Duration(r))
			c := colly.NewCollector()
			c.OnHTML("html", func(e *colly.HTMLElement) {
				articleTitle := e.ChildText("#firstHeading")

				if articleTitle == "" {
					articleTitle = uuid.New().String()
				}
				fileName := strings.ReplaceAll(articleTitle, " ", "_") + ".html"

				err := ensureDirectoryExists(saveDir)
				if err != nil {
					log.Fatal(err)
					return
				}
				filePath := filepath.Join(saveDir, fileName)
				html, err := e.DOM.Html()
				if err != nil {
					log.Printf("can't save html: %s", err)
				}
				err = os.WriteFile(filePath, []byte(html), 0644)
				if err != nil {
					log.Printf("Error saving article '%s': %v\n", articleTitle, err)
				}
				log.Printf("Article '%s' saved as '%s'\n", articleTitle, fileName)
			})

			c.OnError(func(r *colly.Response, err error) {
				log.Printf("Request to %s failed: %v\n", r.Request.URL, err)
			})

			err := c.Visit(wikipediaRandomArticleUrl)
			if err != nil {
				log.Printf("Error visiting URL: %v\n", err)
			}
		}()
	}
	wg.Wait()
}

func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	}
	return nil
}
