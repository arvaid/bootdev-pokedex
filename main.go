package main

import (
	"github.com/arvaid/pokedex/internal"
	"os"
)

func main() {
	cfg := config{
		Pokedex: map[string]internal.Pokemon{},
	}
	loop(&cfg)
	os.Exit(0)
}
