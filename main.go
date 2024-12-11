package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	cleanedInput := strings.Fields(text)
	for idx, word := range cleanedInput {
		cleanedInput[idx] = strings.ToLower(word)
	}
	return cleanedInput
}
