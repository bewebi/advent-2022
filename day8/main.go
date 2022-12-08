package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type tree struct {
	height int
	visible bool
	score int
}

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	grid := map[int]map[int]*tree{}
	row := 0

	for scanner.Scan() {
		grid[row] = map[int]*tree{}
		rowStr := scanner.Text()
		for i, hByte := range rowStr {
			h, _ := strconv.Atoi(string(hByte))
			grid[row][i] = &tree{height: h}
		}
		row++
	}

	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			setCoordsVisibleAndScore(grid, r, c)
		}
	}

	visibleCount := 0
	maxScore := 0
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c].visible {
				visibleCount++
			}
			if grid[r][c].score > maxScore {
				maxScore = grid[r][c].score
			}
		}
	}

	fmt.Printf("total visible trees: %d\n", visibleCount)
	fmt.Printf("max score: %d\n", maxScore)

	setCoordsVisibleAndScore(grid, 3, 2)
}

func setCoordsVisibleAndScore(grid map[int]map[int]*tree, row, col int) {
	tree := grid[row][col]
	fmt.Printf("checking (%d,%d), height: %d\n", row, col, tree.height)
	if row == 0 || col == 0 || row == len(grid) - 1 || col == len(grid[0]) - 1 {
		fmt.Println("edge")
		tree.visible = true
		return
	}

	n, e, s, w := true, true, true, true
	nScore, eScore, sScore, wScore := 0, 0, 0, 0

	// visible from the north
	for i := row-1; i >= 0; i-- {
		nScore++

		if grid[i][col].height >= tree.height {
			n = false
			break
		}
	}
	if n {
		fmt.Println("visible from north")
		tree.visible = true
	}

	// visible from the east
	for i := col+1; i < len(grid[row]); i++ {
		eScore++

		if grid[row][i].height >= tree.height {
			e = false
			break
		}
	}
	if e {
		fmt.Println("visible from east")
		tree.visible = true
	}

	// visible from the south
	for i := row + 1; i < len(grid); i ++ {
		sScore++

		if grid[i][col].height >= tree.height {
			s = false
			break
		}
	}
	if s {
		fmt.Println("visible from south")
		tree.visible = true
	}

	// visible from the west
	for i := col - 1; i >= 0; i-- {
		wScore++

		if grid[row][i].height >= tree.height {
			w = false
			break
		}
	}
	if w {
		fmt.Println("visible from west")
		tree.visible = true
	} else {
		fmt.Println("not visible")
	}

	fmt.Printf("scores: n: %d, e: %d, s: %d, w: %d, total: %d\n", nScore, eScore, sScore, wScore, nScore*eScore*sScore*wScore)
	tree.score = nScore*eScore*sScore*wScore
}