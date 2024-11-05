package engineLogger

import (
	"encoding/json"
	"os"
	"time"
)

type QuickSortLogger struct {
	log *QuickSortLog
}

type QuickSortLog struct {
	StartMessage string                `json:"startMessage,omitempty"`
	StartedAt    time.Time             `json:"startedAt"`
	EndedAt      time.Time             `json:"endedAt"`
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
	l.log.StartedAt = time.Now()
}

func (l *QuickSortLogger) End() {
	l.log.EndedAt = time.Now()
	l.log.Duration = int(l.log.EndedAt.Sub(l.log.StartedAt).Seconds())
}

func (l *QuickSortLogger) AddIteration(info string, currentElements string) {
	return
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

func (l *QuickSortLogger) AddThread() {
	return
}

func (l *QuickSortLogger) ReleaseThread() {
	return
}

func (l *QuickSortLogger) SetResult(result string) {
	l.log.Result = result
}
