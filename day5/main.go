package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type stack []string

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var stacks []*stack

	for scanner.Scan() {
		crateStr := scanner.Text()
		expStacks := ((len(crateStr)-3)/4) + 1
		if len(stacks) < expStacks {
			prevLen := len(stacks)
			fmt.Printf("making %d expected stacks\n", expStacks)
			for i := prevLen; i < expStacks; i++ {
				stacks = append(stacks, &stack{})
			}
		}
		if crateStr == "" {
			break
		}

		for i := 0; i < len(stacks); i++ {
			ind := 1 + 4*i
			if crateStr[ind] != ' ' {
				stacks[i].add(string(crateStr[ind]))
			}
		}
	}
	for _, stack := range stacks {
		stack.remove() // remove number
		stack.reverse() // first item goes on top
	}

	fmt.Printf("have stacks:\n")
	for i, stack := range stacks {
		fmt.Printf("%d: %+v\n", i+1, *stack)
	}

	for scanner.Scan() {
		moveStr := scanner.Text()
		fmt.Println(moveStr)
		moveElems := strings.Split(moveStr, " ")
		n, _ := strconv.Atoi(moveElems[1])
		src, _ := strconv.Atoi(moveElems[3])
		dst, _ := strconv.Atoi(moveElems[5])

		//stacks[src-1].moveNToOther(n, stacks[dst-1])
		stacks[src-1].moveNToOtherInOrder(n, stacks[dst-1])

		fmt.Printf("have stacks:\n")
		for i, stack := range stacks {
			fmt.Printf("%d: %+v\n", i+1, *stack)
		}
	}

	output := ""
	for _, stack := range stacks {
		output += stack.remove()
	}
	fmt.Printf("final output: %s\n", output)
}

func (s *stack) add(elem string) {
	*s = append(*s, elem)
}

func (s *stack) remove() string {
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element
}

func (s *stack) moveNToOther(n int, other *stack) {
	for i := 0; i < n; i++ {
		tmp := s.remove()
		other.add(tmp)
	}
}

func (s *stack) moveNToOtherInOrder(n int, other *stack) {
	tmp := make([]string, n)
	for i := n-1; i >= 0; i-- {
		tmp[i] = s.remove()
	}
	for _, elem := range tmp {
		other.add(elem)
	}
}

func (s *stack) reverse() {
	for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}