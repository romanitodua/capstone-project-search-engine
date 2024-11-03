package engine

import (
	"cli-search-engine/engineLogger"
	"cli-search-engine/models"
	"cli-search-engine/strategies"
	"encoding/json"
	"io"
	"log"
	"math"
	"os"
)

type HtmlEngine struct {
	Documents *models.Documents
}

func NewHtmlEngine(jsonFile string) *HtmlEngine {
	documents, err := fromJson(jsonFile)
	if err != nil {
		log.Panicf("err could not parse json file: %v", err)
	}
	return &HtmlEngine{documents}
}

func (e *HtmlEngine) Search(terms []string, strategy string) string {
	// calculate tf-idf for searched terms
	res := calculateITF(e.Documents, terms)

	switch strategy {
	case BitonicSort:
		logger := engineLogger.NewBitonicSortLogger()
		bitonicSortStrategy := strategies.NewBitonicSort(res, 1, logger)
		return bitonicSortStrategy.Sort()
	}
	return ""
}

func calculateITF(docs *models.Documents, terms []string) []*models.TFIDF {
	tfidfs := []*models.TFIDF{}
	n := len(*docs)

	for _, term := range terms {
		termInCorpus := 0
		for _, doc := range *docs {
			for _, v := range doc {
				if _, ok := v[term]; ok {
					termInCorpus++
					break
				}
			}
		}

		idf := math.Log10(float64(n) / float64(termInCorpus))

		for _, doc := range *docs {
			for fileName, termFreqs := range doc {
				sum := 0
				tf := 0.0
				if freq, ok := termFreqs[term]; ok {
					for _, f := range termFreqs {
						sum += f
					}
					tf = float64(freq) / float64(sum)

					tfidf := &models.TFIDF{
						Document: fileName,
						Tf:       tf,
						Idf:      idf,
						Tfidf:    tf * idf,
					}

					tfidfs = append(tfidfs, tfidf)
				}
			}
		}
	}

	// filter for duplicates and sum the tf-idfs
	uniqueEntries := make(map[string]*models.TFIDF)
	for _, tfidf := range tfidfs {
		if _, ok := uniqueEntries[tfidf.Document]; !ok {
			uniqueEntries[tfidf.Document] = tfidf
			tfidf.Terms = terms
		} else {
			v := uniqueEntries[tfidf.Document]
			v.Tfidf += tfidf.Tfidf
		}
	}

	res := make([]*models.TFIDF, len(uniqueEntries))
	index := 0
	for _, tfidf := range uniqueEntries {
		res[index] = tfidf
		index++
	}
	return res
}

func fromJson(jsonFile string) (*models.Documents, error) {
	docs := models.Documents{}
	file, err := os.OpenFile(jsonFile, os.O_RDONLY, 0622)
	if err != nil {
		return nil, err
	}
	reader, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(reader, &docs)
	if err != nil {
		return nil, err
	}
	return &docs, nil
}
