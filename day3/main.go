package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalPriority := 0
	rucksack := 0
	for scanner.Scan() {
		rucksackStr := scanner.Text()
		rucksackPriority := prioritizeRucksack(rucksackStr)
		fmt.Printf("rucksack %d priority: %d\n", rucksack, rucksackPriority)
		rucksack++
		totalPriority += rucksackPriority
	}

	fmt.Printf("the total priority for phase 1 is: %d\n", totalPriority)

	file2, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file2.Close()

	scanner = bufio.NewScanner(file2)
	totalPriority = 0
	rucksack = 0
	curGroup := []string{}
	groupNum := 0
	for scanner.Scan() {
		rucksackStr := scanner.Text()
		curGroup = append(curGroup, rucksackStr)
		rucksack++

		if len(curGroup) == 3 {
			groupPriority := getBadge(curGroup[0], curGroup[1], curGroup[2])
			fmt.Printf("group %d priority: %d\n", groupNum, groupPriority)

			totalPriority += groupPriority
			groupNum++
			curGroup = []string{}
		}
	}

	fmt.Printf("the total priority for phase 2 is: %d\n", totalPriority)
}

func prioritizeRucksack(rucksackStr string) int {
	compartmentOne := map[uint8]bool{}
	for i := 0; i < len(rucksackStr)/2; i++ {
		compartmentOne[rucksackStr[i]] = true
	}
	for i := len(rucksackStr)/2; i < len(rucksackStr); i++ {
		elem := rucksackStr[i]
		if compartmentOne[elem] {
			return byteToPriority(elem)
		}
	}
	return -1
}

func getBadge(rucksack1, rucksack2, rucksack3 string) int {
	itemCounts := map[uint8]int{}
	for i := 0; i < len(rucksack1); i++ {
		itemCounts[rucksack1[i]] = 1
	}

	for i := 0; i < len(rucksack2); i++ {
		if _, ok := itemCounts[rucksack2[i]]; ok {
			itemCounts[rucksack2[i]] = 2
		}
	}

	for i := 0; i < len(rucksack3); i++ {
		elem := rucksack3[i]
		if itemCounts[elem] == 2 {
			return byteToPriority(elem)
		}
	}
	return -1
}

func byteToPriority(elem uint8) int {
	if elem >= 'a' && elem <= 'z' {
		return int(elem-'a'+1)
	} else if elem >= 'A' && elem <= 'Z' {
		return int(elem-'A'+27)
	}
	return -1
}

