package main

import (
	"fmt"
	"os"

	"github.com/Viinario/OSwithGO/scheduler"
)

func main() {
	s := scheduler.NewScheduler()
	for {
		fmt.Println("Menu:")
		fmt.Println("1. Criar processo")
		fmt.Println("2. Escolher algoritmo de escalonamento")
		fmt.Println("3. Selecionar tempo de quantum")
		fmt.Println("4. Mostrar fila de processos prontos")
		fmt.Println("5. Iniciar a execução e escalonamento")
		fmt.Println("6. Criar Processos Aleatórios")
		fmt.Println("7. Sair")

		var choice int
		fmt.Print("Opção: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Erro ao ler a entrada:", err)
			os.Exit(1)
		}

		switch choice {
		case 1:
			var name string
			var priority int
			var ioBoundInput string
			var totalCPUTime int
			var totalIOTime int

			fmt.Print("Digite o nome do processo: ")
			fmt.Scanln(&name)
			fmt.Print("Digite a prioridade (maior número tem prioridade): ")
			fmt.Scanln(&priority)
			fmt.Print("O processo é I/O bound (S/N)? ")
			fmt.Scanln(&ioBoundInput)
			fmt.Print("Digite o tempo total de CPU (Ex: 10): ")
			fmt.Scanln(&totalCPUTime)
			fmt.Print("Digite o tempo total de IO (Ex: 10): ")
			fmt.Scanln(&totalIOTime)

			s.CreateProcess(name, priority, ioBoundInput, totalCPUTime, totalIOTime)

		case 2:
			var algorithmChoice int
			fmt.Println("Escolha o algoritmo de escalonamento:")
			fmt.Println("1. Round Robin")
			fmt.Println("2. Prioridade")
			fmt.Print("Opção: ")
			_, err := fmt.Scanln(&algorithmChoice)
			if err != nil {
				fmt.Println("Erro ao ler a entrada:", err)
				os.Exit(1)
			}
			s.ChooseAlgorithm(algorithmChoice)

		case 3:
			var quantum int
			fmt.Print("Digite o tempo de quantum da preempção (1 a 10 ms): ")
			_, err := fmt.Scanln(&quantum)
			if err != nil {
				fmt.Println("Erro ao ler a entrada:", err)
				os.Exit(1)
			}
			s.SelectQuantum(quantum)

		case 4:
			s.PrintReadyQueue()

		case 5:
			s.RunSimulation()

		case 6:
			var numProcesses int
			fmt.Print("Digite o número de processos que você gostaria de criar: ")
			_, err := fmt.Scanln(&numProcesses)
			if err != nil {
				fmt.Println("Erro ao ler a entrada:", err)
				os.Exit(1)
			}
			s.CreateRandomProcesses(numProcesses)

		case 7:
			fmt.Println("Saindo...")
			os.Exit(0)

		default:
			fmt.Println("Opção inválida. Tente novamente.")
		}
	}
}
