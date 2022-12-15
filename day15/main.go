package main

import (
	_ "embed"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)


//go:embed input.txt
var input string

var minY, maxY, minX, maxX int

type loc struct {
	beacon, sensor, noBeacon bool
}

type coord struct {
	x, y int
}

type bubble struct {
	center coord
	radius int
}

func main() {
	lineRE := `Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`
	re := regexp.MustCompile(lineRE)

	grid := map[int]map[int]*loc{}
	bubbles := []bubble{}

	lines := strings.Split(input, "\n")
	for i, line := range lines {
		fmt.Printf("adding sensor %d\n", i)
		matches := re.FindStringSubmatch(line)
		sensorX, _ := strconv.Atoi(matches[1])
		sensorY, _ := strconv.Atoi(matches[2])
		beaconX, _ := strconv.Atoi(matches[3])
		beaconY, _ := strconv.Atoi(matches[4])

		addSensor(grid, &bubbles, coord{sensorX,sensorY}, coord{beaconX,beaconY})
	}

	rowCount := 0
	for _, l := range grid[2000000] {
		if l.noBeacon || l.sensor {
			rowCount++
		}
	}

	fmt.Printf("row count: %d\n", rowCount)

	var distressCoord *coord
	for y := 0; y <= 4000000; y++ {
		for x := 0; x <= 4000000; x++ {
			//fmt.Printf("x: %d, y: %d\n", x, y)
			bubbleFound := false
			for _, b := range bubbles {
				if inBubble(b, coord{x,y}) {
					bubbleFound = true
					x = bubbleEastEdgeAtY(b, y) // skip to the other edge of this bubble
					break
				}
			}
			if !bubbleFound {
				distressCoord = &coord{x,y}
				break
			}
		}
		if distressCoord != nil {
			break
		}
	}

	fmt.Printf("distressCoord: (%d, %d); frequency: %d\n", distressCoord.x, distressCoord.y, distressCoord.x*4000000 + distressCoord.y)
}

func addSensor(grid map[int]map[int]*loc, bubbles *[]bubble, sensor, beacon coord) {
	if _, ok := grid[sensor.y]; !ok {
		grid[sensor.y] = map[int]*loc{}
	}
	grid[sensor.y][sensor.x] = &loc{sensor: true}

	if _, ok := grid[beacon.y]; !ok {
		grid[beacon.y] = map[int]*loc{}
	}
	grid[beacon.y][beacon.x] = &loc{beacon: true}

	dist := distance(sensor, beacon)

	*bubbles = append(*bubbles, bubble{center: sensor, radius: dist})

	if dist < int(math.Abs(float64(sensor.y - 2000000))) {
		// no chance this sensor can rule out a location in the desired row
		return
	}

	for xDiff := 0; xDiff <= dist - int(math.Abs(float64(2000000-sensor.y))); xDiff++ {
		if _, ok := grid[2000000][sensor.x+xDiff]; !ok {
			grid[2000000][sensor.x+xDiff] = &loc{noBeacon: true}
		}

		if _, ok := grid[2000000][sensor.x-xDiff]; !ok {
			grid[2000000][sensor.x-xDiff] = &loc{noBeacon: true}
		}
	}
}

func distance(a, b coord) int {
	return int(math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y)))
}

func inBubble(b bubble, c coord) bool {
	return distance(b.center, c) <= b.radius
}

func bubbleEastEdgeAtY(b bubble, y int) int {
	return b.center.x + (b.radius-int(math.Abs(float64(y-b.center.y))))
}