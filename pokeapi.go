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

func GetNextMap(conf *config) error {
	next := "https://pokeapi.co/api/v2/location/"
	if len(conf.currentMap.Next) > 0 {
		next = conf.currentMap.Next
	}

	res, err := http.Get(next)
	if err != nil {
		fmt.Println("error sending get request to location")
		fmt.Println("tried sending request to url: " + next)

		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body from location")
		return err
	}
	defer res.Body.Close()

	var locations pokeapiLocation
	err = json.Unmarshal(body, &locations)
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

	res, err := http.Get(conf.currentMap.Previous)
	if err != nil {
		fmt.Println("error sending get request to previous map")
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading body from location")
		return err
	}
	defer res.Body.Close()

	var locations pokeapiLocation
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("error unmarshalling json data to pokeapiLocation")
		return err
	}

	conf.currentMap = locations

	return nil
}
