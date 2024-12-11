package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zeke-MA/pokedexcli/commands"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := strings.ToLower(scanner.Text())
			words := strings.Fields(input)
			command := words[0]
			key, ok := commands.ValidCommands[command]
			if ok {
				key.Callback()
			} else {
				fmt.Println("Unknown command")
			}

		} else {
			fmt.Printf("Error reading input: %v", scanner.Err())
		}

	}

}

func cleanInput(text string) []string {
	cleanedInput := strings.Fields(text)
	for idx, word := range cleanedInput {
		cleanedInput[idx] = strings.ToLower(word)
	}
	return cleanedInput
}
