package strategies

import (
	"cli-search-engine/engineLogger"
	"cli-search-engine/models"
	"fmt"
)

type QuickSort struct {
	Logger   engineLogger.Logger
	Elements []*models.TFIDF
}

func NewQuickSort(elements []*models.TFIDF, logger engineLogger.Logger) *QuickSort {
	logger.SetStartMessage(fmt.Sprintf("initialized quickSort with elements of length %d", len(elements)))
	return &QuickSort{
		Logger:   logger,
		Elements: elements,
	}
}

func (q *QuickSort) Sort() string {
	q.quickSort(q.Elements, 0, len(q.Elements)-1)
	err := q.Logger.Log()
	if err != nil {
		return fmt.Sprintf("error while logging quickSort with elements of length %d", len(q.Elements))
	}
	return fmt.Sprintf("top results are %s", models.GetTFIDFElements(q.Elements))
}

func (q *QuickSort) partition(arr []*models.TFIDF, low, high int) int {
	pivot := arr[high].Tfidf
	i := low - 1
	for j := low; j <= high-1; j++ {
		if arr[j].Tfidf > pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// Limited can't multithreading without memory overloading
func (q *QuickSort) quickSort(arr []*models.TFIDF, low, high int) {
	if low < high {
		pi := q.partition(arr, low, high)
		q.Logger.AddQuickSortIteration(low, high, models.GetTFIDFElements(arr))
		q.quickSort(arr, low, pi-1)
		q.Logger.AddQuickSortIteration(low, pi-1, models.GetTFIDFElements(arr))
		q.quickSort(arr, pi+1, high)
		q.Logger.AddQuickSortIteration(pi+1, high, models.GetTFIDFElements(arr))
	}
}
