package main

import (
	"bufio"
	"fmt"
	"github.com/arvaid/pokedex/internal"
	"os"
	"strings"
)

func main() {
	commands = map[string]cliCommand{
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Throw Pokeball at pokemon",
			callback:    commandCatch,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "List Pokemon encounters in an area",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			callback:    commandMapB,
		},
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
	}

	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{}
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		input := cleanInput(line)
		command, args := input[0], input[1:]
		if cmd, ok := commands[command]; ok {
			err := cmd.callback(&cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
	os.Exit(0)
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string) error
}

var commands map[string]cliCommand

type config struct {
	Next     string
	Previous string
}

func commandCatch(_ *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing Pokemon argument")
	}
	pokemon := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	catch := internal.Catch(pokemon)
	if catch {
		fmt.Printf("%s was caught!\n", pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}

	return nil
}

func commandExplore(_ *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing location argument")
	}
	area := args[0]
	fmt.Printf("Exploring %s\n", area)

	pokemons, err := internal.Explore(area)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemons {
		fmt.Printf(" - %s\n", pokemon)
	}

	return nil
}

func commandMap(cfg *config, args ...string) error {
	areas, next, prev, err := internal.GetMap(cfg.Next)
	if err != nil {
		return err
	}
	cfg.Next = next
	cfg.Previous = prev
	for _, area := range areas {
		fmt.Println(area)
	}
	return nil
}

func commandMapB(cfg *config, args ...string) error {
	areas, next, prev, err := internal.GetMap(cfg.Previous)
	if err != nil {
		return err
	}
	cfg.Next = next
	cfg.Previous = prev
	for _, area := range areas {
		fmt.Println(area)
	}
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.Trim(text, " ")
	words := strings.Fields(text)
	return words
}
