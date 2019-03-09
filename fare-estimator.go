package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/umahmood/haversine"
)

type Point struct {
	RideId    int64
	Longitude float64
	Latitude  float64
	Timestamp time.Time
}

var (
	// list of channels to communicate with workers
	// workers accessed synchronousely no mutex required
	workers = make(map[string]chan []string)

	// wg is to make sure all workers done before exiting main
	wg = sync.WaitGroup{}

	// mu used only for sequential printing, not relevant for program logic
	mu = sync.Mutex{}
)

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
	defer wg.Wait()

	var ridePoints []Point
	ridePointsChannel := make(chan []Point)

	file, err := os.Open("paths.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

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

	for points := range ridePointsChannel {
		go worker(points)
	}
}

func worker(points []Point) {
	if len(points) < 0 {
		return
	}

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

	fmt.Printf("RideId: %v, Fare: %.2f\n", rideId, fare)
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
