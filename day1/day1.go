package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readFrequencyChanges() (changes []int) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		freqStr := scanner.Text()
		if freq, err := strconv.Atoi(freqStr); err == nil {
			changes = append(changes, freq)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read frequency changes:", err)
		os.Exit(1)
	}
	return changes
}

func analyze(changes []int) (total int, dup int) {
	freq := 0
	freqTracker := map[int]bool{0: true}

	calculatedFreq, foundDup := false, false
	for i := 0; !calculatedFreq || !foundDup; i++ {
		if i == len(changes) {
			i = 0

			if !calculatedFreq {
				total = freq
				calculatedFreq = true
			}
		}

		freq += changes[i]
		if _, isDup := freqTracker[freq]; isDup {
			dup = freq
			foundDup = true
		} else {
			freqTracker[freq] = true
		}
	}
	return total, dup
}

func main() {
	changes := readFrequencyChanges()
	freq, dup := analyze(changes)
	fmt.Println(freq, dup)
}
