package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	pokeCache "pokedexcli/internal"
	"strings"
	"time"
)

func main() {
	var conf config
	conf.pokemon = make(map[string]Pokemon)
	reapingTime := 10 * time.Minute
	conf.cache = *pokeCache.NewCache(reapingTime)

	for true {
		fmt.Print("pokedex > ")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		text := reader.Text()
		inputCommand := strings.Split(text, " ")
		command, exists := getCommand()[inputCommand[0]]

		input := ""
		if len(inputCommand) > 1 {
			input = inputCommand[1]
		}

		if exists {
			command.callback(&conf, input)
		} else {
			continue
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	currentMap       pokeapiLocation
	cache            pokeCache.PokeCache
	currentEncounter PokemonEncounters
	pokemon          map[string]Pokemon
}

func commandHelp(*config, string) error {
	fmt.Println("Command line interface for a pokedex.")
	fmt.Println("\nCommands:")
	fmt.Println("help: to display this message")
	fmt.Println("exit: to exit the interface")
	fmt.Println("map: to display the first or next 20 locations")
	fmt.Println("mapb: to display the last 20 locations")
	fmt.Println("explore <location name>: find pokemon in the area")
	fmt.Println("catch <pokemon name>: try to catch a pokemon")

	return nil
}

func commandExit(*config, string) error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func commandMap(conf *config, input string) error {
	fmt.Println("fetching data...")
	err := GetNextMap(conf)
	if err != nil {
		fmt.Println("error fetching map data")
		return err
	}

	printMap(conf)

	return nil
}

func commandMapb(conf *config, input string) error {
	fmt.Println("fetching data...")
	err := GetLastMap(conf)
	if err != nil {
		fmt.Println("error fetching previous map data")
		return err
	}

	printMap(conf)

	return nil
}

func printMap(conf *config) {
	fmt.Println("You explored and found this locations:")
	for _, r := range conf.currentMap.Results {
		fmt.Printf("    - %v\n", r.Name)
	}
}

func commandExplore(conf *config, name string) error {
	err := GetEncounters(conf, name)
	if err != nil {
		return err
	}

	for _, encounter := range conf.currentEncounter {
		fmt.Printf("    -%v\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(conf *config, name string) error {
	pokemon, err := GetPokemon(conf, name)
	if err != nil {
		fmt.Println("You did not catch the pokemon, as it does not exist")
		return nil
	}

	random := rand.Int() * rand.Int()
	needednumber := pokemon.BaseExperience

	if random >= needednumber {
		conf.pokemon[pokemon.Name] = pokemon
		fmt.Printf("You caught %v!\n", pokemon.Name)
	}
	return nil
}

func getCommand() map[string]cliCommand {

	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "See the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "See the last 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area and find pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon, good luck!",
			callback:    commandCatch,
		},
	}
}
