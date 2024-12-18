package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const baseUrl = "https://pokeapi.co/api/v2/"

type Client struct {
	httpClient *http.Client
	baseURL    string
}
type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ExploreLocation struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
	} `json:"pokemon_encounters"`
}

type APIResponse[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type Config struct {
	Next     *string
	Previous *string
}

type PokemonData struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseUrl,
	}
}

func CreateUrl(c *Client, endpoints ...string) string {
	finalUrl := c.baseURL + strings.Join(endpoints, "/")
	return finalUrl
}

func NewRequest(method, url string, body io.Reader, c *Client) (*http.Request, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to communicate with PokeAPI at this time: %v", err)
	}

	return req, nil
}

func DoRequest(request *http.Request, c *Client) (*http.Response, error) {
	res, err := c.httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("cannot perform request: %v", err)
	}

	return res, nil
}

func GetResponse(response *http.Response) ([]byte, error) {
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to retreive body from response: %v", err)
	}

	if response.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and body: %s", response.StatusCode, body)
	}
	defer response.Body.Close()

	return body, nil

}

// Use for calls where the specific id or name is not given
func Unmarshal[T any](body []byte, config *Config) ([]T, error) {
	var APIResponse APIResponse[T]
	err := json.Unmarshal(body, &APIResponse)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json: %v", err)
	}

	config.Next = &APIResponse.Next
	config.Previous = &APIResponse.Previous

	return APIResponse.Results, nil
}

func UnmarshalExplore(body []byte, config *Config) ([]string, error) {
	var ExploreLocation ExploreLocation
	err := json.Unmarshal(body, &ExploreLocation)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json: %v", err)
	}

	pokemonNames := make([]string, 0)

	for _, pokemon := range ExploreLocation.PokemonEncounters {
		pokemonNames = append(pokemonNames, pokemon.Pokemon.Name)
	}
	return pokemonNames, nil
}

func UnmarshalPokemonData(body []byte, config *Config) (PokemonData, error) {
	var PokemonData PokemonData
	err := json.Unmarshal(body, &PokemonData)
	if err != nil {
		return PokemonData, fmt.Errorf("unable to parse json: %v", err)
	}
	return PokemonData, nil
}
