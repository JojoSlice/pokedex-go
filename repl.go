package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jojoslice/pokedexcli/internal"
)

func cleanInput(input string) []string {
	return strings.Fields(strings.ToLower(input))
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	fmt.Println("Could not exit program")
	return nil
}

func commandHelp(*config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"

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

func commandMapBack(cfg *config) error {
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
