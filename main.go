package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jojoslice/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	next  *string
	prev  *string
	cache *internal.Cache
}

var commands map[string]cliCommand

var coughtPokemon map[string]internal.Pokemon

func init() {
	coughtPokemon = make(map[string]internal.Pokemon)
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "View the map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "View the previous map locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore the pokemon in the choosen area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View the stats of a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View the Pokedex",
			callback:    commandPokedex,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{cache: internal.NewCache(5 * time.Minute)}

	for {
		fmt.Print("Pokedex > ")
		for scanner.Scan() {
			input := scanner.Text()
			text := cleanInput(input)
			if len(text) == 0 {
				fmt.Println("No input")
				break
			}
			fmt.Print("\n")

			if cmd, ok := commands[text[0]]; ok {
				cmd.callback(cfg, text[1:])
			} else {
				fmt.Println("Unknown command")
			}
			break
		}
	}
}
