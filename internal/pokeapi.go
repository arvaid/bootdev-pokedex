package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	API_BASE_URL = "https://pokeapi.co/api/v2"
	API_AREA_URL = "/location-area"
)

var cache *Cache

type MapResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type ExploreResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
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
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetMap(url string) ([]string, string, string, error) {
	if url == "" {
		url = API_BASE_URL + API_AREA_URL
	}
	body, ok := getCachedResponseBody(url)
	if !ok {
		var err error
		body, err = getHTTPResponseBody(url)
		if err != nil {
			return nil, "", "", err
		}
	}
	data := MapResponse{}
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
	url := API_BASE_URL + API_AREA_URL + "/" + area
	body, ok := getCachedResponseBody(url)
	if !ok {
		var err error
		body, err = getHTTPResponseBody(url)
		if err != nil {
			return nil, err
		}
	}
	data := ExploreResponse{}
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
