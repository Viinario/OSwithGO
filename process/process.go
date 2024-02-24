package process

import (
	"strings"
	"sync"

	"github.com/Viinario/OSwithGO/cpu"
	"github.com/Viinario/OSwithGO/io"
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
	wg               sync.WaitGroup // Usado para sincronizar a conclusão da thread
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

// Start inicia a execução da thread
func (t *Thread) Start(cpuTime int, ioTime int) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		if t.IOBound {
			if t.RemainingIOTime > 0 {
				io.UseIO(t.ID, ioTime)
			}
			if t.RemainingCPUTime > 0 {
				cpu.UseCPU(t.ID, cpuTime)
			}
		} else {
			if t.RemainingCPUTime > 0 {
				cpu.UseCPU(t.ID, cpuTime)
			}
			if t.RemainingIOTime > 0 {
				io.UseIO(t.ID, ioTime)
			}
		}
	}()
}

// Wait espera até que a thread termine sua execução
func (t *Thread) Wait() {
	t.wg.Wait()
}

// generateID gera um ID único para a thread
func generateID() int {
	idMutex.Lock()
	defer idMutex.Unlock()
	idCounter++
	return idCounter
}
