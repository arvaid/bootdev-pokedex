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
	words := strings.Split(text, " ")
	for _, word := range words {
		word = strings.Trim(word, " ")
		word = strings.ToLower(word)
		if word == "" || word == " " {
			continue
		}
		slice = append(slice, word)
		fmt.Println(word)
	}
	return slice
}
