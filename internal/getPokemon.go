package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var url = "https://pokeapi.co/api/v2/pokemon/"

func GetPokemon(pokemonName string) (Pokemon, error) {
	var response Pokemon
	fullUrl := url + pokemonName

	res, err := http.Get(fullUrl)
	if err != nil {
		return response, fmt.Errorf("Could not get pokemon.\n Error: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, fmt.Errorf("Could not read body.\n Error: %v", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("Could not decode body.\n Error: %v", err)
	}

	return response, nil
}
