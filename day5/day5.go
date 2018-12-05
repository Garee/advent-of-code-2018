package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func readPolymer() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read polymer:", err)
		os.Exit(1)
	}
	return scanner.Text()
}

func react(polymer string) string {
	finished := false
	for !finished {
		finished = true
		for i := 1; i < len(polymer)-1; i++ {
			a, b := string(polymer[i-1]), string(polymer[i])
			if a != b && strings.ToLower(a) == strings.ToLower(b) {
				polymer = polymer[:i-1] + polymer[i+1:]
				finished = false
			}
		}
	}
	return polymer
}

func set(str string) []string {
	charCounts := make(map[rune]int)
	for _, c := range str {
		charCounts[c]++
	}
	chars := []string{}
	for c := range charCounts {
		chars = append(chars, string(c))
	}
	return chars
}

func getUnits(polymer string) []string {
	return set(strings.ToLower(polymer))
}

func removeUnit(unit string, polymer string) string {
	polymer = strings.Replace(polymer, strings.ToLower(unit), "", -1)
	polymer = strings.Replace(polymer, strings.ToUpper(unit), "", -1)
	return polymer
}

func findShortestPolymerLen(polymer string) int {
	shortest := math.MaxInt64
	for _, u := range getUnits(polymer) {
		p := removeUnit(u, polymer)
		pLen := len(react(p))
		if pLen < shortest {
			shortest = pLen
		}
	}
	return shortest
}

func main() {
	polymer := readPolymer()
	reactedPolymer := react(polymer)
	shortestPolymerLen := findShortestPolymerLen(polymer)
	fmt.Println(len(reactedPolymer), shortestPolymerLen)
}
