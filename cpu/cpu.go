package cpu

import (
	"fmt"
	"os"
	"time"
)

// UseCPU simula o uso da CPU pela thread
func UseCPU(name string, id int, totalTime int) {
	file, err := os.OpenFile("cpu.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening IO file:", err)
		return
	}
	defer file.Close()

	line := fmt.Sprintf("| Thread: %s ID: %d | CPU Time: %d ms | simulating:", name, id, totalTime)
	for i := totalTime; i >= 0; i-- {
		line += fmt.Sprintf(" %d ms \n", i)
		if i > 0 {
			time.Sleep(time.Millisecond) // Simula o Thread utilizando o tempo de CPU
		}
	}

	line += "\n"

	if _, err := file.WriteString(line); err != nil {
		fmt.Println("Error writing thread data to IO file:", err)
		return
	}
}
