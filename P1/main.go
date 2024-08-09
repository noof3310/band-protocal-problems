package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	VALID   = "Good boy"
	INVALID = "Bad boy"
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

		// validate input and return output
		ok, err := isValid(input)
		if err != nil {
			fmt.Println(err.Error())
		} else if ok {
			fmt.Println(VALID)
		} else {
			fmt.Println(INVALID)
		}
	}
}

func isValid(str string) (bool, error) {
	// check obvious case
	if str[0] == 'R' || str[len(str)-1] == 'S' {
		return false, nil
	}

	// declare counting var
	n := 0 // n will be used to count current remaining S
	m := 0 // m will be used to count current group of R

	// check through each letter from left to right
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case 'R':
			// increase R count
			m++
		case 'S':
			// decrease amount of R from remaining S, if R > remaining S, remaining S will be 0
			n -= m
			if n < 0 {
				n = 0
			}
			// increase current remaining S and reset R counter to 0
			n++
			m = 0
		default:
			return false, fmt.Errorf("input invalid: contains character %c", str[i])
		}
	}

	// final check for the last group of R, if current R < current remaining S, then it's lack of revenge
	if m < n {
		return false, nil
	}

	return true, nil
}
