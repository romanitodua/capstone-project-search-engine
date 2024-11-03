package models

import (
	"fmt"
	"math"
	"strings"
)

type DocumentTF = map[string]map[string]int

type Documents = []DocumentTF

type TFIDF struct {
	Terms    []string
	Document string
	Tf       float64
	Idf      float64
	Tfidf    float64
}

func (t *TFIDF) Less(other *TFIDF) bool {
	return t.Tfidf < other.Tfidf
}

func GetTFIDFElements(elements []*TFIDF) string {
	var builder strings.Builder
	for i, elem := range elements {
		builder.WriteString(fmt.Sprintf("%s(%.4f)", elem.Document, elem.Tfidf))
		if i < len(elements)-1 {
			builder.WriteString(", ")
		}
	}
	return builder.String()
}

func GetMaxTFIDF() *TFIDF {
	return &TFIDF{
		Tfidf: math.Inf(1),
	}
}

func GetMinTFIDF() *TFIDF {
	return &TFIDF{
		Tfidf: math.Inf(-1),
	}
}
