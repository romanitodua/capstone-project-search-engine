package engineLogger

import (
	"cli-search-engine/utils"
	"encoding/json"
	"os"
	"time"
)

type QuickSortLogger struct {
	log *QuickSortLog
}

type QuickSortLog struct {
	StartMessage string                `json:"startMessage,omitempty"`
	StartedAt    string                `json:"startedAt"`
	EndedAt      string                `json:"endedAt"`
	Started      time.Time             `json:"-"`
	Ended        time.Time             `json:"-"`
	Duration     int                   `json:"duration"` //seconds
	Result       string                `json:"result"`
	Iterations   []*QuickSortIteration `json:"recursiveCalls"`
}

type QuickSortIteration struct {
	QuickSortCallLow            int    `json:"quickSortCallLow"`
	QuickSortCallHigh           int    `json:"quickSortCallHigh"`
	PartitionFunctionCallNumber int    `json:"partitionFunctionCallNumber"`
	ElementsAfterPartition      string `json:"elementsAfterPartition"`
}

func NewQuickSortLogger() *QuickSortLogger {
	return &QuickSortLogger{
		log: &QuickSortLog{},
	}
}

func (l *QuickSortLogger) Log() error {
	bitonicSortLogFileName := "quickSortLog.json"
	logData, err := json.MarshalIndent(l.log, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(bitonicSortLogFileName, logData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (l *QuickSortLogger) SetStartMessage(msg string) {
	l.log.StartMessage = msg
}

func (l *QuickSortLogger) Start() {
	l.log.Started = time.Now()
	l.log.StartedAt = utils.FormatTime(l.log.Started)
}

func (l *QuickSortLogger) End() {
	l.log.Ended = time.Now()
	l.log.EndedAt = utils.FormatTime(l.log.Ended)
	l.log.Duration = int(l.log.Ended.Sub(l.log.Started).Seconds())
}

func (l *QuickSortLogger) AddQuickSortIteration(low, high int, elements string) {
	i := &QuickSortIteration{
		QuickSortCallLow:            low,
		QuickSortCallHigh:           high,
		PartitionFunctionCallNumber: len(l.log.Iterations) + 1,
		ElementsAfterPartition:      elements,
	}
	l.log.Iterations = append(l.log.Iterations, i)
}

func (l *QuickSortLogger) SetResult(result string) {
	l.log.Result = result
}
