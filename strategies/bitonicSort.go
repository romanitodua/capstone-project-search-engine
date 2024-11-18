package strategies

import (
	"cli-search-engine/engineLogger"
	"cli-search-engine/models"
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"
)

type BitonicSort struct {
	Elements  []*models.TFIDF
	Direction int
	Logger    *engineLogger.BitonicSortLogger
}

// NewBitonicSort descending direction - 0, ascending direction - 1
func NewBitonicSort(elements []*models.TFIDF, dir int, logger engineLogger.Logger) *BitonicSort {
	l := len(elements)
	var startMessage strings.Builder
	startMessage.WriteString(fmt.Sprintf("initialized with elements of length - %d", l))
	if !isPowerOfTwo(l) {
		n := nextPowerOfTwo(l) - l
		startMessage.WriteString(fmt.Sprintf("- is not a power of two,next power of two is - %d", nextPowerOfTwo(l)))
		for i := 0; i < n; i++ {
			if dir == 1 {
				startMessage.WriteString(" appending dummy elements with infinite tfidf")
				elements = append(elements, models.GetMaxTFIDF())
			} else {
				startMessage.WriteString(" appending dummy elements with -infinite tfidf")
				elements = append(elements, models.GetMinTFIDF())
			}
		}
	}
	acquireLogger := logger.(*engineLogger.BitonicSortLogger)
	acquireLogger.SetStartMessage(startMessage.String())
	return &BitonicSort{
		Elements:  elements,
		Direction: dir,
		Logger:    acquireLogger,
	}
}

func (b *BitonicSort) Sort() string {
	b.bitonicSort(b.Elements, b.Direction)
	b.Logger.End()

	if len(b.Elements) == 0 {
		return fmt.Sprintf("search results not found")
	}

	b.Elements = slices.DeleteFunc(b.Elements, func(e *models.TFIDF) bool {
		if b.Direction == 1 {
			return e.Tfidf == math.Inf(1)
		} else {
			return e.Tfidf == math.Inf(-1)
		}
	})
	b.Logger.SetResult(models.GetTFIDFElements(b.Elements))
	err := b.Logger.Log()
	if err != nil {
		fmt.Printf("search ended with an error while saving logs: %v\n", err)
	}

	if len(b.Elements) < 5 {
		return fmt.Sprintf("top matches are - %s", models.GetTFIDFElements(b.Elements))
	}
	return fmt.Sprintf("top matches are - %s,showing 5docs out of %d see more in the logs",
		models.GetTFIDFElements(b.Elements[:5]), len(b.Elements))
}

func (b *BitonicSort) bitonicSort(elements []*models.TFIDF, dir int) {
	l := len(elements)
	b.Logger.Start()
	if l > 1 {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer func() {
				defer wg.Done()
				b.Logger.ReleaseThread()
			}()
			b.Logger.AddThread()
			b.Logger.AddIteration(iterationInfo(elements[:l/2], false, 1))
			b.bitonicSort(elements[:l/2], 1)
		}()

		go func() {
			defer func() {
				defer wg.Done()
				b.Logger.ReleaseThread()
			}()
			b.Logger.AddThread()
			b.Logger.AddIteration(iterationInfo(elements[l/2:], false, 0))
			b.bitonicSort(elements[l/2:], 0)
		}()

		wg.Wait()
		b.bitonicMerge(elements, dir)
	}
}

func (b *BitonicSort) bitonicMerge(elements []*models.TFIDF, dir int) {
	l := len(elements)
	if l > 1 {
		mid := l / 2
		for i := 0; i < mid; i++ {
			compareAndSwap(elements, i, i+mid, dir)
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer func() {
				defer wg.Done()
				b.Logger.ReleaseThread()
			}()
			b.Logger.AddThread()
			b.Logger.AddIteration(iterationInfo(elements[:mid], true, dir))
			b.bitonicMerge(elements[:mid], dir)
		}()

		go func() {
			defer func() {
				defer wg.Done()
				b.Logger.ReleaseThread()
			}()
			b.Logger.AddThread()
			b.Logger.AddIteration(iterationInfo(elements[mid:], true, dir))
			b.bitonicMerge(elements[mid:], dir)
		}()
		wg.Wait()
	}
}

func compareAndSwap(s []*models.TFIDF, i, j, dir int) {
	if (s[i].Less(s[j]) && dir == 0) || (!s[i].Less(s[j]) && dir == 1) {
		s[i], s[j] = s[j], s[i]
	}
}

func isPowerOfTwo(n int) bool {
	return n != 0 && (n&(n-1)) == 0
}

func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

func iterationInfo(elements []*models.TFIDF, isMergeOperation bool, dir int) (string, string) {
	direction := ""
	if dir == 1 {
		direction = "ascending"
	} else {
		direction = "descending"
	}
	if isMergeOperation {
		return fmt.Sprintf(`thread spawned with elements to be merged with %s direction (at this point elements are already compared and swapped in %s direction`,
			direction, direction), models.GetTFIDFElements(elements)
	}
	return fmt.Sprintf("thread spawned with elements to be sorted with direction %s",
		direction), models.GetTFIDFElements(elements)
}
