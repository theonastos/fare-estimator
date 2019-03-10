package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/umahmood/haversine"
)

var mu sync.Mutex

type LockedWriter struct {
	m      sync.Mutex
	Writer io.Writer
}

func (lw *LockedWriter) Write(b []byte) (n int, err error) {
	lw.m.Lock()
	defer lw.m.Unlock()
	return lw.Writer.Write(b)
}

type Point struct {
	RideId    int64
	Longitude float64
	Latitude  float64
	Timestamp time.Time
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func IntToString(input_num int64) string {
	// to convert a int number to a string
	return strconv.FormatInt(input_num, 10)
}

func worker(points []Point, writeChannel chan<- []string) {
	defer wg.Done()

	if len(points) <= 0 {
		return
	}

	var fare = 1.30
	var rideId = points[0].RideId

	for i := 0; i < len(points)-2; i++ {
		point1 := haversine.Coord{
			Lat: (points)[i].Latitude,
			Lon: (points)[i].Longitude}

		point2 := haversine.Coord{
			Lat: (points)[i+1].Latitude,
			Lon: (points)[i+1].Longitude}

		_, d := haversine.Distance(point1, point2)

		t := ((points)[i+1].Timestamp.Sub((points)[i].Timestamp)).Hours()

		u := d / t

		if u > 100 {
			(points) = append((points)[:i+1], (points)[i+2:]...)
			i--
		} else {
			if points[i].Timestamp.Hour() > 00 && points[i].Timestamp.Hour() < 05 {
				fare += calculateFareAmount(u, t, d, true)
			} else {
				fare += calculateFareAmount(u, t, d, false)
			}
		}
	}

	if fare < 3.47 {
		fare = 3.47
	}

	row := []string{IntToString(rideId), FloatToString(fare)}

	writeChannel <- row
}

func appendToFile(writer *csv.Writer, writeChannel <-chan []string, done chan<- bool) {
	for row := range writeChannel {
		fmt.Println(row)
		if err := writer.Write(row); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
		writer.Flush()
	}

	// Using the synchronization approach to signal when done writing
	// done <- true
}

func calculateFareAmount(speed float64, duration float64, distance float64, nightShift bool) float64 {
	if speed > 10 {
		if nightShift {
			return distance * 1.3
		}
		return distance * 0.74
	}

	return duration * 11.9
}
