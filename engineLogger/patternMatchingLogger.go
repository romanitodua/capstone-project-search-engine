package engineLogger

import (
	"cli-search-engine/utils"
	"encoding/json"
	"os"
	"time"
)

type PatternMatchingLogger struct {
	log PatternMatchLog
}

type PatternMatchLog struct {
	StartMessage string    `json:"startMessage,omitempty"`
	StartedAt    string    `json:"startedAt"`
	EndedAt      string    `json:"endedAt"`
	Started      time.Time `json:"-"`
	Ended        time.Time `json:"-"`
	Duration     int       `json:"duration"` //seconds
	Result       string    `json:"result"`
}

func NewPatternMatchingLogger() *PatternMatchingLogger {
	return &PatternMatchingLogger{}
}

func (l *PatternMatchingLogger) Log() error {
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

func (l *PatternMatchingLogger) SetStartMessage(msg string) {
	l.log.StartMessage = msg
}

func (l *PatternMatchingLogger) Start() {
	l.log.Started = time.Now()
	l.log.StartedAt = utils.FormatTime(l.log.Started)
}

func (l *PatternMatchingLogger) End() {
	l.log.Ended = time.Now()
	l.log.EndedAt = utils.FormatTime(l.log.Ended)
	l.log.Duration = int(l.log.Ended.Sub(l.log.Started).Seconds())
}

func (l *PatternMatchingLogger) SetResult(result string) {
	l.log.Result = result
}
