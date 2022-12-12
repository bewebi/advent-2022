package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type monkey struct {
	number int
	items *queue

	operationArg1, operationArg2 *big.Int
	operationFunc func(*big.Int, *big.Int) *big.Int

	testNum int
	trueMonkey, falseMonkey int

	inspections int
}

type queue struct {
	elems []*big.Int
}

func (q *queue) push(n *big.Int) {
	q.elems = append([]*big.Int{n}, q.elems...)
}

func (q *queue) pop() *big.Int {
	n := q.elems[len(q.elems)-1]
	q.elems = q.elems[:(len(q.elems)-1)]
	return n
}

func (q *queue) empty() bool {
	return len(q.elems) == 0
}

// Tried using big package to no avail
// Got hit about LCM from reddit
func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	monkeys := []*monkey{}

	for {
		monkey, done := readMonkey(scanner)
		if done {
			break
		}
		monkeys = append(monkeys, monkey)
	}

	//fmt.Printf("have monkeys:\n")
	//for _, m := range monkeys {
	//	fmt.Printf("> %+v\n", m)
	//}

	lcm := 1
	for _, m := range monkeys {
		lcm *= m.testNum
	}
	bigLcm := big.NewInt(int64(lcm))

	for i := 0; i < 10000; i++ {
		//if i % 10 == 0 {
		//	fmt.Printf("Round %d\n---\n",i)
		//	for _, m := range monkeys {
		//		fmt.Printf("Monkey #%d: %d inspections\n", m.number, m.inspections)
		//	}
		//	fmt.Printf("\n")
		//}
		for _, m := range monkeys {
			for {
				if m.items.empty() {
					break
				}
				item := m.items.pop()
				arg1, arg2 := m.operationArg1, m.operationArg2
				if arg1.Sign() == -1 {
					arg1 = item
				}
				if arg2.Sign() == -1 {
					arg2 = item
				}
				item = m.operationFunc(arg1, arg2)
				item.Mod(item, bigLcm)

				//item = item / 3
				if big.NewInt(0).Mod(item, big.NewInt(int64(m.testNum))).Cmp(big.NewInt(0)) == 0 {
					monkeys[m.trueMonkey].items.push(item)
				} else {
					monkeys[m.falseMonkey].items.push(item)
				}

				m.inspections++
			}
		}
	}

	active1, active2 := &monkey{}, &monkey{}
	for _, m := range monkeys {
		fmt.Printf("Monkey #%d: %d inspections\n", m.number, m.inspections)
		if m.inspections > active1.inspections {
			active2 = active1
			active1 = m
		} else if m.inspections > active2.inspections {
			active2 = m
		}
	}

	fmt.Printf("most active are monkey #%d (%d inspections) and #%d (%d inspections)\nmonkey business: %d\n",
		active1.number, active1.inspections, active2.number, active2.inspections, active1.inspections*active2.inspections)
}

func readMonkey(scanner *bufio.Scanner) (*monkey, bool) {
	fmt.Println("reading monkey")
	m := &monkey{
		items: &queue{
			elems: []*big.Int{},
		},
	}

	// Monkey <n>:
	done := !scanner.Scan()
	if done {
		return nil, true
	}
	monkeyNumStr := string(strings.Split(scanner.Text(), " ")[1][0])
	monkeyNum, _ := strconv.Atoi(monkeyNumStr)
	m.number = monkeyNum

	//   Starting items: x, y, ...
	done = !scanner.Scan()
	if done {
		return nil, true
	}
	startingItemsElems := strings.Split(strings.Split(scanner.Text(), "Starting items: ")[1], ", ")
	for _, siStr := range startingItemsElems {
		si, _ := strconv.Atoi(siStr)
		siBig := big.NewInt(int64(si))
		m.items.push(siBig)
	}

	//   Operation: new = <x/old> <+-*/> <y/old>
	done = !scanner.Scan()
	if done {
		return nil, true
	}
	opElemStrs := strings.Split(strings.Split(scanner.Text(), "Operation: new = ")[1], " ")
	if opElemStrs[0] == "old" {
		m.operationArg1 = big.NewInt(-1)
	} else  {
		n, _ := strconv.Atoi(opElemStrs[0])
		bigN := big.NewInt(int64(n))
		m.operationArg1 = bigN
	}
	if opElemStrs[2] == "old" {
		m.operationArg2 = big.NewInt(-1)
	} else  {
		n, _ := strconv.Atoi(opElemStrs[2])
		bigN := big.NewInt(int64(n))
		m.operationArg2 = bigN
	}

	switch opElemStrs[1] {
	case "+":
		m.operationFunc = func(x, y *big.Int) *big.Int {return big.NewInt(0).Add(x, y)}
	case "-":
		m.operationFunc = func(x, y *big.Int) *big.Int {return big.NewInt(0).Sub(x, y)}
	case "*":
		m.operationFunc = func(x, y *big.Int) *big.Int {return big.NewInt(0).Mul(x, y)}
	case "/":
		m.operationFunc = func(x, y *big.Int) *big.Int {return big.NewInt(0).Div(x, y)}
	}

	//  Test: divisible by n
	done = !scanner.Scan()
	if done {
		return nil, true
	}
	testNumStr := strings.Split(scanner.Text(), "Test: divisible by ")[1]
	testNum, _ := strconv.Atoi(testNumStr)
	m.testNum = testNum

	//    If true: throw to monkey n
	done = !scanner.Scan()
	if done {
		return nil, true
	}
	trueMonkeyStr := strings.Split(scanner.Text(), "If true: throw to monkey ")[1]
	trueMonkey, _ := strconv.Atoi(trueMonkeyStr)
	m.trueMonkey = trueMonkey

	//    If false: throw to monkey n
	done = !scanner.Scan()
	if done {
		return nil, true
	}
	falseMonkeyStr := strings.Split(scanner.Text(), "If false: throw to monkey ")[1]
	falseMonkey, _ := strconv.Atoi(falseMonkeyStr)
	m.falseMonkey = falseMonkey

	// newline
	scanner.Scan()

	return m, false
}
