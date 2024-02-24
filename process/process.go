package process

import (
	"math/rand"
	"strings"
	"time"
)

// Process representa um processo
type Process struct {
	ID               int
	Name             string
	Priority         int
	IOBound          bool
	TotalCPUTime     int
	RemainingCPUTime int
}

// NewProcess cria um novo processo manualmente
func NewProcess(name string, priority int, ioBoundInput string, totalCPUTime int) *Process {
	ioBound := strings.ToLower(ioBoundInput) == "s"
	return &Process{
		ID:               generateID(),
		Name:             name,
		Priority:         priority,
		IOBound:          ioBound,
		TotalCPUTime:     totalCPUTime,
		RemainingCPUTime: totalCPUTime,
	}
}

// generateID gera um ID único para o processo
func generateID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000) + 1
}
