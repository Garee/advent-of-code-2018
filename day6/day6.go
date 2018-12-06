package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Coord struct {
	id int
	x  int
	y  int
}

func strToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func readCoords() (coords []Coord) {
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; scanner.Scan(); i++ {
		tokens := strings.Split(scanner.Text(), ", ")
		coords = append(coords, Coord{
			id: i + 65, // A
			x:  strToInt(tokens[0]),
			y:  strToInt(tokens[1]),
		})
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read coordinates:", err)
		os.Exit(1)
	}
	return coords
}

func findMinMaxXY(coords []Coord) (int, int, int, int) {
	minX, maxX := math.MaxInt32, math.MinInt32
	minY, maxY := math.MaxInt32, math.MinInt32
	for _, c := range coords {
		if c.x < minX {
			minX = c.x
		} else if c.x > maxX {
			maxX = c.x
		}
		if c.y < minY {
			minY = c.y
		} else if c.y > maxY {
			maxY = c.y
		}
	}

	return minX, maxX, minY, maxY
}

func countAreaSizes(minX, maxX, minY, maxY int, coords []Coord) map[Coord]int {
	areaSizes := make(map[Coord]int)
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			var minMhCoord Coord
			multiple := false
			minMhDist := math.MaxInt32
			for _, c := range coords {
				mhDist := abs(j-c.x) + abs(i-c.y)
				if mhDist < minMhDist {
					multiple = false
					minMhDist = mhDist
					minMhCoord = c
				} else if mhDist == minMhDist {
					multiple = true
				}
			}
			if !multiple {
				if i == minY || i == maxY || j == minX || j == maxX {
					areaSizes[minMhCoord] = -1
				} else {
					areaSizes[minMhCoord]++
				}
			}
		}
	}
	return areaSizes
}

func filterBoundaryCoords(minX, maxX, minY, maxY int, coords []Coord) (result []Coord) {
	for _, c := range coords {
		if c.x == minX || c.x == maxX || c.y == minY || c.y == maxY {
			continue
		}

		result = append(result, c)
	}

	return result
}

func findCoordWithLargestArea(coords []Coord, areaSizes map[Coord]int) (coord Coord, size int) {
	max := math.MinInt32
	for _, c := range coords {
		if areaSizes[c] > max {
			max = areaSizes[c]
			coord = c
		}
	}
	return coord, max
}

func calcSizeOfSafeRegion(minX, maxX, minY, maxY int, coords []Coord) (size int) {
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			sum := 0
			for _, c := range coords {
				sum += abs(j-c.x) + abs(i-c.y)
				if sum >= 10000 {
					break
				}
			}
			if sum < 10000 {
				size++
			}
		}
	}
	return size
}

func main() {
	coords := readCoords()
	minX, maxX, minY, maxY := findMinMaxXY(coords)
	areaSizes := countAreaSizes(minX, maxX, minY, maxY, coords)
	finiteCoords := filterBoundaryCoords(minX, maxX, minY, maxY, coords)
	_, size := findCoordWithLargestArea(finiteCoords, areaSizes)
	sizeOfSafeRegion := calcSizeOfSafeRegion(minX, maxX, minY, maxY, coords)
	fmt.Println(size, sizeOfSafeRegion)
}
