package scheduler

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/Viinario/OSwithGO/cpu"
	"github.com/Viinario/OSwithGO/io"
	"github.com/Viinario/OSwithGO/process"
)

var (
	cpuResource = make(chan bool, 1) // Canal para controlar o recurso da CPU
	ioResource  = make(chan bool, 1) // Canal para controlar o recurso de IO
	wg          sync.WaitGroup
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

// PreemptProcess realiza a preempção e adiciona o processo de volta à fila de prontos
func (s *Scheduler) PreemptProcess(process *process.Thread) {
	s.FinishProcess(process)
	//s.ReadyQueue = append(s.ReadyQueue, process) // Adiciona o processo ao final da fila
}

// FinishProcess finaliza um processo
func (s *Scheduler) FinishProcess(process *process.Thread) {
	// Encontrar o índice do processo a ser removido
	index := -1
	for i, p := range s.ReadyQueue {
		if p == process { // Verifica se o ponteiro é o mesmo
			index = i
			break
		}
	}

	// Se o processo foi encontrado, remova-o
	if index != -1 {
		// Remover o processo do slice
		s.ReadyQueue = append(s.ReadyQueue[:index], s.ReadyQueue[index+1:]...)
	}
}

// RoundRobin executa o algoritmo Round Robin
func (s *Scheduler) RoundRobin() {
	cpuBoundQueue := make([]*process.Thread, 0)
	ioBoundQueue := make([]*process.Thread, 0)

	// Separa os processos CPU-bound e I/O-bound
	for _, p := range s.ReadyQueue {
		if p.IOBound {
			ioBoundQueue = append(ioBoundQueue, p)
		} else {
			cpuBoundQueue = append(cpuBoundQueue, p)
		}
	}

	// Verifica se há um processo CPU-bound e um I/O-bound para executar
	if len(cpuBoundQueue) > 0 && len(ioBoundQueue) > 0 {
		s.executeProcess(cpuBoundQueue[0], ioBoundQueue[0])
	} else if len(cpuBoundQueue) > 0 {
		// Se houver apenas um processo CPU-bound, executa-o
		s.executeProcess(cpuBoundQueue[0], nil)
	} else if len(ioBoundQueue) > 0 {
		// Se houver apenas um processo I/O-bound, executa-o
		s.executeProcess(nil, ioBoundQueue[0])
	}
}

func (s *Scheduler) executeCpuProcess(cpuProcess *process.Thread, wg *sync.WaitGroup) {
	defer wg.Done()
	if cpuProcess != nil {
		if cpuProcess.RemainingCPUTime <= s.Quantum { // Se o tempo de CPU for menor que a da preempção, executa o tempo todo:
			if cpuProcess.RemainingCPUTime > 0 {
			Loop_Cpu:
				for {
					select {
					case cpuResource <- true:
						fmt.Printf("Processo CPU-bound %s está na CPU por %d ms.\n", cpuProcess.Name, cpuProcess.RemainingCPUTime)
						cpu.UseCPU(cpuProcess.Name, cpuProcess.ID, cpuProcess.RemainingCPUTime)
						cpuProcess.RemainingCPUTime = 0
						<-cpuResource
						break Loop_Cpu
					default:
						fmt.Printf("CPU Bound: %s Recurso CPU ocupado, esperando...\n", cpuProcess.Name)
						time.Sleep(2 * time.Millisecond)
					}
				}
			}
			if cpuProcess.RemainingIOTime <= s.Quantum { // Se o tempo de IO for menor que a da preempção, executa o tempo todo:
				if cpuProcess.RemainingIOTime > 0 {
					for {
						select {
						case ioResource <- true:
							fmt.Printf("Processo CPU-bound %s está na E/S por %d ms.\n", cpuProcess.Name, cpuProcess.RemainingIOTime)
							io.UseIO(cpuProcess.Name, cpuProcess.ID, cpuProcess.RemainingIOTime)
							cpuProcess.RemainingIOTime = 0
							s.FinishProcess(cpuProcess)
							fmt.Printf("Processo %s finalizado.\n", cpuProcess.Name)
							<-ioResource
							return
						default:
							fmt.Printf("CPU Bound: %s Recurso E/S ocupado, esperando...\n", cpuProcess.Name)
							time.Sleep(2 * time.Millisecond)
						}
					}
				} else {
					s.FinishProcess(cpuProcess)
					fmt.Printf("Processo %s finalizado.\n", cpuProcess.Name)
					return
				}
			} else if cpuProcess.RemainingIOTime > s.Quantum { // Se o tempo de IO for maior que a da preempção, executa o tempo da preempção:
			Loop_Cpu_3:
				for {
					select {
					case ioResource <- true:
						fmt.Printf("Processo CPU-bound %s está na E/S por %d ms.\n", cpuProcess.Name, s.Quantum)
						io.UseIO(cpuProcess.Name, cpuProcess.ID, s.Quantum)
						cpuProcess.RemainingIOTime -= s.Quantum
						s.PreemptProcess(cpuProcess)
						<-ioResource
						break Loop_Cpu_3
					default:
						fmt.Printf("CPU Bound: %s Recurso E/S ocupado, esperando...\n", cpuProcess.Name)
						time.Sleep(2 * time.Millisecond)
					}
				}
			}
		} else if cpuProcess.RemainingCPUTime > s.Quantum { // Se o tempo de CPU for maior que a da preempção, executa o tempo da preempção:
		Loop_Cpu_4:
			for {
				select {
				case cpuResource <- true:
					fmt.Printf("Processo CPU-bound %s está na CPU por %d ms.\n", cpuProcess.Name, s.Quantum)
					cpu.UseCPU(cpuProcess.Name, cpuProcess.ID, s.Quantum)
					cpuProcess.RemainingCPUTime -= s.Quantum
					<-cpuResource
					break Loop_Cpu_4
				default:
					fmt.Printf("CPU Bound: %s Recurso CPU ocupado, esperando...\n", cpuProcess.Name)
					time.Sleep(2 * time.Millisecond)
				}
			}
			if cpuProcess.RemainingIOTime <= s.Quantum { // Se o tempo de IO for menor que a da preempção, executa o tempo todo:
				if cpuProcess.RemainingIOTime > 0 {
				Loop_Cpu_5:
					for {
						select {
						case ioResource <- true:
							fmt.Printf("Processo CPU-bound %s está na E/S por %d ms.\n", cpuProcess.Name, cpuProcess.RemainingIOTime)
							io.UseIO(cpuProcess.Name, cpuProcess.ID, cpuProcess.RemainingIOTime)
							cpuProcess.RemainingIOTime = 0
							s.PreemptProcess(cpuProcess)
							<-ioResource
							break Loop_Cpu_5
						default:
							fmt.Printf("CPU Bound: %s Recurso E/S ocupado, esperando...\n", cpuProcess.Name)
							time.Sleep(2 * time.Millisecond)
						}
					}
				} else {
					s.PreemptProcess(cpuProcess)
				}
			} else if cpuProcess.RemainingIOTime > s.Quantum { // Se o tempo de IO for maior que a da preempção, executa o tempo da preempção:
			Loop_Cpu_6:
				for {
					select {
					case ioResource <- true:
						fmt.Printf("Processo CPU-bound %s está na E/S por %d ms.\n", cpuProcess.Name, s.Quantum)
						io.UseIO(cpuProcess.Name, cpuProcess.ID, s.Quantum)
						cpuProcess.RemainingIOTime -= s.Quantum
						s.PreemptProcess(cpuProcess)
						<-ioResource
						break Loop_Cpu_6
					default:
						fmt.Printf("CPU Bound: %s Recurso E/S ocupado, esperando...\n", cpuProcess.Name)
						time.Sleep(2 * time.Millisecond)
					}
				}
			}
		}
	}
}

func (s *Scheduler) executeIOProcess(ioProcess *process.Thread, wg *sync.WaitGroup) {
	defer wg.Done()
	if ioProcess != nil {
		if ioProcess.RemainingIOTime <= s.Quantum { // se o tempo de IO for menor que da preempção, executa o tempo todo:
			if ioProcess.RemainingIOTime > 0 {
			Loop_IO:
				for {
					select {
					case ioResource <- true:
						fmt.Printf("Processo I/O-bound %s está na E/S por %d ms.\n", ioProcess.Name, ioProcess.RemainingIOTime)
						io.UseIO(ioProcess.Name, ioProcess.ID, ioProcess.RemainingIOTime)
						ioProcess.RemainingIOTime = 0
						<-ioResource
						break Loop_IO
					default:
						fmt.Printf("I/O Bound: %s Recurso E/S ocupado, esperando...\n", ioProcess.Name)
						time.Sleep(2 * time.Millisecond)
					}
				}
			}
			if ioProcess.RemainingCPUTime <= s.Quantum { // se o tempo de CPU for menor que da preempção, executa o tempo todo
				if ioProcess.RemainingCPUTime > 0 {
					for {
						select {
						case cpuResource <- true:
							fmt.Printf("Processo I/O-bound %s está na CPU por %d ms.\n", ioProcess.Name, ioProcess.RemainingCPUTime)
							cpu.UseCPU(ioProcess.Name, ioProcess.ID, ioProcess.RemainingCPUTime)
							ioProcess.RemainingCPUTime = 0
							s.FinishProcess(ioProcess)
							fmt.Printf("Processo %s finalizado.\n", ioProcess.Name)
							<-cpuResource
							return
						default:
							fmt.Printf("I/O Bound: %s Recurso CPU ocupado, esperando...\n", ioProcess.Name)
							time.Sleep(2 * time.Millisecond)
						}
					}
				} else {
					s.FinishProcess(ioProcess)
					fmt.Printf("Processo %s finalizado.\n", ioProcess.Name)
					return
				}
			} else if ioProcess.RemainingCPUTime > s.Quantum { // se o tempo de CPU for maior que da preempção, executa o tempo da preempção
			Loop_IO_3:
				for {
					select {
					case cpuResource <- true:
						fmt.Printf("Processo I/O-bound %s está na CPU por %d ms.\n", ioProcess.Name, s.Quantum)
						cpu.UseCPU(ioProcess.Name, ioProcess.ID, s.Quantum)
						ioProcess.RemainingCPUTime -= s.Quantum
						<-cpuResource
						break Loop_IO_3
					default:
						fmt.Printf("I/O Bound: %s Recurso CPU ocupado, esperando...\n", ioProcess.Name)
						time.Sleep(2 * time.Millisecond)
					}
				}
			}
		} else if ioProcess.RemainingIOTime > s.Quantum { // se o tempo de IO for maior que da preempção, executa o tempo da preempção
		Loop_IO_4:
			for {
				select {
				case ioResource <- true:
					fmt.Printf("Processo I/O-bound %s está na E/S por %d ms.\n", ioProcess.Name, s.Quantum)
					io.UseIO(ioProcess.Name, ioProcess.ID, s.Quantum)
					ioProcess.RemainingIOTime -= s.Quantum
					<-ioResource
					break Loop_IO_4
				default:
					fmt.Printf("I/O Bound: %s Recurso E/S ocupado, esperando...\n", ioProcess.Name)
					time.Sleep(2 * time.Millisecond)
				}
			}
			if ioProcess.RemainingCPUTime <= s.Quantum { // se o tempo de CPU for menor que da preempção, executa o tempo todo
				if ioProcess.RemainingCPUTime > 0 {
				Loop_IO_5:
					for {
						select {
						case cpuResource <- true:
							fmt.Printf("Processo I/O-bound %s está na CPU por %d ms.\n", ioProcess.Name, ioProcess.RemainingCPUTime)
							cpu.UseCPU(ioProcess.Name, ioProcess.ID, ioProcess.RemainingCPUTime)
							ioProcess.RemainingCPUTime = 0
							s.PreemptProcess(ioProcess)
							<-cpuResource
							break Loop_IO_5
						default:
							fmt.Printf("I/O Bound: %s Recurso CPU ocupado, esperando...\n", ioProcess.Name)
							time.Sleep(2 * time.Millisecond)
						}
					}
				} else {
					s.PreemptProcess(ioProcess)
				}
			} else if ioProcess.RemainingCPUTime > s.Quantum { // se o tempo de CPU e IO for maior que da preempção, executa o tempo da preempção
			Loop_IO_6:
				for {
					select {
					case cpuResource <- true:
						fmt.Printf("Processo I/O-bound %s está na CPU por %d ms.\n", ioProcess.Name, s.Quantum)
						cpu.UseCPU(ioProcess.Name, ioProcess.ID, s.Quantum)
						ioProcess.RemainingCPUTime -= s.Quantum
						s.PreemptProcess(ioProcess)
						<-cpuResource
						break Loop_IO_6
					default:
						fmt.Printf("I/O Bound: %s Recurso CPU ocupado, esperando...\n", ioProcess.Name)
						time.Sleep(2 * time.Millisecond)
					}
				}
			}
		}
	}
}

// executeProcess executa um processo CPU-bound e um I/O-bound simultaneamente
func (s *Scheduler) executeProcess(cpuProcess, ioProcess *process.Thread) {
	wg.Add(2)
	// Executa o processo CPU-bound, se disponível
	go s.executeCpuProcess(cpuProcess, &wg)
	// Executa o processo I/O-bound, se disponível
	go s.executeIOProcess(ioProcess, &wg)
	// Aguarda a conclusão de ambos os processos
	wg.Wait()
	// Atualiza os processos na fila de prontos
	s.updateReadyQueue(cpuProcess, ioProcess)
}

// updateReadyQueue atualiza a fila de prontos após a execução dos processos
func (s *Scheduler) updateReadyQueue(cpuProcess, ioProcess *process.Thread) {
	if cpuProcess != nil && cpuProcess.RemainingCPUTime > 0 {
		s.ReadyQueue = append(s.ReadyQueue, cpuProcess)
	}
	if ioProcess != nil && ioProcess.RemainingIOTime > 0 {
		s.ReadyQueue = append(s.ReadyQueue, ioProcess)
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
			//s.FinishProcess()
		} else {
			fmt.Printf("Processo %s está na CPU por %d ms.\n", s.CurrentProcess.Name, s.Quantum)
			s.CurrentProcess.RemainingCPUTime -= s.Quantum
			//s.PreemptProcess()
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
