package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Locations
func GetNextMap(conf *config) error {
	next := "https://pokeapi.co/api/v2/location/"
	if len(conf.currentMap.Next) > 0 {
		next = conf.currentMap.Next
	}

	body, exists := conf.cache.Get(next)
	if !exists {
		fmt.Println("did not find in cache, sending web request")
		res, err := http.Get(next)
		if err != nil {
			fmt.Println("error sending get request to location")
			fmt.Println("tried sending request to url: " + next)

			return err
		}

		body, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			fmt.Println("error reading body from location")
			return err
		}

		conf.cache.Add(next, body)
	}

	var locations pokeapiLocation
	err := json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("error unmarshalling json data to pokeapiLocation")
		return err
	}

	conf.currentMap = locations

	return nil
}

func GetLastMap(conf *config) error {
	if len(conf.currentMap.Previous) == 0 {
		fmt.Println("No previous available")
		return nil
	}

	body, exists := conf.cache.Get(conf.currentMap.Previous)
	if !exists {

		res, err := http.Get(conf.currentMap.Previous)
		if err != nil {
			fmt.Println("error sending get request to previous map")
			return err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("error reading body from location")
			return err
		}
		defer res.Body.Close()

		conf.cache.Add(conf.currentMap.Previous, body)
	}

	var locations pokeapiLocation
	err := json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("error unmarshalling json data to pokeapiLocation")
		return err
	}

	conf.currentMap = locations

	return nil
}

// encounters
func GetEncounters(conf *config, name string) error {
	url, err := GetAreaUrl(conf, name)
	if err != nil {
		return err
	}

	body, exists := conf.cache.Get(url)
	if exists {
		fmt.Println("area found in cache!")
	}

	if !exists {
		fmt.Println("area not found in cache")
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("error sending get request to get encounter info")
			return err
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("error reading encounter info")
			return err
		}
		defer res.Body.Close()
		conf.cache.Add(url, body)
	}

	var area locationArea
	err = json.Unmarshal(body, &area)
	if err != nil {
		fmt.Println("error unmarshalling json data to encounter")
		return err
	}

	conf.currentEncounter = area.PokemonEncounters

	return nil
}

func GetAreaUrl(conf *config, name string) (string, error) {
	url := "https://pokeapi.co/api/v2/location/" + name
	if len(name) == 0 {
		return "", fmt.Errorf("Area name can not be blank")
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error sending get request to get area info")
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading area info")
		return "", err
	}
	defer res.Body.Close()
	var area location
	err = json.Unmarshal(body, &area)
	if err != nil {
		fmt.Println("error unmarshalling json data to location")
		return "", err
	}
	return area.Areas[0].URL, nil
}

// pokemon

func GetPokemon(conf *config, name string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	if len(name) == 0 {
		return Pokemon{}, fmt.Errorf("pokemon can not be empty")
	}
	body, exists := conf.cache.Get(url)
	if !exists {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("error sending get request to get pokemon info")
			return Pokemon{}, err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("error reading area info")
			return Pokemon{}, err
		}
		conf.cache.Add(url, body)
		defer res.Body.Close()
	}
	var pokemon Pokemon
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		fmt.Println("error unmarshalling json data to pokemon")
		return Pokemon{}, err
	}

	conf.pokemon[pokemon.Name] = pokemon
	return pokemon, nil
}
