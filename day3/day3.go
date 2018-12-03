package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Claim struct {
	id       string
	x        int
	y        int
	width    int
	height   int
	overlaps bool
}

func strToInt(str string) int {
	if i, err := strconv.Atoi(str); err == nil {
		return i
	}

	return 0
}

func readClaims() (claims []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		claims = append(claims, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read claims", err)
		os.Exit(1)
	}
	return claims
}

func parseClaims(claimStrings []string) (claims []Claim) {
	for _, str := range claimStrings {
		tokens := strings.Split(str, " ")
		id := tokens[0][1:]
		xyTokens := strings.Split(tokens[2], ",")
		x, y := xyTokens[0], xyTokens[1]
		y = y[:len(y)-1]
		whTokens := strings.Split(tokens[3], "x")
		width, height := whTokens[0], whTokens[1]
		claim := Claim{
			id:       id,
			x:        strToInt(x),
			y:        strToInt(y),
			width:    strToInt(width),
			height:   strToInt(height),
			overlaps: false,
		}
		claims = append(claims, claim)
	}
	return claims
}

func markFabric(claims []Claim) (fabric [1000][1000]int) {
	for _, claim := range claims {
		for row := claim.y; row < claim.height+claim.y; row++ {
			for col := claim.x; col < claim.width+claim.x; col++ {
				fabric[row][col]++
			}
		}
	}
	return fabric
}

func countTwoOrMoreClaimOverlaps(fabric [1000][1000]int) (count int) {
	for _, row := range fabric {
		for _, col := range row {
			if col > 1 {
				count++
			}
		}
	}
	return count
}

func filterOverlaps(claims []Claim, markedFabric [1000][1000]int) (noOverlaps []Claim) {
	for _, claim := range claims {
		nOverlaps := 0
		for row := claim.y; row < claim.height+claim.y; row++ {
			for col := claim.x; col < claim.width+claim.x; col++ {
				if markedFabric[row][col] > 1 {
					nOverlaps++
				}
			}
		}
		if nOverlaps == 0 {
			noOverlaps = append(noOverlaps, claim)
		}
	}
	return noOverlaps
}

func main() {
	claimStrings := readClaims()
	claims := parseClaims(claimStrings)
	fabric := markFabric(claims)
	nOverlaps := countTwoOrMoreClaimOverlaps(fabric)
	fmt.Print(nOverlaps, " ")
	noOverlaps := filterOverlaps(claims, fabric)
	for _, claim := range noOverlaps {
		fmt.Println(claim.id)
	}
}
