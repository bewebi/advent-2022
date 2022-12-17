package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)


//go:embed input.txt
var input string

type valve struct {
	name string
	flow int
	children []string
	distances map[string]int
}

type valveSearch struct {
	v *valve
	dist int
}

type valveQueue struct {
	elems []*valveSearch
}

func (q *valveQueue) push(v *valveSearch) {
	q.elems = append([]*valveSearch{v}, q.elems...)
}

func (q *valveQueue) pop() *valveSearch {
	n := q.elems[len(q.elems)-1]
	q.elems = q.elems[:(len(q.elems)-1)]
	return n
}

func (q *valveQueue) empty() bool {
	return len(q.elems) == 0
}

type searchState struct {
	openValves map[string]int
	totalFlow int
	pTimer, eTimer int
	remainingMinutes int
	curValve, eCurValve *valve
	history []string
}

type searchQueue struct {
	elems []*searchState
}

func (q *searchQueue) push(ss *searchState) {
	q.elems = append([]*searchState{ss}, q.elems...)
}

func (q *searchQueue) pop() *searchState {
	n := q.elems[len(q.elems)-1]
	q.elems = q.elems[:(len(q.elems)-1)]
	return n
}

func (q *searchQueue) empty() bool {
	return len(q.elems) == 0
}

func main() {
	lineRE := `Valve ([A-Z]+) has flow rate=([0-9]+); tunnels? leads? to valves? ([A-Z, ]+)`
	re := regexp.MustCompile(lineRE)

	graph := map[string]*valve{}

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		flowRate, _ := strconv.Atoi(matches[2])

		graph[matches[1]] = &valve{name: matches[1], flow: flowRate, children: strings.Split(matches[3],", "), distances: map[string]int{}}
	}

	for _, v := range graph {
		addDistances(v, graph)
	}

	//fmt.Printf("have valves: %+v\n", graph)
	valvesWithFlowCount := 0
	for _, valve := range graph {
		if valve.flow > 0 {
			valvesWithFlowCount++
		}
	}

	maxFlow := 0

	/* Part 1 (takes a few minutes) */
	//searchQueue := searchQueue{
	//	elems: []*searchState{
	//		{
	//			openValves: map[string]bool{},
	//			remainingMinutes: 30,
	//			curValve: graph["AA"],
	//		},
	//	},
	//}
	//
	//remainingMinutes := 30
	//
	//for !searchQueue.empty() {
	//	state := searchQueue.pop()
	//	fmt.Printf("state: %+v\n", state)
	//
	//	if state.totalFlow > maxFlow {
	//		maxFlow = state.totalFlow
	//	}
	//
	//	if state.remainingMinutes == 0 || len(state.openValves) == valvesWithFlowCount {
	//		//fmt.Printf("out of time or all valves open\n")
	//		continue
	//	}
	//
	//	if state.remainingMinutes < remainingMinutes {
	//		remainingMinutes--
	//		fmt.Printf("remainingMinutes: %d\n", remainingMinutes)
	//	}
	//
	//	if _, ok := state.openValves[state.curValve.name]; !ok && state.curValve.flow > 0 {
	//		//fmt.Printf("haven't opened %s with flow %d, spending 1 minute to do so\n", state.curValve.name, state.curValve.flow)
	//		newOpenValves := copyMap(state.openValves)
	//		newOpenValves[state.curValve.name] = true
	//		searchQueue.push(&searchState{
	//			openValves: newOpenValves,
	//			totalFlow: state.totalFlow + ((state.remainingMinutes-1) * state.curValve.flow),
	//			remainingMinutes: state.remainingMinutes-1,
	//			curValve: state.curValve,
	//		})
	//		continue
	//	}
	//
	//	for other, dist := range state.curValve.distances {
	//		if _, ok := state.openValves[other]; !ok && graph[other].flow > 0 && state.remainingMinutes - dist > 0 {
	//			//fmt.Printf("haven't opened %s with flow %d, spending %d minute(s) to travel there\n", graph[other].name, graph[other].flow, dist)
	//			searchQueue.push(&searchState{
	//				openValves:       state.openValves,
	//				totalFlow:        state.totalFlow,
	//				remainingMinutes: state.remainingMinutes - dist,
	//				curValve:         graph[other],
	//			})
	//		}
	//	}
	//}

	/* Part 2 */
	// Obviously not the most optimized solution (takes 3-4 hours to run) but it does work!
	searchQueue := searchQueue{
		elems: []*searchState{
			{
				openValves: map[string]int{},
				remainingMinutes: 26,
				curValve: graph["AA"],
				eCurValve: graph["AA"],
			},
		},
	}

	//histories := [][]string{}

	visitedMapLock := sync.Mutex{}
	visitedMap := map[string]int{}

	minMinutesLeft := 26

	wg := &sync.WaitGroup{}
	wg.Add(50)

	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			for !searchQueue.empty() {
				state := searchQueue.pop()

				if state.remainingMinutes < minMinutesLeft || len(searchQueue.elems)%10000 == 0 {
					minMinutesLeft = state.remainingMinutes
					fmt.Printf("%d minutes remaining, queue length: %d\n", minMinutesLeft, len(searchQueue.elems))
				}

				visitedMapLock.Lock()
				visited := visitedMap[visitedKey(state.openValves)] > state.totalFlow
				visitedMapLock.Unlock()

				if state.remainingMinutes == 0 || len(state.openValves) == valvesWithFlowCount || visited {
					//histories = append(histories, state.history)
					continue
				}

				//deadEnd := true

				newOpenValves := copyMap(state.openValves)
				newTotalFlow := state.totalFlow
				newPTimer := state.pTimer - 1
				newCurValve := state.curValve
				newETimer := state.eTimer - 1
				newECurValve := state.eCurValve
				newRemainingMinutes := state.remainingMinutes - 1
				//historyEntry := fmt.Sprintf("minutes remaining: %d; human at %s with timer %d, elephant at %s with time %d\n",
				//	state.remainingMinutes, state.curValve.name, state.pTimer, state.eCurValve.name, state.eTimer)

				if _, ok := state.openValves[state.curValve.name]; state.pTimer == 0 && state.curValve.flow > 0 && !ok {
					// person opens valve
					newOpenValves[state.curValve.name] = newRemainingMinutes
					newTotalFlow += newRemainingMinutes * state.curValve.flow
					//historyEntry += fmt.Sprintf("human opens valve %s, adding %d flow per turn, %d lifetime, total: %d\n", state.curValve.name, state.curValve.flow, newRemainingMinutes * state.curValve.flow, newTotalFlow)
					newPTimer = 0
				}
				if _, ok := state.openValves[state.eCurValve.name]; state.eTimer == 0 && state.eCurValve.flow > 0 && !ok {
					// elephant opens valve

					newOpenValves[state.eCurValve.name] = newRemainingMinutes
					newTotalFlow += newRemainingMinutes * state.eCurValve.flow
					//historyEntry += fmt.Sprintf("elephant opens valve %s, adding %d flow per turn, %d lifetime, total: %d\n", state.eCurValve.name, state.eCurValve.flow, newRemainingMinutes * state.eCurValve.flow, newTotalFlow)
					newETimer = 0
				}

				visitedMapLock.Lock()
				bestForKey := visitedMap[visitedKey(newOpenValves)]
				visitedMapLock.Unlock()

				if bestForKey <= newTotalFlow {
					visitedMapLock.Lock()
					visitedMap[visitedKey(newOpenValves)] = newTotalFlow
					visitedMapLock.Unlock()
				} else {
					continue
				}

				if newTotalFlow > maxFlow {
					//fmt.Printf("new max with %d to go; history: \n%s\n", newRemainingMinutes, strings.Join(append(state.history, historyEntry), "\n"))
					maxFlow = newTotalFlow
				}

				if state.pTimer == 0 && (state.openValves[state.curValve.name] > 0 || state.curValve.flow == 0) {
					// person needs new dest
					//fmt.Printf("person needs new dest\n")
					newCurValve = nil
				}
				if state.eTimer == 0 && (state.openValves[state.eCurValve.name] > 0 || state.curValve.flow == 0) {
					// elephant needs new dest
					//fmt.Printf("elephant needs new dest\n")
					newECurValve = nil
				}

				if newCurValve != nil && newECurValve != nil {
					// neither is choosing a new destination
					//fmt.Printf("neither is choosing\n")
					//deadEnd = false
					searchQueue.push(&searchState{
						openValves:       newOpenValves,
						totalFlow:        newTotalFlow,
						pTimer:           newPTimer,
						eTimer:           newETimer,
						remainingMinutes: newRemainingMinutes,
						curValve:         newCurValve,
						eCurValve:        newECurValve,
						//history: append(state.history, historyEntry),
					})
				} else if newCurValve == nil && newECurValve == nil {
					// both are choosing a new destination
					//fmt.Printf("both are choosing\n")
					destinations := []string{}

					for _, v := range graph {
						if _, ok := state.openValves[v.name]; !ok && v.flow > 0 {
							destinations = append(destinations, v.name)
						}
					}

					for _, pDest := range destinations {
						if state.curValve.distances[pDest] > state.remainingMinutes {
							continue
						}
						for _, eDest := range destinations {
							if pDest == eDest || state.eCurValve.distances[eDest] > state.remainingMinutes {
								continue
							}
							//deadEnd = false
							searchQueue.push(&searchState{
								openValves:       newOpenValves,
								totalFlow:        newTotalFlow,
								pTimer:           state.curValve.distances[pDest] - 1,
								eTimer:           state.eCurValve.distances[eDest] - 1,
								remainingMinutes: newRemainingMinutes,
								curValve:         graph[pDest],
								eCurValve:        graph[eDest],
								//history: append(state.history, historyEntry),
							})
						}
					}
				} else if newCurValve == nil {
					// person is choosing a new destination
					//fmt.Printf("person is choosing\n")
					destinations := []string{}

					for _, v := range graph {
						if _, ok := state.openValves[v.name]; !ok && v.flow > 0 && v != newECurValve {
							destinations = append(destinations, v.name)
						}
					}

					if len(destinations) == 0 {
						// elephant going to last open valve, human can wander the rest of the time
						searchQueue.push(&searchState{
							openValves:       newOpenValves,
							totalFlow:        newTotalFlow,
							pTimer:           newETimer + 1,
							eTimer:           newETimer,
							remainingMinutes: newRemainingMinutes,
							curValve:         newECurValve,
							eCurValve:        newECurValve,
							//history: append(state.history, historyEntry),
						})
						continue
					}

					for _, pDest := range destinations {
						if state.curValve.distances[pDest] > state.remainingMinutes {
							continue
						}

						//deadEnd = false
						searchQueue.push(&searchState{
							openValves:       newOpenValves,
							totalFlow:        newTotalFlow,
							pTimer:           state.curValve.distances[pDest] - 1,
							eTimer:           newETimer,
							remainingMinutes: newRemainingMinutes,
							curValve:         graph[pDest],
							eCurValve:        newECurValve,
							//history: append(state.history, historyEntry),
						})
					}
				} else if newECurValve == nil {
					// elephant is choosing a new destination
					//fmt.Printf("elephant is choosing\n")
					destinations := []string{}

					for _, v := range graph {
						if _, ok := state.openValves[v.name]; !ok && v.flow > 0 && v != newCurValve {
							destinations = append(destinations, v.name)
						}
					}

					if len(destinations) == 0 {
						// human going to last open valve, elephant can wander the rest of the time
						searchQueue.push(&searchState{
							openValves:       newOpenValves,
							totalFlow:        newTotalFlow,
							pTimer:           newPTimer,
							eTimer:           newPTimer + 1,
							remainingMinutes: newRemainingMinutes,
							curValve:         newCurValve,
							eCurValve:        newCurValve,
							//history: append(state.history, historyEntry),
						})
						continue
					}

					for _, eDest := range destinations {
						if state.curValve.distances[eDest] > state.remainingMinutes {
							continue
						}

						//deadEnd = false

						searchQueue.push(&searchState{
							openValves:       newOpenValves,
							totalFlow:        newTotalFlow,
							pTimer:           newPTimer,
							eTimer:           state.eCurValve.distances[eDest] - 1,
							remainingMinutes: newRemainingMinutes,
							curValve:         newCurValve,
							eCurValve:        graph[eDest],
							//history: append(state.history, historyEntry),
						})
					}
				}

				//if deadEnd {
				//	histories = append(histories, state.history)
				//}
			}
		}()
	}

	wg.Wait()

	fmt.Printf("maxFlow: %d\n", maxFlow)
}

func copyMap(in map[string]int) map[string]int {
	out := map[string]int{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func addDistances(start *valve, graph map[string]*valve) {
	vq := valveQueue{
		elems: []*valveSearch{{start,0}},
	}

	for !vq.empty() {
		vs := vq.pop()
		if vs.v != start && start.distances[vs.v.name] == 0 {
			start.distances[vs.v.name] = vs.dist
		}

		for _, tunnel := range vs.v.children {
			if _, ok := start.distances[tunnel]; !ok {
				vq.push(&valveSearch{
					graph[tunnel], vs.dist+1,
				})
			}
		}
	}
}

func visitedKey(openValves map[string]int) string {
	open := make([]string, len(openValves))
	i := 0
	for v, _ := range openValves {
		open[i] = v
		i++
	}
	sort.Strings(open)

	return fmt.Sprintf("%+v", open)
}