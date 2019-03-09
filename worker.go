package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/umahmood/haversine"
)

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

func worker(points []Point) {
	if len(points) < 0 {
		return
	}

	resultsChannel := make(chan string)

	var fare = 1.30
	var rideId int64

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

		rideId = points[i].RideId

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

	buf := bytes.Buffer

	resultsChannel, _ := <-fmt.Printf("id_ride: %v, fare_estimate: %.2f\n", rideId, fare)
}

// func writeToFile()

func calculateFareAmount(speed float64, duration float64, distance float64, nightShift bool) float64 {
	if speed > 10 {
		if nightShift {
			return distance * 1.3
		}
		return distance * 0.74
	}

	return duration * 11.9
}
