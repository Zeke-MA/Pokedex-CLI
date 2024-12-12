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
			description: "Displays names of the explorable locations. Limited to 20 results per command call.",
			Callback:    func() error { return commandMap(config) },
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations if available.",
			Callback:    func() error { return commandMapb(config) },
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

	baseUrl := "https://pokeapi.co/api/v2/"

	endpoint := "location-area/"

	finalUrl := baseUrl + endpoint

	if config.Next != nil && *config.Next != "" {
		finalUrl = *config.Next
	}

	req, err := http.NewRequest("GET", finalUrl, nil)
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

	fmt.Println(*config.Next)

	return nil
}

func commandMapb(config *Config) error {

	baseUrl := "https://pokeapi.co/api/v2/"

	endpoint := "location-area/"

	finalUrl := baseUrl + endpoint

	if config.Previous != nil && *config.Previous != "" {
		finalUrl = *config.Previous
	} else if config.Previous == nil || *config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	req, err := http.NewRequest("GET", finalUrl, nil)
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

	fmt.Println(*config.Next)
	fmt.Println(*config.Previous)

	return nil
}
