package internal

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	API_BASE_URL         = "https://pokeapi.co/api/v2"
	API_AREA_ENDPOINT    = "/location-area"
	API_POKEMON_ENDPOINT = "/pokemon"
)

var cache *Cache

func GetMap(url string) ([]string, string, string, error) {
	if url == "" {
		url = API_BASE_URL + API_AREA_ENDPOINT
	}
	body, ok := getCachedResponseBody(url)
	if !ok {
		var err error
		body, err = getHTTPResponseBody(url)
		if err != nil {
			return nil, "", "", err
		}
	}
	data := AreaList{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, "", "", err
	}

	areas := make([]string, 0)
	for _, location := range data.Results {
		areas = append(areas, location.Name)
	}

	return areas, data.Next, data.Previous, nil
}

func Explore(area string) ([]string, error) {
	url := API_BASE_URL + API_AREA_ENDPOINT + "/" + area
	body, ok := getCachedResponseBody(url)
	if !ok {
		var err error
		body, err = getHTTPResponseBody(url)
		if err != nil {
			return nil, err
		}
	}
	data := Area{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	pokemons := make([]string, 0, len(data.PokemonEncounters))
	for _, encounter := range data.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon.Name)
	}

	return pokemons, nil
}

func Catch(pokemonName string) (Pokemon, bool) {
	data := Pokemon{}
	url := API_BASE_URL + API_POKEMON_ENDPOINT + "/" + pokemonName
	body, ok := getCachedResponseBody(url)
	if !ok {
		var err error
		body, err = getHTTPResponseBody(url)
		if err != nil {
			return data, false
		}
	}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return data, false
	}

	xp := data.BaseExperience
	caught := (rand.Intn(xp) / 10) == 0
	return data, caught
}

func getCachedResponseBody(url string) ([]byte, bool) {
	initCache()
	body, ok := cache.Get(url)
	return body, ok
}

func getHTTPResponseBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed, statuscode: %d, body: %s\n", res.StatusCode, res.Body)
	}
	initCache()
	cache.Add(url, body)
	return body, nil
}

func initCache() {
	if cache == nil {
		cache = NewCache(5 * time.Second)
	}
}
