package pokedex

import "github.com/Zeke-MA/pokedexcli/internal/pokeapi"

type Pokedex struct {
	CaughtPokemon map[string]pokeapi.PokemonData
}

func NewPokedex() *Pokedex {
	var Pokedex Pokedex
	Pokedex.CaughtPokemon = map[string]pokeapi.PokemonData{}
	return &Pokedex
}
