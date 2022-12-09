package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	firstKnotGrid, tailGrid := map[int]map[int]bool{}, map[int]map[int]bool{}

	knots := make([]*coordinate, 10) // 0 is head, 1 is first knot, 9 is tail
	for i := 0; i < 10; i ++ {
		knots[i] = &coordinate{}
	}

	for scanner.Scan() {
		directionElems := strings.Split(scanner.Text(), " ")
		num, _ := strconv.Atoi(directionElems[1])

		for i := 0; i < num; i++ {
			moveHeadOnce(knots[0], directionElems[0])
			for k := 0; k < 9; k++ {
				// move each successive knot
				moveTailOnce(knots[k], knots[k+1])
			}

			// part one
			if _, ok := firstKnotGrid[knots[1].x]; !ok {
				firstKnotGrid[knots[1].x] = map[int]bool{}
			}
			firstKnotGrid[knots[1].x][knots[1].y] = true

			// part 2
			if _, ok := tailGrid[knots[9].x]; !ok {
				tailGrid[knots[9].x] = map[int]bool{}
			}
			tailGrid[knots[9].x][knots[9].y] = true
		}
	}

	firstKnotVisitedCnt := 0
	for _, row := range firstKnotGrid {
		for _ = range row {
			firstKnotVisitedCnt++
		}
	}

	tailVisitedCnt := 0
	for _, row := range tailGrid {
		for _ = range row {
			tailVisitedCnt++
		}
	}

	fmt.Printf("the first knot visited %d different locations\n", firstKnotVisitedCnt)
	fmt.Printf("the tail visited %d different locations\n", tailVisitedCnt)

}


type coordinate struct {
	x, y int
}

func moveHeadOnce(head *coordinate, dir string) {
	switch dir {
	case "U":
		head.y++
	case "R":
		head.x++
	case "D":
		head.y--
	case "L":
		head.x--
	}
}

func moveTailOnce(head, tail *coordinate) {
	if math.Abs(float64(head.x-tail.x)) <= 1 && math.Abs(float64(head.y-tail.y)) <= 1 {
		return // already touching
	}

	if head.x == tail.x {
		// same row, must need to move up or down
		if head.y > tail.y {
			tail.y++
		} else {
			tail.y--
		}
		return
	}

	if head.y == tail.y {
		// same column, must need to move left or right
		if head.x > tail.x {
			tail.x++
		} else {
			tail.x--
		}
	}

	// need a diagonal move
	if (head.x - tail.x) == 2 {
		tail.x++
		if head.y > tail.y {
			tail.y++
		} else {
			tail.y--
		}
		return
	}

	if (tail.x - head.x) == 2 {
		tail.x--
		if head.y > tail.y {
			tail.y++
		} else {
			tail.y--
		}
		return
	}

	if (head.y - tail.y) == 2 {
		tail.y++
		if head.x > tail.x {
			tail.x++
		} else {
			tail.x--
		}
		return
	}

	if (tail.y - head.y) == 2 {
		tail.y--
		if head.x > tail.x {
			tail.x++
		} else {
			tail.x--
		}
		return
	}
}
