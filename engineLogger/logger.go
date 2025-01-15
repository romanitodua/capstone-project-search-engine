package engineLogger

import (
	"cli-search-engine/models"
	"fmt"
)

type Logger interface {
	SetInput(input string)
	SetInputSize(size int)
	Log(len int) error
	SetStartMessage(msg string)
	Start()
	End()
	SetResult(result string)
}

func NewLogger(strategy string) (Logger, error) {
	switch strategy {
	case models.BitonicSort:
		return NewBitonicSortLogger(), nil
	case models.QuickSort:
		return NewQuickSortLogger(), nil
	case models.PatternMatch:
		return NewPatternMatchingLogger(), nil

	default:
		return nil, fmt.Errorf("unsupported strategy")
	}
}
