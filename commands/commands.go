package commands

import (
	"fmt"
	"os"

	"github.com/Zeke-MA/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(cfg *pokeapi.Config, client *pokeapi.Client) error
}

var ValidCommands map[string]cliCommand

func init() {

	ValidCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays names of the explorable locations. Limited to 20 results per command call.",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations if available.",
			Callback:    commandMapb,
		},
	}
}

func commandExit(config *pokeapi.Config, client *pokeapi.Client) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config, client *pokeapi.Client) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for key := range ValidCommands {
		fmt.Println(key + ": " + ValidCommands[key].description)
	}

	return nil
}

func commandMap(config *pokeapi.Config, client *pokeapi.Client) error {

	endpoint := "location-area"

	url := pokeapi.CreateUrl(client, endpoint)

	if config.Next != nil && *config.Next != "" {
		url = *config.Next
	}

	request, err := pokeapi.NewRequest("GET", url, nil, client)

	if err != nil {
		return err
	}

	response, err := pokeapi.DoRequest(request, client)

	if err != nil {
		return err
	}

	body, err := pokeapi.GetResponse(response)
	if err != nil {
		return err
	}

	results, err := pokeapi.Unmarshal[pokeapi.LocationArea](body, config)

	if err != nil {
		return err
	}

	for _, location := range results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(config *pokeapi.Config, client *pokeapi.Client) error {

	endpoint := "location-area"

	url := pokeapi.CreateUrl(client, endpoint)

	if config.Previous != nil && *config.Previous != "" {
		url = *config.Previous
	} else if config.Previous == nil || *config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	request, err := pokeapi.NewRequest("GET", url, nil, client)

	if err != nil {
		return err
	}

	response, err := pokeapi.DoRequest(request, client)

	if err != nil {
		return err
	}

	body, err := pokeapi.GetResponse(response)
	if err != nil {
		return err
	}

	results, err := pokeapi.Unmarshal[pokeapi.LocationArea](body, config)

	if err != nil {
		return err
	}

	for _, location := range results {
		fmt.Println(location.Name)
	}

	return nil
}
