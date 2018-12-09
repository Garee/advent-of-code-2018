package main

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func readGameSettings() (int, int) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read game settings:", err)
		os.Exit(1)
	}
	tokens := strings.Split(scanner.Text(), " ")
	nPlayers, _ := strconv.Atoi(tokens[0])
	lastMarble, _ := strconv.Atoi(tokens[6])
	return nPlayers, lastMarble
}

func playGame(nPlayers int, lastMarble int) []int {
	scores := make([]int, nPlayers)
	playerIdx := 0

	circle := list.New()
	current := circle.PushFront(0)
	nextMarble := 1

	for nextMarble <= lastMarble {
		if nextMarble%23 == 0 {
			for i := 0; i < 7; i++ {
				if current == circle.Front() {
					current = circle.Back()
				} else {
					current = current.Prev()
				}
			}
			scores[playerIdx] += nextMarble + current.Value.(int)
			current = current.Next()
			circle.Remove(current.Prev())
		} else {
			for i := 0; i < 2; i++ {
				if current == circle.Back() {
					current = circle.Front()
				} else {
					current = current.Next()
				}
			}
			current = circle.InsertBefore(nextMarble, current)
		}

		nextMarble++
		playerIdx = (playerIdx + 1) % nPlayers
	}

	return scores
}

func findMax(arr []int) int {
	max := math.MinInt32
	for _, n := range arr {
		if n > max {
			max = n
		}
	}
	return max
}

func main() {
	nPlayers, lastMarble := readGameSettings()
	maxScore := findMax(playGame(nPlayers, lastMarble))
	hundredTimesMaxScore := findMax(playGame(nPlayers, lastMarble*100))
	fmt.Println(maxScore, hundredTimesMaxScore)
}
