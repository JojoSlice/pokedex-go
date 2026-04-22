package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

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

func commandCatch(_ *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: Catch <pokemon name>")
		return nil
	}

	pokemonName := args[0]

	res, err := internal.GetPokemon(pokemonName)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("Throwing a Pokeball at " + res.Name + "...")
	time.Sleep(time.Second)

	baseExp := res.BaseExperience
	chance := catchChance(baseExp)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Float64() < chance {
		fmt.Println(res.Name + " escaped!")
		return nil
	}

	fmt.Println(res.Name + " was cought!")
	coughtPokemon[res.Name] = res
	return nil
}

func catchChance(baseExp int) float64 {
	return 1.0 / (1.0 + float64(baseExp)/100.0)
}

func commandInspect(_ *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: Inspect <pokemon name>")
		return nil
	}

	pokemonName := args[0]

	p, ok := coughtPokemon[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Println("Name: " + p.Name)
	fmt.Println("Height: " + strconv.Itoa(p.Height))
	fmt.Println("Weight: " + strconv.Itoa(p.Weight))
	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Println("  -" + s.Stat.Name + ": " + strconv.Itoa(s.BaseStat))
	}
	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Println("  -" + t.Type.Name)
	}

	return nil
}
