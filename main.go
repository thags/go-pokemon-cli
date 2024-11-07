package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var conf config

	for true {
		fmt.Print("pokedex > ")
		reader := bufio.NewScanner(os.Stdin)
		reader.Scan()
		text := reader.Text()
		command, exists := getCommand()[text]
		if exists {
			command.callback(&conf)
		} else {
			continue
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	currentMap pokeapiLocation
}

func commandHelp(*config) error {
	fmt.Println("Command line interface for a pokedex.")
	fmt.Println("\nCommands:")
	fmt.Println("help: to display this message")
	fmt.Println("exit: to exit the interface")
	fmt.Println("map: to display the first or next 20 locations")
	fmt.Println("mapb: to display the last 20 locations")
	return nil
}

func commandExit(*config) error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func commandMap(conf *config) error {
	fmt.Println("fetching data...")
	err := GetNextMap(conf)
	if err != nil {
		fmt.Println("error fetching map data")
		return err
	}

	for _, r := range conf.currentMap.Results {
		fmt.Println(r.Name)
	}

	return nil
}

func commandMapb(conf *config) error {
	fmt.Println("fetching data...")
	err := GetLastMap(conf)
	if err != nil {
		fmt.Println("error fetching previous map data")
		return err
	}

	for _, r := range conf.currentMap.Results {
		fmt.Println(r.Name)
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
	}
}
