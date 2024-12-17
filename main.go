package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Zeke-MA/pokedexcli/commands"
	"github.com/Zeke-MA/pokedexcli/internal/pokeapi"
	"github.com/Zeke-MA/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokeClient := pokeapi.NewClient()
	pokeConfig := &pokeapi.Config{}
	pokeCache := pokecache.NewCache(20 * time.Second)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := strings.ToLower(scanner.Text())
			words := strings.Fields(input)
			command := words[0]
			args := words[1:]
			key, ok := commands.ValidCommands[command]
			if ok {
				key.Callback(pokeConfig, pokeClient, pokeCache, args)
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
