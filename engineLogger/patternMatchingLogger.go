package engineLogger

import (
	"cli-search-engine/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

type PatternMatchingLogger struct {
	log PatternMatchLog
}

type PatternMatchLog struct {
	Input           string                   `json:"-"`
	StartMessage    string                   `json:"startMessage,omitempty"`
	StartedAt       string                   `json:"startedAt"`
	EndedAt         string                   `json:"endedAt"`
	Started         time.Time                `json:"-"`
	Ended           time.Time                `json:"-"`
	Duration        int                      `json:"duration"` //seconds
	Result          string                   `json:"result"`
	FailureFunction []int                    `json:"failureFunction"`
	Iteration       []*PatternMatchIteration `json:"iteration"`
	InputSize       int                      `json:"-"`
}

type PatternMatchIteration struct {
	Document                string `json:"document"`
	TotalPatternOccurrences int    `json:"totalPatternOccurrences"`
	PatternsDetectedIndexes []int  `json:"patternsDetectedIndexes"`
}

func NewPatternMatchingLogger() *PatternMatchingLogger {
	return &PatternMatchingLogger{}
}

func (l *PatternMatchingLogger) SetInputSize(size int) {
	l.log.InputSize = size
}

func (l *PatternMatchingLogger) Log(len int) error {
	patternMatchLogFileName := fmt.Sprintf("%d-PM-%s.json", l.log.InputSize, uuid.NewString())
	logData, err := json.MarshalIndent(l.log, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(patternMatchLogFileName, logData, 0644)
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

func (l *PatternMatchingLogger) SetFailureFunction(f []int) {
	l.log.FailureFunction = f
}

func (l *PatternMatchingLogger) AcquireIteration() *PatternMatchIteration {
	l.log.Iteration = append(l.log.Iteration, &PatternMatchIteration{})
	return l.log.Iteration[len(l.log.Iteration)-1]
}

func (l *PatternMatchingLogger) SetInput(input string) {
	l.log.Input = input
}
