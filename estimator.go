package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {

	writeChannel := make(chan []string)
	doneWritingChannel := make(chan bool)
	ridePointsChannel := make(chan []Point)

	file, err := os.Open("paths.csv")
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.OpenFile("output.csv", os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	writer := csv.NewWriter(outputFile)

	go scanFile(scanner, file, ridePointsChannel)

	go appendToFile(writer, writeChannel, doneWritingChannel)

	for points := range ridePointsChannel {
		wg.Add(1)
		go worker(points, writeChannel)
	}

	wg.Wait()
}
