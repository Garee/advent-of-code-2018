package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func readSerialNumber() int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	serialNumber, _ := strconv.Atoi(scanner.Text())
	return serialNumber
}

func getHundrethDigit(n int) int {
	if n < 100 && n > -100 {
		return 0
	} else if n < 0 {
		n = -n
	}

	s := strconv.Itoa(n)
	sdigit := string(s[len(s)-3])
	digit, _ := strconv.Atoi(sdigit)
	return digit
}

func calcCellPowerLevel(x int, y int, grid [300][300]int, serialNumber int) (powerLevel int) {
	rackID := x + 10
	powerLevel = rackID * y
	powerLevel += serialNumber
	powerLevel *= rackID
	powerLevel = getHundrethDigit(powerLevel)
	powerLevel -= 5
	return powerLevel
}

func calcCellPowerLevels(grid [300][300]int, serialNumber int) [300][300]int {
	for row := 0; row < 300; row++ {
		for col := 0; col < 300; col++ {
			grid[row][col] = calcCellPowerLevel(row, col, grid, serialNumber)
		}
	}
	return grid
}

func calcCellTotalPower(x int, y int, grid [300][300]int, size int) (power int) {
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			power += grid[x+row][y+col]
		}
	}
	return power
}

func calcCellTotalPowerLevels(grid [300][300]int, size int) [300][300]int {
	for row := 0; row <= 300-size; row++ {
		for col := 0; col <= 300-size; col++ {
			grid[row][col] = calcCellTotalPower(row, col, grid, size)
		}
	}
	return grid
}

func findMaxTotalPowerLevel(grid [300][300]int, size int) (x int, y int, power int) {
	max := math.MinInt32
	for row := 0; row <= 300-size; row++ {
		for col := 0; col <= 300-size; col++ {
			totalPower := grid[row][col]
			if totalPower > max {
				max = totalPower
				x = row
				y = col
			}
		}
	}
	return x, y, max
}

func findMaxTotalPowerSquare(grid [300][300]int) (x int, y int, size int) {
	maxPower := math.MinInt32
	for i := 1; i <= 300; i++ {
		totals := calcCellTotalPowerLevels(grid, i)
		a, b, power := findMaxTotalPowerLevel(totals, i)
		if power > maxPower {
			maxPower = power
			x, y, size = a, b, i
		}

	}
	return x, y, size
}

func main() {
	serialNumber := readSerialNumber()
	grid := calcCellPowerLevels([300][300]int{}, serialNumber)
	grid = calcCellTotalPowerLevels(grid, 3)
	x, y, _ := findMaxTotalPowerLevel(grid, 3)
	fmt.Printf("%d,%d", x, y)
	x, y, size := findMaxTotalPowerSquare(grid) // This is SLOW.
	fmt.Printf(" %d,%d,%d\n", x, y, size)
}
