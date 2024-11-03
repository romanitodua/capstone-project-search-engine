package engineLogger

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type BitonicSortLog struct {
	StartMessage         string                 `json:"startMessage,omitempty"`
	StartedAt            time.Time              `json:"startedAt"`
	EndedAt              time.Time              `json:"endedAt"`
	Duration             int                    `json:"duration"` //seconds
	SpawnedThreads       int                    `json:"SpawnedThreads"`
	MaxConcurrentThreads int                    `json:"maxConcurrentThreads"`
	Result               string                 `json:"result"`
	Iterations           []BitonicSortIteration `json:"iterations"`
	CurrentMax           int                    `json:"-"`
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

func (l *BitonicSortLogger) Log() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	bitonicSortLogFileName := "bitonicSortLog.json"
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
	l.log.StartedAt = time.Now()
}

func (l *BitonicSortLogger) End() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.EndedAt = time.Now()
	l.log.Duration = int(l.log.EndedAt.Sub(l.log.StartedAt).Seconds())
}

func (l *BitonicSortLogger) AddIteration(info string, currentElements string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.Iterations = append(l.log.Iterations, BitonicSortIteration{
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
