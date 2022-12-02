package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	elves := [][]int{}

	for elf := getNextElf(scanner); len(elf) > 0; elf = getNextElf(scanner){
		elves = append(elves, elf)
	}

	max := 0
	for _, elf := range elves {
		if cals := calories(elf); cals > max {
			max = cals
		}
	}

	fmt.Printf("the most calories an elf has is: %d\n", max)

	sort.Slice(elves, func(i, j int) bool {
		return calories(elves[i]) > calories(elves[j])
	})

	fmt.Printf("the total of the top three elves is: %d\n", calories(elves[0])+calories(elves[1])+calories(elves[2]))
}

func getNextElf(scanner *bufio.Scanner) []int {
	elf := []int{}
	for scanner.Scan() {
		var cur string
		cur = scanner.Text();
		if cur == "" {
			return elf
		}
		i, err := strconv.Atoi(cur)
		if err != nil {
			log.Fatalf("%v", err)
		}
		elf = append(elf, i)
	}
	return elf
}

func calories(elf []int) int {
	total := 0
	for _, cals := range elf {
		total += cals
	}
	return total
}