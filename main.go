package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	slice := make([]string, 0)
	text = strings.Trim(text, " ")
	words := strings.Fields(text)
	for _, word := range words {
		word = strings.ToLower(word)
		slice = append(slice, word)
	}
	return slice
}
