package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	input := os.Args
	if len(input) < 2 {
		fmt.Println("usage: classify <number>")
		return
	}

	num, error := strconv.Atoi(input[1])
	if error != nil {
		fmt.Println("error: not a number")
		return
	}

	if num > 0 {
		fmt.Println("positive")
	} else if num < 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("zero")
	}
}
