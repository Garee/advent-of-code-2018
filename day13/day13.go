package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func readInputLines() (lines []string, maxLineLen int) {
	maxLen := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		lineLen := len(line)
		if lineLen > maxLen {
			maxLen = lineLen
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read input lines:", err)
		os.Exit(1)
	}
	return lines, maxLen
}

func initTracksGrid(width int, height int) [][]byte {
	grid := make([][]byte, height)
	for row := 0; row < height; row++ {
		grid[row] = make([]byte, width)
		for col := 0; col < width; col++ {
			grid[row][col] = ' '
		}
	}
	return grid
}

func mapTracks(lines []string, width int, height int) ([][]byte, []Cart) {
	tracks := initTracksGrid(width, height)
	carts := []Cart{}
	for y, line := range lines {
		chars := []byte(line)
		for x := 0; x < len(chars); x++ {
			char := chars[x]
			if isCart(char) {
				carts = append(carts, Cart{len(carts), char, x, y, 0})
				tracks[y][x] = getGroundForCart(char)
			} else if char != ' ' {
				tracks[y][x] = char
			}
		}
	}
	return tracks, carts
}

func isCart(char byte) bool {
	return char == '<' || char == '^' || char == '>' || char == 'v'
}

func getGroundForCart(char byte) byte {
	if char == '<' || char == '>' {
		return '-'
	}

	return '|'
}

type Cart struct {
	id   int
	dir  byte
	x    int
	y    int
	turn int
}

type ByPosition []Cart

func (c ByPosition) Len() int {
	return len(c)
}

func (c ByPosition) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ByPosition) Less(i, j int) bool {
	if c[i].y == c[j].y {
		return c[i].x < c[j].x
	}

	return c[i].y < c[j].y
}

func printTracksGrid(tracks [][]byte) {
	fmt.Print("  ")
	for i := 0; i < 10; i++ {
		fmt.Print(i)
	}
	fmt.Println()
	for i, row := range tracks {
		fmt.Print(i, " ")
		for _, col := range row {
			fmt.Print(string(col))
		}
		fmt.Println()
	}
}

func printCarts(carts []Cart) {
	for _, cart := range carts {
		fmt.Printf("[%c %d,%d T%d]\n", cart.dir, cart.x, cart.y, cart.turn)
	}
}

func turn(dir byte, ground byte, turn int) byte {
	if ground == '+' {
		if turn%3 == 1 { // Straight
			return dir
		} else if turn%3 == 2 { // Right
			if dir == '<' {
				return '^'
			} else if dir == '^' {
				return '>'
			} else if dir == '>' {
				return 'v'
			} else if dir == 'v' {
				return '<'
			}
		} else if turn%3 == 0 { // Left
			if dir == '<' {
				return 'v'
			} else if dir == '^' {
				return '<'
			} else if dir == '>' {
				return '^'
			} else if dir == 'v' {
				return '>'
			}
		}
	} else if ground == '/' {
		if dir == '<' {
			return 'v'
		} else if dir == '^' {
			return '>'
		} else if dir == '>' {
			return '^'
		} else if dir == 'v' {
			return '<'
		}
	} else if ground == '\\' {
		if dir == '<' {
			return '^'
		} else if dir == '^' {
			return '<'
		} else if dir == '>' {
			return 'v'
		} else if dir == 'v' {
			return '>'
		}
	}

	return dir
}

func move(cart Cart, ground byte) Cart {
	if ground == '+' || ground == '\\' || ground == '/' {
		cart.dir = turn(cart.dir, ground, cart.turn)
		if ground == '+' {
			cart.turn++
		}
	}

	if cart.dir == '<' {
		cart.x -= 1
	} else if cart.dir == '>' {
		cart.x += 1
	} else if cart.dir == '^' {
		cart.y -= 1
	} else { // v
		cart.y += 1
	}

	return Cart{
		id:   cart.id,
		dir:  cart.dir,
		x:    cart.x,
		y:    cart.y,
		turn: cart.turn,
	}
}

func nextTick(tracks [][]byte, carts []Cart) ([][]byte, []Cart, Cart, bool) {
	sort.Sort(ByPosition(carts))
	for c, cart := range carts {
		ground := tracks[cart.y][cart.x]
		carts[c] = move(cart, ground)
		cart, _, crash := hasCrashOccurred(carts)
		if crash {
			return tracks, carts, cart, crash
		}
	}
	return tracks, carts, Cart{}, false
}

func hasCrashOccurred(carts []Cart) (Cart, Cart, bool) {
	for _, cart := range carts {
		for _, other := range carts {
			if cart != other && cart.x == other.x && cart.y == other.y {
				return cart, other, true
			}
		}
	}
	return Cart{}, Cart{}, false
}

func findFirstCrash(tracks [][]byte, carts []Cart) (x, y int) {
	var cart Cart
	crashed := false
	for !crashed {
		tracks, carts, cart, crashed = nextTick(tracks, carts)
	}
	return cart.x, cart.y
}

func nextTickRmCrashes(tracks [][]byte, carts []Cart) ([][]byte, []Cart) {
	sort.Sort(ByPosition(carts))
	settled := false
	start := 0
	for !settled {
		settled = true
		for c := start; c < len(carts); c++ {
			cart := carts[c]
			ground := tracks[cart.y][cart.x]
			carts[c] = move(cart, ground)
			cart, _, crash := hasCrashOccurred(carts)
			if crash {
				carts = removeCarts(cart.x, cart.y, carts)
				start = c - 1
				if len(carts) > 1 {
					settled = false
				}
				break
			}
		}
	}
	return tracks, carts
}

func findLastCartStanding(tracks [][]byte, carts []Cart) Cart {
	for len(carts) > 1 {
		tracks, carts = nextTickRmCrashes(tracks, carts)
	}
	return carts[0]
}

func removeCarts(x int, y int, carts []Cart) []Cart {
	result := make([]Cart, 0)
	for _, cart := range carts {
		if cart.x != x || cart.y != y {
			result = append(result, cart)
		}
	}
	return result
}

func main() {
	lines, maxLineLen := readInputLines()
	tracks, carts := mapTracks(lines, maxLineLen, len(lines))
	x, y := findFirstCrash(tracks, append([]Cart{}, carts...))
	fmt.Printf("%d,%d\n", x, y)
	lastCart := findLastCartStanding(tracks, carts)
	fmt.Printf("%d,%d\n", lastCart.x, lastCart.y)
}
