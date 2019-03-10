package main

import (
	"encoding/csv"
	"fmt"
	"log"
)

// Receives a csv writer and a write channel which contains the final rows to be written in the file
func appendToFile(writer *csv.Writer, writeChannel <-chan []string, doneWritingChannel chan<- bool) {
	for row := range writeChannel {
		fmt.Println(row)
		if err := writer.Write(row); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
		writer.Flush()
	}
}
