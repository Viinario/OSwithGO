package process

import (
	"strings"
	"sync"
)

var idCounter int
var idMutex sync.Mutex

// Thread representa uma thread
type Thread struct {
	ID               int
	Name             string
	Priority         int
	IOBound          bool
	TotalCPUTime     int
	RemainingCPUTime int
	TotalIOTime      int
	RemainingIOTime  int
}

// NewThread cria uma nova thread manualmente
func NewThread(name string, priority int, ioBoundInput string, totalCPUTime int, totalIOTime int) *Thread {
	ioBound := strings.ToLower(ioBoundInput) == "s"
	return &Thread{
		ID:               generateID(),
		Name:             name,
		Priority:         priority,
		IOBound:          ioBound,
		TotalCPUTime:     totalCPUTime,
		RemainingCPUTime: totalCPUTime,
		TotalIOTime:      totalIOTime,
		RemainingIOTime:  totalIOTime,
	}
}

// generateID gera um ID único para a thread
func generateID() int {
	idMutex.Lock()
	defer idMutex.Unlock()
	idCounter++
	return idCounter
}
