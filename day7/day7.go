package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func readRequirements() (requirements []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		req := scanner.Text()
		requirements = append(requirements, req)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read requirements:", err)
		os.Exit(1)
	}
	return requirements
}

type Step struct {
	id   string
	deps []string
	reqs []string
}

func parseStepIds(requirement string) (string, string) {
	tokens := strings.Split(requirement, " ")
	return tokens[1], tokens[7]
}

func parseSteps(requirements []string) (steps map[string]Step) {
	steps = make(map[string]Step)
	for _, req := range requirements {
		depStepID, forStepID := parseStepIds(req)
		depStep := createStepIfMissing(depStepID, steps)
		steps[depStepID] = Step{
			id:   depStepID,
			deps: depStep.deps,
			reqs: append(depStep.reqs, forStepID),
		}
		forStep := createStepIfMissing(forStepID, steps)
		steps[forStepID] = Step{
			id:   forStepID,
			deps: append(forStep.deps, depStepID),
			reqs: forStep.reqs,
		}
	}
	return steps
}

func createStepIfMissing(id string, steps map[string]Step) Step {
	if _, found := steps[id]; !found {
		steps[id] = Step{
			id:   id,
			deps: []string{},
		}
	}

	return steps[id]
}

func findStartSteps(steps map[string]Step) (noDeps []string) {
	for _, st := range steps {
		if len(st.deps) == 0 {
			noDeps = append(noDeps, st.id)
		}
	}

	sort.Strings(noDeps)
	return noDeps
}

func traverse(avail []string, counts map[string]int, steps map[string]Step) string {
	s := ""
	sort.Strings(avail)
	for i, availID := range avail {
		step := steps[availID]
		sort.Strings(step.reqs)
		sort.Strings(step.deps)

		if count := counts[availID]; count == -1 || count != len(step.deps) {
			continue
		}

		s += availID
		counts[availID] = -1

		sort.Strings(step.reqs)
		for _, depUpon := range step.reqs {
			counts[depUpon]++
		}

		to := append([]string{}, avail...)
		to = append(to[:i], to[i+1:]...)
		to = append(to, step.reqs...)
		to = toSet(to)
		s += traverse(to, counts, steps)
	}
	return s
}

func toSet(strs []string) (set []string) {
	m := make(map[string]string)
	for _, s := range strs {
		if _, found := m[s]; !found {
			m[s] = s
			set = append(set, s)
		}
	}
	return set
}

func printSteps(steps map[string]Step) {
	for k, st := range steps {
		fmt.Println(k)
		fmt.Print("deps => ")
		for _, d := range st.deps {
			fmt.Print(d)
		}
		fmt.Print("\nreqs => ")
		for _, r := range st.reqs {
			fmt.Print(r)
		}
		fmt.Println()
	}
}

func main() {
	requirements := readRequirements()
	steps := parseSteps(requirements)
	noDeps := findStartSteps(steps)
	order := traverse(noDeps, make(map[string]int), steps)
	fmt.Println(order)
}
