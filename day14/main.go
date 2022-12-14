package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)


//go:embed input.txt
var input string

var maxDepth, minX, maxX int

type loc struct {
	rock, sand bool
}

type coord struct {
	x, y int
}

func main() {
	grid := map[int]map[int]*loc{}

	lines := strings.Split(input, "\n")
	for i := 0; i < len(lines); i++ {
		coordStrs := strings.Split(lines[i], " -> ")
		prevCoord := coord{}
		for i, coordStr := range coordStrs {
			coordElems := strings.Split(coordStr, ",")
			x, _ := strconv.Atoi(coordElems[0])
			y, _ := strconv.Atoi(coordElems[1])
			newCoord := coord{x,y}

			if _, ok := grid[y]; !ok {
				grid[y] = map[int]*loc{}
			}
			if y > maxDepth {
				maxDepth = y
			}
			if minX == 0 || x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}

			if i != 0 {
				if newCoord.x > prevCoord.x {
					for n := prevCoord.x+1; n <= newCoord.x; n++ {
						grid[y][n] = &loc{rock: true}
					}
				}
				if newCoord.x < prevCoord.x {
					for n := newCoord.x; n < prevCoord.x; n++ {
						grid[y][n] = &loc{rock: true}
					}
				}
				if newCoord.y > prevCoord.y {
					for n := prevCoord.y+1; n <= newCoord.y; n++ {
						if _, ok := grid[n]; !ok {
							grid[n] = map[int]*loc{}
						}
						grid[n][x] = &loc{rock: true}
					}
				}
				if newCoord.y < prevCoord.y {
					for n := newCoord.y; n < prevCoord.y; n++ {
						if _, ok := grid[n]; !ok {
							grid[n] = map[int]*loc{}
						}
						grid[n][x] = &loc{rock: true}
					}
				}
			}

			grid[y][x] = &loc{rock: true}
			prevCoord = newCoord
		}
	}

	fmt.Printf("Have grid:\n")
	printGrid(grid)

	for {
		placed := addSandForever(grid, 500, 0)
		if !placed {
			break
		}
	}

	fmt.Printf("Poured sand into infinity; new grid:\n")
	printGrid(grid)

	sandCount := 0
	for y := 0; y <= maxDepth; y++ {
		for x := minX; x <= maxX; x++ {
			if loc, ok := grid[y][x]; ok && loc.sand {
				sandCount++
			}
		}
	}

	fmt.Printf("Total sand poured: %d\n", sandCount)

	eraseSand(grid)

	for {
		reachedTop := addSandWithFloorToTop(grid, 500, 0, maxDepth+2, 500, 0)
		if !reachedTop {
			break
		}
	}

	fmt.Printf("Poured sand to the top; new grid:\n")
	printGrid(grid)

	sandCount = 0
	for y := 0; y <= maxDepth+2; y++ {
		for x := minX; x <= maxX; x++ {
			if loc, ok := grid[y][x]; ok && loc.sand {
				sandCount++
			}
		}
	}

	fmt.Printf("Total sand poured to the top: %d\n", sandCount)


}

func addSandForever(grid map[int]map[int]*loc, x, y int) bool  {
	//fmt.Printf("addSandForever (%d,%d)\n", x, y)
	nextRow, ok := grid[y+1]
	if !ok {
		if y+1 > maxDepth {
			//fmt.Printf("next row %d is > maxDepth %d\n", y+1, maxDepth)
			// next row doesn't exist, sand will fall forever
			return false
		}
		// add the empty row and proceed
		//fmt.Printf("Adding empty row at %d\n", y+1)
		grid[y+1] = map[int]*loc{}
	}
	if loc, ok := nextRow[x]; !ok || (!loc.sand && !loc.rock) {
		// if coord doesn't exist, it's open
		//fmt.Printf("moving sand down one to (%d, %d)\n", x, y+1)
		return addSandForever(grid, x,  y+1)
	}
	if loc, ok := nextRow[x-1]; !ok || (!loc.sand && !loc.rock)  {
		//fmt.Printf("moving sand down one and left to (%d, %d)\n", x-1, y+1)
		if x - 1 < minX {
			minX = x - 1
		}
		return addSandForever(grid, x-1, y+1)
	}
	if loc, ok := nextRow[x+1]; !ok || (!loc.sand && !loc.rock) {
		//fmt.Printf("moving sand down one and right to (%d, %d)\n", x+1, y+1)
		if x + 1 > maxX {
			maxX = x + 1
		}
		return addSandForever(grid, x+1, y+1)
	}

	// none of the remaining spots are open, add sand here
	//fmt.Printf("sand has come to a rest at (%d,%d)\n", x, y)
	grid[y][x] = &loc{sand: true}
	return true
}

func addSandWithFloorToTop(grid map[int]map[int]*loc, x, y, floor, topX, topY int) bool  {
	//fmt.Printf("addSandWithFloorToTop (%d,%d)\n", x, y)

	if y+1 == floor {
		// nowhere to go, rest here
		//fmt.Printf("sand has come to a rest on the floor at (%d,%d)\n", x, y)
		grid[y][x] = &loc{sand: true}
		return true
	}

	nextRow, ok := grid[y+1]
	if !ok {
		// add the empty row and proceed
		//fmt.Printf("Adding empty row at %d\n", y+1)
		grid[y+1] = map[int]*loc{}
	}

	if loc, ok := nextRow[x]; !ok || (!loc.sand && !loc.rock) {
		// if coord doesn't exist, it's open
		//fmt.Printf("moving sand down one to (%d, %d)\n", x, y+1)
		return addSandWithFloorToTop(grid, x, y+1, floor, topX, topY)
	}
	if loc, ok := nextRow[x-1]; !ok || (!loc.sand && !loc.rock) {
		//fmt.Printf("moving sand down one and left to (%d, %d)\n", x-1, y+1)
		if x - 1 < minX {
			minX = x - 1
		}
		return addSandWithFloorToTop(grid, x-1, y+1, floor, topX, topY)
	}
	if loc, ok := nextRow[x+1]; !ok || (!loc.sand && !loc.rock) {
		//fmt.Printf("moving sand down one and right to (%d, %d)\n", x+1, y+1)
		if x + 1 > maxX {
			maxX = x + 1
		}
		return addSandWithFloorToTop(grid, x+1, y+1, floor, topX, topY)
	}

	// none of the remaining spots are open, add sand here
	//fmt.Printf("sand has come to a rest at (%d,%d)\n", x, y)
	grid[y][x] = &loc{sand: true}

	if x == topX && y == topY {
		//fmt.Printf("sand has come to a rest at the starting point")
		return false
	}
	return true
}

func printGrid(grid map[int]map[int]*loc) {
	for y := 0; y <= maxDepth+2; y++ {
		if _, ok := grid[y]; !ok {
			grid[y] = map[int]*loc{}
		}
		fmt.Printf("%3d: ", y)
		for x := minX; x <= maxX; x++ {
			if l, ok := grid[y][x]; ok {
				if l.rock {
					fmt.Printf("# ")
				} else if l.sand {
					fmt.Printf("o ")
				} else {
					fmt.Printf(". ")
				}
			} else {
				fmt.Printf(". ")
			}
		}
		fmt.Printf("\n")
	}
}

func eraseSand(grid map[int]map[int]*loc) {
	for _, row := range grid {
		for _, l := range row {
			l.sand = false
		}
	}
}