package cpu

import (
	"fmt"
	"os"
	"time"
)

// UseCPU simula o uso da CPU pela thread
func UseCPU(id int, totalTime int) {
	file, err := os.OpenFile("cpu.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening CPU file:", err)
		return
	}
	defer file.Close()

	line := fmt.Sprintf("Thread ID: %d, CPU Time: %d\n", id, totalTime)
	if _, err := file.WriteString(line); err != nil {
		fmt.Println("Error writing to CPU file:", err)
		return
	}

	time.Sleep(time.Duration(totalTime) * time.Millisecond)
}
