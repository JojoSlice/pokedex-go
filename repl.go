package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jojoslice/pokedexcli/internal"
)

var url = "https://pokeapi.co/api/v2/location-area/"

func cleanInput(input string) []string {
	return strings.Fields(strings.ToLower(input))
}

func commandExit(*config, []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	fmt.Println("Could not exit program")
	return nil
}

func commandHelp(*config, []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, _ []string) error {

	if cfg.next != nil {
		url = *cfg.next
	}

	res, err := internal.GetLocationAreas(url, cfg.cache)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	cfg.next = res.Next
	cfg.prev = res.Previous

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack(cfg *config, _ []string) error {
	if cfg.prev == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	res, err := internal.GetLocationAreas(*cfg.prev, cfg.cache)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	cfg.next = res.Next
	cfg.prev = res.Previous

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(_ *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: Explore <area name>")
		return nil
	}
	areaName := args[0]
	res, err := internal.GetLocationArea(url, areaName)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, encounter := range res.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}
