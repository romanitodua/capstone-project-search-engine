package engine

type PatternMatchingEngine struct {
	data map[string]string
}

func NewPatternMatchingEngine(data map[string]string) *PatternMatchingEngine {
	return &PatternMatchingEngine{data: data}
}

func (e *PatternMatchingEngine) Search(pattern string) map[string]int {
	result := map[string]int{}
	for k, v := range e.data {
		occurrences := e.patternMatch(v, pattern)
		result[k] = occurrences
	}
	return result
}

func (e *PatternMatchingEngine) ComputeFailureFunction(pattern string) []int {
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

func (e *PatternMatchingEngine) patternMatch(originalText, pattern string) int {
	f := e.ComputeFailureFunction(pattern)
	matches := 0
	i, j := 0, 0

	for i < len(originalText) {
		if originalText[i] == pattern[j] {
			i++
			j++
			if j == len(pattern) {
				matches++
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
