package main

import (
	_ "embed"
	"fmt"
	"strings"
)


//go:embed input.txt
var input string

type coordinate struct {
	row, col int
}

type square struct {
	elev int
	steps int
	coord coordinate
	visited bool
}

type queue struct {
	elems []*square
}

func (q *queue) push(s *square) {
	q.elems = append([]*square{s}, q.elems...)
}

func (q *queue) pop() *square {
	n := q.elems[len(q.elems)-1]
	q.elems = q.elems[:(len(q.elems)-1)]
	return n
}

func (q *queue) empty() bool {
	return len(q.elems) == 0
}

func main() {
	grid := map[int]map[int]*square{}
	var start, end coordinate

	lines := strings.Split(input, "\n")
	for row, l := range lines {
		grid[row] = map[int]*square{}
		for col, b := range l {
			s := &square{
				coord: coordinate{row, col},
			}
			switch b {
			case 'S':
				start = coordinate{row: row, col: col}
			case 'E':
				s.elev = 26
				end = coordinate{row: row, col: col}
			default:
				s.elev = int(b) - 96
			}

			grid[row][col] = s
		}
	}

	fmt.Printf("It takes %d steps to get from the original start to the end\n", minStepsFrom(start, end, grid))
	eraseSteps(grid)

	starts := []coordinate{}
	for _, row := range grid {
		for _, sq := range row {
			if sq.elev <= 1 {
				starts = append(starts, sq.coord)
			}
		}
	}

	//fmt.Printf("Potential starting points: %+v\n", starts)

	minSteps := -1
	for _, st := range starts {
		steps := minStepsFrom(st, end, grid)
		if steps == -1 {
			continue
		}

		if minSteps == -1 || minSteps > steps {
			//fmt.Printf("It takes %d steps to get from (%d, %d) to the end\n", steps, st.row, st.col)
			minSteps = steps
		}
		eraseSteps(grid)
	}

	fmt.Printf("The minimum steps from any start is: %d\n", minSteps)
}

func (sq *square) canMoveTo(g map[int]map[int]*square, r, c int) bool {
	return r >= 0 && r < len(g) &&
		c >= 0 && c < len(g[0]) &&
		!g[r][c].visited &&
		g[r][c].elev <= sq.elev+1
}

func minStepsFrom(start, end coordinate, grid map[int]map[int]*square) int {
	moves := queue{
		elems: []*square{},
	}
	moves.push(grid[start.row][start.col])

	for !moves.empty() {
		sq := moves.pop()
		if sq.visited {
			continue
		}
		sq.visited = true
		//fmt.Printf("At (%d, %d)\n", sq.coord.row, sq.coord.col)

		if sq.canMoveTo(grid, sq.coord.row-1, sq.coord.col) {
			//fmt.Printf("can move up to (%d, %d)\n", sq.coord.row-1, sq.coord.col)
			next := grid[sq.coord.row-1][sq.coord.col]
			next.steps = sq.steps + 1
			if next.coord == end {
				next.visited = true
				break
			}
			moves.push(next)
		}
		if sq.canMoveTo(grid, sq.coord.row+1, sq.coord.col) {
			//fmt.Printf("can move down to (%d, %d)\n", sq.coord.row+1, sq.coord.col)
			next := grid[sq.coord.row+1][sq.coord.col]
			next.steps = sq.steps + 1
			if next.coord == end {
				next.visited = true
				break
			}
			moves.push(next)
		}
		if sq.canMoveTo(grid, sq.coord.row, sq.coord.col-1) {
			//fmt.Printf("can move left to (%d, %d)\n", sq.coord.row, sq.coord.col-1)
			next := grid[sq.coord.row][sq.coord.col-1]
			next.steps = sq.steps + 1
			if next.coord == end {
				next.visited = true
				break
			}
			moves.push(next)
		}
		if sq.canMoveTo(grid, sq.coord.row, sq.coord.col+1) {
			//fmt.Printf("can move right to (%d, %d)\n", sq.coord.row, sq.coord.col+1)
			next := grid[sq.coord.row][sq.coord.col+1]
			next.steps = sq.steps + 1
			if next.coord == end {
				next.visited = true
				break
			}
			moves.push(next)
		}
	}

	if grid[end.row][end.col].visited {
		return grid[end.row][end.col].steps
	} else {
		return -1
	}
}

func eraseSteps(grid map[int]map[int]*square) {
	for _, row := range grid {
		for _, sq := range row {
			sq.steps = 0
			sq.visited = false
		}
	}
}