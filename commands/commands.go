package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	Callback    func() error
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type APIResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Config struct {
	Next     *string
	Previous *string
}

var ValidCommands map[string]cliCommand
var config = &Config{}

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
			description: "Displays names of the explorable locations",
			Callback:    func() error { return commandMap(config) },
		},
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for key := range ValidCommands {
		fmt.Println(key + ": " + ValidCommands[key].description)
	}

	return nil
}

func commandMap(config *Config) error {

	endpoint := "https://pokeapi.co/api/v2/location-area/"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("unable to communicate with PokeAPI at this time: %v", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("cannot perform request: %v", err)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("unable to retreive body from response: %v", err)
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}
	defer res.Body.Close()

	var LocationAreaResponse APIResponse
	err = json.Unmarshal(body, &LocationAreaResponse)
	if err != nil {
		return fmt.Errorf("unable to parse json: %v", err)
	}

	config.Next = &LocationAreaResponse.Next
	config.Previous = &LocationAreaResponse.Previous

	for _, location := range LocationAreaResponse.Results {
		fmt.Println(location.Name)
	}

	return nil
}
