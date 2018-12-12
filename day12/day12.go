package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
)

func readInput() (lines []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read input:", err)
		os.Exit(1)
	}
	return lines
}

func parseInitialState(lines []string) string {
	tokens := strings.Split(lines[0], " ")
	return tokens[2]
}

type Note struct {
	cond    string
	toState byte
}

func parseNotes(lines []string) (notes []Note) {
	noteLines := lines[2:]
	for _, line := range noteLines {
		tokens := strings.Split(line, " => ")
		notes = append(notes, Note{tokens[0], tokens[1][0]})
	}
	return notes
}

func createNewGen(length int) []byte {
	next := make([]byte, length)
	for i := 0; i < length; i++ {
		next[i] = '.'
	}
	return next
}

func growGen(gen []byte, centre int) ([]byte, int) {
	emptyPots := []byte(".....")
	if string(gen[:5]) != string(emptyPots) {
		gen = append(emptyPots, gen...)
		centre += 5
	}
	if string(gen[len(gen)-5:]) != string(emptyPots) {
		gen = append(gen, emptyPots...)
	}
	return gen, centre
}

func findIndexes(state []byte, cond string) (indices []int) {
	bState := []byte(state)
	expr := strings.Replace(cond, ".", "[.]", -1)
	re := regexp.MustCompile(expr)
	match := re.FindIndex(bState)
	for match != nil {
		idx := match[0] + 2
		indices = append(indices, len(state)-len(bState)+idx)
		bState = bState[match[0]+1:]
		match = re.FindIndex(bState)
	}
	return indices
}

func evolve(state string, notes []Note, centre int) (string, int) {
	bState, nextCentre := growGen([]byte(state), centre)
	next := createNewGen(len(bState))
	for _, note := range notes {
		indices := findIndexes(bState, note.cond)
		for _, idx := range indices {
			next[idx] = note.toState
		}
	}
	return string(next), nextCentre
}

func evolveFor(state string, notes []Note, gens int) (string, int) {
	next, centre := state, 0
	prev := next
	prevSum := sumPots(prev, centre)
	prevDiff := math.MaxInt32
	for g := 0; g < gens; g++ {
		next, centre = evolve(next, notes, centre)
		if next == prev {
			break
		}
		sum := sumPots(next, centre)
		diff := sum - prevSum

		if diff == prevDiff {
			part2 := ((gens - (g + 1)) * diff) + sum
			fmt.Println(part2)
			return "", 0 // We only care about the sum for part2
		}

		prevSum = sum
		prevDiff = diff
	}
	return next, centre
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func sumPots(pots string, centre int) (sum int) {
	for i, p := range pots {
		if p == '#' {
			sum += (i - centre)
		}
	}
	return sum
}

func main() {
	lines := readInput()
	pots := parseInitialState(lines)
	notes := parseNotes(lines)
	gen20, centre := evolveFor(pots, notes, 20)
	sum := sumPots(gen20, centre)
	fmt.Println(sum)
	gen50b, centre := evolveFor(pots, notes, 50000000000)
	sum = sumPots(gen50b, centre)
	fmt.Println(sum)
}
