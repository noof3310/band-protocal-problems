package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// scan input
	scanner := bufio.NewScanner(os.Stdin)

	// instruct for input
	fmt.Println("Enter input (type 'exit' to quit):")

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			break
		}

		fmt.Print("> ")
		scanner.Scan()
		input2 := scanner.Text()

		// extract input
		n, k := extractInput(input, input2)

		// calculate
		res := maxAmountInRange(n, k)
		fmt.Println(res)
	}
}

// extract the length of roof and chicken position from input string to int and []int
func extractInput(s1, s2 string) (int, []int) {
	n, _ := strconv.Atoi(strings.Split(s1, " ")[1])

	items := strings.Split(s2, " ")
	k := make([]int, len(items))
	for i, item := range items {
		k[i], _ = strconv.Atoi(item)
	}

	return n, k
}

func maxAmountInRange(n int, k []int) int {
	maxCount := 0
	j := 0
	for i := 0; i < len(k); i++ {
		for j < len(k) && k[j] < k[i]+n {
			j++
		}
		if j-i > maxCount {
			maxCount = j - i
		}
	}
	return maxCount
}
