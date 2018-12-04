package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func readGuardEvents() (events []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		events = append(events, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read guard events: ", err)
		os.Exit(-1)
	}
	return events
}

type GuardEvent struct {
	datetime time.Time
	action   string
}

func parseGuardEvents(eventStrs []string) (events []GuardEvent) {
	for _, eventStr := range eventStrs {
		dt, _ := time.Parse("[2006-01-02 15:04]", eventStr[:18])
		event := GuardEvent{
			datetime: dt,
			action:   eventStr[19:],
		}
		events = append(events, event)
	}
	return events
}

type ByDatetime []GuardEvent

func (events ByDatetime) Len() int {
	return len(events)
}

func (events ByDatetime) Less(a, b int) bool {
	return events[a].datetime.Before(events[b].datetime)
}

func (events ByDatetime) Swap(a, b int) {
	events[a], events[b] = events[b], events[a]
}

type Guard struct {
	id     string
	events []GuardEvent
}

func findGuard(id string, guards []Guard) *Guard {
	for i := 0; i < len(guards); i++ {
		if guards[i].id == id {
			return &guards[i]
		}
	}

	return nil
}

func parseGuards(events []GuardEvent) map[string]Guard {
	guards := make(map[string]Guard)
	var guardID string
	for _, event := range events {
		tokens := strings.Split(event.action, " ")
		if len(tokens) > 2 {
			guardID = tokens[1][1:]
			if _, found := guards[guardID]; !found {
				guards[guardID] = Guard{
					id:     guardID,
					events: []GuardEvent{},
				}
			}
		} else if guard, found := guards[guardID]; found {
			guard.events = append(guard.events, event)
			guards[guardID] = guard
		}
	}
	return guards
}

func getGuardWhoSleptMost(guards map[string]Guard) (Guard, int) {
	var sleptMost Guard
	mostMinsSlept := 0
	for _, guard := range guards {
		minsSlept := 0
		for e := 0; e < len(guard.events)-1; e += 2 {
			fellAsleepAt := guard.events[e].datetime
			wokeUpAt := guard.events[e+1].datetime
			sleptFor := wokeUpAt.Sub(fellAsleepAt)
			minsSlept += int(sleptFor.Minutes())
		}
		if minsSlept > mostMinsSlept {
			mostMinsSlept = minsSlept
			sleptMost = guard
		}
	}
	return sleptMost, mostMinsSlept
}

func getMinsSleptForEachMin(guard Guard) (mins map[int]int) {
	mins = make(map[int]int)
	for e := 0; e < len(guard.events)-1; e += 2 {
		fellAsleepAt := guard.events[e].datetime
		wokeUpAt := guard.events[e+1].datetime
		if fellAsleepAt.Day() != wokeUpAt.Day() {
			fmt.Println("Aha!")
		}
		for m := fellAsleepAt.Minute(); m < wokeUpAt.Minute(); m++ {
			mins[m]++
		}
	}

	return mins
}

func getMinSleptMost(guard Guard) int {
	sleptMost, forMins := 0, 0
	for m, t := range getMinsSleptForEachMin(guard) {
		if t > forMins {
			forMins = t
			sleptMost = m
		}
	}
	return sleptMost
}

func getGuardWhoSleptOnSameMinMostFreq(guards map[string]Guard) (Guard, int) {
	guardID, min, freq := "", 0, 0
	for id, g := range guards {
		for m, f := range getMinsSleptForEachMin(g) {
			if f > freq {
				freq = f
				min = m
				guardID = id
			}
		}
	}
	return guards[guardID], min
}

func main() {
	eventStrs := readGuardEvents()
	events := parseGuardEvents(eventStrs)
	sort.Sort(ByDatetime(events))
	guards := parseGuards(events)
	guard, _ := getGuardWhoSleptMost(guards)
	if id, err := strconv.Atoi(guard.id); err == nil {
		minSleptMost := getMinSleptMost(guard)
		fmt.Print(id * minSleptMost)
	}
	sameMinGuard, min := getGuardWhoSleptOnSameMinMostFreq(guards)
	if id, err := strconv.Atoi(sameMinGuard.id); err == nil {
		fmt.Print(" ", id*min)
	}
}
