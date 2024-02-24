package scheduler

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/Viinario/OSwithGO/process"
)

// Scheduler gerencia o escalonamento de processos
type Scheduler struct {
	Processes      []*process.Thread
	ReadyQueue     []*process.Thread
	CurrentProcess *process.Thread
	Algorithm      func()
	Quantum        int
}

// NewScheduler cria uma nova instância de Scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// CreateProcess cria um novo processo manualmente
func (s *Scheduler) CreateProcess(name string, priority int, ioBoundInput string, totalCPUTime int, totalIOTime int) {
	process := process.NewThread(name, priority, ioBoundInput, totalCPUTime, totalIOTime)
	s.Processes = append(s.Processes, process)
	s.ReadyQueue = append(s.ReadyQueue, process)
}

// CreateRandomProcesses cria um número especificado de processos aleatórios
func (s *Scheduler) CreateRandomProcesses(numProcesses int) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numProcesses; i++ {
		name := randString(3)
		priority := rand.Intn(10) + 1
		totalCPUTime := rand.Intn(10) + 1
		totalIOTime := rand.Intn(10) + 1
		var ioBound string
		if totalIOTime > totalCPUTime {
			ioBound = "s"
		} else {
			ioBound = "n"
		}

		process := process.NewThread(name, priority, ioBound, totalCPUTime, totalIOTime)

		s.Processes = append(s.Processes, process)
		s.ReadyQueue = append(s.ReadyQueue, process)
	}
}

// RandString gera uma string aleatória de comprimento n
func randString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ChooseAlgorithm escolhe um algoritmo de escalonamento
func (s *Scheduler) ChooseAlgorithm(choice int) {
	if choice == 1 {
		s.Algorithm = s.RoundRobin
	} else if choice == 2 {
		s.Algorithm = s.Priority
	}
}

// SelectQuantum seleciona o tempo de quantum para o algoritmo Round Robin
func (s *Scheduler) SelectQuantum(quantum int) {
	s.Quantum = quantum
}

// PrintReadyQueue imprime a fila de processos prontos
func (s *Scheduler) PrintReadyQueue() {
	fmt.Println("Fila de processos prontos:")
	for _, process := range s.ReadyQueue {
		ioBoundInfo := "CPU Bound"
		if process.IOBound {
			ioBoundInfo = "I/O Bound"
		}
		fmt.Printf("ID: %d, Nome: %s, Prioridade: %d, CPU Time Restante: %d ms, IO Time Restante: %d ms, Tipo: %s\n", process.ID, process.Name, process.Priority, process.RemainingCPUTime, process.RemainingIOTime, ioBoundInfo)
	}
}

// Aqui o algoritmo ira dar um tempo igual de uso de cpu de io chamado preempção.
// Tera dois semaforos, um para CPU e outro para IO
// Caso um processo já esteja utilizando o IO, ele ira esperar ate o semaforo seja liberado
// a preempção devera ocorrer tanto no cpu bound quanto no io bound

// RoundRobin executa o algoritmo Round Robin
func (s *Scheduler) RoundRobin() {
	if s.CurrentProcess == nil {
		if len(s.ReadyQueue) > 0 {
			s.CurrentProcess = s.ReadyQueue[0]
			s.ReadyQueue = s.ReadyQueue[1:]
		}
	}

	if s.CurrentProcess != nil {
		if s.CurrentProcess.RemainingCPUTime <= s.Quantum {
			fmt.Printf("Processo %s está na CPU por %d ms.\n", s.CurrentProcess.Name, s.CurrentProcess.RemainingCPUTime)
			s.CurrentProcess.RemainingCPUTime = 0
			s.FinishProcess()
		} else {
			fmt.Printf("Processo %s está na CPU por %d ms.\n", s.CurrentProcess.Name, s.Quantum)
			s.CurrentProcess.RemainingCPUTime -= s.Quantum
			s.PreemptProcess()
		}
	}
}

// Priority executa o algoritmo de prioridade
func (s *Scheduler) Priority() {
	if s.CurrentProcess == nil {
		s.ReadyQueue = s.sortByPriority(s.ReadyQueue)
		if len(s.ReadyQueue) > 0 {
			s.CurrentProcess = s.ReadyQueue[0]
			s.ReadyQueue = s.ReadyQueue[1:]
		}
	}

	if s.CurrentProcess != nil {
		if s.CurrentProcess.RemainingCPUTime <= s.Quantum {
			fmt.Printf("Processo %s está na CPU por %d ms.\n", s.CurrentProcess.Name, s.CurrentProcess.RemainingCPUTime)
			s.CurrentProcess.RemainingCPUTime = 0
			s.FinishProcess()
		} else {
			fmt.Printf("Processo %s está na CPU por %d ms.\n", s.CurrentProcess.Name, s.Quantum)
			s.CurrentProcess.RemainingCPUTime -= s.Quantum
			s.PreemptProcess()
		}
	}
}

// sortByPriority ordena a fila de processos por prioridade
func (s *Scheduler) sortByPriority(processes []*process.Thread) []*process.Thread {
	sortedProcesses := make([]*process.Thread, len(processes))
	copy(sortedProcesses, processes)
	sort.Slice(sortedProcesses, func(i, j int) bool {
		return sortedProcesses[i].Priority > sortedProcesses[j].Priority
	})
	return sortedProcesses
}

// PreemptProcess realiza a preempção e adiciona o processo de volta à fila de prontos
func (s *Scheduler) PreemptProcess() {
	s.ReadyQueue = append(s.ReadyQueue, s.CurrentProcess)
	s.CurrentProcess = nil
}

// FinishProcess finaliza um processo
func (s *Scheduler) FinishProcess() {
	fmt.Printf("Processo %s finalizado.\n", s.CurrentProcess.Name)
	s.CurrentProcess = nil
}

// RunSimulation executa a simulação de escalonamento
func (s *Scheduler) RunSimulation() {
	if s.Algorithm == nil {
		fmt.Println("Por favor, escolha um algoritmo de escalonamento antes de iniciar a execução.")
		return
	}
	if s.Quantum == 0 {
		fmt.Println("Por favor, selecione um tempo de quantum antes de iniciar a execução.")
		return
	}
	if len(s.ReadyQueue) == 0 {
		fmt.Println("Não há processos na fila de processos prontos. Crie processos antes de iniciar a execução.")
		return
	}

	for len(s.ReadyQueue) > 0 {
		s.PrintReadyQueue()
		s.Algorithm()
	}

	s.CalculateTurnaroundTime()
}

// CalculateTurnaroundTime calcula o tempo médio de espera dos processos
func (s *Scheduler) CalculateTurnaroundTime() float64 {
	totalTurnaroundTime := 0
	totalCPUTime := 0

	for _, process := range s.Processes {
		totalTurnaroundTime += process.TotalCPUTime - process.RemainingCPUTime
		totalCPUTime += process.TotalCPUTime
	}

	averageWaitingTime := float64(totalTurnaroundTime) / float64(len(s.Processes))
	fmt.Printf("Tempo médio de espera: %.2f ms\n", averageWaitingTime)

	return averageWaitingTime
}
