package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readBoxIDs() (ids []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read box IDs:", err)
		os.Exit(1)
	}
	return ids
}

func getCharCounts(str string) map[rune]int {
	counts := make(map[rune]int)
	for _, c := range str {
		if _, hasCount := counts[c]; hasCount {
			counts[c]++
		} else {
			counts[c] = 1
		}
	}
	return counts
}

func analyzeCounts(charCounts map[rune]int) (foundTwo bool, foundThree bool) {
	for _, count := range charCounts {
		if foundTwo && foundThree {
			break
		}
		if count == 2 && !foundTwo {
			foundTwo = true
		}
		if count == 3 && !foundThree {
			foundThree = true
		}
	}
	return foundTwo, foundThree
}

func computeChecksum(boxIDs []string) int {
	nIDsWithTwoCount := 0
	nIDsWithThreeCount := 0
	for _, id := range boxIDs {
		charCounts := getCharCounts(id)
		foundTwo, foundThree := analyzeCounts(charCounts)
		if foundTwo {
			nIDsWithTwoCount++
		}
		if foundThree {
			nIDsWithThreeCount++
		}
	}
	return nIDsWithTwoCount * nIDsWithThreeCount
}

func findNearlyMatchingBoxID(boxID string, boxIDs []string) (nearlyMatches string, found bool) {
	for _, id := range boxIDs {
		if id == boxID {
			continue
		}
		diffs := 0
		for i := range id {
			if id[i] != boxID[i] {
				diffs++
			}
		}
		if diffs <= 1 {
			return id, true
		}
	}
	return "", false
}

func findPrototypeBoxIDs(boxIDs []string) (string, string, bool) {
	for _, id := range boxIDs {
		if match, found := findNearlyMatchingBoxID(id, boxIDs); found {
			return id, match, true
		}
	}
	return "", "", false
}

func max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

func getCommonChars(str1 string, str2 string) (common string) {
	var builder strings.Builder
	maxLen := max(len(str1), len(str2))
	for i := 0; i < maxLen; i++ {
		if i >= len(str1) || i >= len(str2) {
			break
		}
		if str1[i] == str2[i] {
			builder.WriteByte(str1[i])
		}
	}
	return builder.String()
}

func getCommonCharsOfPrototypeBoxIDs(boxIDs []string) (string, bool) {
	if x, y, found := findPrototypeBoxIDs(boxIDs); found {
		return getCommonChars(x, y), true
	}
	return "", false
}

func main() {
	boxIDs := readBoxIDs()
	checksum := computeChecksum(boxIDs)
	fmt.Print(checksum)
	if commonChars, found := getCommonCharsOfPrototypeBoxIDs(boxIDs); found {
		fmt.Print(" ")
		fmt.Println(commonChars)
	}
}
