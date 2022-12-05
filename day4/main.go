package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type assignment struct {
	hi, lo int
}

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalEclipses := 0
	totalOverlaps := 0
	pair := 0
	for scanner.Scan() {
		pairStr := scanner.Text()
		ass1, ass2 := parseRawStr(pairStr)
		eclipses := containsEclipse(ass1, ass2)
		overlaps := containsOverlap(ass1, ass2)
		fmt.Printf("pair %d eclipses: %v; overlaps: %v\n", pair, eclipses, overlaps)
		pair++
		if eclipses {
			totalEclipses++
		}
		if overlaps {
			totalOverlaps++
		}
	}

	fmt.Printf("the total eclipsing pairs is: %d\nthe total overlapping pairs is: %d\n", totalEclipses, totalOverlaps)
}

func parseRawStr(rawStr string) (assignment, assignment) {
	var ass1, ass2 assignment
	pairElems := strings.Split(rawStr, ",")
	ass1Elems := strings.Split(pairElems[0], "-")
	ass1.lo, _ = strconv.Atoi(ass1Elems[0])
	ass1.hi, _ = strconv.Atoi(ass1Elems[1])
	ass2Elems := strings.Split(pairElems[1], "-")
	ass2.lo, _ = strconv.Atoi(ass2Elems[0])
	ass2.hi, _ = strconv.Atoi(ass2Elems[1])
	return ass1, ass2
}

func containsEclipse(ass1, ass2 assignment) bool {
	return  (ass1.lo <= ass2.lo && ass1.hi >= ass2.hi) ||
		(ass2.lo <= ass1.lo && ass2.hi >= ass1.hi)
}

func containsOverlap(ass1, ass2 assignment) bool {
	return containsEclipse(ass1, ass2) ||
		(ass1.lo >= ass2.lo && ass1.lo <= ass2.hi) ||
		(ass2.lo >= ass1.lo && ass2.lo <= ass1.hi)
}