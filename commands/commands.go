package commands

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"

	"github.com/Zeke-MA/pokedexcli/internal/pokeapi"
	"github.com/Zeke-MA/pokedexcli/internal/pokecache"
	"github.com/Zeke-MA/pokedexcli/internal/pokedex"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(cfg *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error
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
		"catch": {
			name:        "catch",
			description: "Attempts to catch the pokemon in the given location",
			Callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects the pokemon in your pokedex for its stats",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all pokemon collected in your pokedex",
			Callback:    commandPokedex,
		},
	}
}

func commandExit(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for key := range ValidCommands {
		fmt.Println(key + ": " + ValidCommands[key].description)
	}

	return nil
}

func commandMap(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {

	endpoint := "location-area"

	url := pokeapi.CreateUrl(client, endpoint)

	if config.Next != nil && *config.Next != "" {
		url = *config.Next
	}

	val, ok := cache.Get(url)
	if ok {
		results, err := pokeapi.Unmarshal[pokeapi.LocationArea](val, config)

		if err != nil {
			return err
		}

		for _, location := range results {
			fmt.Println("Name: " + location.Name + " ID: " + path.Base(location.Url))
		}
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

	cache.Add(url, body)

	for _, location := range results {
		fmt.Println("Name: " + location.Name + " ID: " + path.Base(location.Url))
	}

	return nil
}

func commandMapb(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {

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
		fmt.Print("Using Cache!!!")
		results, err := pokeapi.Unmarshal[pokeapi.LocationArea](val, config)

		if err != nil {
			return err
		}

		for _, location := range results {
			fmt.Println("Name: " + location.Name + " ID: " + path.Base(location.Url))
		}
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

func commandExplore(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {
	locationVal, err := parseExploreArgs(args)
	if err != nil {
		return err
	}

	endpoint := "location-area"

	url := pokeapi.CreateUrl(client, endpoint, locationVal)

	val, ok := cache.Get(url)
	// need a better way to handle cache
	if ok {
		fmt.Print("Using Cache!!!")
		names, err := pokeapi.UnmarshalExplore(val, config)

		if err != nil {
			return err
		}

		for _, name := range names {
			fmt.Println("Name: " + name)
		}
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

	names, err := pokeapi.UnmarshalExplore(body, config)

	if err != nil {
		return err
	}

	for _, name := range names {
		fmt.Println("Name: " + name)
	}

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

func commandCatch(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {
	endpoint := "pokemon"
	pokemonName := args[0]

	url := pokeapi.CreateUrl(client, endpoint, pokemonName)

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

	pokemon, err := pokeapi.UnmarshalPokemonData(body, config)

	if err != nil {
		return err
	}

	// Determine success rate for RNG gods
	success := 100 - (float64(pokemon.BaseExperience) * 0.08)

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	rng := rand.Intn(100)

	if rng <= int(success) {
		fmt.Printf("You have caught %v!\n", pokemon.Name)
		_, ok := pokedex.CaughtPokemon[pokemonName]
		if ok {
			fmt.Printf("%v is already in your pokedex\n", pokemon.Name)
		} else {
			pokedex.CaughtPokemon[pokemonName] = pokemon
			fmt.Printf("%v has been added to your pokedex!\n", pokemon.Name)
		}

	} else {
		fmt.Printf("%v has escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {
	name := args[0]

	val, ok := pokedex.CaughtPokemon[name]
	if !ok {
		fmt.Printf("%v has not been caught. Unable to inspect\n", name)
		return nil
	}

	fmt.Printf("Height: %v\n", val.Height)
	fmt.Printf("Weight %v\n", val.Weight)
	fmt.Println("Stats: ")
	for _, stat := range val.Stats {
		fmt.Printf("-%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types: ")

	for _, ty := range val.Types {
		fmt.Printf("-%v\n", ty.Type.Name)
	}

	return nil
}

func commandPokedex(config *pokeapi.Config, client *pokeapi.Client, cache *pokecache.Cache, pokedex *pokedex.Pokedex, args []string) error {
	for key := range pokedex.CaughtPokemon {
		fmt.Printf("-%v\n", key)
	}
	return nil
}
