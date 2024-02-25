package io

import (
	"fmt"
	"os"
	"time"
)

// UseIO simula o uso do IO pela thread
func UseIO(name string, id int, totalTime int) {
	file, err := os.OpenFile("io.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening IO file:", err)
		return
	}
	defer file.Close()

	line := fmt.Sprintf("|Thread: %s ID: %d | IO Time: %d ms | simulating: ", name, id, totalTime)
	for i := totalTime; i >= 0; i-- {
		line += fmt.Sprintf(" %d ms \n", i)
		if i > 0 {
			time.Sleep(time.Millisecond) // Simula o Thread utilizando o tempo de IO
		}
	}

	line += "\n"

	if _, err := file.WriteString(line); err != nil {
		fmt.Println("Error writing thread data to IO file:", err)
		return
	}
}
