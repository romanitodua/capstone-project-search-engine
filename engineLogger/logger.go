package engineLogger

import (
	"cli-search-engine/models"
	"fmt"
)

type Logger interface {
	Log() error
	SetStartMessage(msg string)
	Start()
	End()
	AddIteration(info string, currentElements string)
	AddThread()
	ReleaseThread()
	SetResult(result string)
}

func NewLogger(strategy string) (Logger, error) {
	switch strategy {
	case models.BitonicSort:
		return NewBitonicSortLogger(), nil
	default:
		return nil, fmt.Errorf("unsupported strategy")
	}
}
