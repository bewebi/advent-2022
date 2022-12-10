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

	strength := 1
	strengthHistory := []int{}

	for scanner.Scan() {
		instStr := scanner.Text()

		strengthHistory = append(strengthHistory, strength)

		switch instStr{
		case "noop":
		default:
			// extra step to process
			strengthHistory = append(strengthHistory, strength)

			diff, _ := strconv.Atoi(strings.Split(instStr, " ")[1])
			strength += diff
		}
	}

	strengthSum := 0
	for i := 19; i < len(strengthHistory); i+=40 {
		fmt.Printf("cycle %d: strength = %d, adding %d\n", i+1, strengthHistory[i], (i+1)*strengthHistory[i])
		strengthSum += (i+1)*strengthHistory[i]
	}

	fmt.Printf("strengthSum: %d\n", strengthSum)

	var pixelList []string
	for i, v := range strengthHistory {
		//fmt.Printf("cycle: %d; pixel position: %d; sprite position: %d; writing: ", i+1, i%40, v)
		if math.Abs(float64(v - (i%40))) < 2 {
			//fmt.Printf("%s\n", "#")
			pixelList = append(pixelList, "#")
		}  else {
			//fmt.Printf("%s\n", ".")
			pixelList = append(pixelList, ".")
		}
	}

	pixelOutput := ""
	for i, p := range pixelList {
		pixelOutput += p
		if (i+1)%40 == 0 {
			pixelOutput += "\n"
		}
	}

	fmt.Printf("pixel output:\n%s\n", pixelOutput)
}

