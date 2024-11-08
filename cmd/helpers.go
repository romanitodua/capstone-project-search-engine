package cmd

import (
	"cli-search-engine/stemmer"
	"strings"
)

func handleUserInput(input string) ([]string, []string) {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")

	var searchTerms []string
	var flags []string

	for _, part := range parts {
		if strings.HasPrefix(part, "--") {
			flags = append(flags, strings.Replace(part, "--", "", -1))
		} else {
			searchTerms = append(searchTerms, strings.ToLower(stemmer.Stem(part)))
		}
	}
	return searchTerms, flags
}
