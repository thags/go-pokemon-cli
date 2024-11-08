package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type pokeapiLocation struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Results `json:"results"`
}
type Results struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

//Locations

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

type locationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters `json:"pokemon_encounters"`
}
type PokemonEncounters []struct {
	Pokemon struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokemon"`
	VersionDetails []struct {
		EncounterDetails []struct {
			Chance          int   `json:"chance"`
			ConditionValues []any `json:"condition_values"`
			MaxLevel        int   `json:"max_level"`
			Method          struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"method"`
			MinLevel int `json:"min_level"`
		} `json:"encounter_details"`
		MaxChance int `json:"max_chance"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"version_details"`
}

type location struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
	GameIndices []struct {
		GameIndex  int `json:"game_index"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"game_indices"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}
