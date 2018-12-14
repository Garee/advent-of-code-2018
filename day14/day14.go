package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readNRecipes() int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	nRecipes, _ := strconv.Atoi(scanner.Text())
	return nRecipes
}

func createNewRecipes(recipes []int, elves []int) []int {
	sum := 0
	for _, elf := range elves {
		sum += recipes[elf]
	}

	strSum := strconv.Itoa(sum)
	strDigits := strings.Split(strSum, "")

	digits := []int{}
	for _, strDigit := range strDigits {
		digit, _ := strconv.Atoi(strDigit)
		digits = append(digits, digit)
	}

	return append(recipes, digits...)
}

func setCurrentRecipes(recipes []int, elves []int) []int {
	nRecipes := len(recipes)
	for e, elf := range elves {
		elves[e] = (elf + (1 + recipes[elf])) % nRecipes
	}
	return elves
}

func createRecipesUntil(recipes []int, elves []int, nRecipes int) ([]int, []int) {
	for len(recipes) < nRecipes {
		recipes = createNewRecipes(recipes, elves)
		elves = setCurrentRecipes(recipes, elves)
	}
	return recipes, elves
}

func intArrToStr(arr []int) string {
	str := ""
	for _, i := range arr {
		str += strconv.Itoa(i)
	}
	return str
}

func createRecipesUntilEnd(recipes []int, elves []int, end string) ([]int, []int) {
	stop := false
	for !stop {
		recipes = createNewRecipes(recipes, elves)
		elves = setCurrentRecipes(recipes, elves)
		lenRecipes := len(recipes)
		if lenRecipes > len(end) {
			check := recipes[lenRecipes-len(end):]
			cmp := intArrToStr(check)
			if cmp == end {
				stop = true
			}
			// Either one or two new recipes can be added.
			check = recipes[lenRecipes-len(end)-1 : lenRecipes-1]
			cmp = intArrToStr(check)
			if cmp == end {
				stop = true
				recipes = recipes[:len(recipes)-1]
			}
		}
	}
	return recipes, elves
}

func main() {
	nRecipes := readNRecipes()
	recipes, _ := createRecipesUntil([]int{3, 7}, []int{0, 1}, nRecipes+10)
	lastTenRecipes := recipes[nRecipes : nRecipes+10]
	scores := intArrToStr(lastTenRecipes)
	endRecipes := strconv.Itoa(nRecipes)
	recipes, _ = createRecipesUntilEnd([]int{3, 7}, []int{0, 1}, endRecipes)
	leftRecipes := len(recipes) - len(endRecipes)
	fmt.Println(scores, leftRecipes)
}
