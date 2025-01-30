package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		cleaned := cleanInput(input)
		fmt.Printf("Your command was: %s>\n", cleaned[0])
	}
	os.Exit(0)
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.Trim(text, " ")
	words := strings.Fields(text)
	return words
}
