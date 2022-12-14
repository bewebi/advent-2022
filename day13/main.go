package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)


//go:embed input.txt
var input string

func main() {
	packets := [][]interface{}{}

	lines := strings.Split(input, "\n")
	for i := 0; i < len(lines); i +=3 {
		packets = append(packets, parseLine(lines[i]), parseLine(lines[i+1]))
	}

	indexSum := 0
	for i := 0; i < len(packets); i += 2 {
		if ordered, equal := compare(packets[i], packets[i+1]); ordered && !equal {
			indexSum += i+1
		} else if equal {
			fmt.Printf("unexpected equal pair at %d: %+v\n", i+1, packets[i])
		}
	}

	fmt.Printf("Sum of ordered indexes is: %d\n", indexSum)

	divOne, divTwo := []interface{}{[]interface{}{float64(2)}}, []interface{}{[]interface{}{float64(6)}}
	packets = append(packets, divOne, divTwo)

	sort.Slice(packets, func(i, j int) bool {
		ordered, _ := compare(packets[i], packets[j])
		return ordered
	})

	divOneIndex, divTwoIndex := 0, 0
	for i := 0; i < len(packets); i++ {
		if _, equal := compare(packets[i], divOne); equal {
			divOneIndex = i+1
		}
		if _, equal := compare(packets[i], divTwo); equal {
			divTwoIndex = i+1
			break
		}
	}

	fmt.Printf("Seperator packets found at %d and %d; product: %d\n", divOneIndex, divTwoIndex, divOneIndex*divTwoIndex)
}

func parseLine(line string) []interface{} {
	list := []interface{}{}
	err := json.Unmarshal([]byte(line), &list)
	if err != nil {
		//panic(err)
	}
	return list
}

func compare(left, right []interface{}) (bool, bool) {
	//fmt.Printf("comparing %+v and %+v\n", left, right)
	for i, lElem := range left {
		//fmt.Printf("Index %d: ", i)
		if i > len(right) - 1  {
			//fmt.Println("right has more elems")
			// right has more elems
			return false, false
		}
		rElem := right[i]
		switch lElem.(type) {
		case float64:
			if rInt, ok := rElem.(float64); ok {
				//fmt.Printf("both ints; ")
				lInt, _ := lElem.(float64)
				if lInt == rInt {
					//fmt.Printf("equal, skipping\n")
					continue
				}
				if lInt < rInt {
					//fmt.Printf("l < r, returning true\n")
					return true, false
				} else {
					//fmt.Printf("l > r, returning false\n")
					return false, false
				}
			}
			rList, _ := rElem.([]interface{})
			ret, equal := compare([]interface{}{lElem}, rList)
			if equal {
				continue
			}
			return ret, false
		case []interface{}:
			lList, _ := lElem.([]interface{})

			if rList, ok := rElem.([]interface{}); ok {
				ret, equal := compare(lList, rList)
				if equal {
					continue
				}
				return ret, false
			}
			ret, equal := compare(lList, []interface{}{rElem})
			if equal {
				continue
			}
			return ret, false
		default:
			//fmt.Printf("unexpected type: %T\n", lElem)
		}
	}
	if len(left) < len(right) {
		//fmt.Printf("left ran out of elems, returning true\n")
		// left ran out of items first
		return true, false
	}
	// elems were equal
	return false, true
}