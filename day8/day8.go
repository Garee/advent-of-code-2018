package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readLicenseFile() (nums []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		nums = append(nums, tokens...)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read license file: ", err)
		os.Exit(1)
	}
	return nums
}

func toIntArr(strs []string) (arr []int) {
	for _, s := range strs {
		n, _ := strconv.Atoi(s)
		arr = append(arr, n)
	}
	return arr
}

type Node struct {
	id       int
	metadata []int
	children []Node
}

func buildTree(nChildren, nMetadata int, nums []int) (Node, []int) {
	node := Node{
		id:       len(nums),
		children: []Node{},
	}

	if nChildren == 0 {
		node.metadata = nums[:nMetadata]
		return node, nums[nMetadata:]
	}

	next := nums
	for nChildren > 0 {
		child, remain := buildTree(next[0], next[1], next[2:])
		node.children = append(node.children, child)
		nChildren--
		next = append([]int{}, remain...)
	}

	node.metadata = next[:nMetadata]
	return node, next[nMetadata:]
}

func sumMetadata(node Node) (sum int) {
	for _, m := range node.metadata {
		sum += m
	}
	for _, c := range node.children {
		sum += sumMetadata(c)
	}
	return sum
}

func main() {
	nums := toIntArr(readLicenseFile())
	root, _ := buildTree(nums[0], nums[1], nums[2:])
	sum := sumMetadata(root)
	fmt.Println(sum)
}
