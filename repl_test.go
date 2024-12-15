package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/Zeke-MA/pokedexcli/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Unexpected length. Expected: %v Actual: %v for input: %v ", len(c.expected), len(actual), c.input)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Unexpected word found. Expected: %v Actual: %v for input: %v ", c.expected, actual, c.input)
			}
		}
	}
}

// Modify test cases here
func TestAddGetCache(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://pokeapi.co/api/v2/location-area?offset=40&limit=20",
			val: []byte(`{"count":1054,"next":"https://pokeapi.co/api/v2/location-area?offset=60&limit=20","previous":"https://pokeapi.co/api/v2/location-area?offset=20&limit=20","results":[{"name":"solaceon-ruins-b3f-d","url":"https://pokeapi.co/api/v2/location-area/41/"},{"name":"solaceon-ruins-b3f-e","url":"https://pokeapi.co/api/v2/location-area/42/"},{"name":"solaceon-ruins-b4f-a","url":"https://pokeapi.co/api/v2/location-area/43/"},{"name":"solaceon-ruins-b4f-b","url":"https://pokeapi.co/api/v2/location-area/44/"},{"name":"solaceon-ruins-b4f-c","url":"https://pokeapi.co/api/v2/location-area/45/"},{"name":"solaceon-ruins-b4f-d","url":"https://pokeapi.co/api/v2/location-area/46/"},{"name":"solaceon-ruins-b5f","url":"https://pokeapi.co/api/v2/location-area/47/"},{"name":"sinnoh-victory-road-1f","url":"https://pokeapi.co/api/v2/location-area/48/"},{"name":"sinnoh-victory-road-2f","url":"https://pokeapi.co/api/v2/location-area/49/"},{"name":"sinnoh-victory-road-b1f","url":"https://pokeapi.co/api/v2/location-area/50/"},{"name":"sinnoh-victory-road-inside-b1f","url":"https://pokeapi.co/api/v2/location-area/51/"},{"name":"sinnoh-victory-road-inside","url":"https://pokeapi.co/api/v2/location-area/52/"},{"name":"sinnoh-victory-road-inside-exit","url":"https://pokeapi.co/api/v2/location-area/53/"},{"name":"ravaged-path-area","url":"https://pokeapi.co/api/v2/location-area/54/"},{"name":"oreburgh-gate-1f","url":"https://pokeapi.co/api/v2/location-area/55/"},{"name":"oreburgh-gate-b1f","url":"https://pokeapi.co/api/v2/location-area/56/"},{"name":"stark-mountain-area","url":"https://pokeapi.co/api/v2/location-area/57/"},{"name":"stark-mountain-entrance","url":"https://pokeapi.co/api/v2/location-area/58/"},{"name":"stark-mountain-inside","url":"https://pokeapi.co/api/v2/location-area/59/"},{"name":"sendoff-spring-area","url":"https://pokeapi.co/api/v2/location-area/60/"}]}`),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
