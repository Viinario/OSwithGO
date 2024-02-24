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
		fmt.Println("Error opening CPU file:", err)
		return
	}

	line := fmt.Sprintf("Thread: %s ID: %d, IO Time: %d ms\n ", name, id, totalTime)
	if _, err := file.WriteString(line); err != nil {
		fmt.Println("Error writing to CPU file:", err)
		return
	}

	time.Sleep(time.Duration(totalTime) * time.Millisecond) // Simula o Thread utilizando o tempo de CPU
	defer file.Close()
}
