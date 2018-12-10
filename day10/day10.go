package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

func readPositionChanges() (changes []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		changes = append(changes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read position changes:", err)
		os.Exit(1)
	}
	return changes
}

type PosChange struct {
	xPos int
	yPos int
	xVel int
	yVel int
}

func createPosChange(xp, yp, xv, yv string) PosChange {
	xPos, _ := strconv.Atoi(xp)
	yPos, _ := strconv.Atoi(yp)
	xVel, _ := strconv.Atoi(xv)
	yVel, _ := strconv.Atoi(yv)
	return PosChange{xPos, yPos, xVel, yVel}
}

func parsePositionChanges(changesStrs []string) (changes []PosChange) {
	for _, str := range changesStrs {
		re := regexp.MustCompile("-?\\d+")
		tokens := re.FindAllString(str, -1)
		xPos, yPos := tokens[0], tokens[1]
		xVel, yVel := tokens[2], tokens[3]
		change := createPosChange(xPos, yPos, xVel, yVel)
		changes = append(changes, change)
	}
	return changes
}

func findMinMaxXY(changes []PosChange) (int, int, int, int) {
	xMin, yMin := math.MaxInt32, math.MaxInt32
	xMax, yMax := math.MinInt32, math.MinInt32
	for _, change := range changes {
		if change.xPos < xMin {
			xMin = change.xPos
		}
		if change.yPos < yMin {
			yMin = change.yPos
		}
		if change.xPos > xMax {
			xMax = change.xPos
		}
		if change.yPos > yMax {
			yMax = change.yPos
		}
	}
	return xMin, xMax, yMin, yMax
}

func createEmptySky(width int, height int) [][]string {
	sky := make([][]string, height+1)
	for i := range sky {
		row := make([]string, width+1)
		for j := range row {
			row[j] = "."
		}
		sky[i] = row
	}
	return sky
}

func printSky(sky [][]string) {
	for _, row := range sky {
		for _, col := range row {
			fmt.Print(col, " ")
		}
		fmt.Println()
	}
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func addStarsToSky(sky [][]string, changes []PosChange, minX int, minY int) [][]string {
	for _, change := range changes {
		x, y := change.xPos, change.yPos
		if minX < 0 {
			x = change.xPos + absInt(minX)
		}
		if minY < 0 {
			y = change.yPos + absInt(minY)
		}

		sky[y][x] = "#"
	}
	return sky
}

func calcWidth(xMin int, xMax int) int {
	if xMin < 0 {
		return xMax - xMin
	}

	return xMax
}

func calcHeight(yMin int, yMax int) int {
	if yMin < 0 {
		return yMax - yMin
	}

	return yMax
}

func printPosChanges(changes []PosChange) int {
	xMin, xMax, yMin, yMax := findMinMaxXY(changes)
	width := xMax - xMin
	height := yMax - yMin

	secsElapsed := 0
	for {
		nextPos := make([]PosChange, len(changes))
		copy(nextPos, changes)
		nextPos = updatePositions(nextPos)
		secsElapsed++

		nextXMin, nextXMax, nextYMin, nextYMax := findMinMaxXY(nextPos)
		nextWidth := nextXMax - nextXMin
		nextHeight := nextYMax - nextYMin

		if nextWidth >= width && nextHeight >= height {
			break
		}

		changes = nextPos
		width = nextWidth
		height = nextHeight
		xMin, xMax, yMin, yMax = nextXMin, nextXMax, nextYMin, nextYMax
	}

	width = calcWidth(xMin, xMax)
	height = calcHeight(yMin, yMax)
	sky := createEmptySky(width, height)
	sky = addStarsToSky(sky, changes, xMin, yMin)
	printSky(sky)
	return secsElapsed - 1
}

func updatePositions(changes []PosChange) []PosChange {
	for i := 0; i < len(changes); i++ {
		changes[i].xPos += changes[i].xVel
		changes[i].yPos += changes[i].yVel
	}
	return changes
}

func main() {
	changes := parsePositionChanges(readPositionChanges())
	secs := printPosChanges(changes)
	fmt.Println(secs)
}
