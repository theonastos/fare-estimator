package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// Maps a record to a Point instance
func mapStringToPoint(s []string) Point {

	id, err := strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	latitude, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		log.Fatal(err)
	}
	longitude, err := strconv.ParseFloat(s[2], 64)
	if err != nil {
		log.Fatal(err)
	}
	timestampNum, err := strconv.ParseInt(s[3], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	timestamp := time.Unix(timestampNum, 0)

	return Point{
		RideId:    id,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: timestamp}
}

func main() {
	var ridePoints []Point

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

	lastId := int64(-1)

	go func() {
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), ",")
			point := mapStringToPoint(s)

			ridePoints = append(ridePoints, point)

			if point.RideId != lastId {
				if lastId != -1 {
					ridePointsChannel <- ridePoints
				}
				lastId = point.RideId
				// resetting slice - reallocating memory is really slow,
				// so keeping the capacity intact should work better
				ridePoints = nil
				ridePoints = append(ridePoints, point)
			}
		}
		ridePointsChannel <- ridePoints
		close(ridePointsChannel)
	}()

	go appendToFile(writer, writeChannel, doneWritingChannel)

	for points := range ridePointsChannel {
		wg.Add(1)
		go worker(points, writeChannel)
	}

	wg.Wait()
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	dist = dist * 1.609344

	return dist
}
