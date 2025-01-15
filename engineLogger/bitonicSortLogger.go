package engineLogger

import (
	"cli-search-engine/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

type BitonicSortLog struct {
	Input                string                  `json:"input"`
	StartMessage         string                  `json:"startMessage,omitempty"`
	StartedAt            string                  `json:"startedAt"`
	EndedAt              string                  `json:"endedAt"`
	Started              time.Time               `json:"-"`
	Ended                time.Time               `json:"-"`
	Duration             int                     `json:"duration"` //seconds
	SpawnedThreads       int                     `json:"SpawnedThreads"`
	MaxConcurrentThreads int                     `json:"maxConcurrentThreads"`
	Result               string                  `json:"result"`
	Iterations           []*BitonicSortIteration `json:"iterations"`
	CurrentMax           int                     `json:"-"`
	InputSize            int                     `json:"-"`
}

type BitonicSortIteration struct {
	Info            string `json:"info"`
	CurrentElements string `json:"currentElements"`
}

type BitonicSortLogger struct {
	mu  sync.Mutex
	log *BitonicSortLog
}

func NewBitonicSortLogger() *BitonicSortLogger {
	return &BitonicSortLogger{
		log: &BitonicSortLog{},
	}
}

func (l *BitonicSortLogger) SetInputSize(size int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.InputSize = size
}

func (l *BitonicSortLogger) Log(len int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	bitonicSortLogFileName := fmt.Sprintf("%d-BS-%s.json", l.log.InputSize, uuid.NewString())
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

func (l *BitonicSortLogger) SetStartMessage(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.StartMessage = msg
}

func (l *BitonicSortLogger) Start() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Started = time.Now()
	l.log.StartedAt = utils.FormatTime(l.log.Started)
}

func (l *BitonicSortLogger) End() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Ended = time.Now()
	l.log.EndedAt = utils.FormatTime(l.log.Ended)
	l.log.Duration = int(l.log.Ended.Sub(l.log.Started).Seconds())
}

func (l *BitonicSortLogger) AddIteration(info string, currentElements string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Iterations = append(l.log.Iterations, &BitonicSortIteration{
		Info:            info,
		CurrentElements: currentElements,
	})
}

func (l *BitonicSortLogger) AddThread() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.SpawnedThreads++
	l.log.CurrentMax++
	if l.log.CurrentMax > l.log.MaxConcurrentThreads {
		l.log.MaxConcurrentThreads = l.log.CurrentMax
	}
}

func (l *BitonicSortLogger) ReleaseThread() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.log.CurrentMax > 0 {
		l.log.CurrentMax--
	}
}

func (l *BitonicSortLogger) SetResult(result string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Result = result
}

func (l *BitonicSortLogger) SetInput(input string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Input = input
}

func (l *BitonicSortLogger) AddQuickSortIteration(low, high int, elements string) {
	return
}
