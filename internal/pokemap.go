package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	API_BASE_URL = "https://pokeapi.co/api/v2"
	API_MAP_URL  = "/location-area"
)

type MapResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func getMapData(url string) (MapResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return MapResponse{}, err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed, statuscode: %d, body: %s\n", res.StatusCode, res.Body)
	}
	if err != nil {
		return MapResponse{}, err
	}
	resp := MapResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return MapResponse{}, err
	}

	return resp, nil
}

func GetNextMap(url string) (string, string, error) {
	if url == "" {
		url = API_BASE_URL + API_MAP_URL
	}
	mapData, err := getMapData(url)
	if err != nil {
		return "", "", err
	}

	for _, location := range mapData.Results {
		fmt.Println(location.Name)
	}

	return mapData.Next, mapData.Previous, nil
}
