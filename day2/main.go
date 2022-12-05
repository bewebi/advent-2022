package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalScore := 0
	round := 0
	for scanner.Scan() {
		roundStr := scanner.Text()
		roundArr := strings.Split(roundStr, " ")
		roundScore := scoreRoundPhaseOne(roundArr[0], roundArr[1])
		fmt.Printf("round %d score: %d\n", round, roundScore)
		round++
		totalScore += roundScore
	}

	fmt.Printf("the total score for phase 1 is: %d\n", totalScore)

	file2, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file2.Close()
	scanner = bufio.NewScanner(file2)
	totalScore = 0
	round = 0
	for scanner.Scan() {
		roundStr := scanner.Text()
		roundArr := strings.Split(roundStr, " ")
		roundScore := scoreRoundPhaseTwo(roundArr[0], roundArr[1])
		fmt.Printf("round %d score: %d\n", round, roundScore)
		round++
		totalScore += roundScore
	}

	fmt.Printf("the total score for phase 2 is: %d\n", totalScore)
}

func scoreRoundPhaseOne(opp, us string) int {
	score := 0
	switch us {
	case "X":
		score += 1
	case "Y":
		score += 2
	case "Z":
		score += 3
	}

	if (us == "X" && opp == "A") || (us == "Y" && opp == "B") || (us == "Z" && opp == "C") {
		score += 3
	} else if (us == "X" && opp == "C") || (us == "Y" && opp == "A") || (us == "Z" && opp == "B") {
		score += 6
	}

	return score
}

func scoreRoundPhaseTwo(opp, us string) int {
	score := 0
	switch us {
	case "X": // lose
		switch opp {
		case "A":
			return 3 + 0 // scissors, loss
		case "B":
			return 1 + 0 // rock, loss
		case "C":
			return 2 + 0 // paper, loss
		}
	case "Y": // draw
		switch opp {
		case "A":
			return 1 + 3 // rock, draw
		case "B":
			return 2 + 3 // paper, draw
		case "C":
			return 3 + 3 // scissors, draw
		}
	case "Z": // win
		switch opp {
		case "A":
			return 2 + 6 // paper, win
		case "B":
			return 3 + 6 // scissors, win
		case "C":
			return 1 + 6 // rock, win
		}
	}

	if (us == "X" && opp == "A") || (us == "Y" && opp == "B") || (us == "Z" && opp == "C") {
		score += 3
	} else if (us == "X" && opp == "C") || (us == "Y" && opp == "A") || (us == "Z" && opp == "B") {
		score += 6
	}

	return score
}