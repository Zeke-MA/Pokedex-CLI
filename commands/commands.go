package commands

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/Zeke-MA/pokedexcli/internal/pokeapi"
	"github.com/Zeke-MA/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(cfg *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, args []string) error
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
			description: "Displays name and id of the explorable locations. Limited to 20 results per command call.",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location and id if available.",
			Callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explores the given location for all possible pokemon. Use -name flag for the location or -id for the id.",
			Callback:    commandExplore,
		},
	}
}

func commandExit(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for key := range ValidCommands {
		fmt.Println(key + ": " + ValidCommands[key].description)
	}

	return nil
}

func commandMap(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, args []string) error {

	endpoint := "location-area"

	url := pokeapi.CreateUrl(client, endpoint)

	if config.Next != nil && *config.Next != "" {
		url = *config.Next
	}

	val, ok := cache.Get(url)
	if ok {
		fmt.Println(val)
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
		fmt.Println("Name: " + location.Name + " ID: " + path.Base(location.Url))
	}

	return nil
}

func commandMapb(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, args []string) error {

	endpoint := "location-area"

	url := pokeapi.CreateUrl(client, endpoint)

	if config.Previous != nil && *config.Previous != "" {
		url = *config.Previous
	} else if config.Previous == nil || *config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	val, ok := cache.Get(url)
	if ok {
		fmt.Println(val)
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

	cache.Add(url, body)

	results, err := pokeapi.Unmarshal[pokeapi.LocationArea](body, config)

	if err != nil {
		return err
	}

	for _, location := range results {
		fmt.Println("Name: " + location.Name + " ID: " + path.Base(location.Url))
	}

	return nil
}

func commandExplore(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, args []string) error {
	locationVal, err := parseExploreArgs(args)
	if err != nil {
		return err
	}

	endpoint := "location-area"

	url := pokeapi.CreateUrl(client, endpoint, locationVal)

	fmt.Print(url)

	val, ok := cache.Get(url)
	if ok {
		fmt.Println(val)
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

	cache.Add(url, body)

	fmt.Print(string(body))

	// Add in struct for the json structure of the location-area/id(name)
	// Call Unmarshal to get the results slice
	// print out the pokemon from the area

	return nil
}

func parseExploreArgs(args []string) (string, error) {
	flagSet := flag.NewFlagSet("explore", flag.ExitOnError)
	locationName := flagSet.String("name", "", "location area to explore")
	locationId := flagSet.Int("id", 0, "location area ID to explore")

	err := flagSet.Parse(args)

	if err != nil {
		return "", err
	}

	if *locationName != "" {
		return *locationName, nil
	}

	if *locationId != 0 {
		return strconv.Itoa(*locationId), nil
	}

	return "", fmt.Errorf("please provide either -name or -id flag")
}
