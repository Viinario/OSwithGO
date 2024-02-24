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

	line := fmt.Sprintf("Thread: %s ID: %d, IO Time: %d ms\n ", name, id, totalTime)
	if _, err := file.WriteString(line); err != nil {
		fmt.Println("Error writing to IO file:", err)
		return
	}
	time.Sleep(time.Duration(totalTime) * time.Millisecond) // Simula o Thread utilizando o tempo de IO
	defer file.Close()
}
