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
		command := input[0]
		if cmd, ok := commands[command]; ok {
			cmd.callback(&cfg)
		} else {
			fmt.Println("Unknown command")
		}
	}
	os.Exit(0)
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
}

var commands map[string]cliCommand

type config struct {
	Next     string
	Previous string
}

func commandMap(cfg *config) error {
	next, prev, err := internal.GetNextMap(cfg.Next)
	if err != nil {
		return err
	}
	cfg.Next = next
	cfg.Previous = prev
	return nil
}

func commandMapB(cfg *config) error {
	next, prev, err := internal.GetNextMap(cfg.Previous)
	if err != nil {
		return err
	}
	cfg.Next = next
	cfg.Previous = prev
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for cmdName, cmd := range commands {
		fmt.Printf("%s: %s\n", cmdName, cmd.description)
	}
	return nil
}

func commandExit(cfg *config) error {
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
