package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	RideId    int64
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

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

// ScanFile is responsible for the hard job of fetching the ride points data from the initial csv file
func scanFile(scanner *bufio.Scanner, file *os.File, ridePointsChannel chan<- []Point) {

	lastId := int64(-1)

	var ridePoints []Point

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
}
