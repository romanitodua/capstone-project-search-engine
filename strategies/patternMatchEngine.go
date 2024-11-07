package strategies

import (
	"cli-search-engine/engineLogger"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type PatternMatchingStrategy struct {
	data   map[string]string
	logger *engineLogger.PatternMatchingLogger
}

func NewPatternMatchEngine(filePath string) (*PatternMatchingStrategy, error) {
	data := map[string]string{}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0622)
	if err != nil {
		return nil, err
	}
	reader, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(reader, &data)
	if err != nil {
		return nil, err
	}
	acquireLogger := engineLogger.NewPatternMatchingLogger()
	return &PatternMatchingStrategy{data: data, logger: acquireLogger}, nil
}

func (p *PatternMatchingStrategy) Search(pattern string) string {
	result := map[string]int{}
	p.logger.SetStartMessage(fmt.Sprintf(
		"pattern matching initialized with %d documents, search pattern - %s", len(p.data), pattern))
	p.logger.Start()
	for k, v := range p.data {
		failureFunction := p.ComputeFailureFunction(pattern)
		p.logger.SetFailureFunction(failureFunction)
		iteration := p.logger.AcquireIteration()
		iteration.Document = k
		occurrences := p.patternMatch(failureFunction, v, pattern, iteration)
		result[k] = occurrences
	}
	p.logger.End()
	var resultBuilder strings.Builder
	for k, v := range result {
		if v == 0 {
			continue
		}
		resultBuilder.WriteString(fmt.Sprintf("%s(%d)", k, v))
	}
	if resultBuilder.Len() == 0 {
		resultBuilder.WriteString("no matches found")
	}
	p.logger.SetResult(resultBuilder.String())
	err := p.logger.Log()
	if err != nil {
		fmt.Println("error: failed to log the results of pattern matching")
	}
	return resultBuilder.String()
}

func (p *PatternMatchingStrategy) ComputeFailureFunction(pattern string) []int {
	f := make([]int, len(pattern))
	f[0] = 0

	for j := 1; j < len(pattern); j++ {
		i := f[j-1]
		for i > 0 && pattern[j] != pattern[i] {
			i = f[i-1]
		}
		if pattern[j] == pattern[i] {
			f[j] = i + 1
		} else {
			f[j] = 0
		}
	}

	return f
}

func (p *PatternMatchingStrategy) patternMatch(f []int, originalText, pattern string, iteration *engineLogger.PatternMatchIteration) int {
	matches := 0
	i, j := 0, 0

	for i < len(originalText) {
		if originalText[i] == pattern[j] {
			i++
			j++
			if j == len(pattern) {
				matches++
				iteration.TotalPatternOccurrences++
				iteration.PatternsDetectedIndexes = append(iteration.PatternsDetectedIndexes, i-len(pattern))
				j = f[j-1]
			}
		} else {
			if j != 0 {
				j = f[j-1]
			} else {
				i++
			}
		}
	}
	return matches
}
