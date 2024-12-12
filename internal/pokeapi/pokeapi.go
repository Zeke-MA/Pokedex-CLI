package pokeapi

import (
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

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseUrl,
	}
}

func CreateUrl(endpoints ...string) string {
	finalUrl := baseUrl + strings.Join(endpoints, "/")
	return finalUrl
}

func (c *Client) NewRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	url := CreateUrl(c.baseURL, endpoint)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to communicate with PokeAPI at this time: %v", err)
	}

	return req, nil
}

// create do request function

// create get response function

// create unmarshal function
